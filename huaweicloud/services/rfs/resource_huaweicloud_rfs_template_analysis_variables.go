package rfs

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var templateAnalysisVariablesNonUpdatableParams = []string{
	"template_body",
	"template_uri",
}

// @API RFS POST /v1/{project_id}/template-analyses/variables
func ResourceTemplateAnalysisVariables() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTemplateAnalysisVariablesCreate,
		ReadContext:   resourceTemplateAnalysisVariablesRead,
		UpdateContext: resourceTemplateAnalysisVariablesUpdate,
		DeleteContext: resourceTemplateAnalysisVariablesDelete,

		CustomizeDiff: config.FlexibleForceNew(templateAnalysisVariablesNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"template_body": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"variables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataVariablesSchema(),
			},
		},
	}
}

func dataVariablesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Field `default` is object type in API response, but we need to convert it to string.
			"default": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sensitive": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"nullable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"validations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildTemplateAnalysisVariablesBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"template_body": utils.ValueIgnoreEmpty(d.Get("template_body")),
		"template_uri":  utils.ValueIgnoreEmpty(d.Get("template_uri")),
	}
}

func resourceTemplateAnalysisVariablesCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/template-analyses/variables"
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate UUID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId,
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildTemplateAnalysisVariablesBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error analyzing RFS template variables: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(requestId)

	variables := utils.PathSearch("variables", respBody, make([]interface{}, 0)).([]interface{})
	if err := d.Set("variables", flattenVariablesAttributes(variables)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenVariablesAttributes(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"name":        utils.PathSearch("name", v, nil),
			"type":        utils.PathSearch("type", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"default":     flattenVariablesDefaultAttributes(utils.PathSearch("default", v, nil)),
			"sensitive":   utils.PathSearch("sensitive", v, nil),
			"nullable":    utils.PathSearch("nullable", v, nil),
			"validations": flattenVariablesValidationsAttributes(utils.PathSearch("validations", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenVariablesValidationsAttributes(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"condition":     utils.PathSearch("condition", v, nil),
			"error_message": utils.PathSearch("error_message", v, nil),
		})
	}

	return rst
}

func flattenVariablesDefaultAttributes(respValue interface{}) string {
	if respValue == nil {
		return ""
	}

	b, err := json.Marshal(respValue)
	if err != nil {
		log.Printf("[ERROR] failed to marshal RFS template analysis variables default value: %v", err)
		return ""
	}

	return string(b)
}

func resourceTemplateAnalysisVariablesRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Read()' method because resource is a one-time action resource.
	return nil
}

func resourceTemplateAnalysisVariablesUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Update()' method because resource is a one-time action resource.
	return nil
}

func resourceTemplateAnalysisVariablesDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to analyze template variables. Deleting this resource
    will not cancel the analysis operation, but will only remove resource information from
    tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
