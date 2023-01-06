package processor

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	filepath "path/filepath"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

const exifDateTimeFormat = "2006:01:02 15:04:05"

// TODO: following constants should be generalized or configurable
const videoFileSuffix = ".mov"
const expectedMake = "Apple"
const expectedModel = "iPhone 12 mini"

func MakeCollectMetadata(postProcess PostProcessFunc) filepath.WalkFunc {
	exif.RegisterParsers(mknote.All...)

	return func(path string, fileInfo os.FileInfo, err error) error {

		modTime := fileInfo.ModTime()

		if err != nil {

			return postProcess(createUnprocessable(path, modTime, err))
		}

		if fileInfo.IsDir() {
			// not processing path to directory

			return nil
		}

		imgFile, err := os.Open(path)
		if err != nil {

			return postProcess(createUnprocessable(path, modTime, err))
		}

		if isVideo(path) {
			// videos do not have exif metadata, so we are not even trying to load them

			if isLivePhoto(path) {

				return postProcess(createThrash(path, modTime, errors.New("live photo")))
			}

			return postProcess(createRegular(path, modTime))
		}

		exifData, err := exif.Decode(imgFile)
		if err != nil {
			// files without exif metadata are thrash, because iPhone camera photos have exif metadata

			return postProcess(createThrash(path, modTime, err))
		}

		return postProcess(createWithExif(path, exifData, modTime))
	}
}

func createThrash(path string, dateCreated time.Time, err error) *metadata {

	return &metadata{
		Status:       Thrash,
		ThrashReason: fmt.Sprint(err),
		Path:         path,
		Created:      &dateCreated,
	}
}

func createUnprocessable(path string, dateCreated time.Time, err error) *metadata {

	return &metadata{
		Status:              Unprocessable,
		UnprocessableReason: fmt.Sprint(err),
		Path:                path,
		Created:             &dateCreated,
	}
}

func createRegular(path string, dateCreated time.Time) *metadata {

	return &metadata{
		Status:  Regular,
		Path:    path,
		Created: &dateCreated,
	}
}

func createWithExif(path string, exifData *exif.Exif, dateCreatedFallback time.Time) *metadata {
	make := getExifField(exifData, exif.Make)
	model := getExifField(exifData, exif.Model)
	areMakeAndModelExpected := make == expectedMake && model == expectedModel

	var status FileStatus
	var thrashReason string
	if areMakeAndModelExpected {
		status = Regular
	} else {
		status = Thrash
		thrashReason = "unexpected camera make or model"
	}

	return &metadata{
		Status:       status,
		ThrashReason: thrashReason,
		Path:         path,
		Make:         make,
		Model:        model,
		Created:      getExifTimeField(exifData, exif.DateTime, dateCreatedFallback),
	}
}

func getExifField(metaData *exif.Exif, fieldName exif.FieldName) string {
	fieldValue, err := metaData.Get(fieldName)
	if err != nil {
		log.Printf("\terror parsing field %s: %s", fieldName, err)

		return ""
	}

	fieldStringValue, err := fieldValue.StringVal()
	if err != nil {
		log.Printf("\terror converting value of %s to string: %s", fieldName, err)

		return ""
	}

	return fieldStringValue
}

func getExifTimeField(metaData *exif.Exif, fieldName exif.FieldName, fallbackTime time.Time) *time.Time {
	fieldStr := getExifField(metaData, fieldName)

	fieldDate, err := time.ParseInLocation(exifDateTimeFormat, fieldStr, time.Local)

	if err != nil {
		log.Printf("\terror converting value of %s to string: %s, going to use fallback time %s", fieldName, err, fallbackTime)

		return &fallbackTime
	}

	return &fieldDate
}

func isVideo(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), videoFileSuffix)
}

func isLivePhoto(path string) bool {
	filesWithSameNameAndVariousExtensions := path[:len(path)-len(videoFileSuffix)] + "*"
	matches, err := filepath.Glob(filesWithSameNameAndVariousExtensions)

	if err != nil {
		log.Printf("error evaluating if video is live photo, treating it as regular video for safety reasons: %s\n", path)
		log.Println(err)

		return false
	}

	if len(matches) > 1 {
		// video is live photo, because regular photo with the same name (not extension) exists besides the video

		return true
	}

	return false
}
