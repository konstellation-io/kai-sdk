package prediction

type Prediction struct {
	Timestamp int64                  `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
	Metadata  Metadata               `json:"metadata"`
}

type Metadata struct {
	Product   string `json:"product"`
	Version   string `json:"version"`
	Workflow  string `json:"workflow"`
	Process   string `json:"process"`
	RequestID string `json:"requestID"`
}
