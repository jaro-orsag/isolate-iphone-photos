package processor

import filepath "path/filepath"

type TraverseDirTreeFunc func(root string, fn filepath.WalkFunc) error

type MakeCollectMetadataFunc func(postProcessFunc PostProcessFunc) filepath.WalkFunc

type PostProcessFunc func(metadata *metadata) error

type metadata struct {
	Path     string
	Make     string
	Model    string
	DateTime string
}
