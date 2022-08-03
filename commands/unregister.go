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

// Unregister is the handler for the "/unregister" slash command. It deletes a
// workspace from the tacostand service. This is a very destructive command, as
// it will delete all data associated with the workspace. The threads that were
// created in the Slack workspace, however, will not be deleted.
func Unregister(ctx context.Context, e *socketmode.Event, command *slack.SlashCommand) error {
	client, _ := contextutils.SlackSocketClient(ctx)

	if !utils.IsWorkspaceAdmin(ctx, command.UserID) {
		payload := builders.MarkdownText("You are not an admin of this workspace")
		client.Ack(*e.Request, payload)
		return nil
	}

	db, _ := contextutils.Database(ctx)
	logger, _ := contextutils.Logger(ctx)

	var team models.Team
	result := db.Where("id = ?", command.TeamID).Take(&team)

	if db.NoMatch(result) {
		payload := builders.MarkdownText("This workspace has not been registered.")
		client.Ack(*e.Request, payload)
		return nil
	}

	logger.Info(fmt.Sprintf("Unregistering team %s.", command.TeamID))

	db.Delete(&team)

	payload := builders.MarkdownText("Workspace unregistered successfully.")
	client.Ack(*e.Request, payload)

	logger.Success(fmt.Sprintf("Team %s unregistered successfully.", command.TeamID))

	return nil
}
