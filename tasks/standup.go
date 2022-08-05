/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/SomusHQ/tacostand/common"
	"github.com/SomusHQ/tacostand/contextutils"
	"github.com/SomusHQ/tacostand/db/models"
	"github.com/SomusHQ/tacostand/utils"
	"github.com/slack-go/slack"
)

func createReportChannelMessage(ctx context.Context) (string, *slack.Channel, *models.Team, error) {
	api, _ := contextutils.SlackApi(ctx)
	db, _ := contextutils.Database(ctx)

	channel := getReportChannel(ctx)
	if channel == nil {
		return "", nil, nil, common.ErrChannelNotFound.Error()
	}

	teaminfo, err := api.GetTeamInfo()
	if err != nil {
		return "", nil, nil, common.ErrFailedToGetTeamInfo.Error()
	}

	var team models.Team
	result := db.Where("id = ?", teaminfo.ID).Take(&team)
	if db.NoMatch(result) {
		return "", nil, nil, common.ErrTeamNotRegistered.Error()
	}

	_, ts, err := api.PostMessage(channel.ID, slack.MsgOptionText(utils.BuildStandupMessage(), false))
	if err != nil {
		return "", nil, nil, err
	}

	return ts, channel, &team, nil
}

func createDatabaseEntry(ctx context.Context, team *models.Team, threadID string) *models.Summary {
	db, _ := contextutils.Database(ctx)

	summary := team.NewSummary(threadID)
	db.Save(&summary)

	return summary
}

func promptMembers(ctx context.Context, summary *models.Summary, team *models.Team) error {
	api, _ := contextutils.SlackApi(ctx)
	db, _ := contextutils.Database(ctx)
	inquirer, _ := contextutils.Inquirer(ctx)

	var members []models.Member
	db.MemberModel.Where("team_id = ?", team.ID).Find(&members)

	if len(members) == 0 {
		return nil
	}

	var questions []*models.Question
	db.QuestionModel.Where("team_id = ?", team.ID).Find(&questions)

	for _, member := range members {
		_, _, err := api.PostMessage(member.ID, slack.MsgOptionText("Hey! It's time to complete the daily standup 🎉", false))

		if err != nil {
			continue
		}

		report := summary.NewReport(&member)
		db.Save(&report)

		inquirer.Enqueue(member.ID, questions)
		AskQuestion(ctx, &member, "")
	}

	return nil
}

func scheduleWrapUp(ctx context.Context) {
	config, _ := contextutils.Config(ctx)
	time.AfterFunc(time.Duration(config.WrapUpTime)*time.Minute, HandleWrapUp(ctx))
}

// HandleDailyStandup is the handler that will run on a specified interval to
// ask the team for their standup answers.
func HandlePeriodicStandup(ctx context.Context) func() {
	logger, _ := contextutils.Logger(ctx)

	return func() {
		logger.Info("Received call to collect daily standup reports.")

		ts, channel, team, err := createReportChannelMessage(ctx)
		if err != nil {
			logger.Danger(err)
			return
		}

		logger.Success("Successfully sent standup report summary to channel: %s", channel.Name)

		summary := createDatabaseEntry(ctx, team, ts)

		logger.Success(fmt.Sprintf("Created summary with ID #%d (thread ts: %s)", summary.ID, ts))

		if err := promptMembers(ctx, summary, team); err != nil {
			logger.Danger(err)
			return
		}

		logger.Success("Successfully prompted members to answer standup report.")

		scheduleWrapUp(ctx)

		logger.Success("Successfully scheduled wrap up of standup report.")
	}
}
