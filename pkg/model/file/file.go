package file

// IfExistsAction determines what to do if the scaffold file already exists
type IfExistsAction int

const (
	// Skip skips the file and moves to the next one
	Skip IfExistsAction = iota

	// Error returns an error and stops processing
	Error

	// Overwrite truncates and overwrites the existing file
	Overwrite
)

// File describes a file that will be written
type File struct {
	// Path is the file to write
	Path string `json:"path,omitempty"`

	// Contents is the generated output
	Contents string `json:"contents,omitempty"`

	// IfExistsAction determines what to do if the file exists
	IfExistsAction IfExistsAction `json:"ifExistsAction,omitempty"`
}
