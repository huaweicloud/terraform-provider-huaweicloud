package logstreams

import "github.com/huaweicloud/golangsdk"

// Log stream Create response
type CreateResponse struct {
	ID string `json:"log_stream_id"`
}

// CreateResult is a struct that contains all the return parameters of creation
type CreateResult struct {
	golangsdk.Result
}

// Extract from CreateResult
func (r CreateResult) Extract() (*CreateResponse, error) {
	s := new(CreateResponse)
	err := r.Result.ExtractInto(s)
	return s, err
}

// DeleteResult is a struct which contains the result of deletion
type DeleteResult struct {
	golangsdk.ErrResult
}

// Log stream response
type LogStream struct {
	ID           string `json:"log_stream_id"`
	Name         string `json:"log_stream_name"`
	CreationTime int64  `json:"creation_time"`
	FilterCount  int64  `json:"filter_count"`
}

// Log stream list response
type LogStreams struct {
	LogStreams []LogStream `json:"log_streams"`
}

// ListResults contains the body of getting list
type ListResults struct {
	golangsdk.Result
}

// Extract list from GetResult
func (r ListResults) Extract() (*LogStreams, error) {
	s := new(LogStreams)
	err := r.Result.ExtractInto(s)
	return s, err
}
