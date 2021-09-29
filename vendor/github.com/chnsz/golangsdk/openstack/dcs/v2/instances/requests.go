package instances

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dcs/v2/tags"
)

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

type CreateOpts struct {
	Name                  string   `json:"name" required:"true"`
	Engine                string   `json:"engine" required:"true"`
	EngineVersion         string   `json:"engine_version,omitempty"`
	Capacity              float64  `json:"capacity" required:"true"`
	SpecCode              string   `json:"spec_code" required:"true"`
	AzCodes               []string `json:"az_codes" required:"true"`
	VpcId                 string   `json:"vpc_id" required:"true"`
	SubnetId              string   `json:"subnet_id" required:"true"`
	SecurityGroupId       string   `json:"security_group_id,omitempty"`
	PublicIpId            string   `json:"publicip_id,omitempty"`
	EnterpriseProjectId   string   `json:"enterprise_project_id,omitempty"`
	EnterpriseProjectName string   `json:"enterprise_project_name,omitempty"`
	Description           string   `json:"description,omitempty"`
	EnableSsl             *bool    `json:"enable_ssl,omitempty"`
	PrivateIp             string   `json:"private_ip,omitempty"`
	// instance number, the value range is 1-100.
	InstanceNum      int                       `json:"instance_num,omitempty"`
	MaintainBegin    string                    `json:"maintain_begin,omitempty"`
	MaintainEnd      string                    `json:"maintain_end,omitempty"`
	Password         string                    `json:"password,omitempty"`
	NoPasswordAccess *bool                     `json:"no_password_access,omitempty"`
	BssParam         DcsBssParam               `json:"bss_param,omitempty"`
	BackupPolicy     *InstanceBackupPolicyOpts `json:"instance_backup_policy,omitempty"`
	Tags             []tags.ResourceTag        `json:"tags,omitempty"`
	AccessUser       string                    `json:"access_user,omitempty"`
	EnablePublicIp   *bool                     `json:"enable_publicip,omitempty"`
	Port             int                       `json:"port,omitempty"`
	RenameCommands   RedisCommand              `json:"rename_commands,omitempty"`
}

type RedisCommand struct {
	Command  string `json:"command,omitempty"`
	Keys     string `json:"keys,omitempty"`
	Flushdb  string `json:"flushdb,omitempty"`
	Flushall string `json:"flushall,omitempty"`
	Hgetall  string `json:"hgetall,omitempty"`
}

type InstanceBackupPolicyOpts struct {
	BackupType           string     `json:"backup_type" required:"true"`
	SaveDays             int        `json:"save_days,omitempty"`
	PeriodicalBackupPlan BackupPlan `json:"periodical_backup_plan,omitempty"`
}

type BackupPlan struct {
	TimezoneOffset string `json:"timezone_offset,omitempty"`
	BackupAt       []int  `json:"backup_at" required:"true"`
	PeriodType     string `json:"period_type" required:"true"`
	BeginAt        string `json:"begin_at" required:"true"`
}

type DcsBssParam struct {
	ChargingMode string `json:"charging_mode" required:"true"`
	PeriodType   string `json:"period_type,omitempty"`
	// period number, the value range is 1-9.
	PeriodNum   int    `json:"period_num,omitempty"`
	IsAutoRenew string `json:"is_auto_renew,omitempty"`
	IsAutoPay   string `json:"is_auto_pay,omitempty"`
}

func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*CreateResponse, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(rootURL(c), b, &rst.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r CreateResponse
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

type ModifyInstanceOpt struct {
	Name            string                    `json:"name,omitempty"`
	Description     *string                   `json:"description,omitempty"`
	Port            int                       `json:"port,omitempty"`
	MaintainBegin   string                    `json:"maintain_begin,omitempty"`
	MaintainEnd     string                    `json:"maintain_end,omitempty"`
	SecurityGroupId *string                   `json:"security_group_id,omitempty"`
	BackupPolicy    *InstanceBackupPolicyOpts `json:"instance_backup_policy,omitempty"`
}

func Update(c *golangsdk.ServiceClient, id string, opts ModifyInstanceOpt) (*golangsdk.Result, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r golangsdk.Result
	_, err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{204},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &r, err
}

type ResizeInstanceOpts struct {
	SpecCode       string          `json:"spec_code" required:"true"`
	NewCapacity    float64         `json:"new_capacity" required:"true"`
	BssParam       DcsBssParamOpts `json:"bss_param,omitempty"`
	ReservedIp     []string        `json:"reserved_ip,omitempty"`
	ChangeType     string          `json:"change_type,omitempty"`
	AvailableZones []string        `json:"available_zones,omitempty"`
	NodeList       []string        `json:"node_list,omitempty"`
}

type DcsBssParamOpts struct {
	IsAutoPay string `json:"is_auto_pay,omitempty"`
}

func ResizeInstance(c *golangsdk.ServiceClient, id string, opts ResizeInstanceOpts) (*ResizeResponse, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(resizeResourceURL(c, id), b, &rst.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200, 204},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r ResizeResponse
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

func Get(c *golangsdk.ServiceClient, id string) (*DcsInstance, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, id), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r DcsInstance
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

func Delete(c *golangsdk.ServiceClient, id string) error {
	_, err := c.Delete(resourceURL(c, id), &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return err
}

type UpdatePasswordOpts struct {
	OldPassword string `json:"old_password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}

func UpdatePassword(c *golangsdk.ServiceClient, id string, opts UpdatePasswordOpts) (*golangsdk.Result, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r golangsdk.Result
	_, err = c.Put(updatePasswordURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &r, err
}
