package fileutil

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const PermFileDefault = 0644
const PermDirDefault = 0777
const PermOwnerRW = 0600

// MkdirAll creates the dir if it does not exist.
// Parent dirs be created if required
func MkdirAll(dirPath string) error {
	if !PathExists(dirPath) {
		err := os.MkdirAll(dirPath, PermDirDefault)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// truncate creates the file if it does not exist,
// or truncates if it already exists.
// Parent dirs be created if required.
// Remember to call f.Close
func truncate(filePath string, perm fs.FileMode) (f *os.File, err error) {
	dirPath := filepath.Dir(filePath)
	err = MkdirAll(dirPath)
	if err != nil {
		return f, errors.WithStack(err)
	}

	// Create or truncate file
	f, err = os.OpenFile(
		filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return f, errors.WithStack(err)
	}

	return f, nil
}

// append creates the file if it does not exist,
// or opens it for appending data.
// Remember to call f.Close
func append(filePath string) (f *os.File, err error) {
	dirPath := filepath.Dir(filePath)
	err = MkdirAll(dirPath)
	if err != nil {
		return f, errors.WithStack(err)
	}

	// Create or open file
	f, err = os.OpenFile(
		filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, PermFileDefault)
	if err != nil {
		return f, errors.WithStack(err)
	}

	return f, nil
}

func OpenWithPerm(filePath string, perm fs.FileMode) (f *os.File, err error) {
	f, err = os.OpenFile(
		filePath, os.O_RDONLY, perm)
	if err != nil {
		return f, errors.WithStack(err)
	}
	return f, nil
}

// Open file for reading, remember to call f.Close
func Open(filePath string) (f *os.File, err error) {
	return OpenWithPerm(filePath, PermFileDefault)
}

func ReadAll(filePath string) (b []byte, err error) {
	f, err := Open(filePath)
	if err != nil {
		return b, err
	}
	defer f.Close()
	b, err = io.ReadAll(f)
	if err != nil {
		return b, errors.WithStack(err)
	}
	return b, nil
}

func WriteBytesWithPerm(filePath string, b []byte, perm fs.FileMode) (err error) {
	f, err := truncate(filePath, perm)
	if err != nil {
		return err
	}
	// Example here doesn't check the error
	// https://gobyexample.com/writing-files
	defer f.Close()
	_, err = f.Write(b)
	if err != nil {
		return errors.WithStack(err)
	}
	return f.Sync()
}

// WriteBytes to a file.
// Creates the file if it does not exist,
// or truncates existing files before writing
func WriteBytes(filePath string, b []byte) (err error) {
	return WriteBytesWithPerm(filePath, b, PermFileDefault)
}

// WriteBytes to a file.
// Creates the file if it does not exist
func AppendBytes(filePath string, b []byte) (err error) {
	f, err := append(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	if err != nil {
		return errors.WithStack(err)
	}
	return f.Sync()
}

// PathExists returns true if the specified path exists
func PathExists(checkPath string) bool {
	_, err := os.Stat(checkPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		// WARNING Will also return false due to other errors,
		// like insufficient permission to list path
		return false
	}
	return true
}

// IsDir returns true if the path is a directory
func IsDir(checkPath string) bool {
	info, err := os.Stat(checkPath)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return true
	}
	return false
}

// Copy file at source to dest
// TODO For larger files, consider using bufio.NewReader
// and bufio.Writer for better performance?
func Copy(source, dest string) (err error) {
	// Open source file for reading
	sourceFile, err := os.Open(source)
	if err != nil {
		panic(err)
	}
	defer sourceFile.Close()

	// Create destination file for writing
	destFile, err := os.Create(dest)
	if err != nil {
		panic(err)
	}
	defer destFile.Close()

	// Copy the file content
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	return nil
}
