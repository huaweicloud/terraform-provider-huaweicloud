package aad

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var aadInstanceNonUpdatableParams = []string{
	"ip_type",
	"resource_region",
	"instance_access_type",
	"duration",
	"amount",
	"instance_name",
	"period_type",
	"protection_package",
	"protected_domain",
	"forwarding_rule",
	"enterprise_project_id",
}

// Currently, AAD instance do not support deletion or unsubscribing, and the purchase amount is relatively large,
// so the resource code has not been fully tested, also do not provide test case files. Only the import and partial
// parameter modification functions were tested locally.

// @API AAD POST /v2/aad/instance
// @API AAD PUT /v2/aad/instance
// @API AAD GET /v2/aad/instances/{instance_id}
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/subscriptions/resources/query
func ResourceAadInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAadInstanceCreate,
		ReadContext:   resourceAadInstanceRead,
		UpdateContext: resourceAadInstanceUpdate,
		DeleteContext: resourceAadInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(aadInstanceNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			// Query API no return.
			"ip_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			// Query API no return.
			"resource_region": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Query API no return.
			"instance_access_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Query API no return.
			"duration": {
				Type:     schema.TypeInt,
				Required: true,
			},
			// Query API no return.
			"amount": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Query API no return.
			"period_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			// Can be updated.
			"service_bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
			},
			// Can be updated.
			"basic_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			// Can be updated.
			"elastic_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			// Can be updated. Query API no return.
			"basic_qps": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			// Can be updated.
			"elastic_service_bandwidth_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			// Can be updated.
			"elastic_service_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			// Query API no return.
			"protection_package": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Query API no return.
			"protected_domain": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			// Query API no return.
			"forwarding_rule": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Can be updated, update after creation is completed.
			"port_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			// Can be updated, update after creation is completed.
			"bind_domain_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"current_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"product_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isp_spec": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_center": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"spec_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"main_resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"main_resource_spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"main_resource_product_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     aadInstanceConfigSchema(),
			},
			"elastic_service_bw_update_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func aadInstanceConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"freeze_type": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func buildCreateInstanceBodyParams(d *schema.ResourceData, epsId string) map[string]interface{} {
	return map[string]interface{}{
		"ip_type":                        d.Get("ip_type"),
		"resource_region":                d.Get("resource_region"),
		"instance_access_type":           d.Get("instance_access_type"),
		"duration":                       d.Get("duration"),
		"amount":                         d.Get("amount"),
		"instance_name":                  d.Get("instance_name"),
		"period_type":                    d.Get("period_type"),
		"service_bandwidth":              d.Get("service_bandwidth"),
		"basic_bandwidth":                utils.ValueIgnoreEmpty(d.Get("basic_bandwidth")),
		"elastic_bandwidth":              utils.ValueIgnoreEmpty(d.Get("elastic_bandwidth")),
		"basic_qps":                      utils.ValueIgnoreEmpty(d.Get("basic_qps")),
		"elastic_service_bandwidth_type": utils.ValueIgnoreEmpty(d.Get("elastic_service_bandwidth_type")),
		"elastic_service_bandwidth":      utils.ValueIgnoreEmpty(d.Get("elastic_service_bandwidth")),
		"protection_package":             utils.ValueIgnoreEmpty(d.Get("protection_package")),
		"protected_domain":               utils.ValueIgnoreEmpty(d.Get("protected_domain")),
		"forwarding_rule":                utils.ValueIgnoreEmpty(d.Get("forwarding_rule")),
		"enterprise_project_id":          utils.ValueIgnoreEmpty(epsId),
	}
}

func resourceAadInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/aad/instance"
		product = "aad"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildCreateInstanceBodyParams(d, epsId)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating AAD instance: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	orderID := utils.PathSearch("order_id", respBody, "").(string)
	if orderID == "" {
		return diag.Errorf("error creating AAD instance: order ID is not found in API response")
	}

	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}

	if err = common.WaitOrderComplete(ctx, bssClient, orderID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	instanceID, err := common.WaitOrderResourceComplete(ctx, bssClient, orderID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for create AAD instance order (%s) complete: %s", orderID, err)
	}

	d.SetId(instanceID)

	portNum := d.Get("port_num").(int)
	bindDomainNum := d.Get("bind_domain_num").(int)
	if portNum != 0 || bindDomainNum != 0 {
		if err = updateAadInstance(ctx, d, meta); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceAadInstanceRead(ctx, d, meta)
}

func resourceAadInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v2/aad/instances/{instance_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err,
			"error_code", "AAD.10010035"), "error retrieving AAD instance")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_name", utils.PathSearch("instance_name", respBody, nil)),
		d.Set("service_bandwidth", utils.PathSearch(
			"product_spec_data.service_bandwidth", respBody, nil)),
		d.Set("basic_bandwidth", utils.PathSearch(
			"product_spec_data.basic_bandwidth", respBody, nil)),
		d.Set("elastic_bandwidth", utils.PathSearch(
			"product_spec_data.elastic_bandwidth", respBody, nil)),
		d.Set("elastic_service_bandwidth_type", utils.PathSearch(
			"product_spec_data.elastic_service_bandwidth_type", respBody, nil)),
		d.Set("elastic_service_bandwidth", utils.PathSearch(
			"product_spec_data.elastic_service_bandwidth", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch(
			"enterprise_project_id", respBody, nil)),
		d.Set("port_num", utils.PathSearch(
			"product_spec_data.port_num", respBody, nil)),
		d.Set("bind_domain_num", utils.PathSearch(
			"product_spec_data.bind_domain_num", respBody, nil)),
		d.Set("version", utils.PathSearch("version", respBody, nil)),
		d.Set("expire_time", utils.PathSearch("expire_time", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("current_time", utils.PathSearch("current_time", respBody, nil)),
		d.Set("product_uuid", utils.PathSearch(
			"product_spec_data.product_uuid", respBody, nil)),
		d.Set("isp_spec", utils.PathSearch(
			"product_spec_data.isp_spec", respBody, nil)),
		d.Set("data_center", utils.PathSearch(
			"product_spec_data.data_center", respBody, nil)),
		d.Set("spec_type", utils.PathSearch(
			"product_spec_data.spec_type", respBody, nil)),
		d.Set("main_resource_type", utils.PathSearch(
			"product_spec_data.main_resource_type", respBody, nil)),
		d.Set("main_resource_spec_code", utils.PathSearch(
			"product_spec_data.main_resource_spec_code", respBody, nil)),
		d.Set("main_resource_product_id", utils.PathSearch(
			"product_spec_data.main_resource_product_id", respBody, nil)),
		d.Set("instance_config", flattenAadInstanceConfig(utils.PathSearch(
			"instance_config", respBody, nil))),
		d.Set("elastic_service_bw_update_enable", utils.PathSearch(
			"elastic_service_bw_update_enable", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAadInstanceConfig(raw interface{}) []interface{} {
	if raw == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"tags": utils.ExpandToStringList(
				utils.PathSearch("tags", raw, make([]interface{}, 0)).([]interface{})),
			"freeze_type": utils.ExpandToIntList(
				utils.PathSearch("freeze_type", raw, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func buildUpdateInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	upgradeData := map[string]interface{}{
		"basic_bandwidth":                utils.ValueIgnoreEmpty(d.Get("basic_bandwidth")),
		"elastic_bandwidth":              utils.ValueIgnoreEmpty(d.Get("elastic_bandwidth")),
		"service_bandwidth":              utils.ValueIgnoreEmpty(d.Get("service_bandwidth")),
		"port_num":                       utils.ValueIgnoreEmpty(d.Get("forwarding_rule")),
		"bind_domain_num":                utils.ValueIgnoreEmpty(d.Get("protected_domain")),
		"elastic_service_bandwidth_type": utils.ValueIgnoreEmpty(d.Get("elastic_service_bandwidth_type")),
		"elastic_service_bandwidth":      utils.ValueIgnoreEmpty(d.Get("elastic_service_bandwidth")),
		"basic_qps":                      utils.ValueIgnoreEmpty(d.Get("basic_qps")),
	}

	return map[string]interface{}{
		"instance_id":  d.Id(),
		"upgrade_data": utils.RemoveNil(upgradeData),
	}
}

func updateAadInstance(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/aad/instance"
		product = "aad"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return fmt.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildUpdateInstanceBodyParams(d),
	}

	resp, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error updating AAD instance: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	orderID := utils.PathSearch("order_id", respBody, "").(string)
	if orderID != "" {
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating BSS v2 client: %s", err)
		}
		if err = common.WaitOrderComplete(ctx, bssClient, orderID, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
		if _, err = common.WaitOrderResourceComplete(ctx, bssClient, orderID, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("error waiting for update AAD instance order (%s) complete: %s", orderID, err)
		}
	}

	return nil
}

func resourceAadInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := updateAadInstance(ctx, d, meta); err != nil {
		return diag.FromErr(err)
	}

	return resourceAadInstanceRead(ctx, d, meta)
}

func resourceAadInstanceDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This AAD instances do not support deletion and unsubscribing. Deleting this resource will not delete or
    unsubscribe AAD instance, but will only remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
