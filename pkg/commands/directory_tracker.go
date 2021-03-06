package commands

import (
	build "github.com/Azure/acr-builder/pkg"
)

// DirectoryTracker tracks switching of directories
type DirectoryTracker struct {
	path string
}

// ChdirWithTracking performs a chdir and return the tracking object
func ChdirWithTracking(runner build.Runner, chdir string) (*DirectoryTracker, error) {
	if chdir == "" {
		return nil, nil
	}
	fs := runner.GetFileSystem()
	path, err := fs.Getwd()
	if err != nil {
		return nil, err
	}
	err = fs.Chdir(chdir)
	if err != nil {
		return nil, err
	}
	return &DirectoryTracker{path: path}, nil
}

// Return returns to the directory in affect when DirectoryTracker object is created
func (t *DirectoryTracker) Return(runner build.Runner) error {
	fs := runner.GetFileSystem()
	return fs.Chdir(t.path)
}
