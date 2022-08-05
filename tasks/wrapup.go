/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package tasks

import (
	"context"
	"fmt"

	"github.com/SomusHQ/tacostand/contextutils"
	"github.com/SomusHQ/tacostand/db/models"
	"github.com/SomusHQ/tacostand/utils"
	"github.com/slack-go/slack"
)

func informUserAboutWrapUp(ctx context.Context, memberID string, reportID uint64) {
	api, _ := contextutils.SlackApi(ctx)
	logger, _ := contextutils.Logger(ctx)

	opts := slack.MsgOptionBlocks(utils.BuildWrapUpBlock(reportID)...)
	_, _, err := api.PostMessage(memberID, opts)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to inform user about wrap up: %s", err))
	}
}

// HandleWrapUp will be run N minutes after the standup session has started. It
// checks if any user failed to submit their answers until this time, and if so,
// it will inform them about it.
func HandleWrapUp(ctx context.Context) func() {
	logger, _ := contextutils.Logger(ctx)
	api, _ := contextutils.SlackApi(ctx)
	db, _ := contextutils.Database(ctx)

	return func() {
		logger.Info("Received request to wrap up standup reports.")

		teaminfo, err := api.GetTeamInfo()
		if err != nil {
			logger.Danger(fmt.Sprintf("Failed to get team info: %s", err))
			return
		}

		var summary models.Summary
		result := db.Where("team_id = ?", teaminfo.ID).Last(&summary)
		if db.NoMatch(result) {
			logger.Danger("No standup reports found.")
			return
		}

		var questions []models.Question
		db.Where("team_id = ?", teaminfo.ID).Find(&questions)

		var reports []models.Report
		db.ReportModel.Preload("Answers").Where("summary_id = ? AND ongoing = ?", summary.ID, true).Find(&reports)

		for _, report := range reports {
			report.Ongoing = false
			db.Save(&report)

			if len(report.Answers) != len(questions) {
				logger.Info(fmt.Sprintf("Report %d has %d/%d answers. Informing user...", report.ID, len(report.Answers), len(questions)))
				informUserAboutWrapUp(ctx, report.MemberID, report.ID)
			}
		}

		db.Save(&summary)

		logger.Success("Wrap up complete.")
	}
}
