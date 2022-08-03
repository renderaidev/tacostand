/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SomusHQ/tacostand/config"
	"github.com/SomusHQ/tacostand/db/migrations"
	"github.com/SomusHQ/tacostand/db/models"
	logs "github.com/SomusHQ/tacostand/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is a wrapper around the GORM database object. It also contains references
// to the various models.
type DB struct {
	*gorm.DB

	TeamModel     *gorm.DB
	MemberModel   *gorm.DB
	QuestionModel *gorm.DB
	AnswerModel   *gorm.DB
	ReportModel   *gorm.DB
	SummaryModel  *gorm.DB
}

// New creates a new database connection and runs the auto-migrations.
func New(config *config.Config, appLogger *logs.Logger) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s database=%s sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
	)

	dbLogger := logger.New(
		log.New(os.Stdout, "db: ", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
		},
	)

	appLogger.Info("Opening connection to Postgres database...")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})

	if err != nil {
		return nil, err
	}

	appLogger.Success("Database connection established successfully.")
	appLogger.Info("Running database migrations...")

	migrations.Migrate(db)

	appLogger.Success("Database migrations completed successfully.")

	return &DB{
		DB: db,

		TeamModel:     db.Model(&models.Team{}),
		MemberModel:   db.Model(&models.Member{}),
		QuestionModel: db.Model(&models.Question{}),
		AnswerModel:   db.Model(&models.Answer{}),
		ReportModel:   db.Model(&models.Report{}),
		SummaryModel:  db.Model(&models.Summary{}),
	}, nil
}

func (db *DB) NoMatch(tx *gorm.DB) bool {
	return errors.Is(tx.Error, gorm.ErrRecordNotFound)
}
