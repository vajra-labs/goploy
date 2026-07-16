package handler

import (
	"runtime"
	"time"

	"goploy/src/apis/dtos"
	"goploy/src/conf"

	"github.com/dustin/go-humanize"
	"github.com/gofiber/fiber/v3"
)

// HealthHandler manages system diagnostics and server health endpoints.
type HealthHandler struct {
	cfg *conf.Config
}

// NewHealthHandler constructs a HealthHandler dependency.
func NewHealthHandler(cfg *conf.Config) *HealthHandler {
	return &HealthHandler{cfg: cfg}
}

// Ping returns a simple pong message to test server availability.
func (h *HealthHandler) Ping(ctx fiber.Ctx) error {
	return ctx.SendString("Pong!")
}

// Pong returns a simple ping message.
func (h *HealthHandler) Pong(ctx fiber.Ctx) error {
	return ctx.SendString("Ping!")
}

// Health gathers runtime memory metrics and returns diagnostic stats.
func (h *HealthHandler) Health(ctx fiber.Ctx) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	uptime := humanize.Time(h.cfg.START_TIME)
	result := dtos.HealthRes{
		Uptime:    uptime,
		Version:   conf.VERSION,
		GoEnv:     h.cfg.GO_ENV,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Memory: dtos.MemoryUsage{
			Alloc: dtos.MemoryInfo{
				Bytes: m.Alloc,
				Human: humanize.Bytes(m.Alloc),
			},
			TotalAlloc: dtos.MemoryInfo{
				Bytes: m.TotalAlloc,
				Human: humanize.Bytes(m.TotalAlloc),
			},
			Sys: dtos.MemoryInfo{
				Bytes: m.Sys,
				Human: humanize.Bytes(m.Sys),
			},
			HeapAlloc: dtos.MemoryInfo{
				Bytes: m.HeapAlloc,
				Human: humanize.Bytes(m.HeapAlloc),
			},
			HeapSys: dtos.MemoryInfo{
				Bytes: m.HeapSys,
				Human: humanize.Bytes(m.HeapSys),
			},
		},
	}
	return ctx.JSON(result)
}
