package firehose_server


import (
	"net"
	"log"
	"fmt"
	"encoding/json"
	"bytes"
)

func Run () {
	con, err := net.ListenPacket("udp", ":7531");
	if err != nil {
		log.Fatal(err);
	}
	defer con.Close()

	println("here")
	var n int
	var addr net.Addr

	var msg Msg
		bs := make([]byte, 1024)
	for {
		n, addr, err = con.ReadFrom(bs)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("addr: %v bs: ?? num: %d\n", addr, n)
		reader := bytes.NewReader(bs[0:n])
		decoder := json.NewDecoder(reader)
		for decoder.More() {
			err = decoder.Decode(&msg)
			if (err != nil) {
				fmt.Printf("-> %v\n", (string)(bs[0:n]))
				
				log.Fatal(err)
			}
		fmt.Printf("-> %v\n", msg)

		}
	}
}
