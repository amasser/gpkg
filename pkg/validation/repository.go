package validation

import (
	"fmt"
	"net/url"
)

// -----------------------------------------------------------------------------
// Public Functions - Validation
// -----------------------------------------------------------------------------

// ValidateRepositoryURL parses and validates the URL for usability
func ValidateRepositoryURL(u string) (*url.URL, error) {
	// parse the provided string URL into *url.URL
	repoURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	if repoURL.Scheme == "" {
		repoURL.Scheme = "https"
	}
	if repoURL.Scheme != "https" {
		return nil, fmt.Errorf("received %s, only https:// scheme is supported", repoURL)
	}

	if repoURL.Path == "" {
		return nil, fmt.Errorf("no path to the repository was provided, only received %s", repoURL)
	}

	return repoURL, nil
}
