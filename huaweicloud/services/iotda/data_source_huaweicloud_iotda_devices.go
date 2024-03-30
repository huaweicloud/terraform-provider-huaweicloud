package iotda

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

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

func dataSourceDevicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)

	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	var (
		allDevices []model.QueryDeviceSimplify
		limit      = int32(50)
		offset     int32
	)

	for {
		listOpts := model.ListDevicesRequest{
			AppId:          utils.StringIgnoreEmpty(d.Get("space_id").(string)),
			ProductId:      utils.StringIgnoreEmpty(d.Get("product_id").(string)),
			GatewayId:      utils.StringIgnoreEmpty(d.Get("gateway_id").(string)),
			IsCascadeQuery: utils.Bool(d.Get("is_cascade").(bool)),
			NodeId:         utils.StringIgnoreEmpty(d.Get("node_id").(string)),
			DeviceName:     utils.StringIgnoreEmpty(d.Get("name").(string)),
			StartTime:      utils.StringIgnoreEmpty(d.Get("start_time").(string)),
			EndTime:        utils.StringIgnoreEmpty(d.Get("end_time").(string)),
			Limit:          utils.Int32(limit),
			Offset:         &offset,
		}

		listResp, listErr := client.ListDevices(&listOpts)
		if listErr != nil {
			return diag.Errorf("error querying IoTDA devices: %s", listErr)
		}

		if len(*listResp.Devices) == 0 {
			break
		}

		allDevices = append(allDevices, *listResp.Devices...)
		offset += limit
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuId)

	targetDevices := filterListDevices(allDevices, d)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("devices", flattenDevices(targetDevices)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterListDevices(devices []model.QueryDeviceSimplify, d *schema.ResourceData) []model.QueryDeviceSimplify {
	if len(devices) == 0 {
		return nil
	}

	rst := make([]model.QueryDeviceSimplify, 0, len(devices))
	for _, v := range devices {
		if deviceID, ok := d.GetOk("device_id"); ok &&
			fmt.Sprint(deviceID) != utils.StringValue(v.DeviceId) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenDevices(devices []model.QueryDeviceSimplify) []interface{} {
	if len(devices) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(devices))
	for _, v := range devices {
		rst = append(rst, map[string]interface{}{
			"space_id":     v.AppId,
			"space_name":   v.AppName,
			"product_id":   v.ProductId,
			"product_name": v.ProductName,
			"gateway_id":   v.GatewayId,
			"id":           v.DeviceId,
			"name":         v.DeviceName,
			"node_id":      v.NodeId,
			"node_type":    v.NodeType,
			"description":  v.Description,
			"status":       v.Status,
			"fw_version":   v.FwVersion,
			"sw_version":   v.SwVersion,
			"sdk_version":  v.DeviceSdkVersion,
			"tags":         flattenTags(v.Tags),
		})
	}

	return rst
}
