/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package builders

import (
	"fmt"
	"strings"

	"github.com/SomusHQ/tacostand/db/models"
)

// StandupReportMessage will concatenate the questions and answers into a single
// message string.
func StandupReportMessage(answers []*models.Answer) string {
	output := make([]string, 0)

	for _, answer := range answers {
		output = append(output, fmt.Sprintf("*%s*\n%s", answer.Question.Contents, answer.Contents))
	}

	return strings.Join(output, "\n\n")
}
