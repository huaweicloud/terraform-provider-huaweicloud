// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CFW
// ---------------------------------------------------------------

package cfw

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	FirewallNotExistsCode = "CFW.00200005"
)

// @API CFW POST /v2/{project_id}/firewall
// @API CFW GET /v3/{project_id}/jobs/{id}
// @API CFW POST /v1/{project_id}/firewall/east-west
// @API CFW GET /v1/{project_id}/firewall/east-west
// @API CFW GET /v1/{project_id}/firewall/exist
// @API CFW POST /v1/{project_id}/firewall/east-west/protect
// @API CFW DELETE /v2/{project_id}/firewall/{id}
// @API CFW POST /v1/{project_id}/ips/switch
// @API CFW GET /v1/{project_id}/ips/switch
// @API CFW POST /v1/{project_id}/ips/protect
// @API CFW GET /v1/{project_id}/ips/protect
// @API CFW POST /v2/{project_id}/cfw-cfw/{fw_instance_id}/tags/create
// @API CFW DELETE /v2/{project_id}/cfw-cfw/{fw_instance_id}/tags/delete

func ResourceFirewall() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallCreate,
		ReadContext:   resourceFirewallRead,
		UpdateContext: resourceFirewallUpdate,
		DeleteContext: resourceFirewallDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			customdiff.ValidateChange("east_west_firewall_inspection_cidr", func(_ context.Context, old, new, _ any) error {
				// can only update from empty
				if old.(string) != new.(string) && old.(string) != "" {
					return fmt.Errorf("east_west_firewall_inspection_cidr can't be updated")
				}
				return nil
			}),
			customdiff.ValidateChange("east_west_firewall_er_id", func(_ context.Context, old, new, _ any) error {
				// can only update from empty
				if old.(string) != new.(string) && old.(string) != "" {
					return fmt.Errorf("east_west_firewall_er_id can't be updated")
				}
				return nil
			}),
			customdiff.ValidateChange("east_west_firewall_mode", func(_ context.Context, old, new, _ any) error {
				// can only update from empty
				if old.(string) != new.(string) && old.(string) != "" {
					return fmt.Errorf("east_west_firewall_mode can't be updated")
				}
				return nil
			}),
			config.MergeDefaultTags(),
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the firewall name.`,
			},
			"flavor": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        firewallFlavorSchema(),
				ForceNew:    true,
				Description: `Specifies the flavor of the firewall.`,
			},
			"tags": common.TagsSchema(),
			"east_west_firewall_inspection_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"east_west_firewall_er_id", "east_west_firewall_mode"},
				Description:  `Specifies the inspection cidr of the east-west firewall.`,
			},
			"east_west_firewall_er_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"east_west_firewall_inspection_cidr", "east_west_firewall_mode"},
				Description:  `Specifies the ER ID of the east-west firewall.`,
			},
			"east_west_firewall_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"east_west_firewall_inspection_cidr", "east_west_firewall_er_id"},
				Description:  `Specifies the mode of the east-west firewall.`,
			},
			"east_west_firewall_status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the protection statue of the east-west firewall.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the enterprise project ID of the firewall.`,
			},
			"charging_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "postPaid",
				ForceNew:    true,
				Description: `Specifies the charging mode.`,
			},
			"period_unit": common.SchemaPeriodUnit(nil),
			"period":      common.SchemaPeriod(nil),
			"auto_renew":  common.SchemaAutoRenew(nil),
			"ips_switch_status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the IPS patch switch status of the firewall.`,
			},
			"ips_protection_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the IPS protection mode of the firewall.`,
			},
			"engine_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The engine type`,
			},
			"ha_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The HA type.`,
			},
			"service_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The service type.`,
			},
			"protect_objects": {
				Type:        schema.TypeList,
				Elem:        firewallsGetFirewallInstanceResponseRecordProtectObjectVOSchema(),
				Computed:    true,
				Description: `The protect objects list.`,
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The firewall status.`,
			},
			"support_ipv6": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether IPv6 is supported.`,
			},
			"east_west_firewall_inspection_vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The east-west firewall inspection VPC ID.`,
			},
			"east_west_firewall_er_attachment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Enterprise Router and Firewall Connection ID`,
			},
		},
	}
}

func firewallFlavorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the version of the firewall.`,
			},
			"extend_eip_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the extend EIP number of the firewall.`,
			},
			"extend_bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the extend bandwidth of the firewall.`,
			},
			"extend_vpc_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the extend VPC number of the firewall.`,
			},
			"eip_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the EIP number of the firewall.`,
			},
			"vpc_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the VPC number of the firewall.`,
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the bandwidth of the firewall.`,
			},
			"log_storage": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the log storage of the firewall.`,
			},
			"default_eip_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the default EIP number of the firewall.`,
			},
			"default_vpc_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the default VPC number of the firewall.`,
			},
			"default_bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the default bandwidth of the firewall.`,
			},
			"default_log_storage": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the default log storage of the firewall.`,
			},
			"vpc_bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the VPC bandwidth of the firewall.`,
			},
			"used_rule_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the used rule count of the firewall.`,
			},
			"total_rule_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the total rule count of the firewall.`,
			},
		},
	}
	return &sc
}

func resourceFirewallCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createFirewall: Create a CFW firewall.
	var (
		createFirewallHttpUrl = "v2/{project_id}/firewall"
		createFirewallProduct = "cfw"
	)
	createFirewallClient, err := cfg.NewServiceClient(createFirewallProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	createFirewallPath := createFirewallClient.Endpoint + createFirewallHttpUrl
	createFirewallPath = strings.ReplaceAll(createFirewallPath, "{project_id}", createFirewallClient.ProjectID)

	createFirewallOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createFirewallOpt.JSONBody = utils.RemoveNil(buildCreateFirewallBodyParams(d, cfg))
	createFirewallResp, err := createFirewallClient.Request("POST", createFirewallPath, &createFirewallOpt)
	if err != nil {
		return diag.Errorf("error creating Firewall: %s", err)
	}

	createFirewallRespBody, err := utils.FlattenResponse(createFirewallResp)
	if err != nil {
		return diag.FromErr(err)
	}

	orderID := utils.PathSearch("order_id", createFirewallRespBody, "")
	jobID := utils.PathSearch("job_id", createFirewallRespBody, "")

	if orderID != "" {
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}

		err = common.WaitOrderComplete(ctx, bssClient, orderID.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderID.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(resourceId)
	} else {
		d.SetId(jobID.(string))

		err = createFirewallWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for the Firewall (%s) creation to complete: %s", d.Id(), err)
		}
	}

	return resourceFirewallUpdate(ctx, d, meta)
}

func buildCreateFirewallBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                  d.Get("name"),
		"flavor":                buildCreateFirewallRequestBodyFlavor(d.Get("flavor")),
		"tags":                  utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		"charge_info":           buildCreateFirewallRequestBodyChargeInfo(d),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}
	return bodyParams
}

func buildCreateFirewallRequestBodyChargeInfo(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"charge_mode":   utils.ValueIgnoreEmpty(d.Get("charging_mode")),
		"period_type":   utils.ValueIgnoreEmpty(d.Get("period_unit")),
		"period_num":    utils.ValueIgnoreEmpty(d.Get("period")),
		"is_auto_renew": utils.ValueIgnoreEmpty(d.Get("auto_renew")),
		"is_auto_pay":   true,
	}
	return params
}

func buildCreateFirewallRequestBodyFlavor(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"version":          raw["version"],
			"extend_eip_count": utils.ValueIgnoreEmpty(raw["extend_eip_count"]),
			"extend_bandwidth": utils.ValueIgnoreEmpty(raw["extend_bandwidth"]),
			"extend_vpc_count": utils.ValueIgnoreEmpty(raw["extend_vpc_count"]),
		}
		return params
	}
	return nil
}

func createEastWestFirewall(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		createEastWestFirewallHttpUrl = "v1/{project_id}/firewall/east-west"
		createEastWestFirewallProduct = "cfw"
	)
	createEastWestFirewallClient, err := cfg.NewServiceClient(createEastWestFirewallProduct, region)
	if err != nil {
		return fmt.Errorf("error creating CFW client: %s", err)
	}

	createEastWestFirewallPath := createEastWestFirewallClient.Endpoint + createEastWestFirewallHttpUrl
	createEastWestFirewallPath = strings.ReplaceAll(createEastWestFirewallPath, "{project_id}", createEastWestFirewallClient.ProjectID)
	createEastWestFirewallPath += fmt.Sprintf("?fw_instance_id=%s", d.Id())

	createEastWestFirewallOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createEastWestFirewallOpt.JSONBody = utils.RemoveNil(buildCreateEastWestFirewallBodyParams(d))
	_, err = createEastWestFirewallClient.Request("POST", createEastWestFirewallPath, &createEastWestFirewallOpt)
	if err != nil {
		return fmt.Errorf("error creating east-west firewall: %s", err)
	}

	err = createEastWestFirewallWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error waiting for the east-west firewall (%s) creation to complete: %s", d.Id(), err)
	}

	return nil
}

func buildCreateEastWestFirewallBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"er_id":           utils.ValueIgnoreEmpty(d.Get("east_west_firewall_er_id")),
		"inspection_cidr": utils.ValueIgnoreEmpty(d.Get("east_west_firewall_inspection_cidr")),
		"mode":            utils.ValueIgnoreEmpty(d.Get("east_west_firewall_mode")),
	}
	return bodyParams
}

func buildUpdateEastWestFirewallStatusBodyParams(d *schema.ResourceData, objectID string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"object_id": objectID,
		"status":    d.Get("east_west_firewall_status"),
	}
	return bodyParams
}

func createFirewallWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			// createFirewallWaiting: Create a CFW firewall.
			var (
				createFirewallWaitingHttpUrl = "v3/{project_id}/jobs/{id}"
				createFirewallWaitingProduct = "cfw"
			)
			createFirewallWaitingClient, err := cfg.NewServiceClient(createFirewallWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating CFW client: %s", err)
			}

			createFirewallWaitingPath := createFirewallWaitingClient.Endpoint + createFirewallWaitingHttpUrl
			createFirewallWaitingPath = strings.ReplaceAll(createFirewallWaitingPath, "{project_id}", createFirewallWaitingClient.ProjectID)
			createFirewallWaitingPath = strings.ReplaceAll(createFirewallWaitingPath, "{id}", d.Id())

			createFirewallWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}

			createFirewallWaitingResp, err := createFirewallWaitingClient.Request("GET", createFirewallWaitingPath, &createFirewallWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createFirewallWaitingRespBody, err := utils.FlattenResponse(createFirewallWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`data.status`, createFirewallWaitingRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parsing %s from response body", `data.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"Success",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createFirewallWaitingRespBody, "COMPLETED", nil
			}

			pendingStatus := []string{
				"Running",
			}
			if utils.StrSliceContains(pendingStatus, status) {
				return createFirewallWaitingRespBody, "PENDING", nil
			}

			return createFirewallWaitingRespBody, status, nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func createEastWestFirewallWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			// createEastWestFirewallWaiting: Create a CFW firewall.
			var (
				createEastWestFirewallWaitingHttpUrl = "v1/{project_id}/firewall/east-west"
				createEastWestFirewallWaitingProduct = "cfw"
			)
			createEastWestFirewallWaitingClient, err := cfg.NewServiceClient(createEastWestFirewallWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating CFW client: %s", err)
			}

			createEastWestFirewallWaitingPath := createEastWestFirewallWaitingClient.Endpoint + createEastWestFirewallWaitingHttpUrl
			createEastWestFirewallWaitingPath = strings.ReplaceAll(createEastWestFirewallWaitingPath, "{project_id}",
				createEastWestFirewallWaitingClient.ProjectID)
			createEastWestFirewallWaitingPath += fmt.Sprintf("?offset=0&limit=10&fw_instance_id=%s", d.Id())

			createEastWestFirewallWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}

			createEastWestFirewallWaitingResp, err := createEastWestFirewallWaitingClient.Request("GET", createEastWestFirewallWaitingPath,
				&createEastWestFirewallWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createEastWestFirewallWaitingRespBody, err := utils.FlattenResponse(createEastWestFirewallWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`data.status`, createEastWestFirewallWaitingRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parsing %s from response body", `data.status`)
			}

			status := fmt.Sprintf("%v", int(statusRaw.(float64)))

			targetStatus := []string{
				"0", "1",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createEastWestFirewallWaitingRespBody, "COMPLETED", nil
			}

			pendingStatus := []string{
				"2",
			}
			if utils.StrSliceContains(pendingStatus, status) {
				return createEastWestFirewallWaitingRespBody, "PENDING", nil
			}

			return createEastWestFirewallWaitingRespBody, status, nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceFirewallRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getFirewall: Query the List of CFW firewalls
	var (
		getFirewallHttpUrl = "v1/{project_id}/firewall/exist"
		getFirewallProduct = "cfw"
	)
	getFirewallClient, err := cfg.NewServiceClient(getFirewallProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	getFirewallPath := getFirewallClient.Endpoint + getFirewallHttpUrl
	getFirewallPath = strings.ReplaceAll(getFirewallPath, "{project_id}", getFirewallClient.ProjectID)
	getFirewallPath += fmt.Sprintf("?offset=0&limit=10&service_type=0&fw_instance_id=%s", d.Id())

	getFirewallOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getFirewallsResp, err := getFirewallClient.Request("GET", getFirewallPath, &getFirewallOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", FirewallNotExistsCode),
			"error retrieving firewalls",
		)
	}

	getFirewallRespBody, err := utils.FlattenResponse(getFirewallsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("data.records[?fw_instance_id=='%s']|[0]", d.Id())
	getFirewallRespBody = utils.PathSearch(jsonPath, getFirewallRespBody, nil)
	if getFirewallRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	chargMode := utils.PathSearch("charge_mode", getFirewallRespBody, 0)
	chargingMode := "prePaid"
	if chargMode.(float64) == 1 {
		chargingMode = "postPaid"
	}

	internetBorderObjectID := utils.PathSearch("protect_objects[?type==`0`]|[0].object_id", getFirewallRespBody, "").(string)
	mode, err := getIpsProtectMode(getFirewallClient, internetBorderObjectID)
	if err != nil {
		return diag.Errorf("error retrieving IPS protect mode: %s", err)
	}
	status, err := getIpsSwitchStatus(getFirewallClient, internetBorderObjectID)
	if err != nil {
		return diag.Errorf("error retrieving IPS patch switch status: %s", err)
	}

	tags := utils.PathSearch("tags", getFirewallRespBody, "")

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("fw_instance_name", getFirewallRespBody, nil)),
		d.Set("charging_mode", chargingMode),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", getFirewallRespBody, nil)),
		d.Set("flavor", flattenGetFirewallResponseBodyFlavor(getFirewallRespBody, chargingMode)),
		setTagsToState(d, tags.(string)),
		d.Set("engine_type", utils.PathSearch("engine_type", getFirewallRespBody, nil)),
		d.Set("ha_type", utils.PathSearch("ha_type", getFirewallRespBody, nil)),
		d.Set("protect_objects", flattenGetFirewallResponseBodyProtectObject(getFirewallRespBody)),
		d.Set("service_type", utils.PathSearch("service_type", getFirewallRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getFirewallRespBody, nil)),
		d.Set("support_ipv6", utils.PathSearch("support_ipv6", getFirewallRespBody, nil)),
		d.Set("ips_protection_mode", mode),
		d.Set("ips_switch_status", status),
	)

	// get east-west firewall
	var (
		getEastWestFirewallHttpUrl = "v1/{project_id}/firewall/east-west"
	)

	getEastWestFirewallPath := getFirewallClient.Endpoint + getEastWestFirewallHttpUrl
	getEastWestFirewallPath = strings.ReplaceAll(getEastWestFirewallPath, "{project_id}", getFirewallClient.ProjectID)
	getEastWestFirewallPath += fmt.Sprintf("?offset=0&limit=10&fw_instance_id=%s", d.Id())

	getEastWestFirewallOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getEastWestFirewallsResp, err := getFirewallClient.Request("GET", getEastWestFirewallPath, &getEastWestFirewallOpt)
	if err != nil {
		return diag.Errorf("error retrieving east-west firewall %s", err)
	}

	getEastWestFirewallRespBody, err := utils.FlattenResponse(getEastWestFirewallsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("east_west_firewall_inspection_cidr", utils.PathSearch("data.inspection_vpc.cidr", getEastWestFirewallRespBody, nil)),
		d.Set("east_west_firewall_er_id", utils.PathSearch("data.er.id", getEastWestFirewallRespBody, nil)),
		d.Set("east_west_firewall_mode", utils.PathSearch("data.mode", getEastWestFirewallRespBody, nil)),
		d.Set("east_west_firewall_status", utils.PathSearch("data.status", getEastWestFirewallRespBody, nil)),
		d.Set("east_west_firewall_inspection_vpc_id", utils.PathSearch("data.inspection_vpc.id", getEastWestFirewallRespBody, nil)),
		d.Set("east_west_firewall_er_attachment_id", utils.PathSearch("data.er.attachment_id", getEastWestFirewallRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func hasErrorCode(err error, expectCode string) bool {
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var response interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &response); jsonErr == nil {
			errorCode := utils.PathSearch("error_code", response, nil)
			if errorCode == nil {
				log.Printf("[WARN] failed to parse error_code from response body")
			}

			if errorCode == expectCode {
				return true
			}
		}
	}

	return false
}

func flattenGetFirewallResponseBodyFlavor(resp interface{}, chargingMode string) []interface{} {
	curJson := utils.PathSearch("flavor", resp, nil)
	if curJson == nil {
		return nil
	}

	v := int(utils.PathSearch("version", curJson, float64(-1)).(float64))
	version := ""
	if v == 0 {
		version = "Standard"
	} else if v == 1 {
		version = "Professional"
	}

	eipCount := utils.PathSearch("eip_count", curJson, float64(0)).(float64)
	vpcCount := utils.PathSearch("vpc_count", curJson, float64(0)).(float64)
	bandwidth := utils.PathSearch("bandwidth", curJson, float64(0)).(float64)
	defaultEipCount := utils.PathSearch("default_eip_count", curJson, float64(0)).(float64)
	defaultVpcCount := utils.PathSearch("default_vpc_count", curJson, float64(0)).(float64)
	defaultBandwidth := utils.PathSearch("default_bandwidth", curJson, float64(0)).(float64)
	extendEipCount, extendVpcCount, extendBandwidth := 0, 0, 0

	if chargingMode == "prePaid" {
		extendEipCount = int(eipCount) - int(defaultEipCount)
		extendVpcCount = int(vpcCount) - int(defaultVpcCount)
		extendBandwidth = int(bandwidth) - int(defaultBandwidth)
	}

	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"version":             version,
		"extend_eip_count":    extendEipCount,
		"extend_bandwidth":    extendBandwidth,
		"extend_vpc_count":    extendVpcCount,
		"eip_count":           eipCount,
		"vpc_count":           vpcCount,
		"bandwidth":           bandwidth,
		"log_storage":         utils.PathSearch("log_storage", curJson, 0),
		"default_eip_count":   defaultEipCount,
		"default_vpc_count":   defaultVpcCount,
		"default_bandwidth":   defaultBandwidth,
		"default_log_storage": utils.PathSearch("default_log_storage", curJson, 0),
		"vpc_bandwidth":       utils.PathSearch("vpc_bandwidth", curJson, 0),
		"used_rule_count":     utils.PathSearch("used_rule_count", curJson, 0),
		"total_rule_count":    utils.PathSearch("total_rule_count", curJson, 0),
	})
	return rst
}

func setTagsToState(d *schema.ResourceData, tags string) error {
	if tags == "" {
		return nil
	}

	var rst map[string]string
	err := json.Unmarshal([]byte(tags), &rst)
	if err != nil {
		return fmt.Errorf("error parsing tags from API response: %s", err)
	}

	return d.Set("tags", rst)
}

func flattenGetFirewallResponseBodyProtectObject(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("protect_objects", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"object_id":   utils.PathSearch("object_id", v, nil),
			"object_name": utils.PathSearch("object_name", v, nil),
			"type":        utils.PathSearch("type", v, nil),
		})
	}
	return rst
}

func getIpsProtectMode(client *golangsdk.ServiceClient, objectID string) (interface{}, error) {
	getIpsProtectModeHttpUrl := "v1/{project_id}/ips/protect"
	getIpsProtectModePath := client.Endpoint + getIpsProtectModeHttpUrl
	getIpsProtectModePath = strings.ReplaceAll(getIpsProtectModePath, "{project_id}", client.ProjectID)
	getIpsProtectModePath += fmt.Sprintf("?object_id=%s", objectID)

	getIpsProtectModeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getIpsProtectModeResp, err := client.Request("GET", getIpsProtectModePath, &getIpsProtectModeOpt)
	if err != nil {
		return nil, err
	}

	getIpsProtectModeRespBody, err := utils.FlattenResponse(getIpsProtectModeResp)
	if err != nil {
		return nil, err
	}

	mode := utils.PathSearch("data.mode", getIpsProtectModeRespBody, nil)
	if mode == nil {
		return nil, fmt.Errorf("error parsing data.mode from response body")
	}

	return mode, nil
}

func getIpsSwitchStatus(client *golangsdk.ServiceClient, objectID string) (interface{}, error) {
	getIpsSwitchStatusHttpUrl := "v1/{project_id}/ips/switch"
	getIpsSwitchStatusPath := client.Endpoint + getIpsSwitchStatusHttpUrl
	getIpsSwitchStatusPath = strings.ReplaceAll(getIpsSwitchStatusPath, "{project_id}", client.ProjectID)
	getIpsSwitchStatusPath += fmt.Sprintf("?object_id=%s", objectID)

	getIpsSwitchStatusOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getIpsSwitchStatusResp, err := client.Request("GET", getIpsSwitchStatusPath, &getIpsSwitchStatusOpt)
	if err != nil {
		return nil, err
	}

	getIpsSwitchStatusRespBody, err := utils.FlattenResponse(getIpsSwitchStatusResp)
	if err != nil {
		return nil, err
	}

	status := utils.PathSearch("data.virtual_patches_status", getIpsSwitchStatusRespBody, nil)
	if status == nil {
		return nil, fmt.Errorf("error parsing data.virtual_patches_status from response body")
	}
	return status, nil
}

func resourceFirewallUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	if d.HasChanges("east_west_firewall_inspection_cidr", "east_west_firewall_er_id", "east_west_firewall_mode") {
		// create east west firewall
		err := createEastWestFirewall(ctx, d, meta)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	var vpcBoderObjectID, internetBorderObjectID string
	if d.IsNewResource() || d.HasChanges("east_west_firewall_inspection_cidr", "east_west_firewall_er_id", "east_west_firewall_mode") {
		// getFirewall: Query the List of CFW firewalls
		var (
			getFirewallHttpUrl = "v1/{project_id}/firewall/exist"
			getFirewallProduct = "cfw"
		)
		getFirewallClient, err := cfg.NewServiceClient(getFirewallProduct, region)
		if err != nil {
			return diag.Errorf("error creating CFW client: %s", err)
		}

		getFirewallPath := getFirewallClient.Endpoint + getFirewallHttpUrl
		getFirewallPath = strings.ReplaceAll(getFirewallPath, "{project_id}", getFirewallClient.ProjectID)
		getFirewallPath += fmt.Sprintf("?offset=0&limit=10&service_type=0&fw_instance_id=%s", d.Id())

		getFirewallOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getFirewallsResp, err := getFirewallClient.Request("GET", getFirewallPath, &getFirewallOpt)
		if err != nil {
			return common.CheckDeletedDiag(d,
				common.ConvertExpected400ErrInto404Err(err, "error_code", FirewallNotExistsCode),
				"error retrieving firewalls",
			)
		}

		getFirewallRespBody, err := utils.FlattenResponse(getFirewallsResp)
		if err != nil {
			return diag.FromErr(err)
		}

		jsonPath := fmt.Sprintf("data.records[?fw_instance_id=='%s']|[0]", d.Id())
		getFirewallRespBody = utils.PathSearch(jsonPath, getFirewallRespBody, nil)
		if getFirewallRespBody == nil {
			return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
		}

		vpcBoderObjectID = utils.PathSearch("protect_objects[?type==`1`]|[0].object_id", getFirewallRespBody, "").(string)
		internetBorderObjectID = utils.PathSearch("protect_objects[?type==`0`]|[0].object_id", getFirewallRespBody, "").(string)
	} else {
		protectObjects := d.Get("protect_objects").([]interface{})
		for _, protectObject := range protectObjects {
			p := protectObject.(map[string]interface{})
			if p["type"].(int) == 1 {
				vpcBoderObjectID = p["object_id"].(string)
			}
			if p["type"].(int) == 0 {
				internetBorderObjectID = p["object_id"].(string)
			}
		}
	}

	if d.IsNewResource() || d.HasChanges("east_west_firewall_status", "east_west_firewall_inspection_cidr",
		"east_west_firewall_er_id", "east_west_firewall_mode") {
		if vpcBoderObjectID != "" {
			err := updateEastWestFirewallStatus(d, meta, vpcBoderObjectID)
			if err != nil {
				return diag.Errorf("error updating east-west firewall status: %s", err)
			}
		}
	}

	if d.HasChanges("ips_switch_status") {
		err := updateIpsSwitchStatus(d, meta, internetBorderObjectID)
		if err != nil {
			return diag.Errorf("error updating IPS patch switch status: %s", err)
		}
	}

	if d.HasChange("ips_protection_mode") {
		err := updateIpsProtectMode(d, meta, internetBorderObjectID)
		if err != nil {
			return diag.Errorf("error updating IPS protection mode: %s", err)
		}
	}

	if !d.IsNewResource() && d.HasChange("tags") {
		err := updateTags(d, meta)
		if err != nil {
			return diag.Errorf("error updating tags: %s", err)
		}
	}

	return resourceFirewallRead(ctx, d, meta)
}

func updateEastWestFirewallStatus(d *schema.ResourceData, meta interface{}, objectID string) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateEastWestFirewallHttpUrl = "v1/{project_id}/firewall/east-west/protect"
		updateEastWestFirewallProduct = "cfw"
	)

	updateEastWestFirewallClient, err := cfg.NewServiceClient(updateEastWestFirewallProduct, region)
	if err != nil {
		return fmt.Errorf("error creating CFW client: %s", err)
	}

	updateEastWestFirewallPath := updateEastWestFirewallClient.Endpoint + updateEastWestFirewallHttpUrl
	updateEastWestFirewallPath = strings.ReplaceAll(updateEastWestFirewallPath, "{project_id}", updateEastWestFirewallClient.ProjectID)
	updateEastWestFirewallPath += fmt.Sprintf("?fw_instance_id=%s", d.Id())

	updateEastWestFirewallOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	updateEastWestFirewallOpt.JSONBody = utils.RemoveNil(buildUpdateEastWestFirewallStatusBodyParams(d, objectID))
	updateEastWestFirewallResp, err := updateEastWestFirewallClient.Request("POST", updateEastWestFirewallPath, &updateEastWestFirewallOpt)
	if err != nil {
		return fmt.Errorf("error updating east-west firewall status: %s", err)
	}

	_, err = utils.FlattenResponse(updateEastWestFirewallResp)

	return err
}

func updateIpsSwitchStatus(d *schema.ResourceData, meta interface{}, objectID string) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateSwitchStatusHttpUrl = "v1/{project_id}/ips/switch"
		updateSwitchStatusProduct = "cfw"
	)
	updateSwitchStatusClient, err := cfg.NewServiceClient(updateSwitchStatusProduct, region)
	if err != nil {
		return fmt.Errorf("error creating CFW client: %s", err)
	}

	updateSwitchStatusPath := updateSwitchStatusClient.Endpoint + updateSwitchStatusHttpUrl
	updateSwitchStatusPath = strings.ReplaceAll(updateSwitchStatusPath, "{project_id}", updateSwitchStatusClient.ProjectID)
	updateSwitchStatusPath += fmt.Sprintf("?fw_instance_id=%s", d.Id())

	status := d.Get("ips_switch_status").(int)
	updateSwitchStatusOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateIpsSwitchStatusBodyParams(status, objectID)),
	}
	_, err = updateSwitchStatusClient.Request("POST", updateSwitchStatusPath, &updateSwitchStatusOpt)

	return err
}

func buildUpdateIpsSwitchStatusBodyParams(status int, objectID string) map[string]interface{} {
	// ips_type is the patch type, only supports virtual patch, the value is 2.
	return map[string]interface{}{
		"object_id": objectID,
		"ips_type":  2,
		"status":    status,
	}
}

func updateIpsProtectMode(d *schema.ResourceData, meta interface{}, objectID string) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateProtectModeHttpUrl = "v1/{project_id}/ips/protect"
		updateProtectModeProduct = "cfw"
	)
	updateProtectModeClient, err := cfg.NewServiceClient(updateProtectModeProduct, region)
	if err != nil {
		return fmt.Errorf("error creating CFW client: %s", err)
	}

	updateProtectModePath := updateProtectModeClient.Endpoint + updateProtectModeHttpUrl
	updateProtectModePath = strings.ReplaceAll(updateProtectModePath, "{project_id}", updateProtectModeClient.ProjectID)

	mode := d.Get("ips_protection_mode").(int)
	updateProtectModeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateIpsProtectModeBodyParams(mode, objectID)),
	}
	_, err = updateProtectModeClient.Request("POST", updateProtectModePath, &updateProtectModeOpt)

	return err
}

func buildUpdateIpsProtectModeBodyParams(mode int, objectID string) map[string]interface{} {
	return map[string]interface{}{
		"object_id": objectID,
		"mode":      mode,
	}
}

func updateTags(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return fmt.Errorf("error creating CFW client: %s", err)
	}

	oRaw, nRaw := d.GetChange("tags")
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})

	// Remove old tags.
	if len(oMap) > 0 {
		if err := deleteTags(client, oMap, d.Id()); err != nil {
			return err
		}
	}

	// Set new tags.
	if len(nMap) > 0 {
		if err := createTags(client, nMap, d.Id()); err != nil {
			return err
		}
	}
	return nil
}

func createTags(createTagsClient *golangsdk.ServiceClient, tags map[string]interface{}, id string) error {
	if len(tags) == 0 {
		return nil
	}

	createTagsHttpUrl := "v2/{project_id}/cfw-cfw/{fw_instance_id}/tags/create"
	createTagsPath := createTagsClient.Endpoint + createTagsHttpUrl
	createTagsPath = strings.ReplaceAll(createTagsPath, "{project_id}", createTagsClient.ProjectID)
	createTagsPath = strings.ReplaceAll(createTagsPath, "{fw_instance_id}", id)
	createTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"tags": utils.ExpandResourceTags(tags),
		},
	}

	_, err := createTagsClient.Request("POST", createTagsPath, &createTagsOpt)
	if err != nil {
		return fmt.Errorf("error creating tags: %s", err)
	}

	return nil
}

func deleteTags(deleteTagsClient *golangsdk.ServiceClient, tags map[string]interface{}, id string) error {
	if len(tags) == 0 {
		return nil
	}

	deleteTagsHttpUrl := "v2/{project_id}/cfw-cfw/{fw_instance_id}/tags/delete"
	deleteTagsPath := deleteTagsClient.Endpoint + deleteTagsHttpUrl
	deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{project_id}", deleteTagsClient.ProjectID)
	deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{fw_instance_id}", id)
	deleteTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	deleteTagsOpt.JSONBody = map[string]interface{}{
		"tags": utils.ExpandResourceTags(tags),
	}

	_, err := deleteTagsClient.Request("DELETE", deleteTagsPath, &deleteTagsOpt)
	if err != nil {
		return fmt.Errorf("error deleting tags: %s", err)
	}

	return nil
}

func resourceFirewallDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	if d.Get("charging_mode") == "prePaid" {
		if err := common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
			return common.CheckDeletedDiag(d,
				common.ConvertExpected400ErrInto404Err(err, "error_code", "CBC.30000067"),
				"error unsubscribing CFW firewall",
			)
		}
	} else {
		// deleteFirewall: Delete an existing CFW firewall
		var (
			deleteFirewallHttpUrl = "v2/{project_id}/firewall/{id}"
			deleteFirewallProduct = "cfw"
		)
		deleteFirewallClient, err := cfg.NewServiceClient(deleteFirewallProduct, region)
		if err != nil {
			return diag.Errorf("error creating CFW client: %s", err)
		}

		deleteFirewallPath := deleteFirewallClient.Endpoint + deleteFirewallHttpUrl
		deleteFirewallPath = strings.ReplaceAll(deleteFirewallPath, "{project_id}", deleteFirewallClient.ProjectID)
		deleteFirewallPath = strings.ReplaceAll(deleteFirewallPath, "{id}", d.Id())

		deleteFirewallOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		_, err = deleteFirewallClient.Request("DELETE", deleteFirewallPath, &deleteFirewallOpt)
		if err != nil {
			return common.CheckDeletedDiag(d,
				common.ConvertExpected400ErrInto404Err(err, "error_code", FirewallNotExistsCode),
				"error deleting CFW firewall",
			)
		}
	}

	err := deleteFirewallWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the delete of CFW firewall (%s) to complete: %s", d.Id(), err)
	}

	return nil
}

func deleteFirewallWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			// deleteFirewallWaiting: delete a CFW firewall.
			var (
				deleteFirewallWaitingHttpUrl = "v1/{project_id}/firewall/exist"
				deleteFirewallWaitingProduct = "cfw"
			)
			deleteFirewallWaitingClient, err := cfg.NewServiceClient(deleteFirewallWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating CFW client: %s", err)
			}

			deleteFirewallWaitingPath := deleteFirewallWaitingClient.Endpoint + deleteFirewallWaitingHttpUrl
			deleteFirewallWaitingPath = strings.ReplaceAll(deleteFirewallWaitingPath, "{project_id}", deleteFirewallWaitingClient.ProjectID)
			deleteFirewallWaitingPath += fmt.Sprintf("?offset=0&limit=10&service_type=0&fw_instance_id=%s", d.Id())

			deleteFirewallWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}

			deleteFirewallWaitingResp, err := deleteFirewallWaitingClient.Request("GET", deleteFirewallWaitingPath, &deleteFirewallWaitingOpt)
			if err != nil {
				if hasErrorCode(err, FirewallNotExistsCode) {
					return deleteFirewallWaitingResp, "COMPLETED", nil
				}
				return nil, "ERROR", err
			}

			deleteFirewallWaitingRespBody, err := utils.FlattenResponse(deleteFirewallWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			jsonPath := fmt.Sprintf("data.records[?fw_instance_id=='%s']|[0]", d.Id())
			deleteFirewallWaitingRespBody = utils.PathSearch(jsonPath, deleteFirewallWaitingRespBody, nil)
			if deleteFirewallWaitingRespBody == nil {
				return deleteFirewallWaitingRespBody, "COMPLETED", nil
			}

			statusRaw := utils.PathSearch(`status`, deleteFirewallWaitingRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parsing %s from response body", `data.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"4",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return deleteFirewallWaitingRespBody, "COMPLETED", nil
			}

			pendingStatus := []string{
				"1",
				"2",
			}
			if utils.StrSliceContains(pendingStatus, status) {
				return deleteFirewallWaitingRespBody, "PENDING", nil
			}

			return deleteFirewallWaitingRespBody, status, nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
