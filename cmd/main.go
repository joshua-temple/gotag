package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/joshua-temple/gotag/internal/manipulator"
	"github.com/joshua-temple/gotag/internal/parser"
)

var (
	target        string
	recursive     bool
	interactive   bool
	force         bool
	caseStyle     string
	addTags       []string
	deleteTags    []string
	overwriteTags []string
)

func splitCommaSeparated(input []string) []string {
	var result []string
	for _, v := range input {
		parts := strings.Split(v, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				result = append(result, part)
			}
		}
	}
	return result
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "gostructtag",
		Short: "A CLI tool to manipulate Go struct tags.",
		Long:  `A CLI tool to add, delete, or overwrite struct tags in Go source files.`,
		Run: func(cmd *cobra.Command, args []string) {
			if target == "" {
				fmt.Println("Error: target is required")
				os.Exit(1)
			}

			// Process comma-separated flags.
			addList := splitCommaSeparated(addTags)
			deleteList := splitCommaSeparated(deleteTags)
			overwriteList := splitCommaSeparated(overwriteTags)

			targetInfo, err := parser.ParseTarget(target)
			if err != nil {
				log.Fatalf("Error parsing target: %v", err)
			}

			files := []string{}
			switch targetInfo.Type {
			case parser.TargetTypeStruct, parser.TargetTypeFile:
				files = append(files, targetInfo.FilePath)
			case parser.TargetTypeDirectory:
				files, err = parser.GetGoFilesFromDir(targetInfo.FilePath, recursive)
				if err != nil {
					log.Fatalf("Error reading directory: %v", err)
				}
			}

			for _, file := range files {
				fmt.Printf("Processing file: %s\n", file)
				err := manipulator.ProcessFile(file, targetInfo.StructName, addList, deleteList, overwriteList, caseStyle, interactive, force)
				if err != nil {
					log.Printf("Error processing file %s: %v", file, err)
				} else {
					fmt.Printf("Updated file: %s\n", file)
				}
			}
		},
	}

	rootCmd.Flags().StringVarP(&target, "target", "t", "", "Target file/directory/struct (e.g. pkg/here/struct.go@StructName, pkg/here/struct.go, or pkg/here/)")
	rootCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively scan directories")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode for tag overrides (works with add or overwrite)")
	rootCmd.Flags().BoolVarP(&force, "force", "f", false, "Bypass confirmations and force changes")
	rootCmd.Flags().StringVarP(&caseStyle, "case", "c", "camelCase", "Case style for tag values (camel(Case)?, snake(_case)?, kebab(-case)?)")
	rootCmd.Flags().StringArrayVarP(&addTags, "add", "a", []string{}, "Tags to add (format: key). For multiple tags, use a comma-separated list e.g. -a json,xml,db")
	rootCmd.Flags().StringArrayVarP(&deleteTags, "delete", "d", []string{}, "Tags to delete (format: key). For multiple tags, use a comma-separated list e.g. -d json,xml")
	rootCmd.Flags().StringArrayVarP(&overwriteTags, "overwrite", "o", []string{}, "Tags to overwrite (format: key=newValue, use empty value for interactive/defaults)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
