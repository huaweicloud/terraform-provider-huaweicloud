package tracker

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type Tracker struct {
	Status                    string                    `json:"status"`
	BucketName                string                    `json:"bucket_name"`
	FilePrefixName            string                    `json:"file_prefix_name"`
	TrackerName               string                    `json:"tracker_name"`
	SimpleMessageNotification SimpleMessageNotification `json:"smn"`
}

// Extract will get the tracker object from the commonResult
func (r commonResult) Extract() (*Tracker, error) {
	var s Tracker
	err := r.ExtractInto(&s)
	return &s, err
}

type TrackerPage struct {
	pagination.LinkedPageBase
}

// ExtractTracker accepts a Page struct, specifically a TrackerPage struct,
// and extracts the elements into a slice of Tracker structs. In other words,
// a generic collection is mapped into a relevant slice.
func (r commonResult) ExtractTracker() ([]Tracker, error) {
	var s []Tracker
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}

	return s, nil

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
