package blockchains

import (
	"github.com/chnsz/golangsdk"
)

//CreateOpts is a struct which will be used to create a bcs instance
type CreateOpts struct {
	Name                string          `json:"name" required:"true"`
	ClusterType         string          `json:"cluster_type" required:"true"`
	CreateNewCluster    *bool           `json:"create_new_cluster" required:"true"`
	EnterpriseProjectId string          `json:"enterprise_project_id" required:"true"`
	FabricVersion       string          `json:"fabric_version" required:"true"`
	Password            string          `json:"resource_password" required:"true"`
	VersionType         int             `json:"version_type" required:"true"`
	BlockChainType      string          `json:"blockchain_type,omitempty"`
	Consensus           string          `json:"consensus,omitempty"`
	SignAlgorithm       string          `json:"sign_algorithm,omitempty"`
	VolumeType          string          `json:"volume_type,omitempty"`
	EvsDiskType         string          `json:"evs_disk_type,omitempty"`
	OrgDiskSize         int             `json:"org_disk_size,omitempty"`
	DatabaseType        string          `json:"database_type,omitempty"`
	OrdererNodeNumber   int             `json:"orderer_node_number,omitempty"`
	EIPEnable           bool            `json:"use_eip,omitempty"`
	BandwidthSize       int             `json:"bandwidth_size,omitempty"`
	CCEClusterInfo      *CCEClusterInfo `json:"cce_cluster_info,omitempty"`
	CCECreateInfo       *CCECreateInfo  `json:"cce_create_info,omitempty"`
	IEFDeployMode       int             `json:"ief_deploy_mode,omitempty"`
	IEFNodesInfo        []IEFNode       `json:"ief_nodes_info,omitempty"`
	PeerOrgs            []PeerOrg       `json:"peer_orgs,omitempty"`
	Channels            []ChannelInfo   `json:"channels,omitempty"`
	CouchDBInfo         *CouchDBInfo    `json:"couchdb_info,omitempty"`
	SFSTurbo            *SFSTurbo       `json:"turbo_info,omitempty"`
	Block               *BlockInfo      `json:"block_info,omitempty"`
	Kafka               *KafkaInfo      `json:"kafka_create_info,omitempty"`
	TC3Need             bool            `json:"tc3_need,omitempty"`
	RestfulAPISupport   bool            `json:"restful_api_support,omitempty"`
	IsInvitee           bool            `json:"is_invitee,omitempty"`
	InvitorInfo         *InvitorInfo    `json:"invitor_infos,omitempty"`
}

//CCEClusterInfo is the CCE cluster struct that will be used to associate when creating a bcs instance
type CCEClusterInfo struct {
	ID   string `json:"cluster_id" required:"true"`
	Name string `json:"cluster_name" required:"true"`
}

//CCECreateInfo is the struct that will be used to specify the creation of a new CCE cluster
//when creating a bcs instance
type CCECreateInfo struct {
	NodeNum          int    `json:"node_num" required:"true"`
	Flavor           string `json:"node_flavor" required:"true"`
	ClusterFlavor    string `json:"cce_flavor" required:"true"`
	Password         string `json:"init_node_pwd" required:"true"`
	AvailabilityZone string `json:"az" required:"true"`
	PlatformType     string `json:"cluster_platform_type" required:"true"`
}

//IEFNode is the IEF node struct that will be used to associate when creating a bcs instance
type IEFNode struct {
	ID        string `json:"id" required:"true"`
	Status    string `json:"status" required:"true"`
	IPAddress string `json:"public_ip_address" required:"true"`
}

//PeerOrg is the peer organization struct that will be used to creating a bcs instance
type PeerOrg struct {
	Name      string `json:"name" required:"true"`
	NodeCount int    `json:"node_count" required:"true"`
}

//ChannelInfo is the channel struct that will be used to creating a bcs instance
type ChannelInfo struct {
	Name        string   `json:"name" required:"true"`
	OrgNames    []string `json:"org_names" required:"true"`
	Description string   `json:"desctiption,omitempty"`
}

//CouchDBInfo is the couch database struct that will be used to creating a bcs instance
type CouchDBInfo struct {
	UserName string `json:"user_name" required:"true"`
	Password string `json:"password" required:"true"`
}

//SFSTurbo is the turbo struct that will be used to creating a bcs instance
type SFSTurbo struct {
	ShareType        string `json:"share_type" required:"true"`
	Type             string `json:"type" required:"true"`
	AvailabilityZone string `json:"availability_zone" required:"true"`
	Flavor           string `json:"resource_spec_code" required:"true"`
}

//BlockInfo is the turbo struct that will be used to creating a bcs instance
type BlockInfo struct {
	BatchTimeout      int `json:"batch_timeout,omitempty"`
	MaxMessageCount   int `json:"max_message_count,omitempty"`
	PreferredMaxbytes int `json:"preferred_maxbytes,omitempty"`
}

//KafkaInfo is the block generation struct that be used to config when creating a bcs instance
type KafkaInfo struct {
	Flavor           string `json:"spec" required:"true"`
	Storage          int    `json:"storage" required:"true"`
	AvailabilityZone string `json:"az" required:"true"`
}

//InvitorInfo is the invitor struct that be used to config when creating a bcs instance
type InvitorInfo struct {
	TenantID     string `json:"tenant_id" required:"true"`
	ProjectID    string `json:"project_id" required:"true"`
	BlockchainID string `json:"blockchain_id" required:"true"`
}

type CreateOptsBuilder interface {
	ToInstancesCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToInstancesCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

//Create is a method by which can be able to access the create function that create a bcs instance
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToInstancesCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(rootURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}

//DeleteOpts is a struct which will be used to delete an existing bcs instance
type DeleteOpts struct {
	IsDeleteStorage  bool `q:"is_delete_storage"`
	IsDeleteOBS      bool `q:"is_delete_obs"`
	IsDeleteResource bool `q:"is_delete_resource"`
}

type DeleteOptsBuilder interface {
	ToInstanceDeleteQuery() (string, error)
}

func (opts DeleteOpts) ToInstanceDeleteQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

//Delete is a method to delete an existing bcs instance
func Delete(client *golangsdk.ServiceClient, opts DeleteOptsBuilder, id string) (r DeleteResult) {
	url := resourceURL(client, id)
	if opts != nil {
		query, err := opts.ToInstanceDeleteQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = client.Delete(url, &golangsdk.RequestOpts{
		OkCodes:      []int{200, 202, 204},
		JSONResponse: nil,
		MoreHeaders:  map[string]string{"Content-Type": "application/json"},
	})
	return
}

//Get is a method to obtain the detailed information of an existing bcs instance
func Get(client *golangsdk.ServiceClient, id string) (r ShowResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)
	return
}

//GetStatus is a method to obtain all block status of an existing bcs instance
func GetStatus(client *golangsdk.ServiceClient, id string) (r StatusResult) {
	_, r.Err = client.Get(extraURL(client, id, "status"), &r.Body, nil)
	return
}

//List is a method to obtain the detailed information list of all existing bcs instance
func List(client *golangsdk.ServiceClient) (r ListResult) {
	_, r.Err = client.Get(rootURL(client), &r.Body, nil)
	return
}

//GetNodes is a method to obtain the node information list of an existing bcs instance
func GetNodes(client *golangsdk.ServiceClient, id string) (r NodesResult) {
	_, r.Err = client.Get(extraURL(client, id, "nodes"), &r.Body, nil)
	return
}

//UpdateOpts is a struct which will be used to update an existing bcs instance
type UpdateOpts struct {
	NodePeer  []NodePeer `json:"node_orgs" required:"true"`
	PublicIPs []IEFNode  `json:"publicips,omitempty"`
}

//NodePeer is the peer organization struct that will be used to add a peer organization to an existing bcs instance
type NodePeer struct {
	Name    string `json:"name" required:"true"`
	Count   int    `json:"node_count" required:"true"`
	PVCName string `json:"pvc_name,omitempty"`
}

type UpdateOptsBuilder interface {
	ToInstancesUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateOpts) ToInstancesUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

//Update is a method to update an existing bcs instance
func Update(client *golangsdk.ServiceClient, opts UpdateOptsBuilder, id string) (r UpdateResult) {
	b, err := opts.ToInstancesUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(resourceURL(client, id), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}
