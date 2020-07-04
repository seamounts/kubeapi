package file

import "github.com/seamounts/kubeapi/pkg/model/resource"

// PathMixin provides file builders with a path field
type PathMixin struct {
	// Path is the of the file
	Path string
}

// GetPath implements Builder
func (t *PathMixin) GetPath() string {
	return t.Path
}

// IfExistsActionMixin provides file builders with a if-exists-action field
type IfExistsActionMixin struct {
	// IfExistsAction determines what to do if the file exists
	IfExistsAction IfExistsAction
}

// GetIfExistsAction implements Builder
func (t *IfExistsActionMixin) GetIfExistsAction() IfExistsAction {
	return t.IfExistsAction
}

// TemplateMixin is the mixin that should be embedded in Template builders
type TemplateMixin struct {
	PathMixin
	IfExistsActionMixin

	// TemplateBody is the template body to execute
	TemplateBody string
}

// GetBody implements Template
func (t *TemplateMixin) GetBody() string {
	return t.TemplateBody
}

// InserterMixin is the mixin that should be embedded in Inserter builders
type InserterMixin struct {
	PathMixin
}

// GetIfExistsAction implements Builder
func (t *InserterMixin) GetIfExistsAction() IfExistsAction {
	// Inserter builders always need to overwrite previous files
	return Overwrite
}

// RepositoryMixin provides templates with a injectable repository field
type RepositoryMixin struct {
	// Repo is the go project package path
	Repo string
}

// InjectRepository implements HasRepository
func (m *RepositoryMixin) InjectRepository(repository string) {
	if m.Repo == "" {
		m.Repo = repository
	}
}

// ResourceMixin provides templates with a injectable resource field
type ResourceMixin struct {
	Resource *resource.Resource
}

// InjectResource implements HasResource
func (m *ResourceMixin) InjectResource(res *resource.Resource) {
	if m.Resource == nil {
		m.Resource = res
	}
}

// BoilerplateMixin provides templates with a injectable boilerplate field
type BoilerplateMixin struct {
	// Boilerplate is the contents of a Boilerplate go header file
	Boilerplate string
}

// InjectBoilerplate implements HasBoilerplate
func (m *BoilerplateMixin) InjectBoilerplate(boilerplate string) {
	if m.Boilerplate == "" {
		m.Boilerplate = boilerplate
	}
}
