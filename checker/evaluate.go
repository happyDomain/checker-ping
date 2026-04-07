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
	"fmt"
	"strings"
	"time"
)

// Metrics extracts time-series metrics from ping data.
func Metrics(data *PingData, collectedAt time.Time) []Metric {
	var metrics []Metric
	for _, t := range data.Targets {
		labels := map[string]string{"address": t.Address}
		metrics = append(metrics,
			Metric{Name: "ping_rtt_avg", Value: t.RTTAvg, Unit: "ms", Labels: labels, Timestamp: collectedAt},
			Metric{Name: "ping_rtt_min", Value: t.RTTMin, Unit: "ms", Labels: labels, Timestamp: collectedAt},
			Metric{Name: "ping_rtt_max", Value: t.RTTMax, Unit: "ms", Labels: labels, Timestamp: collectedAt},
			Metric{Name: "ping_packet_loss", Value: t.PacketLoss, Unit: "%", Labels: labels, Timestamp: collectedAt},
			Metric{Name: "ping_packets_sent", Value: float64(t.Sent), Unit: "count", Labels: labels, Timestamp: collectedAt},
			Metric{Name: "ping_packets_received", Value: float64(t.Received), Unit: "count", Labels: labels, Timestamp: collectedAt},
		)
	}
	return metrics
}

// Metric represents a single time-series metric.
type Metric struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Unit      string            `json:"unit,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// EvaluateResult holds the evaluation outcome.
type EvaluateResult struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

const (
	StatusOK   = 1
	StatusWarn = 3
	StatusCrit = 4
)

// Evaluate checks the ping data against the given thresholds.
func Evaluate(data *PingData, warningRTT, criticalRTT, warningPacketLoss, criticalPacketLoss float64) EvaluateResult {
	overallStatus := StatusOK
	var summaryParts []string

	for _, target := range data.Targets {
		if target.PacketLoss >= criticalPacketLoss || target.RTTAvg >= criticalRTT {
			overallStatus = StatusCrit
		} else if (target.PacketLoss >= warningPacketLoss || target.RTTAvg >= warningRTT) && overallStatus < StatusWarn {
			overallStatus = StatusWarn
		}

		summaryParts = append(summaryParts, fmt.Sprintf("%s: %.1fms avg, %.0f%% loss", target.Address, target.RTTAvg, target.PacketLoss))
	}

	return EvaluateResult{
		Status:  overallStatus,
		Message: strings.Join(summaryParts, " | "),
		Code:    "ping_result",
	}
}
