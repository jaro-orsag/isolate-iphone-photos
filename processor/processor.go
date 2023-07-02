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

	thrash := make(chan int)
	regular := make(chan int)
	livePhotoVideo := make(chan int)
	quit := make(chan int)

	collectMetadata := args.MakeCollectMetadata(&makeCollectMetadataArgs{
		ExpectedCameraMake:  args.ExpectedCameraMake,
		ExpectedCameraModel: args.ExpectedCameraModel,
		VideoFileExtension:  args.VideoFileExtension,
		PostProcess:         args.MakeWriteFile(outputFolder, thrash, regular, livePhotoVideo),
	})

	go func() {
		err := args.TraverseDirTree(args.InputFolder, collectMetadata)
		if err != nil {
			log.Panicf("error traversing directory tree %s", err)
		}

		quit <- 0
	}()

	calculateStats(thrash, regular, livePhotoVideo, quit)
}

func calculateStats(thrash chan int, regular chan int, livePhotoVideo chan int, quit chan int) {
	thrashCount := 0
	regularCount := 0
	livePhotoVideoCount := 0

	for {
		select {
		case <-thrash:
			thrashCount++
		case <-regular:
			regularCount++
		case <-livePhotoVideo:
			livePhotoVideoCount++
		case <-quit:
			log.Printf("thrash: %v", thrashCount)
			log.Printf("regular: %v", regularCount)
			log.Printf("livePhotoVideo: %v", livePhotoVideoCount)
			
			return
		}
	}
}
