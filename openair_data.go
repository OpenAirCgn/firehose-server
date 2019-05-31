package firehose_server

import (
	"fmt"
	"math"
)

type Tag uint32
type NetworkEvent uint32

const (
	OA_Time             Tag = iota
	OA_Alpha_1          Tag = iota
	OA_Alpha_2          Tag = iota
	OA_Alpha_3          Tag = iota
	OA_Alpha_4          Tag = iota
	OA_Alpha_5          Tag = iota
	OA_Alpha_6          Tag = iota
	OA_Alpha_7          Tag = iota
	OA_Alpha_8          Tag = iota
	OA_BME_Pressure_Raw Tag = iota
	OA_BME_Pressure     Tag = iota
	OA_BME_Temp_Raw     Tag = iota
	OA_BME_Temp         Tag = iota
	OA_BME_Humidity_Raw Tag = iota
	OA_BME_Humidity     Tag = iota
	OA_BME_SDS_PM25     Tag = iota
	OA_BME_SDS_PM10     Tag = iota
	OA_Network_Events   Tag = math.MaxUint32
)

const (
	CONNECT    NetworkEvent = iota
	DISCONNECT NetworkEvent = iota
)

func (t Tag) String() string {
	switch t {
	case OA_Time:
		return "OA_Time"
	case OA_Alpha_1:
		return "OA_Alpha_1"
	case OA_Alpha_2:
		return "OA_Alpha_2"
	case OA_Alpha_3:
		return "OA_Alpha_3"
	case OA_Alpha_4:
		return "OA_Alpha_4"
	case OA_Alpha_5:
		return "OA_Alpha_5"
	case OA_Alpha_6:
		return "OA_Alpha_6"
	case OA_Alpha_7:
		return "OA_Alpha_7"
	case OA_Alpha_8:
		return "OA_Alpha_8"
	case OA_BME_Pressure_Raw:
		return "OA_BME_Pressure_Raw"
	case OA_BME_Pressure:
		return "OA_BME_Pressure"
	case OA_BME_Temp_Raw:
		return "OA_BME_Temp_Raw"
	case OA_BME_Temp:
		return "OA_BME_Temp"
	case OA_BME_Humidity_Raw:
		return "OA_BME_Humidity_Raw"
	case OA_BME_Humidity:
		return "OA_BME_Humidity"
	case OA_BME_SDS_PM25:
		return "OA_BME_SDS_PM25"
	case OA_BME_SDS_PM10:
		return "OA_BME_SDS_PM10"
	case OA_Network_Events:
		return "OA_Network_Events"
	default:
		return "UNKNOWN TAG"
	}
}

type Msg struct {
	Timestamp uint32 `json:"ts"`
	DeviceId  string `json:"device_id"`
	Tag       Tag    `json:"tag"`
	Value     uint32 `json:"value"`
}

func AnnotateValue(m Msg) string {
	switch m.Tag {
	case OA_Time:
		return fmt.Sprintf("%d s", m.Value)
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
		fallthrough
	case OA_BME_Pressure_Raw:
		fallthrough
	case OA_BME_Temp_Raw:
		fallthrough
	case OA_BME_Humidity_Raw:
		return "raw"
	case OA_BME_Pressure:
		hPa := float64(m.Value) / 100.0
		return fmt.Sprintf("%.2 hPa", hPa)
	case OA_BME_Temp:
		temp := float64(m.Value)/1000.0 - 273.15
		return fmt.Sprintf("%.2 C", temp)
	case OA_BME_Humidity:
		hum := float64(m.Value) / 10.0
		return fmt.Sprintf("%.2 %%RH", hum)
	case OA_BME_SDS_PM25:
		pm25 := float64(m.Value) / 1000.0
		return fmt.Sprintf("%.2 ug/m3", pm25)
	case OA_BME_SDS_PM10:
		pm10 := float64(m.Value) / 1000.0
		return fmt.Sprintf("%.2 ug/m3", pm10)
	case OA_Network_Events:
		switch NetworkEvent(m.Value) {
		case CONNECT:
			return "CONNECT"
		case DISCONNECT:
			return "DISCONNECT"
		default:
			return "UNKNOWN NETWORK EVENT"
		}
	default:
		return "UNKNOWN TAG"

	}
}
