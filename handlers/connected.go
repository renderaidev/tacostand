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

// Connected is the event handler that will run when the service connects to the
// Slack socket API.
func Connected(ctx context.Context, event *socketmode.Event) error {
	logger, _ := contextutils.Logger(ctx)
	logger.Success("Connected to Slack successfully!")

	return nil
}
