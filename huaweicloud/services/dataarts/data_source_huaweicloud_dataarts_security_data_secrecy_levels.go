package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/security/data-classification/secrecy-level
func DataSourceSecurityDataSecrecyLevels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityDataSecrecyLevelsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the data secrecy levels are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the data secrecy levels belong.`,
			},

			// Optional parameters.
			"order_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the field used to sort the data secrecy levels.`,
			},
			"desc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to sort the data secrecy levels in descending order.`,
			},

			// Attributes.
			"secrecy_levels": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of data secrecy levels that matched the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the data secrecy level.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the data secrecy level.`,
						},
						"level_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The level number of the data secrecy level.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the data secrecy level.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance ID to which the data secrecy level belongs.`,
						},
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the data secrecy level.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the data secrecy level, in RFC3339 format.`,
						},
						"updated_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The updater of the data secrecy level.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the data secrecy level, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

func buildSecurityDataSecrecyLevelsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("order_by"); ok {
		res = fmt.Sprintf("%s&order_by=%v", res, v)
	}
	if v, ok := d.GetOk("desc"); ok {
		res = fmt.Sprintf("%s&desc=%v", res, v)
	}

	return res
}

func listSecurityDataSecrecyLevels(client *golangsdk.ServiceClient, d *schema.ResourceData, workspaceId string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/security/data-classification/secrecy-level?limit={limit}"
		// limit: maximum is `100`, default is `10`.
		limit  = 100
		offset = 0
		result = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildSecurityDataSecrecyLevelsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		levels := utils.PathSearch("content", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, levels...)
		if len(levels) < limit {
			break
		}

		offset += len(levels)
	}

	return result, nil
}

func dataSourceSecurityDataSecrecyLevelsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	secrecyLevels, err := listSecurityDataSecrecyLevels(client, d, workspaceId)
	if err != nil {
		return diag.Errorf("error querying DataArts Security data secrecy levels: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("secrecy_levels", flattenSecurityDataSecrecyLevels(secrecyLevels)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSecurityDataSecrecyLevels(levels []interface{}) []interface{} {
	if len(levels) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(levels))
	for _, level := range levels {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("secrecy_level_id", level, nil),
			"name":         utils.PathSearch("secrecy_level_name", level, nil),
			"level_number": utils.PathSearch("secrecy_level_number", level, nil),
			"description":  utils.PathSearch("description", level, nil),
			"instance_id":  utils.PathSearch("instance_id", level, nil),
			"created_by":   utils.PathSearch("created_by", level, nil),
			"created_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("created_at",
				level, float64(0)).(float64))/1000, false),
			"updated_by": utils.PathSearch("updated_by", level, nil),
			"updated_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("updated_at",
				level, float64(0)).(float64))/1000, false),
		})
	}

	return result
}
