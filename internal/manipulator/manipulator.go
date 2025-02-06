package manipulator

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/iancoleman/strcase"
)

// ProcessFile processes a Go source file, updating struct tags according to directives.
// filePath: path to the file
// targetStruct: if non-empty, only update the struct with that name.
// addTags: slice of tag keys to add
// deleteTags: slice of tag keys to delete
// overwriteTags: slice of strings "key=newValue" (if newValue is empty, interactive prompt or default is used)
// caseStyle: style for tag value conversion ("camelCase", "snake_case", "kebab-case")
// interactive: if true, prompt for overrides/confirmations for add or overwrite operations
// force: if true, bypass any confirmation prompts
func ProcessFile(filePath, targetStruct string, addTags, deleteTags, overwriteTags []string, caseStyle string, interactive bool, force bool) error {
	src, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)
	if err != nil {
		return err
	}

	// Parse overwriteTags into a map.
	overwriteMap := make(map[string]string)
	for _, ot := range overwriteTags {
		parts := strings.SplitN(ot, "=", 2)
		key := parts[0]
		value := ""
		if len(parts) > 1 {
			value = parts[1]
		}
		overwriteMap[key] = value
	}

	changed := false

	// Traverse AST for type declarations.
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// If targetStruct is specified, skip other structs.
			if targetStruct != "" && typeSpec.Name.Name != targetStruct {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// Process struct fields.
			for _, field := range structType.Fields.List {
				var tagValue string
				if field.Tag != nil {
					tagValue = field.Tag.Value // e.g. "`json:\"name\"`"
					tagValue = strings.Trim(tagValue, "`")
				}
				// Convert tag string into a map.
				tags := parseTag(tagValue)

				// Process delete directives.
				for _, key := range deleteTags {
					if _, exists := tags[key]; exists {
						delete(tags, key)
						changed = true
					}
				}

				// Process add directives.
				for _, key := range addTags {
					defaultVal := defaultTagValue(field, caseStyle)
					if current, exists := tags[key]; exists {
						// If tag exists and the current value does not match the default, then if interactive and not forced, confirm update.
						if current != defaultVal {
							update := force
							if !force && interactive {
								update = confirmUpdate(fmt.Sprintf("Field '%s': tag '%s' value is '%s'. Update to '%s'? (y/n): ", fieldName(field), key, current, defaultVal))
							}
							if update {
								tags[key] = defaultVal
								changed = true
							}
						}
					} else {
						// Tag does not exist: prompt if interactive, else add default.
						newVal := ""
						if interactive {
							fmt.Printf("Enter value for new tag '%s' in field '%s' (default %s): ", key, fieldName(field), defaultVal)
							scanner := bufio.NewScanner(os.Stdin)
							if scanner.Scan() {
								newVal = scanner.Text()
							}
						}
						if newVal == "" {
							newVal = defaultVal
						}
						tags[key] = newVal
						changed = true
					}
				}

				// Process overwrite directives.
				for key, newVal := range overwriteMap {
					if current, exists := tags[key]; exists {
						// If interactive and not forced and newVal is empty, confirm update.
						if interactive && !force && newVal == "" && current != defaultTagValue(field, caseStyle) {
							if !confirmUpdate(fmt.Sprintf("Field '%s': tag '%s' value is '%s'. Overwrite with default '%s'? (y/n): ", fieldName(field), key, current, defaultTagValue(field, caseStyle))) {
								continue
							}
						}
						if newVal == "" {
							newVal = defaultTagValue(field, caseStyle)
						}
						// Even if not interactive, if new value differs from current, update.
						if current != newVal {
							tags[key] = newVal
							changed = true
						}
					}
				}

				// Rebuild tag string from map with sorted keys for uniform ordering.
				newTag := buildTag(tags)
				if newTag == "" {
					field.Tag = nil
				} else {
					field.Tag = &ast.BasicLit{
						Kind:  token.STRING,
						Value: "`" + newTag + "`",
					}
				}
			}
		}
	}

	// If changes were made, write the modified AST back to file.
	if changed {
		var buf bytes.Buffer
		if err := printer.Fprint(&buf, fset, file); err != nil {
			return err
		}
		formatted, err := format.Source(buf.Bytes())
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filePath, formatted, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

// confirmUpdate prompts the user to confirm an update. Returns true if confirmed.
func confirmUpdate(prompt string) bool {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		response := strings.ToLower(strings.TrimSpace(scanner.Text()))
		return response == "y" || response == "yes"
	}
	return false
}

// parseTag converts a struct tag string into a map of key-value pairs.
// Example: `json:"name" xml:"name"` becomes map["json"] = "name", map["xml"] = "name".
func parseTag(tag string) map[string]string {
	result := make(map[string]string)
	re := regexp.MustCompile(`(\w+):"([^"]*)"`)
	matches := re.FindAllStringSubmatch(tag, -1)
	for _, m := range matches {
		if len(m) == 3 {
			result[m[1]] = m[2]
		}
	}
	return result
}

// buildTag rebuilds a tag string from a map of key-value pairs with keys sorted alphabetically.
func buildTag(tags map[string]string) string {
	var parts []string
	var keys []string
	for key := range tags {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf(`%s:"%s"`, key, tags[key]))
	}
	return strings.Join(parts, " ")
}

// defaultTagValue returns the default tag value for a field based on its name and the given case style.
func defaultTagValue(field *ast.Field, caseStyle string) string {
	name := fieldName(field)
	switch strings.ToLower(caseStyle) {
	case "snake", "snake_case":
		return strcase.ToSnake(name)
	case "kebab", "kebab-case":
		return strcase.ToKebab(name)
	case "camel", "camelCase":
		fallthrough
	default: // default to camelCase
		return strcase.ToLowerCamel(name)
	}
}

// fieldName returns the name of the field.
func fieldName(field *ast.Field) string {
	if len(field.Names) > 0 {
		return field.Names[0].Name
	}
	return ""
}
