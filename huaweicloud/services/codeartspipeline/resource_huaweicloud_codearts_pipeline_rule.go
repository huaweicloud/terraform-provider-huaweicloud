package codeartspipeline

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var ruleNonUpdatableParams = []string{
	"type", "layout_content",
}

// @API CodeArtsPipeline POST /v2/{domain_id}/rules/create
// @API CodeArtsPipeline GET /v2/{domain_id}/rules/{rule_id}/detail
// @API CodeArtsPipeline PUT /v2/{domain_id}/rules/{rule_id}/update
// @API CodeArtsPipeline DELETE /v2/{domain_id}/rules/{rule_id}/delete
// @API CodeArtsPipeline GET /v2/{domainId}/rules/{ruleId}/related/query
func ResourceCodeArtsPipelineRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelineRuleCreate,
		ReadContext:   resourcePipelineRuleRead,
		UpdateContext: resourcePipelineRuleUpdate,
		DeleteContext: resourcePipelineRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(ruleNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the rule name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the rule type.`,
			},
			"layout_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the layout content.`,
			},
			"content": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: `Specifies the rule attribute group list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the group name.`,
						},
						"properties": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: `Specifies the rule attribute list.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Specifies the display name.`,
									},
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Specifies the type.`,
									},
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Specifies the attribute key.`,
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Specifies the attribute value.`,
									},
									"value_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Specifies the value type.`,
									},
									"operator": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the comparison operators.`,
									},
									"is_valid": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `Specifies wether the property is valid.`,
									},
								},
							},
						},
						"can_modify_when_inherit": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Specifies whether thresholds of an inherited policy can be modified.`,
						},
					},
				},
			},
			"plugin_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the plugin ID.`,
			},
			"plugin_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the plugin name.`,
			},
			"plugin_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the plugin version.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the rule version.`,
			},
			"creator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creator.`,
			},
			"updater": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the updater.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the create time.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the update time.`,
			},
			"rule_set_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of policies.`,
			},
			"project_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of projects.`,
			},
			"pipeline_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of pipelines.`,
			},
		},
	}
}

func resourcePipelineRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v2/{domain_id}/rules/create"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{domain_id}", cfg.DomainID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdatePipelineRuleBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline rule: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(createRespBody, ""); err != nil {
		return diag.Errorf("error creating CodeArts Pipeline rule: %s", err)
	}

	id := utils.PathSearch("rule_id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the CodeArts Pipeline rule ID from the API response")
	}

	d.SetId(id)

	return resourcePipelineRuleRead(ctx, d, meta)
}

func buildCreateOrUpdatePipelineRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":           d.Get("name"),
		"type":           d.Get("type"),
		"layout_content": d.Get("layout_content"),
		"content":        buildPipelineRuleContent(d),
		"plugin_id":      utils.ValueIgnoreEmpty(d.Get("plugin_id")),
		"plugin_name":    utils.ValueIgnoreEmpty(d.Get("plugin_name")),
		"plugin_version": utils.ValueIgnoreEmpty(d.Get("plugin_version")),
	}

	return bodyParams
}

func buildPipelineRuleContent(d *schema.ResourceData) interface{} {
	rawParams := d.Get("content").(*schema.Set).List()
	rst := make([]map[string]interface{}, 0, len(rawParams))
	for _, p := range rawParams {
		if params, ok := p.(map[string]interface{}); ok {
			m := map[string]interface{}{
				"group_name":              params["group_name"],
				"properties":              buildPipelineRuleContentProperties(params["properties"].(*schema.Set).List()),
				"can_modify_when_inherit": params["can_modify_when_inherit"],
			}
			rst = append(rst, m)
		}
	}

	return rst
}

func buildPipelineRuleContentProperties(rawParams []interface{}) interface{} {
	rst := make([]map[string]interface{}, 0, len(rawParams))
	for _, p := range rawParams {
		if params, ok := p.(map[string]interface{}); ok {
			m := map[string]interface{}{
				"name":       params["name"],
				"type":       params["type"],
				"key":        params["key"],
				"value":      params["value"],
				"value_type": params["value_type"],
				"is_valid":   params["is_valid"],
				"operator":   utils.ValueIgnoreEmpty(params["operator"]),
			}
			rst = append(rst, m)
		}
	}

	return rst
}

func resourcePipelineRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v2/{domain_id}/rules/{rule_id}/detail"
	getRespBody, err := GetPipelineRuleItems(client, httpUrl, cfg.DomainID, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CodeArts Pipeline rule")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getRespBody, nil)),
		d.Set("content", flattenPipelineRuleContent(getRespBody)),
		d.Set("plugin_id", utils.PathSearch("plugin_id", getRespBody, nil)),
		d.Set("plugin_name", utils.PathSearch("plugin_name", getRespBody, nil)),
		d.Set("plugin_version", utils.PathSearch("plugin_version", getRespBody, nil)),
		d.Set("creator", utils.PathSearch("creator", getRespBody, nil)),
		d.Set("updater", utils.PathSearch("updater", getRespBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", getRespBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", getRespBody, nil)),
		d.Set("version", utils.PathSearch("version", getRespBody, nil)),
	)

	getPipelineRuleUsageHttpURl := "v2/{domain_id}/rules/{rule_id}/related/query"
	if usage, err := GetPipelineRuleItems(client, getPipelineRuleUsageHttpURl, cfg.DomainID, d.Id()); err != nil {
		log.Println("[WARN] error retrieving rule relation details")
	} else {
		mErr = multierror.Append(mErr,
			d.Set("rule_set_count", utils.PathSearch("rule_set_count", usage, nil)),
			d.Set("project_count", utils.PathSearch("project_count", usage, nil)),
			d.Set("pipeline_count", utils.PathSearch("pipeline_count", usage, nil)),
		)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetPipelineRuleItems(client *golangsdk.ServiceClient, httpUrl, domainId, id string) (interface{}, error) {
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", domainId)
	getPath = strings.ReplaceAll(getPath, "{rule_id}", id)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	if err := checkResponseError(getRespBody, ""); err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func flattenPipelineRuleContent(resp interface{}) []interface{} {
	contentList, ok := utils.PathSearch("content", resp, make([]interface{}, 0)).([]interface{})
	if ok && len(contentList) > 0 {
		result := make([]interface{}, 0, len(contentList))
		for _, v := range contentList {
			content := v.(map[string]interface{})
			m := map[string]interface{}{
				"group_name":              utils.PathSearch("group_name", content, nil),
				"properties":              flattenPipelineRuleContentProperties(content),
				"can_modify_when_inherit": utils.PathSearch("can_modify_when_inherit", content, nil),
			}
			result = append(result, m)
		}
		return result
	}

	return nil
}

func flattenPipelineRuleContentProperties(resp interface{}) []interface{} {
	propertiesList, ok := utils.PathSearch("properties", resp, make([]interface{}, 0)).([]interface{})
	if ok && len(propertiesList) > 0 {
		result := make([]interface{}, 0, len(propertiesList))
		for _, v := range propertiesList {
			property := v.(map[string]interface{})
			m := map[string]interface{}{
				"name":       utils.PathSearch("name", property, nil),
				"type":       utils.PathSearch("type", property, nil),
				"key":        utils.PathSearch("key", property, nil),
				"value":      utils.PathSearch("value", property, nil),
				"value_type": utils.PathSearch("value_type", property, nil),
				"operator":   utils.PathSearch("operator", property, nil),
				"is_valid":   utils.PathSearch("is_valid", property, nil),
			}

			result = append(result, m)
		}
		return result
	}

	return nil
}

func resourcePipelineRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v2/{domain_id}/rules/{rule_id}/update"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{domain_id}", cfg.DomainID)
	updatePath = strings.ReplaceAll(updatePath, "{rule_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdatePipelineRuleBodyParams(d)),
	}

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating CodeArts Pipeline rule: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(updateRespBody, ""); err != nil {
		return diag.Errorf("error updating CodeArts Pipeline rule: %s", err)
	}

	return resourcePipelineRuleRead(ctx, d, meta)
}

func resourcePipelineRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v2/{domain_id}/rules/{rule_id}/delete"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{domain_id}", cfg.DomainID)
	deletePath = strings.ReplaceAll(deletePath, "{rule_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CodeArts Pipeline rule")
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(deleteRespBody, ""); err != nil {
		return diag.Errorf("error deleting CodeArts Pipeline rule: %s", err)
	}

	return nil
}
