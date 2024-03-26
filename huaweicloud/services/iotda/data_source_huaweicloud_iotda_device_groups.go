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

// @API IoTDA GET /v5/iot/{project_id}/device-group
func DataSourceDeviceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeviceGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDeviceGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	var (
		allDeviceGroups []model.DeviceGroupResponseSummary
		limit           = int32(50)
		offset          int32
	)

	for {
		listOpts := model.ListDeviceGroupsRequest{
			Name:      utils.StringIgnoreEmpty(d.Get("name").(string)),
			GroupType: utils.StringIgnoreEmpty(d.Get("type").(string)),
			AppId:     utils.StringIgnoreEmpty(d.Get("space_id").(string)),
			Limit:     utils.Int32(limit),
			Offset:    &offset,
		}

		listResp, listErr := client.ListDeviceGroups(&listOpts)
		if listErr != nil {
			return diag.Errorf("error querying IoTDA device groups: %s", listErr)
		}

		if len(*listResp.DeviceGroups) == 0 {
			break
		}

		allDeviceGroups = append(allDeviceGroups, *listResp.DeviceGroups...)
		offset += limit
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuId)

	targetDeviceGroups := filterListDeviceGroups(allDeviceGroups, d)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("groups", flattenDeviceGroups(targetDeviceGroups)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterListDeviceGroups(deviceGroups []model.DeviceGroupResponseSummary, d *schema.ResourceData) []model.DeviceGroupResponseSummary {
	if len(deviceGroups) == 0 {
		return nil
	}

	rst := make([]model.DeviceGroupResponseSummary, 0, len(deviceGroups))
	for _, v := range deviceGroups {
		if deviceGroupId, ok := d.GetOk("group_id"); ok &&
			fmt.Sprint(deviceGroupId) != utils.StringValue(v.GroupId) {
			continue
		}

		if parentGroupId, ok := d.GetOk("parent_group_id"); ok &&
			fmt.Sprint(parentGroupId) != utils.StringValue(v.SuperGroupId) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenDeviceGroups(deviceGroups []model.DeviceGroupResponseSummary) []interface{} {
	if len(deviceGroups) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(deviceGroups))
	for _, v := range deviceGroups {
		rst = append(rst, map[string]interface{}{
			"id":              v.GroupId,
			"name":            v.Name,
			"description":     v.Description,
			"parent_group_id": v.SuperGroupId,
			"type":            v.GroupType,
		})
	}

	return rst
}
