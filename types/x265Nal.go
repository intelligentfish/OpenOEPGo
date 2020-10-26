package types

//X265Nal
type X265Nal struct {
	Type    int    `json:"type"`    // nal type
	Size    int    `json:"size"`    // nal size
	Payload []byte `json:"payload"` // nal payload
}
