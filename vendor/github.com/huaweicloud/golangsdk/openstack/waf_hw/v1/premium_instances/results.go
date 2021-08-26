/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package premium_instances

// CreationgRst the struct of returned by creating.
type CreationgRst struct {
	Instances []IdNameEntry `json:"instances"`
}

type IdNameEntry struct {
	Id   string
	Name string
}

// DedicatedInstance the dedicated waf instance detail.
type DedicatedInstance struct {
	Id                 string            `json:"id"`
	InstanceName       string            `json:"instancename"`
	ServerId           string            `json:"serverId"`
	Region             string            `json:"region"`
	Zone               string            `json:"zone"`
	Arch               string            `json:"arch"`
	CupFlavor          string            `json:"cpu_flavor"`
	VpcId              string            `json:"vpc_id"`
	SubnetId           string            `json:"subnet_id"`
	ServiceIp          string            `json:"service_ip"`
	ServiceIpv6        string            `json:"service_ipv6"`
	FloatIp            string            `json:"floatIp"`
	SecurityGroupIds   []string          `json:"security_group_ids"`
	MgrSecurityGroupId string            `json:"mgrSecurityGroupId"`
	Status             int               `json:"status"`
	RunStatus          int               `json:"run_status"`
	AccessStatus       int               `json:"access_status"`
	Upgradable         int               `json:"upgradable"`
	CloudServiceType   string            `json:"cloudServiceType"`
	ResourceType       string            `json:"resourceType"`
	ResourceSpecCode   string            `json:"resourceSpecCode"`
	Specification      string            `json:"specification"`
	Hosts              []IdHostnameEntry `json:"hosts"`
	VolumeType         string            `json:"volume_type"`
	ClusterId          string            `json:"cluster_id"`
	PoolId             string            `json:"pool_id"`
}

type IdHostnameEntry struct {
	Id       string `json:"id"`
	HostName string `json:"hostname"`
}

// DedicatedInstanceList the struct of returned by querying list.
type DedicatedInstanceList struct {
	Total     int                 `json:"total"`
	Purchased bool                `json:"purchased"`
	Items     []DedicatedInstance `json:"items"`
}
