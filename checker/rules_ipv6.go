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

	sdk "git.happydns.org/checker-sdk-go/checker"
)

// ipv6ReachabilityRule verifies at least one IPv6 target replied to a probe.
// The rule is skipped (StatusUnknown) when no IPv6 target was pinged, which
// is expected for IPv4-only hosts.
type ipv6ReachabilityRule struct{}

func (r *ipv6ReachabilityRule) Name() string { return "ping.ipv6_reachable" }
func (r *ipv6ReachabilityRule) Description() string {
	return "Verifies that at least one IPv6 target replied to an ICMP probe."
}

func (r *ipv6ReachabilityRule) Evaluate(ctx context.Context, obs sdk.ObservationGetter, _ sdk.CheckerOptions) []sdk.CheckState {
	data, errSt := loadPingData(ctx, obs)
	if errSt != nil {
		return []sdk.CheckState{*errSt}
	}

	var ipv6Total, ipv6Reachable int
	for _, t := range data.Targets {
		probed := t.ResolvedIP
		if probed == "" {
			probed = t.Address
		}
		if !isIPv6(probed) {
			continue
		}
		ipv6Total++
		if t.Received > 0 {
			ipv6Reachable++
		}
	}

	switch {
	case ipv6Total == 0:
		return []sdk.CheckState{{
			Status:  sdk.StatusUnknown,
			Message: "No IPv6 target pinged.",
			Code:    "ping.ipv6_reachable.skipped",
		}}
	case ipv6Reachable == 0:
		return []sdk.CheckState{{
			Status:  sdk.StatusWarn,
			Message: "No IPv6 target replied to ICMP probes.",
			Code:    "ping.ipv6_reachable.unreachable",
		}}
	default:
		return []sdk.CheckState{{
			Status:  sdk.StatusOK,
			Message: "At least one IPv6 target is reachable.",
			Code:    "ping.ipv6_reachable.ok",
		}}
	}
}
