package instances

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

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
	// only letters, digits, and hyphens (-).
	Name string `json:"name" required:"true"`

	// Indicates the description of an instance.
	// It is a character string containing not more than 1024 characters.
	Description string `json:"description,omitempty"`

	// Indicates a message engine.
	// Currently, only RabbitMQ is supported.
	Engine string `json:"engine" required:"true"`

	// Indicates the version of a message engine.
	EngineVersion string `json:"engine_version,omitempty"`

	// Indicates the message storage space.
	StorageSpace int `json:"storage_space" required:"true"`

	// Indicates the password of an instance.
	// An instance password must meet the following complexity requirements:
	// Must be 6 to 32 characters long.
	// Must contain at least two of the following character types:
	// Lowercase letters
	// Uppercase letters
	// Digits
	// Special characters (`~!@#$%^&*()-_=+\|[{}]:'",<.>/?)
	Password string `json:"password,omitempty"`

	// Indicates a username.
	// A username consists of 1 to 64 characters
	// and supports only letters, digits, and hyphens (-).
	AccessUser string `json:"access_user,omitempty"`

	// Indicates the ID of a VPC.
	VPCID string `json:"vpc_id" required:"true"`

	// Indicates the ID of a security group.
	SecurityGroupID string `json:"security_group_id" required:"true"`

	// Indicates the ID of a subnet.
	SubnetID string `json:"subnet_id" required:"true"`

	// Indicates the ID of an AZ.
	// The parameter value can be left blank or an empty array.
	AvailableZones []string `json:"available_zones,omitempty"`

	// Indicates a product ID.
	ProductID string `json:"product_id" required:"true"`

	// Indicates the time at which a maintenance time window starts.
	// Format: HH:mm:ss
	MaintainBegin string `json:"maintain_begin,omitempty"`

	// Indicates the time at which a maintenance time window ends.
	// Format: HH:mm:ss
	MaintainEnd string `json:"maintain_end,omitempty"`

	//This parameter is mandatory when a Kafka instance is created.
	//Indicates the maximum number of topics in a Kafka instance.
	PartitionNum int `json:"partition_num,omitempty"`

	//Indicates whether to enable SSL-encrypted access.
	SslEnable bool `json:"ssl_enable"`

	//This parameter is mandatory if the engine is kafka.
	//Indicates the baseline bandwidth of a Kafka instance, that is,
	//the maximum amount of data transferred per unit time. Unit: byte/s.
	Specification string `json:"specification,omitempty"`

	//Indicates the storage I/O specification. For details on how to select a disk type
	StorageSpecCode string `json:"storage_spec_code,omitempty"`
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

// Delete an instance by id
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

//UpdateOptsBuilder is an interface which can build the map paramter of update function
type UpdateOptsBuilder interface {
	ToInstanceUpdateMap() (map[string]interface{}, error)
}

//UpdateOpts is a struct which represents the parameters of update function
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

type ListDmsInstanceOpts struct {
	Id                  string `q:"id"`
	Name                string `q:"name"`
	Engine              string `q:"engine"`
	Status              string `q:"status"`
	IncludeFailure      string `q:"includeFailure"`
	ExactMatchName      string `q:"exactMatchName"`
	EnterpriseProjectID int    `q:"enterprise_project_id"`
}

type ListDmsBuilder interface {
	ToDmsListDetailQuery() (string, error)
}

func (opts ListDmsInstanceOpts) ToDmsListDetailQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, opts ListDmsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToDmsListDetailQuery()

		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageDmsList := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return DmsPage{pagination.SinglePageBase(r)}
	})

	dmsheader := map[string]string{"Content-Type": "application/json"}
	pageDmsList.Headers = dmsheader
	return pageDmsList
}
