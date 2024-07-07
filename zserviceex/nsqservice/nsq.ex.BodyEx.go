package nsqservice

type BodyEx struct {
	S2S  string `json:"s2s"`  // S2S数据
	Body []byte `json:"body"` // body数据
}
