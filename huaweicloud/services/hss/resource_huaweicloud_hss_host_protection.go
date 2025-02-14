package hss

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

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

func buildProtectionHostQueryParams(epsId, hostId string) string {
	// When calling the query API, if the enterprise project ID is not set, all enterprise projects will be queried.
	if epsId == "" {
		epsId = QueryAllEpsValue
	}

	return fmt.Sprintf("?enterprise_project_id=%v&host_id=%v", epsId, hostId)
}

func getProtectionHost(client *golangsdk.ServiceClient, epsId, hostId string) (interface{}, error) {
	getPath := client.Endpoint + "v5/{project_id}/host-management/hosts"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildProtectionHostQueryParams(epsId, hostId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HSS host, %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	hostResp := utils.PathSearch("data_list[0]", getRespBody, nil)
	if hostResp == nil {
		return nil, fmt.Errorf("%s", getProtectionHostNeedRetryMsg)
	}

	return hostResp, nil
}

func checkHostAvailable(client *golangsdk.ServiceClient, epsId, hostId string) error {
	host, err := getProtectionHost(client, epsId, hostId)
	if err != nil {
		return err
	}

	agentStatus := utils.PathSearch("agent_status", host, "").(string)
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

func buildSwitchHostProtectionQueryParams(epsId string) string {
	queryParams := ""
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return queryParams
}

func buildCloseHostProtectionBodyParams(hostId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"version":      protectionVersionNull,
		"host_id_list": []string{hostId},
	}
	return bodyParams
}

func closeHostProtection(client *golangsdk.ServiceClient, region, epsId, hostId string) error {
	requestPath := client.Endpoint + "v5/{project_id}/host-management/protection"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildSwitchHostProtectionQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
		JSONBody:         buildCloseHostProtectionBodyParams(hostId),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)

	return err
}

func buildSwitchHostProtectionStatusBodyParams(d *schema.ResourceData, hostId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"version":       d.Get("version").(string),
		"charging_mode": convertChargingModeRequest(d.Get("charging_mode").(string)),
		"resource_id":   utils.StringIgnoreEmpty(d.Get("quota_id").(string)),
		"host_id_list":  []string{hostId},
	}
	return bodyParams
}

func switchHostProtectionStatus(client *golangsdk.ServiceClient, d *schema.ResourceData, region, epsId, hostId string) error {
	requestPath := client.Endpoint + "v5/{project_id}/host-management/protection"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildSwitchHostProtectionQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
		JSONBody:         buildSwitchHostProtectionStatusBodyParams(d, hostId),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)

	return err
}

func waitingForHostAvailable(ctx context.Context, client *golangsdk.ServiceClient, epsId, hostId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			host, err := getProtectionHost(client, epsId, hostId)
			if err != nil {
				if err.Error() == getProtectionHostNeedRetryMsg {
					return nil, "PENDING", nil
				}

				return nil, "ERROR", err
			}

			agentStatus := utils.PathSearch("agent_status", host, "").(string)
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
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		hostId  = d.Get("host_id").(string)
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	if d.Get("is_wait_host_available").(bool) {
		if err := waitingForHostAvailable(ctx, client, epsId, hostId, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("error waiting for host (%s) agent status to become online: %s", hostId, err)
		}
	} else {
		checkHostAvailableErr := checkHostAvailable(client, epsId, hostId)
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

	err = switchHostProtectionStatus(client, d, region, epsId, hostId)
	if err != nil {
		return diag.Errorf("error opening HSS host (%s) protection: %s", hostId, err)
	}

	d.SetId(hostId)

	return resourceHostProtectionRead(ctx, d, meta)
}

func resourceHostProtectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		id      = d.Id()
		epsId   = cfg.GetEnterpriseProjectID(d, QueryAllEpsValue)
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + "v5/{project_id}/host-management/hosts"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildProtectionHostQueryParams(epsId, id)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving HSS host, %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	hostResp := utils.PathSearch("data_list[0]", getRespBody, nil)
	protectStatus := utils.PathSearch("protect_status", hostResp, "").(string)
	if hostResp == nil || protectStatus == string(ProtectStatusClosed) {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "HSS host protection")
	}

	openTime := utils.PathSearch("open_time", hostResp, float64(0)).(float64)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("host_id", utils.PathSearch("host_id", hostResp, nil)),
		d.Set("version", utils.PathSearch("version", hostResp, nil)),
		d.Set("charging_mode", flattenChargingMode(utils.PathSearch("charging_mode", hostResp, "").(string))),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", hostResp, nil)),
		d.Set("host_name", utils.PathSearch("host_name", hostResp, nil)),
		d.Set("host_status", utils.PathSearch("host_status", hostResp, nil)),
		d.Set("private_ip", utils.PathSearch("private_ip", hostResp, nil)),
		d.Set("agent_id", utils.PathSearch("agent_id", hostResp, nil)),
		d.Set("agent_status", utils.PathSearch("agent_status", hostResp, nil)),
		d.Set("os_type", utils.PathSearch("os_type", hostResp, nil)),
		d.Set("status", utils.PathSearch("protect_status", hostResp, nil)),
		d.Set("detect_result", utils.PathSearch("detect_result", hostResp, nil)),
		d.Set("asset_value", utils.PathSearch("asset_value", hostResp, nil)),
		d.Set("open_time", utils.FormatTimeStampRFC3339(int64(openTime)/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenChargingMode(chargingMode string) string {
	switch chargingMode {
	case chargingModePacketCycle:
		return chargingModePrePaid
	case chargingModeOnDemand:
		return chargingModePostPaid
	default:
		return ""
	}
}

func resourceHostProtectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		id      = d.Id()
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	checkHostAvailableErr := checkHostAvailable(client, epsId, id)
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

		err = switchHostProtectionStatus(client, d, region, epsId, id)
		if err != nil {
			return diag.Errorf("error updating HSS host (%s) protection: %s", id, err)
		}
	}

	return resourceHostProtectionRead(ctx, d, meta)
}

func resourceHostProtectionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		id      = d.Id()
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
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
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		id      = d.Id()
		epsId   = cfg.GetEnterpriseProjectID(d, QueryAllEpsValue)
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating HSS client: %s", err)
	}

	checkHostAvailableErr := checkHostAvailable(client, epsId, id)

	return []*schema.ResourceData{d}, checkHostAvailableErr
}
