package firehose_server

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

var writer *csv.Writer = csv.NewWriter(os.Stdout)

func uInt32ToA(v uint32) string {
	return strconv.FormatInt(int64(v), 10)
}

func DumpCSV(msg Msg) {
	_time := time.Now().Unix()
	time := strconv.FormatInt(_time, 10)
	ts := uInt32ToA(msg.Timestamp)
	dev := msg.DeviceId
	tag := fmt.Sprintf("0x%02x", msg.Tag)
	value := fmt.Sprintf("0x%08x", msg.Value)

	record := []string{
		time, ts, dev, tag, value,
	}

	writer.Write(record)

}
