package instances

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
)

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOpsBuilder is used for creating instance parameters.
// any struct providing the parameters should implement this interface
type CreateOpsBuilder interface {
	ToInstanceCreateMap() (map[string]interface{}, error)
}

// CreateOps is a struct that contains all the parameters.
type CreateOps struct {
	// Indicates the name of an instance.
	// An instance name starts with a letter,
	// consists of 4 to 64 characters, and supports
	// only letters, digits, hyphens (-), and underscores (_).
	Name string `json:"name" required:"true"`

	// Indicates the description of an instance.
	// It is a character string containing not more than 1024 characters.
	Description string `json:"description,omitempty"`

	// Indicates a message engine.
	Engine string `json:"engine" required:"true"`

	// Indicates the version of a message engine.
	EngineVersion string `json:"engine_version" required:"true"`

	// Indicates the message storage space.
	StorageSpace int `json:"storage_space" required:"true"`

	// Indicates the baseline bandwidth of a Kafka instance, that is,
	// the maximum amount of data transferred per unit time. Unit: byte/s.
	Specification string `json:"specification,omitempty"`

	// Indicates the maximum number of brokers in a Kafka instance.
	BrokerNum int `json:"broker_num,omitempty"`

	// Indicates the maximum number of topics in a Kafka instance.
	PartitionNum int `json:"partition_num,omitempty"`

	// Indicates a username.
	// A username consists of 1 to 64 characters
	// and supports only letters, digits, and hyphens (-).
	AccessUser string `json:"access_user,omitempty"`

	// Indicates the password of an instance.
	// An instance password must meet the following complexity requirements:
	// Must be 6 to 32 characters long.
	// Must contain at least two of the following character types:
	// Lowercase letters
	// Uppercase letters
	// Digits
	// Special characters (`~!@#$%^&*()-_=+\|[{}]:'",<.>/?)
	Password string `json:"password,omitempty"`

	// Indicates the ID of a VPC.
	VPCID string `json:"vpc_id" required:"true"`

	// Indicates the ID of a security group.
	SecurityGroupID string `json:"security_group_id" required:"true"`

	// Indicates the ID of a subnet.
	SubnetID string `json:"subnet_id" required:"true"`

	// Whether to enable the IPv6.
	Ipv6Enable bool `json:"ipv6_enable,omitempty"`

	// Indicates the ID of an AZ.
	// The parameter value can be left blank or an empty array.
	AvailableZones []string `json:"available_zones" required:"true"`

	// Indicates a product ID.
	ProductID string `json:"product_id" required:"true"`

	// CPU architecture.
	ArchType string `json:"arch_type,omitempty"`

	// Indicates the username for logging in to the Kafka Manager.
	// The username consists of 4 to 64 characters and can contain
	//letters, digits, hyphens (-), and underscores (_).
	KafkaManagerUser string `json:"kafka_manager_user,omitempty"`

	// Indicates the password for logging in to the Kafka Manager.
	// The password must meet the following complexity requirements:
	// Must be a string consisting of 8 to 32 characters.
	// Contains at least three of the following characters:
	// Lowercase letters
	// Uppercase letters
	// Digits
	// Special characters `~!@#$%^&*()-_=+\|[{}];:',<.>/?
	KafkaManagerPassword string `json:"kafka_manager_password,omitempty"`

	// Indicates the time at which a maintenance time window starts.
	// Format: HH:mm:ss
	MaintainBegin string `json:"maintain_begin,omitempty"`

	// Indicates the time at which a maintenance time window ends.
	// Format: HH:mm:ss
	MaintainEnd string `json:"maintain_end,omitempty"`

	// Indicates whether to open the public network access function. Default to false.
	EnablePublicIP bool `json:"enable_publicip,omitempty"`

	// Indicates the bandwidth of the public network.
	PublicBandWidth int `json:"public_bandwidth,omitempty"`

	// Indicates the ID of the Elastic IP address bound to the instance.
	PublicIpID string `json:"publicip_id,omitempty"`

	// Indicates whether to enable SSL-encrypted access.
	SslEnable bool `json:"ssl_enable,omitempty"`

	// Whether to enable vpc client plain
	VpcClientPlain bool `json:"vpc_client_plain,omitempty"`

	// Indicates the protocol to use after SASL is enabled.
	KafkaSecurityProtocol string `json:"kafka_security_protocol,omitempty"`

	// Indicates the authentication mechanisms to use after SASL is enabled.
	SaslEnabledMechanisms []string `json:"sasl_enabled_mechanisms,omitempty"`

	// Indicates the action to be taken when the memory usage reaches the disk capacity threshold. Options:
	// time_base: Automatically delete the earliest messages.
	// produce_reject: Stop producing new messages.
	RetentionPolicy string `json:"retention_policy,omitempty"`

	// Indicates whether to enable dumping.
	ConnectorEnalbe bool `json:"connector_enable,omitempty"`

	// Indicates whether to enable automatic topic creation.
	EnableAutoTopic bool `json:"enable_auto_topic,omitempty"`

	//Indicates the storage I/O specification. For details on how to select a disk type
	StorageSpecCode string `json:"storage_spec_code,omitempty"`

	// Indicates the enterprise project ID.
	EnterpriseProjectID string `json:"enterprise_project_id,omitempty"`

	// Indicates the tags of the instance
	Tags []tags.ResourceTag `json:"tags,omitempty"`

	// Indicates the parameter related to the yearly/monthly billing mode.
	BssParam BssParam `json:"bss_param,omitempty"`

	// The access mode of the Kafka instance.
	PortProtocol *PortProtocol `json:"port_protocol,omitempty"`

	// The private IP addresses of the Kafka instance.
	TenantIps []string `json:"tenant_ips,omitempty"`

	// Indicates whether to enable disk encryption.
	DiskEncryptedEnable bool `json:"disk_encrypted_enable,omitempty"`

	// The key ID of the disk encryption.
	DiskEncryptedKey string `json:"disk_encrypted_key,omitempty"`
}

type BssParam struct {
	// Indicates the charging mode of the instance.
	ChargingMode string `json:"charging_mode" required:"true"`

	// Indicates the charging period unit of the instance
	PeriodType string `json:"period_type,omitempty"`

	// Indicates the charging period of the instance.
	PeriodNum int `json:"period_num,omitempty"`

	// Indicates whether auto renew is enabled.
	IsAutoRenew *bool `json:"is_auto_renew,omitempty"`

	// Indicates whether the order is automatically or manually paid.
	IsAutoPay *bool `json:"is_auto_pay,omitempty"`
}

type PortProtocol struct {
	// Whether to enable private plaintext access. Default to `false`.
	// + true: Enable private plaintext access. Connection address: ip:9092, access protocol: PLAINTEXT.
	// + false: Disable private plaintext access.
	PrivatePlainEnable *bool `json:"private_plain_enable,omitempty"`

	// Whether to enable private ciphertext access mode with the SASL_SSL security protocol. Default to `false`.
	// `private_sasl_ssl_enable` and `private_sasl_plaintext_enable` cannot be set to `true` at the same time.
	PrivateSaslSslEnable *bool `json:"private_sasl_ssl_enable,omitempty"`

	// Whether to enable private SASL plaintext access. Default to `false`.
	// Connection address: ip:9093, access protocol: SASL_PLAINTEXT.
	// `private_sasl_plaintext_enable` and `private_sasl_ssl_enable` cannot be set to `true` at the same time.
	PrivateSaslPlaintextEnable *bool `json:"private_sasl_plaintext_enable,omitempty"`

	// Whether to enable public plaintext access. Default to `false`.
	// Before enabling public plaintext access, you must enable public network access.
	// Connection address: ip:9094, access protocol: PLAINTEXT.
	PublicPlainEnable *bool `json:"public_plain_enable,omitempty"`

	// Whether to enable public SASL SSL access. Default to `false`.
	// Connection address: ip:9095, access protocol: SASL_SSL.
	// `public_sasl_ssl_enable` and `public_sasl_plaintext_enable` cannot be set to `true` at the same time.
	// When this parameter is set to `true`, the instance must be enabled for public network access.
	PublicSaslSslEnable *bool `json:"public_sasl_ssl_enable,omitempty"`

	// Whether to enable public SASL plaintext access. Default to `false`.
	// Connection address: ip:9095, access protocol: SASL_PLAINTEXT.
	// `public_sasl_plaintext_enable` and `public_sasl_ssl_enable` cannot be set to `true` at the same time.
	// When this parameter is set to `true`, the instance must be enabled for public network access.
	PublicSaslPlaintextEnable *bool `json:"public_sasl_plaintext_enable,omitempty"`
}

// ToInstanceCreateMap is used for type convert
func (ops CreateOps) ToInstanceCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(ops, "")
}

// Create an instance with given parameters.
func Create(client *golangsdk.ServiceClient, ops CreateOpsBuilder) (r CreateResult) {
	b, err := ops.ToInstanceCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})

	return
}

func CreateWithEngine(client *golangsdk.ServiceClient, ops CreateOpsBuilder, engine string) (r CreateResult) {
	b, err := ops.ToInstanceCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURLWithEngine(engine, client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})

	return
}

func CreateInstance(client *golangsdk.ServiceClient, ops CreateOpsBuilder) (r CreateResult) {
	b, err := ops.ToInstanceCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createInstanceURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})

	return
}

// Delete an instance by id
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

// UpdateOptsBuilder is an interface which can build the map paramter of update function
type UpdateOptsBuilder interface {
	ToInstanceUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is a struct which represents the parameters of update function
type UpdateOpts struct {
	// Indicates the name of an instance.
	// An instance name starts with a letter,
	// consists of 4 to 64 characters,
	// and supports only letters, digits, and hyphens (-).
	Name string `json:"name,omitempty"`

	// Indicates the description of an instance.
	// It is a character string containing not more than 1024 characters.
	Description *string `json:"description,omitempty"`

	// Indicates the time at which a maintenance time window starts.
	// Format: HH:mm:ss
	MaintainBegin string `json:"maintain_begin,omitempty"`

	// Indicates the time at which a maintenance time window ends.
	// Format: HH:mm:ss
	MaintainEnd string `json:"maintain_end,omitempty"`

	// Indicates the ID of a security group.
	SecurityGroupID string `json:"security_group_id,omitempty"`

	// Indicates the action to be taken when the memory usage reaches the disk capacity threshold. Options:
	// time_base: Automatically delete the earliest messages.
	// produce_reject: Stop producing new messages.
	RetentionPolicy string `json:"retention_policy,omitempty"`

	// Indicates the enterprise project ID.
	EnterpriseProjectID string `json:"enterprise_project_id,omitempty"`
}

// ToInstanceUpdateMap is used for type convert
func (opts UpdateOpts) ToInstanceUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update is a method which can be able to update the instance
// via accessing to the service with Put method and parameters
func Update(client *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	body, err := opts.ToInstanceUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updateURL(client, id), body, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

// Get a instance with detailed information by id
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

type ListOpts struct {
	InstanceId          string `q:"instance_id"`
	Name                string `q:"name"`
	Engine              string `q:"engine"`
	Status              string `q:"status"`
	IncludeFailure      string `q:"include_failure"`
	ExactMatchName      string `q:"exact_match_name"`
	EnterpriseProjectID string `q:"enterprise_project_id"`
}

type ListOpsBuilder interface {
	ToListDetailQuery() (string, error)
}

func (opts ListOpts) ToListDetailQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, opts ListOpsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToListDetailQuery()

		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageList := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return Page{pagination.SinglePageBase(r)}
	})

	return pageList
}

type ResizeInstanceOpts struct {
	NewSpecCode     *string  `json:"new_spec_code,omitempty"`
	NewStorageSpace *int     `json:"new_storage_space,omitempty"`
	OperType        *string  `json:"oper_type,omitempty"`
	NewBrokerNum    *int     `json:"new_broker_num,omitempty"`
	NewProductID    *string  `json:"new_product_id,omitempty"`
	PublicIpID      *string  `json:"publicip_id,omitempty"`
	TenantIps       []string `json:"tenant_ips,omitempty"`
}

func Resize(client *golangsdk.ServiceClient, id string, opts ResizeInstanceOpts) (string, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return "", err
	}

	var rst golangsdk.Result
	_, err = client.Post(extend(client, id), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	if err == nil {
		var r struct {
			JobID string `json:"job_id"`
		}
		if err = rst.ExtractInto(&r); err != nil {
			return "", err
		}
		return r.JobID, nil
	}
	return "", err
}

func ExtendInstance(client *golangsdk.ServiceClient, instanceId string, opts ResizeInstanceOpts) (string, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return "", err
	}

	var rst golangsdk.Result
	_, err = client.Post(extendInstanceURL(client, instanceId), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	if err == nil {
		var r struct {
			JobID string `json:"job_id"`
		}
		if err = rst.ExtractInto(&r); err != nil {
			return "", err
		}
		return r.JobID, nil
	}
	return "", err
}

// CrossVpcUpdateOpts is the structure required by the UpdateCrossVpc method to update the internal IP address for
// cross-VPC access.
type CrossVpcUpdateOpts struct {
	// User-defined advertised IP contents key-value pair.
	// The key is the listeners IP.
	// The value is advertised.listeners IP, or domain name.
	Contents map[string]string `json:"advertised_ip_contents" required:"true"`
}

// UpdateCrossVpc is a method to update the internal IP address for cross-VPC access using given parameters.
func UpdateCrossVpc(c *golangsdk.ServiceClient, instanceId string, opts CrossVpcUpdateOpts) (*CrossVpc, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r CrossVpc
	_, err = c.Post(crossVpcURL(c, instanceId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// AutoTopicOpts is a struct which represents the parameter of UpdateAutoTopic function
type AutoTopicOpts struct {
	// Indicates whether to enable automatic topic creation.
	EnableAutoTopic *bool `json:"enable_auto_topic" required:"true"`
}

// UpdateAutoTopic is used to enable or disable automatic topic creation.
// via accessing to the service with POST method and parameters
func UpdateAutoTopic(client *golangsdk.ServiceClient, id string,
	opts AutoTopicOpts) (r AutoTopicResult) {
	body, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(autoTopicURL(client, id), body, nil, &golangsdk.RequestOpts{})
	return
}

// ResetPasswordOpts is a struct which represents the parameter of ResetPassword function
type ResetPasswordOpts struct {
	// Indicates the new password of an instance.
	NewPassword string `json:"new_password" required:"true"`
}

// ToResetPasswordMap is used for type convert
func (opts ResetPasswordOpts) ToResetPasswordMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// ResetPasswordOptsBuilder is an interface which can build the map parameter of ResetPassword function
type ResetPasswordOptsBuilder interface {
	ToResetPasswordMap() (map[string]interface{}, error)
}

// ResetPassword is used to reset password for the instance
// via accessing to the service with POST method and parameters
func ResetPassword(client *golangsdk.ServiceClient, id string, opts ResetPasswordOptsBuilder) (r ResetPasswordResult) {
	body, err := opts.ToResetPasswordMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(resetPasswordURL(client, id), body, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

type ConfigParam struct {
	Name  string `json:"name" required:"true"`
	Value string `json:"value" required:"true"`
}

type KafkaConfigs struct {
	KafkaConfigs []ConfigParam `json:"kafka_configs"`
}

type RestartInstanceOpts struct {
	Action    string   `json:"action" required:"true"`
	Instances []string `json:"instances" required:"true"`
}

func ModifyConfiguration(c *golangsdk.ServiceClient, instanceID string, opts KafkaConfigs) (r ModifyConfigurationResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Put(configurationsURL(c, instanceID), b, &r.Body, &golangsdk.RequestOpts{})
	return
}

func GetConfigurations(c *golangsdk.ServiceClient, instanceID string) (r GetConfigurationResult) {
	_, r.Err = c.Get(configurationsURL(c, instanceID), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	})
	return
}

func RebootInstance(c *golangsdk.ServiceClient, params RestartInstanceOpts) (r RebootResult) {
	_, r.Err = c.Post(actionURL(c), params, &r.Body, &golangsdk.RequestOpts{})
	return
}

func GetTasks(c *golangsdk.ServiceClient, instanceID string) (r GetTasksResult) {
	_, r.Err = c.Get(tasksURL(c, instanceID), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	})
	return
}

func GetTask(c *golangsdk.ServiceClient, instanceID, taskID string) (r GetTasksResult) {
	_, r.Err = c.Get(taskURL(c, instanceID, taskID), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	})
	return
}
