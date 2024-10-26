package instances

import (
	"fmt"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

var requestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

type ChargeInfoOpt struct {
	ChargingMode string `json:"charge_mode,omitempty"`
	PeriodType   string `json:"period_type,omitempty"`
	PeriodNum    int    `json:"period_num,omitempty"`
	IsAutoRenew  string `json:"is_auto_renew,omitempty"`
	IsAutoPay    string `json:"is_auto_pay,omitempty"`
}

type VolumeOpt struct {
	Size int `json:"size" required:"true"`
}

type DataStoreOpt struct {
	Type    string `json:"type" required:"true"`
	Version string `json:"version" required:"true"`
}

type BackupStrategyOpt struct {
	StartTime string `json:"start_time" required:"true"`
	KeepDays  string `json:"keep_days" required:"true"`
}

type CreateTaurusDBOpts struct {
	Name                string             `json:"name" required:"true"`
	Region              string             `json:"region" required:"true"`
	Mode                string             `json:"mode" required:"true"`
	Flavor              string             `json:"flavor_ref" required:"true"`
	VpcId               string             `json:"vpc_id" required:"true"`
	SubnetId            string             `json:"subnet_id" required:"true"`
	SecurityGroupId     string             `json:"security_group_id,omitempty"`
	Password            string             `json:"password" required:"true"`
	TimeZone            string             `json:"time_zone" required:"true"`
	AZMode              string             `json:"availability_zone_mode" required:"true"`
	SlaveCount          int                `json:"slave_count" required:"true"`
	MasterAZ            string             `json:"master_availability_zone,omitempty"`
	ConfigurationId     string             `json:"configuration_id,omitempty"`
	EnterpriseProjectId string             `json:"enterprise_project_id,omitempty"`
	DedicatedResourceId string             `json:"dedicated_resource_id,omitempty"`
	LowerCaseTableNames *int               `json:"lower_case_table_names,omitempty"`
	DataStore           DataStoreOpt       `json:"datastore" required:"true"`
	BackupStrategy      *BackupStrategyOpt `json:"backup_strategy,omitempty"`
	ChargeInfo          *ChargeInfoOpt     `json:"charge_info,omitempty"`
	Volume              *VolumeOpt         `json:"volume,omitempty"`
}

type CreateTaurusDBBuilder interface {
	ToInstancesCreateMap() (map[string]interface{}, error)
}

func (opts CreateTaurusDBOpts) ToInstancesCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Create(client *golangsdk.ServiceClient, opts CreateTaurusDBBuilder) (r CreateResult) {
	b, err := opts.ToInstancesCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{201, 202},
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return
}

type CreateReplicaOpts struct {
	Priorities []int  `json:"priorities" required:"true"`
	IsAutoPay  string `json:"is_auto_pay,omitempty"`
}

type CreateReplicaBuilder interface {
	ToReplicaCreateMap() (map[string]interface{}, error)
}

func (opts CreateReplicaOpts) ToReplicaCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func CreateReplica(client *golangsdk.ServiceClient, instanceId string, opts CreateReplicaBuilder) (r JobResult) {
	b, err := opts.ToReplicaCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(updateURL(client, instanceId, "nodes/enlarge"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{201, 202},
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return
}

func DeleteReplica(client *golangsdk.ServiceClient, instanceId, nodeId string) (r JobResult) {
	url := deleteReplicaURL(client, instanceId, nodeId)

	_, r.Err = client.DeleteWithResponse(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200, 202},
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return
}

func Delete(client *golangsdk.ServiceClient, instanceId string) (r DeleteResult) {
	url := deleteURL(client, instanceId)

	_, r.Err = client.Delete(url, &golangsdk.RequestOpts{
		OkCodes:     []int{200, 202},
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return
}

func Get(client *golangsdk.ServiceClient, instanceId string) (r GetResult) {
	url := getURL(client, instanceId)

	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return
}

type ListTaurusDBInstanceOpts struct {
	Id            string `q:"id"`
	Name          string `q:"name"`
	Type          string `q:"type"`
	DataStoreType string `q:"datastore_type"`
	VpcId         string `q:"vpc_id"`
	SubnetId      string `q:"subnet_id"`
	Offset        int    `q:"offset"`
	Limit         int    `q:"limit"`
}

type ListTaurusDBBuilder interface {
	ToTaurusDBListDetailQuery() (string, error)
}

func (opts ListTaurusDBInstanceOpts) ToTaurusDBListDetailQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, opts ListTaurusDBBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToTaurusDBListDetailQuery()

		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageList := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return TaurusDBPage{pagination.SinglePageBase(r)}
	})
	// Headers supplies additional HTTP headers to populate on each paged request
	pageList.Headers = requestOpts.MoreHeaders

	return pageList
}

func GetInstanceByName(client *golangsdk.ServiceClient, name string) (ListTaurusDBInstance, error) {
	var instance ListTaurusDBInstance

	opts := ListTaurusDBInstanceOpts{
		Name: name,
	}

	pages, err := List(client, &opts).AllPages()
	if err != nil {
		return instance, err
	}

	all, err := ExtractTaurusDBInstances(pages)
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

func UpdateName(client *golangsdk.ServiceClient, instanceId string, opts UpdateNameBuilder) (r JobResult) {
	b, err := opts.ToNameUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updateURL(client, instanceId, "name"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: requestOpts.MoreHeaders,
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

func UpdatePass(client *golangsdk.ServiceClient, instanceId string, opts UpdatePassBuilder) (r JobResult) {
	b, err := opts.ToPassUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(updateURL(client, instanceId, "password"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: requestOpts.MoreHeaders,
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

func ExtendVolume(client *golangsdk.ServiceClient, instanceId string, opts ExtendVolumeBuilder) (r JobResult) {
	b, err := opts.ToVolumeExtendMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(updateURL(client, instanceId, "volume/extend"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{201},
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return
}

type ResizeOpt struct {
	Spec string `json:"spec_code" required:"true"`
}

type ResizeOpts struct {
	Resize    ResizeOpt `json:"resize_flavor" required:"true"`
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

func Resize(client *golangsdk.ServiceClient, instanceId string, opts ResizeBuilder) (r JobResult) {
	b, err := opts.ToResizeMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(updateURL(client, instanceId, "action"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return
}

type ProxyOpts struct {
	Flavor  string `json:"flavor_ref" required:"true"`
	NodeNum int    `json:"node_num" required:"true"`
}

type ProxyBuilder interface {
	ToProxyMap() (map[string]interface{}, error)
}

func (opts ProxyOpts) ToProxyMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func EnableProxy(client *golangsdk.ServiceClient, instanceId string, opts ProxyBuilder) (r JobResult) {
	b, err := opts.ToProxyMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(updateURL(client, instanceId, "proxy"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{201},
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return
}

type EnlargeProxyOpts struct {
	NodeNum int `json:"node_num" required:"true"`
}

type EnlargeProxyBuilder interface {
	ToEnlargeProxyMap() (map[string]interface{}, error)
}

func (opts EnlargeProxyOpts) ToEnlargeProxyMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func EnlargeProxy(client *golangsdk.ServiceClient, instanceId string, opts EnlargeProxyBuilder) (r JobResult) {
	b, err := opts.ToEnlargeProxyMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(updateURL(client, instanceId, "proxy/enlarge"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{201},
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return
}

func DeleteProxy(client *golangsdk.ServiceClient, instanceId string) (r JobResult) {
	url := proxyURL(client, instanceId)

	_, r.Err = client.DeleteWithResponse(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200, 202},
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return
}

func GetProxy(client *golangsdk.ServiceClient, instanceId string) (r GetProxyResult) {
	url := proxyURL(client, instanceId)

	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
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

type RestartOpts struct {
	Delay bool `json:"delay"`
}

type RestartBuilder interface {
	ToRestartMap() (map[string]interface{}, error)
}

func (opts RestartOpts) ToRestartMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Restart(client *golangsdk.ServiceClient, instanceId string, opts RestartBuilder) (r JobResult) {
	b, err := opts.ToRestartMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(updateURL(client, instanceId, "restart"), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return
}

type UpdatePrivateIpOpts struct {
	InternalIp string `json:"internal_ip" required:"true"`
}

type UpdatePrivateIpBuilder interface {
	ToPrivateIpUpdateMap() (map[string]interface{}, error)
}

func (opts UpdatePrivateIpOpts) ToPrivateIpUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UpdatePrivateIp(client *golangsdk.ServiceClient, instanceId string, opts UpdatePrivateIpBuilder) (r JobResult) {
	b, err := opts.ToPrivateIpUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updateURL(client, instanceId, "internal-ip"), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return
}

type UpdatePortOpts struct {
	Port int `json:"port" required:"true"`
}

type UpdatePortBuilder interface {
	ToPortUpdateMap() (map[string]interface{}, error)
}

func (opts UpdatePortOpts) ToPortUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UpdatePort(client *golangsdk.ServiceClient, instanceId string, opts UpdatePortBuilder) (r JobResult) {
	b, err := opts.ToPortUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updateURL(client, instanceId, "port"), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return
}

type UpdateSecurityGroupOpts struct {
	SecurityGroupId string `json:"security_group_id" required:"true"`
}

type UpdateSecurityGroupBuilder interface {
	ToSecurityGroupUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateSecurityGroupOpts) ToSecurityGroupUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UpdateSecurityGroup(client *golangsdk.ServiceClient, instanceId string, opts UpdateSecurityGroupBuilder) (r JobResult) {
	b, err := opts.ToSecurityGroupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updateURL(client, instanceId, "security-group"), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return
}

type UpdateSslOptionOpts struct {
	SslOption bool `json:"ssl_option"`
}

type UpdateSslOptionBuilder interface {
	ToSslOptionUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateSslOptionOpts) ToSslOptionUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UpdateSslOption(client *golangsdk.ServiceClient, instanceId string, opts UpdateSslOptionBuilder) (r JobResult) {
	b, err := opts.ToSslOptionUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updateURL(client, instanceId, "ssl-option"), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return
}

type UpdateAliasOpts struct {
	Alias string `json:"alias"`
}

type UpdateAliasBuilder interface {
	ToAliasUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateAliasOpts) ToAliasUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UpdateAlias(client *golangsdk.ServiceClient, instanceId string, opts UpdateAliasBuilder) (r UpdateAliasResult) {
	b, err := opts.ToAliasUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updateURL(client, instanceId, "alias"), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return
}

type UpdateMaintenanceWindowOpts struct {
	StartTime string `json:"start_time" required:"true"`
	EndTime   string `json:"end_time" required:"true"`
}

type UpdateMaintenanceWindowBuilder interface {
	ToMaintenanceWindowUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateMaintenanceWindowOpts) ToMaintenanceWindowUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UpdateMaintenanceWindow(client *golangsdk.ServiceClient, instanceId string,
	opts UpdateMaintenanceWindowBuilder) (r UpdateMaintenanceWindowResult) {
	b, err := opts.ToMaintenanceWindowUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updateURL(client, instanceId, "ops-window"), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return
}

type UpdateSecondLevelMonitoringOpts struct {
	MonitorSwitch bool `json:"monitor_switch"`
	Period        int  `json:"period,omitempty"`
}

type UpdateSecondLevelMonitoringBuilder interface {
	ToSecondLevelMonitoringUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateSecondLevelMonitoringOpts) ToSecondLevelMonitoringUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UpdateSecondLevelMonitoring(client *golangsdk.ServiceClient, instanceId string,
	opts UpdateSecondLevelMonitoringBuilder) (r JobResult) {
	b, err := opts.ToSecondLevelMonitoringUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updateURL(client, instanceId, "monitor-policy"), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return
}

func GetSecondLevelMonitoring(client *golangsdk.ServiceClient, instanceId string) (r GetSecondLevelMonitoringResult) {
	url := secondLevelMonitoringURL(client, instanceId)

	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return
}

type ApplyPrivateDnsNameOpts struct {
	DnsType string `json:"dns_type" required:"true"`
}

type ApplyPrivateDnsNameBuilder interface {
	ToPrivateDnsNameApplyMap() (map[string]interface{}, error)
}

func (opts ApplyPrivateDnsNameOpts) ToPrivateDnsNameApplyMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ApplyPrivateDnsName(client *golangsdk.ServiceClient, instanceId string, opts ApplyPrivateDnsNameBuilder) (r JobResult) {
	b, err := opts.ToPrivateDnsNameApplyMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(updateURL(client, instanceId, "dns"), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return
}

type UpdatePrivateDnsNameOpts struct {
	DnsName string `json:"dns_name" required:"true"`
}

type UpdatePrivateDnsNameBuilder interface {
	ToPrivateDnsNameUpdateMap() (map[string]interface{}, error)
}

func (opts UpdatePrivateDnsNameOpts) ToPrivateDnsNameUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UpdatePrivateDnsName(client *golangsdk.ServiceClient, instanceId string, opts UpdatePrivateDnsNameBuilder) (r JobResult) {
	b, err := opts.ToPrivateDnsNameUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updateURL(client, instanceId, "dns"), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return
}

func GetVersion(client *golangsdk.ServiceClient, instanceId string) (r GetVersionResult) {
	url := versionURL(client, instanceId)

	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return
}

type UpdateSlowLogShowOriginalSwitchOpts struct {
	OpenSlowLogSwitch bool `json:"open_slow_log_switch"`
}

type UpdateSlowLogShowOriginalSwitchBuilder interface {
	ToSlowLogShowOriginalSwitchUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateSlowLogShowOriginalSwitchOpts) ToSlowLogShowOriginalSwitchUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UpdateSlowLogShowOriginalSwitch(client *golangsdk.ServiceClient, instanceId string,
	opts UpdateSlowLogShowOriginalSwitchBuilder) (r UpdateSlowLogShowOriginalSwitchResult) {
	b, err := opts.ToSlowLogShowOriginalSwitchUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(updateURL(client, instanceId, "slowlog/modify"), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return
}

func GetSlowLogShowOriginalSwitch(client *golangsdk.ServiceClient, instanceId string) (r GetSlowLogShowOriginalSwitchResult) {
	url := slowLogShowOriginalSwitchURL(client, instanceId)

	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	return
}
