package cluster

import "github.com/huaweicloud/golangsdk"

type Cluster struct {
	Clusterid             string       `json:"clusterId"`
	Clustername           string       `json:"clusterName"`
	Masternodenum         string       `json:"masterNodeNum"`
	Corenodenum           string       `json:"coreNodeNum"`
	Clusterstate          string       `json:"clusterState"`
	Createat              string       `json:"createAt"`
	Updateat              string       `json:"updateAt"`
	Billingtype           string       `json:"billingType"`
	Datacenter            string       `json:"dataCenter"`
	Vpc                   string       `json:"vpc"`
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
	Masternodeproductid   string       `json:"masterNodeProductId"`
	Masternodespecid      string       `json:"masterNodeSpecId"`
	Corenodeproductid     string       `json:"coreNodeProductId"`
	Corenodespecid        string       `json:"coreNodeSpecId"`
	Azname                string       `json:"azName"`
	Instanceid            string       `json:"instanceId"`
	Vnc                   string       `json:"vnc"`
	Tenantid              string       `json:"tenantId"`
	Volumesize            int          `json:"volumeSize"`
	Subnetname            string       `json:"subnetName"`
	Securitygroupsid      string       `json:"securityGroupsId"`
	Slavesecuritygroupsid string       `json:"slaveSecurityGroupsId"`
	Safemode              int          `json:"safeMode"`
	Clusterversion        string       `json:"clusterVersion"`
	Nodepubliccertname    string       `json:"nodePublicCertName"`
	Masternodeip          string       `json:"masterNodeIp"`
	Privateipfirst        string       `json:"privateIpFirst"`
	Errorinfo             string       `json:"errorInfo"`
	Chargingstarttime     string       `json:"chargingStartTime"`
	LogCollection         int          `json:"log_collection"`
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
