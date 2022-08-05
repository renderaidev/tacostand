/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/SomusHQ/tacostand/builders"
	"github.com/SomusHQ/tacostand/contextutils"
	"github.com/SomusHQ/tacostand/db/models"
	"github.com/SomusHQ/tacostand/utils"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// AddQuestion is the handler for the "/addquestion" slash command. It adds a
// question to a workspace's list of questions.
func AddQuestion(ctx context.Context, e *socketmode.Event, command *slack.SlashCommand) error {
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

	contents := strings.TrimSpace(command.Text)

	if len(contents) == 0 {
		payload := builders.PlainText("Please provide a question.")
		client.Ack(*e.Request, payload)
		return nil
	}

	logger.Info(fmt.Sprintf("Adding question \"%s\" to team %s.", contents, command.TeamID))

	question := team.NewQuestion(contents)
	db.Save(question)

	logger.Success(fmt.Sprintf("Question added successfully to team %s.", command.TeamID))

	payload := builders.PlainText("Question added successfully!")
	client.Ack(*e.Request, payload)

	return nil
}
