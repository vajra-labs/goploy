-- +goose Up
-- create "groups" table
CREATE TABLE `groups` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NOT NULL,
  `created_at` integer NOT NULL DEFAULT (strftime('%s', 'now')),
  `updated_at` integer NOT NULL DEFAULT (strftime('%s', 'now'))
) STRICT;
-- create index "groups_name" to table: "groups"
CREATE UNIQUE INDEX `groups_name` ON `groups` (`name`);
-- create trigger "groups_updated_at"
-- +goose StatementBegin
CREATE TRIGGER `groups_updated_at` AFTER UPDATE ON `groups` FOR EACH ROW BEGIN
	UPDATE groups
	SET updated_at = strftime('%s', 'now')
	WHERE id = OLD.id;
END;
-- +goose StatementEnd
-- create "policy" table
CREATE TABLE `policy` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `action` text NOT NULL,
  `created_at` integer NOT NULL DEFAULT (strftime('%s', 'now'))
) STRICT;
-- create index "policy_action" to table: "policy"
CREATE UNIQUE INDEX `policy_action` ON `policy` (`action`);
-- create "group_policy" table
CREATE TABLE `group_policy` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `group_id` integer NOT NULL,
  `policy_id` integer NOT NULL,
  `created_at` integer NOT NULL DEFAULT (strftime('%s', 'now')),
  CONSTRAINT `0` FOREIGN KEY (`policy_id`) REFERENCES `policy` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `1` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) STRICT;
-- create "users" table
CREATE TABLE `users` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `email` text NULL,
  `last_name` text NULL,
  `first_name` text NULL,
  `avatar` text NOT NULL,
  `role` text NULL DEFAULT 'OWNER',
  `about_me` text NULL,
  `password` text NOT NULL,
  `is_email_verify` integer NULL DEFAULT 0,
  `email_verify_at` integer NULL,
  `two_factor_enable` integer NULL DEFAULT 0,
  `is_registered` integer NOT NULL DEFAULT 0,
  `added_by` integer NULL DEFAULT (NULL),
  `group_id` integer NOT NULL,
  `created_at` integer NOT NULL DEFAULT (strftime('%s', 'now')),
  `updated_at` integer NOT NULL DEFAULT (strftime('%s', 'now')),
  CONSTRAINT `0` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `1` FOREIGN KEY (`added_by`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `role_check` CHECK (role IN ('OWNER', 'ADMIN', 'MEMBER'))
) STRICT;
-- create index "users_email" to table: "users"
CREATE UNIQUE INDEX `users_email` ON `users` (`email`);
-- create trigger "users_updated_at"
-- +goose StatementBegin
CREATE TRIGGER `users_updated_at` AFTER UPDATE ON `users` FOR EACH ROW BEGIN
	UPDATE users
	SET updated_at = strftime('%s', 'now')
	WHERE id = OLD.id;
END;
-- +goose StatementEnd
-- create "two_factor" table
CREATE TABLE `two_factor` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `secret` text NOT NULL,
  `backup_codes` text NOT NULL,
  `user_id` integer NOT NULL,
  CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) STRICT;
-- create "jwt_tokens" table
CREATE TABLE `jwt_tokens` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `jti` text NOT NULL,
  `role` text NOT NULL,
  `user_id` integer NOT NULL,
  `is_blacklist` integer NULL DEFAULT 0,
  `blacklist_at` integer NULL,
  `expired_at` integer NULL,
  `created_at` integer NOT NULL DEFAULT (strftime('%s', 'now')),
  `updated_at` integer NOT NULL DEFAULT (strftime('%s', 'now')),
  CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `role_check` CHECK (role IN ('OWNER', 'ADMIN', 'MEMBER'))
) STRICT;
-- create trigger "jwt_tokens_updated_at"
-- +goose StatementBegin
CREATE TRIGGER `jwt_tokens_updated_at` AFTER UPDATE ON `jwt_tokens` FOR EACH ROW BEGIN
	UPDATE jwt_tokens
	SET updated_at = strftime('%s', 'now')
	WHERE id = OLD.id;
END;
-- +goose StatementEnd
-- create "activity_logs" table
CREATE TABLE `activity_logs` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `user_id` integer NOT NULL,
  `activity` text NOT NULL,
  `source` text NOT NULL,
  `client_ip` text NOT NULL,
  `created_at` integer NOT NULL DEFAULT (strftime('%s', 'now')),
  CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `activity_check` CHECK (activity IN ('LOGIN', 'LOGOUT', 'REGISTER'))
) STRICT;
-- create "server_metrics" table
CREATE TABLE `server_metrics` (
  `timestamp` integer NULL,
  `cpu` real NOT NULL,
  `cpu_model` text NOT NULL,
  `cpu_cores` integer NOT NULL,
  `cpu_physical_cores` integer NOT NULL,
  `cpu_speed` real NOT NULL,
  `os` text NOT NULL,
  `distro` text NOT NULL,
  `kernel` text NOT NULL,
  `arch` text NOT NULL,
  `mem_used` real NOT NULL,
  `mem_used_gb` real NOT NULL,
  `mem_total` real NOT NULL,
  `uptime` integer NOT NULL,
  `disk_used` real NOT NULL,
  `total_disk` real NOT NULL,
  `network_in` real NOT NULL,
  `network_out` real NOT NULL,
  PRIMARY KEY (`timestamp`)
) STRICT;
-- create "container_metrics" table
CREATE TABLE `container_metrics` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `timestamp` integer NOT NULL,
  `container_id` text NOT NULL,
  `container_name` text NOT NULL,
  `metrics_json` text NOT NULL
) STRICT;
-- create index "idx_container_metrics_timestamp" to table: "container_metrics"
CREATE INDEX `idx_container_metrics_timestamp` ON `container_metrics` (`timestamp`);
-- create index "idx_container_metrics_name" to table: "container_metrics"
CREATE INDEX `idx_container_metrics_name` ON `container_metrics` (`container_name`);

-- +goose Down
-- reverse: create index "idx_container_metrics_name" to table: "container_metrics"
DROP INDEX `idx_container_metrics_name`;
-- reverse: create index "idx_container_metrics_timestamp" to table: "container_metrics"
DROP INDEX `idx_container_metrics_timestamp`;
-- reverse: create "container_metrics" table
DROP TABLE `container_metrics`;
-- reverse: create "server_metrics" table
DROP TABLE `server_metrics`;
-- reverse: create "activity_logs" table
DROP TABLE `activity_logs`;
-- reverse: create trigger "jwt_tokens_updated_at"
DROP TRIGGER `jwt_tokens_updated_at`;
-- reverse: create "jwt_tokens" table
DROP TABLE `jwt_tokens`;
-- reverse: create "two_factor" table
DROP TABLE `two_factor`;
-- reverse: create trigger "users_updated_at"
DROP TRIGGER `users_updated_at`;
-- reverse: create index "users_email" to table: "users"
DROP INDEX `users_email`;
-- reverse: create "users" table
DROP TABLE `users`;
-- reverse: create "group_policy" table
DROP TABLE `group_policy`;
-- reverse: create index "policy_action" to table: "policy"
DROP INDEX `policy_action`;
-- reverse: create "policy" table
DROP TABLE `policy`;
-- reverse: create trigger "groups_updated_at"
DROP TRIGGER `groups_updated_at`;
-- reverse: create index "groups_name" to table: "groups"
DROP INDEX `groups_name`;
-- reverse: create "groups" table
DROP TABLE `groups`;
