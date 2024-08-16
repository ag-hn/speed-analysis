// Package filesystem is a collection of various different filesystem
// helper functions.
package filesystem

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	thousand    = 1000
	ten         = 10
	fivePercent = 0.0499
)

// Different types of listings.
const (
	DirectoriesListingType = "directories"
	FilesListingType       = "files"
)

// GetDirectoryListing returns a list of files and directories within a given directory.
func GetDirectoryListing(dir string, showHidden bool) ([]fs.DirEntry, error) {
	index := 0

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Unwrap(err)
	}

	if !showHidden {
		for _, file := range files {
			// If the file or directory starts with a dot,
			// we know its hidden so dont add it to the array
			// of files to return.
			if !strings.HasPrefix(file.Name(), ".") {
				files[index] = file
				index++
			}
		}

		// Set files to the list that does not include hidden files.
		files = files[:index]
	}

	return files, nil
}

// GetDirectoryListingByType returns a directory listing based on type (directories | files).
func GetDirectoryListingByType(dir, listingType string, showHidden bool) ([]fs.DirEntry, error) {
	index := 0

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Unwrap(err)
	}

	for _, file := range files {
		switch {
		case file.IsDir() && listingType == DirectoriesListingType && !showHidden:
			if !strings.HasPrefix(file.Name(), ".") {
				files[index] = file
				index++
			}
		case file.IsDir() && listingType == DirectoriesListingType && showHidden:
			files[index] = file
			index++
		case !file.IsDir() && listingType == FilesListingType && !showHidden:
			if !strings.HasPrefix(file.Name(), ".") {
				files[index] = file
				index++
			}
		case !file.IsDir() && listingType == FilesListingType && showHidden:
			files[index] = file
			index++
		}
	}

	return files[:index], nil
}

// GetHomeDirectory returns the users home directory.
func GetHomeDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Unwrap(err)
	}

	return home, nil
}

// GetWorkingDirectory returns the current working directory.
func GetWorkingDirectory() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", errors.Unwrap(err)
	}

	return workingDir, nil
}

// ReadFileContent returns the contents of a file given a name.
func ReadFileContent(name string) (string, error) {
	fileContent, err := os.ReadFile(filepath.Clean(name))
	if err != nil {
		return "", errors.Unwrap(err)
	}

	return string(fileContent), nil
}

// GetDirectoryItemSize calculates the size of a directory or file.
func GetDirectoryItemSize(path string) (int64, error) {
	var size int64

	curFile, err := os.Stat(path)
	if err != nil {
		return 0, errors.Unwrap(err)
	}

	if !curFile.IsDir() {
		return curFile.Size(), nil
	}

	err = filepath.WalkDir(path, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return errors.Unwrap(err)
		}

		fileInfo, err := entry.Info()
		if err != nil {
			return errors.Unwrap(err)
		}

		if !entry.IsDir() {
			size += fileInfo.Size()
		}

		return errors.Unwrap(err)
	})

	return size, errors.Unwrap(err)
}

// FindFilesByName returns files found based on a name.
func FindFilesByName(name, dir string) ([]string, []fs.DirEntry, error) {
	var paths []string
	var entries []fs.DirEntry

	err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		if strings.Contains(entry.Name(), name) {
			paths = append(paths, path)
			entries = append(entries, entry)
		}

		return errors.Unwrap(err)
	})

	return paths, entries, errors.Unwrap(err)
}

// WriteToFile writes content to a file, overwriting content if it exists.
func WriteToFile(path, content string) error {
	file, err := os.OpenFile(filepath.Clean(path), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return errors.Unwrap(err)
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return errors.Unwrap(err)
	}

	_, err = file.WriteString(fmt.Sprintf("%s\n", filepath.Join(workingDir, content)))
	if err != nil {
		err = file.Close()
		if err != nil {
			return errors.Unwrap(err)
		}

		return errors.Unwrap(err)
	}

	err = file.Close()
	if err != nil {
		return errors.Unwrap(err)
	}

	return errors.Unwrap(err)
}

// ConvertBytesToSizeString converts a byte count to a human readable string.
func ConvertBytesToSizeString(size int64) string {
	if size < thousand {
		return fmt.Sprintf("%dB", size)
	}

	suffix := []string{
		"K", // kilo
		"M", // mega
		"G", // giga
		"T", // tera
		"P", // peta
		"E", // exa
		"Z", // zeta
		"Y", // yotta
	}

	curr := float64(size) / thousand
	for _, s := range suffix {
		if curr < ten {
			return fmt.Sprintf("%.1f%s", curr-fivePercent, s)
		} else if curr < thousand {
			return fmt.Sprintf("%d%s", int(curr), s)
		}
		curr /= thousand
	}

	return ""
}
