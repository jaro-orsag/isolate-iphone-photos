package processor

import (
	"log"
	"os"

	filepath "path/filepath"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func MakeCollectMetadata(postProcess PostProcessFunc) filepath.WalkFunc {
	exif.RegisterParsers(mknote.All...)

	return func(path string, fileInfo os.FileInfo, err error) error {

		log.Printf("processing %s\n", path)

		if err != nil {
			log.Printf("\terror processing path: %s\n", err)

			return nil
		}

		if fileInfo.IsDir() {
			log.Printf("\tnot processing path to directory\n")

			return nil
		}

		imgFile, err := os.Open(path)
		if err != nil {
			log.Printf("\terror opening file: %s\n", err)

			return nil
		}

		exifData, err := exif.Decode(imgFile)
		if err != nil {
			log.Printf("\terror parsing file: %s\n", err)

			return nil
		}

		metadata := &metadata{
			Path:     path,
			Make:     getField(exifData, exif.Make),
			Model:    getField(exifData, exif.Model),
			DateTime: getField(exifData, exif.DateTime),
		}

		log.Printf("\tsuccess parsing file: %#v\n", metadata)

		return postProcess(metadata)
	}
}

func getField(metaData *exif.Exif, fieldName exif.FieldName) string {
	fieldValue, err := metaData.Get(fieldName)
	if err != nil {
		log.Printf("\terror parsing field %s: %s", fieldName, err)

		return "N/A"
	}

	fieldStringValue, err := fieldValue.StringVal()
	if err != nil {
		log.Printf("\terror converting value of %s to string: %s", fieldName, err)

		return "N/A"
	}

	return fieldStringValue
}
