package instances

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
)

// InstanceCreate response
type InstanceCreate struct {
	InstanceID string `json:"instance_id"`
}

// CreateResult is a struct that contains all the return parameters of creation
type CreateResult struct {
	golangsdk.Result
}

// Extract from CreateResult
func (r CreateResult) Extract() (*InstanceCreate, error) {
	var s InstanceCreate
	err := r.Result.ExtractInto(&s)
	return &s, err
}

// DeleteResult is a struct which contains the result of deletion
type DeleteResult struct {
	golangsdk.ErrResult
}

type ListResponse struct {
	Instances  []Instance `json:"instances"`
	TotalCount int        `json:"instance_num"`
}

// Instance response
type Instance struct {
	Name                     string   `json:"name"`
	Description              string   `json:"description"`
	Engine                   string   `json:"engine"`
	EngineVersion            string   `json:"engine_version"`
	Specification            string   `json:"specification"`
	StorageSpace             int      `json:"storage_space"`
	PartitionNum             string   `json:"partition_num"`
	BrokerNum                int      `json:"broker_num"`
	NodeNum                  int      `json:"node_num"`
	UsedStorageSpace         int      `json:"used_storage_space"`
	ConnectAddress           string   `json:"connect_address"`
	Port                     int      `json:"port"`
	Status                   string   `json:"status"`
	InstanceID               string   `json:"instance_id"`
	ResourceSpecCode         string   `json:"resource_spec_code"`
	ChargingMode             int      `json:"charging_mode"`
	VPCID                    string   `json:"vpc_id"`
	VPCName                  string   `json:"vpc_name"`
	CreatedAt                string   `json:"created_at"`
	UserID                   string   `json:"user_id"`
	UserName                 string   `json:"user_name"`
	OrderID                  string   `json:"order_id"`
	MaintainBegin            string   `json:"maintain_begin"`
	MaintainEnd              string   `json:"maintain_end"`
	EnablePublicIP           bool     `json:"enable_publicip"`
	ManagementConnectAddress string   `json:"management_connect_address"`
	SslEnable                bool     `json:"ssl_enable"`
	KafkaSecurityProtocol    string   `json:"kafka_security_protocol"`
	SaslEnabledMechanisms    []string `json:"sasl_enabled_mechanisms"`
	EnterpriseProjectID      string   `json:"enterprise_project_id"`
	IsLogicalVolume          bool     `json:"is_logical_volume"`
	ExtendTimes              int      `json:"extend_times"`
	EnableAutoTopic          bool     `json:"enable_auto_topic"`
	Type                     string   `json:"type"`
	ProductID                string   `json:"product_id"`
	SecurityGroupID          string   `json:"security_group_id"`
	SecurityGroupName        string   `json:"security_group_name"`
	SubnetID                 string   `json:"subnet_id"`
	SubnetName               string   `json:"subnet_name"`
	AvailableZones           []string `json:"available_zones"`
	TotalStorageSpace        int      `json:"total_storage_space"`
	PublicConnectionAddress  string   `json:"public_connect_address"`
	StorageResourceID        string   `json:"storage_resource_id"`
	StorageSpecCode          string   `json:"storage_spec_code"`
	ServiceType              string   `json:"service_type"`
	StorageType              string   `json:"storage_type"`
	RetentionPolicy          string   `json:"retention_policy"`
	KafkaPublicStatus        string   `json:"kafka_public_status"`
	PublicBandWidth          int      `json:"public_bandwidth"`
	KafkaManagerUser         string   `json:"kafka_manager_user"`
	EnableLogCollect         bool     `json:"enable_log_collection"`
	CrossVpcInfo             string   `json:"cross_vpc_info"`
	Ipv6Enable               bool     `json:"ipv6_enable"`
	Ipv6ConnectAddresses     []string `json:"ipv6_connect_addresses"`
	ConnectorEnalbe          bool     `json:"connector_enable"`
	ConnectorID              string   `json:"connector_id"`
	ConnectorNodeNum         int      `json:"connector_node_num"`
	RestEnable               bool     `json:"rest_enable"`
	RestConnectAddress       string   `json:"rest_connect_address"`
	MessageQueryInstEnable   bool     `json:"message_query_inst_enable"`
	VpcClientPlain           bool     `json:"vpc_client_plain"`
	SupportFeatures          string   `json:"support_features"`
	Task                     Task     `json:"task"`
	TraceEnable              bool     `json:"trace_enable"`
	PodConnectAddress        string   `json:"pod_connect_address"`
	// Whether disk encryption is enabled.
	DiskEncrypted bool `json:"disk_encrypted"`
	// The key ID of the disk encryption.
	DiskEncryptedKey           string             `json:"disk_encrypted_key"`
	KafkaPrivateConnectAddress string             `json:"kafka_private_connect_address"`
	CesVersion                 string             `json:"ces_version"`
	AccessUser                 string             `json:"access_user"`
	Tags                       []tags.ResourceTag `json:"tags"`
	CertReplaced               bool               `json:"cert_replaced"`
	SslTwoWayEnable            bool               `json:"ssl_two_way_enable"`
	PortProtocols              PortProtocols      `json:"port_protocols"`
}

type Task struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type PortProtocols struct {
	PrivatePlainEnable             bool   `json:"private_plain_enable"`
	PrivatePlainAddress            string `json:"private_plain_address"`
	PrivatePlainDomainName         string `json:"private_plain_domain_name"`
	PrivateSaslSslEnable           bool   `json:"private_sasl_ssl_enable"`
	PrivateSaslSslAddress          string `json:"private_sasl_ssl_address"`
	PrivateSaslSslDomainName       string `json:"private_sasl_ssl_domain_name"`
	PrivateSaslPlaintextEnable     bool   `json:"private_sasl_plaintext_enable"`
	PrivateSaslPlaintextAddress    string `json:"private_sasl_plaintext_address"`
	PrivateSaslPlaintextDomainName string `json:"private_sasl_plaintext_domain_name"`
	PublicPlainEnable              bool   `json:"public_plain_enable"`
	PublicPlainAddress             string `json:"public_plain_address"`
	PublicPlainDomainName          string `json:"public_plain_domain_name"`
	PublicSaslSslEnable            bool   `json:"public_sasl_ssl_enable"`
	PublicSaslSslAddress           string `json:"public_sasl_ssl_address"`
	PublicSaslSslDomainName        string `json:"public_sasl_ssl_domain_name"`
	PublicSaslPlaintextEnable      bool   `json:"public_sasl_plaintext_enable"`
	PublicSaslPlaintextAddress     string `json:"public_sasl_plaintext_address"`
	PublicSaslPlaintextDomainName  string `json:"public_sasl_plaintext_domain_name"`
}

// UpdateResult is a struct from which can get the result of update method
type UpdateResult struct {
	golangsdk.Result
}

// GetResult contains the body of getting detailed
type GetResult struct {
	golangsdk.Result
}

// Extract from GetResult
func (r GetResult) Extract() (*Instance, error) {
	var s Instance
	err := r.Result.ExtractInto(&s)
	return &s, err
}

type Page struct {
	pagination.SinglePageBase
}

func (r Page) IsEmpty() (bool, error) {
	data, err := ExtractInstances(r)
	if err != nil {
		return false, err
	}
	return len(data.Instances) == 0, err
}

// ExtractCloudServers is a function that takes a ListResult and returns the services' information.
func ExtractInstances(r pagination.Page) (ListResponse, error) {
	var s ListResponse
	err := (r.(Page)).ExtractInto(&s)
	return s, err
}

// CrossVpc is the structure that represents the API response of 'UpdateCrossVpc' method.
type CrossVpc struct {
	// The result of cross-VPC access modification.
	Success bool `json:"success"`
	// The result list of broker cross-VPC access modification.
	Connections []Connection `json:"results"`
}

// Connection is the structure that represents the detail of the cross-VPC access.
type Connection struct {
	// advertised.listeners IP/domain name.
	AdvertisedIp string `json:"advertised_ip"`
	// The status of broker cross-VPC access modification.
	Success bool `json:"success"`
	// Listeners IP.
	ListenersIp string `json:"ip"`
}

// AutoTopicResult is a struct that contains all the return parameters of UpdateAutoTopic function
type AutoTopicResult struct {
	golangsdk.Result
}

// ResetPasswordResult is a struct that contains all the return parameters of ResetPassword function
type ResetPasswordResult struct {
	golangsdk.Result
}

type commonResult struct {
	golangsdk.Result
}

type ModifyConfigurationResult struct {
	commonResult
}

type GetConfigurationResult struct {
	commonResult
}

type RebootResult struct {
	commonResult
}

type GetTasksResult struct {
	commonResult
}

type ModifyConfigurationResp struct {
	JobId         string `json:"job_id"`
	DynamicConfig int    `json:"dynamic_config"`
	StaticConfig  int    `json:"static_config"`
}

func (r ModifyConfigurationResult) Extract() (*ModifyConfigurationResp, error) {
	var response ModifyConfigurationResp
	err := r.ExtractInto(&response)
	return &response, err
}

type Result struct {
	Result   string `json:"result"`
	Instance string `json:"instance"`
}

type RebootResp struct {
	Results []Result `json:"results"`
}

func (r RebootResult) Extract() (*RebootResp, error) {
	var response RebootResp
	err := r.ExtractInto(&response)
	return &response, err
}

type KafkaParam struct {
	Name         string `json:"name"`
	Value        string `json:"value"`
	DefaultValue string `json:"default_value"`
	ConfigType   string `json:"config_type"`
	ValidValues  string `json:"valid_values"`
	ValueType    string `json:"value_type"`
}

type GetConfigurationResp struct {
	KafkaConfigs []KafkaParam `json:"kafka_configs"`
}

func (r GetConfigurationResult) Extract() (*GetConfigurationResp, error) {
	var response GetConfigurationResp
	err := r.ExtractInto(&response)
	return &response, err
}

type TaskParams struct {
	Name   string `json:"name"`
	Params string `json:"params"`
	Status string `json:"status"`
}

type GetTaskResp struct {
	Tasks []TaskParams `json:"tasks"`
}

func (r GetTasksResult) Extract() (*GetTaskResp, error) {
	var response GetTaskResp
	err := r.ExtractInto(&response)
	return &response, err
}
