/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package common

// BlockInteractionIdentifier represents the ID of a block interaction, which
// will be mapped to a handler.
type BlockInteractionIdentifier string

const (
	ADD_ANSWERS_BLOCK_ID  BlockInteractionIdentifier = "add_answers_block"
	ADD_ANSWERS_BUTTON_ID BlockInteractionIdentifier = "add_answers"
)
