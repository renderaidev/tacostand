/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package commands

import (
	"context"
	"fmt"

	"github.com/SomusHQ/tacostand/builders"
	"github.com/SomusHQ/tacostand/contextutils"
	"github.com/SomusHQ/tacostand/db/models"
	"github.com/SomusHQ/tacostand/utils"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// Register is the handler for the "/register" slash command. It registers a
// workspace with the tacostand service.
func Register(ctx context.Context, e *socketmode.Event, command *slack.SlashCommand) error {
	client, _ := contextutils.SlackSocketClient(ctx)

	if !utils.IsWorkspaceAdmin(ctx, command.UserID) {
		payload := builders.MarkdownText("You are not an admin of this workspace")
		client.Ack(*e.Request, payload)
		return nil
	}

	db, _ := contextutils.Database(ctx)
	logger, _ := contextutils.Logger(ctx)

	var ct int64
	db.TeamModel.Where("id = ?", command.TeamID).Count(&ct)

	if ct != 0 {
		payload := builders.MarkdownText("This workspace is already registered. You can add questions using `/addquestion <question>` or use `/unregister` to remove the workspace.")
		client.Ack(*e.Request, payload)
		return nil
	}

	logger.Info(fmt.Sprintf("Registering team %s.", command.TeamID))

	team := models.NewTeam(command.TeamID)
	db.Create(&team)

	payload := builders.MarkdownText("Workspace registered successfully! Run `/addquestion <question>` to add a question to the daily standup report.")
	client.Ack(*e.Request, payload)

	logger.Success(fmt.Sprintf("Team %s registered successfully.", command.TeamID))

	return nil
}
