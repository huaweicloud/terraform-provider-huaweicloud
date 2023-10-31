package logtank

import (
	"github.com/chnsz/golangsdk"
)

type LogtankGet struct {
	ID          string `json:"id"`
	LogGroupID  string `json:"log_group_id"`
	LogStreamID string `json:"log_stream_id"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
}

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type ListResult struct {
	commonResult
}

// Extract will get the logtank object out of the commonResult object.
func (r commonResult) Extract() (LogtankGet, error) {
	var s LogtankGet
	err := r.Result.ExtractInto(&s)
	return s, err
}

// Extract will get the logtank array object out of the ListResult object.
func (lr ListResult) Extract() ([]LogtankGet, error) {
	var l struct {
		Logtanks []LogtankGet `json:"logtanks"`
	}
	err := lr.Result.ExtractInto(&l)
	return l.Logtanks, err
}
