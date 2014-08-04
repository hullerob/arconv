// See LICENSE file for copyright and license details.

package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	zipch := make(chan *zip.File)
	msg := make(chan message)
	go readZip(flag.Args()[0], zipch, msg)
	ifch := make(chan ifFile)
	go fileConv(ifch, zipch, msg)
	go writeTar(os.Stdout, ifch, msg)
	ret := progress(msg)
	os.Stdout.Sync()
	os.Stdout.Close()
	if ret != 0 {
		os.Exit(ret)
	}
}

const (
	msgDone = iota
	msgTotalFiles
	msgProcFile
	msgErrZip
	msgErrConv
	msgErrTar
)

type message struct {
	msg  int
	ival int
	sval string
	err  error
}

func progress(messages <-chan message) int {
	tfiles := 0
	pfiles := 0
	var procfile string
LOOP:
	msg := <-messages
	switch msg.msg {
	case msgDone:
		fmt.Fprintf(os.Stderr, "\nall done\n")
		return 0
	case msgTotalFiles:
		tfiles = msg.ival
	case msgProcFile:
		pfiles++
		procfile = msg.sval
	case msgErrZip:
		fmt.Fprintf(os.Stderr, "\r%s: %v\033[K\n", os.Args[0], msg.err)
		fmt.Fprintf(os.Stderr, "error not recoverable, exiting\n")
		return 1
	case msgErrConv:
		fmt.Fprintf(os.Stderr, "\r%s: error while converting file '%s'\033[K\n", os.Args[0], msg.sval)
		fmt.Fprintf(os.Stderr, "%v\n", msg.err)
	case msgErrTar:
		fmt.Fprintf(os.Stderr, "\r%s: %v\033[K\n", os.Args[0], msg.err)
		fmt.Fprintf(os.Stderr, "error not recoverable, exiting\n")
		return 1
	default:
		panic("erroneous message")
	}
	if len(procfile) > 0 {
		fmt.Fprintf(os.Stderr, "\r[%4d/%4d] converting file '%s'\033[K", pfiles, tfiles, procfile)
	} else {
		fmt.Fprintf(os.Stderr, "\r[%4d/%4d]\033[K", pfiles, tfiles)
	}
	goto LOOP
}

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options] archive.zip\n", os.Args[0])
	flag.PrintDefaults()
}
