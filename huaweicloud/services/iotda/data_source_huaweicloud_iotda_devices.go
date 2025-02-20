package iotda

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA GET /v5/iot/{project_id}/devices
func DataSourceDevices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDevicesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_cascade": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"device_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"devices": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"space_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"space_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fw_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sw_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sdk_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func buildDevicesQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=50"
	if v, ok := d.GetOk("space_id"); ok {
		queryParams = fmt.Sprintf("%s&app_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("product_id"); ok {
		queryParams = fmt.Sprintf("%s&product_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("gateway_id"); ok {
		queryParams = fmt.Sprintf("%s&gateway_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("is_cascade"); ok {
		queryParams = fmt.Sprintf("%s&is_cascade_query=%v", queryParams, v)
	}
	if v, ok := d.GetOk("node_id"); ok {
		queryParams = fmt.Sprintf("%s&node_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&device_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("start_time"); ok {
		queryParams = fmt.Sprintf("%s&start_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		queryParams = fmt.Sprintf("%s&end_time=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceDevicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
		httpUrl   = "v5/iot/{project_id}/devices"
		offset    = 0
		result    = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildDevicesQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving IoTDA devices: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		devicesResp := utils.PathSearch("devices", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(devicesResp) == 0 {
			break
		}

		result = append(result, devicesResp...)
		offset += len(devicesResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("devices", flattenDevices(filterListDevices(result, d))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterListDevices(devicesResp []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(devicesResp))
	for _, v := range devicesResp {
		if deviceID, ok := d.GetOk("device_id"); ok &&
			fmt.Sprint(deviceID) != utils.PathSearch("device_id", v, "").(string) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenDevices(devicesResp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(devicesResp))
	for _, v := range devicesResp {
		rst = append(rst, map[string]interface{}{
			"space_id":     utils.PathSearch("app_id", v, nil),
			"space_name":   utils.PathSearch("app_name", v, nil),
			"product_id":   utils.PathSearch("product_id", v, nil),
			"product_name": utils.PathSearch("product_name", v, nil),
			"gateway_id":   utils.PathSearch("gateway_id", v, nil),
			"id":           utils.PathSearch("device_id", v, nil),
			"name":         utils.PathSearch("device_name", v, nil),
			"node_id":      utils.PathSearch("node_id", v, nil),
			"node_type":    utils.PathSearch("node_type", v, nil),
			"description":  utils.PathSearch("description", v, nil),
			"status":       utils.PathSearch("status", v, nil),
			"fw_version":   utils.PathSearch("fw_version", v, nil),
			"sw_version":   utils.PathSearch("sw_version", v, nil),
			"sdk_version":  utils.PathSearch("device_sdk_version", v, nil),
			"tags":         flattenDeviceTags(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}
