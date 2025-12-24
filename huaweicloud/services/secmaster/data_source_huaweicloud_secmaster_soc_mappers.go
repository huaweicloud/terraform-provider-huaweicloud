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

// @API Secmaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/mappings/mappers/search
func DataSourceSocMapper() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSocMapperRead,

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
			"mapping_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"has_preprocess_rule": {
				Type:     schema.TypeBool,
				Optional: true,
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
						"dataclass_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dataclass_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mapper_type_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mapping_id": {
							Type:     schema.TypeString,
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
					},
				},
			},
		},
	}
}

func buildSocMapperParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"mapping_id": d.Get("mapping_id"),
		"name":       utils.ValueIgnoreEmpty(d.Get("name")),
		"start_time": utils.ValueIgnoreEmpty(d.Get("start_time")),
		"end_time":   utils.ValueIgnoreEmpty(d.Get("end_time")),
		"limit":      100,
	}
}

func dataSourceSocMapperRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg               = meta.(*config.Config)
		region            = cfg.GetRegion(d)
		httpUrl           = "v1/{project_id}/workspaces/{workspace_id}/soc/mappings/mappers/search"
		workspaceId       = d.Get("workspace_id").(string)
		hasPreprocessRule = d.Get("has_preprocess_rule").(bool)
		offset            = 0
		allData           = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	bodyParams := utils.RemoveNil(buildSocMapperParams(d))
	if hasPreprocessRule {
		bodyParams["has_preprocess_rule"] = true
	}

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
			return diag.Errorf("error retrieving mapper list: %s", err)
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
		d.Set("data", flattenSocMapper(allData)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSocMapper(mapperInfos []interface{}) []interface{} {
	if len(mapperInfos) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(mapperInfos))
	for _, v := range mapperInfos {
		rst = append(rst, map[string]interface{}{
			"id":             utils.PathSearch("id", v, nil),
			"name":           utils.PathSearch("name", v, nil),
			"project_id":     utils.PathSearch("project_id", v, nil),
			"workspace_id":   utils.PathSearch("workspace_id", v, nil),
			"dataclass_id":   utils.PathSearch("dataclass_id", v, nil),
			"dataclass_name": utils.PathSearch("dataclass_name", v, nil),
			"mapper_type_id": utils.PathSearch("mapper_type_id", v, nil),
			"mapping_id":     utils.PathSearch("mapping_id", v, nil),
			"create_time":    utils.PathSearch("create_time", v, nil),
			"creator_id":     utils.PathSearch("creator_id", v, nil),
			"creator_name":   utils.PathSearch("creator_name", v, nil),
			"update_time":    utils.PathSearch("update_time", v, nil),
			"modifier_id":    utils.PathSearch("modifier_id", v, nil),
			"modifier_name":  utils.PathSearch("modifier_name", v, nil),
		})
	}

	return rst
}
