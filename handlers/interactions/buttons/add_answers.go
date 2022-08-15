/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package buttons

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/SomusHQ/tacostand/contextutils"
	"github.com/SomusHQ/tacostand/db/models"
	"github.com/SomusHQ/tacostand/tasks"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func handlePersonalWrapUp(ctx context.Context, reportID uint64) func() {
	api, _ := contextutils.SlackApi(ctx)
	db, _ := contextutils.Database(ctx)
	logger, _ := contextutils.Logger(ctx)
	inquirer, _ := contextutils.Inquirer(ctx)

	return func() {
		teaminfo, err := api.GetTeamInfo()
		if err != nil {
			return
		}

		var questions []models.Question
		db.Where("team_id = ?", teaminfo.ID).Find(&questions)

		var report models.Report
		result := db.Preload("Answers").Where("id = ?", reportID).Take(&report)
		if db.NoMatch(result) {
			logger.Warn(fmt.Sprintf("Failed to find report with ID %d while scheduling personal wrap up.", reportID))
			return
		}

		report.Ongoing = false
		db.Save(&report)

		if inquirer.Exists(report.MemberID) {
			inquirer.Destroy(report.MemberID)
		}

		if len(report.Answers) != len(questions) {
			tasks.InformUserAboutWrapUp(ctx, report.MemberID, report.ID)
		}
	}
}

func scheduleWrapUp(ctx context.Context, report *models.Report) {
	config, _ := contextutils.Config(ctx)
	time.AfterFunc(time.Duration(config.WrapUpTime)*time.Minute, handlePersonalWrapUp(ctx, report.ID))
}

// AddAnswers is the handler that's called on the `add_answers` button
// interaction. The button's value is the numeric ID of the report which needs
// to be reopened. The handler will reopen the report and start prompting the
// user for answers. It will also schedule a personal wrap-up in case the user
// fails to answer the questions again in time. This callback will also perform
// a soft-delete on the already existing answers.
func AddAnswers(ctx context.Context, event *socketmode.Event, data *slack.InteractionCallback, action *slack.BlockAction) error {
	logger, _ := contextutils.Logger(ctx)

	logger.Info(fmt.Sprintf("Received request to add answers from user %s.", data.User.ID))

	reportID, err := strconv.Atoi(action.Value)
	if err != nil {
		return err
	}

	db, _ := contextutils.Database(ctx)
	api, _ := contextutils.SlackApi(ctx)
	inquirer, _ := contextutils.Inquirer(ctx)

	var member models.Member
	result := db.Where("id = ?", data.User.ID).Take(&member)
	if db.NoMatch(result) {
		logger.Warn(fmt.Sprintf("User %s is not a member of the team.", data.User.ID))
		return nil
	}

	// clean up stale answers
	db.Where("report_id = ?", reportID).Delete(&models.Answer{})

	// reopen report
	var report models.Report
	result = db.Where("id = ? AND ongoing = ?", reportID, false).Take(&report)
	if db.NoMatch(result) {
		logger.Warn(fmt.Sprintf("Report %d does not exist or is currently open.", reportID))
		return nil
	}

	report.Ongoing = true
	db.Save(&report)

	// prompt user for new answers
	var questions []*models.Question
	db.Where("team_id = ?", member.TeamID).Find(&questions)

	inquirer.Enqueue(report.MemberID, questions)
	err = tasks.AskQuestion(ctx, &member, "")

	if err == nil {
		// schedule personal wrap-up
		logger.Info(fmt.Sprintf("Scheduling personal wrap-up for user %s and report %d.", data.User.ID, reportID))
		scheduleWrapUp(ctx, &report)

		return nil
	}

	logger.Danger(err)

	// something went wrong, inform user and close the report again
	report.Ongoing = false
	db.Save(&report)

	message := slack.MsgOptionText("Something went wrong while rescheduling your questions. Please try again later or inform an administrator.", false)
	api.PostMessage(report.MemberID, message)

	return err
}
