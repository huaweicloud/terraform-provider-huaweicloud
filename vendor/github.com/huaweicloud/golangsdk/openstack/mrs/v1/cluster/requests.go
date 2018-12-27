package cluster

import (
	"log"

	"github.com/huaweicloud/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

type CreateOpts struct {
	BillingType        int             `json:"billing_type" required:"true"`
	DataCenter         string          `json:"data_center" required:"true"`
	MasterNodeNum      int             `json:"master_node_num" required:"true"`
	MasterNodeSize     string          `json:"master_node_size" required:"true"`
	CoreNodeNum        int             `json:"core_node_num" required:"true"`
	CoreNodeSize       string          `json:"core_node_size" required:"true"`
	AvailableZoneID    string          `json:"available_zone_id" required:"true"`
	ClusterName        string          `json:"cluster_name" required:"true"`
	Vpc                string          `json:"vpc" required:"true"`
	VpcID              string          `json:"vpc_id" required:"true"`
	SubnetID           string          `json:"subnet_id" required:"true"`
	SubnetName         string          `json:"subnet_name" required:"true"`
	ClusterVersion     string          `json:"cluster_version,omitempty"`
	ClusterType        int             `json:"cluster_type,omitempty"`
	VolumeType         string          `json:"volume_type" required:"true"`
	VolumeSize         int             `json:"volume_size" required:"true"`
	NodePublicCertName string          `json:"node_public_cert_name" required:"true"`
	SafeMode           int             `json:"safe_mode"`
	ClusterAdminSecret string          `json:"cluster_admin_secret,omitempty"`
	LogCollection      int             `json:"log_collection,omitempty"`
	ComponentList      []ComponentOpts `json:"component_list" required:"true"`
	AddJobs            []JobOpts       `json:"add_jobs,omitempty"`
}

type ComponentOpts struct {
	ComponentName string `json:"component_name" required:"true"`
}

type JobOpts struct {
	JobType                 int    `json:"job_type" required:"true"`
	JobName                 string `json:"job_name" required:"true"`
	JarPath                 string `json:"jar_path" required:"true"`
	Arguments               string `json:"arguments,omitempty"`
	Input                   string `json:"input,omitempty"`
	Output                  string `json:"output,omitempty"`
	JobLog                  string `json:"job_log,omitempty"`
	ShutdownCluster         bool   `json:"shutdown_cluster,omitempty"`
	FileAction              string `json:"file_action,omitempty"`
	SubmitJobOnceClusterRun bool   `json:"submit_job_once_cluster_run" required:"true"`
	Hql                     string `json:"hql,omitempty"`
	HiveScriptPath          string `json:"hive_script_path" required:"true"`
}

type CreateOptsBuilder interface {
	ToClusterCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToClusterCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToClusterCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	log.Printf("[DEBUG] create url:%q, body=%#v", createURL(c), b)
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(createURL(c), b, &r.Body, reqOpt)
	return
}

func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{204},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = c.Delete(deleteURL(c, id), reqOpt)
	return
}
