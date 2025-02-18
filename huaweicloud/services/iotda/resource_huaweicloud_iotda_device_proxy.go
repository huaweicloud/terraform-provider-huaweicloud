package iotda

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA POST /v5/iot/{project_id}/device-proxies
// @API IoTDA GET /v5/iot/{project_id}/device-proxies/{proxy_id}
// @API IoTDA PUT /v5/iot/{project_id}/device-proxies/{proxy_id}
// @API IoTDA DELETE /v5/iot/{project_id}/device-proxies/{proxy_id}
func ResourceDeviceProxy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeviceProxyCreate,
		ReadContext:   resourceDeviceProxyRead,
		UpdateContext: resourceDeviceProxyUpdate,
		DeleteContext: resourceDeviceProxyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"devices": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"effective_time_range": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func buildEffectiveTimeRangeBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawEffectiveTimeRange, ok := d.Get("effective_time_range").([]interface{})
	if !ok {
		return nil
	}

	if len(rawEffectiveTimeRange) < 1 {
		return nil
	}

	effectiveTimeRange := rawEffectiveTimeRange[0].(map[string]interface{})
	effectiveTimeRangeParams := map[string]interface{}{
		"start_time": effectiveTimeRange["start_time"],
		"end_time":   effectiveTimeRange["end_time"],
	}

	return effectiveTimeRangeParams
}

func buildDeviceProxyCreateParams(d *schema.ResourceData) map[string]interface{} {
	proxyParams := map[string]interface{}{
		"proxy_name":           d.Get("name"),
		"proxy_devices":        utils.ExpandToStringList(d.Get("devices").(*schema.Set).List()),
		"effective_time_range": buildEffectiveTimeRangeBodyParams(d),
		"app_id":               d.Get("space_id"),
	}

	return proxyParams
}

func resourceDeviceProxyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/device-proxies"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildDeviceProxyCreateParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating IoTDA device proxy: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	proxyId := utils.PathSearch("proxy_id", respBody, "").(string)
	if proxyId == "" {
		return diag.Errorf("error creating IoTDA device proxy: ID is not found in API response")
	}

	d.SetId(proxyId)

	return resourceDeviceProxyRead(ctx, d, meta)
}

func resourceDeviceProxyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/device-proxies/{proxy_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{proxy_id}", d.Id())
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		// When the resource does not exist, query API will return `404` error code.
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA device proxy")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("space_id", utils.PathSearch("app_id", getRespBody, nil)),
		d.Set("name", utils.PathSearch("proxy_name", getRespBody, nil)),
		d.Set("devices", utils.PathSearch("proxy_devices", getRespBody, make([]interface{}, 0))),
		d.Set("effective_time_range", flattenEffectiveTimeRange(utils.PathSearch("effective_time_range", getRespBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEffectiveTimeRange(effectiveTimeRangeResp interface{}) []map[string]interface{} {
	if effectiveTimeRangeResp == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"start_time": utils.PathSearch("start_time", effectiveTimeRangeResp, nil),
			"end_time":   utils.PathSearch("end_time", effectiveTimeRangeResp, nil),
		},
	}
}

func buildDeviceProxyUpdateParams(d *schema.ResourceData) map[string]interface{} {
	updateProxyParams := map[string]interface{}{
		"proxy_name":           d.Get("name"),
		"proxy_devices":        utils.ExpandToStringList(d.Get("devices").(*schema.Set).List()),
		"effective_time_range": buildEffectiveTimeRangeBodyParams(d),
	}
	return updateProxyParams
}

func resourceDeviceProxyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/device-proxies/{proxy_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{proxy_id}", d.Id())
	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildDeviceProxyUpdateParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return diag.Errorf("error updating IoTDA device proxy: %s", err)
	}

	return resourceDeviceProxyRead(ctx, d, meta)
}

func resourceDeviceProxyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/device-proxies/{proxy_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{proxy_id}", d.Id())
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		// When the resource does not exist, delete API will return `404` error code.
		return common.CheckDeletedDiag(d, err, "error deleting IoTDA device proxy")
	}

	return nil
}
