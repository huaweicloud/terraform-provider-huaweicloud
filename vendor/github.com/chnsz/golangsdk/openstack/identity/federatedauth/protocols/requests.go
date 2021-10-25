package protocols

import "github.com/chnsz/golangsdk"

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

type ProtocolOption struct {
	MappingID string `json:"mapping_id,omitempty"`
}

func Create(c *golangsdk.ServiceClient, idpID, protocolID string, mappingID string) (*IdentityProtocol, error) {
	opts := ProtocolOption{MappingID: mappingID}
	b, err := golangsdk.BuildRequestBody(opts, "protocol")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Put(resourceURL(c, idpID, protocolID), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r struct {
			Protocol IdentityProtocol `json:"protocol"`
		}
		rst.ExtractInto(&r)
		return &r.Protocol, nil
	}
	return nil, err
}

func Get(c *golangsdk.ServiceClient, idpID string, protocolID string) (*IdentityProtocol, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, idpID, protocolID), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	if err == nil {
		var r struct {
			Protocol IdentityProtocol `json:"protocol"`
		}
		rst.ExtractInto(&r)
		return &r.Protocol, nil
	}
	return nil, err
}

func List(c *golangsdk.ServiceClient, idpID string) ([]IdentityProtocol, error) {
	var rst golangsdk.Result
	_, err := c.Get(root(c, idpID), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	if err == nil {
		var r struct {
			Protocols []IdentityProtocol `json:"protocols"`
		}
		rst.ExtractInto(&r)
		return r.Protocols, nil
	}
	return nil, err
}

func Update(c *golangsdk.ServiceClient, idpID string, id string, mappingID string) (*IdentityProtocol, error) {
	opts := ProtocolOption{MappingID: mappingID}
	b, err := golangsdk.BuildRequestBody(opts, "protocol")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Patch(resourceURL(c, idpID, id), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r struct {
			Protocol IdentityProtocol `json:"protocol"`
		}
		rst.ExtractInto(&r)
		return &r.Protocol, nil
	}
	return nil, err
}

func Delete(c *golangsdk.ServiceClient, idpID string, protocolID string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult
	_, r.Err = c.Delete(resourceURL(c, idpID, protocolID), &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &r
}
