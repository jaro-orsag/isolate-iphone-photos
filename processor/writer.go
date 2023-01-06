package processor

import "log"

func WriteFile(metadata *metadata) error {
	log.Printf("Writer: %#v", metadata)

	return nil
}
