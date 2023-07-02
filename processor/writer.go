package processor

import (
	"log"
	"path/filepath"

	cp "github.com/otiai10/copy"
)

const dateFormat = "2006-01-02"

func MakeWriteFile(outputRoot string, thrash chan int, regular chan int, livePhotoVideo chan int) PostProcessFunc {

	return func(metadata *metadata) error {
		log.Printf("%#v", metadata)

		if metadata.Status == Unprocessable {
			log.Print("\tnot processing")

			return nil
		}

		created := metadata.Created.Format(dateFormat)
		_, fileName := filepath.Split(metadata.Path)
		targetPath := filepath.Join(outputRoot, metadata.Status.String(), created, fileName)

		log.Printf("\ttargetPath %s", targetPath)

		err := cp.Copy(metadata.Path, targetPath)
		if err != nil {
			log.Println(err)

			return err
		}

		if metadata.Status == Thrash {
			thrash <- 0
		} else if  metadata.Status == Regular {
			regular <- 0
		} else if  metadata.Status == LivePhotoVideo {
			livePhotoVideo <- 0
		}
		
		return nil
	}
}
