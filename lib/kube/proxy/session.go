/*
Copyright 2021 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package proxy

import (
	"github.com/pborman/uuid"
	"sync"
)

type SessionState int

const (
	SessionPending SessionState = 1
	SessionRunning SessionState = 2
)

type Session struct {
	mu sync.Mutex
	state SessionState
	uuid string
	owner *authContext
	participants []*authContext
	notifier chan struct{}
}

func NewSession(ctx *authContext) *Session {
	participants := make([]*authContext, 0)
	participants = append(participants, ctx)

	return &Session {
		state: SessionPending,
		uuid: uuid.New(),
		owner: ctx,
		participants: participants,
		notifier: make(chan struct{}),
	}
}

func (s *Session) WaitOnStart(ctx *authContext) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.participants = append(s.participants, ctx)

	if s.state == SessionRunning {
		return
	}

	if sessionClearanceAcquired(s.participants) {
		s.state = SessionRunning
		close(s.notifier)
	} else {
		_, _ = <- s.notifier
	}
}

func sessionClearanceAcquired(participants []*authContext) bool {
	for _, participant := range participants {
		roles := participant.User.GetRoles()

		for _, role := range roles {
			if role == "admin" {
				return true
			}
		}
	}

	return false
}