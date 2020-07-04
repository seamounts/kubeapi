package config

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

// Scaffolding versions
const (
	Version1 = "1"
)

// Config is the unmarshalled representation of the configuration file
type Config struct {
	// Repo is the go package name of the project root
	Repo string `json:"repo,omitempty"`

	Resources []GVK `json:"resources,omitempty"`

	// Version is the project version, defaults to "1" (backwards compatibility)
	Version string `json:"version,omitempty"`

	// Domain is the domain associated with the project and used for API groups
	Domain string `json:"domain,omitempty"`
}

// IsV1 returns true if it is a v1 project
func (c Config) IsV1() bool {
	return c.Version == Version1
}

// Marshal returns the bytes of c.
func (c Config) Marshal() ([]byte, error) {
	// Ignore extra fields at first.
	cfg := c

	content, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("error marshalling project configuration: %v", err)
	}
	// Empty config strings are "{}" due to the map field.
	if strings.TrimSpace(string(content)) == "{}" {
		content = []byte{}
	}

	return content, nil
}

// Unmarshal unmarshals the bytes of a Config into c.
func (c *Config) Unmarshal(b []byte) error {
	if err := yaml.UnmarshalStrict(b, c); err != nil {
		return fmt.Errorf("error unmarshalling project configuration: %v", err)
	}

	return nil
}

// HasResource returns true if API resource is already tracked
func (c Config) HasResource(target GVK) bool {
	// Return true if the target resource is found in the tracked resources
	for _, r := range c.Resources {
		if r.isEqualTo(target) {
			return true
		}
	}

	// Return false otherwise
	return false
}

// AddResource appends the provided resource to the tracked ones
// It returns if the configuration was modified
// NOTE: in v1 resources are not tracked, so we return false
func (c *Config) AddResource(gvk GVK) bool {
	// Short-circuit v1
	if c.IsV1() {
		return false
	}

	// No-op if the resource was already tracked, return false
	if c.HasResource(gvk) {
		return false
	}

	// Append the resource to the tracked ones, return true
	c.Resources = append(c.Resources, gvk)
	return true
}

// HasGroup returns true if group is already tracked
func (c Config) HasGroup(group string) bool {
	// Return true if the target group is found in the tracked resources
	for _, r := range c.Resources {
		if strings.EqualFold(group, r.Group) {
			return true
		}
	}

	// Return false otherwise
	return false
}

// GVK contains information about scaffolded resources
type GVK struct {
	Group   string `json:"group,omitempty"`
	Version string `json:"version,omitempty"`
	Kind    string `json:"kind,omitempty"`
}

// isEqualTo compares it with another resource
func (r GVK) isEqualTo(other GVK) bool {
	return r.Group == other.Group &&
		r.Version == other.Version &&
		r.Kind == other.Kind
}
