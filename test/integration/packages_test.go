package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"

	"github.com/git-pkg/gpkg/pkg/constants"
	"github.com/git-pkg/gpkg/pkg/packages"
)

// -----------------------------------------------------------------------------
// Tests - Vars
// -----------------------------------------------------------------------------

const unusableDir = "/tmp/should-not-exist"

var (
	testDir  string
	testPkgs = map[string]string{
		"cli":      "cli",
		"neovim":   "neovim",
		"gohugoio": "hugo",
	}
	testVers = map[string]semver.Version{
		"cli":    semver.MustParse("0.11.0"),
		"hugo":   semver.MustParse("0.74.3"),
		"neovim": constants.LatestVersion,
	}
)

// -----------------------------------------------------------------------------
// Tests - Setup
// -----------------------------------------------------------------------------

func TestSetup(t *testing.T) {
	// setup the temporary directory we we use for testing
	var err error
	testDir, err = ioutil.TempDir(os.TempDir(), "gpkg-tests-")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("INFO: created temporary directory %s for testing", testDir)
}

// -----------------------------------------------------------------------------
// Tests - Package - Install/Remove
// -----------------------------------------------------------------------------

func TestPackageInstall(t *testing.T) {
	for owner, repo := range testPkgs {
		p, err := packages.NewPackageFromGithub(owner, repo)
		assert.NoError(t, err)

		v := testVers[repo]
		assert.NoError(t, p.Install(testDir, &v))
	}
}

func TestPackageReinstall(t *testing.T) {
	p, err := packages.NewPackageFromGithub("cli", "cli")
	assert.NoError(t, err)
	assert.Error(t, p.Install(testDir, &constants.LatestVersion))
}

func TestPackageUninstall(t *testing.T) {
	if os.Getenv("DEBUG") != "" {
		t.Skip()
	}

	for owner, repo := range testPkgs {
		p, err := packages.NewPackageFromGithub(owner, repo)
		assert.NoError(t, err)
		assert.NoError(t, p.Uninstall(testDir))
	}
}

// -----------------------------------------------------------------------------
// Tests - Cleanup
// -----------------------------------------------------------------------------

func TestCleanup(t *testing.T) {
	if os.Getenv("DEBUG") != "" {
		t.Skip()
	}

	if err := os.RemoveAll(testDir); err != nil {
		t.Fatal(fmt.Errorf("could not remove tempdir created for testing: %s", err))
	}

	t.Logf("INFO: removed temporary testing directory %s", testDir)
}
