/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package utils

import (
	"fmt"
	"time"

	"github.com/slack-go/slack"
)

// BuildStandupMessage will return a string that will be sent to the standup
// channel when a standup session is created.
func BuildStandupMessage() string {
	today := time.Now().Format("Monday, January 2")
	return fmt.Sprintf("Here are the results of the `Standup Report` for *%s*", today)
}

// BuildWrapUpBlock will return a Slack block that will be sent to a member if
// they failed to submit their responses in time. The block will include an
// informational text and a button that will allow them to reopen their report
// and submit the answers anyway.
func BuildWrapUpBlock(reportID uint64) []slack.Block {
	id := fmt.Sprintf("%d", reportID)

	blocks := make([]slack.Block, 0)

	textblock := slack.NewSectionBlock(
		slack.NewTextBlockObject(slack.MarkdownType, "Aw, you're out of time... 😔 I will wrap up the responses now.", false, false),
		nil,
		nil,
	)

	buttonblock := slack.NewActionBlock(
		"add_answers_block",
		slack.NewButtonBlockElement("add_answers", id, slack.NewTextBlockObject(slack.PlainTextType, "Add answers", false, false)),
	)

	blocks = append(blocks, textblock, buttonblock)

	return blocks
}
