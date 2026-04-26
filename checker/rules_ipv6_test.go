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

func TestIPv6ReachabilityEvaluate(t *testing.T) {
	r := &ipv6ReachabilityRule{}
	ctx := context.Background()

	cases := []struct {
		name     string
		targets  []PingTargetResult
		wantStat sdk.Status
		wantCode string
	}{
		{
			name:     "no v6 targets",
			targets:  []PingTargetResult{{Address: "1.1.1.1", Sent: 5, Received: 5}},
			wantStat: sdk.StatusUnknown,
			wantCode: "ping.ipv6_reachable.skipped",
		},
		{
			name: "v6 reachable",
			targets: []PingTargetResult{
				{Address: "1.1.1.1", Sent: 5, Received: 5},
				{Address: "2001:db8::1", Sent: 5, Received: 5},
			},
			wantStat: sdk.StatusOK,
			wantCode: "ping.ipv6_reachable.ok",
		},
		{
			name: "all v6 unreachable",
			targets: []PingTargetResult{
				{Address: "2001:db8::1", Sent: 5, Received: 0},
				{Address: "2001:db8::2", Sent: 5, Received: 0},
			},
			wantStat: sdk.StatusWarn,
			wantCode: "ping.ipv6_reachable.unreachable",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			states := r.Evaluate(ctx, obsWith(c.targets...), sdk.CheckerOptions{})
			if len(states) != 1 {
				t.Fatalf("got %d states, want 1", len(states))
			}
			if states[0].Status != c.wantStat {
				t.Errorf("status = %v, want %v", states[0].Status, c.wantStat)
			}
			if states[0].Code != c.wantCode {
				t.Errorf("code = %q, want %q", states[0].Code, c.wantCode)
			}
		})
	}
}
