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

func TestRTTValidateOptions(t *testing.T) {
	r := &rttRule{}
	cases := []struct {
		name    string
		opts    sdk.CheckerOptions
		wantErr bool
	}{
		{"defaults", sdk.CheckerOptions{}, false},
		{"valid", sdk.CheckerOptions{"warningRTT": 50.0, "criticalRTT": 200.0}, false},
		{"warn zero", sdk.CheckerOptions{"warningRTT": 0.0, "criticalRTT": 100.0}, true},
		{"crit negative", sdk.CheckerOptions{"warningRTT": 50.0, "criticalRTT": -1.0}, true},
		{"crit equal warn", sdk.CheckerOptions{"warningRTT": 100.0, "criticalRTT": 100.0}, true},
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

func TestRTTEvaluate(t *testing.T) {
	r := &rttRule{}
	ctx := context.Background()
	opts := sdk.CheckerOptions{"warningRTT": 100.0, "criticalRTT": 500.0}

	states := r.Evaluate(ctx, obsWith(
		PingTargetResult{Address: "ok", RTTAvg: 50, Sent: 5, Received: 5},
		PingTargetResult{Address: "warn", RTTAvg: 200, Sent: 5, Received: 5},
		PingTargetResult{Address: "crit", RTTAvg: 600, Sent: 5, Received: 5},
		PingTargetResult{Address: "dead", RTTAvg: 0, Sent: 5, Received: 0},
	), opts)

	if len(states) != 4 {
		t.Fatalf("got %d states, want 4", len(states))
	}
	want := []sdk.Status{sdk.StatusOK, sdk.StatusWarn, sdk.StatusCrit, sdk.StatusUnknown}
	wantCode := []string{"ping.rtt.ok", "ping.rtt.warning", "ping.rtt.critical", "ping.rtt.no_replies"}
	for i, s := range states {
		if s.Status != want[i] {
			t.Errorf("state[%d] status = %v, want %v", i, s.Status, want[i])
		}
		if s.Code != wantCode[i] {
			t.Errorf("state[%d] code = %q, want %q", i, s.Code, wantCode[i])
		}
	}
}
