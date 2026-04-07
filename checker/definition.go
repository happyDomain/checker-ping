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
	"time"

	happydns "git.happydns.org/checker-sdk-go/checker"
)

// Version is the checker version reported in CheckerDefinition.Version.
//
// It defaults to "built-in", which is appropriate when the checker package is
// imported directly (built-in or plugin mode). Standalone binaries (like
// main.go) should override this from their own Version variable at the start
// of main(), which makes it easy for CI to inject a version with a single
// -ldflags "-X main.Version=..." flag instead of targeting the nested
// package path.
var Version = "built-in"

// Definition returns the CheckerDefinition for the ping checker.
func Definition() *happydns.CheckerDefinition {
	return &happydns.CheckerDefinition{
		ID:      "ping",
		Name:    "Ping (ICMP)",
		Version: Version,
		Availability: happydns.CheckerAvailability{
			ApplyToService:  true,
			LimitToServices: []string{"abstract.Server"},
		},
		ObservationKeys: []happydns.ObservationKey{ObservationKeyPing},
		Options: happydns.CheckerOptionsDocumentation{
			UserOpts: []happydns.CheckerOptionDocumentation{
				{
					Id:      "warningRTT",
					Type:    "number",
					Label:   "Warning RTT threshold (ms)",
					Default: float64(100),
				},
				{
					Id:      "criticalRTT",
					Type:    "number",
					Label:   "Critical RTT threshold (ms)",
					Default: float64(500),
				},
				{
					Id:      "warningPacketLoss",
					Type:    "number",
					Label:   "Warning packet loss threshold (%)",
					Default: float64(10),
				},
				{
					Id:      "criticalPacketLoss",
					Type:    "number",
					Label:   "Critical packet loss threshold (%)",
					Default: float64(50),
				},
				{
					Id:      "count",
					Type:    "uint",
					Label:   "Number of pings to send",
					Default: float64(5),
				},
			},
			ServiceOpts: []happydns.CheckerOptionDocumentation{
				{
					Id:       "service",
					Label:    "Service",
					AutoFill: happydns.AutoFillService,
				},
			},
		},
		Rules: []happydns.CheckRule{
			Rule(),
		},
		Interval: &happydns.CheckIntervalSpec{
			Min:     1 * time.Minute,
			Max:     1 * time.Hour,
			Default: 5 * time.Minute,
		},
		HasMetrics: true,
	}
}
