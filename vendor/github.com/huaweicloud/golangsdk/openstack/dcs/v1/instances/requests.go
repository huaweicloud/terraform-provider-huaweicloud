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
	// DCS instance name.
	// An instance name is a string of 4–64 characters
	// that contain letters, digits, underscores (_), and hyphens (-).
	// An instance name must start with letters.
	Name string `json:"name" required:"true"`

	// Brief description of the DCS instance.
	// A brief description supports up to 1024 characters.
	Description string `json:"description,omitempty"`

	// Cache engine, which is Redis.
	Engine string `json:"engine" required:"true"`

	// Cache engine version, which is 3.0.7.
	EngineVersion string `json:"engine_version"`

	// Indicates the message storage space.
	// Cache capacity.

	// Unit: GB.
	// For a DCS Redis instance in single-node or master/standby mode,
	// the cache capacity can be 2 GB, 4 GB, 8 GB, 16 GB, 32 GB, or 64 GB.
	// For a DCS Redis instance in cluster mode, the cache capacity can be
	// 64, 128, 256, 512, or 1024 GB.
	Capacity int `json:"capacity" required:"true"`

	// Indicate if no password visit cache instance is allowed.
	NoPasswordAccess string `json:"no_password_access,omitempty"`

	// Indicates the password of an instance.
	// An instance password must meet the following complexity requirements:

	// Password of a DCS instance.
	// The password of a DCS Redis instance must meet
	// the following complexity requirements:
	// A string of 6–32 characters.
	// Contains at least two of the following character types:
	// Uppercase letters
	// Lowercase letters
	// Digits
	// Special characters, such as `~!@#$%^&*()-_=+\|[{}]:'",<.>/?
	Password string `json:"password,omitempty"`

	// When NoPasswordAccess is flase, the AccessUser is enabled.
	AccessUser string `json:"access_user,omitempty"`

	// Tenant's VPC ID.
	VPCID string `json:"vpc_id" required:"true"`

	// Tenant's security group ID.
	SecurityGroupID string `json:"security_group_id,omitempty"`

	// Subnet ID.
	SubnetID string `json:"subnet_id" required:"true"`

	// IDs of the AZs where cache nodes reside.
	// In the current version, only one AZ ID can be set in the request.
	AvailableZones []string `json:"available_zones" required:"true"`

	// Product ID used to differentiate DCS instance types.
	ProductID string `json:"product_id" required:"true"`

	// Backup policy.
	// This parameter is available for master/standby DCS instances.
	InstanceBackupPolicy *InstanceBackupPolicy `json:"instance_backup_policy,omitempty"`

	// Indicates the time at which a maintenance time window starts.
	// Format: HH:mm:ss
	MaintainBegin string `json:"maintain_begin,omitempty"`

	// Indicates the time at which a maintenance time window ends.
	// Format: HH:mm:ss
	MaintainEnd string `json:"maintain_end,omitempty"`
}

// InstanceBackupPolicy for dcs
type InstanceBackupPolicy struct {
	// Retention time.
	// Unit: day.
	// Range: 1–7.
	SaveDays int `json:"save_days" required:"true"`

	// Backup type. Options:
	// auto: automatic backup.
	// manual: manual backup.
	BackupType string `json:"backup_type" required:"true"`

	// Backup plan.
	PeriodicalBackupPlan PeriodicalBackupPlan `json:"periodical_backup_plan" required:"true"`
}

// PeriodicalBackupPlan for dcs
type PeriodicalBackupPlan struct {
	// Time at which backup starts.
	// "00:00-01:00" indicates that backup starts at 00:00:00.
	BeginAt string `json:"begin_at" required:"true"`

	// Interval at which backup is performed.
	// Currently, only weekly backup is supported.
	PeriodType string `json:"period_type" required:"true"`

	// Day in a week on which backup starts.
	// Range: 1–7. Where: 1 indicates Monday; 7 indicates Sunday.
	BackupAt []int `json:"backup_at" required:"true"`
}

type ListDcsInstanceOpts struct {
	Id            string `q:"id"`
	Name          string `q:"name"`
	Type          string `q:"type"`
	DataStoreType string `q:"datastore_type"`
	VpcId         string `q:"vpc_id"`
	SubnetId      string `q:"subnet_id"`
	Offset        int    `q:"offset"`
	Limit         int    `q:"limit"`
}

type ListDcsBuilder interface {
	ToDcsListDetailQuery() (string, error)
}

func (opts ListDcsInstanceOpts) ToDcsListDetailQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
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
	// DCS instance name.
	// An instance name is a string of 4–64 characters
	// that contain letters, digits, underscores (_), and hyphens (-).
	// An instance name must start with letters.
	Name string `json:"name,omitempty"`

	// Brief description of the DCS instance.
	// A brief description supports up to 1024 characters.
	Description *string `json:"description,omitempty"`

	// Backup policy.
	// This parameter is available for master/standby DCS instances.
	InstanceBackupPolicy *InstanceBackupPolicy `json:"instance_backup_policy,omitempty"`

	// Time at which the maintenance time window starts.
	// Format: HH:mm:ss
	MaintainBegin string `json:"maintain_begin,omitempty"`

	// Time at which the maintenance time window ends.
	// Format: HH:mm:ss
	MaintainEnd string `json:"maintain_end,omitempty"`

	// Security group ID.
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

//UpdatePasswordOptsBuilder is an interface which can build the map paramter of update password function
type UpdatePasswordOptsBuilder interface {
	ToPasswordUpdateMap() (map[string]interface{}, error)
}

//UpdatePasswordOpts is a struct which represents the parameters of update function
type UpdatePasswordOpts struct {
	// Old password. It may be empty.
	OldPassword string `json:"old_password" required:"true"`

	// New password.
	// Password complexity requirements:
	// A string of 6–32 characters.
	// Must be different from the old password.
	// Contains at least two types of the following characters:
	// Uppercase letters
	// Lowercase letters
	// Digits
	// Special characters `~!@#$%^&*()-_=+\|[{}]:'",<.>/?
	NewPassword string `json:"new_password" required:"true"`
}

// ToPasswordUpdateMap is used for type convert
func (opts UpdatePasswordOpts) ToPasswordUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// UpdatePassword is updating password for a dcs instance
func UpdatePassword(client *golangsdk.ServiceClient, id string, opts UpdatePasswordOptsBuilder) (r UpdatePasswordResult) {

	body, err := opts.ToPasswordUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(passwordURL(client, id), body, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

//ExtendOptsBuilder is an interface which can build the map paramter of extend function
type ExtendOptsBuilder interface {
	ToExtendMap() (map[string]interface{}, error)
}

//ExtendOpts is a struct which represents the parameters of extend function
type ExtendOpts struct {
	// New specifications (memory space) of the DCS instance.
	// The new specification value to which the DCS instance
	// will be scaled up must be greater than the current specification value.
	// Unit: GB.
	NewCapacity int `json:"new_capacity" required:"true"`

	// New order ID.
	OrderID string `json:"order_id,omitempty"`
}

// ToExtendMap is used for type convert
func (opts ExtendOpts) ToExtendMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Extend is extending for a dcs instance
func Extend(client *golangsdk.ServiceClient, id string, opts ExtendOptsBuilder) (r ExtendResult) {

	body, err := opts.ToExtendMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(extendURL(client, id), body, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

func List(client *golangsdk.ServiceClient, opts ListDcsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToDcsListDetailQuery()

		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageDcsList := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return DcsPage{pagination.SinglePageBase(r)}
	})

	dcsheader := map[string]string{"Content-Type": "application/json"}
	pageDcsList.Headers = dcsheader
	return pageDcsList
}
