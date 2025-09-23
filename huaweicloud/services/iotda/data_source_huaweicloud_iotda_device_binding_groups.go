package iotda

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
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
		product   = "iotda"
		httpUrl   = "v5/iot/{project_id}/devices/{device_id}/list-device-group"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{device_id}", d.Get("device_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving IoTDA device groups: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	deviceGroupsResp := utils.PathSearch("device_groups", getRespBody, make([]interface{}, 0)).([]interface{})
	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("groups", flattenListDeviceGroups(deviceGroupsResp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListDeviceGroups(deviceGroupsResp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(deviceGroupsResp))
	for _, v := range deviceGroupsResp {
		rst = append(rst, map[string]interface{}{
			"id":              utils.PathSearch("group_id", v, nil),
			"name":            utils.PathSearch("name", v, nil),
			"description":     utils.PathSearch("description", v, nil),
			"parent_group_id": utils.PathSearch("super_group_id", v, nil),
			"type":            utils.PathSearch("group_type", v, nil),
		})
	}

	return rst
}
