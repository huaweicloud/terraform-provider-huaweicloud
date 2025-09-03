package secmaster

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

// @API Secmaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/assetcredentials
func DataSourceOperationConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOperationConnectionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"component_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"creator_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"modifier_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"update_start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"update_end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_defense_type": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// The field in the API documentation is `data`, but it is actually returned as `assets`.
			// However, `assets` cannot accurately express the actual meaning of this field,
			// and all use of `data`should be consistent with the API documentation.
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"component_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"component_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"component_version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modifier_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modifier_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"defense_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_enterprise_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_enterprise_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildOperationConnectionsQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("component_name"); ok {
		queryParams = fmt.Sprintf("%s&component_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("creator_name"); ok {
		queryParams = fmt.Sprintf("%s&creator_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("modifier_name"); ok {
		queryParams = fmt.Sprintf("%s&modifier_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("description"); ok {
		queryParams = fmt.Sprintf("%s&description=%v", queryParams, v)
	}
	if v, ok := d.GetOk("create_start_time"); ok {
		queryParams = fmt.Sprintf("%s&create_start_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("create_end_time"); ok {
		queryParams = fmt.Sprintf("%s&create_end_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("update_start_time"); ok {
		queryParams = fmt.Sprintf("%s&update_start_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("update_end_time"); ok {
		queryParams = fmt.Sprintf("%s&update_end_time=%v", queryParams, v)
	}
	if d.Get("is_defense_type").(bool) {
		queryParams = fmt.Sprintf("%s&is_defense_type=%v", queryParams, true)
	}

	return queryParams
}

func dataSourceOperationConnectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/assetcredentials"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath += buildOperationConnectionsQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json",
		},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		requestResp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster operation connections: %s", err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("assets", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		result = append(result, dataResp...)
		offset += len(dataResp)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenOperationConnectionsData(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOperationConnectionsData(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		rst = append(rst, map[string]interface{}{
			"id":                     utils.PathSearch("id", v, nil),
			"project_id":             utils.PathSearch("project_id", v, nil),
			"workspace_id":           utils.PathSearch("workspace_id", v, nil),
			"name":                   utils.PathSearch("name", v, nil),
			"component_id":           utils.PathSearch("component_id", v, nil),
			"component_name":         utils.PathSearch("component_name", v, nil),
			"component_version_id":   utils.PathSearch("component_version_id", v, nil),
			"type":                   utils.PathSearch("type", v, nil),
			"status":                 utils.PathSearch("status", v, nil),
			"config":                 utils.PathSearch("config", v, nil),
			"description":            utils.PathSearch("description", v, nil),
			"enabled":                utils.PathSearch("enabled", v, nil),
			"create_time":            utils.PathSearch("create_time", v, nil),
			"creator_id":             utils.PathSearch("creator_id", v, nil),
			"creator_name":           utils.PathSearch("creator_name", v, nil),
			"update_time":            utils.PathSearch("update_time", v, nil),
			"modifier_id":            utils.PathSearch("modifier_id", v, nil),
			"modifier_name":          utils.PathSearch("modifier_name", v, nil),
			"defense_type":           utils.PathSearch("defense_type", v, nil),
			"target_project_id":      utils.PathSearch("target_project_id", v, nil),
			"target_project_name":    utils.PathSearch("target_project_name", v, nil),
			"target_enterprise_id":   utils.PathSearch("target_enterprise_id", v, nil),
			"target_enterprise_name": utils.PathSearch("target_enterprise_name", v, nil),
			"region_id":              utils.PathSearch("region_id", v, nil),
			"region_name":            utils.PathSearch("region_name", v, nil),
		})
	}

	return rst
}
