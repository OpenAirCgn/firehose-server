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
	MsgChan chan Msg
}

func (s *TCPServer) handleTCPConnection(c net.Conn) {
	defer c.Close()
	decoder := json.NewDecoder(c)
	var msg Msg
	networkMsg := Msg{Tag: OA_Network_Events, Value: uint32(CONNECT)}

	for decoder.More() {
		err := decoder.Decode(&msg)
		if err != nil {
			log.Print(err)
			return
		}
		if networkMsg.Value == uint32(CONNECT) {
			networkMsg.DeviceId = msg.DeviceId
			s.MsgChan <- networkMsg
			networkMsg.Value = uint32(DISCONNECT)
		}
		s.MsgChan <- msg
	}

	s.MsgChan <- networkMsg
}

func (s *TCPServer) Run(doneChan chan<- bool) {
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
			log.Print(err)
			close(s.MsgChan)
			doneChan <- true
			return
		}
		fmt.Fprintf(os.Stderr, "Received connection: %v\n", con)
		go s.handleTCPConnection(con)
	}
}
