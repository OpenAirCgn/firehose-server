package firehose_server

import (
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"time"
)

// skip unknown tag value per Default, but allow
// the possibility to dump everything
var DontSkipUnknown bool

type MsgWriter func(Msg)

type CSVMsgWriter struct {
	writer  *csv.Writer
	msgChan <-chan Msg
}

func NewCSVWriter(w io.Writer, ch <-chan Msg) *CSVMsgWriter {
	writer := new(CSVMsgWriter)
	writer.writer = csv.NewWriter(w)
	writer.msgChan = ch
	return writer
}

func uInt32ToA(v uint32) string {
	return strconv.FormatInt(int64(v), 10)
}

func hexToSignedInt32(hexS string) int32 {
	// there has to be a better way, but it's late.
	bs, err := hex.DecodeString(hexS)
	if err != nil {
		println(err)
	}
	var i int32 = 0
	for _, b := range bs {
		i <<= 8
		i = i + int32(b)
	}
	return i
}

func isAlpha(tag Tag) bool {
	switch tag {
	case OA_Alpha_1:
		fallthrough
	case OA_Alpha_2:
		fallthrough
	case OA_Alpha_3:
		fallthrough
	case OA_Alpha_4:
		fallthrough
	case OA_Alpha_5:
		fallthrough
	case OA_Alpha_6:
		fallthrough
	case OA_Alpha_7:
		fallthrough
	case OA_Alpha_8:
		return true
	default:
		return false

	}
}

func getRecord(msg Msg) []string {

	_time := time.Now().Unix()
	time := strconv.FormatInt(_time, 10)
	ts := uInt32ToA(msg.Timestamp)
	dev := msg.DeviceId
	tag := fmt.Sprintf("0x%08x", uint32(msg.Tag))
	valueHex := fmt.Sprintf("0x%08x", msg.Value)
	valueDec := fmt.Sprintf("%d", msg.Value)

	// stupid hackaround: values are transmitted as unsigned values,
	// convert alphasense values to hex then back to signed.

	if isAlpha(msg.Tag) {
		valueDec = fmt.Sprintf("%d", hexToSignedInt32(valueHex[2:]))
	}

	tagAnnotation := msg.Tag.String()
	valueAnnotation := AnnotateValue(msg)

	return []string{
		time, ts, dev, tag, valueHex, valueDec, tagAnnotation, valueAnnotation,
	}
}
func (w *CSVMsgWriter) DumpCSV(msg Msg) {
	if !DontSkipUnknown && msg.Tag.Unknown() {
		return
	}
	record := getRecord(msg)
	w.writer.Write(record)
	w.writer.Flush()
}

func (w *CSVMsgWriter) Run(doneChan chan<- bool) {
	header := []string{
		"server_time",
		"timestamp",
		"device_id",
		"tag",
		"value(hex)",
		"value(decimal)",
		"tag_annotation",
		"value_annotation",
	}
	w.writer.Write(header)

	for {
		msg, ok := <-w.msgChan
		if !ok {
			w.writer.Flush()
			break
		}
		w.DumpCSV(msg)
	}
	doneChan <- true
}
