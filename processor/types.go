package processor

import (
	filepath "path/filepath"
)

type TraverseDirTreeFunc func(root string, fn filepath.WalkFunc) error

type MakeCollectMetadataFunc func(*makeCollectMetadataArgs) filepath.WalkFunc

type PostProcessFunc func(*metadata) error

type MakePostProcessFunc func(root string) PostProcessFunc
