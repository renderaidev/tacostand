/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package interactions

import (
	"context"

	"github.com/SomusHQ/tacostand/common"
	"github.com/SomusHQ/tacostand/handlers/interactions/buttons"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// ButtonHandler is a generic handler for a Slack button interaction event. The
// handler must receive an application context, the Slack socket event object
// itself, the interaction callback data, and the button block action object.
// The returned value is an error.
type ButtonHandler func(context.Context, *socketmode.Event, *slack.InteractionCallback, *slack.BlockAction) error

var BUTTON_HANDLERS = map[common.BlockInteractionIdentifier]ButtonHandler{
	common.ADD_ANSWERS_BUTTON_ID: buttons.AddAnswers,
}

// BlockActions is called whenever the application receives a block interaction
// from Slack. Typically, this is used for button interactions, so it will
// essentially map a button interaction identifier to a handler, and call it if
// it exists. It returns an error.
func BlockActions(ctx context.Context, event *socketmode.Event, data *slack.InteractionCallback) error {
	if len(data.ActionCallback.BlockActions) == 0 {
		return nil
	}

	firstAction := data.ActionCallback.BlockActions[0]

	if handler, ok := BUTTON_HANDLERS[common.BlockInteractionIdentifier(firstAction.ActionID)]; ok {
		return handler(ctx, event, data, firstAction)
	}

	return nil
}
