package firehose_server

import (
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"time"
)

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

	tag_annotation := msg.Tag.String()
	value_annotation := AnnotateValue(msg)

	return []string{
		time, ts, dev, tag, valueHex, valueDec, tag_annotation, value_annotation,
	}
}
func (w *CSVMsgWriter) DumpCSV(msg Msg) {
	record := getRecord(msg)
	w.writer.Write(record)
	w.handleAlphaAnnotation(msg)
}

const adc2v = 0.000031356811523

var last_alpha1 int64
var last_alpha3 int64
var last_alpha5 int64
var last_alpha7 int64

func (w *CSVMsgWriter) handleAlphaAnnotation(msg Msg) {
	var alpha float64
	var tag Tag

	switch msg.Tag {
	case OA_Alpha_1:
		last_alpha1 = int64(msg.Value)
		return
	case OA_Alpha_2:
		// ignore corner case that last_alpha1 is not set
		alpha = float64(last_alpha1-int64(msg.Value)) * adc2v
		tag = OA_AlphaCalc_1
	case OA_Alpha_3:
		last_alpha3 = int64(msg.Value)
		return
	case OA_Alpha_4:
		alpha = float64(last_alpha3-int64(msg.Value)) * adc2v
		tag = OA_AlphaCalc_2
	case OA_Alpha_5:
		last_alpha5 = int64(msg.Value)
		return
	case OA_Alpha_6:
		alpha = float64(last_alpha5-int64(msg.Value)) * adc2v
		tag = OA_AlphaCalc_3
	case OA_Alpha_7:
		last_alpha7 = int64(msg.Value)
		return
	case OA_Alpha_8:
		alpha = float64(last_alpha7-int64(msg.Value)) * adc2v
		tag = OA_AlphaCalc_4
	default:
		return
	}

	msg.Tag = tag
	msg.Value = 0

	record := getRecord(msg)
	record[len(record)-1] = fmt.Sprintf("%f V (derived)", alpha)
	w.writer.Write(record)
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
