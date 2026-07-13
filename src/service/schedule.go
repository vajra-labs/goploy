package service

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	"goploy/src/conf"
	"goploy/src/core/errorx"
	"goploy/src/db/repos"
	"goploy/src/utility/docker"
	"goploy/src/utility/shell"
)

type ScheduleService struct {
	cfg     *conf.Config
	queries *repos.Queries
	paths   *docker.AppPaths
}

func NewScheduleService(cfg *conf.Config, queries *repos.Queries, paths *docker.AppPaths) *ScheduleService {
	return &ScheduleService{
		cfg:     cfg,
		queries: queries,
		paths:   paths,
	}
}

// FindScheduleByID finds a schedule by ID with relations loaded.
func (s *ScheduleService) FindScheduleByID(ctx context.Context, id int64) (*repos.GetScheduleWithRelationsRow, error) {
	schedule, err := s.queries.GetScheduleWithRelations(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorx.NotFoundError("Schedule not found", "SCHEDULE_NOT_FOUND")
		}
		return nil, errorx.InternalServerError("Failed to find schedule", "SCHEDULE_FIND_ERROR", errorx.WithCause(err))
	}
	return &schedule, nil
}

// ListEnabledSchedules lists all enabled schedules.
func (s *ScheduleService) ListEnabledSchedules(ctx context.Context) ([]repos.ListEnabledSchedulesRow, error) {
	schedulesList, err := s.queries.ListEnabledSchedules(ctx)
	if err != nil {
		return nil, errorx.InternalServerError("Failed to list enabled schedules", "SCHEDULE_LIST_ERROR", errorx.WithCause(err))
	}
	return schedulesList, nil
}

// CreateSchedule inserts a new schedule and writes scripts if server type.
func (s *ScheduleService) CreateSchedule(ctx context.Context, arg repos.CreateScheduleParams) (*repos.Schedule, error) {
	newSchedule, err := s.queries.CreateSchedule(ctx, arg)
	if err != nil {
		return nil, errorx.InternalServerError("Failed to create schedule", "SCHEDULE_CREATE_ERROR", errorx.WithCause(err))
	}

	if newSchedule.ScheduleType == "Goploy-SERVER" || newSchedule.ScheduleType == "SERVER" {
		if err := s.handleScript(ctx, newSchedule); err != nil {
			return nil, err
		}
	}

	return &newSchedule, nil
}

// UpdateSchedule updates a schedule and updates script file if server type.
func (s *ScheduleService) UpdateSchedule(ctx context.Context, arg repos.UpdateScheduleParams) (*repos.Schedule, error) {
	updatedSchedule, err := s.queries.UpdateSchedule(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorx.NotFoundError("Schedule not found", "SCHEDULE_NOT_FOUND")
		}
		return nil, errorx.InternalServerError("Failed to update schedule", "SCHEDULE_UPDATE_ERROR", errorx.WithCause(err))
	}

	if updatedSchedule.ScheduleType == "Goploy-SERVER" || updatedSchedule.ScheduleType == "SERVER" {
		if err := s.handleScript(ctx, updatedSchedule); err != nil {
			return nil, err
		}
	}

	return &updatedSchedule, nil
}

// DeleteSchedule deletes the schedule and cleans up its directory.
func (s *ScheduleService) DeleteSchedule(ctx context.Context, id int64) error {
	schedule, err := s.queries.GetScheduleWithRelations(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errorx.NotFoundError("Schedule not found", "SCHEDULE_NOT_FOUND")
		}
		return errorx.InternalServerError("Failed to find schedule", "SCHEDULE_FIND_ERROR", errorx.WithCause(err))
	}

	// Identify target server ID (order of priority: schedule.ServerID, application.ServerID, compose.ServerID)
	var targetServerID *int64
	if schedule.ServerID != nil {
		targetServerID = schedule.ServerID
	} else if schedule.ApplicationServerID != nil {
		targetServerID = schedule.ApplicationServerID
	} else if schedule.ComposeServerID != nil {
		targetServerID = schedule.ComposeServerID
	}

	folderPath := filepath.Join(s.paths.SCHEDULES_PATH, schedule.AppName)
	cleanupCmd := fmt.Sprintf("rm -rf %s", folderPath)

	if targetServerID != nil {
		// Remote cleanup
		client, err := s.getOrRegisterSSHClient(ctx, *targetServerID)
		if err != nil {
			return err
		}
		resChan := client.Exec(ctx, cleanupCmd, nil)
		res := <-resChan
		if res.Err != nil {
			return errorx.InternalServerError("Failed to remove remote script directory", "REMOTE_SCRIPT_CLEANUP_ERROR", errorx.WithCause(res.Err))
		}
	} else {
		// Local cleanup
		resChan := shell.Exec(ctx, cleanupCmd, shell.WithShell("/bin/bash"))
		res := <-resChan
		if res.Err != nil {
			return errorx.InternalServerError("Failed to remove local script directory", "LOCAL_SCRIPT_CLEANUP_ERROR", errorx.WithCause(res.Err))
		}
	}

	err = s.queries.DeleteSchedule(ctx, id)
	if err != nil {
		return errorx.InternalServerError("Failed to delete schedule row", "SCHEDULE_DELETE_ERROR", errorx.WithCause(err))
	}

	return nil
}

// handleScript deploys the script shell file either locally or remotely.
func (s *ScheduleService) handleScript(ctx context.Context, schedule repos.Schedule) error {
	folderPath := filepath.Join(s.paths.SCHEDULES_PATH, schedule.AppName)

	scriptVal := ""
	if schedule.Script != nil {
		scriptVal = *schedule.Script
	}

	// Suffix PID and ID to script output just like Dokploy
	scriptWithPid := fmt.Sprintf("echo \"PID: $$ | Schedule ID: %d\"\n%s", schedule.ID, scriptVal)
	encodedContent := base64.StdEncoding.EncodeToString([]byte(scriptWithPid))

	deployCmd := fmt.Sprintf(`
mkdir -p %s
rm -f %s/script.sh
touch %s/script.sh
chmod +x %s/script.sh
echo "%s" | base64 -d > %s/script.sh
`, folderPath, folderPath, folderPath, folderPath, encodedContent, folderPath)

	shellPath := "/bin/bash"
	if schedule.ShellType == "SH" {
		shellPath = "/bin/sh"
	}

	if schedule.ScheduleType == "Goploy-SERVER" {
		// Run deployment script locally
		resChan := shell.Exec(ctx, deployCmd, shell.WithShell(shellPath))
		res := <-resChan
		if res.Err != nil {
			return errorx.InternalServerError("Failed to deploy script locally", "LOCAL_SCRIPT_DEPLOY_ERROR", errorx.WithCause(res.Err))
		}
	} else if schedule.ScheduleType == "SERVER" {
		if schedule.ServerID == nil {
			return errorx.BadRequestError("Server ID is required for SERVER schedule type", "SERVER_ID_REQUIRED")
		}

		// Run deployment script remotely via SSH client
		client, err := s.getOrRegisterSSHClient(ctx, *schedule.ServerID)
		if err != nil {
			return err
		}

		resChan := client.Exec(ctx, deployCmd, nil)
		res := <-resChan
		if res.Err != nil {
			return errorx.InternalServerError("Failed to deploy script remotely", "REMOTE_SCRIPT_DEPLOY_ERROR", errorx.WithCause(res.Err))
		}
	}

	return nil
}

// getOrRegisterSSHClient resolves or registers an SSH client for remote operations.
func (s *ScheduleService) getOrRegisterSSHClient(ctx context.Context, serverID int64) (*shell.SSHClient, error) {
	serverIDStr := strconv.FormatInt(serverID, 10)
	client, err := shell.GetSSHClient(serverIDStr)
	if err == nil {
		return client, nil
	}

	// Client not in pool, register it now
	serverData, err := s.queries.GetServerWithSSHKey(ctx, serverID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorx.NotFoundError("Target server not found", "SERVER_NOT_FOUND")
		}
		return nil, errorx.InternalServerError("Failed to fetch target server credentials", "SERVER_FETCH_ERROR", errorx.WithCause(err))
	}

	if serverData.PrivateKey == nil || *serverData.PrivateKey == "" {
		return nil, errorx.BadRequestError("Server SSH private key is missing", "SSH_KEY_MISSING")
	}

	sshCfg := shell.SSHConfig{
		Host:       serverData.IpAddress,
		Port:       strconv.FormatInt(serverData.Port, 10),
		User:       serverData.Username,
		PrivateKey: *serverData.PrivateKey,
	}

	err = shell.SetSSHClient(serverIDStr, sshCfg)
	if err != nil {
		return nil, errorx.InternalServerError("Failed to register and dial remote SSH connection", "SSH_DIAL_ERROR", errorx.WithCause(err))
	}

	client, err = shell.GetSSHClient(serverIDStr)
	if err != nil {
		return nil, errorx.InternalServerError("Failed to fetch registered SSH connection from pool", "SSH_POOL_ERROR", errorx.WithCause(err))
	}

	return client, nil
}
