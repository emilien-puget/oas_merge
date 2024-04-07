package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pb33f/libopenapi"

	"github.com/emilien-puget/oas_merge/pkg/merge"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

const (
	defaultOutputFile   = "result"
	defaultOutputFormat = "yaml"
)

func main() {
	if len(os.Args) != 3 {
		panic("missing args, usage : main_oas dirpath")
	}
	mainOas := os.Args[1]
	dirPath := os.Args[2]

	files, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	specs := make([]*v3.Document, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			path := filepath.Join(dirPath, file.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}

			document, err := libopenapi.NewDocument(data)
			if err != nil {
				panic(err)
			}
			model, errs := document.BuildV3Model()
			if len(errs) != 0 {
				panic(errors.Join(errs...))
			}

			specs = append(specs, &model.Model)
		}
	}

	outputFileFlag := flag.String("output", defaultOutputFile, "output file")
	outputFileFormatFlag := flag.String("output_format", defaultOutputFormat, "output file format")
	flag.Parse()

	outputFile := defaultOutputFile
	if outputFileFlag != nil && *outputFileFlag != "" {
		outputFile = *outputFileFlag
	}
	outputFileFormat := "yaml"
	if outputFileFormatFlag != nil && *outputFileFormatFlag != "" {
		outputFileFormat = *outputFileFormatFlag
	}

	data, err := os.ReadFile(mainOas)
	if err != nil {
		panic(err)
	}
	document, err := libopenapi.NewDocument(data)
	if err != nil {
		panic(err)
	}
	model, errs := document.BuildV3Model()
	if len(errs) != 0 {
		panic(errors.Join(errs...))
	}

	mergedSpec, err := merge.OpenAPISpecsWithMain(&model.Model, specs)
	if err != nil {
		panic(err)
	}

	if err := writeOutput(outputFile, outputFileFormat, mergedSpec); err != nil {
		panic(err)
	}
}

func writeOutput(outputFile string, format string, mergedSpec *v3.Document) error {
	var body []byte
	switch format {
	case "json":
		body = mergedSpec.RenderJSON(" ")
	default:
		body = mergedSpec.RenderWithIndention(2)
	}

	err := os.WriteFile(outputFile+"."+format, body, 0o644)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}
