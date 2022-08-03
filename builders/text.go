/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package builders

import "github.com/slack-go/slack"

// PlainText generates a payload to send a plaintext response to the Slack
// socket service.
func PlainText(text string) map[string]interface{} {
	return map[string]interface{}{
		"text": text,
	}
}

// MarkdownText generates a payload to send a markdown response to the Slack
// socket service.
func MarkdownText(text string) map[string]interface{} {
	return map[string]interface{}{
		"text": text,
		"type": slack.MarkdownType,
	}
}
