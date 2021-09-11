package tags

import "github.com/chnsz/golangsdk"

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

//ResourceTag is in key-value format
type ResourceTag struct {
	Key string `json:"key" required:"ture"`
	// Value is not required, but it must be included in the parameter.
	Value string `json:"value"`
}

type tagsActionOpts struct {
	Action string        `json:"action" required:"true"`
	Tags   []ResourceTag `json:"tags,omitempty"`
}

func Create(c *golangsdk.ServiceClient, instanceID string, tags []ResourceTag) error {
	opts := tagsActionOpts{
		Action: "create",
		Tags:   tags,
	}
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Post(actionURL(c, instanceID), b, nil, &golangsdk.RequestOpts{
		OkCodes:     []int{204},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return err
}

func Delete(c *golangsdk.ServiceClient, instanceID string, tags []ResourceTag) error {
	opts := tagsActionOpts{
		Action: "delete",
		Tags:   tags,
	}
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Post(actionURL(c, instanceID), b, nil, &golangsdk.RequestOpts{
		OkCodes:     []int{204},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return err
}

func Get(c *golangsdk.ServiceClient, instanceID string) ([]ResourceTag, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, instanceID), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r []ResourceTag
		rst.ExtractInto(&r)
		return r, nil
	}
	return nil, err
}

func List(c *golangsdk.ServiceClient) (*ResourceTagList, error) {
	var rst golangsdk.Result
	_, err := c.Get(listURL(c), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r ResourceTagList
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}
