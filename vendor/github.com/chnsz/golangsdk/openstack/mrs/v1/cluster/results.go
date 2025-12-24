package cluster

import "github.com/chnsz/golangsdk"

type Cluster struct {
	Clusterid             string            `json:"clusterId"`
	Clustername           string            `json:"clusterName"`
	Masternodenum         string            `json:"masterNodeNum"`
	Corenodenum           string            `json:"coreNodeNum"`
	Totalnodenum          string            `json:"totalNodeNum"`
	Clusterstate          string            `json:"clusterState"`
	Createat              string            `json:"createAt"`
	Updateat              string            `json:"updateAt"`
	Billingtype           string            `json:"billingType"`
	Datacenter            string            `json:"dataCenter"`
	Duration              string            `json:"duration"`
	Fee                   string            `json:"fee"`
	Hadoopversion         string            `json:"hadoopVersion"`
	Masternodesize        string            `json:"masterNodeSize"`
	Corenodesize          string            `json:"coreNodeSize"`
	Componentlist         []Component       `json:"componentList"`
	Externalip            string            `json:"externalIp"`
	Externalalternateip   string            `json:"externalAlternateIp"`
	Internalip            string            `json:"internalIp"`
	Deploymentid          string            `json:"deploymentId"`
	Remark                string            `json:"remark"`
	Orderid               string            `json:"orderId"`
	AvailabilityZone      string            `json:"azCode"`
	Azid                  string            `json:"azId"`
	Azname                string            `json:"azName"`
	Masternodeproductid   string            `json:"masterNodeProductId"`
	Masternodespecid      string            `json:"masterNodeSpecId"`
	Corenodeproductid     string            `json:"coreNodeProductId"`
	Corenodespecid        string            `json:"coreNodeSpecId"`
	Instanceid            string            `json:"instanceId"`
	Vnc                   string            `json:"vnc"`
	Tenantid              string            `json:"tenantId"`
	Volumesize            int               `json:"volumeSize"`
	Vpc                   string            `json:"vpc"`
	Vpcid                 string            `json:"vpcId"`
	Subnetid              string            `json:"subnetId"`
	Subnetname            string            `json:"subnetName"`
	Securitygroupsid      string            `json:"securityGroupsId"`
	Slavesecuritygroupsid string            `json:"slaveSecurityGroupsId"`
	Stagedesc             string            `json:"stageDesc"`
	Safemode              int               `json:"safeMode"`
	Clusterversion        string            `json:"clusterVersion"`
	ClusterType           int               `json:"clusterType"`
	Nodepubliccertname    string            `json:"nodePublicCertName"`
	Masternodeip          string            `json:"masterNodeIp"`
	Privateipfirst        string            `json:"privateIpFirst"`
	Errorinfo             string            `json:"errorInfo"`
	Chargingstarttime     string            `json:"chargingStartTime"`
	LogCollection         int               `json:"logCollection"`
	TaskNodeGroups        []NodeGroup       `json:"taskNodeGroups"`
	NodeGroups            []NodeGroup       `json:"nodeGroups"`
	MasterDataVolumeType  string            `json:"masterDataVolumeType"`
	MasterDataVolumeSize  int               `json:"masterDataVolumeSize"`
	MasterDataVolumeCount int               `json:"masterDataVolumeCount"`
	CoreDataVolumeType    string            `json:"coreDataVolumeType"`
	CoreDataVolumeSize    int               `json:"coreDataVolumeSize"`
	CoreDataVolumeCount   int               `json:"coreDataVolumeCount"`
	BootstrapScripts      []BootStrapScript `json:"bootstrapScripts"`
	EnterpriseProjectId   string            `json:"enterpriseProjectId"`
	IsMrsManagerFinish    bool              `json:"ismrsManagerFinish"`
	PeriodType            int               `json:"periodType"`
	Scale                 string            `json:"scale"`
	EipId                 string            `json:"eipId"`
	EipAddress            string            `json:"eipAddress"`
	Eipv6Address          string            `json:"eipv6Address"`
	Tags                  string            `json:"tags"`
	// The default agency name bound to the cluster node.
	MrsEcsDefaultAgency string `json:"mrsEcsDefaultAgency"`
}

type Component struct {
	Componentid      string `json:"componentId"`
	Componentname    string `json:"componentName"`
	Componentversion string `json:"componentVersion"`
	Componentdesc    string `json:"componentDesc"`
}

type NodeGroup struct {
	GroupName                  string   `json:"groupName"`
	NodeNum                    int      `json:"nodeNum"`
	NodeSize                   string   `json:"nodeSize"`
	NodeSpecId                 string   `json:"nodeSpecId"`
	NodeProductId              string   `json:"nodeProductId"`
	VMProductId                string   `json:"vmProductId"`
	VMSpecCode                 string   `json:"vmSpecCode"`
	RootVolumeSize             int      `json:"rootVolumeSize"`
	RootVolumeType             string   `json:"rootVolumeType"`
	RootVolumeProductId        string   `json:"rootVolumeProductId"`
	RootVolumeResourceSpecCode string   `json:"rootVolumeResourceSpecCode"`
	DataVolumeType             string   `json:"dataVolumeType"`
	DataVolumeSize             int      `json:"dataVolumeSize"`
	DataVolumeCount            int      `json:"dataVolumeCount"`
	DataVolumeResourceSpecCode string   `json:"dataVolumeResourceSpecCode"`
	DataVolumeResourceType     string   `json:"dataVolumeResourceType"`
	AssignedRoles              []string `json:"assignedRoles"`
}

type BootStrapScript struct {
	Name                 string   `json:"name"`
	URI                  string   `json:"uri"`
	Parameters           string   `json:"parameters"`
	Nodes                []string `json:"nodes"`
	ActiveMaster         bool     `json:"active_master"`
	BeforeComponentStart bool     `json:"before_component_start"`
	ExecuteNeedSudoRoot  bool     `json:"execute_need_sudo_root"`
	FailAction           string   `json:"fail_action"`
	StartTime            int      `json:"start_time"`
	State                string   `json:"state"`
}

type ClusterResult struct {
	ClusterID string `json:"cluster_id"`
	Result    bool   `json:"result"`
	Msg       string `json:"msg"`
}

type HostListResult struct {
	Hosts []Host `json:"hosts"`
	Total int    `json:"total"`
}

type Host struct {
	// VM ID
	Id string `json:"id"`
	// VM IP address
	Ip string `json:"ip"`
	// VM flavor ID
	Flavor string `json:"flavor"`
	// VM type
	// Currently, MasterNode, CoreNode, and TaskNode are supported.
	Type string `json:"type"`
	// VM name
	Name string `json:"name"`
	// Current VM state
	Status string `json:"status"`
	// Memory
	Mem string `json:"mem"`
	// Number of CPU cores
	Cpu string `json:"cpu"`
	// OS disk capacity
	RootVolumeSize string `json:"root_volume_size"`
	// Data disk type
	DataVolumeType string `json:"data_volume_type"`
	// Data disk capacity
	DataVolumeSize int `json:"data_volume_size"`
	// Number of data disks
	DataVolumeCount int `json:"data_volume_count"`
}

type CreateResult struct {
	golangsdk.Result
}

func (r CreateResult) Extract() (*ClusterResult, error) {
	var s ClusterResult
	err := r.ExtractInto(&s)
	return &s, err
}

type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (*Cluster, error) {
	var s Cluster
	err := r.ExtractInto(&s)
	return &s, err
}

func (r GetResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "cluster")
}

// UpdateResult represents a result of the Update method.
type UpdateResult struct {
	golangsdk.Result
}

// UpdateResp is an object struct that represents an result of node group resize operation.
type UpdateResp struct {
	// Operation result
	// succeeded: The operation is successful.
	// Table 8 describes the error codes returned upon operation failures.
	Result string `json:"result"`
}

// Extract is a method which to extract the response of the resize operation.
func (r UpdateResult) Extract() (*UpdateResp, error) {
	var s UpdateResp
	err := r.ExtractInto(&s)
	return &s, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}
