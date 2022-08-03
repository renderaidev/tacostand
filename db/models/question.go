/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package models

import "gorm.io/gorm"

// Question is a question that will be asked in a standup report.
type Question struct {
	gorm.Model

	ID       uint64 `gorm:"primaryKey" json:"id"`
	Contents string `json:"contents"`

	TeamID string `json:"team_id"`
	Team   *Team  `gorm:"constraint:OnDelete:CASCADE" json:"team"`
}

// NewQuestion is a utility function that will create a new question that's
// associated with a workspace.
func NewQuestion(contents string, team *Team) *Question {
	return &Question{
		Contents: contents,
		Team:     team,
	}
}
