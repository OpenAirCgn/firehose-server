package firehose_server

type Msg struct {
	Timestamp uint32 `json:"ts"`
	DeviceId  string `json:"device_id"`
	Tag       uint32 `json:"tag"`
	Value     uint32 `json:"value"`
}
