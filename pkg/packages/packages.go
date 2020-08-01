package packages

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/blang/semver"

	"github.com/git-pkg/gpkg/internal/github"
	"github.com/git-pkg/gpkg/pkg/constants"
)

// -----------------------------------------------------------------------------
// Package - Public Functions
// -----------------------------------------------------------------------------

// NewPackageFromGithub produces a *Package given a valid owner and repo
func NewPackageFromGithub(owner, repo string) (*Package, error) {
	return NewPackageFromURL(fmt.Sprintf("%s/github.com/%s/%s/%s/%s.json", constants.DefaultPackageRepository, owner, repo, runtime.GOARCH, runtime.GOOS))
}

// NewPackageFromURL retrieves a package from the provided URL and returns a *Package
func NewPackageFromURL(u string) (*Package, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return NewPackageFromJSON(buf.String())
}

// NewPackageFromJSON unmarshals a provided JSON string into a *Package
func NewPackageFromJSON(p string) (*Package, error) {
	pkg := &Package{}

	if err := json.Unmarshal([]byte(p), pkg); err != nil {
		return nil, fmt.Errorf("%s was not a valid package: %w", p, err)
	}

	return pkg, nil
}

// -----------------------------------------------------------------------------
// Package - Public Types
// -----------------------------------------------------------------------------

// Package represents the metadata about an application and indicates how to install it to the local system and where to install it from
type Package struct {
	// URL the repository URL for the package
	URL string `json:"url,required"`

	// Method represents the methodology that will need to be imployed to install the relevant artifacts
	Method string `json:"method,required"`

	// Binaries indicates the paths to the binaries in the artifact which are meant to be installed for the end-user
	Binaries []string `json:"binaries,required"`

	// FIXME: Hackargs will go away using it for testing right now
	Hackargs string `json:"hackargs,required"`

	// Pattern indicates the regexp pattern to use to determine the correct asset from the release to use
	Pattern string `json:"pattern,required"`

	// Compression indicates the type of compression used (if any) for the relevant artifacts
	Compression string `json:"compression,omit_empty"`

	// SkipBaseDir is used for "tar" method artifacts and indicates all files in the archive will be nested under a base directory, which can be safely ignored.
	SkipBaseDir bool `json:"skipbasedir,omit_empty"`
}

// -----------------------------------------------------------------------------
// Package - Public Methods
// -----------------------------------------------------------------------------

func (p *Package) Uninstall(basedir string) error {
	repoURL, err := url.Parse(p.URL)
	if err != nil {
		return err
	}

	owner, repo, err := github.GetOwnerRepo(repoURL)
	if err != nil {
		return err
	}

	for _, bin := range p.Binaries {
		bin := fmt.Sprintf("%s/bin/%s", basedir, filepath.Base(bin))
		fmt.Printf("removing symlink %s...\n", bin)
		if err := os.Remove(bin); err != nil {
			return err
		}
	}

	installDir := fmt.Sprintf("%s/%s/%s/%s", basedir, repoURL.Host, owner, repo)
	fmt.Printf("removing %s installation...\n", repoURL)
	return os.RemoveAll(installDir)
}

// Install installs the *Package to the local system
func (p *Package) Install(basedir string, v *semver.Version) error {
	repoURL, err := url.Parse(p.URL)
	if err != nil {
		return err
	}

	owner, repo, err := github.GetOwnerRepo(repoURL)
	if err != nil {
		return err
	}

	installDir := fmt.Sprintf("%s/%s/%s/%s", basedir, repoURL.Host, owner, repo)
	if _, err := os.Stat(installDir); err == nil {
		return fmt.Errorf("%s already installed, try uninstalling first", repoURL)
	}

	if err := os.MkdirAll(installDir, 0755); err != nil {
		return err
	}

	re, err := regexp.Compile(p.Pattern)
	if err != nil {
		return err
	}

	artifactURL, err := github.GetReleaseArtifactURL(basedir, owner, repo, v, re)
	if err != nil {
		return err
	}

	r, err := http.Get(artifactURL.String())
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", r.Status)
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	if mediaType != constants.DefaultContentType {
		return fmt.Errorf("expected download content type to be %s, received %s", constants.DefaultContentType, mediaType)
	}

	if p.Method != constants.TarMethod {
		return fmt.Errorf("package %s was configured for %s method, but only %s method is currently supported", p.URL, p.Method, constants.TarMethod)
	}

	// TODO: we only support tar archives and gunzip compression for now, will support more later
	var ar *tar.Reader
	if p.Compression != "" {
		if p.Compression != constants.GunzipCompression {
			return fmt.Errorf("package %s had compression type %s, but only %s is currently supported", p.URL, p.Compression, constants.GunzipCompression)
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return err
		}
		ar = tar.NewReader(gz)
	} else {
		ar = tar.NewReader(r.Body)
	}

	skipDir := ""
	for {
		hdr, err := ar.Next()
		if err == io.EOF {
			break // end of archive
		}

		if err != nil {
			return err
		}

		filename := filepath.Base(hdr.Name)
		dirname := filepath.Dir(hdr.Name)
		if p.SkipBaseDir {
			dir := strings.Split(dirname, "/")
			if len(dir) == 1 {
				skipDir = dir[0]
				continue
			}
		}

		subpath := fmt.Sprintf("%s/%s", installDir, dirname)
		if skipDir != "" {
			re := regexp.MustCompile(skipDir)
			subpath = re.ReplaceAllString(subpath, "")
		}

		if err := os.MkdirAll(subpath, 0755); err != nil {
			return err
		}

		f, err := os.OpenFile(fmt.Sprintf("%s/%s", subpath, filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, hdr.FileInfo().Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(f, ar)
		if err != nil {
			return err
		}

		f.Close()
	}

	for _, bin := range p.Binaries {
		binDir := fmt.Sprintf("%s/bin", basedir)
		if err := os.MkdirAll(binDir, 0755); err != nil {
			return err
		}
		binPath := fmt.Sprintf("%s/%s", installDir, bin)
		binName := filepath.Base(binPath)
		if err := os.Symlink(binPath, fmt.Sprintf("%s/%s", binDir, binName)); err != nil {
			return err
		}
	}

	return nil
}
