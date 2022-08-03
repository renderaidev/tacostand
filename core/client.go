/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package core

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/SomusHQ/tacostand/config"
	"github.com/SomusHQ/tacostand/db"
	"github.com/SomusHQ/tacostand/logger"
	"github.com/SomusHQ/tacostand/tasks"
	"github.com/go-co-op/gocron"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// Client is the main application object.
type Client struct {
	Logger            *logger.Logger
	Config            *config.Config
	SlackAPI          *slack.Client
	SlackSocketClient *socketmode.Client
	Scheduler         *gocron.Scheduler
	Database          *db.DB
}

// NewClient creates a new application client and instantiates all the required
// services (database, cron scheduler, Slack API and socket clients). Returns an
// error on invalid configuration or invalid database credentials.
func NewClient() (*Client, error) {
	config, err := config.Init()
	if err != nil {
		return nil, err
	}

	logger := logger.NewLogger(log.New(os.Stdout, "tacostand: ", log.LstdFlags), config.DebugMode)

	api := slack.New(
		config.Slack.BotToken,
		slack.OptionDebug(config.DebugMode),
		slack.OptionLog(log.New(os.Stdout, "slack api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(config.Slack.AppToken),
	)

	socketClient := socketmode.New(
		api,
		socketmode.OptionLog(log.New(os.Stdout, "slack socket client: ", log.Lshortfile|log.LstdFlags)),
	)

	scheduler := gocron.NewScheduler(time.UTC)

	client := &Client{
		Logger:            logger,
		Config:            config,
		SlackAPI:          api,
		SlackSocketClient: socketClient,
		Scheduler:         scheduler,
	}

	database, err := db.New(config, logger)
	if err != nil {
		return nil, err
	}

	client.Database = database

	return client, nil
}

// Run will start the goroutine for handling Slack socket events and will
// register the tasks for the cron scheduler. Both of these will run in an
// asynchronous fashion.
func (c *Client) Run(ctx context.Context) error {
	go c.HandleEvents(ctx)

	c.Scheduler.Cron(c.Config.CronExpression).Do(tasks.HandlePeriodicStandup(ctx))
	c.Scheduler.StartAsync()

	return c.SlackSocketClient.Run()
}
