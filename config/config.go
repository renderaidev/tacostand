/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/SomusHQ/tacostand/common"
	dotenv "github.com/joho/godotenv"
)

// Init will load the configuration from the dotenv file and the environment
// variables and return a configuration object.
func Init() (*Config, error) {
	dotenv.Load()

	slackAppToken := os.Getenv("SLACK_APP_TOKEN")
	if !strings.HasPrefix(slackAppToken, "xapp-") {
		return nil, common.ErrInvalidSlackAppToken.Error()
	}

	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	if !strings.HasPrefix(slackBotToken, "xoxb-") {
		return nil, common.ErrInvalidSlackBotToken.Error()
	}

	slackReportChannelName := os.Getenv("SLACK_REPORT_CHANNEL_NAME")
	if slackReportChannelName == "" {
		return nil, common.ErrSlackReportChannelNameNotSet.Error()
	}

	cronExpression := os.Getenv("CRON_EXPRESSION")
	if cronExpression == "" {
		return nil, common.ErrInvalidCronExpression.Error()
	}

	wrapUpTime, err := strconv.Atoi(os.Getenv("WRAP_UP_TIME"))
	if err != nil {
		wrapUpTime = 30
	}

	databaseHost := os.Getenv("PGHOST")
	databasePort, err := strconv.Atoi(os.Getenv("PGPORT"))
	if err != nil {
		databasePort = 5432
	}

	databaseUser := os.Getenv("PGUSER")
	databasePassword := os.Getenv("PGPASSWORD")
	databaseName := os.Getenv("PGDATABASE")

	debugModeRaw := os.Getenv("DEBUG_MODE")
	debugMode := debugModeRaw == "true"

	skipMigrationsRaw := os.Getenv("SKIP_MIGRATIONS")
	skipMigrations := skipMigrationsRaw == "true"

	return &Config{
		Slack: &SlackConfig{
			AppToken:          slackAppToken,
			BotToken:          slackBotToken,
			ReportChannelName: slackReportChannelName,
		},

		Database: &DatabaseConfig{
			Host:     databaseHost,
			Port:     databasePort,
			User:     databaseUser,
			Password: databasePassword,
			Name:     databaseName,
		},

		CronExpression: cronExpression,
		WrapUpTime:     uint32(wrapUpTime),

		DebugMode:      debugMode,
		SkipMigrations: skipMigrations,
	}, nil
}
