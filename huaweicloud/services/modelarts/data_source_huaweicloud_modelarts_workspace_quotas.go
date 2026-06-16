package modelarts

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v1/{project_id}/workspaces/{workspace_id}/quotas
func DataSourceWorkspaceQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkspaceQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the workspace quotas are located.`,
			},

			// Required parameters
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to be queried.`,
			},

			// Attributes
			"quotas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of workspace quotas that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The unique identifier of the resource.`,
						},
						"quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The current quota value.`,
						},
						"min_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The minimum value allowed for the quota.`,
						},
						"max_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum value allowed for the quota.`,
						},
						"used_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The used quota value.`,
						},
						"name_cn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the quota in Chinese.`,
						},
						"name_en": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the quota in English.`,
						},
						"unit_cn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The unit of the quota in Chinese.`,
						},
						"unit_en": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The unit of the quota in English.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The last update time of the quota, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

func getWorkspaceQuotas(client *golangsdk.ServiceClient, workspaceId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/workspaces/{workspace_id}/quotas"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func flattenWorkspaceQuotas(quotas []interface{}) []map[string]interface{} {
	if len(quotas) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(quotas))
	for _, quota := range quotas {
		result = append(result, map[string]interface{}{
			"resource":   utils.PathSearch("resource", quota, nil),
			"quota":      utils.PathSearch("quota", quota, nil),
			"min_quota":  utils.PathSearch("min_quota", quota, nil),
			"max_quota":  utils.PathSearch("max_quota", quota, nil),
			"used_quota": utils.PathSearch("used_quota", quota, nil),
			"name_cn":    utils.PathSearch("name_cn", quota, nil),
			"name_en":    utils.PathSearch("name_en", quota, nil),
			"unit_cn":    utils.PathSearch("unit_cn", quota, nil),
			"unit_en":    utils.PathSearch("unit_en", quota, nil),
			"updated_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", quota,
				float64(0)).(float64))/1000, false),
		})
	}

	return result
}

func dataSourceWorkspaceQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	workspaceId := d.Get("workspace_id").(string)
	respBody, err := getWorkspaceQuotas(client, workspaceId)
	if err != nil {
		return diag.Errorf("error querying workspace quotas: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("quotas", flattenWorkspaceQuotas(utils.PathSearch("quotas", respBody,
			make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
