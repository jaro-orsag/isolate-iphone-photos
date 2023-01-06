package processor

import (
	"log"
	filepath "path/filepath"

	"github.com/google/uuid"
)

type Args struct {
	Root                string
	TraverseDirTree     TraverseDirTreeFunc
	MakeCollectMetadata MakeCollectMetadataFunc
	MakeWriteFile       MakePostProcessFunc
}

func Run(args *Args) {
	outputRoot := filepath.Join(args.Root, uuid.NewString())
	writeFile := args.MakeWriteFile(outputRoot)
	collectMetadata := args.MakeCollectMetadata(writeFile)

	err := args.TraverseDirTree(args.Root, collectMetadata)

	if err != nil {
		log.Panicf("error traversing directory tree %s", err)
	}
}
