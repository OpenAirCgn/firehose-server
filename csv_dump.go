package firehose_server

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

type MsgWriter func(Msg)

type CSVMsgWriter struct {
	writer *csv.Writer
}

func NewCSVWriter(w io.Writer) *CSVMsgWriter {
	writer := new(CSVMsgWriter)
	writer.writer = csv.NewWriter(w)
	return writer
}

func uInt32ToA(v uint32) string {
	return strconv.FormatInt(int64(v), 10)
}

func (w *CSVMsgWriter) DumpCSV(msg Msg) {
	_time := time.Now().Unix()
	time := strconv.FormatInt(_time, 10)
	ts := uInt32ToA(msg.Timestamp)
	dev := msg.DeviceId
	tag := fmt.Sprintf("0x%02x", msg.Tag)
	value := fmt.Sprintf("0x%08x", msg.Value)

	record := []string{
		time, ts, dev, tag, value,
	}
	w.writer.Write(record)
	w.writer.Flush()

}
