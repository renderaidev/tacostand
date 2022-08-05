/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package models

import "gorm.io/gorm"

// Answer represents a user's answer to a question in a specific report.
type Answer struct {
	gorm.Model

	ID uint64 `gorm:"primaryKey" json:"id"`

	QuestionID uint64    `json:"question_id"`
	Question   *Question `gorm:"constraint:OnDelete:CASCADE" json:"question"`

	ReportID uint64  `json:"report_id"`
	Report   *Report `gorm:"constraint:OnDelete:CASCADE" json:"report"`

	Contents string `json:"contents"`
}

func NewAnswer(question *Question, report *Report, contents string) *Answer {
	return &Answer{
		Question: question,
		Report:   report,
		Contents: contents,
	}
}
