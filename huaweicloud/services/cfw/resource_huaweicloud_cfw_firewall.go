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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	FirewallNotExistsCode = "CFW.00200005"
)

// API: CFW POST v2/{project_id}/firewall
// API: CFW GET v3/{project_id}/jobs/{id}
// API: CFW POST v1/{project_id}/firewall/east-west
// API: CFW GET v1/{project_id}/firewall/east-west
// API: CFW GET v1/{project_id}/firewall/exist
// API: CFW POST v1/{project_id}/firewall/east-west/protect
// API: CFW DELETE v2/{project_id}/firewall/{id}

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
			"tags": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the key/value pairs to associate with the firewall.`,
			},
			"east_west_firewall_inspection_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the inspection cidr of the east-west firewall.`,
			},
			"east_west_firewall_er_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the ER ID of the east-west firewall.`,
			},
			"east_west_firewall_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the mode of the east-west firewall.`,
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

	if _, ok := d.GetOk("east_west_firewall_inspection_cidr"); ok {
		// create east west firewall
		var (
			createEastWestFirewallHttpUrl = "v1/{project_id}/firewall/east-west"
		)

		createEastWestFirewallPath := createFirewallClient.Endpoint + createEastWestFirewallHttpUrl
		createEastWestFirewallPath = strings.ReplaceAll(createEastWestFirewallPath, "{project_id}", createFirewallClient.ProjectID)
		createEastWestFirewallPath += fmt.Sprintf("?fw_instance_id=%s", d.Id())

		createEastWestFirewallOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		createEastWestFirewallOpt.JSONBody = utils.RemoveNil(buildCreateEastWestFirewallBodyParams(d))
		_, err := createFirewallClient.Request("POST", createEastWestFirewallPath, &createEastWestFirewallOpt)
		if err != nil {
			return diag.Errorf("error creating east-west firewall: %s", err)
		}

		err = createEastWestFirewallWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for the east-west firewall (%s) creation to complete: %s", d.Id(), err)
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
		"enterprise_project_id": utils.ValueIngoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}
	return bodyParams
}

func buildCreateFirewallRequestBodyChargeInfo(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"charge_mode":   utils.ValueIngoreEmpty(d.Get("charging_mode")),
		"period_type":   utils.ValueIngoreEmpty(d.Get("period_unit")),
		"period_num":    utils.ValueIngoreEmpty(d.Get("period")),
		"is_auto_renew": utils.ValueIngoreEmpty(d.Get("auto_renew")),
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
			"extend_eip_count": utils.ValueIngoreEmpty(raw["extend_eip_count"]),
			"extend_bandwidth": utils.ValueIngoreEmpty(raw["extend_bandwidth"]),
			"extend_vpc_count": utils.ValueIngoreEmpty(raw["extend_vpc_count"]),
		}
		return params
	}
	return nil
}

func buildCreateEastWestFirewallBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"er_id":           utils.ValueIngoreEmpty(d.Get("east_west_firewall_er_id")),
		"inspection_cidr": utils.ValueIngoreEmpty(d.Get("east_west_firewall_inspection_cidr")),
		"mode":            utils.ValueIngoreEmpty(d.Get("east_west_firewall_mode")),
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
			statusRaw, err := jmespath.Search(`data.status`, createFirewallWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `data.status`)
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
			statusRaw, err := jmespath.Search(`data.status`, createEastWestFirewallWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `data.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

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
		if hasErrorCode(err, FirewallNotExistsCode) {
			err = golangsdk.ErrDefault404{}
		}
		return common.CheckDeletedDiag(d, err, "error retrieving firewalls")
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

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("fw_instance_name", getFirewallRespBody, nil)),
		d.Set("charging_mode", chargingMode),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", getFirewallRespBody, nil)),
		d.Set("engine_type", utils.PathSearch("engine_type", getFirewallRespBody, nil)),
		d.Set("ha_type", utils.PathSearch("ha_type", getFirewallRespBody, nil)),
		d.Set("protect_objects", flattenGetFirewallResponseBodyProtectObject(getFirewallRespBody)),
		d.Set("service_type", utils.PathSearch("service_type", getFirewallRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getFirewallRespBody, nil)),
		d.Set("support_ipv6", utils.PathSearch("support_ipv6", getFirewallRespBody, nil)),
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
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func hasErrorCode(err error, expectCode string) bool {
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var response interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &response); jsonErr == nil {
			errorCode, parseErr := jmespath.Search("error_code", response)
			if parseErr != nil {
				log.Printf("[WARN] failed to parse error_code from response body: %s", parseErr)
			}

			if errorCode == expectCode {
				return true
			}
		}
	}

	return false
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

func resourceFirewallUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var objectID string
	if d.IsNewResource() {
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
			if hasErrorCode(err, FirewallNotExistsCode) {
				err = golangsdk.ErrDefault404{}
			}
			return common.CheckDeletedDiag(d, err, "error retrieving firewalls")
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

		objectID = utils.PathSearch("protect_objects[?type==`1`]|[0].object_id", getFirewallRespBody, "").(string)
	} else {
		protectObjects := d.Get("protect_objects").([]interface{})
		for _, protectObject := range protectObjects {
			p := protectObject.(map[string]interface{})
			if p["type"].(int) == 1 {
				objectID = p["object_id"].(string)
			}
		}
	}

	if objectID != "" {
		var (
			updateEastWestFirewallHttpUrl = "v1/{project_id}/firewall/east-west/protect"
			updateEastWestFirewallProduct = "cfw"
		)

		updateEastWestFirewallClient, err := cfg.NewServiceClient(updateEastWestFirewallProduct, region)
		if err != nil {
			return diag.Errorf("error creating CFW client: %s", err)
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
			return diag.Errorf("error updating east-west Firewall status: %s", err)
		}

		_, err = utils.FlattenResponse(updateEastWestFirewallResp)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceFirewallRead(ctx, d, meta)
}

func resourceFirewallDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	if d.Get("charging_mode") == "prePaid" {
		if err := common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
			return diag.Errorf("error unsubscribing CFW firewall: %s", err)
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
			return diag.Errorf("error deleting Firewall: %s", err)
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

			statusRaw, err := jmespath.Search(`status`, deleteFirewallWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `data.status`)
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
