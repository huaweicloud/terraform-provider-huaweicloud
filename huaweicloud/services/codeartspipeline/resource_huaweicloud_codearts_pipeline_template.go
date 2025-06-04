package codeartspipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// CodeArtsPipelineTemplate is a tenant-level resource, do not involve CodeArts project
// @API CodeArtsPipeline POST /v5/{tenant_id}/api/pipeline-templates
// @API CodeArtsPipeline GET /v5/{tenant_id}/api/pipeline-templates/{template_id}
// @API CodeArtsPipeline PUT /v5/{tenant_id}/api/pipeline-templates/{template_id}
// @API CodeArtsPipeline DELETE /v5/{tenant_id}/api/pipeline-templates/{template_id}
// @API CodeArtsPipeline POST /v5/{tenant_id}/api/pipeline-templates/{templateId}/favorite
func ResourceCodeArtsPipelineTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelineTemplateCreate,
		ReadContext:   resourcePipelineTemplateRead,
		UpdateContext: resourcePipelineTemplateUpdate,
		DeleteContext: resourcePipelineTemplateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
				Description: `Specifies the template name.`,
			},
			"language": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the template language.`,
			},
			"definition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the template definition JSON.`,
			},
			"is_show_source": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies whether to display the pipeline source.`,
			},
			"variables": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the custom variables.`,
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
						"is_runtime": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Specifies whether to set parameters at runtime.`,
						},
						"is_reset": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Specifies whether to reset.`,
						},
						"latest_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the last parameter value.`,
						},
						"runtime_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the value passed in at runtime.`,
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the template description.`,
			},
			"is_favorite": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether it is a favorite template.`,
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the creation time.`,
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the last update time.`,
			},
			"creator_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creator.`,
			},
			"updater_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the last updater.`,
			},
			"icon": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the template icon.`,
			},
			"manifest_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the manifest version.`,
			},
			"is_system": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the template is a system template.`,
			},
		},
	}
}

func resourcePipelineTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v5/{tenant_id}/api/pipeline-templates"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{tenant_id}", cfg.DomainID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdatePipelineTemplateBodyParams(d, cfg.DomainID)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline template: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(createRespBody, ""); err != nil {
		return diag.Errorf("error creating CodeArts Pipeline template: %s", err)
	}

	id := utils.PathSearch("templateId", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the CodeArts Pipeline template ID from the API response")
	}

	d.SetId(id)

	if d.Get("is_favorite").(bool) {
		if err := updateCodeArtsTemplateIsFavorite(client, d, cfg.DomainID); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourcePipelineTemplateRead(ctx, d, meta)
}

func buildCreateOrUpdatePipelineTemplateBodyParams(d *schema.ResourceData, domainId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":           d.Get("name"),
		"language":       d.Get("language"),
		"definition":     d.Get("definition"),
		"variables":      buildPipelineTemplateVariables(d),
		"domain_id":      domainId,
		"is_show_source": d.Get("is_show_source"),
		"description":    d.Get("description"),

		// only support to create custom template
		"is_system": false,
	}

	return bodyParams
}

func buildPipelineTemplateVariables(d *schema.ResourceData) interface{} {
	rawVariables := d.Get("variables").([]interface{})
	if len(rawVariables) == 0 {
		return nil
	}

	variables := make([]map[string]interface{}, 0, len(rawVariables))
	for _, v := range rawVariables {
		if variable, ok := v.(map[string]interface{}); ok {
			customVar := map[string]interface{}{
				"name":          utils.ValueIgnoreEmpty(variable["name"]),
				"sequence":      utils.ValueIgnoreEmpty(variable["sequence"]),
				"type":          utils.ValueIgnoreEmpty(variable["type"]),
				"value":         utils.ValueIgnoreEmpty(variable["value"]),
				"is_secret":     utils.ValueIgnoreEmpty(variable["is_secret"]),
				"description":   utils.ValueIgnoreEmpty(variable["description"]),
				"is_runtime":    utils.ValueIgnoreEmpty(variable["is_runtime"]),
				"is_reset":      utils.ValueIgnoreEmpty(variable["is_reset"]),
				"latest_value":  utils.ValueIgnoreEmpty(variable["latest_value"]),
				"runtime_value": utils.ValueIgnoreEmpty(variable["runtime_value"]),
			}
			variables = append(variables, customVar)
		}
	}

	return variables
}

func updateCodeArtsTemplateIsFavorite(client *golangsdk.ServiceClient, d *schema.ResourceData, domainId string) error {
	httpUrl := "v5/{tenant_id}/api/pipeline-templates/{templateId}/favorite?flag={flag}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{tenant_id}", domainId)
	updatePath = strings.ReplaceAll(updatePath, "{templateId}", d.Id())
	updatePath = strings.ReplaceAll(updatePath, "{flag}", fmt.Sprintf("%t", d.Get("is_favorite").(bool)))
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating CodeArts Pipeline template is_favorite: %s", err)
	}
	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	if err := checkResponseError(updateRespBody, ""); err != nil {
		return fmt.Errorf("error updating CodeArts Pipeline template is_favorite: %s", err)
	}

	return nil
}

func GetPipelineTemplate(client *golangsdk.ServiceClient, domainId, id string) (interface{}, error) {
	httpUrl := "v5/{tenant_id}/api/pipeline-templates/{template_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{tenant_id}", domainId)
	getPath = strings.ReplaceAll(getPath, "{template_id}", id)
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

	if err := checkResponseError(getRespBody, templateNotFoundError); err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func resourcePipelineTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getRespBody, err := GetPipelineTemplate(client, cfg.DomainID, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CodeArts Pipeline template")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("language", utils.PathSearch("language", getRespBody, nil)),
		d.Set("definition", utils.PathSearch("definition", getRespBody, nil)),
		d.Set("variables", flattenPipelineTemplateVariables(getRespBody)),
		d.Set("is_system", utils.PathSearch("is_system", getRespBody, nil)),
		d.Set("is_show_source", utils.PathSearch("is_show_source", getRespBody, nil)),
		d.Set("creator_id", utils.PathSearch("creator_id", getRespBody, nil)),
		d.Set("updater_id", utils.PathSearch("updater_id", getRespBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", getRespBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", getRespBody, nil)),
		d.Set("icon", utils.PathSearch("icon", getRespBody, nil)),
		d.Set("manifest_version", utils.PathSearch("manifest_version", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPipelineTemplateVariables(resp interface{}) []interface{} {
	variables := utils.PathSearch("variables", resp, nil)
	if variables == nil {
		return nil
	}

	variablesList, ok := variables.([]interface{})
	if !ok || len(variablesList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(variablesList))
	for _, v := range variablesList {
		variable := v.(map[string]interface{})
		customVar := map[string]interface{}{
			"name":          utils.PathSearch("name", variable, nil),
			"sequence":      utils.PathSearch("sequence", variable, nil),
			"type":          utils.PathSearch("type", variable, nil),
			"value":         utils.PathSearch("value", variable, nil),
			"is_secret":     utils.PathSearch("is_secret", variable, nil),
			"description":   utils.PathSearch("description", variable, nil),
			"is_runtime":    utils.PathSearch("is_runtime", variable, nil),
			"is_reset":      utils.PathSearch("is_reset", variable, nil),
			"latest_value":  utils.PathSearch("latest_value", variable, nil),
			"runtime_value": utils.PathSearch("runtime_value", variable, nil),
		}
		result = append(result, customVar)
	}

	return result
}

func resourcePipelineTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	changes := []string{
		"name", "language", "definition", "is_show_source", "variables", "description",
	}

	if d.HasChanges(changes...) {
		httpUrl := "v5/{tenant_id}/api/pipeline-templates/{template_id}"
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{tenant_id}", cfg.DomainID)
		updatePath = strings.ReplaceAll(updatePath, "{template_id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildCreateOrUpdatePipelineTemplateBodyParams(d, cfg.DomainID)),
		}

		updateResp, err := client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating CodeArts Pipeline template: %s", err)
		}
		updateRespBody, err := utils.FlattenResponse(updateResp)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := checkResponseError(updateRespBody, ""); err != nil {
			return diag.Errorf("error updating CodeArts Pipeline template: %s", err)
		}
	}

	if d.HasChange("is_favorite") {
		if err := updateCodeArtsTemplateIsFavorite(client, d, cfg.DomainID); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourcePipelineTemplateRead(ctx, d, meta)
}

func resourcePipelineTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v5/{tenant_id}/api/pipeline-templates/{template_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{tenant_id}", cfg.DomainID)
	deletePath = strings.ReplaceAll(deletePath, "{template_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CodeArts Pipeline template")
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(deleteRespBody, templateNotFoundError); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CodeArts Pipeline template")
	}

	return nil
}
