package main

import (
	"fmt"
	"os"
)

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
	file string
	conv string
	err  error
}

func progress(messages <-chan message) int {
	tfiles := 0
	pfiles := 0
	var procfile string
	conv := make(map[string]int)
LOOP:
	msg := <-messages
	switch msg.msg {
	case msgDone:
		printDone(conv)
		return 0
	case msgTotalFiles:
		tfiles = msg.ival
	case msgProcFile:
		pfiles++
		procfile = msg.file
		conv[msg.conv]++
	case msgErrZip:
		printError("%v", msg.err)
		printError("error not recoverable, exiting")
		return 1
	case msgErrConv:
		printError("error while converting file '%s'", msg.file)
		printError("%v", msg.err)
	case msgErrTar:
		printError("%v", msg.err)
		printError("error not recoverable, exiting")
		return 1
	default:
		panic("erroneous message")
	}
	printProgress(pfiles, tfiles, procfile)
	goto LOOP
}

func printError(format string, arg ...interface{}) {
	if flagVerbose {
		fmt.Fprint(os.Stderr, "\r\033[K")
	}
	msg := fmt.Sprintf(format, arg...)
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], msg)
}

func printProgress(current, total int, file string) {
	if !flagVerbose {
		return
	}
	if len(file) > 0 {
		fmt.Fprintf(os.Stderr, "\r[%4d/%4d] converting file '%s'\033[K", current, total, file)
	} else {
		fmt.Fprintf(os.Stderr, "\r[%4d/%4d]\033[K", current, total)
	}
}

func printDone(conv map[string]int) {
	if !flagVerbose {
		return
	}
	fmt.Fprintf(os.Stderr, "\nconverions:\n")
	for k, v := range conv {
		if k == "copy" {
			continue
		}
		fmt.Fprintf(os.Stderr, "%s: %d\n", k, v)
	}
	fmt.Fprintf(os.Stderr, "no conversion: %d\n", conv["copy"])
	fmt.Fprintf(os.Stderr, "all done\n")
}
