// This file is part of the happyDomain (R) project.
// Copyright (c) 2020-2026 happyDomain
// Authors: Pierre-Olivier Mercier, et al.
//
// This program is offered under a commercial and under the AGPL license.
// For commercial licensing, contact us at <contact@happydomain.org>.
//
// For AGPL licensing:
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package checker

import (
	"context"
	"testing"

	sdk "git.happydns.org/checker-sdk-go/checker"
)

func TestReachabilityEvaluate(t *testing.T) {
	r := &reachabilityRule{}
	states := r.Evaluate(context.Background(), obsWith(
		PingTargetResult{Address: "up", Sent: 5, Received: 5},
		PingTargetResult{Address: "down", Sent: 5, Received: 0},
	), sdk.CheckerOptions{})

	if len(states) != 2 {
		t.Fatalf("got %d states, want 2", len(states))
	}
	if states[0].Status != sdk.StatusOK || states[0].Code != "ping.reachable.ok" {
		t.Errorf("up: %+v", states[0])
	}
	if states[1].Status != sdk.StatusCrit || states[1].Code != "ping.reachable.unreachable" {
		t.Errorf("down: %+v", states[1])
	}
}

func TestReachabilityNoTargets(t *testing.T) {
	r := &reachabilityRule{}
	states := r.Evaluate(context.Background(), obsWith(), sdk.CheckerOptions{})
	if len(states) != 1 || states[0].Status != sdk.StatusUnknown {
		t.Errorf("expected single Unknown state, got %+v", states)
	}
}
