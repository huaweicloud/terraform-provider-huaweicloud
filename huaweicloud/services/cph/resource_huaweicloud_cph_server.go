// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CPH
// ---------------------------------------------------------------

package cph

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CPH PUT /v1/{project_id}/cloud-phone/servers/{server_id}
// @API CPH GET /v1/{project_id}/cloud-phone/servers/{server_id}
// @API CPH POST /v2/{project_id}/cloud-phone/servers
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceCphServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCphServerCreate,
		UpdateContext: resourceCphServerUpdate,
		ReadContext:   resourceCphServerRead,
		DeleteContext: resourceCphServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return newValue+"-1" == oldValue
				},
				Description: `Server name.`,
			},
			"server_flavor": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The CPH server flavor.`,
			},
			"phone_flavor": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The cloud phone flavor.`,
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The cloud phone image ID.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of VPC which the cloud server belongs to`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of subnet which the cloud server belongs to`,
			},
			"eip_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eip_type", "bandwidth"},
				Description:   `The ID of an **existing** EIP assigned to the server.`,
			},
			"eip_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eip_id"},
				RequiredWith:  []string{"bandwidth"},
				Description:   `The type of an EIP that will be automatically assigned to the cloud server.`,
			},
			"bandwidth": {
				Type:          schema.TypeList,
				MaxItems:      1,
				ConflictsWith: []string{"eip_id"},
				RequiredWith:  []string{"eip_type"},
				Elem:          cphServerBandWidthSchema(),
				Optional:      true,
				ForceNew:      true,
				Description:   `The bandwidth used by the cloud phone.`,
			},
			"period_unit": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The charging period unit.`,
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, false),
			},
			"period": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `The charging period.`,
			},
			"auto_renew": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Whether auto renew is enabled. Valid values are "true" and "false".`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The name of the AZ where the cloud server is located.`,
			},
			"keypair_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The key pair name, which is used for logging in to the cloud phone through ADB.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `the enterprise project ID.`,
			},
			"ports": {
				Type:        schema.TypeList,
				Elem:        cphServerApplicationPortSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The application port enabled by the cloud phone.`,
			},
			"order_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The order ID.`,
			},
			"addresses": {
				Type:        schema.TypeList,
				Elem:        cphServerAddressSchema(),
				Computed:    true,
				Description: `The IP addresses of the CPH server.`,
			},
			"security_groups": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The list of the security groups bound to the extension NIC of the CPH server.`,
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The CPH server status.`,
			},
		},
	}
}

func cphServerBandWidthSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"share_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The bandwidth type.`,
				ValidateFunc: validation.StringInSlice([]string{
					"0", "1",
				}, false),
			},
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   `The bandwidth ID.`,
				ConflictsWith: []string{"bandwidth.0.size", "bandwidth.0.charge_mode"},
			},
			"size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  `The bandwidth (Mbit/s).`,
				RequiredWith: []string{"bandwidth.0.charge_mode"},
			},
			"charge_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Which the bandwidth used by the CPH server is billed.`,
				ValidateFunc: validation.StringInSlice([]string{
					"0", "1",
				}, false),
				RequiredWith: []string{"bandwidth.0.size"},
			},
		},
	}
	return &sc
}

func cphServerApplicationPortSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The application port name, which can contain a maximum of 16 bytes.`,
			},
			"listen_port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The port number, which ranges from 10000 to 50000.`,
			},
			"internet_accessible": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Whether public network access is mapped.`,
			},
		},
	}
	return &sc
}

func cphServerAddressSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"server_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The internal IP address of the CPH server.`,
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The public IP address of the CPH server.`,
			},
		},
	}
	return &sc
}

func resourceCphServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createCphServer: create a CPH.
	var (
		createCphServerHttpUrl = "v2/{project_id}/cloud-phone/servers"
		createCphServerProduct = "cph"
	)
	createCphServerClient, err := cfg.NewServiceClient(createCphServerProduct, region)
	if err != nil {
		return diag.Errorf("error creating CPH Client: %s", err)
	}

	createCphServerPath := createCphServerClient.Endpoint + createCphServerHttpUrl
	createCphServerPath = strings.ReplaceAll(createCphServerPath, "{project_id}", createCphServerClient.ProjectID)

	createCphServerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createCphServerOpt.JSONBody = utils.RemoveNil(buildCreateCphServerBodyParams(d, cfg))
	createCphServerResp, err := createCphServerClient.Request("POST", createCphServerPath, &createCphServerOpt)
	if err != nil {
		return diag.Errorf("error creating CPH Server: %s", err)
	}

	createCphServerRespBody, err := utils.FlattenResponse(createCphServerResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("server_ids[0]", createCphServerRespBody)
	if err != nil {
		return diag.Errorf("error creating CPH Server: ID is not found in API response")
	}

	d.SetId(id.(string))

	err = createCphServerWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Create of CPH server (%s) to complete: %s", d.Id(), err)
	}

	orderId, err := jmespath.Search("order_id", createCphServerRespBody)
	if err != nil {
		return diag.Errorf("error creating CPH Server: order_id is not found in API response")
	}

	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}
	err = common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Create of CPH server (%s) to complete: %s", d.Id(), err)
	}

	return resourceCphServerRead(ctx, d, meta)
}

func buildCreateCphServerBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"server_name":       utils.ValueIgnoreEmpty(d.Get("name")),
		"server_model_name": utils.ValueIgnoreEmpty(d.Get("server_flavor")),
		"phone_model_name":  utils.ValueIgnoreEmpty(d.Get("phone_flavor")),
		"image_id":          utils.ValueIgnoreEmpty(d.Get("image_id")),
		"count":             1,
		"tenant_vpc_id":     utils.ValueIgnoreEmpty(d.Get("vpc_id")),
		"nics": []map[string]interface{}{
			{
				"subnet_id": utils.ValueIgnoreEmpty(d.Get("subnet_id")),
			},
		},
		"public_ip":         buildCreateCphServerRequestBodyPublicIp(d),
		"band_width":        buildCreateCphServerRequestBodyBandWidth(d.Get("bandwidth")),
		"keypair_name":      utils.ValueIgnoreEmpty(d.Get("keypair_name")),
		"availability_zone": utils.ValueIgnoreEmpty(d.Get("availability_zone")),
		"ports":             buildCreateCphServerRequestBodyApplicationPort(d.Get("ports")),
	}

	extendParam := map[string]interface{}{
		"charging_mode":         0,
		"is_auto_pay":           1,
		"period_num":            utils.ValueIgnoreEmpty(d.Get("period")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}

	periodUnit := d.Get("period_unit").(string)
	if periodUnit == "month" {
		extendParam["period_type"] = 2
	} else {
		extendParam["period_type"] = 3
	}

	autoRenew := d.Get("auto_renew").(string)
	if autoRenew == "true" {
		extendParam["is_auto_renew"] = 1
	} else {
		extendParam["is_auto_renew"] = 0
	}

	bodyParams["extend_param"] = extendParam

	return bodyParams
}

func buildCreateCphServerRequestBodyBandWidth(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"band_width_id":   utils.ValueIgnoreEmpty(raw["id"]),
			"band_width_size": utils.ValueIgnoreEmpty(raw["size"]),
		}

		shareType, _ := strconv.Atoi(utils.ValueIgnoreEmpty(raw["share_type"]).(string))
		params["band_width_share_type"] = shareType

		chargeMode := utils.ValueIgnoreEmpty(raw["charge_mode"]).(string)
		if len(chargeMode) > 0 {
			chargeModeInteger, _ := strconv.Atoi(chargeMode)
			params["band_width_charge_mode"] = chargeModeInteger
		}

		return params
	}
	return nil
}

func buildCreateCphServerRequestBodyPublicIp(d *schema.ResourceData) map[string]interface{} {
	params := make(map[string]interface{})
	if value, ok := d.GetOk("eip_type"); ok {
		params["eip"] = map[string]interface{}{
			"type": utils.ValueIgnoreEmpty(value),
		}
	}
	if value, ok := d.GetOk("eip_id"); ok {
		params["ids"] = []string{value.(string)}
	}
	return params
}

func buildCreateCphServerRequestBodyApplicationPort(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name":                utils.ValueIgnoreEmpty(raw["name"]),
				"listen_port":         utils.ValueIgnoreEmpty(raw["listen_port"]),
				"internet_accessible": utils.ValueIgnoreEmpty(raw["internet_accessible"]),
			}
		}
		return rst
	}
	return nil
}

func createCphServerWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			var (
				createCphServerHttpUrl = "v1/{project_id}/cloud-phone/servers/{server_id}"
				createCphServerProduct = "cph"
			)
			createCphServerClient, err := cfg.NewServiceClient(createCphServerProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating CPH Client: %s", err)
			}

			createCphServerPath := createCphServerClient.Endpoint + createCphServerHttpUrl
			createCphServerPath = strings.ReplaceAll(createCphServerPath, "{project_id}", createCphServerClient.ProjectID)
			createCphServerPath = strings.ReplaceAll(createCphServerPath, "{server_id}", d.Id())

			createCphServerOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			createCphServerResp, err := createCphServerClient.Request("GET", createCphServerPath, &createCphServerOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createCphServerRespBody, err := utils.FlattenResponse(createCphServerResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`status`, createCphServerRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"5", "8", "10",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createCphServerRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"2",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return createCphServerRespBody, status, nil
			}

			return createCphServerRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        120 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceCphServerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getCphServer: Query the CPH instance
	var (
		getCphServerHttpUrl = "v1/{project_id}/cloud-phone/servers/{server_id}"
		getCphServerProduct = "cph"
	)
	getCphServerClient, err := cfg.NewServiceClient(getCphServerProduct, region)
	if err != nil {
		return diag.Errorf("error creating CPH Client: %s", err)
	}

	getCphServerPath := getCphServerClient.Endpoint + getCphServerHttpUrl
	getCphServerPath = strings.ReplaceAll(getCphServerPath, "{project_id}", getCphServerClient.ProjectID)
	getCphServerPath = strings.ReplaceAll(getCphServerPath, "{server_id}", d.Id())

	getCphServerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getCphServerResp, err := getCphServerClient.Request("GET", getCphServerPath, &getCphServerOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CPH Server")
	}

	getCphServerRespBody, err := utils.FlattenResponse(getCphServerResp)
	if err != nil {
		return diag.FromErr(err)
	}
	statusRaw := utils.PathSearch("status", getCphServerRespBody, nil)

	if fmt.Sprint(statusRaw) == "6" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving CPH Server")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("server_name", getCphServerRespBody, nil)),
		d.Set("server_flavor", utils.PathSearch("server_model_name", getCphServerRespBody, nil)),
		d.Set("phone_flavor", utils.PathSearch("phone_model_name", getCphServerRespBody, nil)),
		d.Set("keypair_name", utils.PathSearch("keypair_name", getCphServerRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", getCphServerRespBody, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", getCphServerRespBody, nil)),
		d.Set("order_id", utils.PathSearch("metadata.order_id", getCphServerRespBody, nil)),
		d.Set("addresses", flattenGetCphServerResponseBodyAddress(getCphServerRespBody)),
		d.Set("bandwidth", flattenGetCphServerResponseBodyBandWidth(getCphServerRespBody)),
		d.Set("availability_zone", utils.PathSearch("availability_zone", getCphServerRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", getCphServerRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getCphServerRespBody, nil)),
		d.Set("security_groups", utils.PathSearch("security_groups", getCphServerRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetCphServerResponseBodyAddress(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("addresses", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"server_ip": utils.PathSearch("server_ip", v, nil),
			"public_ip": utils.PathSearch("public_ip", v, nil),
		})
	}
	return rst
}

func flattenGetCphServerResponseBodyBandWidth(resp interface{}) []interface{} {
	var rst []interface{}
	curJson, err := jmespath.Search("band_widths[0]", resp)
	if err != nil {
		log.Printf("[ERROR] error parsing band_widths from response= %#v", resp)
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"share_type":  fmt.Sprint(utils.PathSearch("band_width_share_type", curJson, nil)),
			"id":          utils.PathSearch("band_width_id", curJson, nil),
			"size":        utils.PathSearch("band_width_size", curJson, nil),
			"charge_mode": fmt.Sprint(utils.PathSearch("band_width_charge_mode", curJson, nil)),
		},
	}
	return rst
}

func resourceCphServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateCphServerNameChanges := []string{
		"name",
	}

	if d.HasChanges(updateCphServerNameChanges...) {
		// updateCphServerName: update CPH server name
		var (
			updateCphServerNameHttpUrl = "v1/{project_id}/cloud-phone/servers/{server_id}"
			updateCphServerNameProduct = "cph"
		)
		updateCphServerNameClient, err := cfg.NewServiceClient(updateCphServerNameProduct, region)
		if err != nil {
			return diag.Errorf("error creating CPH Client: %s", err)
		}

		updateCphServerNamePath := updateCphServerNameClient.Endpoint + updateCphServerNameHttpUrl
		updateCphServerNamePath = strings.ReplaceAll(updateCphServerNamePath, "{project_id}", updateCphServerNameClient.ProjectID)
		updateCphServerNamePath = strings.ReplaceAll(updateCphServerNamePath, "{server_id}", d.Id())

		updateCphServerNameOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateCphServerNameOpt.JSONBody = utils.RemoveNil(buildUpdateCphServerNameBodyParams(d, cfg))
		_, err = updateCphServerNameClient.Request("PUT", updateCphServerNamePath, &updateCphServerNameOpt)
		if err != nil {
			return diag.Errorf("error updating CPH Server: %s", err)
		}
	}
	return resourceCphServerRead(ctx, d, meta)
}

func buildUpdateCphServerNameBodyParams(d *schema.ResourceData, _ *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"server_name": utils.ValueIgnoreEmpty(d.Get("name")),
	}
	return bodyParams
}

func resourceCphServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	if err := common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
		return diag.Errorf("error unsubscribing CPH server: %s", err)
	}

	err := deleteCphServerWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the delete of CPH server (%s) to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteCphServerWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			var (
				getCphServerHttpUrl = "v1/{project_id}/cloud-phone/servers/{server_id}"
				getCphServerProduct = "cph"
			)
			getCphServerClient, err := cfg.NewServiceClient(getCphServerProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating CPH Client: %s", err)
			}

			getCphServerPath := getCphServerClient.Endpoint + getCphServerHttpUrl
			getCphServerPath = strings.ReplaceAll(getCphServerPath, "{project_id}", getCphServerClient.ProjectID)
			getCphServerPath = strings.ReplaceAll(getCphServerPath, "{server_id}", d.Id())

			getCphServerOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			getCphServerResp, err := getCphServerClient.Request("GET", getCphServerPath, &getCphServerOpt)

			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return getCphServerResp, "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			getCphServerRespBody, err := utils.FlattenResponse(getCphServerResp)
			if err != nil {
				return nil, "ERROR", err
			}

			statusRaw, err := jmespath.Search(`status`, getCphServerRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"6",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return getCphServerRespBody, "COMPLETED", nil
			}

			return getCphServerResp, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
