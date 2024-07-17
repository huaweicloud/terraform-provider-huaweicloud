package hss

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	hssv5 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/hss/v5"
	hssv5model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/hss/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	QueryAllEpsValue              string = "all_granted_eps"
	protectionVersionNull         string = "hss.version.null"
	hostAgentStatusOnline         string = "online"
	chargingModePacketCycle       string = "packet_cycle"
	chargingModeOnDemand          string = "on_demand"
	chargingModePrePaid           string = "prePaid"
	chargingModePostPaid          string = "postPaid"
	getProtectionHostNeedRetryMsg string = "The host cannot be found temporarily, please try again later."
)

// @API HSS GET /v5/{project_id}/host-management/hosts
// @API HSS POST /v5/{project_id}/host-management/protection
func ResourceHostProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostProtectionCreate,
		ReadContext:   resourceHostProtectionRead,
		UpdateContext: resourceHostProtectionUpdate,
		DeleteContext: resourceHostProtectionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceHostProtectionImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
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
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"quota_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_wait_host_available": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			// Attributes
			"host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"agent_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"agent_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"detect_result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"asset_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"open_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getProtectionHost(client *hssv5.HssClient, region, epsId, hostId string) (*hssv5model.Host, error) {
	request := hssv5model.ListHostStatusRequest{
		Region:              &region,
		EnterpriseProjectId: utils.StringIgnoreEmpty(epsId),
		HostId:              utils.String(hostId),
	}

	resp, err := client.ListHostStatus(&request)
	if err != nil {
		return nil, fmt.Errorf("error querying HSS hosts: %s", err)
	}

	if resp == nil || resp.DataList == nil {
		return nil, fmt.Errorf("the host (%s) for HSS host protection does not exist", hostId)
	}

	hostList := *resp.DataList
	if len(hostList) == 0 {
		return nil, fmt.Errorf("%s", getProtectionHostNeedRetryMsg)
	}

	return &hostList[0], nil
}

func checkHostAvailable(client *hssv5.HssClient, region, epsId, hostId string) error {
	request := hssv5model.ListHostStatusRequest{
		Region:              &region,
		EnterpriseProjectId: utils.StringIgnoreEmpty(epsId),
		HostId:              utils.String(hostId),
	}

	resp, err := client.ListHostStatus(&request)
	if err != nil {
		return fmt.Errorf("error querying HSS hosts: %s", err)
	}

	if resp == nil || resp.DataList == nil {
		return fmt.Errorf("the host (%s) for HSS host protection does not exist", hostId)
	}

	hostList := *resp.DataList
	if len(hostList) == 0 {
		return fmt.Errorf("the host (%s) does not exist", hostId)
	}

	agentStatus := *hostList[0].AgentStatus
	if agentStatus != hostAgentStatusOnline {
		return fmt.Errorf("the host anget status for HSS protection must be: %s,"+
			" but the current host (%s) agent status is: %s ", hostAgentStatusOnline, hostId, agentStatus)
	}

	return nil
}

func convertChargingModeRequest(chargingMode string) string {
	switch chargingMode {
	case chargingModePrePaid:
		return chargingModePacketCycle
	case chargingModePostPaid:
		return chargingModeOnDemand
	default:
		return chargingMode
	}
}

func closeHostProtection(client *hssv5.HssClient, region, epsId, hostId string) error {
	closeOpts := hssv5model.SwitchHostsProtectStatusRequest{
		Region:              region,
		EnterpriseProjectId: utils.StringIgnoreEmpty(epsId),
		Body: &hssv5model.SwitchHostsProtectStatusRequestInfo{
			Version:    protectionVersionNull,
			HostIdList: []string{hostId},
		},
	}

	_, err := client.SwitchHostsProtectStatus(&closeOpts)

	return err
}

func switchHostsProtectStatus(client *hssv5.HssClient, region, epsId, hostId string, d *schema.ResourceData) error {
	var (
		version      = d.Get("version").(string)
		chargingMode = d.Get("charging_mode").(string)
		quotaId      = d.Get("quota_id").(string)
	)

	switchOpts := hssv5model.SwitchHostsProtectStatusRequest{
		Region:              region,
		EnterpriseProjectId: utils.StringIgnoreEmpty(epsId),
		Body: &hssv5model.SwitchHostsProtectStatusRequestInfo{
			Version:      version,
			ChargingMode: utils.String(convertChargingModeRequest(chargingMode)),
			ResourceId:   utils.StringIgnoreEmpty(quotaId),
			HostIdList:   []string{hostId},
		},
	}

	_, err := client.SwitchHostsProtectStatus(&switchOpts)
	if err != nil {
		return err
	}

	return nil
}

func waitingForHostAvailable(ctx context.Context, client *hssv5.HssClient, region, epsId, hostId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			host, err := getProtectionHost(client, region, epsId, hostId)
			if err != nil {
				if err.Error() == getProtectionHostNeedRetryMsg {
					return nil, "PENDING", nil
				}

				return nil, "ERROR", err
			}

			agentStatus := utils.StringValue(host.AgentStatus)
			if agentStatus == hostAgentStatusOnline {
				return host, "COMPLETED", nil
			}

			return host, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func resourceHostProtectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if d.Get("is_wait_host_available").(bool) {
		if err := waitingForHostAvailable(ctx, client, region, epsId, hostId, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("error waiting for host (%s) agent status to become online: %s", hostId, err)
		}
	} else {
		checkHostAvailableErr := checkHostAvailable(client, region, epsId, hostId)
		if checkHostAvailableErr != nil {
			return diag.FromErr(checkHostAvailableErr)
		}
	}

	// Due to API limitations, when switching host protection for the first time, protection needs to be closed first.
	err = closeHostProtection(client, region, epsId, hostId)
	if err != nil {
		return diag.Errorf("error closing host protection before opening HSS host (%s) protection: %s",
			hostId, err)
	}

	err = switchHostsProtectStatus(client, region, epsId, hostId, d)
	if err != nil {
		return diag.Errorf("error opening HSS host (%s) protection: %s", hostId, err)
	}

	d.SetId(hostId)

	return resourceHostProtectionRead(ctx, d, meta)
}

func resourceHostProtectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		id     = d.Id()
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
	listHostOpts := hssv5model.ListHostStatusRequest{
		Region:              &region,
		EnterpriseProjectId: utils.String(epsId),
		HostId:              utils.String(id),
	}

	resp, err := client.ListHostStatus(&listHostOpts)
	if err != nil {
		return diag.Errorf("error querying HSS hosts: %s", err)
	}

	if resp == nil || resp.DataList == nil {
		return diag.Errorf("the host (%s) for HSS host protection does not exist", id)
	}

	hostList := *resp.DataList
	if len(hostList) == 0 || utils.StringValue(hostList[0].ProtectStatus) == string(ProtectStatusClosed) {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "HSS host protection")
	}

	host := hostList[0]
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("host_id", host.HostId),
		d.Set("version", host.Version),
		d.Set("charging_mode", convertChargingMode(host.ChargingMode)),
		d.Set("enterprise_project_id", host.EnterpriseProjectId),
		d.Set("host_name", host.HostName),
		d.Set("host_status", host.HostStatus),
		d.Set("private_ip", host.PrivateIp),
		d.Set("agent_id", host.AgentId),
		d.Set("agent_status", host.AgentStatus),
		d.Set("os_type", host.OsType),
		d.Set("status", host.ProtectStatus),
		d.Set("detect_result", host.DetectResult),
		d.Set("asset_value", host.AssetValue),
		d.Set("open_time", convertOpenTime(host.OpenTime)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func convertChargingMode(chargingMode *string) string {
	if utils.StringValue(chargingMode) == chargingModePacketCycle {
		return chargingModePrePaid
	}

	return chargingModePostPaid
}

func convertOpenTime(openTime *int64) string {
	if openTime == nil {
		return ""
	}

	return utils.FormatTimeStampRFC3339(*openTime/1000, false)
}

func resourceHostProtectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
		id     = d.Id()
	)

	client, err := cfg.HcHssV5Client(region)
	if err != nil {
		return diag.Errorf("error creating HSS v5 client: %s", err)
	}

	checkHostAvailableErr := checkHostAvailable(client, region, epsId, id)
	if checkHostAvailableErr != nil {
		return diag.FromErr(checkHostAvailableErr)
	}

	if d.HasChanges("version", "charging_mode", "quota_id") {
		// Due to API limitations, when switching host protection for the first time, protection needs to be closed first.
		err = closeHostProtection(client, region, epsId, id)
		if err != nil {
			return diag.Errorf("error closing host protection before updating HSS host (%s) protection: %s",
				id, err)
		}

		err = switchHostsProtectStatus(client, region, epsId, id, d)
		if err != nil {
			return diag.Errorf("error updating HSS host (%s) protection: %s", id, err)
		}
	}

	return resourceHostProtectionRead(ctx, d, meta)
}

func resourceHostProtectionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
		id     = d.Id()
	)

	client, err := cfg.HcHssV5Client(region)
	if err != nil {
		return diag.Errorf("error creating HSS v5 client: %s", err)
	}

	err = closeHostProtection(client, region, epsId, id)
	if err != nil {
		// Repeatedly closing host protection, API will not report errors.
		// If the host does not exist, closing host protection will result in an error as follows:
		// {"error_code": "00000010","error_description": "拒绝访问"}
		// The API documentation does not provide any explanatory information about this error,
		// so the logic of checkDeleted is not added.
		return diag.Errorf("error closing HSS host (%s) protection: %s", id, err)
	}

	return nil
}

func resourceHostProtectionImportState(_ context.Context, d *schema.ResourceData, meta interface{}) (
	[]*schema.ResourceData, error) {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		id     = d.Id()
	)

	client, err := cfg.HcHssV5Client(region)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating HSS v5 client: %s", err)
	}

	checkHostAvailableErr := checkHostAvailable(client, region, QueryAllEpsValue, id)

	return []*schema.ResourceData{d}, checkHostAvailableErr
}
