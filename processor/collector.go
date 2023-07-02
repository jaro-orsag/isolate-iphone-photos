package processor

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	filepath "path/filepath"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)



const (
	exifDateTimeFormat = "2006:01:02 15:04:05"
	heicExtension      = ".heic"
)

type fileStatus int

const (
	Unprocessable fileStatus = iota
	Thrash
	Regular
	LivePhotoVideo
)

func (fs fileStatus) String() string {
	return [...]string{"Unprocessable", "Thrash", "Regular", "LivePhotoVideo"}[fs]
}

func (fs fileStatus) EnumIndex() int {
	return int(fs)
}

type metadata struct {
	Status              fileStatus
	UnprocessableReason string
	ThrashReason        string
	Path                string
	Make                string
	Model               string
	Created             *time.Time
}

type makeCollectMetadataArgs struct {
	ExpectedCameraMake  string
	ExpectedCameraModel string
	VideoFileExtension  string
	PostProcess         PostProcessFunc
}

func MakeCollectMetadata(args *makeCollectMetadataArgs) filepath.WalkFunc {
	exif.RegisterParsers(mknote.All...)

	return func(path string, fileInfo os.FileInfo, err error) error {

		modTime := fileInfo.ModTime()

		if err != nil {

			return args.PostProcess(createUnprocessable(path, modTime, err))
		}

		if fileInfo.IsDir() {
			// not processing path to directory

			return nil
		}

		imgFile, err := os.Open(path)
		if err != nil {

			return args.PostProcess(createUnprocessable(path, modTime, err))
		}

		if isVideo(path, args.VideoFileExtension) {
			// videos do not have exif metadata, so we are not even trying to load exif metadata

			if isLivePhoto(path, args.VideoFileExtension) {

				return args.PostProcess(createLivePhotoVideo(path, modTime))
			}

			return args.PostProcess(createRegular(path, modTime))
		}

		if isHeic(path) {
			// exif lib we are using does not have exif support, so we assume heic images are from expected cammera

			return args.PostProcess(createRegular(path, modTime))
		}

		exifData, err := exif.Decode(imgFile)
		if err != nil {
			// files without exif metadata are thrash, because iPhone camera photos have exif metadata

			return args.PostProcess(createThrash(path, modTime, err))
		}

		return args.PostProcess(createWithExif(&createWithExifArgs{
			Path:                path,
			ExifData:            exifData,
			DateCreatedFallback: modTime,
			ExpectedCameraMake:  args.ExpectedCameraMake,
			ExpectedCameraModel: args.ExpectedCameraModel,
		}))
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

func createLivePhotoVideo(path string, dateCreated time.Time) *metadata {

	return &metadata{
		Status:  LivePhotoVideo,
		Path:    path,
		Created: &dateCreated,
	}
}

type createWithExifArgs struct {
	Path                string
	ExifData            *exif.Exif
	DateCreatedFallback time.Time
	ExpectedCameraMake  string
	ExpectedCameraModel string
}

func createWithExif(args *createWithExifArgs) *metadata {
	make := getExifField(args.ExifData, exif.Make)
	model := getExifField(args.ExifData, exif.Model)
	areMakeAndModelExpected := make == args.ExpectedCameraMake && model == args.ExpectedCameraModel

	var status fileStatus
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
		Path:         args.Path,
		Make:         make,
		Model:        model,
		Created:      getExifTimeField(args.ExifData, exif.DateTime, args.DateCreatedFallback),
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

func isVideo(path string, videoFileExt string) bool {
	return strings.ToLower(filepath.Ext(path)) == videoFileExt
}

func isLivePhoto(path string, videoFileExt string) bool {
	filesWithSameNameAndVariousExtensions := path[:len(path)-len(videoFileExt)] + "*"
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

func isHeic(path string) bool {

	return strings.ToLower(filepath.Ext(path)) == heicExtension
}
