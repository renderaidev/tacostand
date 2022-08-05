/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package callbacks

import (
	"context"
	"fmt"

	"github.com/SomusHQ/tacostand/contextutils"
	"github.com/SomusHQ/tacostand/db/models"
	"github.com/SomusHQ/tacostand/tasks"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// Message is the callback handler that's used on a "message" Slack event type.
// This is the event handler that's used for communicating with a user during a
// stand-up session.
func Message(ctx context.Context, event *socketmode.Event, innerEvent *slackevents.EventsAPIInnerEvent) error {
	logger, _ := contextutils.Logger(ctx)
	db, _ := contextutils.Database(ctx)

	data, ok := innerEvent.Data.(*slackevents.MessageEvent)
	if !ok {
		return nil
	}

	if data.ChannelType != "im" {
		return nil
	}

	var member models.Member
	result := db.Where("id = ?", data.User).Take(&member)
	if db.NoMatch(result) {
		return nil
	}

	err := tasks.AskQuestion(ctx, &member, data.Text)
	if err != nil {
		logger.Danger(fmt.Sprintf("Error asking question: %s", err))
		return err
	}

	return nil
}
