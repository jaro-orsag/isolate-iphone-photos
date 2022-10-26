# sort-iphone-photos
Command line utility that helps with separation of iPhone camera photos from other pictures on iPhone. For example 
pictures from WhatsApp are separated out.

## How it works?
1. Import pictures from your iPhone to photo library. See one of the following articles for more info
    * https://mackeeper.com/blog/how-to-import-photos-iphone-to-mac/
    * https://support.apple.com/guide/photos/create-additional-libraries-pht6d60b524/mac

2. Export pictures from your photo library
    * Ideally export unmodified originals
    * Export with `Subfolder Format: None`. `sort-iphone-photos` will create it's own folder structure.

3. [TODO] Run `sort-iphone-photos`

4. `sort-iphone-photos` will produce following folder structure based on picture metedata
    * `device_name` - name of device that captured the picture
      * `yyyy-mm-dd-moment-name`
        * the picture itself

## Running `sort-iphone-photos`

First you have to install golang on your machine https://go.dev/dl/

And install dependencies
```
$ go get github.com/rwcarlsen/goexif/exif
$ go get github.com/tidwall/gjson
```

To run the program
```
$ go run sort-iphone-photos.go
```

To build and run the binary
```
$ go build sort-iphone-photos.go
$ ./sort-iphone-photos
```



Make: "Apple"
DateTimeOriginal: "2022:07:28 12:02:22"
LensModel: "iPhone 12 mini back dual wide camera 4.2mm f/1.6"
DateTimeDigitized: "2022:07:28 12:02:22"
LensMake: "Apple"
Model: "iPhone 12 mini"