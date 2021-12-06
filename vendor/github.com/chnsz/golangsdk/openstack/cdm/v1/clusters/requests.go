package clusters

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

type ClusterCreateOpts struct {
	Cluster ClusterRequest `json:"cluster" required:"true"`
	// Whether to enable message notification. If this function is enabled, a maximum of five mobile numbers or email
	// addresses can be set. When a table/file migration job fails or an EIP exception occurs,
	// you will receive a notification by SMS message or email.
	AutoRemind *bool `json:"auto_remind,omitempty"`
	// Mobile number for receiving notifications
	PhoneNum string `json:"phone_num,omitempty"`
	// Email address for receiving notifications
	Email string `json:"email,omitempty"`
}

type ClusterRequest struct {
	// Time for scheduled startup of a CDM cluster. The CDM cluster starts at this time every day.
	ScheduleBootTime string `json:"scheduleBootTime,omitempty"`
	// Whether to enable the scheduled startup/shutdown function.
	// The scheduled startup/shutdown and auto shutdown functions cannot be enabled at the same time.
	IsScheduleBootOff *bool `json:"isScheduleBootOff,omitempty"`
	// Node list.
	Instances []InstanceReq `json:"instances,omitempty"`
	// Cluster information. For details, see the description of the datastore parameter.
	Datastore Datastore `json:"datastore,omitempty"`
	// Time for scheduled shutdown of a CDM cluster. The system shuts down directly at this time every day without waiting for unfinished jobs to complete.
	ScheduleOffTime string `json:"scheduleOffTime,omitempty"`
	// VPC ID, which is used for configuring a network for the cluster
	VpcId string `json:"vpcId,omitempty"`
	// Cluster name
	Name string `json:"name,omitempty"`
	// Enterprise project information. For details, see the description of the sys_tags parameter.
	SysTags []tags.ResourceTag `json:"sys_tags,omitempty"`
	// Whether to enable auto shutdown. The auto shutdown and scheduled startup/shutdown functions cannot be enabled at the same time. When auto shutdown is enabled, if no job is running in the cluster and no scheduled job is created, a cluster will be automatically shut down 15 minutes after it starts running to reduce costs.
	IsAutoOff *bool `json:"isAutoOff,omitempty"`
}

type InstanceReq struct {
	AvailabilityZone string `json:"availability_zone" required:"true"`
	// NIC list. A maximum of two NICs are supported. For details, see the description of the nics parameter.
	Nics      []Nics `json:"nics" required:"true"`
	FlavorRef string `json:"flavorRef" required:"true"`
	// Node type. Currently, only cdm is available.
	Type string `json:"type" required:"true"`
}

type Nics struct {
	SecurityGroupId string `json:"securityGroupId" required:"true"`
	// Subnet ID
	NetId string `json:"net-id" required:"true"`
}

type Datastore struct {
	// Type. Generally, the value is cdm.
	Type string `json:"type,omitempty"`
	// Cluster version
	Version string `json:"version,omitempty"`
}

type ClusterDeleteOpts struct {
	// Number of backup log files. Retain the default value 0.
	KeepLastManualBackup int `q:"keep_last_manual_backup"`
}

type RestartReq struct {
	Restart RestartConfig `json:"restart" required:"true"`
}

type RestartConfig struct {
	// Restart delay, in seconds
	RestartDelayTime *int `json:"restartDelayTime,omitempty"`
	// Restart mode. The options are as follows:
	// IMMEDIATELY: immediate restart
	// GRACEFULL: graceful restart
	// FORCELY: forcible restart
	// SOFTLY: normal restart
	// The default value is IMMEDIATELY.
	RestartMode string `json:"restartMode,omitempty"`
	// Restart level. The options are as follows:
	// SERVICE: service restart
	// VM: VM restart
	// The default value is SERVICE.
	RestartLevel string `json:"restartLevel,omitempty"`
	Type         string `json:"type,omitempty"`
	// Reserved field. When restartLevel is set to SERVICE,
	// this parameter is mandatory and an empty string should be entered.
	Instance string `json:"instance,omitempty"`
	// Reserved field. When restartLevel is set to SERVICE,
	// this parameter is mandatory and an empty string should be entered.
	Group string `json:"group,omitempty"`
}

type StartOpts struct {
	Start map[string]string `json:"start" required:"true"`
}

type StopOpts struct {
	Stop StopConfig `json:"stop" required:"true"`
}

type StopConfig struct {
	// Stop type. The options are as follows:
	// IMMEDIATELY: immediate stop
	// GRACEFULLY: graceful stop
	// Enumeration values:
	// IMMEDIATELY
	// GRACEFULLY
	StopMode string `json:"stopMode"`
	// Stop delay, in seconds. This parameter is valid only when stopMode is set to GRACEFULLY.
	// If the value of this parameter is set to -1, the system waits for all jobs to complete and stops accepting
	// new jobs. If the value of this parameter is greater than 0, the system stops the cluster after
	// the specified time and stops accepting new jobs.
	DelayTime *int `json:"delayTime"`
}

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

func Create(c *golangsdk.ServiceClient, opts ClusterCreateOpts) (*ClusterCreateResult, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r ClusterCreateResult
	_, err = c.Post(createURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &r, err
}

func Delete(c *golangsdk.ServiceClient, clusterId string, opts ClusterDeleteOpts) *golangsdk.ErrResult {
	var r golangsdk.ErrResult

	url := deleteURL(c, clusterId)
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return &r
	}

	_, r.Err = c.DeleteWithBody(url, b, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &r
}

func List(c *golangsdk.ServiceClient) (*ClustersRepsonse, error) {
	var rst ClustersRepsonse
	_, err := c.Get(listURL(c), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Restart(c *golangsdk.ServiceClient, clusterId string, opts RestartReq) (*Job, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst Job
	_, err = c.Post(restartURL(c, clusterId), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Get(c *golangsdk.ServiceClient, clusterId string) (*Cluster, error) {
	var rst Cluster
	_, err := c.Get(getURL(c, clusterId), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Start(c *golangsdk.ServiceClient, clusterId string, opts StartOpts) (*ActionResponse, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst ActionResponse
	_, err = c.Post(actionURL(c, clusterId), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Stop(c *golangsdk.ServiceClient, clusterId string, opts StopConfig) (*ActionResponse, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst ActionResponse
	_, err = c.Post(actionURL(c, clusterId), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}
