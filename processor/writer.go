package processor

import (
	"log"
	"path/filepath"

	cp "github.com/otiai10/copy"
)

// TODO: following constants should be generalized or configurable
const dateFormat = "2006-01-02"

// TODO: implement shared counters and stats
func MakeWriteFile(outputRoot string) PostProcessFunc {

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

		return nil
	}
}
