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

// Rule returns a new ping check rule for local evaluation.
func Rule() sdk.CheckRule {
	return &pingRule{}
}

type pingRule struct{}

func (r *pingRule) Name() string { return "ping_check" }
func (r *pingRule) Description() string {
	return "Checks ICMP ping reachability, round-trip time, and packet loss"
}

func (r *pingRule) ValidateOptions(opts sdk.CheckerOptions) error {
	warningRTT := float64(100)
	criticalRTT := float64(500)
	warningPacketLoss := float64(10)
	criticalPacketLoss := float64(50)

	if v, ok := opts["warningRTT"]; ok {
		d, ok := v.(float64)
		if !ok {
			return fmt.Errorf("warningRTT must be a number")
		}
		if d <= 0 {
			return fmt.Errorf("warningRTT must be positive")
		}
		warningRTT = d
	}
	if v, ok := opts["criticalRTT"]; ok {
		d, ok := v.(float64)
		if !ok {
			return fmt.Errorf("criticalRTT must be a number")
		}
		if d <= 0 {
			return fmt.Errorf("criticalRTT must be positive")
		}
		criticalRTT = d
	}
	if v, ok := opts["warningPacketLoss"]; ok {
		d, ok := v.(float64)
		if !ok {
			return fmt.Errorf("warningPacketLoss must be a number")
		}
		if d < 0 || d > 100 {
			return fmt.Errorf("warningPacketLoss must be between 0 and 100")
		}
		warningPacketLoss = d
	}
	if v, ok := opts["criticalPacketLoss"]; ok {
		d, ok := v.(float64)
		if !ok {
			return fmt.Errorf("criticalPacketLoss must be a number")
		}
		if d < 0 || d > 100 {
			return fmt.Errorf("criticalPacketLoss must be between 0 and 100")
		}
		criticalPacketLoss = d
	}
	if v, ok := opts["count"]; ok {
		d, ok := v.(float64)
		if !ok {
			return fmt.Errorf("count must be a number")
		}
		if d < 1 || d > 20 {
			return fmt.Errorf("count must be between 1 and 20")
		}
	}

	if criticalRTT <= warningRTT {
		return fmt.Errorf("criticalRTT (%v) must be greater than warningRTT (%v)", criticalRTT, warningRTT)
	}
	if criticalPacketLoss <= warningPacketLoss {
		return fmt.Errorf("criticalPacketLoss (%v) must be greater than warningPacketLoss (%v)", criticalPacketLoss, warningPacketLoss)
	}

	return nil
}

func (r *pingRule) Evaluate(ctx context.Context, obs sdk.ObservationGetter, opts sdk.CheckerOptions) []sdk.CheckState {
	var data PingData
	if err := obs.Get(ctx, ObservationKeyPing, &data); err != nil {
		return []sdk.CheckState{{
			Status:  sdk.StatusError,
			Message: fmt.Sprintf("Failed to get ping data: %v", err),
			Code:    "ping_error",
		}}
	}

	warningRTT := sdk.GetFloatOption(opts, "warningRTT", 100)
	criticalRTT := sdk.GetFloatOption(opts, "criticalRTT", 500)
	warningPacketLoss := sdk.GetFloatOption(opts, "warningPacketLoss", 10)
	criticalPacketLoss := sdk.GetFloatOption(opts, "criticalPacketLoss", 50)

	results := Evaluate(&data, warningRTT, criticalRTT, warningPacketLoss, criticalPacketLoss)
	if len(results) == 0 {
		return []sdk.CheckState{{
			Status:  sdk.StatusInfo,
			Message: "No targets to ping",
			Code:    "ping_no_targets",
		}}
	}

	targetByAddr := make(map[string]PingTargetResult, len(data.Targets))
	for _, t := range data.Targets {
		targetByAddr[t.Address] = t
	}

	out := make([]sdk.CheckState, 0, len(results))
	for _, r := range results {
		var status sdk.Status
		switch r.Status {
		case StatusOK:
			status = sdk.StatusOK
		case StatusWarn:
			status = sdk.StatusWarn
		case StatusCrit:
			status = sdk.StatusCrit
		default:
			status = sdk.StatusUnknown
		}

		state := sdk.CheckState{
			Status:  status,
			Subject: r.Address,
			Message: r.Message,
			Code:    r.Code,
		}
		if t, ok := targetByAddr[r.Address]; ok {
			state.Meta = map[string]any{"target": t}
		}
		out = append(out, state)
	}
	return out
}
