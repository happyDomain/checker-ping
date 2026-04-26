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
	"encoding/json"
	"fmt"
	"net"
	"time"

	probing "github.com/prometheus-community/pro-bing"

	sdk "git.happydns.org/checker-sdk-go/checker"
	happydns "git.happydns.org/happyDomain/model"
	"git.happydns.org/happyDomain/services/abstract"
)

// Collect performs ICMP ping and returns PingData.
// Addresses are resolved from opts: "addresses" ([]string), "address" (string),
// or "service" (*ServiceMessage of type abstract.Server).
func (p *pingProvider) Collect(ctx context.Context, opts sdk.CheckerOptions) (any, error) {
	addresses, err := resolveAddresses(opts)
	if err != nil {
		return nil, err
	}

	const minCount, maxCount = 1, 20
	count := sdk.GetIntOption(opts, "count", 5)
	if count < minCount || count > maxCount {
		return nil, fmt.Errorf("count must be between %d and %d, got %d", minCount, maxCount, count)
	}

	data := &PingData{}
	var errs []string

	for _, addr := range addresses {
		pinger, err := probing.NewPinger(addr)
		if err != nil {
			errs = append(errs, fmt.Sprintf("failed to create pinger for %s: %v", addr, err))
			continue
		}

		pinger.Count = count
		pinger.Timeout = time.Duration(count)*time.Second + 5*time.Second

		if p.Privileged {
			pinger.SetPrivileged(true)
		}

		if err = pinger.RunWithContext(ctx); err != nil {
			errs = append(errs, fmt.Sprintf("ping failed for %s: %v", addr, err))
			continue
		}

		stats := pinger.Statistics()
		var resolved string
		if ip := pinger.IPAddr(); ip != nil {
			resolved = ip.IP.String()
		}
		data.Targets = append(data.Targets, PingTargetResult{
			Address:    addr,
			ResolvedIP: resolved,
			RTTMin:     float64(stats.MinRtt.Microseconds()) / 1000.0,
			RTTAvg:     float64(stats.AvgRtt.Microseconds()) / 1000.0,
			RTTMax:     float64(stats.MaxRtt.Microseconds()) / 1000.0,
			PacketLoss: stats.PacketLoss,
			Sent:       stats.PacketsSent,
			Received:   stats.PacketsRecv,
		})
	}

	if len(data.Targets) == 0 {
		return nil, fmt.Errorf("all %d ping(s) failed; first error: %s", len(errs), errs[0])
	}

	return data, nil
}

// resolveAddresses extracts target IP addresses from the options.
func resolveAddresses(opts sdk.CheckerOptions) ([]string, error) {
	// Direct addresses (from HTTP server).
	if v, ok := opts["addresses"]; ok {
		switch addrs := v.(type) {
		case []any:
			var result []string
			for _, a := range addrs {
				if s, ok := a.(string); ok && s != "" {
					result = append(result, s)
				}
			}
			if len(result) > 0 {
				return result, nil
			}
		case []string:
			if len(addrs) > 0 {
				return addrs, nil
			}
		}
	}

	// Single address.
	if v, ok := opts["address"]; ok {
		if s, ok := v.(string); ok && s != "" {
			return []string{s}, nil
		}
	}

	// From auto-filled service (plugin provider path or HTTP JSON).
	if svc, ok := sdk.GetOption[happydns.ServiceMessage](opts, "service"); ok {
		if svc.Type != "abstract.Server" {
			return nil, fmt.Errorf("service is %s, expected abstract.Server", svc.Type)
		}
		ips, err := ipsFromService(&svc)
		if err != nil {
			return nil, fmt.Errorf("decode service payload: %w", err)
		}
		if len(ips) > 0 {
			addrs := make([]string, len(ips))
			for i, ip := range ips {
				addrs[i] = ip.String()
			}
			return addrs, nil
		}
		return nil, fmt.Errorf("no IP addresses found in the service")
	}

	return nil, fmt.Errorf("no addresses provided: set 'addresses', 'address', or 'service' in options")
}

func ipsFromService(svc *happydns.ServiceMessage) ([]net.IP, error) {
	var server abstract.Server
	if err := json.Unmarshal(svc.Service, &server); err != nil {
		return nil, err
	}

	var ips []net.IP
	if server.A != nil && len(server.A.A) > 0 {
		ips = append(ips, server.A.A)
	}
	if server.AAAA != nil && len(server.AAAA.AAAA) > 0 {
		ips = append(ips, server.AAAA.AAAA)
	}
	return ips, nil
}
