// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package version

import "github.com/hashicorp/packer-plugin-sdk/version"

var (
	// Version is set via ldflags
	Version string

	// Build version is set via ldflags
	Build string

	// Build date is set via ldflags
	Date string

	// VersionPrerelease is A pre-release marker for the Version. If this is ""
	// (empty string) then it means that it is a final release. Otherwise, this
	// is a pre-release such as "dev" (in development), "beta", "rc1", etc.
	VersionPrerelease = "dev"

	// PluginVersion is used by the plugin set to allow Packer to recognize
	// what version this plugin is.
	PluginVersion = version.InitializePluginVersion(Version, VersionPrerelease)
)
