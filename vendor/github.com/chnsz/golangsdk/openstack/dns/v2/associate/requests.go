package associate

import "github.com/chnsz/golangsdk"

type RouterOpts struct {
	RouterID     string `json:"router_id" required:"true"`
	RouterRegion string `json:"router_region,omitempty"`
}

func Associate(client *golangsdk.ServiceClient, resolverRuleID string, opts RouterOpts) (r AssociateResult) {
	b, err := golangsdk.BuildRequestBody(opts, "router")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(associateURL(client, resolverRuleID), b, &r.Body, nil)
	return
}

func DisAssociate(client *golangsdk.ServiceClient, resolverRuleID string, opts RouterOpts) (r DisAssociateResult) {
	b, err := golangsdk.BuildRequestBody(opts, "router")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(disAssociateURL(client, resolverRuleID), b, &r.Body, nil)
	return
}
