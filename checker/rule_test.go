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
	"encoding/json"
	"errors"
	"testing"

	sdk "git.happydns.org/checker-sdk-go/checker"
)

// stubObs implements sdk.ObservationGetter for tests. If err is non-nil, Get
// returns it; otherwise it JSON-roundtrips data into dest so callers see the
// same shape they would get over HTTP.
type stubObs struct {
	data any
	err  error
}

func (s stubObs) Get(_ context.Context, _ sdk.ObservationKey, dest any) error {
	if s.err != nil {
		return s.err
	}
	raw, err := json.Marshal(s.data)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, dest)
}

func (s stubObs) GetRelated(_ context.Context, _ sdk.ObservationKey) ([]sdk.RelatedObservation, error) {
	return nil, nil
}

func obsWith(targets ...PingTargetResult) stubObs {
	return stubObs{data: PingData{Targets: targets}}
}

func TestIsIPv6(t *testing.T) {
	cases := []struct {
		addr string
		want bool
	}{
		{"::1", true},
		{"2001:db8::1", true},
		{"127.0.0.1", false},
		{"::ffff:192.0.2.1", false}, // IPv4-mapped is treated as IPv4
		{"example.com", false},
		{"", false},
		{"not-an-ip", false},
	}
	for _, c := range cases {
		if got := isIPv6(c.addr); got != c.want {
			t.Errorf("isIPv6(%q) = %v, want %v", c.addr, got, c.want)
		}
	}
}

func TestLoadPingDataError(t *testing.T) {
	_, st := loadPingData(context.Background(), stubObs{err: errors.New("boom")})
	if st == nil {
		t.Fatal("expected error CheckState, got nil")
	}
	if st.Status != sdk.StatusError {
		t.Errorf("status = %v, want StatusError", st.Status)
	}
	if st.Code != "ping.observation_error" {
		t.Errorf("code = %q, want ping.observation_error", st.Code)
	}
}

func TestLoadPingDataOK(t *testing.T) {
	d, st := loadPingData(context.Background(), obsWith(PingTargetResult{Address: "1.1.1.1"}))
	if st != nil {
		t.Fatalf("unexpected error state: %+v", st)
	}
	if len(d.Targets) != 1 || d.Targets[0].Address != "1.1.1.1" {
		t.Errorf("unexpected data: %+v", d)
	}
}

func TestRulesContainsAll(t *testing.T) {
	names := map[string]bool{}
	for _, r := range Rules() {
		names[r.Name()] = true
	}
	for _, want := range []string{"ping.reachable", "ping.packet_loss", "ping.rtt", "ping.ipv6_reachable"} {
		if !names[want] {
			t.Errorf("Rules() missing %q", want)
		}
	}
}
