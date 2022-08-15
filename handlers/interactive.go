/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package handlers

import (
	"context"

	"github.com/SomusHQ/tacostand/contextutils"
	"github.com/SomusHQ/tacostand/handlers/interactions"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// InteractionHandler is a generic handler for a Slack interaction. A handler
// must receive an application context, the Slack socket event object itself, as
// well as the interaction callback data. The returned value is an error.
type InteractionHandler func(context.Context, *socketmode.Event, *slack.InteractionCallback) error

var INTERACTION_HANDLERS = map[interface{}]InteractionHandler{
	slack.InteractionTypeBlockActions: interactions.BlockActions,
}

// Interactive is called whenever the application receives an interaction
// payload from Slack. It will map the interaction type to a handler and call it
// if there is one. It returns an error.
func Interactive(ctx context.Context, e *socketmode.Event) error {
	client, _ := contextutils.SlackSocketClient(ctx)

	data, ok := e.Data.(slack.InteractionCallback)
	if !ok {
		return nil
	}

	client.Ack(*e.Request)

	if handler, ok := INTERACTION_HANDLERS[data.Type]; ok {
		return handler(ctx, e, &data)
	}

	return nil
}
