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
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/apps"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	isDefault := d.Get("is_default").(string)
	if isDefault == "true" || isDefault == "false" {
		listPath = fmt.Sprintf("%s?default_app=%s", listPath, isDefault)
	}

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving IoTDA spaces: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	listSpaceArray := utils.PathSearch("applications", listRespBody, make([]interface{}, 0)).([]interface{})

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("spaces", flattenListSpaces(filterListSpaces(listSpaceArray, d))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterListSpaces(spaces []interface{}, d *schema.ResourceData) []interface{} {
	if len(spaces) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(spaces))
	for _, v := range spaces {
		if spaceId, ok := d.GetOk("space_id"); ok &&
			fmt.Sprint(spaceId) != utils.PathSearch("app_id", v, "").(string) {
			continue
		}

		if spaceName, ok := d.GetOk("space_name"); ok &&
			fmt.Sprint(spaceName) != utils.PathSearch("app_name", v, "").(string) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenListSpaces(spaces []interface{}) []interface{} {
	if len(spaces) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(spaces))
	for _, v := range spaces {
		result = append(result, map[string]interface{}{
			"id":         utils.PathSearch("app_id", v, nil),
			"name":       utils.PathSearch("app_name", v, nil),
			"created_at": utils.PathSearch("create_time", v, nil),
			"is_default": utils.PathSearch("default_app", v, nil),
		})
	}
	return result
}
