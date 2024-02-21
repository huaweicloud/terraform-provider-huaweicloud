package baremetalservers

import (
	"encoding/base64"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
)

type CreateOpts struct {
	ImageRef string `json:"imageRef" required:"true"`

	FlavorRef string `json:"flavorRef" required:"true"`

	Name string `json:"name" required:"true"`

	MetaData MetaData `json:"metadata" required:"true"`

	UserData []byte `json:"-"`

	AdminPass string `json:"adminPass,omitempty"`

	KeyName string `json:"key_name,omitempty"`

	SecurityGroups []SecurityGroup `json:"security_groups,omitempty"`

	Nics []Nic `json:"nics" required:"true"`

	AvailabilityZone string `json:"availability_zone" required:"true"`

	VpcId string `json:"vpcid" required:"true"`

	PublicIp *PublicIp `json:"publicip,omitempty"`

	Count int `json:"count,omitempty"`

	RootVolume *RootVolume `json:"root_volume,omitempty"`

	DataVolumes []DataVolume `json:"data_volumes,omitempty"`

	ExtendParam ServerExtendParam `json:"extendparam" required:"true"`

	SchedulerHints *SchedulerHints `json:"os:scheduler_hints,omitempty"`

	ServerTags []tags.ResourceTag `json:"server_tags,omitempty"`
}

type MetaData struct {
	OpSvcUserId string `json:"op_svc_userid" required:"true"`
	BYOL        string `json:"BYOL,omitempty"`
	AdminPass   string `json:"admin_pass,omitempty"`
	AgencyName  string `json:"agency_name,omitempty"`
}

type SecurityGroup struct {
	ID string `json:"id" required:"true"`
}

type Nic struct {
	SubnetId  string `json:"subnet_id" required:"true"`
	IpAddress string `json:"ip_address,omitempty"`
}

type DeleteNic struct {
	ID string `json:"id" required:"true"`
}

type PublicIp struct {
	Id  string `json:"id,omitempty"`
	Eip *Eip   `json:"eip,omitempty"`
}

type RootVolume struct {
	VolumeType  string `json:"volumetype,omitempty"`
	Size        int    `json:"size,omitempty"`
	ClusterID   string `json:"cluster_id,omitempty"`
	ClusterType string `json:"cluster_type,omitempty"`
}

type DataVolume struct {
	VolumeType  string `json:"volumetype" required:"true"`
	Size        int    `json:"size" required:"true"`
	Shareable   bool   `json:"shareable,omitempty"`
	ClusterID   string `json:"cluster_id,omitempty"`
	ClusterType string `json:"cluster_type,omitempty"`
}

type ServerExtendParam struct {
	ChargingMode string `json:"chargingMode,omitempty"`

	RegionID string `json:"regionID,omitempty"`

	PeriodType string `json:"periodType,omitempty"`

	PeriodNum int `json:"periodNum,omitempty"`

	IsAutoRenew string `json:"isAutoRenew,omitempty"`

	IsAutoPay string `json:"isAutoPay,omitempty"`

	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
}

type SchedulerHints struct {
	DecBaremetal string `json:"dec_baremetal,omitempty"`
}

type Eip struct {
	IpType      string         `json:"iptype" required:"true"`
	BandWidth   BandWidth      `json:"bandwidth" required:"true"`
	ExtendParam EipExtendParam `json:"extendparam" required:"true"`
}

type BandWidth struct {
	Name       string `json:"name,omitempty"`
	ShareType  string `json:"sharetype" required:"true"`
	Id         string `json:"id,omitempty"`
	Size       int    `json:"size" required:"true"`
	ChargeMode string `json:"chargemode,omitempty"`
}

type EipExtendParam struct {
	ChargingMode string `json:"chargingMode" required:"true"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToServerCreateMap() (map[string]interface{}, error)
}

// ToServerCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToServerCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if opts.UserData != nil {
		var userData string
		if _, err := base64.StdEncoding.DecodeString(string(opts.UserData)); err != nil {
			userData = base64.StdEncoding.EncodeToString(opts.UserData)
		} else {
			userData = string(opts.UserData)
		}
		b["user_data"] = &userData
	}

	return map[string]interface{}{"server": b}, nil
}

// CreatePrePaid requests a server to be provisioned to the user in the current tenant.
func CreatePrePaid(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r OrderResult) {
	reqBody, err := opts.ToServerCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), reqBody, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

// Get retrieves a particular Server based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

type UpdateOpts struct {
	Name string `json:"name,omitempty"`
}

type ListOpts struct {
	// Specifies the ID of the BMS flavor.
	FlavorId string `q:"flavor"`
	// Specifies the BMS name.
	Name string `q:"name"`
	// Specifies the BMS status.
	// The value can be: ACTIVE, BUILD, ERROR, HARD_REBOOT, REBOOT or SHUTOFF.
	Status string `q:"status"`
	// Number of records to be queried.
	// The valid value is range from 25 to 1000, defaults to 25.
	Limit int `q:"limit"`
	// Specifies the index position, which starts from the next data record specified by offset.
	Offset int `q:"offset"`
	// Specifies the BMS tag.
	Tags string `q:"tags"`
	// Specifies the reserved ID, which can be used to query BMSs created in a batch.
	ReservationId string `q:"reservation_id"`
	// Specifies the level for details about BMS query results.
	// A higher level indicates more details about BMS query results.
	// Available levels include 4, 3, 2, and 1. The default level is 4.
	Detail string `q:"detail"`
	// Specifies the enterprise project ID of the BMS instance.
	EnterpriseProjectId string `q:"enterprise_project_id"`
}

// List is a method to query the list of BMS instances with **pagination**.
func List(client *golangsdk.ServiceClient, opts ListOpts) ([]CloudServer, error) {
	url := listURL(client)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := InstanceDetailPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractServers(pages)
}

type DeleteNicsOpts struct {
	Nics []DeleteNic `json:"nics" required:"true"`
}

type AddNicsOpts struct {
	Nics []Nic `json:"nics" required:"true"`
}

type UpdateOptsBuilder interface {
	ToServerUpdateMap() (map[string]interface{}, error)
}

type DeleteNicsOptsBuilder interface {
	ToServerDeleteNicsMap() (map[string]interface{}, error)
}

type AddNicsOptsBuilder interface {
	ToServerAddNicsMap() (map[string]interface{}, error)
}

func (opts StopServerOps) ToServerStopServerMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "os-stop")
}

func (opts StartServerOps) ToServerStartServerMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "os-start")
}

func (opts RebootServerOps) ToServerRebootServerMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "reboot")
}

type StartServerOps struct {
	Servers []Servers `json:"servers" required:"true"`
}

type StopServerOps struct {
	// The value can be: HARD and SOFT. Only HARD takes effect.
	Type    string    `json:"type" required:"true"`
	Servers []Servers `json:"servers" required:"true"`
}

type RebootServerOps struct {
	// The value can be: HARD and SOFT. Only HARD takes effect.
	Type    string    `json:"type" required:"true"`
	Servers []Servers `json:"servers" required:"true"`
}

type Servers struct {
	ID string `json:"id" required:"true"`
}

type RebootServerOpsBuilder interface {
	ToServerRebootServerMap() (map[string]interface{}, error)
}

type StartServerOpsBuilder interface {
	ToServerStartServerMap() (map[string]interface{}, error)
}

type StopServerOpsBuilder interface {
	ToServerStopServerMap() (map[string]interface{}, error)
}

func (opts UpdateOpts) ToServerUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"server": b}, nil
}

func Update(client *golangsdk.ServiceClient, id string, ops UpdateOptsBuilder) (r UpdateResult) {
	b, err := ops.ToServerUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(putURL(client, id), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// UpdateMetadata updates (or creates) all the metadata specified by opts for
// the given server ID. This operation does not affect already-existing metadata
// that is not specified by opts.
func UpdateMetadata(client *golangsdk.ServiceClient, id string, opts map[string]interface{}) (r UpdateMetadataResult) {
	b := map[string]interface{}{"metadata": opts}
	_, r.Err = client.Post(metadataURL(client, id), b, &r.Body, nil)
	return
}

func (opts DeleteNicsOpts) ToServerDeleteNicsMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func (opts AddNicsOpts) ToServerAddNicsMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func DeleteNics(client *golangsdk.ServiceClient, id string, ops DeleteNicsOptsBuilder) (r JobResult) {
	reqBody, err := ops.ToServerDeleteNicsMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(deleteNicsURL(client, id), reqBody, &r.Body, nil)
	return
}

func AddNics(client *golangsdk.ServiceClient, id string, ops AddNicsOptsBuilder) (r JobResult) {
	reqBody, err := ops.ToServerAddNicsMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(addNicsURL(client, id), reqBody, &r.Body, nil)
	return
}

func GetJobDetail(client *golangsdk.ServiceClient, jobID string) (r JobResult) {
	_, r.Err = client.Get(jobURL(client, jobID), &r.Body, nil)
	return
}

func RebootServer(client *golangsdk.ServiceClient, ops RebootServerOpsBuilder) (r JobResult) {
	reqBody, err := ops.ToServerRebootServerMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(serverStatusPostURL(client), reqBody, &r.Body, nil)
	return
}

func StartServer(client *golangsdk.ServiceClient, ops StartServerOpsBuilder) (r JobResult) {
	reqBody, err := ops.ToServerStartServerMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(serverStatusPostURL(client), reqBody, &r.Body, nil)
	return
}

func StopServer(client *golangsdk.ServiceClient, ops StopServerOpsBuilder) (r JobResult) {
	reqBody, err := ops.ToServerStopServerMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(serverStatusPostURL(client), reqBody, &r.Body, nil)
	return
}
