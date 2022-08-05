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

// DeleteUser is the handler for the "/deluser" slash command. It deletes a
// user from the list of daily standup notifications in a workspace.
func DeleteUser(ctx context.Context, e *socketmode.Event, command *slack.SlashCommand) error {
	client, _ := contextutils.SlackSocketClient(ctx)

	if !utils.IsWorkspaceAdmin(ctx, command.UserID) {
		payload := builders.MarkdownText("You are not an admin of this workspace")
		client.Ack(*e.Request, payload)
		return nil
	}

	db, _ := contextutils.Database(ctx)
	logger, _ := contextutils.Logger(ctx)

	userID, err := utils.ExtractUserID(command.Text)
	if err != nil {
		payload := builders.MarkdownText(fmt.Sprintf("Failed to parse your input: `%s`", err))
		client.Ack(*e.Request, payload)
		return nil
	}

	var member models.Member
	result := db.Where("id = ? AND team_id = ?", userID, command.TeamID).Take(&member)

	if db.NoMatch(result) {
		payload := builders.MarkdownText("This user has not been added to the workspace.")
		client.Ack(*e.Request, payload)
		return nil
	}

	logger.Info(fmt.Sprintf("Deleting %s from team %s.", userID, command.TeamID))

	db.Delete(&member)

	payload := builders.MarkdownText("Member deleted successfully.")
	client.Ack(*e.Request, payload)

	logger.Success(fmt.Sprintf("Member %s deleted successfully.", command.TeamID))

	return nil
}
