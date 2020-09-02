package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	firehose "github.com/openaircgn/firehose_server"
)

var outfile string
var addr string
var help bool
var printIP bool

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
		f, err := os.OpenFile(outfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			flag.Usage()
			os.Exit(1)
		}
		defer f.Close()
		w = firehose.NewCSVWriter(f, msgChan)
	}
	server := firehose.TCPServer{
		Address: addr,
		MsgChan: msgChan,
	}

	println("Welcome to Firehose! Press Ctl-c to end")

	go server.Run(doneChan)
	go w.Run(doneChan)

	for i := 0; i != 2; i++ {
		<-doneChan
	}
}
