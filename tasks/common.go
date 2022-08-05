/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package tasks

import (
	"context"

	"github.com/SomusHQ/tacostand/contextutils"
	"github.com/SomusHQ/tacostand/db/models"
	"github.com/slack-go/slack"
)

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

		return err
	}

	_, _, err := api.PostMessage(member.ID, slack.MsgOptionText(question.Contents, false))
	return err
}
