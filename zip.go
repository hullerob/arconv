// See LICENSE file for copyright and license details.

package main

import (
	"archive/zip"
	"sort"
)

type sorter []*zip.File

func (s sorter) Len() int {
	return len(s)
}

func (s sorter) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func (s sorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func readZip(name string, files chan<- *zip.File, msg chan<- message) {
	defer close(files)
	zr, err := zip.OpenReader(name)
	if err != nil {
		msg <- message{msg: msgErrZip, err: err}
		return
	}
	msg <- message{msg: msgTotalFiles, ival: len(zr.File)}
	s := sorter(zr.File)
	sort.Sort(s)
	for _, f := range s {
		files <- f
	}
}
