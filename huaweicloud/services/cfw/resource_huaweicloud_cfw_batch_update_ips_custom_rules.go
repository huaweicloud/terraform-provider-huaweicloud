package cfw

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

var batchUpdateIpsCustomRulesActionNonUpdatableParams = []string{
	"fw_instance_id", "ips_ids", "action_type"}

// @API CFW POST /v1/{project_id}/ips/custom-rule/action
func ResourceBatchUpdateIpsCustomRulesAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBatchUpdateIpsCustomRulesActionCreate,
		ReadContext:   resourceBatchUpdateIpsCustomRulesActionRead,
		UpdateContext: resourceBatchUpdateIpsCustomRulesActionUpdate,
		DeleteContext: resourceBatchUpdateIpsCustomRulesActionDelete,

		CustomizeDiff: config.FlexibleForceNew(batchUpdateIpsCustomRulesActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ips_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildBatchUpdateIpsCustomRulesActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"fw_instance_id": d.Get("fw_instance_id"),
		"action_type":    d.Get("action_type"),
		"ips_ids":        utils.ExpandToStringList(d.Get("ips_ids").([]interface{})),
	}
}

func resourceBatchUpdateIpsCustomRulesActionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		fwInstanceId = d.Get("fw_instance_id").(string)
		httpUrl      = "v1/{project_id}/ips/custom-rule/action"
		product      = "cfw"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildBatchUpdateIpsCustomRulesActionBodyParams(d),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch updating CFW IPS custom rules action: %s", err)
	}

	d.SetId(fwInstanceId)

	return nil
}

func resourceBatchUpdateIpsCustomRulesActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchUpdateIpsCustomRulesActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchUpdateIpsCustomRulesActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch update IPS custom rule actions. 
Deleting this resource will not revert the actions of the IPS rules, but will only remove the resource 
information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
