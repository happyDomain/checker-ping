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

func TestPacketLossValidateOptions(t *testing.T) {
	r := &packetLossRule{}
	cases := []struct {
		name    string
		opts    sdk.CheckerOptions
		wantErr bool
	}{
		{"defaults", sdk.CheckerOptions{}, false},
		{"valid", sdk.CheckerOptions{"warningPacketLoss": 5.0, "criticalPacketLoss": 25.0}, false},
		{"warn negative", sdk.CheckerOptions{"warningPacketLoss": -1.0, "criticalPacketLoss": 50.0}, true},
		{"crit over 100", sdk.CheckerOptions{"warningPacketLoss": 10.0, "criticalPacketLoss": 150.0}, true},
		{"crit equal warn", sdk.CheckerOptions{"warningPacketLoss": 20.0, "criticalPacketLoss": 20.0}, true},
		{"crit below warn", sdk.CheckerOptions{"warningPacketLoss": 30.0, "criticalPacketLoss": 10.0}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := r.ValidateOptions(c.opts)
			if (err != nil) != c.wantErr {
				t.Errorf("err=%v, wantErr=%v", err, c.wantErr)
			}
		})
	}
}

func TestPacketLossEvaluate(t *testing.T) {
	r := &packetLossRule{}
	ctx := context.Background()
	opts := sdk.CheckerOptions{"warningPacketLoss": 10.0, "criticalPacketLoss": 50.0}

	states := r.Evaluate(ctx, obsWith(
		PingTargetResult{Address: "ok", PacketLoss: 0, Sent: 5, Received: 5},
		PingTargetResult{Address: "warn", PacketLoss: 20, Sent: 5, Received: 4},
		PingTargetResult{Address: "crit", PacketLoss: 80, Sent: 5, Received: 1},
	), opts)

	if len(states) != 3 {
		t.Fatalf("got %d states, want 3", len(states))
	}
	want := []sdk.Status{sdk.StatusOK, sdk.StatusWarn, sdk.StatusCrit}
	for i, s := range states {
		if s.Status != want[i] {
			t.Errorf("state[%d] status = %v, want %v", i, s.Status, want[i])
		}
	}
}

func TestPacketLossEvaluateNoTargets(t *testing.T) {
	r := &packetLossRule{}
	states := r.Evaluate(context.Background(), obsWith(), sdk.CheckerOptions{})
	if len(states) != 1 || states[0].Status != sdk.StatusUnknown {
		t.Errorf("expected single Unknown state, got %+v", states)
	}
}
