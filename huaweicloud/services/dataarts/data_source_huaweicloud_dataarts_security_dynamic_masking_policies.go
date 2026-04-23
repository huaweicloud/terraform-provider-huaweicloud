package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/security/masking/dynamic/policies
func DataSourceSecurityDynamicMaskingPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityDynamicMaskingPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the dynamic masking policies are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the dynamic masking policies belong.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the dynamic masking policy to be queried.`,
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the cluster to be queried.`,
			},
			"database_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the database to be queried.`,
			},
			"table_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the data table to be queried.`,
			},

			// Attributes.
			"policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of dynamic masking policies that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the dynamic masking policy.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the dynamic masking policy.`,
						},
						"datasource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data source type of the dynamic masking policy.`,
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the cluster corresponding to the data source.`,
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the cluster corresponding to the data source.`,
						},
						"database_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the database.`,
						},
						"table_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the data table.`,
						},
						"user_groups": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The user groups of the dynamic masking policy, separated by commas (,).`,
						},
						"users": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The users of the dynamic masking policy, separated by commas (,).`,
						},
						"sync_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The synchronization status of the dynamic masking policy.`,
						},
						"sync_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The synchronization time of the dynamic masking policy, in RFC3339 format.`,
						},
						"sync_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The synchronization log of the dynamic masking policy.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the dynamic masking policy, in RFC3339 format.`,
						},
						"create_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the dynamic masking policy.`,
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the dynamic masking policy, in RFC3339 format.`,
						},
						"update_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest updater of the dynamic masking policy.`,
						},
					},
				},
			},
		},
	}
}

func buildSecurityDynamicMaskingPoliciesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		res = fmt.Sprintf("%s&cluster_name=%v", res, v)
	}

	if v, ok := d.GetOk("database_name"); ok {
		res = fmt.Sprintf("%s&database_name=%v", res, v)
	}

	if v, ok := d.GetOk("table_name"); ok {
		res = fmt.Sprintf("%s&table_name=%v", res, v)
	}

	return res
}

func listSecurityDynamicMaskingPolicies(client *golangsdk.ServiceClient, workspaceId string, queryParams ...string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/security/masking/dynamic/policies?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	if len(queryParams) > 0 && queryParams[0] != "" {
		listPath += queryParams[0]
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		policies := utils.PathSearch("policies", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, policies...)
		if len(policies) < limit {
			break
		}

		offset += len(policies)
	}

	return result, nil
}

func flattenSecurityDynamicMaskingPolicies(policies []interface{}) []map[string]interface{} {
	if len(policies) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		result = append(result, map[string]interface{}{
			"id":              utils.PathSearch("id", policy, nil),
			"name":            utils.PathSearch("name", policy, nil),
			"datasource_type": utils.PathSearch("datasource_type", policy, nil),
			"cluster_id":      utils.PathSearch("cluster_id", policy, nil),
			"cluster_name":    utils.PathSearch("cluster_name", policy, nil),
			"database_name":   utils.PathSearch("database_name", policy, nil),
			"table_name":      utils.PathSearch("table_name", policy, nil),
			"user_groups":     utils.PathSearch("user_groups", policy, nil),
			"users":           utils.PathSearch("users", policy, nil),
			"sync_status":     utils.PathSearch("sync_status", policy, nil),
			"sync_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("sync_time",
				policy, float64(0)).(float64))/1000, false),
			"sync_msg": utils.PathSearch("sync_msg", policy, nil),
			"create_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time",
				policy, float64(0)).(float64))/1000, false),
			"create_user": utils.PathSearch("create_user", policy, nil),
			"update_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time",
				policy, float64(0)).(float64))/1000, false),
			"update_user": utils.PathSearch("update_user", policy, nil),
		})
	}

	return result
}

func dataSourceSecurityDynamicMaskingPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	policies, err := listSecurityDynamicMaskingPolicies(client, workspaceId, buildSecurityDynamicMaskingPoliciesQueryParams(d))
	if err != nil {
		return diag.Errorf("error querying DataArts Security dynamic masking policies: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("policies", flattenSecurityDynamicMaskingPolicies(policies)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
