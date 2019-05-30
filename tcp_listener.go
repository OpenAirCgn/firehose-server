package firehose_server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

type TCPServer struct {
	Address string
	Writer  MsgWriter
}

func (s *TCPServer) handleTCPConnection(c net.Conn) {
	decoder := json.NewDecoder(c)
	var msg Msg
	for decoder.More() {
		err := decoder.Decode(&msg)
		if err != nil {
			log.Fatal(err)
		}
		s.Writer(msg)
	}
}

func (s *TCPServer) Run() {
	listener, err := net.Listen("tcp", s.Address)
	fmt.Fprintf(os.Stderr, "Listening on %s\n", s.Address)

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
		fmt.Fprintf(os.Stderr, "Received connection: %v\n", con)
		go s.handleTCPConnection(con)
	}
}
