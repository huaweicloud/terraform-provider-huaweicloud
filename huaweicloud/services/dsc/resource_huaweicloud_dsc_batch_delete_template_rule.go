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

var batchDeleteTemplateRuleNonUpdatableParams = []string{"template_id", "rule_ids"}

// @API DSC DELETE /v1/{project_id}/scan-templates/{template_id}/scan-rules/{rule_ids}
func ResourceBatchDeleteTemplateRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBatchDeleteTemplateRuleCreate,
		ReadContext:   resourceBatchDeleteTemplateRuleRead,
		UpdateContext: resourceBatchDeleteTemplateRuleUpdate,
		DeleteContext: resourceBatchDeleteTemplateRuleDelete,

		CustomizeDiff: config.FlexibleForceNew(batchDeleteTemplateRuleNonUpdatableParams),

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
			"rule_ids": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the rule IDs to be deleted.",
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

func resourceBatchDeleteTemplateRuleCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		templateId = d.Get("template_id").(string)
		ruleIds    = d.Get("rule_ids").(string)
		httpUrl    = "v1/{project_id}/scan-templates/{template_id}/scan-rules/{rule_ids}"
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{template_id}", templateId)
	requestPath = strings.ReplaceAll(requestPath, "{rule_ids}", ruleIds)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch deleting DSC template rules: %s", err)
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

func resourceBatchDeleteTemplateRuleRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchDeleteTemplateRuleUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchDeleteTemplateRuleDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch delete DSC template rule associations.
Deleting this resource will not restore the deleted rule associations or undo the delete action, but will only
remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
