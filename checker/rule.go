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
	"fmt"
	"net"

	sdk "git.happydns.org/checker-sdk-go/checker"
)

// Rules returns the full list of CheckRules exposed by the ping checker.
// Each rule covers a single concern so callers see at a glance which
// aspects passed and which did not, instead of sharing a single monolithic
// rule result.
func Rules() []sdk.CheckRule {
	return []sdk.CheckRule{
		&reachabilityRule{},
		&packetLossRule{},
		&rttRule{},
		&ipv6ReachabilityRule{},
	}
}

// validateThresholdPair checks that warn and crit are within [min, max] and
// that crit is strictly greater than warn. The names are used in error
// messages so callers get diagnostics naming their actual options.
func validateThresholdPair(warnName, critName string, warn, crit, min, max float64) error {
	if warn < min || warn > max {
		return fmt.Errorf("%s must be between %v and %v", warnName, min, max)
	}
	if crit < min || crit > max {
		return fmt.Errorf("%s must be between %v and %v", critName, min, max)
	}
	if crit <= warn {
		return fmt.Errorf("%s (%v) must be greater than %s (%v)", critName, crit, warnName, warn)
	}
	return nil
}

// loadPingData fetches the ping observation. On error, returns a CheckState
// the caller should emit to short-circuit its rule.
func loadPingData(ctx context.Context, obs sdk.ObservationGetter) (*PingData, *sdk.CheckState) {
	var data PingData
	if err := obs.Get(ctx, ObservationKeyPing, &data); err != nil {
		return nil, &sdk.CheckState{
			Status:  sdk.StatusError,
			Message: fmt.Sprintf("failed to load ping observation: %v", err),
			Code:    "ping.observation_error",
		}
	}
	return &data, nil
}

// noTargetsState is returned when the observation has no targets at all.
func noTargetsState(code string) sdk.CheckState {
	return sdk.CheckState{
		Status:  sdk.StatusUnknown,
		Message: "No targets to ping",
		Code:    code,
	}
}

// isIPv6 reports whether addr parses as an IPv6 address (excluding
// IPv4-mapped representations).
func isIPv6(addr string) bool {
	ip := net.ParseIP(addr)
	if ip == nil {
		return false
	}
	return ip.To4() == nil
}
