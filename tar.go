// See LICENSE file for copyright and license details.

package main

import (
	"archive/tar"
	"io"
)

type ifFile struct {
	header tar.Header
	reader io.ReadCloser
}

func writeTar(w io.Writer, ifch <-chan ifFile, msg chan<- message) {
	tw := tar.NewWriter(w)
	for file := range ifch {
		err := tw.WriteHeader(&file.header)
		if err != nil {
			msg <- message{msg: msgErrTar, err: err}
			return
		}
		_, err = io.Copy(tw, file.reader)
		if err != nil {
			msg <- message{msg: msgErrTar, err: err}
			return
		}
		file.reader.Close()
		msg <- message{msg: msgProcFile, sval: file.header.Name}
	}
	tw.Close()
	msg <- message{msg: msgDone}
}
