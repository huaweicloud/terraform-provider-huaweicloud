package instances

import (
	"fmt"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type AvailabilityZoneDetailOpt struct {
	PrimaryAvailabilityZone   string `json:"primary_availability_zone" required:"true"`
	SecondaryAvailabilityZone string `json:"secondary_availability_zone" required:"true"`
}

type FlavorOpt struct {
	Num      string `json:"num" required:"true"`
	Size     int    `json:"size" required:"true"`
	Storage  string `json:"storage" required:"true"`
	SpecCode string `json:"spec_code" required:"true"`
}

type BackupStrategyOpt struct {
	StartTime string `json:"start_time,omitempty"`
	KeepDays  string `json:"keep_days,omitempty"`
}

type ChargeInfoOpt struct {
	ChargingMode string `json:"charge_mode,omitempty"`
	PeriodType   string `json:"period_type,omitempty"`
	PeriodNum    int    `json:"period_num,omitempty"`
	IsAutoRenew  string `json:"is_auto_renew,omitempty"`
	IsAutoPay    string `json:"is_auto_pay,omitempty"`
}

type CreateGeminiDBOpts struct {
	Name                   string                     `json:"name" required:"true"`
	Region                 string                     `json:"region" required:"true"`
	AvailabilityZone       string                     `json:"availability_zone" required:"true"`
	AvailabilityZoneDetail *AvailabilityZoneDetailOpt `json:"availability_zone_detail,omitempty"`
	VpcId                  string                     `json:"vpc_id" required:"true"`
	SubnetId               string                     `json:"subnet_id" required:"true"`
	SecurityGroupId        string                     `json:"security_group_id,omitempty"`
	Password               string                     `json:"password" required:"true"`
	Mode                   string                     `json:"mode" required:"true"`
	ConfigurationId        string                     `json:"configuration_id,omitempty"`
	EnterpriseProjectId    string                     `json:"enterprise_project_id,omitempty"`
	DedicatedResourceId    string                     `json:"dedicated_resource_id,omitempty"`
	Ssl                    string                     `json:"ssl_option,omitempty"`
	Port                   string                     `json:"port,omitempty"`
	DataStore              DataStore                  `json:"datastore" required:"true"`
	Flavor                 []FlavorOpt                `json:"flavor" required:"true"`
	BackupStrategy         *BackupStrategyOpt         `json:"backup_strategy,omitempty"`
	ChargeInfo             *ChargeInfoOpt             `json:"charge_info,omitempty"`
}

type CreateGeminiDBBuilder interface {
	ToInstancesCreateMap() (map[string]interface{}, error)
}

func (opts CreateGeminiDBOpts) ToInstancesCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Create(client *golangsdk.ServiceClient, opts CreateGeminiDBBuilder) (r CreateResult) {
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

type ExtendVolumeOpts struct {
	Size      int    `json:"size" required:"true"`
	IsAutoPay string `json:"is_auto_pay,omitempty"`
}

type ExtendVolumeBuilder interface {
	ToVolumeExtendMap() (map[string]interface{}, error)
}

func (opts ExtendVolumeOpts) ToVolumeExtendMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ExtendVolume(client *golangsdk.ServiceClient, instanceId string, opts ExtendVolumeBuilder) (r ExtendResult) {
	b, err := opts.ToVolumeExtendMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(extendURL(client, instanceId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{202},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})
	return
}

type EnlargeNodeOpts struct {
	Num       int    `json:"num" required:"true"`
	IsAutoPay string `json:"is_auto_pay,omitempty"`
}

type EnlargeNodeBuilder interface {
	ToNodeEnlargeMap() (map[string]interface{}, error)
}

func (opts EnlargeNodeOpts) ToNodeEnlargeMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func EnlargeNode(client *golangsdk.ServiceClient, instanceId string, opts EnlargeNodeBuilder) (r ExtendResult) {
	b, err := opts.ToNodeEnlargeMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(enlargeNodeURL(client, instanceId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{202},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})
	return
}

type ReduceNodeOpts struct {
	Num int `json:"num" required:"true"`
}

type ReduceNodeBuilder interface {
	ToNodeReduceMap() (map[string]interface{}, error)
}

func (opts ReduceNodeOpts) ToNodeReduceMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ReduceNode(client *golangsdk.ServiceClient, instanceId string, opts ReduceNodeBuilder) (r ExtendResult) {
	b, err := opts.ToNodeReduceMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(reduceNodeURL(client, instanceId), b, &r.Body, &golangsdk.RequestOpts{
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

type ListGeminiDBInstanceOpts struct {
	Id            string `q:"id"`
	Name          string `q:"name"`
	Mode          string `q:"mode"`
	DataStoreType string `q:"datastore_type"`
	VpcId         string `q:"vpc_id"`
	SubnetId      string `q:"subnet_id"`
	Offset        int    `q:"offset"`
	Limit         int    `q:"limit"`
}

type ListGeminiDBBuilder interface {
	ToGeminiDBListDetailQuery() (string, error)
}

func (opts ListGeminiDBInstanceOpts) ToGeminiDBListDetailQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, opts ListGeminiDBBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToGeminiDBListDetailQuery()

		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageList := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return GeminiDBPage{pagination.SinglePageBase(r)}
	})
	// Headers supplies additional HTTP headers to populate on each paged request
	pageList.Headers = map[string]string{"Content-Type": "application/json"}

	return pageList
}

func GetInstanceByID(client *golangsdk.ServiceClient, instanceId string) (GeminiDBInstance, error) {
	var instance GeminiDBInstance

	opts := ListGeminiDBInstanceOpts{
		Id: instanceId,
	}

	pages, err := List(client, &opts).AllPages()
	if err != nil {
		return instance, err
	}

	all, err := ExtractGeminiDBInstances(pages)
	if err != nil {
		return instance, err
	}
	if all.TotalCount < 1 {
		return instance, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte(fmt.Sprintf("the database instance (%s) does not exist", instanceId)),
			},
		}
	}

	instance = all.Instances[0]
	return instance, nil
}

func GetInstanceByName(client *golangsdk.ServiceClient, name string) (GeminiDBInstance, error) {
	var instance GeminiDBInstance

	opts := ListGeminiDBInstanceOpts{
		Name: name,
	}

	pages, err := List(client, &opts).AllPages()
	if err != nil {
		return instance, err
	}

	all, err := ExtractGeminiDBInstances(pages)
	if err != nil {
		return instance, err
	}
	if all.TotalCount < 1 {
		return instance, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte(fmt.Sprintf("the database instance (%s) does not exist", name)),
			},
		}
	}

	instance = all.Instances[0]
	return instance, nil
}

type UpdateNameOpts struct {
	Name string `json:"name" required:"true"`
}

type UpdateNameBuilder interface {
	ToNameUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateNameOpts) ToNameUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UpdateName(client *golangsdk.ServiceClient, instanceId string, opts UpdateNameBuilder) (r UpdateResult) {
	b, err := opts.ToNameUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updateNameURL(client, instanceId), b, nil, &golangsdk.RequestOpts{
		OkCodes:     []int{204},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})
	return
}

type UpdatePassOpts struct {
	Password string `json:"password" required:"true"`
}

type UpdatePassBuilder interface {
	ToPassUpdateMap() (map[string]interface{}, error)
}

func (opts UpdatePassOpts) ToPassUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UpdatePass(client *golangsdk.ServiceClient, instanceId string, opts UpdatePassBuilder) (r UpdateResult) {
	b, err := opts.ToPassUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updatePassURL(client, instanceId), b, nil, &golangsdk.RequestOpts{
		OkCodes:     []int{204},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})
	return
}

type ResizeOpt struct {
	InstanceID string `json:"target_id" required:"true"`
	SpecCode   string `json:"target_spec_code" required:"true"`
}

type ResizeOpts struct {
	Resize    ResizeOpt `json:"resize" required:"true"`
	IsAutoPay string    `json:"is_auto_pay,omitempty"`
}

type ResizeBuilder interface {
	ToResizeMap() (map[string]interface{}, error)
}

func (opts ResizeOpts) ToResizeMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Resize(client *golangsdk.ServiceClient, instanceId string, opts ResizeBuilder) (r ExtendResult) {
	b, err := opts.ToResizeMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(resizeURL(client, instanceId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{202},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})
	return
}

type UpdateSgOpts struct {
	SecurityGroupID string `json:"security_group_id" required:"true"`
}

type UpdateSgBuilder interface {
	ToSgUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateSgOpts) ToSgUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UpdateSg(client *golangsdk.ServiceClient, instanceId string, opts UpdateSgBuilder) (r ExtendResult) {
	b, err := opts.ToSgUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updateSgURL(client, instanceId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{202},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})
	return
}

func ListDeh(client *golangsdk.ServiceClient) pagination.Pager {
	pageList := pagination.NewPager(client, listDehURL(client), func(r pagination.PageResult) pagination.Page {
		return DehResourcePage{pagination.SinglePageBase(r)}
	})
	// Headers supplies additional HTTP headers to populate on each paged request
	pageList.Headers = map[string]string{"Content-Type": "application/json"}

	return pageList
}

type UpdateSslOpts struct {
	Ssl string `json:"ssl_option" required:"true"`
}

type UpdateSslBuilder interface {
	ToSslUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateSslOpts) ToSslUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UpdateSsl(client *golangsdk.ServiceClient, instanceId string, opts UpdateSslBuilder) (r UpdateResult) {
	b, err := opts.ToSslUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(updateSslURL(client, instanceId), b, nil, &golangsdk.RequestOpts{
		OkCodes:     []int{202},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})
	return
}

type UpdatePublicIpOpts struct {
	Action     string `json:"action" required:"true"`
	PublicIp   string `json:"public_ip,omitempty"`
	PublicIpId string `json:"public_ip_id,omitempty"`
}

type UpdatePublicIpBuilder interface {
	ToPublicIpUpdateMap() (map[string]interface{}, error)
}

func (opts UpdatePublicIpOpts) ToPublicIpUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UpdatePublicIp(client *golangsdk.ServiceClient, instanceId, nodeId string, opts UpdatePublicIpOpts) (r PublicIpResult) {
	b, err := opts.ToPublicIpUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(updatePublicIpURL(client, instanceId, nodeId), b, nil, &golangsdk.RequestOpts{
		OkCodes:     []int{202},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})
	return
}
