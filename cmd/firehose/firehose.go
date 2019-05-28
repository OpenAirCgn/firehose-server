package main

import "fmt"

import firehose "github.com/openaircgn/firehose_server"

func main () {
	fmt.Printf("hello");
	firehose.TCPRun();
}
