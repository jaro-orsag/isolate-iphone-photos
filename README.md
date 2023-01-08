# isolate-iphone-photos
Command line utility that helps with separation of iPhone camera media from other media on iPhone. It separates them by generating following folder structure:
- Regular
    - YYYY-MM-DD
        - iPhone camera photos, including selfies
        - iPhone camera videos, including selfies

- LivePhotoVideos
    - YYYY-MM-DD
        - Live photo videos

- Thrash
    - YYYY-MM-DD
        - Whatsapp media
        - iPhone screenshots
        - All the pictures without metadata
        - All the pictures with metadata, but with unexpected make and model. Make and model can be configured using flags.

This project has two goals:
- To fulful the needs of its author related to processing of iPhone photos
- It is an attempt to practice golang. Therefore, please use it with care and love. And double check the results.

For more info see https://github.com/jaro-orsag/isolate-iphone-photos.

## Workflow
1. Import pictures from your iPhone to photo library. See one of the following articles for more info
    * https://mackeeper.com/blog/how-to-import-photos-iphone-to-mac/
    * https://support.apple.com/guide/photos/create-additional-libraries-pht6d60b524/mac

2. Export pictures from your photo library
    * Ideally export unmodified originals
    * Export with `Subfolder Format: None`

3. Run `isolate-iphone-photos`

## Running `isolate-iphone-photos`
First you have to install golang on your machine https://go.dev/dl/

Let's assume our photo library export is located in `../_examples/photo-library-export` folder.

To run the program
```
go run isolate-iphone-photos.go ../_examples/photo-library-export
```

To build and run the binary
```
go build isolate-iphone-photos.go
./isolate-iphone-photos ../_examples/photo-library-export
```

```
USAGE:
   isolate-iphone-photos [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --inputFolder value, -i value               Input folder. Should contain media library export. Required.
   --outputFolder value, -o value              Output folder. If not provided, random-guid-named folder is created under input folder and is used as output folder.
   --expectedCameraMake value, --make value    Only picutres that have expected exif metadata 'make' are copied to 'Regular' folder. (default: "Apple")
   --expectedCameraModel value, --model value  Only picutres that have expected exif metadata 'model' are copied to 'Regular' folder. (default: "iPhone 12 mini")
   --videoFileExtension value, --video value   Extension used by video files. (default: ".mov")
   --help, -h                                  show help (default: false)
```

# Roadmap
Functional
* Recognize multiple video format, so that we cover situation when output format is changed in iPhone settings
* Dry run
* Verbose and silent mode
* Counters and statistics
* Use moment name from media library in target folder names
* Include simple usage gif in this readme

Non functional
* Spaces instead of tabs
* Cover with unit tests. Shame on me - I should have implemented tests first. That would also affect the design.
* Enable force pushing to branches
* Make the code more [golang idiomatic](https://go.dev/doc/effective_go)
* Introduce concurrency and channels for communication between collector and writer