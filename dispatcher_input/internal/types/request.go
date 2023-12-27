package types

type Request struct {
	Id             uint32  `json:"id"`
	GenerationTime float32 `json:"generationTime"`
	CurrentTime    float32 `json:"currentTime"`
	SourceId       int32   `json:"sourceId"`
}
