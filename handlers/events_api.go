/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package handlers

import (
	"context"

	"github.com/SomusHQ/tacostand/contextutils"
	"github.com/SomusHQ/tacostand/handlers/callbacks"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// CallbackHandler is a generic handler for a Slack callback event. A handler
// must receive an application context, the raw event, and the Slack inner event
// as arguments. The returned value is an error.
type CallbackHandler func(context.Context, *socketmode.Event, *slackevents.EventsAPIInnerEvent) error

var CALLBACK_HANDLERS = map[interface{}]CallbackHandler{
	"message": callbacks.Message,
}

// CallbackEvent is used for handling all callback events that are received from
// the Slack socket API.
func CallbackEvent(ctx context.Context, e *socketmode.Event, innerEvent *slackevents.EventsAPIInnerEvent) error {
	if handler, ok := CALLBACK_HANDLERS[innerEvent.Type]; ok {
		return handler(ctx, e, innerEvent)
	}

	return nil
}

// EventsAPIEvent is the handler that is called on a Slack events API event.
func EventsAPIEvent(ctx context.Context, e *socketmode.Event) error {
	client, _ := contextutils.SlackSocketClient(ctx)

	data, ok := e.Data.(slackevents.EventsAPIEvent)
	if !ok {
		return nil
	}

	client.Ack(*e.Request)

	switch data.Type {
	case slackevents.CallbackEvent:
		innerEvent := data.InnerEvent
		return CallbackEvent(ctx, e, &innerEvent)
	default:
		return nil
	}
}
