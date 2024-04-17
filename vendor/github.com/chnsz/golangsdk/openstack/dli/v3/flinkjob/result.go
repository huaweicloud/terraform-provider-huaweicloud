package flinkjob

type streamGraphResp struct {
	// Indicates whether the request is successfully executed. Value true indicates that the request is successfully executed.
	IsSuccess bool `json:"is_success"`
	// System prompt. If execution succeeds, the message may be left blank.
	Message string `json:"message"`
	// Error codes.
	ErrorCode string `json:"error_code"`
	// Description of a static stream graph.
	StreamGraph string `json:"stream_graph"`
}
