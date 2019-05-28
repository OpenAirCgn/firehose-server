package firehose_server

import (
	"encoding/json"
	"log"
	"net"
)

func handleTCPConnection(c net.Conn) {
	decoder := json.NewDecoder(c)
	var msg Msg
	for decoder.More() {
		err := decoder.Decode(&msg)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("-> %v\n", msg)
		DumpCSV(msg)
	}
}
func TCPRun() {
	listener, err := net.Listen("tcp", ":7531")
	println("Listening ... ")

	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	var con net.Conn
	for {
		con, err = listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		println("received connection")
		go handleTCPConnection(con)
	}
}
