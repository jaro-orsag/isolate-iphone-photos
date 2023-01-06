package processor

import "log"

func WriteFile(metadata *metadata) error {
	log.Printf("%#v", metadata)

	return nil
}
