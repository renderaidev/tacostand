/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package utils

import (
	"context"
	"regexp"
	"strings"

	"github.com/SomusHQ/tacostand/common"
	"github.com/SomusHQ/tacostand/contextutils"
)

// IsWorkspaceAdmin checks if a provided user is an administrator of the
// workspace.
func IsWorkspaceAdmin(ctx context.Context, userID string) bool {
	client, err := contextutils.SlackSocketClient(ctx)
	if err != nil {
		return false
	}

	user, err := client.GetUserInfo(userID)
	if err != nil {
		return false
	}

	return user.IsAdmin
}

// ExtractUserID will attempt to extract a user ID from a given string. It will
// match any string that satsifies the <@user_id> or <@user_id|user_name>
// pattern. Only the `user_id` part of the match will be returned.
func ExtractUserID(content string) (string, error) {
	normalizedContent := strings.TrimSpace(content)

	if len(normalizedContent) == 0 {
		return "", common.ErrMessageEmpty.Error()
	}

	re := regexp.MustCompile(`<@(.*?)>`)
	matches := re.FindStringSubmatch(normalizedContent)

	if len(matches) == 0 {
		return "", common.ErrMessageDoesNotContainUser.Error()
	}

	return strings.Split(matches[1], "|")[0], nil
}
