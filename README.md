# sort-iphone-photos
Command line utility that helps with separation of iPhone camera photos and videos from other media on iPhone. 

Media that are separated out are
* Whatsapp media
* iPhone screenshots
* Live photo videos
* All the pictures without metadata
* All the pictures with metadata, but with unexpected make and model. Expected make is `Apple` and expected model is 
`iPhone 12 mini`.

Use with care and love. And double check the results.

## How it works?
1. Import pictures from your iPhone to photo library. See one of the following articles for more info
    * https://mackeeper.com/blog/how-to-import-photos-iphone-to-mac/
    * https://support.apple.com/guide/photos/create-additional-libraries-pht6d60b524/mac

2. Export pictures from your photo library
    * Ideally export unmodified originals
    * Export with `Subfolder Format: None`

3. Run `./sort-iphone-photos <input-folder>`

4. `sort-iphone-photos` will produce following folder structure based on picture metedata
    * `random-guid` - subfolder of `input-folder`
        * `Regular` - pictures and videos from iPhone without excluded media
            * `yyyy-mm-dd`
            * the pictures
        * `Thrash` - excluded media such as live photo videos, screenshots, WhatsApp files, etc..
            * `yyyy-mm-dd`
            * the pictures

## Running `sort-iphone-photos`
First you have to install golang on your machine https://go.dev/dl/

Let's assume our photo library export is located in `../_examples/photo-library-export` folder.

To run the program
```
go run sort-iphone-photos.go ../_examples/photo-library-export
```

To build and run the binary
```
go build sort-iphone-photos.go
./sort-iphone-photos ../_examples/photo-library-export
```

# Roadmap
Functional
* Live photo videos have separate status and subfolder
* Dry run
* Verbose and silent mode
* Counters and statistics
* Configurable output folder
* Use moment name from media library in target folder names
* Include simple usage gif in this readme

Non functional
* Use `github.com/urfave/cli/v2` for handling of command line params
* Spaces instead of tabs
* Cover with unit tests. Shame on me - I should have implemented tests first. That would also affect the design.
* Enable force pushing to branches
* Extract some of the constants to cmd flags with defaults
    * output date format
    * video file extension
    * expected device make and model
* Make the code more [golang idiomatic](https://go.dev/doc/effective_go)
* Introduce concurrency and channels for communication between collector and writer