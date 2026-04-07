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
	"encoding/json"
	"time"

	happydns "git.happydns.org/checker-sdk-go/checker"
)

// Provider returns a new ping observation provider for local execution.
func Provider() happydns.ObservationProvider {
	return &pingProvider{}
}

// ProviderWithPrivileged returns a provider with privileged ICMP mode enabled.
func ProviderWithPrivileged(privileged bool) happydns.ObservationProvider {
	return &pingProvider{Privileged: privileged}
}

type pingProvider struct {
	// Privileged controls whether raw ICMP sockets are used (requires CAP_NET_RAW or root).
	Privileged bool
}

func (p *pingProvider) Key() happydns.ObservationKey {
	return ObservationKeyPing
}

// Definition implements happydns.CheckerDefinitionProvider.
func (p *pingProvider) Definition() *happydns.CheckerDefinition {
	return Definition()
}

// ExtractMetrics implements happydns.CheckerMetricsReporter.
func (p *pingProvider) ExtractMetrics(raw json.RawMessage, collectedAt time.Time) ([]happydns.CheckMetric, error) {
	var data PingData
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	metrics := Metrics(&data, collectedAt)
	result := make([]happydns.CheckMetric, len(metrics))
	for i, m := range metrics {
		result[i] = happydns.CheckMetric{
			Name:      m.Name,
			Value:     m.Value,
			Unit:      m.Unit,
			Labels:    m.Labels,
			Timestamp: m.Timestamp,
		}
	}
	return result, nil
}
