package cluster

import (
	"github.com/chnsz/golangsdk"
)

var requestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

type CreateOpts struct {
	BillingType           int             `json:"billing_type" required:"true"`
	DataCenter            string          `json:"data_center" required:"true"`
	AvailableZoneID       string          `json:"available_zone_id" required:"true"`
	ClusterName           string          `json:"cluster_name" required:"true"`
	Vpc                   string          `json:"vpc" required:"true"`
	VpcID                 string          `json:"vpc_id" required:"true"`
	SubnetID              string          `json:"subnet_id" required:"true"`
	SubnetName            string          `json:"subnet_name" required:"true"`
	SecurityGroupsID      string          `json:"security_groups_id,omitempty"`
	ClusterVersion        string          `json:"cluster_version" required:"true"`
	ClusterType           int             `json:"cluster_type"`
	MasterNodeNum         int             `json:"master_node_num,omitempty"`
	MasterNodeSize        string          `json:"master_node_size,omitempty"`
	CoreNodeNum           int             `json:"core_node_num,omitempty"`
	CoreNodeSize          string          `json:"core_node_size,omitempty"`
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
	NodeGroups            []NodeGroupOpts `json:"node_groups,omitempty"`
	ComponentList         []ComponentOpts `json:"component_list" required:"true"`
	AddJobs               []JobOpts       `json:"add_jobs,omitempty"`
	BootstrapScripts      []ScriptOpts    `json:"bootstrap_scripts,omitempty"`
}

type NodeGroupOpts struct {
	GroupName       string `json:"group_name" required:"true"`
	NodeSize        string `json:"node_size" required:"true"`
	NodeNum         int    `json:"node_num" required:"true"`
	RootVolumeType  string `json:"root_volume_type" required:"true"`
	RootVolumeSize  int    `json:"root_volume_size" required:"true"`
	DataVolumeType  string `json:"data_volume_type" required:"true"`
	DataVolumeSize  int    `json:"data_volume_size" required:"true"`
	DataVolumeCount int    `json:"data_volume_count" required:"true"`
}

type ComponentOpts struct {
	ComponentName string `json:"component_name" required:"true"`
}

type JobOpts struct {
	JobType                 int    `json:"job_type" required:"true"`
	JobName                 string `json:"job_name" required:"true"`
	JarPath                 string `json:"jar_path,omitempty"`
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

type HostOpts struct {
	// Maximum number of clusters displayed on a page
	// Value range: [1-2147483646]. The default value is 10.
	PageSize int `q:"pageSize"`
	// Current page number The default value is 1.
	CurrentPage int `q:"currentPage"`
}

type HostOptsBuilder interface {
	ToHostsListQuery() (string, error)
}

func (opts HostOpts) ToHostsListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
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
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: requestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

// UpdateOpts is a struct which will be used to update a node group of the mrs cluster.
type UpdateOpts struct {
	// Core parameters.
	Parameters ResizeParameters `json:"parameters" required:"true"`
	// Extension parameter.
	PreviousValues PreviousValues `json:"previous_values,omitempty"`
}

// ResizeParameters is an object which is a node group resize operations.
type ResizeParameters struct {
	// Number of nodes to be added or removed.
	Instances string `json:"instances" required:"true"`
	// When expanding or shrinking the capacity, the ID of the node is added or reduced,
	// and the parameter value is fixed as 'node_orderadd'.
	NodeId string `json:"node_id" required:"true"`
	// scale_in: cluster scale-in
	// scale_out: cluster scale-out
	ScaleType string `json:"scale_type" required:"true"`
	// Node group to be scaled out or in.
	// If the value of node_group is core_node_default_group, the node group is a Core node group.
	// If the value of node_group is task_node_default_group, the node group is a Task node group.
	// If it is left blank, the default value core_node_default_group is used.
	NodeGroup *string `json:"node_group,omitempty"`
	// This parameter is valid only when a bootstrap action is configured during cluster creation and takes effect
	// during scale-out. It indicates whether the bootstrap action specified during cluster creation is performed on
	// nodes added during scale-out. The default value is false, indicating that the bootstrap action is performed.
	SkipBootstrapScripts *bool `json:"skip_bootstrap_scripts,omitempty"`
	// Whether to start components on the added nodes after cluster scale-out
	// true: Do not start components after scale-out.
	// false: Start components after scale-out.
	ScaleWithoutStart *bool `json:"scale_without_start,omitempty"`
	// ID list of Task nodes to be deleted during task node scale-in.
	// This parameter does not take effect when scale_type is set to scale-out.
	// If scale_type is set to scale-in and cannot be left blank, the system deletes the specified Task nodes.
	// When scale_type is set to scale-in and server_ids is left blank, the system automatically deletes the Task nodes
	// based on the system rules.
	ServerIds []string `json:"server_ids,omitempty"`
	// Task node specifications.
	// When the number of Task nodes is 0, this parameter is used to specify Task node specifications.
	// When the number of Task nodes is greater than 0, this parameter is unavailable.
	TaskNodeInfo *TaskNodeInfo `json:"task_node_info,omitempty"`
}

// PreviousValues is an object which is a extension parameter.
type PreviousValues struct {
	// Reserve the parameter for extending APIs.
	// You do not need to set the parameter.
	PlanId string `json:"plan_id,omitempty"`
}

// TaskNodeInfo is an object which is a informations of the task node group creation.
type TaskNodeInfo struct {
	// Instance specifications of a Task node, for example, c3.4xlarge.2.linux.bigdata
	NodeSize string `json:"node_size" required:"true"`
	// Data disk storage type of the Task node, supporting SATA, SAS, and SSD currently
	// SATA: Common I/O
	// SAS: High I/O
	// SSD: Ultra-high I/O
	DataVolumeType string `json:"data_volume_type,omitempty"`
	// Number of data disks of a Task node
	// Value range: 1 to 10
	DataVolumeCount int `json:"data_volume_count,omitempty"`
	// Data disk storage space of a Task node
	// Value range: 100 GB to 32,000 GB
	DataVolumeSize int `json:"data_volume_size,omitempty"`
}

// UpdateOptsBuilder is an interface which to support request body build of the node group updation.
type UpdateOptsBuilder interface {
	ToUpdateOptsMap() (map[string]interface{}, error)
}

// ToUpdateOptsMap is a method which to build a request body by the UpdateOpts.
func (opts UpdateOpts) ToUpdateOptsMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update is a method to resize a node group.
func Update(client *golangsdk.ServiceClient, clusterId string, opts UpdateOptsBuilder) (r UpdateResult) {
	reqBody, err := opts.ToUpdateOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, clusterId), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
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

func ListHosts(client *golangsdk.ServiceClient, clusterId string, hostOpts HostOptsBuilder) (*HostListResult, error) {
	url := listHostsURL(client, clusterId)
	listResult := new(HostListResult)

	if hostOpts != nil {
		query, err := hostOpts.ToHostsListQuery()
		if err != nil {
			return nil, err
		}
		url += query
	}

	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}, MoreHeaders: map[string]string{"Content-Type": "application/json;charset=utf8"}}
	_, err := client.Get(url, &listResult, reqOpt)
	return listResult, err

}
