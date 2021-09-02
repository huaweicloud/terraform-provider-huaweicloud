package blockchains

import "github.com/chnsz/golangsdk"

type commonResult struct {
	golangsdk.Result
}

//CreateResult is a struct that represents the result of CreateNewBlockchain
type CreateResult struct {
	commonResult
}

type CreateResponse struct {
	ID   string `json:"blockchain_id"`
	Name string `json:"blockchain_name"`
}

func (r CreateResult) Extract() (*CreateResponse, error) {
	var res CreateResponse
	err := r.ExtractInto(&res)
	return &res, err
}

//DeleteResult is a struct that represents the result of DeleteBlockchain
type DeleteResult struct {
	commonResult
}

func (r DeleteResult) Extract() error {
	return r.Err
}

//ShowResult is a struct that represents the result of ShowBlockchainDetail
type ShowResult struct {
	commonResult
}

type BCSInstance struct {
	Basic      Basic     `json:"basic_info"`
	Channels   []Channel `json:"channels"`
	Peer       []Peer    `json:"peer_info"`
	LightPeer  []Peer    `json:"loght_peer_info"`
	Orderer    Peer      `json:"orderer_info"`
	CouchDB    CouchDB   `json:"couch_db_info"`
	DMSKafka   DMSKafka  `json:"dms_kafka_info"`
	IEF        IEF       `json:"ief_info"`
	SFS        SFS       `json:"sfs_info"`
	Agent      Peer      `json:"agent_info"`
	RestfulAPI Peer      `json:"restapi_info"`
	PVC        PVC       `json:"evs_pvc_info"`
	TaskServer Peer      `json:"tc3_taskserver_info"`
	OBS        OBS       `json:"obs_bucket_info"`
}

type Basic struct {
	ID                       string       `json:"id"`
	Name                     string       `json:"name"`
	KernelType               string       `json:"kernel_type"`
	Version                  string       `json:"version"`
	VersionType              int          `json:"version_type"`
	VolumeType               string       `json:"volume_type"`
	ServiceType              string       `json:"service_type"`
	PurchaseType             string       `json:"purchase_type"`
	SignAlgorithm            string       `json:"sign_algorithm"`
	Consensus                string       `json:"consensus"`
	ChargingMode             int          `json:"charging_mode"`
	DatabaseType             string       `json:"database_type"`
	ClusterID                string       `json:"cluster_id"`
	ClusterName              string       `json:"cluster_name"`
	ClusterType              string       `json:"cluster_type"`
	ClusterPlatformType      string       `json:"cluster_platform_type"`
	ClusterAvailabilityZone  string       `json:"cluster_az"`
	CreatedTime              string       `json:"created_time"`
	DeployType               string       `json:"deploy_type"`
	DeployScale              int          `json:"deploy_scale"`
	DeployStatus             int          `json:"deploy_status"`
	DetailStatus             DetailStatus `json:"detail_status"`
	IsCrossRegion            bool         `json:"is_cross_region"`
	IsSupportRollback        bool         `json:"is_support_rollback"`
	IsSupportRestful         bool         `json:"is_support_restful"`
	IsSupportTc3             bool         `json:"is_support_tc3"`
	IsOldService             bool         `json:"is_old_service"`
	OldServiceVersion        string       `json:"old_service_version"`
	AgentPortalAddress       []string     `json:"agent_portal_addrs"`
	Status                   string       `json:"status"`
	ProcessStatus            string       `json:"process_status"`
	Tc3TaskServerPortalAddrs []string     `json:"tc3_taskserver_portal_addrs"`
	TotalDeployPeer          int          `json:"total_deploy_peer"`
	OrderStatus              int          `json:"order_status"`
	OrderInfo                OrderInfo    `json:"order_info"`
	OrderFadeCache           int          `json:"order_fade_cache"`
	OrderFadeEnable          bool         `json:"order_fade_enable"`
	IEFClusterInfo           IEFCluster   `json:"ief_cluster_info"`
	IEFAPIVersion            string       `json:"iefapi_version"`
}

type Channel struct {
	Name        string              `json:"name"`
	OrgNames    []string            `json:"org_names"`
	OrgNameHash []string            `json:"org_name_hash"`
	Peers       map[string][]string `json:"peers"`
}

type Peer struct {
	Name         string        `json:"name"`
	NodeCount    int           `json:"node_cnt"`
	Status       string        `json:"status"`
	StatusDetail string        `json:"status_detail"`
	PVCName      string        `json:"pvc_name"`
	Address      []PeerAddress `json:"address"`
}

type PeerAddress struct {
	DomainPort string `json:"domain_port"`
	IPPort     string `json:"ip_port"`
}

type CouchDB struct {
	User string `json:"user"`
}

type DMSKafka struct {
	Address         []string `json:"addr"`
	Name            string   `json:"name"`
	Status          string   `json:"status"`
	NodeCount       int      `json:"node_cnt"`
	OrderFadeEnable bool     `json:"order_fade_enable"`
	OrderFadeCache  int      `json:"order_fade_cache"`
}

type IEF struct {
	DeployMode int `json:"deploy_mode"`
}

type IEFCluster struct {
	GroupID      string `json:"group_id"`
	GroupName    string `json:"group_name"`
	InstanceID   string `json:"instance_id"`
	InstanceName string `json:"instance_name"`
}

type SFS struct {
	Name    string `json:"name"`
	PVCName string `json:"pvc_name"`
	Address string `json:"addr"`
	Type    string `json:"type"`
}

type PVC struct {
	DeployMode int `json:"deploy_mode"`
}

type OBS struct {
	Name    string `json:"name"`
	Address string `json:"addr"`
}

type DetailStatus struct {
	AgentStatus     string `json:"agent_status"`
	ConsensusStatus string `json:"consensus_status"`
	OrgStatus       string `json:"org_status"`
	PeerStatus      string `json:"peer_status"`
	PluginStatus    string `json:"plugin_status"`
}

type OrderInfo struct {
	Delete            int    `json:"delete"`
	Operate           int    `json:"operate"`
	OrderID           string `json:"order_id"`
	OrderStatus       int    `json:"order_status"`
	OrderType         int    `json:"order_type"`
	Release           int    `json:"release"`
	ResourceErrorCode string `json:"resource_error_code"`
	ResourceStatus    int    `json:"resource_status"`
}

func (r ShowResult) Extract() (*BCSInstance, error) {
	var res BCSInstance
	err := r.ExtractInto(&res)
	return &res, err
}

//ListResult is a struct that represents the result of ListBlockchain
type ListResult struct {
	commonResult
}

type BlockChain struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (r ListResult) Extract() (*[]BlockChain, error) {
	var s struct {
		BlockChains []BlockChain `json:"blockchains"`
	}
	err := r.ExtractInto(&s)
	return &s.BlockChains, err
}

//StatusResult is a struct that represents the result of ShowBlockchainStatus
type StatusResult struct {
	commonResult
}

type Status struct {
	BCSStatus   StatusDetail `json:"bcs"`
	EIPStatus   StatusDetail `json:"eip"`
	SFSStatus   StatusDetail `json:"sfs"`
	OBSStatus   StatusDetail `json:"obs"`
	KafkaStatus StatusDetail `json:"kafka"`
	CCEStatus   CCEEngine    `json:"cce"`
}

type CCEEngine struct {
	Cluster       StatusDetail `json:"cluster"`
	Network       StatusDetail `json:"network"`
	SecurityGroup StatusDetail `json:"security_group"`
}

type StatusDetail struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
	Detail    string `json:"detail"`
}

func (r StatusResult) Extract() (*Status, error) {
	var s Status
	err := r.ExtractInto(&s)
	return &s, err
}

//NodesResult is a struct that represents the result of ShowBlockchainNode
type NodesResult struct {
	commonResult
}

type Org struct {
	OrgMSPID  string          `json:"org_msp_id"`
	OrgDomain string          `json:"org_domain"`
	Peers     map[string]Node `json:"peers"`
}

type Node struct {
	Port     string   `json:"ip_port"`
	Channels []string `json:"channels"`
}

func (r NodesResult) Extract() (*map[string]Org, error) {
	var s struct {
		NodeOrgs map[string]Org `json:"node_orgs"`
	}
	err := r.ExtractInto(&s)
	return &s.NodeOrgs, err
}

//UpdateResult is a struct which represents the result of UpdateBlockchain
type UpdateResult struct {
	commonResult
}

func (r UpdateResult) Extract() error {
	return r.Err
}
