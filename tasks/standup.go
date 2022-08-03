/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package tasks

import (
	"context"

	"github.com/SomusHQ/tacostand/contextutils"
)

// HandleDailyStandup is the handler that will run on a specified interval to
// ask the team for their standup answers.
//
// TODO: implement
func HandlePeriodicStandup(ctx context.Context) func() {
	logger, _ := contextutils.Logger(ctx)

	return func() {
		logger.Info("Received call to collect daily standup reports.")
	}
}
