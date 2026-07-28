package dsc

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var dscBatchUpdateTemplateRulesClassificationNonUpdatableParams = []string{"template_id", "classification_id", "rule_id_list"}

// @API DSC PUT /v1/{project_id}/scan-templates/{template_id}/rules-classification-mapping
func ResourceDscBatchUpdateTemplateRulesClassification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDscBatchUpdateTemplateRulesClassificationCreate,
		ReadContext:   resourceDscBatchUpdateTemplateRulesClassificationRead,
		UpdateContext: resourceDscBatchUpdateTemplateRulesClassificationUpdate,
		DeleteContext: resourceDscBatchUpdateTemplateRulesClassificationDelete,

		CustomizeDiff: config.FlexibleForceNew(dscBatchUpdateTemplateRulesClassificationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the scan template ID.",
			},
			"classification_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the classification ID for batch updating the rule classification.",
			},
			"rule_id_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the list of rule IDs to be batch updated.",
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"msg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The returned message.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The returned status.",
			},
		},
	}
}

func resourceDscBatchUpdateTemplateRulesClassificationCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		templateId = d.Get("template_id").(string)
		httpUrl    = "v1/{project_id}/scan-templates/{template_id}/rules-classification-mapping"
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{template_id}", templateId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildDscBatchUpdateTemplateRulesClassificationBodyParams(d)),
	}

	resp, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch updating DSC template rules classification: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("msg", utils.PathSearch("msg", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDscBatchUpdateTemplateRulesClassificationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"classification_id": d.Get("classification_id"),
		"rule_id_list":      utils.ExpandToStringList(d.Get("rule_id_list").([]interface{})),
	}
	return bodyParams
}

func resourceDscBatchUpdateTemplateRulesClassificationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDscBatchUpdateTemplateRulesClassificationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDscBatchUpdateTemplateRulesClassificationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch update DSC template rules classification mapping.
Deleting this resource will not restore the updated classification mapping or undo the update action, but will only
remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
