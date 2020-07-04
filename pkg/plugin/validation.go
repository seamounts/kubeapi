package plugin

import (
	"errors"
	"fmt"

	"github.com/blang/semver"

	"github.com/seamounts/kubeapi/pkg/internal/validation"
)

// ValidateVersion ensures version adheres to the plugin version format,
// which is tolerant semver.
func ValidateVersion(version string) error {
	if version == "" {
		return errors.New("plugin version is empty")
	}
	// ParseTolerant allows versions with a "v" prefix or shortened versions,
	// ex. "3" or "v3.0".
	if _, err := semver.ParseTolerant(version); err != nil {
		return fmt.Errorf("failed to validate plugin version %q: %v", version, err)
	}
	return nil
}

// ValidateName ensures name is a valid DNS 1123 subdomain.
func ValidateName(name string) error {
	if errs := validation.IsDNS1123Subdomain(name); len(errs) != 0 {
		return fmt.Errorf("plugin name %q is invalid: %v", name, errs)
	}
	return nil
}
