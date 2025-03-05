package command

import "os"

// FileWriter defines an interface for writing files.
type FileWriter interface {
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

// DefaultFileWriter implements FileWriter using os.WriteFile.
type DefaultFileWriter struct{}

func (d *DefaultFileWriter) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}
