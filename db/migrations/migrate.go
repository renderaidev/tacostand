/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package migrations

import (
	"github.com/SomusHQ/tacostand/db/models"
	"gorm.io/gorm"
)

// Migrate will run GORM's AutoMigrate on every model.
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Team{})
	db.AutoMigrate(&models.Member{})
	db.AutoMigrate(&models.Question{})
	db.AutoMigrate(&models.Answer{})
	db.AutoMigrate(&models.Summary{})
	db.AutoMigrate(&models.Report{})
}
