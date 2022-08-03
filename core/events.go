/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package core

import (
	"context"

	"github.com/SomusHQ/tacostand/handlers"
	"github.com/slack-go/slack/socketmode"
)

// EventHandler is a generic handler for a Slack socket event. A handler must
// receive an application context and the Slack socket event object itself. The
// returned value is an error.
type EventHandler func(context.Context, *socketmode.Event) error

var EVENT_HANDLERS = map[socketmode.EventType]EventHandler{
	socketmode.EventTypeConnecting:      handlers.Connecting,
	socketmode.EventTypeConnected:       handlers.Connected,
	socketmode.EventTypeConnectionError: handlers.ConnectionError,
	socketmode.EventTypeSlashCommand:    handlers.SlashCommand,
}

// HandleEvents is called in the goroutine channel that receives socket events.
// It will try to map the event type to a handler and call it. It will also log
// any error that occurs in the event handler.
func (c *Client) HandleEvents(ctx context.Context) {
	for event := range c.SlackSocketClient.Events {
		if handler, ok := EVENT_HANDLERS[event.Type]; ok {
			err := handler(ctx, &event)
			if err != nil {
				c.Logger.Danger(err.Error())
			}
		}
	}
}
