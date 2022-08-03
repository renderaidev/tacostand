/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package models

import "gorm.io/gorm"

// Report is a singular standup report of a workspace member. This report is a
// collection of all the answers that the member provided to the questions in
// a standup report session. For daily standups, this object is a given user's
// answers for a specific day.
type Report struct {
	gorm.Model

	ID uint64 `gorm:"primaryKey" json:"id"`

	SummaryID uint64   `json:"summary_id"`
	Summary   *Summary `gorm:"constraint:OnDelete:CASCADE" json:"summary"`

	MemberID string  `json:"member_id"`
	Member   *Member `gorm:"constraint:OnDelete:CASCADE" json:"member"`

	Answers []*Answer `json:"answers"`
}
