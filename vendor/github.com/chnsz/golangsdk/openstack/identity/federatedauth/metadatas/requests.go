package metadatas

import "github.com/chnsz/golangsdk"

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

type ImportOpts struct {
	DomainID string `json:"domain_id" required:"true"`
	Metadata string `json:"metadata" required:"true"`
	// It is not a required field, but this argument must be present in the parameter.
	XAccountType string `json:"xaccount_type"`
}

func Import(c *golangsdk.ServiceClient, idpID string, protocolID string, opts ImportOpts) (*MetadataResult, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(resourceURL(c, idpID, protocolID), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r MetadataResult
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

func Get(c *golangsdk.ServiceClient, idpID string, protocolID string) (*Metadata, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, idpID, protocolID), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r Metadata
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}
