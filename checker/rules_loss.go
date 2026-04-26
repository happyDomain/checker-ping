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

// packetLossRule evaluates the observed packet-loss ratio of each target
// against the configured warning/critical thresholds.
type packetLossRule struct{}

func (r *packetLossRule) Name() string { return "ping.packet_loss" }
func (r *packetLossRule) Description() string {
	return "Flags targets whose packet-loss ratio crosses the warning or critical threshold."
}

func (r *packetLossRule) ValidateOptions(opts sdk.CheckerOptions) error {
	warn := sdk.GetFloatOption(opts, "warningPacketLoss", 10)
	crit := sdk.GetFloatOption(opts, "criticalPacketLoss", 50)
	if warn < 0 || warn > 100 {
		return fmt.Errorf("warningPacketLoss must be between 0 and 100")
	}
	if crit < 0 || crit > 100 {
		return fmt.Errorf("criticalPacketLoss must be between 0 and 100")
	}
	if crit <= warn {
		return fmt.Errorf("criticalPacketLoss (%v) must be greater than warningPacketLoss (%v)", crit, warn)
	}
	return nil
}

func (r *packetLossRule) Evaluate(ctx context.Context, obs sdk.ObservationGetter, opts sdk.CheckerOptions) []sdk.CheckState {
	data, errSt := loadPingData(ctx, obs)
	if errSt != nil {
		return []sdk.CheckState{*errSt}
	}
	if len(data.Targets) == 0 {
		return []sdk.CheckState{noTargetsState("ping.packet_loss.no_targets")}
	}

	warn := sdk.GetFloatOption(opts, "warningPacketLoss", 10)
	crit := sdk.GetFloatOption(opts, "criticalPacketLoss", 50)

	out := make([]sdk.CheckState, 0, len(data.Targets))
	for _, t := range data.Targets {
		state := sdk.CheckState{
			Subject: t.Address,
			Message: fmt.Sprintf("Packet loss %.0f%% (warn=%.0f%%, crit=%.0f%%).", t.PacketLoss, warn, crit),
			Meta:    map[string]any{"target": t},
		}
		switch {
		case t.PacketLoss >= crit:
			state.Status = sdk.StatusCrit
			state.Code = "ping.packet_loss.critical"
		case t.PacketLoss >= warn:
			state.Status = sdk.StatusWarn
			state.Code = "ping.packet_loss.warning"
		default:
			state.Status = sdk.StatusOK
			state.Code = "ping.packet_loss.ok"
		}
		out = append(out, state)
	}
	return out
}
