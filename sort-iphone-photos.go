package main

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func main() {
	inputDir := os.Args[1]

	fmt.Printf("Traversing directory %s\n\n", inputDir)

	exif.RegisterParsers(mknote.All...)

	err := filepath.Walk(inputDir, func(path string, fileInfo os.FileInfo, err error) error {

		if err != nil {
			fmt.Printf("%s\n\terror processing path\n", path, err)
			return nil
		}

		if fileInfo.IsDir() {
			return nil
		}

		imgFile, err := os.Open(path)
		if err != nil {
			fmt.Printf("%s\n\terror opening file (%s)\n", path, err)
			return nil
		}

		metaData, err := exif.Decode(imgFile)
		if err != nil {
			fmt.Printf("%s\n\terror parsing file (%s)\n", path, err)
			return nil
		}

		fmt.Printf("%s\n", path)

		make := getField(metaData, exif.Make)
		fmt.Printf("\tMake: %s\n", make)

		model := getField(metaData, exif.Model)
		fmt.Printf("\tModel: %s\n", model)

		return nil
	})

	if err != nil {
		fmt.Printf("error iterating directory structure (%s)", err)
	}
}

func getField(metaData *exif.Exif, fieldName exif.FieldName) string {
	fieldValue, err := metaData.Get(fieldName)
	if err != nil {
		fmt.Printf("\terror parsing %s (%s)\n", fieldName, err)
		return "N/A"
	}

	fieldStringValue, err := fieldValue.StringVal()
	if err != nil {
		fmt.Printf("\terror converting %s to string (%s)\n", fieldName, err)
		return "N/A"
	}

	return fieldStringValue
}
