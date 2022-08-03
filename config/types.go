/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package config

// SlackConfig contains the app and bot tokens for the Slack integration, as
// well as the name of the channel where the reports will be collected.
type SlackConfig struct {
	AppToken string
	BotToken string

	ReportChannelName string
}

// DatabaseConfig contains the configuration for a Postgres database.
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// Config wraps around the Slack and database configurations and contains
// a cron expression for running the stand-up notifications periodically.
type Config struct {
	Slack          *SlackConfig
	Database       *DatabaseConfig
	CronExpression string

	DebugMode bool
}
