/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package models

import "time"

// Team represents a singular Slack workspace.
type Team struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewTeam is a utility function that will add a workspace to the database.
func NewTeam(id string) *Team {
	return &Team{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewQuestion is a wrapper function that will associate a new question to the
// team object.
func (t *Team) NewQuestion(question string) *Question {
	return NewQuestion(question, t)
}

// NewMember is a wrapper function that will associate a new member to the team
// object.
func (t *Team) NewMember(userID string) *Member {
	return NewMember(userID, t)
}

// NewSummary creates a new stand-up summary for a team.
func (t *Team) NewSummary(threadID string) *Summary {
	return NewSummary(t, threadID)
}
