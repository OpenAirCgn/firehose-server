package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	mrand "math/rand"
	"net"
	"os"
	"time"

	fh "github.com/openaircgn/firehose_server"
)

var firehoseAddr = flag.String("firehose_addr", "localhost:7531", "address and port to connect to")
var printVersion = flag.Bool("version", false, "print version info and exit")

var version string = "unknown"

var deviceId string = generateDeviceId()
var startTime = time.Now()

func generatePacket() []byte {
	guard := fh.OA_FINAL_SPECIAL_GUARD
	msg := fh.Msg{}
	msg.Timestamp = uint32(time.Since(startTime) * time.Second)
	msg.DeviceId = deviceId
	msg.Tag = fh.Tag(mrand.Int31n(int32(guard + 1)))
	msg.Value = mrand.Uint32()
	bs, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bs
}

func generateDeviceId() string {
	bs := make([]byte, 3)
	rand.Read(bs)
	return fmt.Sprintf("esp32_%X", bs)
}

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", *firehoseAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Printf("Welcome to Firehose(%s)!", version)
	if *printVersion {
		println()
		os.Exit(0)
	}

	fmt.Printf("Press Ctl-c to end\n")
	fmt.Printf("   using simulated device_id:%s\n", deviceId)
	fmt.Printf("   using firehose_addr:%s\n", *firehoseAddr)
	// generate MAC
	// save initial TS
	for {
		// generate packet
		packet := generatePacket()
		// send
		conn.Write(packet)
		// sleep
		time.Sleep(100 * time.Millisecond)
	}

}
