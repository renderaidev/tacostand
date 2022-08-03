/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package contextutils

import (
	"context"
	"fmt"

	"github.com/SomusHQ/tacostand/config"
	"github.com/SomusHQ/tacostand/db"
	"github.com/SomusHQ/tacostand/logger"
	"github.com/go-co-op/gocron"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type contextKey struct {
	Name string
}

func (ck *contextKey) String() string {
	return ck.Name
}

var (
	ctxConfig            = &contextKey{"config"}
	ctxLogger            = &contextKey{"logger"}
	ctxSlackApi          = &contextKey{"slackApi"}
	ctxSlackSocketClient = &contextKey{"slackSocketClient"}
	ctxScheduler         = &contextKey{"scheduler"}
	ctxDatabase          = &contextKey{"database"}
)

// WithConfig returns a Go context with a configuration object.
func WithConfig(ctx context.Context, config *config.Config) context.Context {
	return context.WithValue(ctx, ctxConfig, config)
}

// Config attempts to extract the configuration object from a Go context.
func Config(ctx context.Context) (*config.Config, error) {
	config, ok := ctx.Value(ctxConfig).(*config.Config)
	if !ok {
		return nil, fmt.Errorf("config not found in context")
	}

	return config, nil
}

// WithConfig returns a Go context with a logger.
func WithLogger(ctx context.Context, logger *logger.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger, logger)
}

// Config attempts to extract the logger from a Go context.
func Logger(ctx context.Context) (*logger.Logger, error) {
	logger, ok := ctx.Value(ctxLogger).(*logger.Logger)
	if !ok {
		return nil, fmt.Errorf("missing context: %s", ctxLogger)
	}

	return logger, nil
}

// WithConfig returns a Go context with a Slack API client.
func WithSlackApi(ctx context.Context, api *slack.Client) context.Context {
	return context.WithValue(ctx, ctxSlackApi, api)
}

// Config attempts to extract the Slack API client from a Go context.
func SlackApi(ctx context.Context) (*slack.Client, error) {
	api, ok := ctx.Value(ctxSlackApi).(*slack.Client)
	if !ok {
		return nil, fmt.Errorf("missing context: %s", ctxSlackApi)
	}

	return api, nil
}

// WithConfig returns a Go context with a Slack socket client.
func WithSlackSocketClient(ctx context.Context, client *socketmode.Client) context.Context {
	return context.WithValue(ctx, ctxSlackSocketClient, client)
}

// Config attempts to extract the Slack socket client from a Go context.
func SlackSocketClient(ctx context.Context) (*socketmode.Client, error) {
	client, ok := ctx.Value(ctxSlackSocketClient).(*socketmode.Client)
	if !ok {
		return nil, fmt.Errorf("missing context: %s", ctxSlackSocketClient)
	}

	return client, nil
}

// WithConfig returns a Go context with a cron scheduler instance.
func WithScheduler(ctx context.Context, scheduler *gocron.Scheduler) context.Context {
	return context.WithValue(ctx, ctxScheduler, scheduler)
}

// Config attempts to extract the cron scheduler instance from a Go context.
func Scheduler(ctx context.Context) (*gocron.Scheduler, error) {
	scheduler, ok := ctx.Value(ctxScheduler).(*gocron.Scheduler)
	if !ok {
		return nil, fmt.Errorf("missing context: %s", ctxScheduler)
	}

	return scheduler, nil
}

// WithConfig returns a Go context with a database instance.
func WithDatabase(ctx context.Context, db *db.DB) context.Context {
	return context.WithValue(ctx, ctxDatabase, db)
}

// Config attempts to extract the database instance from a Go context.
func Database(ctx context.Context) (*db.DB, error) {
	db, ok := ctx.Value(ctxDatabase).(*db.DB)
	if !ok {
		return nil, fmt.Errorf("missing context: %s", ctxDatabase)
	}

	return db, nil
}
