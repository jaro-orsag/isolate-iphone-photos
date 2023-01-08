package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	processor "isolate-iphone-photos/processor"
)

const (
	inputFolderFlagName         = "inputFolder"
	outputFolderFlagName        = "outputFolder"
	expectedCameraMakeFlagName  = "expectedCameraMake"
	expectedCameraModelFlagName = "expectedCameraModel"
	videoFileExtensionFlagName  = "videoFileExtension"
)

func main() {
	log.SetFlags(log.Ltime)

	app := &cli.App{
		Name: "isolate-iphone-photos",
		Usage: "Command line utility that helps with separation of iPhone camera media from other media " +
			"on iPhone. It separates them by generating following folder structure: \n\n" +

			" - Regular \n" +
			"     - YYYY-MM-DD \n" +
			"         - iPhone camera photos, including selfies\n" +
			"         - iPhone camera videos, including selfies\n\n" +

			" - LivePhotoVideos \n" +
			"     - YYYY-MM-DD \n" +
			"         - Live photo videos\n\n" +

			" - Thrash \n" +
			"     - YYYY-MM-DD \n" +
			"         - Whatsapp media\n" +
			"         - iPhone screenshots\n" +
			"         - All the pictures without metadata\n" +
			"         - All the pictures with metadata, but with unexpected make and model. Make and model " +
			"can be configured using flags.\n\n" +

			"This project has two goals:\n\n" +

			" - To fulful the needs of its author related to processing of iPhone photos \n" +
			" - It is an attempt to practice golang. Therefore, please use it with care and love. And double check " +
			"the results.\n\n" +

			"For more info see https://github.com/jaro-orsag/isolate-iphone-photos.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     inputFolderFlagName,
				Aliases:  []string{"i"},
				Usage:    "Input folder. Should contain media library export. Required.",
				Required: true,
			},
			&cli.StringFlag{
				Name:    outputFolderFlagName,
				Aliases: []string{"o"},
				Usage: "Output folder. If not provided, random-guid-named folder is created under " +
					"input folder and is used as output folder.",
			},
			&cli.StringFlag{
				Name:    expectedCameraMakeFlagName,
				Aliases: []string{"make"},
				Usage:   "Only picutres that have expected exif metadata 'make' are copied to 'Regular' folder.",
				Value:   "Apple",
			},
			&cli.StringFlag{
				Name:    expectedCameraModelFlagName,
				Aliases: []string{"model"},
				Usage:   "Only picutres that have expected exif metadata 'model' are copied to 'Regular' folder.",
				Value:   "iPhone 12 mini",
			},
			&cli.StringFlag{
				Name:    videoFileExtensionFlagName,
				Aliases: []string{"video"},
				Usage:   "Extension used by video files.",
				Value:   ".mov",
			},
		},
		Action: func(c *cli.Context) error {
			log.Printf("processing photo library export from %s\n", c.String(inputFolderFlagName))

			processor.Run(&processor.Args{
				InputFolder:         c.String(inputFolderFlagName),
				OutputFolder:        c.String(outputFolderFlagName),
				ExpectedCameraMake:  c.String(expectedCameraMakeFlagName),
				ExpectedCameraModel: c.String(expectedCameraModelFlagName),
				VideoFileExtension:  c.String(videoFileExtensionFlagName),
				TraverseDirTree:     filepath.Walk,
				MakeCollectMetadata: processor.MakeCollectMetadata,
				MakeWriteFile:       processor.MakeWriteFile,
			})

			log.Print("finished processing photo library export")

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
