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

	sdk "git.happydns.org/checker-sdk-go/checker"
)

// reachabilityRule verifies each target replied to at least one probe.
// A target with 100% packet loss is considered unreachable.
type reachabilityRule struct{}

func (r *reachabilityRule) Name() string { return "ping.reachable" }
func (r *reachabilityRule) Description() string {
	return "Verifies every target replied to at least one ICMP probe."
}

func (r *reachabilityRule) Evaluate(ctx context.Context, obs sdk.ObservationGetter, _ sdk.CheckerOptions) []sdk.CheckState {
	data, errSt := loadPingData(ctx, obs)
	if errSt != nil {
		return []sdk.CheckState{*errSt}
	}
	if len(data.Targets) == 0 {
		return []sdk.CheckState{noTargetsState("ping.reachable.no_targets")}
	}

	out := make([]sdk.CheckState, 0, len(data.Targets))
	for _, t := range data.Targets {
		if t.Received == 0 {
			out = append(out, sdk.CheckState{
				Status:  sdk.StatusCrit,
				Subject: t.Address,
				Message: fmt.Sprintf("Target unreachable (0/%d replies).", t.Sent),
				Code:    "ping.reachable.unreachable",
				Meta:    map[string]any{"target": t},
			})
		} else {
			out = append(out, sdk.CheckState{
				Status:  sdk.StatusOK,
				Subject: t.Address,
				Message: fmt.Sprintf("Target reachable (%d/%d replies).", t.Received, t.Sent),
				Code:    "ping.reachable.ok",
				Meta:    map[string]any{"target": t},
			})
		}
	}
	return out
}
