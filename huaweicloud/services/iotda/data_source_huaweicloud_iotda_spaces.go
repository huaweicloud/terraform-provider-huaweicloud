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

// @API IoTDA GET /v5/iot/{project_id}/apps
func DataSourceSpaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSpacesRead,

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
			"space_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_default": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"spaces": {
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
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSpacesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	listOpts := &model.ShowApplicationsRequest{}
	isDefault := d.Get("is_default").(string)
	if isDefault == "true" || isDefault == "false" {
		listOpts.DefaultApp = utils.StringToBool(isDefault)
	}

	listResp, listErr := client.ShowApplications(listOpts)
	if listErr != nil {
		return diag.Errorf("error querying IoTDA spaces: %s", listErr)
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuId)

	targetSpaces := filterListSpaces(*listResp.Applications, d)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("spaces", flattenSpaces(targetSpaces)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterListSpaces(spaces []model.ApplicationDto, d *schema.ResourceData) []model.ApplicationDto {
	if len(spaces) == 0 {
		return nil
	}

	rst := make([]model.ApplicationDto, 0, len(spaces))
	for _, v := range spaces {
		if spaceID, ok := d.GetOk("space_id"); ok &&
			fmt.Sprint(spaceID) != utils.StringValue(v.AppId) {
			continue
		}

		if spaceName, ok := d.GetOk("space_name"); ok &&
			fmt.Sprint(spaceName) != utils.StringValue(v.AppName) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenSpaces(spaces []model.ApplicationDto) []interface{} {
	if len(spaces) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(spaces))
	for _, v := range spaces {
		rst = append(rst, map[string]interface{}{
			"id":         v.AppId,
			"name":       v.AppName,
			"created_at": v.CreateTime,
			"is_default": v.DefaultApp,
		})
	}

	return rst
}
