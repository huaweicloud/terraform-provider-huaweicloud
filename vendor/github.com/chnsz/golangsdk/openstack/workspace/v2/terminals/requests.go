package terminals

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure required by the Create method to batch bind the terminals and desktop.
type CreateOpts struct {
	// The list of MAC binding VM policy information that needs to be added.
	BindList []TerminalBindingInfo `json:"bind_list,omitempty"`
}

// TerminalBindingInfo is the structure that represents the configuation of terminal binding.
type TerminalBindingInfo struct {
	// Line number, used for batch import.
	Line int `json:"line,omitempty"`
	// Terminal MAC address.
	Mac string `json:"mac,omitempty"`
	// Desktop name, used for batch import.
	DesktopName string `json:"desktop_name,omitempty"`
	// Description.
	Description string `json:"description,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to batch create terminal binding using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Post(rootURL(c), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// ListOpts is the structure that used to query terminal binding list.
type ListOpts struct {
	// Computer name.
	ComputerName string `q:"computer_name"`
	// MAC address.
	Mac string `q:"mac"`
	// The offset number.
	// Valid range is 0-8000.
	Offset *int `q:"offset" required:"true"`
	// Number of records to be queried.
	// Valid range is 0-8000.
	Limit int `q:"limit" required:"true"`
	// Whether query total count.
	CountOnly bool `q:"count_only"`
}

// List is the method that used to query terminal binding list using given parameters.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]TerminalBindingResp, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pager := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := TerminalBindingPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	})
	pager.Headers = requestOpts.MoreHeaders
	pages, err := pager.AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractTerminalBindings(pages)
}

// UpdateOpts is the structure that used to change the specified binding configuration.
type UpdateOpts struct {
	// Bind ID.
	ID string `json:"id" required:"true"`
	// MAC address.
	MAC *string `json:"mac" required:"true"`
	// Desktop name.
	DesktopName string `json:"desktop_name" required:"true"`
	// Description.
	Description string `json:"description,omitempty"`
}

// Update is a method used to change the specified binding configuration using givin parameters.
func Update(c *golangsdk.ServiceClient, userId string, opts UpdateOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Put(rootURL(c), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// DeleteOpts is the structure required by the Update method to delete some binding configuration.
type DeleteOpts struct {
	// ID list of bind configuration.
	IDs []string `json:"id_list,omitempty"`
}

// Delete is a method to remove an existing user using given parameters.
func Delete(c *golangsdk.ServiceClient, opts DeleteOpts) ([]DeleteResult, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r deleteResp
	_, err = c.Post(deleteURL(c), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return r.ResultList, err
}

// UpdateConfigOpts is the structure that represents the bind switch status.
type UpdateConfigOpts struct {
	// Disable of enable for the bind switch.
	TcBindSwitch string `json:"tc_bind_switch" required:"true"`
}

// UpdateConfig is the method that used to update the bind switch status.
func UpdateConfig(c *golangsdk.ServiceClient, opts UpdateConfigOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}
	_, err = c.Post(configURL(c), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// GetConfig is the method that used to query the current bind switch status.
func GetConfig(c *golangsdk.ServiceClient) (string, error) {
	var r struct {
		TcBindSwitch string `json:"tc_bind_switch"`
	}
	_, err := c.Get(configURL(c), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return r.TcBindSwitch, err
}
