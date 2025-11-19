package workspace

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

// @API Workspace GET /v2/{project_id}/common/timezones
func DataSourceTimeZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTimeZonesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the time zones are located.`,
			},

			// Attributes.
			"time_zones": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of time zones.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the time zone.`,
						},
						"offset": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The offset of the time zone.`,
						},
						"us_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The English description of the time zone.`,
						},
						"cn_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The Chinese description of the time zone.`,
						},
					},
				},
			},
		},
	}
}

// Currently, the filter parameter time_zone_name is not available.
func listTimeZones(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/common/timezones"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("time_zones", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func dataSourceTimeZonesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	resp, err := listTimeZones(client)
	if err != nil {
		return diag.Errorf("error querying Workspace time zones: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("time_zones", flattenTimeZones(resp)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTimeZones(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"name":           utils.PathSearch("time_zone_name", item, nil),
			"offset":         utils.PathSearch("time_zone", item, nil),
			"us_description": utils.PathSearch("time_zone_desc_us", item, nil),
			"cn_description": utils.PathSearch("time_zone_desc_cn", item, nil),
		})
	}

	return result
}
