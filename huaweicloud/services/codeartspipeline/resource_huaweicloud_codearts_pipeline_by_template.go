package codeartspipeline

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var pipelineByTemplateNonUpdatableParams = []string{
	"project_id", "component_id", "template_id",
}

// @API CodeArtsPipeline POST /v5/{project_id}/api/pipeline-templates/{template_id}/create-pipeline
// @API CodeArtsPipeline GET /v5/{project_id}/api/pipelines/{pipeline_id}
// @API CodeArtsPipeline PUT /v5/{project_id}/api/pipelines/{pipeline_id}
// @API CodeArtsPipeline DELETE /v5/{project_id}/api/pipelines/{pipeline_id}
// @API CodeArtsPipeline PUT /v5/{project_id}/api/pipelines/{pipeline_id}/unban
// @API CodeArtsPipeline PUT /v5/{project_id}/api/pipelines/{pipeline_id}/ban
func ResourceCodeArtsPipelineByTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelineByTemplateCreate,
		ReadContext:   resourcePipelineRead,
		UpdateContext: resourcePipelineUpdate,
		DeleteContext: resourcePipelineDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceImportStateFuncWithProjectIdAndId,
		},

		CustomizeDiff: config.FlexibleForceNew(pipelineByTemplateNonUpdatableParams),

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
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the CodeArts template ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the pipeline name.`,
			},
			"is_publish": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies whether it is a change-triggered pipeline.`,
			},
			"sources": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: `Specifies the pipeline source information.`,
				Elem:        resourceSchemePipelineSources(),
			},
			"component_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the microservice ID.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the pipeline description.`,
			},
			"manifest_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the pipeline structure definition version.`,
			},
			"definition": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the pipeline definition JSON.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the pipeline group ID.`,
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the project name.`,
			},
			"banned": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether the pipeline is banned.`,
			},
			"variables": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the custom variables.`,
				Elem:        resourceSchemePipelineByTemplateVariables(),
			},
			"schedules": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the pipeline schedule settings.`,
				Elem:        resourceSchemePipelineSchedules(),
			},
			"triggers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the pipeline trigger settings.`,
				Elem:        resourceSchemePipelineTriggers(),
			},
			"concurrency_control": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Specifies the pipeline concurrency control information.`,
				Elem:        resourceSchemePipelineConcurrencyControl(),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"creator_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creator ID.`,
			},
			"creator_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creator name.`,
			},
			"updater_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the last updater ID.`,
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
			"is_collect": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the current user has collected it.`,
			},
		},
	}
}

func resourceSchemePipelineByTemplateVariables() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the custom variable name.`,
			},
			"sequence": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the parameter sequence, starting from 1.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the custom parameter type.`,
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the custom parameter default value.`,
			},
			"is_secret": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether it is a private parameter.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the parameter description.`,
			},
			"is_runtime": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to set parameters at runtime.`,
			},
			"is_reset": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
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
				Computed:    true,
				Description: `Specifies the value passed in at runtime.`,
			},
			"limits": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of enumerated values.`,
			},
		},
	}
}

func resourcePipelineByTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	projectId := d.Get("project_id").(string)
	templateId := d.Get("template_id").(string)

	httpUrl := "v5/{project_id}/api/pipeline-templates/{template_id}/create-pipeline"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", projectId)
	createPath = strings.ReplaceAll(createPath, "{template_id}", templateId)
	if v, ok := d.GetOk("component_id"); ok {
		createPath += fmt.Sprintf("?component_id=%v", v)
	}

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdatePipelineByTemplateBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline by template: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(createRespBody, ""); err != nil {
		return diag.Errorf("error creating CodeArts Pipeline by template: %s", err)
	}

	id := utils.PathSearch("pipeline_id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the CodeArts Pipeline ID from the API response")
	}

	d.SetId(id)

	changes := []string{
		"manifest_version", "definition", "group_id", "project_name",
		"schedules", "triggers", "concurrency_control",
	}

	if d.HasChanges(changes...) {
		if err := updatePipelineAfterCreateByPipeline(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.Get("banned").(bool) {
		if err := updatePipelineBanned(client, banHttpUrl, projectId, d.Id()); err != nil {
			return diag.Errorf("error banning pipeline: %s", err)
		}
	}

	return resourcePipelineRead(ctx, d, meta)
}

func buildCreateOrUpdatePipelineByTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
		"is_publish":  d.Get("is_publish"),
		"sources":     buildPipelineSources(d),
		"variables":   buildPipelineTemplateVariables(d),
	}

	return bodyParams
}

func updatePipelineAfterCreateByPipeline(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	projectId := d.Get("project_id").(string)
	definition := d.Get("definition").(string)
	variablesLen := d.Get("variables").(*schema.Set).Len()
	if definition == "" || variablesLen == 0 {
		getRespBody, err := GetPipeline(client, projectId, d.Id())
		if err != nil {
			return errors.New("error retrieving CodeArts Pipeline")
		}
		if definition == "" {
			d.Set("definition", utils.PathSearch("definition", getRespBody, nil))
		}
		if variablesLen == 0 {
			d.Set("variables", flattenPipelineTemplateVariables(getRespBody))
		}
	}

	return updatePipeline(client, d)
}
