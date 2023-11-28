package prediction

type Prediction struct {
	CreationDate int64    `json:"creationDate"`
	LastModified int64    `json:"lastModified"`
	Payload      Payload  `json:"payload"`
	Metadata     Metadata `json:"metadata"`
}

type Payload map[string]interface{}

type Metadata struct {
	Product      string `json:"product"`
	Version      string `json:"version"`
	Workflow     string `json:"workflow"`
	WorkflowType string `json:"workflowType"`
	Process      string `json:"process"`
	RequestID    string `json:"requestID"`
}
