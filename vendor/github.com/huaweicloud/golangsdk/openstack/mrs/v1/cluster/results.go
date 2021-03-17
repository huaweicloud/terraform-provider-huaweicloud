package cluster

import "github.com/huaweicloud/golangsdk"

type Cluster struct {
	Clusterid             string       `json:"clusterId"`
	Clustername           string       `json:"clusterName"`
	Masternodenum         string       `json:"masterNodeNum"`
	Corenodenum           string       `json:"coreNodeNum"`
	Totalnodenum          string       `json:"totalNodeNum"`
	Clusterstate          string       `json:"clusterState"`
	Createat              string       `json:"createAt"`
	Updateat              string       `json:"updateAt"`
	Billingtype           string       `json:"billingType"`
	Datacenter            string       `json:"dataCenter"`
	Duration              string       `json:"duration"`
	Fee                   string       `json:"fee"`
	Hadoopversion         string       `json:"hadoopVersion"`
	Masternodesize        string       `json:"masterNodeSize"`
	Corenodesize          string       `json:"coreNodeSize"`
	Componentlist         []Component  `json:"componentList"`
	Externalip            string       `json:"externalIp"`
	Externalalternateip   string       `json:"externalAlternateIp"`
	Internalip            string       `json:"internalIp"`
	Deploymentid          string       `json:"deploymentId"`
	Remark                string       `json:"remark"`
	Orderid               string       `json:"orderId"`
	Azid                  string       `json:"azId"`
	Azname                string       `json:"azName"`
	Masternodeproductid   string       `json:"masterNodeProductId"`
	Masternodespecid      string       `json:"masterNodeSpecId"`
	Corenodeproductid     string       `json:"coreNodeProductId"`
	Corenodespecid        string       `json:"coreNodeSpecId"`
	Instanceid            string       `json:"instanceId"`
	Vnc                   string       `json:"vnc"`
	Tenantid              string       `json:"tenantId"`
	Volumesize            int          `json:"volumeSize"`
	Vpc                   string       `json:"vpc"`
	Vpcid                 string       `json:"vpcId"`
	Subnetid              string       `json:"subnetId"`
	Subnetname            string       `json:"subnetName"`
	Securitygroupsid      string       `json:"securityGroupsId"`
	Slavesecuritygroupsid string       `json:"slaveSecurityGroupsId"`
	Stagedesc             string       `json:"stageDesc"`
	Safemode              int          `json:"safeMode"`
	Clusterversion        string       `json:"clusterVersion"`
	Nodepubliccertname    string       `json:"nodePublicCertName"`
	Masternodeip          string       `json:"masterNodeIp"`
	Privateipfirst        string       `json:"privateIpFirst"`
	Errorinfo             string       `json:"errorInfo"`
	Chargingstarttime     string       `json:"chargingStartTime"`
	LogCollection         int          `json:"logCollection"`
	TaskNodeGroups        []NodeGroup  `json:"taskNodeGroups"`
	NodeGroups            []NodeGroup  `json:"nodeGroups"`
	MasterDataVolumeType  string       `json:"masterDataVolumeType"`
	MasterDataVolumeSize  int          `json:"masterDataVolumeSize"`
	MasterDataVolumeCount int          `json:"masterDataVolumeCount"`
	CoreDataVolumeType    string       `json:"coreDataVolumeType"`
	CoreDataVolumeSize    int          `json:"coreDataVolumeSize"`
	CoreDataVolumeCount   int          `json:"coreDataVolumeCount"`
	BootstrapScripts      []ScriptOpts `json:"bootstrapScripts"`
}

type Component struct {
	Componentid      string `json:"componentId"`
	Componentname    string `json:"componentName"`
	Componentversion string `json:"componentVersion"`
	Componentdesc    string `json:"componentDesc"`
}

type NodeGroup struct {
	GroupName       string `json:"groupName"`
	NodeNum         int    `json:"nodeNum"`
	NodeSize        string `json:"nodeSize"`
	NodeSpecId      string `json:"nodeSpecId"`
	NodeProductId   string `json:"nodeProductId"`
	VMProductId     string `json:"vmProductId"`
	VMSpecCode      string `json:"vmSpecCode"`
	RootVolumeSize  int    `json:"rootVolumeSize"`
	RootVolumeType  string `json:"rootVolumeType"`
	DataVolumeType  string `json:"dataVolumeType"`
	DataVolumeSize  int    `json:"dataVolumeSize"`
	DataVolumeCount int    `json:"dataVolumeCount"`
}

type ClusterResult struct {
	ClusterID string `json:"cluster_id"`
	Result    bool   `json:"result"`
	Msg       string `json:"msg"`
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

type DeleteResult struct {
	golangsdk.ErrResult
}
