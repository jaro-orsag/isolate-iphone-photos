package processor

import (
	"log"
	filepath "path/filepath"

	"github.com/google/uuid"
)

type Args struct {
	InputFolder         string
	OutputFolder        string
	ExpectedCameraMake  string
	ExpectedCameraModel string
	VideoFileExtension  string
	TraverseDirTree     TraverseDirTreeFunc
	MakeCollectMetadata MakeCollectMetadataFunc
	MakeWriteFile       MakePostProcessFunc
}

func Run(args *Args) {
	var outputFolder string
	if args.OutputFolder != "" {
		outputFolder = args.OutputFolder
	} else {
		outputFolder = filepath.Join(args.InputFolder, uuid.NewString())
	}

	collectMetadata := args.MakeCollectMetadata(&makeCollectMetadataArgs{
		ExpectedCameraMake:  args.ExpectedCameraMake,
		ExpectedCameraModel: args.ExpectedCameraModel,
		VideoFileExtension:  args.VideoFileExtension,
		PostProcess:         args.MakeWriteFile(outputFolder),
	})

	err := args.TraverseDirTree(args.InputFolder, collectMetadata)

	if err != nil {
		log.Panicf("error traversing directory tree %s", err)
	}
}
