/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package handlers

import (
	"context"
	"fmt"

	"github.com/SomusHQ/tacostand/builders"
	"github.com/SomusHQ/tacostand/commands"
	"github.com/SomusHQ/tacostand/contextutils"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// CommandHandler is a generic handler for a Slack slash command. A handler must
// receive an application context, the raw event, and the slash command data.
// The returned value is an error.
type CommandHandler func(context.Context, *socketmode.Event, *slack.SlashCommand) error

var SLASH_COMMAND_HANDLERS = map[string]CommandHandler{
	"/register":    commands.Register,
	"/unregister":  commands.Unregister,
	"/addquestion": commands.AddQuestion,
	"/adduser":     commands.AddUser,
	"/deluser":     commands.DeleteUser,
}

// SlashCommand is a handler for a slash command. It will try to map the command
// to a handler and call it. It will also log any error that occurs in the
// command handler.
func SlashCommand(ctx context.Context, e *socketmode.Event) error {
	client, _ := contextutils.SlackSocketClient(ctx)

	data, ok := e.Data.(slack.SlashCommand)
	if !ok {
		return nil
	}

	if handler, ok := SLASH_COMMAND_HANDLERS[data.Command]; ok {
		return handler(ctx, e, &data)
	}

	payload := builders.MarkdownText(fmt.Sprintf("*Unknown command*: `%s`", data.Command))
	client.Ack(*e.Request, payload)

	return nil
}
