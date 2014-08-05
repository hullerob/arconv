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

var (
	flagVerbose bool
	flagJpg     bool
)

func init() {
	flag.BoolVar(&flagVerbose, "v", false, "print progress")
	flag.BoolVar(&flagJpg, "jpg", false, "convert jpg files")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options] archive.zip\n", os.Args[0])
	flag.PrintDefaults()
}
