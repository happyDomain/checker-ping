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

// Command plugin is the happyDomain plugin entrypoint for the ping checker.
//
// It is built as a Go plugin (`go build -buildmode=plugin`) and loaded at
// runtime by happyDomain.
package main

import (
	ping "git.happydns.org/checker-ping/checker"
	sdk "git.happydns.org/checker-sdk-go/checker"
)

// Version is the plugin's version. It defaults to "custom-build" and is
// meant to be overridden by the CI at link time:
//
//	go build -buildmode=plugin -ldflags "-X main.Version=1.2.3" -o checker-ping.so ./plugin
var Version = "custom-build"

// NewCheckerPlugin is the symbol resolved by happyDomain when loading the
// .so file. It returns the checker definition and the observation provider
// that the host will register in its global registries.
func NewCheckerPlugin() (*sdk.CheckerDefinition, sdk.ObservationProvider, error) {
	// Propagate the plugin's version to the checker package so it shows up
	// in CheckerDefinition.Version.
	ping.Version = Version
	prvd := ping.Provider()
	return prvd.(sdk.CheckerDefinitionProvider).Definition(), prvd, nil
}
