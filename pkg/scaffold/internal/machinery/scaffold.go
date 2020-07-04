package machinery

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/seamounts/kubeapi/pkg/model"
	"github.com/seamounts/kubeapi/pkg/model/file"
	"github.com/seamounts/kubeapi/pkg/scaffold/internal/filesystem"
	"golang.org/x/tools/imports"
)

var options = imports.Options{
	Comments:   true,
	TabIndent:  true,
	TabWidth:   8,
	FormatOnly: true,
}

// Scaffold uses templates to scaffold new files
type Scaffold interface {
	// Execute writes to disk the provided files
	Execute(universe *model.Universe, files ...file.Builder) error
}

// scaffold implements Scaffold interface
type scaffold struct {
	// fs allows to mock the file system for tests
	fs filesystem.FileSystem
}

// NewScaffold returns a new Scaffold with the provided plugins
func NewScaffold() Scaffold {
	return &scaffold{
		fs: filesystem.New(),
	}
}

// Execute implements Scaffold.Execute
func (s *scaffold) Execute(universe *model.Universe, files ...file.Builder) error {
	// Initialize the universe files
	universe.Files = make(map[string]*file.File, len(files))

	// Set the repo as the local prefix so that it knows how to group imports
	if universe.Config != nil {
		imports.LocalPrefix = universe.Config.Repo
	}

	for _, f := range files {
		// Inject common fields
		universe.InjectInto(f)

		// Build models for Template builders
		if t, isTemplate := f.(file.Template); isTemplate {
			if err := s.buildFileModel(t, universe.Files); err != nil {
				return err
			}
		}
	}

	// Persist the files to disk
	for _, f := range universe.Files {
		if err := s.writeFile(f); err != nil {
			return err
		}
	}

	return nil
}

func (s *scaffold) buildFileModel(t file.Template, models map[string]*file.File) error {
	// Set the template default values
	err := t.SetTemplateDefaults()
	if err != nil {
		return err
	}

	// Handle already existing models
	if _, found := models[t.GetPath()]; found {
		switch t.GetIfExistsAction() {
		case file.Skip:
			return nil
		case file.Error:
			return modelAlreadyExistsError{t.GetPath()}
		case file.Overwrite:
		default:
			return unknownIfExistsActionError{t.GetPath(), t.GetIfExistsAction()}
		}
	}

	m := &file.File{
		Path:           t.GetPath(),
		IfExistsAction: t.GetIfExistsAction(),
	}

	b, err := doTemplate(t)
	if err != nil {
		return err
	}
	m.Contents = string(b)

	models[m.Path] = m

	return nil

}

// doTemplate executes the template for a file using the input
func doTemplate(t file.Template) ([]byte, error) {
	temp, err := newTemplate(t).Parse(t.GetBody())
	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}
	err = temp.Execute(out, t)
	if err != nil {
		return nil, err
	}
	b := out.Bytes()

	// TODO(adirio): move go-formatting to write step
	// gofmt the imports
	if filepath.Ext(t.GetPath()) == ".go" {
		b, err = imports.Process(t.GetPath(), b, &options)
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

// newTemplate a new template with common functions
func newTemplate(t file.Template) *template.Template {
	fm := file.DefaultFuncMap()
	useFM, ok := t.(file.UseCustomFuncMap)
	if ok {
		fm = useFM.GetFuncMap()
	}
	return template.New(fmt.Sprintf("%T", t)).Funcs(fm)
}

func (s *scaffold) writeFile(f *file.File) error {
	// Check if the file to write already exists
	exists, err := s.fs.Exists(f.Path)
	if err != nil {
		return err
	}
	if exists {
		switch f.IfExistsAction {
		case file.Overwrite:
			// By not returning, the file is written as if it didn't exist
		case file.Skip:
			// By returning nil, the file is not written but the process will carry on
			return nil
		case file.Error:
			// By returning an error, the file is not written and the process will fail
			return fileAlreadyExistsError{f.Path}
		}
	}

	writer, err := s.fs.Create(f.Path)
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(f.Contents))

	return err
}
