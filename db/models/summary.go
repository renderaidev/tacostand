/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package models

import (
	"time"

	"gorm.io/gorm"
)

// Summary is a collection of standup reports that are collected from all the
// members that need to be notified. For example, in daily standups, this will
// represent all reports that were sent on a specific day. It also contains the
// ID of the thread that was created for collecting the responses.
type Summary struct {
	gorm.Model

	ID uint64 `gorm:"primaryKey" json:"id"`

	TeamID string `json:"team_id"`
	Team   *Team  `gorm:"constraint:OnDelete:CASCADE" json:"team"`

	ThreadID string `json:"thread_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewSummary creates a new summary instance for a team. It also saves the ID
// of the thread that was created for collecting the report summaries.
func NewSummary(team *Team, threadID string) *Summary {
	return &Summary{
		Team:      team,
		ThreadID:  threadID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewReport adds a new report for a specific member to the current stand-up
// summary.
func (s *Summary) NewReport(member *Member) *Report {
	return NewReport(s, member)
}
