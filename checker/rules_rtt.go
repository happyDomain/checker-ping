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

// rttRule evaluates the average round-trip time of each target against
// the configured warning/critical thresholds. Unreachable targets (no
// reply at all) are skipped; the reachability rule handles those.
type rttRule struct{}

func (r *rttRule) Name() string { return "ping.rtt" }
func (r *rttRule) Description() string {
	return "Flags targets whose average round-trip time crosses the warning or critical threshold."
}

func (r *rttRule) ValidateOptions(opts sdk.CheckerOptions) error {
	warn := sdk.GetFloatOption(opts, "warningRTT", 100)
	crit := sdk.GetFloatOption(opts, "criticalRTT", 500)
	if warn <= 0 {
		return fmt.Errorf("warningRTT must be positive")
	}
	if crit <= 0 {
		return fmt.Errorf("criticalRTT must be positive")
	}
	if crit <= warn {
		return fmt.Errorf("criticalRTT (%v) must be greater than warningRTT (%v)", crit, warn)
	}
	return nil
}

func (r *rttRule) Evaluate(ctx context.Context, obs sdk.ObservationGetter, opts sdk.CheckerOptions) []sdk.CheckState {
	data, errSt := loadPingData(ctx, obs)
	if errSt != nil {
		return []sdk.CheckState{*errSt}
	}
	if len(data.Targets) == 0 {
		return []sdk.CheckState{noTargetsState("ping.rtt.no_targets")}
	}

	warn := sdk.GetFloatOption(opts, "warningRTT", 100)
	crit := sdk.GetFloatOption(opts, "criticalRTT", 500)

	out := make([]sdk.CheckState, 0, len(data.Targets))
	for _, t := range data.Targets {
		if t.Received == 0 {
			out = append(out, sdk.CheckState{
				Status:  sdk.StatusUnknown,
				Subject: t.Address,
				Message: "RTT not measurable (no replies).",
				Code:    "ping.rtt.no_replies",
				Meta:    map[string]any{"target": t},
			})
			continue
		}

		state := sdk.CheckState{
			Subject: t.Address,
			Message: fmt.Sprintf("Average RTT %.1fms (warn=%.0fms, crit=%.0fms).", t.RTTAvg, warn, crit),
			Meta:    map[string]any{"target": t},
		}
		switch {
		case t.RTTAvg >= crit:
			state.Status = sdk.StatusCrit
			state.Code = "ping.rtt.critical"
		case t.RTTAvg >= warn:
			state.Status = sdk.StatusWarn
			state.Code = "ping.rtt.warning"
		default:
			state.Status = sdk.StatusOK
			state.Code = "ping.rtt.ok"
		}
		out = append(out, state)
	}
	return out
}
