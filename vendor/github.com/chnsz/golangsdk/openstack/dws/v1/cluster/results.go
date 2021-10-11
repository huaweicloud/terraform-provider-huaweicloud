package cluster

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

type Cluster struct {
	Status              string             `json:"status"`
	SubStatus           string             `json:"sub_status"`
	Updated             string             `json:"updated"` //ISO8601:YYYY-MM-DDThh:mm:ssZ
	Endpoints           []Endpoints        `json:"endPoints"`
	Name                string             `json:"name"`
	NumberOfNode        int                `json:"number_of_node"`
	AvailabilityZone    string             `json:"availability_zone"`
	SubnetID            string             `json:"subnet_id"`
	PublicEndpoints     []PublicEndpoints  `json:"public_endpoints"`
	Created             string             `json:"created"` //ISO8601:YYYY-MM-DDThh:mm:ssZ
	SecurityGroupID     string             `json:"security_group_id"`
	Port                int                `json:"port"`
	NodeType            string             `json:"node_type"`
	Version             string             `json:"version"`
	PublicIp            *PublicIp          `json:"public_ip"`
	FailedReasons       *FailInfo          `json:"failed_reasons"`
	VpcID               string             `json:"vpc_id"`
	TaskStatus          string             `json:"task_status"`
	UserName            string             `json:"user_name"`
	ID                  string             `json:"id"`
	ActionProgress      map[string]string  `json:"action_progress"`
	RecentEvent         int                `json:"recent_event"`
	Tags                []tags.ResourceTag `json:"tags"`
	EnterpriseProjectId string             `json:"enterprise_project_id"`
}

type ClusterDetail struct {
	Cluster
	PrivateIp      []string        `json:"private_ip"`
	ParameterGroup *ParameterGroup `json:"parameter_group"`
	NodeTypeId     string          `json:"node_type_id"`
	NodeDetail     *NodeDetail     `json:"node_detail"`
	MaintainWindow *MaintainWindow `json:"maintain_window"`
	ResizeInfo     *ResizeInfo     `json:"resize_info"`
}

type ClusterDetailsRst struct {
	Cluster ClusterDetail `json:"cluster"`
}

type CreateClusterRst struct {
	Cluster IdObject `json:"cluster"`
}

type IdObject struct {
	Id string `json:"id"`
}

type ListClustersRst struct {
	Clusters []Cluster `json:"clusters"`
}

type FailInfo struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

type PublicIp struct {
	EipID          string `json:"eip_id"`
	PublicBindType string `json:"public_bind_type"`
}

type FailedReasons struct {
	FailInfo FailInfo `json:"fail_info"`
}

type Endpoints struct {
	ConnectInfo string `json:"connect_info"`
	JdbcUrl     string `json:"jdbc_url"`
}

type PublicEndpoints struct {
	PublicConnectInfo string `json:"public_connect_info"`
	JdbcUrl           string `json:"jdbc_url"`
}

type ParameterGroup struct {
	Name   string `json:"name"`
	Id     string `json:"id"`
	Status string `json:"status"` // In-Sync,Applying.Pending-Reboot,Sync-Failure
}

type NodeDetail struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type MaintainWindow struct {
	StartTime string `json:"start_time"` // HH:mm,timezoue:GMT+0。
	EndTime   string `json:"end_time"`
	Day       string `json:"day"`
}

type ResizeInfo struct {
	ResizeStatus  string `json:"resize_status"` //GROWING,RESIZE_FAILURE
	StartTime     string `json:"start_time"`    //ISO8601:YYYY-MM-DDThh:mm:ss
	TargetNodeNum int    `json:"target_node_num"`
	OriginNodeNum int    `json:"origin_node_num"`
}

type GetResult struct {
	golangsdk.Result
}

type CreateRsp struct {
	ID string `json:"id"`
}

type CreateResult struct {
	golangsdk.Result
}

type DeleteResult struct {
	golangsdk.ErrResult
}
