package processor

import (
	filepath "path/filepath"
	"time"
)

type TraverseDirTreeFunc func(root string, fn filepath.WalkFunc) error

type MakeCollectMetadataFunc func(postProcessFunc PostProcessFunc) filepath.WalkFunc

type PostProcessFunc func(metadata *metadata) error

type FileStatus int

const (
	Unprocessable FileStatus = iota
	Thrash
	Regular
)

func (fs FileStatus) String() string {
	return [...]string{"Unprocessable", "Thrash", "Regular"}[fs-1]
}

func (fs FileStatus) EnumIndex() int {
	return int(fs)
}

type metadata struct {
	Status              FileStatus
	UnprocessableReason string
	ThrashReason        string
	Path                string
	Make                string
	Model               string
	Created             *time.Time
}
