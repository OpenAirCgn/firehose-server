package main

import (
	"flag"
	"fmt"
	"os"
)
import firehose "github.com/openaircgn/firehose_server"

var outfile string
var addr string
var help bool

const (
	defaultAddr    = ":7531"
	defaultOutFile = "-"
)

func init() {
	flag.StringVar(&addr, "a", defaultAddr, "address for server to listen on")
	flag.StringVar(&addr, "addr", defaultAddr, "address for server to listen on")
	flag.StringVar(&outfile, "o", defaultOutFile, "filename to save output to")
	flag.StringVar(&outfile, "outfile", defaultOutFile, "filename to save output to")
	flag.BoolVar(&help, "h", false, "print usage")
	flag.BoolVar(&help, "help", false, "print usage")
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	var w *firehose.CSVMsgWriter
	if outfile == defaultOutFile {
		w = firehose.NewCSVWriter(os.Stdout)
	} else {
		f, err := os.OpenFile(outfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			flag.Usage()
			os.Exit(1)
		}
		defer f.Close()
		w = firehose.NewCSVWriter(f)
	}
	server := firehose.TCPServer{
		Address: addr,
		Writer:  w.DumpCSV,
	}
	server.Run()
}
