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
	"net"
	"strings"
	"testing"

	sdk "git.happydns.org/checker-sdk-go/checker"
	happydns "git.happydns.org/happyDomain/model"
	"git.happydns.org/happyDomain/services/abstract"
	"github.com/miekg/dns"
)

func TestResolveAddressesAddressesSlice(t *testing.T) {
	cases := []sdk.CheckerOptions{
		{"addresses": []string{"1.1.1.1", "2.2.2.2"}},
		{"addresses": []any{"1.1.1.1", "2.2.2.2"}},
		{"addresses": []any{"1.1.1.1", "", "2.2.2.2"}}, // empties dropped
	}
	for i, opts := range cases {
		got, err := resolveAddresses(opts)
		if err != nil {
			t.Fatalf("[%d] err: %v", i, err)
		}
		if len(got) != 2 || got[0] != "1.1.1.1" || got[1] != "2.2.2.2" {
			t.Errorf("[%d] got %v", i, got)
		}
	}
}

func TestResolveAddressesSingle(t *testing.T) {
	got, err := resolveAddresses(sdk.CheckerOptions{"address": "8.8.8.8"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "8.8.8.8" {
		t.Errorf("got %v", got)
	}
}

func TestResolveAddressesMissing(t *testing.T) {
	if _, err := resolveAddresses(sdk.CheckerOptions{}); err == nil {
		t.Error("expected error for missing addresses")
	}
}

func TestResolveAddressesEmptyStringIgnored(t *testing.T) {
	if _, err := resolveAddresses(sdk.CheckerOptions{"address": ""}); err == nil {
		t.Error("expected error when address is empty string")
	}
}

func TestResolveAddressesWrongServiceType(t *testing.T) {
	svc := happydns.ServiceMessage{
		ServiceMeta: happydns.ServiceMeta{Type: "abstract.NotAServer"},
		Service:     json.RawMessage(`{}`),
	}
	_, err := resolveAddresses(sdk.CheckerOptions{"service": svc})
	if err == nil || !strings.Contains(err.Error(), "expected abstract.Server") {
		t.Errorf("got err=%v", err)
	}
}

func TestIpsFromService(t *testing.T) {
	srv := abstract.Server{
		A:    &dns.A{A: net.ParseIP("192.0.2.1").To4()},
		AAAA: &dns.AAAA{AAAA: net.ParseIP("2001:db8::1")},
	}
	raw, err := json.Marshal(srv)
	if err != nil {
		t.Fatal(err)
	}
	msg := &happydns.ServiceMessage{Service: raw}
	ips, err := ipsFromService(msg)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(ips) != 2 {
		t.Fatalf("got %d ips, want 2: %v", len(ips), ips)
	}
}

func TestIpsFromServiceMalformed(t *testing.T) {
	msg := &happydns.ServiceMessage{Service: json.RawMessage(`{not json`)}
	if _, err := ipsFromService(msg); err == nil {
		t.Error("expected error on malformed payload")
	}
}

func TestIpsFromServiceEmpty(t *testing.T) {
	msg := &happydns.ServiceMessage{Service: json.RawMessage(`{}`)}
	ips, err := ipsFromService(msg)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(ips) != 0 {
		t.Errorf("expected 0 ips, got %v", ips)
	}
}
