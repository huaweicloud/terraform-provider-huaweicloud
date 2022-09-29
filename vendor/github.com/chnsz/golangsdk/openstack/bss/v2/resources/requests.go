package resources

import (
	"github.com/chnsz/golangsdk"
)

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// List of resource IDs.
	ResourceIds []string `json:"resource_ids,omitempty"`
	// Order ID.
	OrderId string `json:"order_id,omitempty"`
	// Whether to query only the main resource, this parameter is invalid when the request parameter is the ID of the
	// sub-resource. If the resource_ids is the ID of the sub-resource, it can only query itself.
	OnlyMainResource int `json:"only_main_resource,omitempty"`
	// resource status.
	StatusList []int `json:"status_list,omitempty"`
	// Query the list of resources that have expired within the specified time period, the start time of the time
	// period, and the UTC time.
	ExpireTimeBegin string `json:"expire_time_begin,omitempty"`
	// Query the list of resources that have expired within the specified time period, the end time of the time period,
	// and the UTC time.
	ExpireTimeEnd string `json:"expire_time_end,omitempty"`
}

func List(c *golangsdk.ServiceClient, opts ListOpts) (*QueryResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r QueryResp
	_, err = c.Post(queryURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// EnableAutoRenew is a method to enable the auto-renew of the prepaid resource.
func EnableAutoRenew(c *golangsdk.ServiceClient, resourceId string) error {
	_, err := c.Post(autoRenewURL(c, resourceId), nil, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
		OkCodes:     []int{204},
	})
	return err
}

// DisableAutoRenew is a method to disable the auto-renew of the prepaid resource.
func DisableAutoRenew(c *golangsdk.ServiceClient, resourceId string) error {
	_, err := c.Delete(autoRenewURL(c, resourceId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
		OkCodes:     []int{204},
	})
	return err
}
