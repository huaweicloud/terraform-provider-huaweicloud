package instances

import (
	"fmt"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type CreateRdsOpts struct {
	Name                string          `json:"name"  required:"true"`
	Datastore           *Datastore      `json:"datastore" required:"true"`
	Ha                  *Ha             `json:"ha,omitempty"`
	ConfigurationId     string          `json:"configuration_id,omitempty"`
	Port                string          `json:"port,omitempty"`
	Password            string          `json:"password" required:"true"`
	BackupStrategy      *BackupStrategy `json:"backup_strategy,omitempty"`
	EnterpriseProjectId string          `json:"enterprise_project_id,omitempty"`
	DiskEncryptionId    string          `json:"disk_encryption_id,omitempty"`
	FlavorRef           string          `json:"flavor_ref" required:"true"`
	Volume              *Volume         `json:"volume" required:"true"`
	Region              string          `json:"region" required:"true"`
	AvailabilityZone    string          `json:"availability_zone" required:"true"`
	VpcId               string          `json:"vpc_id" required:"true"`
	SubnetId            string          `json:"subnet_id" required:"true"`
	SecurityGroupId     string          `json:"security_group_id" required:"true"`
	ChargeInfo          *ChargeInfo     `json:"charge_info,omitempty"`
	TimeZone            string          `json:"time_zone,omitempty"`
}

type CreateReplicaOpts struct {
	Name                string      `json:"name"  required:"true"`
	ReplicaOfId         string      `json:"replica_of_id" required:"true"`
	EnterpriseProjectId string      `json:"enterprise_project_id,omitempty"`
	DiskEncryptionId    string      `json:"disk_encryption_id,omitempty"`
	FlavorRef           string      `json:"flavor_ref" required:"true"`
	Volume              *Volume     `json:"volume" required:"true"`
	Region              string      `json:"region,omitempty"`
	AvailabilityZone    string      `json:"availability_zone" required:"true"`
	ChargeInfo          *ChargeInfo `json:"charge_info,omitempty"`
}

type Datastore struct {
	Type    string `json:"type" required:"true"`
	Version string `json:"version" required:"true"`
}

type Ha struct {
	Mode            string `json:"mode" required:"true"`
	ReplicationMode string `json:"replication_mode,omitempty"`
}

type BackupStrategy struct {
	StartTime string `json:"start_time" required:"true"`
	KeepDays  int    `json:"keep_days,omitempty"`
}

type Volume struct {
	Type string `json:"type" required:"true"`
	Size int    `json:"size" required:"true"`
}

type ChargeInfo struct {
	ChargeMode  string `json:"charge_mode" required:"true"`
	PeriodType  string `json:"period_type,omitempty"`
	PeriodNum   int    `json:"period_num,omitempty"`
	IsAutoRenew string `json:"is_auto_renew,omitempty"`
	IsAutoPay   string `json:"is_auto_pay,omitempty"`
}

type CreateRdsBuilder interface {
	ToInstancesCreateMap() (map[string]interface{}, error)
}

func (opts CreateRdsOpts) ToInstancesCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Create(client *golangsdk.ServiceClient, opts CreateRdsBuilder) (r CreateResult) {
	b, err := opts.ToInstancesCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})
	return
}

type CreateReplicaBuilder interface {
	ToCreateReplicaMap() (map[string]interface{}, error)
}

func (opts CreateReplicaOpts) ToCreateReplicaMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func CreateReplica(client *golangsdk.ServiceClient, opts CreateReplicaBuilder) (r CreateResult) {
	b, err := opts.ToCreateReplicaMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})

	return
}

type DeleteInstance struct {
	InstanceId string `json:"instance_id" required:"true"`
}

type DeleteInstanceBuilder interface {
	ToInstancesDeleteMap() (map[string]interface{}, error)
}

func (opts DeleteInstance) ToInstancesDeleteMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(&opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Delete(client *golangsdk.ServiceClient, instanceId string) (r DeleteInstanceRdsResult) {

	url := deleteURL(client, instanceId)

	_, r.Err = client.Delete(url, &golangsdk.RequestOpts{JSONResponse: &r.Body, MoreHeaders: map[string]string{"Content-Type": "application/json"}})
	return
}

type RestartRdsInstanceOpts struct {
	Restart string `json:"restart" required:"true"`
}

type RestartRdsInstanceBuilder interface {
	ToRestartRdsInstanceMap() (map[string]interface{}, error)
}

func (opts RestartRdsInstanceOpts) ToRestartRdsInstanceMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(&opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Restart(client *golangsdk.ServiceClient, opts RestartRdsInstanceBuilder, instanceId string) (r RestartRdsInstanceResult) {
	b, err := opts.ToRestartRdsInstanceMap()
	if err != nil {
		r.Err = err
		return
	}
	fmt.Println("restart Rds instance body = ", b)
	_, r.Err = client.Post(updateURL(client, instanceId, "action"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})
	return
}

type RenameRdsInstanceOpts struct {
	Name string `json:"name" required:"true"`
}

type RenameRdsInstanceBuilder interface {
	ToRenameRdsInstanceMap() (map[string]interface{}, error)
}

func (opts RenameRdsInstanceOpts) ToRenameRdsInstanceMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(&opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Rename(client *golangsdk.ServiceClient, opts RenameRdsInstanceBuilder, instanceId string) (r golangsdk.Result) {
	b, err := opts.ToRenameRdsInstanceMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(updateURL(client, instanceId, "name"), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}

type ListRdsInstanceOpts struct {
	Id            string `q:"id"`
	Name          string `q:"name"`
	Type          string `q:"type"`
	DataStoreType string `q:"datastore_type"`
	VpcId         string `q:"vpc_id"`
	SubnetId      string `q:"subnet_id"`
	Offset        int    `q:"offset"`
	Limit         int    `q:"limit"`
}

type ListRdsBuilder interface {
	ToRdsListDetailQuery() (string, error)
}

func (opts ListRdsInstanceOpts) ToRdsListDetailQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, opts ListRdsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToRdsListDetailQuery()

		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageRdsList := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return RdsPage{pagination.SinglePageBase(r)}
	})

	rdsheader := map[string]string{"Content-Type": "application/json"}
	pageRdsList.Headers = rdsheader
	return pageRdsList
}

type SingleToHaRdsOpts struct {
	SingleToHa *SingleToHaRds `json:"single_to_ha" required:"true"`
}

type SingleToHaRds struct {
	AzCodeNewNode string `json:"az_code_new_node" required:"true"`
	Password      string `json:"password,omitempty"`
}

type SingleToRdsHaBuilder interface {
	ToSingleToRdsHaMap() (map[string]interface{}, error)
}

func (opts SingleToHaRdsOpts) ToSingleToRdsHaMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(&opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func SingleToHa(client *golangsdk.ServiceClient, opts SingleToRdsHaBuilder, instanceId string) (r SingleToHaRdsInstanceResult) {
	b, err := opts.ToSingleToRdsHaMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(updateURL(client, instanceId, "action"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})

	return
}

type SpecCode struct {
	Speccode string `json:"spec_code" required:"true"`
}

type ResizeFlavorOpts struct {
	ResizeFlavor *SpecCode `json:"resize_flavor" required:"true"`
}

type ResizeFlavorBuilder interface {
	ResizeFlavorMap() (map[string]interface{}, error)
}

func (opts ResizeFlavorOpts) ResizeFlavorMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(&opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Resize(client *golangsdk.ServiceClient, opts ResizeFlavorBuilder, instanceId string) (r ResizeFlavorResult) {
	b, err := opts.ResizeFlavorMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(updateURL(client, instanceId, "action"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})

	return
}

type EnlargeVolumeRdsOpts struct {
	EnlargeVolume *EnlargeVolumeSize `json:"enlarge_volume" required:"true"`
}

type EnlargeVolumeSize struct {
	Size int `json:"size" required:"true"`
}

type EnlargeVolumeBuilder interface {
	ToEnlargeVolumeMap() (map[string]interface{}, error)
}

func (opts EnlargeVolumeRdsOpts) ToEnlargeVolumeMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(&opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func EnlargeVolume(client *golangsdk.ServiceClient, opts EnlargeVolumeBuilder, instanceId string) (r EnlargeVolumeResult) {
	b, err := opts.ToEnlargeVolumeMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(updateURL(client, instanceId, "action"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})

	return
}

type DbErrorlogOpts struct {
	StartDate string `q:"start_date"`
	EndDate   string `q:"end_date"`
	Offset    string `q:"offset"`
	Limit     string `q:"limit"`
	Level     string `q:"level"`
}

type DbErrorlogBuilder interface {
	DbErrorlogQuery() (string, error)
}

func (opts DbErrorlogOpts) DbErrorlogQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func ListErrorLog(client *golangsdk.ServiceClient, opts DbErrorlogBuilder, instanceID string) pagination.Pager {
	url := updateURL(client, instanceID, "errorlog")
	if opts != nil {
		query, err := opts.DbErrorlogQuery()

		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageRdsList := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ErrorLogPage{pagination.SinglePageBase(r)}
	})

	rdsheader := map[string]string{"Content-Type": "application/json"}
	pageRdsList.Headers = rdsheader
	return pageRdsList
}

type DbSlowLogOpts struct {
	StartDate string `q:"start_date"`
	EndDate   string `q:"end_date"`
	Offset    string `q:"offset"`
	Limit     string `q:"limit"`
	Level     string `q:"level"`
}

type DbSlowLogBuilder interface {
	ToDbSlowLogListQuery() (string, error)
}

func (opts DbSlowLogOpts) ToDbSlowLogListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func ListSlowLog(client *golangsdk.ServiceClient, opts DbSlowLogBuilder, instanceID string) pagination.Pager {
	url := updateURL(client, instanceID, "slowlog")
	if opts != nil {
		query, err := opts.ToDbSlowLogListQuery()

		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageRdsList := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ErrorLogPage{pagination.SinglePageBase(r)}
	})

	rdsheader := map[string]string{"Content-Type": "application/json"}
	pageRdsList.Headers = rdsheader
	return pageRdsList
}
