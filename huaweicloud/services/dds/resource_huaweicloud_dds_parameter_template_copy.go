package dds

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDS POST /v3/{project_id}/configurations/{config_id}/copy
// @API DDS DELETE /v3/{project_id}/configurations/{config_id}
// @API DDS GET /v3/{project_id}/configurations/{config_id}
// @API DDS PUT /v3/{project_id}/configurations/{config_id}
func ResourceDDSParameterTemplateCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceParameterTemplateCopyCreate,
		UpdateContext: resourceDdsParameterTemplateUpdate,
		ReadContext:   resourceDdsParameterTemplateRead,
		DeleteContext: resourceDdsParameterTemplateDelete,

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
			"configuration_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the parameter template ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of replicated parameter template.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of replicated parameter template.`,
			},
			"parameter_values": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the mapping between parameter names and parameter values.`,
			},
			"node_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the database version.`,
			},
			"parameters": {
				Type:        schema.TypeList,
				Elem:        ParameterTemplateParameterSchema(),
				Computed:    true,
				Description: `Indicates the parameters defined by users based on the default parameter templates.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the create time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the update time.`,
			},
		},
	}
}

func resourceParameterTemplateCopyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	configurationId := d.Get("configuration_id").(string)

	httpUrl := "v3/{project_id}/configurations/{config_id}/copy"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{config_id}", configurationId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateParameterTemplateCopyBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error copying DDS parameter template: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("configuration_id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find ID from the API response")
	}

	d.SetId(id)

	if _, ok := d.GetOk("parameter_values"); ok {
		if err := updateParameterTemplate(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDdsParameterTemplateRead(ctx, d, meta)
}

func buildCreateParameterTemplateCopyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}
