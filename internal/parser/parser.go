package parser

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type TargetType int

const (
	TargetTypeStruct TargetType = iota
	TargetTypeFile
	TargetTypeDirectory
)

type TargetInfo struct {
	Type       TargetType
	FilePath   string
	StructName string // Only set for TargetTypeStruct
}

// ParseTarget parses a target string into a TargetInfo.
// Format can be: file@StructName, file.go, or a directory.
func ParseTarget(target string) (*TargetInfo, error) {
	if target == "" {
		return nil, errors.New("target cannot be empty")
	}

	if strings.Contains(target, "@") {
		parts := strings.Split(target, "@")
		if len(parts) != 2 {
			return nil, errors.New("invalid target format for struct")
		}
		filePath := parts[0]
		structName := parts[1]
		if _, err := os.Stat(filePath); err != nil {
			return nil, err
		}
		return &TargetInfo{
			Type:       TargetTypeStruct,
			FilePath:   filePath,
			StructName: structName,
		}, nil
	}

	// Check if target is a file.
	info, err := os.Stat(target)
	if err != nil {
		return nil, err
	}
	if info.Mode().IsRegular() && filepath.Ext(target) == ".go" {
		return &TargetInfo{
			Type:     TargetTypeFile,
			FilePath: target,
		}, nil
	}

	// Otherwise, assume directory.
	if info.IsDir() {
		return &TargetInfo{
			Type:     TargetTypeDirectory,
			FilePath: target,
		}, nil
	}

	return nil, errors.New("invalid target type")
}

// GetGoFilesFromDir returns a list of .go files in a directory.
// If recursive is true, it traverses subdirectories.
func GetGoFilesFromDir(dir string, recursive bool) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// If not recursive and path is not the same as directory, skip subdirectories.
		if !recursive && info.IsDir() && path != dir {
			return filepath.SkipDir
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
