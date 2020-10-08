package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	firehose "github.com/openaircgn/firehose_server"
)

var outfile string
var addr string
var help bool
var printIP bool
var printVersion bool
var useDateTree bool

var logrotationIntervalMinutes int

var version string = "unknown" // set by linker, see xcompile.sh

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
	flag.BoolVar(&printIP, "printIP", false, "print IP Address(es) of host")
	flag.BoolVar(&printVersion, "version", false, "print version banner and exit")
	flag.BoolVar(&useDateTree, "dateDir", false, "rotate logs into YYYY/MM/DD directory structure")
	flag.BoolVar(&firehose.DontSkipUnknown, "dontSkipUnknown", false, "output value even if the tag is unknown")

	flag.IntVar(&logrotationIntervalMinutes, "csvAgeMinutes", 10, "after how many minutes to rotate the csv file")

}

func addrs() []string {
	itfs, err := net.Interfaces()
	if err != nil {
		panic(err.Error())
	}
	var ips []string
	for _, itf := range itfs {
		addrs, err := itf.Addrs()
		if err != nil {
			panic(err.Error())
		}
		for _, addr := range addrs {
			//fmt.Printf("> %v\n", addr)
			if ip, ok := addr.(*net.IPNet); ok {
				if !ip.IP.IsLoopback() && ip.IP.To4() != nil {
					ips = append(ips, ip.IP.String())
				}
			} else {
				println(ok)
			}
		}

	}
	return ips

}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}

	if printIP {
		ips := addrs()
		switch len(ips) {
		case 0:
			println("Can't determine IP address of host")
		case 1:
			println("IP of host:", ips[0])
		default:
			ips_list := strings.Join(ips, ",")
			println("IPs of host:", ips_list)
		}
		os.Exit(0)
	}

	msgChan := make(chan firehose.Msg)
	doneChan := make(chan bool)

	var w *firehose.CSVMsgWriter
	if outfile == defaultOutFile {
		w = firehose.NewCSVWriter(os.Stdout, msgChan)
	} else {
		//f, err := os.OpenFile(outfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		writer := firehose.Logrotation{
			BaseFilename: outfile,
			UseDateTree:  useDateTree,
			Interval:     time.Duration(logrotationIntervalMinutes) * time.Minute,
		}
		banner := fmt.Sprintf("# Firehose ver: %s Starttime: %v\n", version, time.Now())
		_, err := writer.Write([]byte(banner))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			flag.Usage()
			os.Exit(1)
		}
		defer writer.Close()
		w = firehose.NewCSVWriter(&writer, msgChan)
	}
	server := firehose.TCPServer{
		Address: addr,
		MsgChan: msgChan,
	}

	fmt.Printf("Welcome to Firehose(%s)!", version)
	if printVersion {
		println()
		os.Exit(0)
	}

	fmt.Printf("Listening on: %s. Press Ctl-c to end\n", addr)

	go server.Run(doneChan)
	go w.Run(doneChan)

	for i := 0; i != 2; i++ {
		<-doneChan
	}
}
