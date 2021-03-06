// See LICENSE file for copyright and license details.

package main

import (
	"archive/tar"
	"archive/zip"
	"github.com/hullerob/go.imagefile"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"
	"time"
)

type convFunc func(zf *zip.File, msg chan<- message) ifFile

func fileConv(ifch chan<- ifFile, zipch <-chan *zip.File, msg chan<- message) {
	defer close(ifch)
	for zf := range zipch {
		conv := getConvFuncByName(zf.Name)
		tf := conv(zf, msg)
		ifch <- tf
	}
}

func getConvFuncByName(name string) convFunc {
	lname := strings.ToLower(name)
	switch {
	case strings.HasSuffix(lname, ".png"):
		return imgConv(3)
	case flagJpg && strings.HasSuffix(lname, ".jpg"):
		return imgConv(3)
	case flagJpg && strings.HasSuffix(lname, ".jpeg"):
		return imgConv(4)
	default:
		return noConv
	}
}

func noConv(zf *zip.File, msg chan<- message) ifFile {
	tfh := tar.Header{
		Name:       zf.Name,
		Size:       int64(zf.UncompressedSize64),
		ChangeTime: time.Now(),
		AccessTime: time.Now(),
		ModTime:    time.Now(),
		Typeflag:   tar.TypeReg,
		Mode:       0644,
	}
	zfr, err := zf.Open()
	if err != nil {
		msg <- message{msg: msgErrZip, err: err}
	}
	tf := ifFile{
		header: tfh,
		reader: zfr,
		format: "copy",
	}
	return tf
}

func imgConv(sufLen int) convFunc {
	return func(zf *zip.File, msg chan<- message) ifFile {
		tf := noConv(zf, msg)
		img, format, err := image.Decode(tf.reader)
		if err != nil {
			msg <- message{
				msg:  msgErrConv,
				err:  err,
				file: tf.header.Name,
			}
			tf.reader.Close()
			return noConv(zf, msg)
		}
		nl := len(tf.header.Name) - sufLen
		tf.header.Name = tf.header.Name[:nl] + "if"
		size := img.Bounds().Dx()*img.Bounds().Dy()*4 + 17
		r, w := io.Pipe()
		go func() {
			imagefile.Encode(w, img)
			w.Close()
		}()
		tf.header.Size = int64(size)
		tf.reader = r
		tf.format = format
		return tf
	}
}
