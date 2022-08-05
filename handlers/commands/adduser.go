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

// AddUser is the handler for the "/adduser" slash command. It adds a user to
// the list of users who will be notified by the daily standup in a workspace.
func AddUser(ctx context.Context, e *socketmode.Event, command *slack.SlashCommand) error {
	client, _ := contextutils.SlackSocketClient(ctx)

	if !utils.IsWorkspaceAdmin(ctx, command.UserID) {
		payload := builders.MarkdownText("You are not an admin of this workspace.")
		client.Ack(*e.Request, payload)
		return nil
	}

	db, _ := contextutils.Database(ctx)
	logger, _ := contextutils.Logger(ctx)

	var team models.Team
	result := db.Where("id = ?", command.TeamID).Take(&team)

	if db.NoMatch(result) {
		payload := builders.MarkdownText("This workspace has not been registered. Please run `/register` first.")
		client.Ack(*e.Request, payload)
		return nil
	}

	userID, err := utils.ExtractUserID(command.Text)
	if err != nil {
		payload := builders.MarkdownText(fmt.Sprintf("Failed to parse your input: `%s`", err))
		client.Ack(*e.Request, payload)
		return nil
	}

	var ct int64
	db.MemberModel.Where("id = ? AND team_id = ?", userID, command.TeamID).Count(&ct)

	if ct != 0 {
		payload := builders.MarkdownText("This user is already added.")
		client.Ack(*e.Request, payload)
		return nil
	}

	logger.Info(fmt.Sprintf("Adding member %s to team %s.", userID, command.TeamID))

	member := team.NewMember(userID)
	db.Save(member)

	logger.Success(fmt.Sprintf("Member %s added successfully to team %s.", userID, command.TeamID))

	payload := builders.PlainText("Member added successfully!")
	client.Ack(*e.Request, payload)

	return nil
}
