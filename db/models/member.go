/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package models

import "gorm.io/gorm"

// Member is a user in a workspace that will be notified about standup reports.
type Member struct {
	gorm.Model

	ID string `gorm:"primaryKey" json:"id"`

	TeamID string `json:"team_id"`
	Team   *Team  `gorm:"constraint:OnDelete:CASCADE" json:"team"`
}

// NewMember is a utility function that will create a new member that's
// associated with a workspace.
func NewMember(userID string, team *Team) *Member {
	return &Member{
		ID:   userID,
		Team: team,
	}
}
