package processor

import (
	"log"
)

type Args struct {
	Root                string
	TraverseDirTree     TraverseDirTreeFunc
	MakeCollectMetadata MakeCollectMetadataFunc
	WriteFile           PostProcessFunc
}

func Run(args *Args) {
	err := args.TraverseDirTree(args.Root, args.MakeCollectMetadata(args.WriteFile))

	if err != nil {
		log.Panicf("error traversing directory tree %s", err)
	}
}
