package hss

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	hssv5 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/hss/v5"
	hssv5model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/hss/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS POST /v5/{project_id}/webtamper/static/status
// @API HSS POST /v5/{project_id}/webtamper/rasp/status
// @API HSS GET /v5/{project_id}/webtamper/hosts
// @API HSS GET /v5/{project_id}/host-management/hosts
func ResourceWebTamperProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWebTamperProtectionCreate,
		ReadContext:   resourceWebTamperProtectionRead,
		UpdateContext: resourceWebTamperProtectionUpdate,
		DeleteContext: resourceWebTamperProtectionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceWebTamperProtectionImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"is_dynamics_protect": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_bit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protect_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rasp_protect_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"anti_tampering_times": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"detect_tampering_times": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func openWebTamperProtection(client *hssv5.HssClient, cfg *config.Config, hostId string, d *schema.ResourceData) error {
	var (
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
	)

	requestBody := hssv5model.SetWtpProtectionStatusRequestInfo{
		Status:       true,
		HostIdList:   []string{hostId},
		ChargingMode: utils.String(chargingModePacketCycle),
		ResourceId:   utils.StringIgnoreEmpty(d.Get("quota_id").(string)),
	}
	opts := hssv5model.SetWtpProtectionStatusInfoRequest{
		Region:              region,
		EnterpriseProjectId: &epsId,
		Body:                &requestBody,
	}

	_, err := client.SetWtpProtectionStatusInfo(&opts)

	return err
}

func openOrCloseDynamicProtection(client *hssv5.HssClient, cfg *config.Config, hostId string, d *schema.ResourceData) error {
	var (
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
	)

	requestBody := hssv5model.SetRaspSwitchRequestInfo{
		HostIdList: &[]string{hostId},
		Status:     utils.Bool(d.Get("is_dynamics_protect").(bool)),
	}
	opts := hssv5model.SetRaspSwitchRequest{
		Region:              region,
		EnterpriseProjectId: &epsId,
		Body:                &requestBody,
	}

	_, err := client.SetRaspSwitch(&opts)

	return err
}

func resourceWebTamperProtectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
		hostId = d.Get("host_id").(string)
	)

	client, err := cfg.HcHssV5Client(region)
	if err != nil {
		return diag.Errorf("error creating HSS v5 client: %s", err)
	}

	checkHostAvailableErr := checkHostAvailable(client, region, epsId, hostId)
	if checkHostAvailableErr != nil {
		return diag.FromErr(checkHostAvailableErr)
	}

	err = openWebTamperProtection(client, cfg, hostId, d)
	if err != nil {
		return diag.Errorf("error opening HSS web tamper protection: %s", err)
	}

	d.SetId(hostId)

	if d.Get("is_dynamics_protect").(bool) {
		err = openOrCloseDynamicProtection(client, cfg, hostId, d)
		if err != nil {
			return diag.Errorf("error opening HSS dynamic web tamper protection: %s", err)
		}
	}

	return resourceWebTamperProtectionRead(ctx, d, meta)
}

func resourceWebTamperProtectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		hostId = d.Id()
		epsId  = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.HcHssV5Client(region)
	if err != nil {
		return diag.Errorf("error creating HSS v5 client: %s", err)
	}

	// If the enterprise project ID is not set during query, query all enterprise projects.
	if epsId == "" {
		epsId = QueryAllEpsValue
	}
	listOpts := hssv5model.ListWtpProtectHostRequest{
		Region:              region,
		EnterpriseProjectId: utils.String(epsId),
		HostId:              utils.String(hostId),
		ProtectStatus:       utils.String(string(ProtectStatusOpened)),
	}

	resp, err := client.ListWtpProtectHost(&listOpts)
	if err != nil {
		return diag.Errorf("error querying HSS web tamper protection hosts: %s", err)
	}

	if resp == nil || resp.DataList == nil {
		return diag.Errorf("the host (%s) for HSS web tamper protection does not exist", hostId)
	}

	wtpProtectHostList := *resp.DataList
	if len(wtpProtectHostList) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "HSS web tamper protection")
	}

	host := wtpProtectHostList[0]
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("host_id", host.HostId),
		d.Set("host_name", host.HostName),
		d.Set("public_ip", host.PublicIp),
		d.Set("private_ip", host.PrivateIp),
		d.Set("group_name", host.GroupName),
		d.Set("os_bit", host.OsBit),
		d.Set("os_type", host.OsType),
		d.Set("protect_status", host.ProtectStatus),
		d.Set("rasp_protect_status", host.RaspProtectStatus),
		d.Set("anti_tampering_times", host.AntiTamperingTimes),
		d.Set("detect_tampering_times", host.DetectTamperingTimes),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceWebTamperProtectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
		hostId = d.Id()
	)

	client, err := cfg.HcHssV5Client(region)
	if err != nil {
		return diag.Errorf("error creating HSS v5 client: %s", err)
	}

	checkHostAvailableErr := checkHostAvailable(client, region, epsId, hostId)
	if checkHostAvailableErr != nil {
		return diag.FromErr(checkHostAvailableErr)
	}

	if d.HasChange("is_dynamics_protect") {
		err = openOrCloseDynamicProtection(client, cfg, hostId, d)
		if err != nil {
			return diag.Errorf("error updating HSS dynamic web tamper protection: %s", err)
		}
	}

	return resourceWebTamperProtectionRead(ctx, d, meta)
}

func resourceWebTamperProtectionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
		hostId = d.Id()
	)

	client, err := cfg.HcHssV5Client(region)
	if err != nil {
		return diag.Errorf("error creating HSS v5 client: %s", err)
	}

	opts := hssv5model.SetWtpProtectionStatusInfoRequest{
		Region:              region,
		EnterpriseProjectId: &epsId,
		Body: &hssv5model.SetWtpProtectionStatusRequestInfo{
			Status:     false,
			HostIdList: []string{hostId},
		},
	}

	_, err = client.SetWtpProtectionStatusInfo(&opts)
	if err != nil {
		// Repeatedly closing web tamper protection, API will not report errors.
		// If the host does not exist, closing web tamper protection will result in an error as follows:
		// {"error_code": "00000010","error_description": "拒绝访问"}
		// The API documentation does not provide any explanatory information about this error,
		// so the logic of checkDeleted is not added.
		return diag.Errorf("error closing HSS web tamper protection: %s", err)
	}

	return nil
}

func resourceWebTamperProtectionImportState(_ context.Context, d *schema.ResourceData, meta interface{}) (
	[]*schema.ResourceData, error) {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		hostId = d.Id()
	)

	client, err := cfg.HcHssV5Client(region)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating HSS v5 client: %s", err)
	}

	checkHostAvailableErr := checkHostAvailable(client, region, QueryAllEpsValue, hostId)

	return []*schema.ResourceData{d}, checkHostAvailableErr
}
