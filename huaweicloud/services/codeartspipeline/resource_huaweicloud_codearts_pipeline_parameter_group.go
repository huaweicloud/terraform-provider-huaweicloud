package codeartspipeline

import (
	"context"
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

var parameterGroupNonUpdatableParams = []string{
	"project_id",
}

// @API CodeArtsPipeline POST /v5/{project_id}/api/pipeline/variable/group/create
// @API CodeArtsPipeline GET /v5/{project_id}/api/pipeline/variable/group/{id}
// @API CodeArtsPipeline PUT /v5/{project_id}/api/pipeline/variable/group/update
// @API CodeArtsPipeline DELETE /v5/{project_id}/api/pipeline/variable/group/delete
func ResourceCodeArtsPipelineParameterGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelineParameterGroupCreate,
		ReadContext:   resourcePipelineParameterGroupRead,
		UpdateContext: resourcePipelineParameterGroupUpdate,
		DeleteContext: resourcePipelineParameterGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceImportStateFuncWithProjectIdAndId,
		},

		CustomizeDiff: config.FlexibleForceNew(parameterGroupNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the CodeArts project ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the parameter group name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the parameter group description.`,
			},
			"variables": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the permission information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the custom variable name.`,
						},
						"sequence": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `Specifies the parameter sequence, starting from 1.`,
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the custom parameter type.`,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the custom parameter default value.`,
						},
						"is_secret": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Specifies whether it is a private parameter.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the parameter description.`,
						},
					},
				},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"related_pipelines": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the associated pipeline.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the pipeline name.`,
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the pipeline ID.`,
						},
					},
				},
			},
			"creator_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creator ID.`,
			},
			"updater_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the updater ID.`,
			},
			"creator_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creator name.`,
			},
			"updater_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the updater name.`,
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the create time.`,
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the update time.`,
			},
		},
	}
}

func resourcePipelineParameterGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v5/{project_id}/api/pipeline/variable/group/create"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", d.Get("project_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdatePipelineParameterGroupBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline parameter group: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(createRespBody, ""); err != nil {
		return diag.Errorf("error creating CodeArts Pipeline parameter group: %s", err)
	}

	id := utils.PathSearch("pipeline_variable_group_id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the CodeArts Pipeline parameter group ID from the API response")
	}

	d.SetId(id)

	return resourcePipelineParameterGroupRead(ctx, d, meta)
}

func buildCreateOrUpdatePipelineParameterGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"project_id":  d.Get("project_id"),
		"name":        d.Get("name"),
		"description": d.Get("description"),
		"variables":   buildPipelineParameterGroupVariables(d),
		"id":          utils.ValueIgnoreEmpty(d.Id()),
	}

	return bodyParams
}

func buildPipelineParameterGroupVariables(d *schema.ResourceData) interface{} {
	rawVariables := d.Get("variables").(*schema.Set).List()
	if len(rawVariables) == 0 {
		return nil
	}

	variables := make([]map[string]interface{}, 0, len(rawVariables))
	for _, v := range rawVariables {
		if variable, ok := v.(map[string]interface{}); ok {
			customVar := map[string]interface{}{
				"name":        utils.ValueIgnoreEmpty(variable["name"]),
				"sequence":    utils.ValueIgnoreEmpty(variable["sequence"]),
				"type":        utils.ValueIgnoreEmpty(variable["type"]),
				"value":       utils.ValueIgnoreEmpty(variable["value"]),
				"is_secret":   utils.ValueIgnoreEmpty(variable["is_secret"]),
				"description": utils.ValueIgnoreEmpty(variable["description"]),
			}
			variables = append(variables, customVar)
		}
	}

	return variables
}

func resourcePipelineParameterGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getRespBody, err := GetPipelineParameterGroup(client, d.Get("project_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CodeArts Pipeline parameter group")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("project_id", utils.PathSearch("project_id", getRespBody, nil)),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("creator_id", utils.PathSearch("creator_id", getRespBody, nil)),
		d.Set("updater_id", utils.PathSearch("updater_id", getRespBody, nil)),
		d.Set("creator_name", utils.PathSearch("creator_name", getRespBody, nil)),
		d.Set("updater_name", utils.PathSearch("updater_name", getRespBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", getRespBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", getRespBody, nil)),
		d.Set("variables", flattenPipelineParameterGroupVariables(getRespBody)),
		d.Set("related_pipelines", flattenPipelineParameterGroupRelatedPipelines(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetPipelineParameterGroup(client *golangsdk.ServiceClient, projectId, id string) (interface{}, error) {
	httpUrl := "v5/{project_id}/api/pipeline/variable/group/{id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", projectId)
	getPath = strings.ReplaceAll(getPath, "{id}", id)
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
	if err := checkResponseError(getRespBody, parameterGroupNotFoundError); err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func flattenPipelineParameterGroupVariables(resp interface{}) []interface{} {
	variablesList, ok := utils.PathSearch("variables", resp, make([]interface{}, 0)).([]interface{})
	if ok && len(variablesList) > 0 {
		result := make([]interface{}, 0, len(variablesList))
		for _, v := range variablesList {
			variable := v.(map[string]interface{})
			customVar := map[string]interface{}{
				"name":        utils.PathSearch("name", variable, nil),
				"sequence":    utils.PathSearch("sequence", variable, nil),
				"type":        utils.PathSearch("type", variable, nil),
				"value":       utils.PathSearch("value", variable, nil),
				"is_secret":   utils.PathSearch("is_secret", variable, nil),
				"description": utils.PathSearch("description", variable, nil),
			}
			result = append(result, customVar)
		}
		return result
	}

	return nil
}

func flattenPipelineParameterGroupRelatedPipelines(resp interface{}) []interface{} {
	relatedPipelinesList, ok := utils.PathSearch("related_pipelines", resp, make([]interface{}, 0)).([]interface{})
	if ok && len(relatedPipelinesList) > 0 {
		result := make([]interface{}, 0, len(relatedPipelinesList))
		for _, v := range relatedPipelinesList {
			relatedPipeline := v.(map[string]interface{})
			m := map[string]interface{}{
				"name": utils.PathSearch("name", relatedPipeline, nil),
				"id":   utils.PathSearch("id", relatedPipeline, nil),
			}
			result = append(result, m)
		}
		return result
	}

	return nil
}

func resourcePipelineParameterGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v5/{project_id}/api/pipeline/variable/group/update"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", d.Get("project_id").(string))
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdatePipelineParameterGroupBodyParams(d)),
	}

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating CodeArts Pipeline parameter group: %s", err)
	}
	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(updateRespBody, ""); err != nil {
		return diag.Errorf("error updating CodeArts Pipeline parameter group: %s", err)
	}

	return resourcePipelineParameterGroupRead(ctx, d, meta)
}

func resourcePipelineParameterGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	projectId := d.Get("project_id").(string)
	httpUrl := "v5/{project_id}/api/pipeline/variable/group/delete?id={id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", projectId)
	deletePath = strings.ReplaceAll(deletePath, "{id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         map[string]interface{}{},
	}

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CodeArts Pipeline parameter group")
	}
	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(deleteRespBody, parameterGroupNotFoundError); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CodeArts Pipeline parameter group")
	}

	return nil
}
