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

package main

import (
	"flag"
	"log"

	ping "git.happydns.org/checker-ping/checker"
	"git.happydns.org/checker-sdk-go/checker/server"
)

var (
	listenAddr = flag.String("listen", ":8080", "HTTP listen address")
	privileged = flag.Bool("privileged", false, "Use privileged ICMP (requires CAP_NET_RAW or root)")
)

// Version is the standalone binary's version. It defaults to "custom-build"
// and is meant to be overridden by the CI at link time:
//
//	go build -ldflags "-X main.Version=1.2.3" .
var Version = "custom-build"

func main() {
	flag.Parse()

	// Propagate the binary version to the checker package so it shows up in
	// CheckerDefinition.Version.
	ping.Version = Version

	srv := server.New(ping.ProviderWithPrivileged(*privileged))
	if err := srv.ListenAndServe(*listenAddr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
