package task

import (
	"github.com/huaweicloud/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

type CreateOpts struct {
	SrcNode     SrcNodeOpts  `json:"src_node" required:"true"`
	DstNode     DstNodeOpts  `json:"dst_node" required:"true"`
	EnableKMS   *bool        `json:"enableKMS" required:"true"`
	ThreadNum   int          `json:"thread_num" required:"true"`
	Description string       `json:"description,omitempty"`
	SmnInfo     *SmnInfoOpts `json:"smnInfo,omitempty"`
}

type SrcNodeOpts struct {
	Region    string `json:"region" required:"true"`
	AK        string `json:"ak" required:"true"`
	SK        string `json:"sk" required:"true"`
	ObjectKey string `json:"object_key" required:"true"`
	Bucket    string `json:"bucket" required:"true"`
	CloudType string `json:"cloud_type,omitempty"`
}

type DstNodeOpts struct {
	Region    string `json:"region" required:"true"`
	AK        string `json:"ak" required:"true"`
	SK        string `json:"sk" required:"true"`
	ObjectKey string `json:"object_key,omitempty"`
	Bucket    string `json:"bucket" required:"true"`
}

type SmnInfoOpts struct {
	TopicUrn          string   `json:"topicUrn" required:"true"`
	Language          string   `json:"language,omitempty"`
	TriggerConditions []string `json:"triggerConditions" required:"true"`
}

type CreateOptsBuilder interface {
	ToTaskCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToTaskCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToTaskCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, reqOpt)
	return
}

func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = c.Delete(resourceURL(c, id), reqOpt)
	return
}
