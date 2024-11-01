package iotda

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IoTDA POST /v5/iot/{project_id}/devices/{device_id}/list-device-group
func DataSourceDeviceBindingGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeviceBindingGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"device_id": {
				Type:     schema.TypeString,
				Required: true,
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

func dataSourceDeviceBindingGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	listOpts := model.ListDeviceGroupsByDeviceRequest{
		DeviceId: d.Get("device_id").(string),
	}

	allDeviceGroups := make([]model.ListDeviceGroupSummary, 0)
	listResp, listErr := client.ListDeviceGroupsByDevice(&listOpts)
	if listErr != nil {
		return diag.Errorf("error querying IoTDA device groups: %s", listErr)
	}

	if listResp != nil && listResp.DeviceGroups != nil {
		allDeviceGroups = *listResp.DeviceGroups
	}

	uuID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("groups", flattenListDeviceGroups(allDeviceGroups)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListDeviceGroups(deviceGroups []model.ListDeviceGroupSummary) []interface{} {
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
