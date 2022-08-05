/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package common

import "fmt"

// ConfigError represents an error with the service configuration.
type ConfigError int

const (
	ErrInvalidSlackAppToken ConfigError = iota
	ErrInvalidSlackBotToken
	ErrSlackReportChannelNameNotSet

	ErrInvalidCronExpression
)

func (e ConfigError) Error() error {
	switch e {
	case ErrInvalidSlackAppToken:
		return fmt.Errorf("invalid Slack app token (should start with `xapp-`)")
	case ErrInvalidSlackBotToken:
		return fmt.Errorf("invalid Slack bot token (should start with `xoxb-`)")
	case ErrSlackReportChannelNameNotSet:
		return fmt.Errorf("slack report channel name has not been specified")
	case ErrInvalidCronExpression:
		return fmt.Errorf("the provided cron expression is invalid")
	default:
		return fmt.Errorf("configuration error")
	}
}

// ParseError is a generic error that represents an error that is raised when
// parsing any kind of data in the application.
type ParseError int

const (
	ErrMessageEmpty ParseError = iota
	ErrMessageDoesNotContainUser
)

func (e ParseError) Error() error {
	switch e {
	case ErrMessageEmpty:
		return fmt.Errorf("message is empty")
	case ErrMessageDoesNotContainUser:
		return fmt.Errorf("message does not contain user")
	default:
		return fmt.Errorf("parse error")
	}
}

// StandupSchedulerError is an error that is raised when the standup scheduler
// fails to perform its job.
type StandupSchedulerError int

const (
	ErrChannelNotFound StandupSchedulerError = iota
	ErrFailedToGetTeamInfo
	ErrTeamNotRegistered
)

func (e StandupSchedulerError) Error() error {
	switch e {
	case ErrChannelNotFound:
		return fmt.Errorf("standup report channel not found")
	case ErrFailedToGetTeamInfo:
		return fmt.Errorf("failed to get team info")
	case ErrTeamNotRegistered:
		return fmt.Errorf("team not registered")
	default:
		return fmt.Errorf("standup scheduler error")
	}
}
