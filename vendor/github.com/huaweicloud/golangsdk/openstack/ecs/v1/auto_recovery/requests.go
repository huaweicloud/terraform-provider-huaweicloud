package auto_recovery

import (
	"log"

	"github.com/huaweicloud/golangsdk"
)

type UpdateOpts struct {
	SupportAutoRecovery string `json:"support_auto_recovery" required:"true"`
}

type UpdateOptsBuilder interface {
	ToAutoRecoveryUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateOpts) ToAutoRecoveryUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) error {
	b, err := opts.ToAutoRecoveryUpdateMap()
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] update url:%q, body=%#v", updateURL(c, id), b)
	_, err = c.Put(updateURL(c, id), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return err
}

func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}
