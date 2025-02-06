package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTarget_Struct(t *testing.T) {
	tmpFile := "temp_test.go"
	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(tmpFile)

	target := tmpFile + "@MyStruct"
	info, err := ParseTarget(target)
	if err != nil {
		t.Fatal(err)
	}
	if info.Type != TargetTypeStruct {
		t.Errorf("expected TargetTypeStruct, got %v", info.Type)
	}
	if info.StructName != "MyStruct" {
		t.Errorf("expected StructName 'MyStruct', got %s", info.StructName)
	}
}

func TestParseTarget_File(t *testing.T) {
	tmpFile := "temp_test.go"
	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(tmpFile)

	info, err := ParseTarget(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	if info.Type != TargetTypeFile {
		t.Errorf("expected TargetTypeFile, got %v", info.Type)
	}
}

func TestParseTarget_Directory(t *testing.T) {
	tmpDir := "temp_test_dir"
	err := os.Mkdir(tmpDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	info, err := ParseTarget(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if info.Type != TargetTypeDirectory {
		t.Errorf("expected TargetTypeDirectory, got %v", info.Type)
	}
}

func TestGetGoFilesFromDir(t *testing.T) {
	tmpDir := "temp_test_dir"
	err := os.Mkdir(tmpDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	goFile := filepath.Join(tmpDir, "a.go")
	nonGoFile := filepath.Join(tmpDir, "a.txt")

	if err := os.WriteFile(goFile, []byte("package main"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(nonGoFile, []byte("text"), 0644); err != nil {
		t.Fatal(err)
	}

	files, err := GetGoFilesFromDir(tmpDir, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 go file, got %d", len(files))
	}
}
