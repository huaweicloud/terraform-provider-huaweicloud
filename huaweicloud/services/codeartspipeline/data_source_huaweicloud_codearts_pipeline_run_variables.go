package codeartspipeline

import (
	"context"
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

// @API CodeArtsPipeline GET /v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipelineRunId}/run-variables
func DataSourceCodeArtsPipelineRunVariables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelineRunVariablesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the CodeArts project ID.`,
			},
			"pipeline_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the pipeline ID.`,
			},
			"pipeline_run_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the pipeline run ID.`,
			},
			"mode": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the query mode.`,
			},
			"variables": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the pipeline variables list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the custom variable name.`,
						},
						"sequence": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the parameter sequence, starting from 1.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the custom parameter type.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the custom parameter default value.`,
						},
						"is_secret": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether it is a private parameter.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the parameter description.`,
						},
						"is_runtime": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether to set parameters at runtime.`,
						},
						"is_reset": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether to reset.`,
						},
						"latest_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the last parameter value.`,
						},
						"limits": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Indicates the list of enumerated values.`,
						},
						"required": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the variable is required.`,
						},
						"variable_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the variable group name.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsPipelineRunVariablesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getHttpUrl := "v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipelineRunId}/run-variables?mode={mode}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", d.Get("project_id").(string))
	getPath = strings.ReplaceAll(getPath, "{pipeline_id}", d.Get("pipeline_id").(string))
	getPath = strings.ReplaceAll(getPath, "{pipelineRunId}", d.Get("pipeline_run_id").(string))
	getPath = strings.ReplaceAll(getPath, "{mode}", strconv.Itoa(d.Get("mode").(int)))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving pipeline run variables: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.Errorf("error flatten response: %s", err)
	}

	if err := checkResponseError(getRespBody, ""); err != nil {
		return diag.Errorf("error retrieving pipeline run variables: %s", err)
	}

	variables := getRespBody.([]interface{})
	rst := make([]map[string]interface{}, 0, len(variables))
	for _, variable := range variables {
		rst = append(rst, map[string]interface{}{
			"name":                utils.PathSearch("name", variable, nil),
			"sequence":            utils.PathSearch("sequence", variable, nil),
			"type":                utils.PathSearch("type", variable, nil),
			"value":               utils.PathSearch("value", variable, nil),
			"is_secret":           utils.PathSearch("is_secret", variable, nil),
			"description":         utils.PathSearch("description", variable, nil),
			"is_runtime":          utils.PathSearch("is_runtime", variable, nil),
			"is_reset":            utils.PathSearch("is_reset", variable, nil),
			"latest_value":        utils.PathSearch("latest_value", variable, nil),
			"limits":              utils.PathSearch("limits", variable, nil),
			"required":            utils.PathSearch("required", variable, nil),
			"variable_group_name": utils.PathSearch("variableGroupName", variable, nil),
		})
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("variables", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
