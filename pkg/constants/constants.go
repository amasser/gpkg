package constants

import "github.com/blang/semver"

// -----------------------------------------------------------------------------
// Public Vars - Defaults
// -----------------------------------------------------------------------------

var (
	// DefaultPackageRepository indicates the URL of the default repository used to fetch packages
	DefaultPackageRepository = "https://packages.git-pkg.dev"

	// DefaultContentType is the mime type expected for downloaded artifacts
	DefaultContentType = "application/octet-stream"

	// DefaultBaseDir is the default directory where packages will be installed and maintained
	DefaultBaseDir = ".local/gpkg"

	// LatestVersion can be provided as the desired package version when you just want the latest release/prerelease/tag
	LatestVersion = semver.MustParse("0.0.0")
)

// -----------------------------------------------------------------------------
// Public Vars - Packages
// -----------------------------------------------------------------------------

const (
	// TarMethod indicates that the artifact is in TAR archive format
	TarMethod = "tar"

	// GunzipCompression indicates that the artifact is compressed with gunzip
	GunzipCompression = "gz"
)

// -----------------------------------------------------------------------------
// Public Vars - Statuses
// -----------------------------------------------------------------------------

const (
	// ReadyStatus indicates that all actions on a resource are resolved
	ReadyStatus = "ready"

	// FailedStatus indicates that some actions on a resource have failed and it can't be automatically recovered
	FailedStatus = "failed"

	// ProcessingStatus indicates that some actions are currently underway on a resource and have not yet resolved
	ProcessingStatus = "processing"
)

// -----------------------------------------------------------------------------
// Public Vars - Release Targets
// -----------------------------------------------------------------------------

const (
	// Releases indicates that tags which have an associated release are valid for app installation and upgrade
	Releases = "releases"

	// PreReleases indicates that tags which are marked as pre-release are valid for app installation and upgrade
	PreReleases = "pre-releases"

	// Tags indicates that tags without releases associated with them are valid for app installation and upgrade
	Tags = "tags"
)
