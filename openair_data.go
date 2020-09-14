package firehose_server

import (
	"fmt"
	"math"
)

type Tag uint32
type NetworkEvent uint32

const (
	OA_Time                Tag = iota
	OA_Alpha_1             Tag = iota
	OA_Alpha_2             Tag = iota
	OA_Alpha_3             Tag = iota
	OA_Alpha_4             Tag = iota
	OA_Alpha_5             Tag = iota
	OA_Alpha_6             Tag = iota
	OA_Alpha_7             Tag = iota
	OA_Alpha_8             Tag = iota
	OA_BME0_Pressure_Raw   Tag = iota
	OA_BME0_Pressure       Tag = iota
	OA_BME0_Temp_Raw       Tag = iota
	OA_BME0_Temp           Tag = iota
	OA_BME0_Humidity_Raw   Tag = iota
	OA_BME0_Humidity       Tag = iota
	OA_BME1_Pressure_Raw   Tag = iota
	OA_BME1_Pressure       Tag = iota
	OA_BME1_Temp_Raw       Tag = iota
	OA_BME1_Temp           Tag = iota
	OA_BME1_Humidity_Raw   Tag = iota
	OA_BME1_Humidity       Tag = iota
	OA_SDS_PM25            Tag = iota
	OA_SDS_PM10            Tag = iota
	OA_SI7006_Temp_Raw     Tag = iota
	OA_SI7006_Temp         Tag = iota
	OA_SI7006_RH_Raw       Tag = iota
	OA_SI7006_RH           Tag = iota
	OA_MICS4514_VRED       Tag = iota
	OA_MICS4514_VOX        Tag = iota
	OA_NOISE_DBA           Tag = iota
	OA_NOISE_DBC           Tag = iota
	OA_FINAL_SPECIAL_GUARD Tag = iota
	OA_Network_Events      Tag = math.MaxUint32
	OA_AlphaCalc_1         Tag = OA_Network_Events - 1
	OA_AlphaCalc_2         Tag = OA_Network_Events - 2
	OA_AlphaCalc_3         Tag = OA_Network_Events - 3
	OA_AlphaCalc_4         Tag = OA_Network_Events - 4
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
	case OA_BME0_Pressure_Raw:
		return "OA_BME0_Pressure_Raw"
	case OA_BME0_Pressure:
		return "OA_BME0_Pressure"
	case OA_BME0_Temp_Raw:
		return "OA_BME0_Temp_Raw"
	case OA_BME0_Temp:
		return "OA_BME0_Temp"
	case OA_BME0_Humidity_Raw:
		return "OA_BME0_Humidity_Raw"
	case OA_BME0_Humidity:
		return "OA_BME0_Humidity"
	case OA_BME1_Pressure_Raw:
		return "OA_BME1_Pressure_Raw"
	case OA_BME1_Pressure:
		return "OA_BME1_Pressure"
	case OA_BME1_Temp_Raw:
		return "OA_BME1_Temp_Raw"
	case OA_BME1_Temp:
		return "OA_BME1_Temp"
	case OA_BME1_Humidity_Raw:
		return "OA_BME1_Humidity_Raw"
	case OA_BME1_Humidity:
		return "OA_BME1_Humidity"
	case OA_SDS_PM25:
		return "OA_BME1_SDS_PM25"
	case OA_SDS_PM10:
		return "OA_BME1_SDS_PM10"
	case OA_SI7006_Temp_Raw:
		return "OA_SI7006_Temp_Raw"
	case OA_SI7006_Temp:
		return "OA_SI7006_Temp"
	case OA_SI7006_RH_Raw:
		return "OA_SI7006_RH_Raw"
	case OA_SI7006_RH:
		return "OA_SI7006_RH"
	case OA_MICS4514_VRED:
		return "OA_MICS4514_VRED"
	case OA_MICS4514_VOX:
		return "OA_MICS4514_VOX"
	case OA_Network_Events:
		return "OA_Network_Events"
	case OA_AlphaCalc_1:
		return "OA_AlphaCalc_1"
	case OA_AlphaCalc_2:
		return "OA_AlphaCalc_2"
	case OA_AlphaCalc_3:
		return "OA_AlphaCalc_3"
	case OA_AlphaCalc_4:
		return "OA_AlphaCalc_4"
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
	case OA_BME0_Pressure_Raw:
		fallthrough
	case OA_BME0_Temp_Raw:
		fallthrough
	case OA_BME0_Humidity_Raw:
		fallthrough
	case OA_BME1_Pressure_Raw:
		fallthrough
	case OA_BME1_Temp_Raw:
		fallthrough
	case OA_BME1_Humidity_Raw:
		fallthrough
	case OA_SI7006_Temp_Raw:
		fallthrough
	case OA_SI7006_RH_Raw:
		return fmt.Sprintf("(raw %d)", m.Value)
	case OA_BME0_Pressure:
		fallthrough
	case OA_BME1_Pressure:
		hPa := float64(m.Value) / 100.0
		return fmt.Sprintf("%.2f hPa", hPa)
	case OA_BME0_Temp:
		fallthrough
	case OA_BME1_Temp:
		temp := float64(m.Value)/1000.0 - 273.15
		return fmt.Sprintf("%.2f C", temp)
	case OA_SI7006_Temp:
		temp := float64(m.Value) / 1000.0
		return fmt.Sprintf("%.2f C", temp)
	case OA_BME0_Humidity:
		fallthrough
	case OA_BME1_Humidity:
		fallthrough
	case OA_SI7006_RH:
		hum := float64(m.Value) / 100.0
		return fmt.Sprintf("%.2f %%RH", hum)
	case OA_SDS_PM25:
		pm25 := float64(m.Value) / 1000.0
		return fmt.Sprintf("%.2f ug/m3", pm25)
	case OA_SDS_PM10:
		pm10 := float64(m.Value) / 1000.0
		return fmt.Sprintf("%.2f ug/m3", pm10)
	case OA_MICS4514_VRED:
		fallthrough
	case OA_MICS4514_VOX:
		return fmt.Sprintf("%d mV", m.Value)
	case OA_Network_Events:
		switch NetworkEvent(m.Value) {
		case CONNECT:
			return "CONNECT"
		case DISCONNECT:
			return "DISCONNECT"
		default:
			return "UNKNOWN NETWORK EVENT"
		}
	case OA_AlphaCalc_1:
		fallthrough
	case OA_AlphaCalc_2:
		fallthrough
	case OA_AlphaCalc_3:
		fallthrough
	case OA_AlphaCalc_4:
		return "derived (V)"
	default:
		return "UNKNOWN TAG"

	}
}
