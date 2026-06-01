package geminidb

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var geminiDBParameterTemplateCopyNonUpdatableParams = []string{
	"config_id",
}

// @API GeminiDB POST /v3/{project_id}/configurations/{config_id}/copy
// @API GeminiDB DELETE /v3/{project_id}/configurations/{config_id}
// @API GeminiDB GET /v3/{project_id}/configurations/{config_id}
// @API GeminiDB PUT /v3/{project_id}/configurations/{config_id}
func ResourceGeminiDBParameterTemplateCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeminiDBParameterTemplateCopyCreate,
		UpdateContext: resourceGeminiDBParameterTemplateUpdate,
		ReadContext:   resourceGeminiDBParameterTemplateRead,
		DeleteContext: resourceGeminiDBParameterTemplateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(geminiDBParameterTemplateCopyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"config_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"values": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"datastore_version_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"configuration_parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     geminiDBParameterTemplateParameterSchema(),
			},
		},
	}
}

func resourceGeminiDBParameterTemplateCopyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/configurations/{config_id}/copy"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{config_id}", d.Get("config_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildGeminiDBParameterTemplateCopyBody(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error copying GeminiDB configuration: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	configID := utils.PathSearch("config_id", respBody, "").(string)
	if configID == "" {
		return diag.Errorf("unable to find configuration ID from API response")
	}

	d.SetId(configID)

	if _, ok := d.GetOk("values"); ok {
		if err := updateGeminiDBParameterTemplate(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGeminiDBParameterTemplateRead(ctx, d, meta)
}

func buildGeminiDBParameterTemplateCopyBody(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return body
}
