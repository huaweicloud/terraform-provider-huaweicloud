package loggroups

import "github.com/huaweicloud/golangsdk"

// Log group Create response
type CreateResponse struct {
	ID string `json:"log_group_id"`
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

// UpdateResult contains the response body and error from an Update request.
type UpdateResult struct {
	golangsdk.Result
}

// Extract from UpdateResult
func (r UpdateResult) Extract() (*LogGroup, error) {
	s := new(LogGroup)
	err := r.Result.ExtractInto(s)
	return s, err
}

// DeleteResult is a struct which contains the result of deletion
type DeleteResult struct {
	golangsdk.ErrResult
}

// Log group response
type LogGroup struct {
	ID           string `json:"log_group_id"`
	Name         string `json:"log_group_name"`
	CreationTime int64  `json:"creation_time"`
	TTLinDays    int    `json:"ttl_in_days"`
}

// Log group list response
type LogGroups struct {
	LogGroups []LogGroup `json:"log_groups"`
}

// ListResults contains the body of getting list
type ListResults struct {
	golangsdk.Result
}

// Extract list from GetResult
func (r ListResults) Extract() (*LogGroups, error) {
	s := new(LogGroups)
	err := r.Result.ExtractInto(s)
	return s, err
}
