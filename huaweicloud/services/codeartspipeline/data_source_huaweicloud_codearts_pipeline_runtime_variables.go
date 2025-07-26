package codeartspipeline

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

// @API CodeArtsPipeline GET /v5/{project_id}/api/pipelines/{pipeline_id}/list-runtime-vars
func DataSourceCodeArtsPipelineRuntimeVariables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelineRuntimeVariablesRead,

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
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsPipelineRuntimeVariablesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getHttpUrl := "v5/{project_id}/api/pipelines/{pipeline_id}/list-runtime-vars?limit=10"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", d.Get("project_id").(string))
	getPath = strings.ReplaceAll(getPath, "{pipeline_id}", d.Get("pipeline_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	offset := 0
	rst := make([]map[string]interface{}, 0)
	for {
		currentPath := getPath + fmt.Sprintf("&offset=%d", offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error getting pipeline runtime variables: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flattening response: %s", err)
		}
		if err := checkResponseError(getRespBody, ""); err != nil {
			return diag.Errorf("error getting pipeline runtime variables: %s", err)
		}

		variables := utils.PathSearch("variables", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(variables) == 0 {
			break
		}

		for _, variable := range variables {
			rst = append(rst, map[string]interface{}{
				"name":         utils.PathSearch("name", variable, nil),
				"sequence":     utils.PathSearch("sequence", variable, nil),
				"type":         utils.PathSearch("type", variable, nil),
				"value":        utils.PathSearch("value", variable, nil),
				"is_secret":    utils.PathSearch("is_secret", variable, nil),
				"description":  utils.PathSearch("description", variable, nil),
				"is_runtime":   utils.PathSearch("is_runtime", variable, nil),
				"is_reset":     utils.PathSearch("is_reset", variable, nil),
				"latest_value": utils.PathSearch("latest_value", variable, nil),
				"limits":       utils.PathSearch("limits", variable, nil),
			})
		}

		offset += len(variables)
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
