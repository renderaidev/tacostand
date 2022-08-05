/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package inquirer

import (
	"github.com/SomusHQ/tacostand/db/models"
)

// QuestionQueue is a slice of questions to be asked from a specific user.
// Within the Inquirer, these queues will be mapped to every member that has to
// answer standup questions.
type QuestionQueue []*models.Question

// Inquirer is a service that will queue up questions to be asked from members.
type Inquirer struct {
	queue        map[string]QuestionQueue
	lastQuestion map[string]*models.Question
}

// NewInquirer creates a new inquirer instance.
func NewInquirer() *Inquirer {
	return &Inquirer{
		queue:        make(map[string]QuestionQueue),
		lastQuestion: make(map[string]*models.Question),
	}
}

// Enqueue will add a list of questions to the queue for a specific member,
// creating a new queue if one doesn't exist.
func (inq *Inquirer) Enqueue(memberID string, questions []*models.Question) {
	if _, ok := inq.queue[memberID]; !ok {
		inq.queue[memberID] = make(QuestionQueue, 0)
	}

	inq.queue[memberID] = append(inq.queue[memberID], questions...)
}

// Pop will return the next question in the queue for a specific member. It will
// also remove it from the queue. If the queue is empty or the member doesn't
// have a queue, it will return nil.
func (inq *Inquirer) Pop(memberID string) *models.Question {
	if _, ok := inq.queue[memberID]; !ok {
		return nil
	}

	if len(inq.queue[memberID]) == 0 {
		return nil
	}

	question := inq.queue[memberID][0]
	inq.lastQuestion[memberID] = question
	inq.queue[memberID] = inq.queue[memberID][1:]

	return question
}

// Destroy will remove the queue for a specific member.
func (inq *Inquirer) Destroy(memberID string) {
	delete(inq.queue, memberID)
	delete(inq.lastQuestion, memberID)
}

// Exists checks if a given member has a queue.
func (inq *Inquirer) Exists(memberID string) bool {
	_, ok := inq.queue[memberID]
	return ok
}

// Last will return the last popped question for a specific member, or nil, if
// no such question exists.
func (inq *Inquirer) Last(memberID string) *models.Question {
	if _, ok := inq.lastQuestion[memberID]; !ok {
		return nil
	}

	return inq.lastQuestion[memberID]
}
