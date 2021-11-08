package sms

import(
	
	"github.com/chnsz/golangsdk"
)

type CreateOpts struct {
	Name                string `json:"name,omitempty"`
	Is_template         bool 
	Region              string `json:"Region,omitempty"`
	Projectid           string `json:"project_id,omitempty"`
}

type CreateOptsBuilder interface {
	ToSmsTemplateCreateMap() (map[string]interface{}, error)
}


func (opts CreateOpts) ToSmsTemplateCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "sms")
}


func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSmsTemplateCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}