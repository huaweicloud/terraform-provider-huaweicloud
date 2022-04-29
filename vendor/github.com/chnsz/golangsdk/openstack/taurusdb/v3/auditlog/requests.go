package auditlog

import (
	"github.com/chnsz/golangsdk"
)

type UpdateAuditlogOpts struct {
	SwitchStatus string `json:"switch_status" required:"true"`
}

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

func Update(c *golangsdk.ServiceClient, instanceId string, opts UpdateAuditlogOpts) (*UpdateResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst UpdateResp
	_, err = c.Post(updateURL(c, instanceId), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Get(c *golangsdk.ServiceClient, instanceId string) (*AuditLogStatus, error) {
	var rst AuditLogStatus
	_, err := c.Get(getURL(c, instanceId), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}
