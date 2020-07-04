package templates

import (
	"path/filepath"

	"github.com/seamounts/kubeapi/pkg/model/file"
)

type Doc struct {
	file.TemplateMixin
	file.ResourceMixin
	file.RepositoryMixin
}

// GetBody implements Template
func (d *Doc) GetBody() string {
	return d.TemplateBody
}

func (d *Doc) SetTemplateDefaults() error {
	d.Path = filepath.Join("apis", "%[group]", "%[version]", "doc.go")
	d.Path = d.Resource.Replacer().Replace(d.Path)

	d.TemplateBody = docTemplate

	d.IfExistsAction = file.Error

	return nil
}

const docTemplate = `
// Package v1 contains API Schema definitions for the ecs v1 API group
// +k8s:deepcopy-gen=package,register
// +groupName={{ .Resource.Group }}
package v1

const (
	GroupName = {{ .Resource.Group }}
	Version = {{ .Resource.Version }}
)
`
