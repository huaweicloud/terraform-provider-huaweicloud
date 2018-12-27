package job

import (
	"log"

	"github.com/huaweicloud/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

type CreateOpts struct {
	JobType        int    `json:"job_type" required:"true"`
	JobName        string `json:"job_name" required:"true"`
	ClusterID      string `json:"cluster_id" required:"true"`
	JarPath        string `json:"jar_path" required:"true"`
	Arguments      string `json:"arguments,omitempty"`
	Input          string `json:"input,omitempty"`
	Output         string `json:"output,omitempty"`
	JobLog         string `json:"job_log,omitempty"`
	HiveScriptPath string `json:"hive_script_path,omitempty"`
	IsProtected    bool   `json:"is_protected,omitempty"`
	IsPublic       bool   `json:"is_public,omitempty"`
}

type CreateOptsBuilder interface {
	ToJobCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToJobCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToJobCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	log.Printf("[DEBUG] create url:%q, body=%#v", createURL(c), b)
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = c.Post(createURL(c), b, &r.Body, reqOpt)
	return
}

func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = c.Get(getURL(c, id), &r.Body, reqOpt)
	return
}

func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{204},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = c.Delete(deleteURL(c, id), reqOpt)
	return
}
