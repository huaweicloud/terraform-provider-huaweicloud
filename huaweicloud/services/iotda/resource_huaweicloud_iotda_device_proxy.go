package iotda

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

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

func buildEffectiveTimeRangeBodyParams(d *schema.ResourceData) *model.EffectiveTimeRange {
	effectiveTimeRangeMap, ok := d.Get("effective_time_range").([]interface{})[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return &model.EffectiveTimeRange{
		StartTime: utils.String(effectiveTimeRangeMap["start_time"].(string)),
		EndTime:   utils.String(effectiveTimeRangeMap["end_time"].(string)),
	}
}

func buildDeviceProxyCreateParams(d *schema.ResourceData) *model.CreateDeviceProxyRequest {
	createOptsBody := model.AddDeviceProxy{
		ProxyName:          d.Get("name").(string),
		ProxyDevices:       utils.ExpandToStringList(d.Get("devices").(*schema.Set).List()),
		EffectiveTimeRange: buildEffectiveTimeRangeBodyParams(d),
		AppId:              d.Get("space_id").(string),
	}

	return &model.CreateDeviceProxyRequest{
		Body: &createOptsBody,
	}
}

func resourceDeviceProxyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	createOpts := buildDeviceProxyCreateParams(d)
	resp, err := client.CreateDeviceProxy(createOpts)
	if err != nil {
		return diag.Errorf("error creating IoTDA device proxy: %s", err)
	}

	if resp == nil || resp.ProxyId == nil {
		return diag.Errorf("error creating IoTDA device proxy: ID is not found in API response")
	}

	d.SetId(*resp.ProxyId)

	return resourceDeviceProxyRead(ctx, d, meta)
}

func buildDeviceProxyQueryParams(d *schema.ResourceData) *model.ShowDeviceProxyRequest {
	return &model.ShowDeviceProxyRequest{
		ProxyId: d.Id(),
	}
}

func resourceDeviceProxyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	response, err := client.ShowDeviceProxy(buildDeviceProxyQueryParams(d))
	// When the resource does not exist, query API will return `404` error code.
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA device proxy")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("space_id", response.AppId),
		d.Set("name", response.ProxyName),
		d.Set("devices", response.ProxyDevices),
		d.Set("effective_time_range", flattenEffectiveTimeRange(response.EffectiveTimeRange)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEffectiveTimeRange(effectiveTimeRangeResp *model.EffectiveTimeRangeResponseDto) []map[string]interface{} {
	if effectiveTimeRangeResp == nil {
		return nil
	}

	return []map[string]interface{}{{
		"start_time": effectiveTimeRangeResp.StartTime,
		"end_time":   effectiveTimeRangeResp.EndTime,
	}}
}

func buildDeviceProxyUpdateParams(d *schema.ResourceData) *model.UpdateDeviceProxyRequest {
	updateRequest := model.UpdateDeviceProxyRequest{
		ProxyId: d.Id(),
		Body: &model.UpdateDeviceProxy{
			ProxyName:          utils.String(d.Get("name").(string)),
			ProxyDevices:       utils.ExpandToStringListPointer(d.Get("devices").(*schema.Set).List()),
			EffectiveTimeRange: buildEffectiveTimeRangeBodyParams(d),
		},
	}

	return &updateRequest
}

func resourceDeviceProxyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	updateOpts := buildDeviceProxyUpdateParams(d)
	_, err = client.UpdateDeviceProxy(updateOpts)
	if err != nil {
		return diag.Errorf("error updating IoTDA device proxy: %s", err)
	}

	return resourceDeviceProxyRead(ctx, d, meta)
}

func buildDeviceProxyDeleteParams(d *schema.ResourceData) *model.DeleteDeviceProxyRequest {
	return &model.DeleteDeviceProxyRequest{
		ProxyId: d.Id(),
	}
}

func resourceDeviceProxyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	_, err = client.DeleteDeviceProxy(buildDeviceProxyDeleteParams(d))
	// When the resource does not exist, delete API will return `404` error code.
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IoTDA device proxy")
	}

	return nil
}
