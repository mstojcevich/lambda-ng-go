// +build !windows

package upload

import(
	bimg "gopkg.in/h2non/bimg.v1"
)

var thumbnailOptions = bimg.Options{
	Width:     128,
	Height:    128,
	Crop:      true,
	Quality:   75,
	Interlace: true,
}

/**
Returns the bytes of the thumbnail file, or null if none was made
 */
func createThumb(b []byte) []byte {
	img := bimg.NewImage(b)
	thumb, err := img.Process(thumbnailOptions)

	if err != nil {
		return nil
	}

	return thumb
}