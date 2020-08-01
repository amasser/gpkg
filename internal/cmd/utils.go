package cmd

import (
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/git-pkg/gpkg/pkg/constants"
	"github.com/git-pkg/gpkg/pkg/packages"
	"github.com/git-pkg/gpkg/pkg/validation"
)

// -----------------------------------------------------------------------------
// Private Functions - Utils
// -----------------------------------------------------------------------------

func installDir() (string, error) {
	h, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	dir := fmt.Sprintf("%s/%s", h, constants.DefaultBaseDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	return dir, nil
}

func packageForURL(strURL string) (*packages.Package, error) {
	u, err := validation.ValidateRepositoryURL(strURL)
	if err != nil {
		return nil, err
	}

	trimmedPath := strings.TrimSuffix(strings.TrimPrefix(u.Path, "/"), "/")
	path := strings.Split(trimmedPath, "/")
	if len(path) != 3 {
		return nil, fmt.Errorf("expected <host>/<owner>/<repo> got %s", u)
	}
	owner, repo := path[1], path[2]

	pkg, err := packages.NewPackageFromGithub(owner, repo)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}
