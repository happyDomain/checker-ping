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
	"errors"
	"net/http"
	"strconv"
	"strings"

	sdk "git.happydns.org/checker-sdk-go/checker"
)

// RenderForm implements sdk.CheckerInteractive. It exposes a minimal form
// letting a human submit one or more ping targets (hostnames or IPs) along
// with the usual threshold knobs.
func (p *pingProvider) RenderForm() []sdk.CheckerOptionField {
	return []sdk.CheckerOptionField{
		{
			Id:          "addresses",
			Type:        "string",
			Label:       "Targets",
			Placeholder: "example.com, 192.0.2.1",
			Description: "Comma- or newline-separated list of hostnames or IP addresses.",
			Required:    true,
		},
		{
			Id:      "count",
			Type:    "uint",
			Label:   "Number of pings to send",
			Default: float64(5),
		},
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
	}
}

// ParseForm implements sdk.CheckerInteractive. It converts the HTML form
// inputs into a CheckerOptions that Collect can consume directly — pinging
// resolves hostnames on its own, so no extra lookups are needed here.
func (p *pingProvider) ParseForm(r *http.Request) (sdk.CheckerOptions, error) {
	raw := strings.TrimSpace(r.FormValue("addresses"))
	if raw == "" {
		return nil, errors.New("at least one target is required")
	}

	var addresses []string
	for _, part := range strings.FieldsFunc(raw, func(c rune) bool {
		return c == ',' || c == '\n' || c == '\r' || c == ' ' || c == '\t' || c == ';'
	}) {
		if part = strings.TrimSpace(part); part != "" {
			addresses = append(addresses, part)
		}
	}
	if len(addresses) == 0 {
		return nil, errors.New("at least one target is required")
	}

	opts := sdk.CheckerOptions{"addresses": addresses}
	for _, k := range []string{"count", "warningRTT", "criticalRTT", "warningPacketLoss", "criticalPacketLoss"} {
		if v := strings.TrimSpace(r.FormValue(k)); v != "" {
			if n, err := strconv.ParseFloat(v, 64); err == nil {
				opts[k] = n
			}
		}
	}
	return opts, nil
}
