package cluster

import "github.com/huaweicloud/golangsdk"

var requestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

type CreateOpts struct {
	BillingType           int             `json:"billing_type" required:"true"`
	DataCenter            string          `json:"data_center" required:"true"`
	MasterNodeNum         int             `json:"master_node_num" required:"true"`
	MasterNodeSize        string          `json:"master_node_size" required:"true"`
	CoreNodeNum           int             `json:"core_node_num" required:"true"`
	CoreNodeSize          string          `json:"core_node_size" required:"true"`
	AvailableZoneID       string          `json:"available_zone_id" required:"true"`
	ClusterName           string          `json:"cluster_name" required:"true"`
	Vpc                   string          `json:"vpc" required:"true"`
	VpcID                 string          `json:"vpc_id" required:"true"`
	SubnetID              string          `json:"subnet_id" required:"true"`
	SubnetName            string          `json:"subnet_name" required:"true"`
	SecurityGroupsID      string          `json:"security_groups_id,omitempty"`
	ClusterVersion        string          `json:"cluster_version" required:"true"`
	ClusterType           int             `json:"cluster_type"`
	MasterDataVolumeType  string          `json:"master_data_volume_type,omitempty"`
	MasterDataVolumeSize  int             `json:"master_data_volume_size,omitempty"`
	MasterDataVolumeCount int             `json:"master_data_volume_count,omitempty"`
	CoreDataVolumeType    string          `json:"core_data_volume_type,omitempty"`
	CoreDataVolumeSize    int             `json:"core_data_volume_size,omitempty"`
	CoreDataVolumeCount   int             `json:"core_data_volume_count,omitempty"`
	VolumeType            string          `json:"volume_type,omitempty"`
	VolumeSize            int             `json:"volume_size,omitempty"`
	SafeMode              int             `json:"safe_mode"`
	ClusterAdminSecret    string          `json:"cluster_admin_secret" required:"true"`
	LoginMode             int             `json:"login_mode"`
	ClusterMasterSecret   string          `json:"cluster_master_secret,omitempty"`
	NodePublicCertName    string          `json:"node_public_cert_name,omitempty"`
	LogCollection         int             `json:"log_collection,omitempty"`
	ComponentList         []ComponentOpts `json:"component_list" required:"true"`
	AddJobs               []JobOpts       `json:"add_jobs,omitempty"`
	BootstrapScripts      []ScriptOpts    `json:"bootstrap_scripts,omitempty"`
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

type ScriptOpts struct {
	Name                 string   `json:"name" required:"true"`
	Uri                  string   `json:"uri" required:"true"`
	Parameters           string   `json:"parameters,omitempty"`
	Nodes                []string `json:"nodes" required:"true"`
	ActiveMaster         bool     `json:"active_master,omitempty"`
	BeforeComponentStart bool     `json:"before_component_start,omitempty"`
	FailAction           string   `json:"fail_action" required:"true"`
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

	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(createURL(c), b, &r.Body, reqOpt)
	return
}

func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: requestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{
		OkCodes:     []int{204},
		MoreHeaders: requestOpts.MoreHeaders,
	}
	_, r.Err = c.Delete(deleteURL(c, id), reqOpt)
	return
}
