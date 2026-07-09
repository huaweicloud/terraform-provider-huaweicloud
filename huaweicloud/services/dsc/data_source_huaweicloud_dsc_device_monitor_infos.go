package dsc

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v1/{project_id}/devices/monitor-info
func DataSourceDscDeviceMonitorInfos() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscDeviceMonitorInfosRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"monitor_infos": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The device monitor information list.",
				Elem:        dscMonitorInfoSchema(),
			},
		},
	}
}

func dscMonitorInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device name.",
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device IP address.",
			},
			"type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The device type.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device status.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device description.",
			},
			"license_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The license information.",
				Elem:        dscLicenseInfoSchema(),
			},
			"os_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VM resource usage information.",
				Elem:        dscOsInfoSchema(),
			},
		},
	}
}

func dscLicenseInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"license": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The license information.",
			},
			"license_start": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The license effective time.",
			},
			"license_end": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The license expiration time.",
			},
		},
	}
}

func dscOsInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cpu": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The number of CPU cores.",
			},
			"disk": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The disk size.",
			},
			"mem": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The memory size.",
			},
		},
	}
}

func dataSourceDscDeviceMonitorInfosRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/devices/monitor-info"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC device monitor infos: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	monitorInfos := utils.PathSearch("monitor_infos", respBody, make([]interface{}, 0)).([]interface{})

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("monitor_infos", flattenDscMonitorInfos(monitorInfos)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscMonitorInfos(monitorInfos []interface{}) []interface{} {
	if len(monitorInfos) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(monitorInfos))
	for _, v := range monitorInfos {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"ip":          utils.PathSearch("ip", v, nil),
			"type":        utils.PathSearch("type", v, nil),
			"status":      utils.PathSearch("status", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"license_info": flattenDscLicenseInfo(
				utils.PathSearch("license_info", v, nil)),
			"os_info": flattenDscOsInfo(
				utils.PathSearch("os_info", v, nil)),
		})
	}

	return rst
}

func flattenDscLicenseInfo(licenseInfo interface{}) []interface{} {
	if licenseInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"license":       utils.PathSearch("license", licenseInfo, nil),
			"license_start": utils.PathSearch("license_start", licenseInfo, nil),
			"license_end":   utils.PathSearch("license_end", licenseInfo, nil),
		},
	}
}

func flattenDscOsInfo(osInfo interface{}) []interface{} {
	if osInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"cpu":  utils.PathSearch("cpu", osInfo, nil),
			"disk": utils.PathSearch("disk", osInfo, nil),
			"mem":  utils.PathSearch("mem", osInfo, nil),
		},
	}
}
