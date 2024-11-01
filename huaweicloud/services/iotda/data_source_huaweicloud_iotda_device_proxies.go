package iotda

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA GET /v5/iot/{project_id}/device-proxies
func DataSourceDeviceProxies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeviceProxiesRead,

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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proxies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"space_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"effective_time_range": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceDeviceProxiesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)

	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	var (
		allProxies []model.QueryDeviceProxySimplify
		limit      = int32(50)
		offset     int32
	)

	for {
		listOpts := model.ListDeviceProxiesRequest{
			AppId:     utils.StringIgnoreEmpty(d.Get("space_id").(string)),
			ProxyName: utils.StringIgnoreEmpty(d.Get("name").(string)),
			Limit:     utils.Int32(limit),
			Offset:    &offset,
		}

		listResp, listErr := client.ListDeviceProxies(&listOpts)
		if listErr != nil {
			return diag.Errorf("error querying IoTDA device proxies: %s", listErr)
		}

		if listResp == nil || listResp.DeviceProxies == nil {
			break
		}

		if len(*listResp.DeviceProxies) == 0 {
			break
		}

		allProxies = append(allProxies, *listResp.DeviceProxies...)
		//nolint:gosec
		offset += int32(len(*listResp.DeviceProxies))
	}

	uuID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("proxies", flattenDeviceProxies(allProxies)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDeviceProxies(proxies []model.QueryDeviceProxySimplify) []interface{} {
	if len(proxies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(proxies))
	for _, v := range proxies {
		rst = append(rst, map[string]interface{}{
			"id":                   v.ProxyId,
			"name":                 v.ProxyName,
			"space_id":             v.AppId,
			"effective_time_range": flattenEffectiveTimeRange(v.EffectiveTimeRange),
		})
	}

	return rst
}
