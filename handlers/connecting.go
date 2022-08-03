/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package handlers

import (
	"context"

	"github.com/SomusHQ/tacostand/contextutils"
	"github.com/slack-go/slack/socketmode"
)

// Connecting is the event handler that will run when the service attempts to
// connect to the Slack socket API.
func Connecting(ctx context.Context, event *socketmode.Event) error {
	logger, _ := contextutils.Logger(ctx)
	if logger != nil {
		logger.Debug("Connecting to Slack...")
	}

	return nil
}
