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

// ConnectionError is the event handler that will run when the service
// encounters an error while connecting to the Slack socket API.
func ConnectionError(ctx context.Context, e *socketmode.Event) error {
	logger, _ := contextutils.Logger(ctx)
	logger.Danger("Failed to connect to Slack, will retry...")

	return nil
}
