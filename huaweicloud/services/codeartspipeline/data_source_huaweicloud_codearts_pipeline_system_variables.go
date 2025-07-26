package codeartspipeline

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

// @API CodeArtsPipeline GET /v5/{project_id}/api/pipelines/{pipeline_id}/list-system-vars
func DataSourceCodeArtsPipelineSystemVariables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelineSystemVariablesRead,

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
							Description: `Indicates the system variable name.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the system parameter value.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the system parameter type.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the parameter description.`,
						},
						"is_show": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether it is showed.`,
						},
						"ordinal": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the parameter ordinal.`,
						},
						"is_alias": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the name is alias.`,
						},
						"kind": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the parameter context type.`,
						},
						"context_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the context name.`,
						},
						"source_identifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the source identifier.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsPipelineSystemVariablesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getHttpUrl := "v5/{project_id}/api/pipelines/{pipeline_id}/list-system-vars"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", d.Get("project_id").(string))
	getPath = strings.ReplaceAll(getPath, "{pipeline_id}", d.Get("pipeline_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving pipeline system variables: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.Errorf("error flatten response: %s", err)
	}

	if err := checkResponseError(getRespBody, ""); err != nil {
		return diag.Errorf("error retrieving pipeline system variables: %s", err)
	}

	variables := getRespBody.([]interface{})
	rst := make([]map[string]interface{}, 0, len(variables))
	for _, variable := range variables {
		rst = append(rst, map[string]interface{}{
			"name":              utils.PathSearch("name", variable, nil),
			"type":              utils.PathSearch("type", variable, nil),
			"value":             utils.PathSearch("value", variable, nil),
			"description":       utils.PathSearch("description", variable, nil),
			"is_show":           utils.PathSearch("isShow", variable, nil),
			"ordinal":           utils.PathSearch("ordinal", variable, nil),
			"is_alias":          utils.PathSearch("isAlias", variable, nil),
			"kind":              utils.PathSearch("kind", variable, nil),
			"context_name":      utils.PathSearch("contextName", variable, nil),
			"source_identifier": utils.PathSearch("sourceIdentifier", variable, nil),
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
