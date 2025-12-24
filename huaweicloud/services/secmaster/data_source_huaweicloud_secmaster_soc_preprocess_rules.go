package secmaster

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

// @API Secmaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/mappings/preprocess-rules/search
func DataSourceSocPreprocessRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSocPreprocessRulesRead,

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
			// This field is `Required` in the API documentation, but actual is `Optional`.
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mapping_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mapper_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
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
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mapping_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mapper_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mapper_type_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expression": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildSocPreprocessRulesParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":       utils.ValueIgnoreEmpty(d.Get("name")),
		"mapping_id": utils.ValueIgnoreEmpty(d.Get("mapping_id")),
		"mapper_ids": utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("mapper_ids").([]interface{}))),
		"start_time": utils.ValueIgnoreEmpty(d.Get("start_time")),
		"end_time":   utils.ValueIgnoreEmpty(d.Get("end_time")),
		"limit":      100,
	}
}

func dataSourceSocPreprocessRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/mappings/preprocess-rules/search"
		workspaceId = d.Get("workspace_id").(string)
		offset      = 0
		allData     = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	bodyParams := utils.RemoveNil(buildSocPreprocessRulesParams(d))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		bodyParams["offset"] = offset
		requestOpt.JSONBody = bodyParams

		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving preprocess rules: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		data := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		if len(data) == 0 {
			break
		}

		allData = append(allData, data...)
		totalCount := int(utils.PathSearch("total", respBody, float64(0)).(float64))
		if len(allData) >= totalCount {
			break
		}

		offset += len(data)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenSocPreprocessRules(allData)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSocPreprocessRules(rules []interface{}) []interface{} {
	if len(rules) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rules))
	for _, v := range rules {
		rst = append(rst, map[string]interface{}{
			"id":             utils.PathSearch("id", v, nil),
			"name":           utils.PathSearch("name", v, nil),
			"project_id":     utils.PathSearch("project_id", v, nil),
			"workspace_id":   utils.PathSearch("workspace_id", v, nil),
			"mapping_id":     utils.PathSearch("mapping_id", v, nil),
			"mapper_id":      utils.PathSearch("mapper_id", v, nil),
			"mapper_type_id": utils.PathSearch("mapper_type_id", v, nil),
			"action":         utils.PathSearch("action", v, nil),
			"expression":     utils.PathSearch("expression", v, nil),
			"create_time":    utils.PathSearch("create_time", v, nil),
			"update_time":    utils.PathSearch("update_time", v, nil),
		})
	}

	return rst
}
