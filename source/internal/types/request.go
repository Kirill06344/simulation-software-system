package types

type Request struct {
	Id          string  `json:"id"`
	CurrentTime float32 `json:"currentTime"`
	SourceId    int32   `json:"sourceId"`
}
