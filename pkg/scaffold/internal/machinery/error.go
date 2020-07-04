package machinery

import (
	"errors"
	"fmt"

	"github.com/seamounts/kubeapi/pkg/model/file"
)

// fileAlreadyExistsError is returned if the file is expected not to exist but it does
type fileAlreadyExistsError struct {
	path string
}

// Error implements error interface
func (e fileAlreadyExistsError) Error() string {
	return fmt.Sprintf("failed to create %s: file already exists", e.path)
}

// IsFileAlreadyExistsError checks if the returned error is because the file already existed when expected not to
func IsFileAlreadyExistsError(err error) bool {
	return errors.As(err, &fileAlreadyExistsError{})
}

// modelAlreadyExistsError is returned if the file is expected not to exist but a previous model does
type modelAlreadyExistsError struct {
	path string
}

// Error implements error interface
func (e modelAlreadyExistsError) Error() string {
	return fmt.Sprintf("failed to create %s: model already exists", e.path)
}

// IsModelAlreadyExistsError checks if the returned error is because the model already existed when expected not to
func IsModelAlreadyExistsError(err error) bool {
	return errors.As(err, &modelAlreadyExistsError{})
}

// unknownIfExistsActionError is returned if the if-exists-action is unknown
type unknownIfExistsActionError struct {
	path           string
	ifExistsAction file.IfExistsAction
}

// Error implements error interface
func (e unknownIfExistsActionError) Error() string {
	return fmt.Sprintf("unknown behavior if file exists (%d) for %s", e.ifExistsAction, e.path)
}

// IsUnknownIfExistsActionError checks if the returned error is because the if-exists-action is unknown
func IsUnknownIfExistsActionError(err error) bool {
	return errors.As(err, &unknownIfExistsActionError{})
}
