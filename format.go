package main

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
)

type SupportedFormat struct {
	Title      string
	Extensions []string
	Function   func(w io.ReadWriter, image image.Image, encode bool) (image.Image, error)
	Encodable  bool
}

var formats []SupportedFormat = []SupportedFormat{
	{Title: "Bitmap Files", Extensions: []string{".bmp"},
		Function: func(w io.ReadWriter, image image.Image, encode bool) (image.Image, error) {
			if encode {
				bmp.Encode(w, image)
				return image, nil
			} else {
				return bmp.Decode(w)
			}
		}, Encodable: true},
	{Title: "PNG", Extensions: []string{".png"},
		Function: func(w io.ReadWriter, image image.Image, encode bool) (image.Image, error) {
			if encode {
				png.Encode(w, image)
				return image, nil
			} else {
				return png.Decode(w)
			}
		}, Encodable: true},
	{Title: "JPEG", Extensions: []string{".jpg", ".jpeg"},
		Function: func(w io.ReadWriter, image image.Image, encode bool) (image.Image, error) {
			if encode {
				jpeg.Encode(w, image, &jpeg.Options{Quality: 100})
				return image, nil
			} else {
				return jpeg.Decode(w)
			}
		}, Encodable: true},
	{Title: "TIFF", Extensions: []string{".tif", ".tiff"},
		Function: func(w io.ReadWriter, image image.Image, encode bool) (image.Image, error) {
			if encode {
				tiff.Encode(w, image, &tiff.Options{Compression: tiff.CCITTGroup4})
				return image, nil
			} else {
				return tiff.Decode(w)
			}
		}, Encodable: true},
	{Title: "WEBP", Extensions: []string{".webp"},
		Function: func(w io.ReadWriter, image image.Image, encode bool) (image.Image, error) {
			if encode {
				return nil, errors.New("'webp' format encoding not supported")
			} else {
				return webp.Decode(w)
			}
		}, Encodable: false},
	{Title: "GIF", Extensions: []string{".gif"},
		Function: func(w io.ReadWriter, image image.Image, encode bool) (image.Image, error) {
			if encode {
				//gif.Encode(w, image, &gif.Options{})
				//return image, nil
				return nil, errors.New("'gif' format encoding not supported")
			} else {
				return gif.Decode(w)
			}
		}, Encodable: false},
}

func FindFormatFromExt(extension string) *SupportedFormat {
	for i := range formats {
		for _, ext := range formats[i].Extensions {
			if strings.EqualFold(ext, extension) {
				return &formats[i]
			}
		}
	}
	return nil
}

func FindFormatIndexFromExt(extension string) int {
	for i := range formats {
		for _, ext := range formats[i].Extensions {
			if strings.EqualFold(ext, extension) {
				return i
			}
		}
	}
	return 0
}

func GetFormatCount() int {
	return len(formats)
}

func GetDialogFilters(onlyEncodable bool) string {
	filter := ""
	for _, format := range formats {
		if onlyEncodable && !format.Encodable {
			continue
		}
		filter += format.Title + " (*"
		for j, ext := range format.Extensions {
			if j != 0 {
				filter += ",  *"
			}
			filter += ext
		}
		filter += ")|*"
		for j, ext := range format.Extensions {
			if j != 0 {
				filter += ";*"
			}
			filter += ext
		}
		filter += "|"
	}
	return filter
}

func GetSaveFileDialogFilters() string {
	return GetDialogFilters(true)
}

func GetOpenFileDialogFilters() string {
	filter := GetDialogFilters(false)
	filter += "All Picture Files|*"
	for i, format := range formats {
		if i != 0 {
			filter += ";*"
		}
		for j, ext := range format.Extensions {
			if j != 0 {
				filter += ";*"
			}
			filter += ext
		}
	}
	filter += "|All Files|*.*|"
	return filter
}
