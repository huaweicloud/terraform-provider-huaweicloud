package subscriptions

import (
	"github.com/huaweicloud/golangsdk"
)

type Subscription struct {
	RequestId       string `json:"request_id"`
	SubscriptionUrn string `json:"subscription_urn"`
}

type SubscriptionGet struct {
	TopicUrn        string `json:"topic_urn"`
	Protocol        string `json:"protocol"`
	SubscriptionUrn string `json:"subscription_urn"`
	Owner           string `json:"owner"`
	Endpoint        string `json:"endpoint"`
	Remark          string `json:"remark"`
	Status          int    `json:"status"`
}

// Extract will get the subscription object out of the commonResult object.
func (r commonResult) Extract() (*Subscription, error) {
	var s Subscription
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

type ListResult struct {
	golangsdk.Result
}

func (lr ListResult) Extract() ([]SubscriptionGet, error) {
	var a struct {
		Subscriptions []SubscriptionGet `json:"subscriptions"`
	}
	err := lr.Result.ExtractInto(&a)
	return a.Subscriptions, err
}
