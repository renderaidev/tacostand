/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package tasks

import (
	"context"

	"github.com/SomusHQ/tacostand/builders"
	"github.com/SomusHQ/tacostand/common"
	"github.com/SomusHQ/tacostand/contextutils"
	"github.com/SomusHQ/tacostand/db/models"
	"github.com/slack-go/slack"
)

func getReportChannel(ctx context.Context) *slack.Channel {
	api, _ := contextutils.SlackApi(ctx)
	config, _ := contextutils.Config(ctx)

	var channel *slack.Channel

	if channels, _, err := api.GetConversations(&slack.GetConversationsParameters{
		ExcludeArchived: true,
	}); err == nil {
		for _, c := range channels {
			if c.Name == config.Slack.ReportChannelName {
				channel = &c
				break
			}
		}
	}

	return channel
}

func wrapReport(ctx context.Context, user *slack.User, report *models.Report) error {
	db, _ := contextutils.Database(ctx)
	api, _ := contextutils.SlackApi(ctx)

	var summary models.Summary
	result := db.SummaryModel.Where("id = ?", report.SummaryID).Take(&summary)
	if db.NoMatch(result) {
		return common.ErrSummaryNotFound.Error()
	}

	var answers []*models.Answer
	result = db.AnswerModel.Preload("Question").Where("report_id = ?", report.ID).Find(&answers)
	if db.NoMatch(result) || len(answers) == 0 {
		return common.ErrAnswersNotFound.Error()
	}

	channel := getReportChannel(ctx)

	message := builders.StandupReportMessage(answers)

	username := user.RealName
	if len(username) == 0 {
		username = user.Name
	}

	_, _, err := api.PostMessage(
		channel.ID,
		slack.MsgOptionText(message, false),
		slack.MsgOptionTS(summary.ThreadID),
		slack.MsgOptionUsername(username),
		slack.MsgOptionIconURL(user.Profile.ImageOriginal),
	)

	return err
}

// AskQuestion will ask a single question from a member, and if the `answer`
// argument is provided, it will also answer the last popped question. It will
// do nothing if the provided member doesn't have a current question queue and
// if the queue became empty, it will send an informational message to the user
// and wrap up their report.
func AskQuestion(ctx context.Context, member *models.Member, answer string) error {
	db, _ := contextutils.Database(ctx)
	api, _ := contextutils.SlackApi(ctx)
	inquirer, _ := contextutils.Inquirer(ctx)

	if !inquirer.Exists(member.ID) {
		return nil
	}

	var report models.Report
	result := db.Where("member_id = ? AND ongoing = ?", member.ID, true).Last(&report)
	if db.NoMatch(result) {
		inquirer.Destroy(member.ID)
		return nil
	}

	if len(answer) != 0 {
		last := inquirer.Last(member.ID)
		if last != nil {
			a := report.NewAnswer(last, answer)
			db.Save(&a)
		}
	}

	question := inquirer.Pop(member.ID)
	if question == nil {
		_, _, err := api.PostMessage(member.ID, slack.MsgOptionText("Woo-hoo! 🎉 Thank you for completing your check-in. Have a taco: 🌮", false))

		report.Ongoing = false
		db.Save(&report)

		inquirer.Destroy(member.ID)

		if err != nil {
			return err
		}

		user, err := api.GetUserInfo(member.ID)
		if err != nil {
			return err
		}

		return wrapReport(ctx, user, &report)
	}

	_, _, err := api.PostMessage(member.ID, slack.MsgOptionText(question.Contents, false))
	return err
}
