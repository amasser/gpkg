package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/blang/semver"
	ghc "github.com/google/go-github/v32/github"

	"github.com/git-pkg/gpkg/pkg/constants"
)

// -----------------------------------------------------------------------------
// Github - Public Functions
// -----------------------------------------------------------------------------

// GetOwnerRepo provides the owner and repo names given a Github URL
func GetOwnerRepo(u *url.URL) (owner, repo string, err error) {
	path := strings.Split(u.Path, "/")
	if len(path) < 2 {
		err = fmt.Errorf("invalid repository %s: must follow the pattern https://github.com/<owner>/<repo>", u)
		return
	}

	owner, repo = path[1], path[2]
	return
}

// GetReleaseArtifactURL locates the appropriate release based on repo and version and produces a *url.URL to the artifact for retrieval
func GetReleaseArtifactURL(basedir, owner, repo string, v *semver.Version, re *regexp.Regexp) (*url.URL, error) {
	// if no version was provided, simply use the latest version
	if v == nil {
		v = &constants.LatestVersion
	}

	// find the latest release if requested
	if v.String() == constants.LatestVersion.String() {
		release, err := latestRelease(owner, repo)
		if err != nil {
			return nil, err
		}
		return artifactURL(release, re)
	}

	// retrieve a specific release
	release, err := releaseByVersion(owner, repo, v)
	if err != nil {
		return nil, err
	}

	return artifactURL(release, re)
}

// -----------------------------------------------------------------------------
// Github - Private Functions
// -----------------------------------------------------------------------------

func artifactURL(release *ghc.RepositoryRelease, re *regexp.Regexp) (*url.URL, error) {
	found := []*ghc.ReleaseAsset{}
	for _, asset := range release.Assets {
		if re.MatchString(asset.GetName()) {
			found = append(found, asset)
		}
	}

	if len(found) < 1 {
		return nil, fmt.Errorf("no release artifact could be found for package %s", release)
	}

	if len(found) > 1 {
		return nil, fmt.Errorf("package for %s matches multiple release artifacts (report this as a bug for that package!)", release)
	}

	asset := found[0]
	artifactURL, err := url.Parse(asset.GetBrowserDownloadURL())
	if err != nil {
		return nil, err
	}

	return artifactURL, nil
}

func latestRelease(owner, repo string) (*ghc.RepositoryRelease, error) {
	c := ghc.NewClient(nil)
	releases, resp, err := c.Repositories.ListReleases(context.TODO(), owner, repo, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}

	if len(releases) < 1 {
		return nil, fmt.Errorf("no usable releases found for github.com/%s/%s", owner, repo)
	}

	prefound := 0
	for _, r := range releases {
		// FIXME - enable flagged use of github prereleases https://github.com/git-pkg/cli/issues/8
		if r.GetPrerelease() {
			prefound++
			continue
		}

		return r, nil
	}

	return nil, fmt.Errorf("couldn't find latest release for github.com/%s/%s", owner, repo)
}

func releaseByVersion(owner, repo string, v *semver.Version) (*ghc.RepositoryRelease, error) {
	var err error
	var resp *ghc.Response

	c := ghc.NewClient(nil)
	tag := fmt.Sprintf("v%s", v.String())
	release, resp, err := c.Repositories.GetReleaseByTag(context.TODO(), owner, repo, tag)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return release, nil
}
