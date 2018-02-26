package topics

import (
	"github.com/huaweicloud/golangsdk"
)

type Topic struct {
	RequestId string `json:"request_id"`
	TopicUrn  string `json:"topic_urn"`
}

type TopicGet struct {
	TopicUrn    string `json:"topic_urn"`
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
	PushPolicy  int    `json:"push_policy"`
	UpdateTime  string `json:"update_time"`
	CreateTime  string `json:"create_time"`
}

// Extract will get the topic object out of the commonResult object.
func (r commonResult) Extract() (*Topic, error) {
	var s Topic
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractGet() (*TopicGet, error) {
	var s TopicGet
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

type commonResult struct {
	golangsdk.Result
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type ListResult struct {
	golangsdk.Result
}

func (lr ListResult) Extract() ([]TopicGet, error) {
	var a struct {
		Topics []TopicGet `json:"topics"`
	}
	err := lr.Result.ExtractInto(&a)
	return a.Topics, err
}
