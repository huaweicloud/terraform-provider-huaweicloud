package subscriptions

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type Subscription struct {
	RequestId       string `json:"request_id"`
	SubscriptionUrn string `json:"subscription_urn"`
}

type SubscriptionGet struct {
	TopicUrn        string         `json:"topic_urn"`
	Protocol        string         `json:"protocol"`
	SubscriptionUrn string         `json:"subscription_urn"`
	Owner           string         `json:"owner"`
	Endpoint        string         `json:"endpoint"`
	Remark          string         `json:"remark"`
	Status          int            `json:"status"`
	FilterPolicies  []FilterPolicy `json:"filter_polices"`
}

type FilterPolicy struct {
	Name         string   `json:"name"`
	StringEquals []string `json:"string_equals"`
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

type SubscriptionPage struct {
	pagination.OffsetPageBase
}

func (b SubscriptionPage) IsEmpty() (bool, error) {
	arr, err := ExtractSubscriptions(b)
	return len(arr) == 0, err
}

func ExtractSubscriptions(r pagination.Page) ([]SubscriptionGet, error) {
	var s struct {
		Subscriptions []SubscriptionGet `json:"subscriptions"`
	}
	err := (r.(SubscriptionPage)).ExtractInto(&s)
	return s.Subscriptions, err
}
