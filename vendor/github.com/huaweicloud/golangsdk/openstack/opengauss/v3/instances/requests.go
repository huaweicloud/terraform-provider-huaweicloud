package instances

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type DataStoreOpt struct {
	Type    string `json:"type" required:"true"`
	Version string `json:"version" required:"true"`
}

type BackupStrategyOpt struct {
	StartTime string `json:"start_time" required:"true"`
	KeepDays  int    `json:"keep_days,omitempty"`
}

type HaOpt struct {
	Mode            string `json:"mode" required:"true"`
	ReplicationMode string `json:"replication_mode" required:"true"`
	Consistency     string `json:"consistency,omitempty"`
}

type VolumeOpt struct {
	Type string `json:"type" required:"true"`
	Size int    `json:"size" required:"true"`
}

type RestorePointOpt struct {
	InstanceId   string            `json:"instance_id" required:"true"`
	Type         string            `json:"type" required:"true"`
	BackupId     string            `json:"backup_id,omitempty"`
	RestoreTime  int               `json:"restore_time,omitempty"`
	DatabaseName map[string]string `json:"database_name,omitempty"`
}

type CreateGaussDBOpts struct {
	Name                string             `json:"name" required:"true"`
	Region              string             `json:"region,omitempty"`
	Flavor              string             `json:"flavor_ref" required:"true"`
	VpcId               string             `json:"vpc_id,omitempty"`
	SubnetId            string             `json:"subnet_id,omitempty"`
	SecurityGroupId     string             `json:"security_group_id,omitempty"`
	Password            string             `json:"password" required:"true"`
	Port                string             `json:"port,omitempty"`
	DiskEncryptionId    string             `json:"disk_encryption_id,omitempty"`
	TimeZone            string             `json:"time_zone,omitempty"`
	AvailabilityZone    string             `json:"availability_zone" required:"true"`
	ConfigurationId     string             `json:"configuration_id,omitempty"`
	DsspoolId           string             `json:"dsspool_id,omitempty"`
	ReplicaOfId         string             `json:"replica_of_id,omitempty"`
	ShardingNum         int                `json:"sharding_num,omitempty"`
	CoordinatorNum      int                `json:"coordinator_num,omitempty"`
	EnterpriseProjectId string             `json:"enterprise_project_id,omitempty"`
	DataStore           DataStoreOpt       `json:"datastore" required:"true"`
	Volume              VolumeOpt          `json:"volume" required:"true"`
	Ha                  *HaOpt             `json:"ha,omitempty"`
	BackupStrategy      *BackupStrategyOpt `json:"backup_strategy,omitempty"`
	RestorePoint        *RestorePointOpt   `json:"restore_point,omitempty"`
}

type CreateGaussDBBuilder interface {
	ToInstancesCreateMap() (map[string]interface{}, error)
}

func (opts CreateGaussDBOpts) ToInstancesCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Create(client *golangsdk.ServiceClient, opts CreateGaussDBBuilder) (r CreateResult) {
	b, err := opts.ToInstancesCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{202},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})
	return
}

type UpdateVolumeOptsBuilder interface {
	ToVolumeUpdateMap() (map[string]interface{}, error)
}

type UpdateVolumeOpts struct {
	Size int `json:"size" required:"true"`
}

func (opts UpdateVolumeOpts) ToVolumeUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "enlarge_volume")
}

func UpdateVolume(client *golangsdk.ServiceClient, opts UpdateVolumeOptsBuilder, id string) (r UpdateResult) {
	b, err := opts.ToVolumeUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(updateURL(client, id, "action"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{202},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})
	return
}

type UpdateClusterOptsBuilder interface {
	ToClusterUpdateMap() (map[string]interface{}, error)
}

type Shard struct {
	Count int `json:"count" required:"true"`
}

type Coordinator struct {
	AzCode string `json:"az_code" required:"true"`
}

type UpdateClusterOpts struct {
	Shard        *Shard        `json:"shard,omitempty"`
	Coordinators []Coordinator `json:"coordinators,omitempty"`
}

func (opts UpdateClusterOpts) ToClusterUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "expand_cluster")
}

func UpdateCluster(client *golangsdk.ServiceClient, opts UpdateClusterOptsBuilder, id string) (r UpdateResult) {
	b, err := opts.ToClusterUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(updateURL(client, id, "action"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{202},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})
	return
}

func Delete(client *golangsdk.ServiceClient, instanceId string) (r DeleteResult) {
	url := deleteURL(client, instanceId)

	_, r.Err = client.Delete(url, &golangsdk.RequestOpts{
		OkCodes:     []int{202},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})

	return
}

type ListGaussDBInstanceOpts struct {
	Id            string `q:"id"`
	Name          string `q:"name"`
	Type          string `q:"type"`
	DataStoreType string `q:"datastore_type"`
	VpcId         string `q:"vpc_id"`
	SubnetId      string `q:"subnet_id"`
	Offset        int    `q:"offset"`
	Limit         int    `q:"limit"`
}

type ListGaussDBBuilder interface {
	ToGaussDBListDetailQuery() (string, error)
}

func (opts ListGaussDBInstanceOpts) ToGaussDBListDetailQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, opts ListGaussDBBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToGaussDBListDetailQuery()

		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageList := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return GaussDBPage{pagination.SinglePageBase(r)}
	})
	// Headers supplies additional HTTP headers to populate on each paged request
	pageList.Headers = map[string]string{"Content-Type": "application/json"}

	return pageList
}

func GetInstanceByID(client *golangsdk.ServiceClient, instanceId string) (GaussDBInstance, error) {
	var instance GaussDBInstance

	opts := ListGaussDBInstanceOpts{
		Id: instanceId,
	}

	pages, err := List(client, &opts).AllPages()
	if err != nil {
		return instance, err
	}

	all, err := ExtractGaussDBInstances(pages)
	if err != nil {
		return instance, err
	}
	if all.TotalCount == 0 {
		return instance, nil
	}

	instance = all.Instances[0]
	return instance, nil
}

func GetInstanceByName(client *golangsdk.ServiceClient, name string) (GaussDBInstance, error) {
	var instance GaussDBInstance

	opts := ListGaussDBInstanceOpts{
		Name: name,
	}

	pages, err := List(client, &opts).AllPages()
	if err != nil {
		return instance, err
	}

	all, err := ExtractGaussDBInstances(pages)
	if err != nil {
		return instance, err
	}
	if all.TotalCount == 0 {
		return instance, nil
	}

	instance = all.Instances[0]
	return instance, nil
}

type RenameOptsBuilder interface {
	ToRenameMap() (map[string]interface{}, error)
}

type RenameOpts struct {
	Name string `json:"name" required:"true"`
}

func (opts RenameOpts) ToRenameMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func Rename(client *golangsdk.ServiceClient, opts RenameOptsBuilder, id string) (r RenameResult) {
	b, err := opts.ToRenameMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updateURL(client, id, "name"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})
	return
}

type RestorePasswordOptsBuilder interface {
	ToRestorePasswordMap() (map[string]interface{}, error)
}

type RestorePasswordOpts struct {
	Password string `json:"password" required:"true"`
}

func (opts RestorePasswordOpts) ToRestorePasswordMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func RestorePassword(client *golangsdk.ServiceClient, opts RestorePasswordOptsBuilder, id string) (r golangsdk.Result) {
	b, err := opts.ToRestorePasswordMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(updateURL(client, id, "password"), b, nil, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})
	return
}
