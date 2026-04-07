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

// ObservationKeyPing is the observation key for ICMP ping data.
const ObservationKeyPing = "ping"

// PingData holds the collected ping results for all targets.
type PingData struct {
	Targets []PingTargetResult `json:"targets"`
}

// PingTargetResult contains the ping statistics for a single IP address.
type PingTargetResult struct {
	Address    string  `json:"address"`
	RTTMin     float64 `json:"rtt_min"`
	RTTAvg     float64 `json:"rtt_avg"`
	RTTMax     float64 `json:"rtt_max"`
	PacketLoss float64 `json:"packet_loss"`
	Sent       int     `json:"sent"`
	Received   int     `json:"received"`
}
