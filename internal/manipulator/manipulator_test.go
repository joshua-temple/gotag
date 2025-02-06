package manipulator

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const sampleSourceNoTags = "package sample\n\n" +
	"type Person struct {\n" +
	"\tName string\n" +
	"\tAge  int\n" +
	"}\n"

func TestProcessFile_AddTag(t *testing.T) {
	// Create a temporary file with a sample source that has no tags.
	tmpFile, err := ioutil.TempFile("", "sample_*.go")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(sampleSourceNoTags)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Force-add the "json" tag in non-interactive mode.
	err = ProcessFile(tmpFile.Name(), "Person", []string{"json"}, nil, nil, "camelCase", false, true)
	if err != nil {
		t.Fatal(err)
	}

	updated, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	content := string(updated)
	if !strings.Contains(content, "json:\"name\"") {
		t.Errorf("expected tag json:\"name\", got: %s", content)
	}
	if !strings.Contains(content, "json:\"age\"") {
		t.Errorf("expected tag json:\"age\", got: %s", content)
	}
}

func TestProcessFile_DeleteTag(t *testing.T) {
	// Define a sample source with existing tags.
	source := "package sample\n\n" +
		"type Person struct {\n" +
		"\tName string `json:\"name\" xml:\"name\"`\n" +
		"\tAge  int    `json:\"age\"`\n" +
		"}\n"
	tmpFile, err := ioutil.TempFile("", "sample_*.go")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(source)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Force-delete the "xml" tag.
	err = ProcessFile(tmpFile.Name(), "Person", nil, []string{"xml"}, nil, "camelCase", false, true)
	if err != nil {
		t.Fatal(err)
	}

	updated, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	content := string(updated)
	if strings.Contains(content, "xml:") {
		t.Errorf("expected xml tag to be removed, got: %s", content)
	}
}

func TestProcessFile_OverwriteTag(t *testing.T) {
	// Define a sample source with an existing json tag.
	source := "package sample\n\n" +
		"type Person struct {\n" +
		"\tName string `json:\"name\"`\n" +
		"\tAge  int\n" +
		"}\n"
	tmpFile, err := ioutil.TempFile("", "sample_*.go")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(source)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Force-overwrite the "json" tag to "full_name".
	err = ProcessFile(tmpFile.Name(), "Person", nil, nil, []string{"json=full_name"}, "camelCase", false, true)
	if err != nil {
		t.Fatal(err)
	}

	updated, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	content := string(updated)
	if !strings.Contains(content, "json:\"full_name\"") {
		t.Errorf("expected tag json:\"full_name\", got: %s", content)
	}
}

func TestProcessFile_MultipleTagsUniformOrder(t *testing.T) {
	// Define a sample source with no tags.
	source := "package sample\n\n" +
		"type Person struct {\n" +
		"\tName string\n" +
		"}\n"
	tmpFile, err := ioutil.TempFile("", "sample_*.go")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(source)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Force-add multiple tags: json, xml, and db.
	err = ProcessFile(tmpFile.Name(), "Person", []string{"json", "xml", "db"}, nil, nil, "camelCase", false, true)
	if err != nil {
		t.Fatal(err)
	}

	updated, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	content := string(updated)

	// Expected uniform ordering: keys sorted alphabetically: db, json, xml.
	expectedOrder := []string{`db:"name"`, `json:"name"`, `xml:"name"`}
	prevIdx := -1
	for _, exp := range expectedOrder {
		idx := strings.Index(content, exp)
		if idx == -1 {
			t.Errorf("expected to find %s in content: %s", exp, content)
		}
		if idx < prevIdx {
			t.Errorf("tags are not in uniform sorted order, expected order: %v", expectedOrder)
		}
		prevIdx = idx
	}
}
