package fixer

import (
	"log"
	"os"
	"time"
	"syscall"
	"path/filepath"
)

func FixupMetadata(inputFolder string, lookupFolder string) error {

	filepath.Walk(inputFolder, func(inputPath string, inputFileInfo os.FileInfo, err error) error {

		if inputFileInfo.IsDir() {

			return nil
		}

		filepath.Walk(lookupFolder, func(lookupPath string, lookupFileInfo os.FileInfo, err error) error {

			if lookupFileInfo.IsDir() {

				return nil
			}

			if (lookupFileInfo.Name() == inputFileInfo.Name()) {
				log.Printf("processing %v", inputFileInfo.Name())
				
				stat := lookupFileInfo.Sys().(*syscall.Stat_t)
				aTime := time.Unix(stat.Atimespec.Sec, stat.Atimespec.Nsec)
				modTime := lookupFileInfo.ModTime

				err := os.Chtimes(inputPath, aTime, modTime())
				if err != nil {
					log.Printf("%v", err)
					
					return nil
				}
			}

			return nil
		})
		
		return nil
	})

	return nil
}