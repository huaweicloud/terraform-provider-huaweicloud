package coc

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

var scriptAcceptNonUpdatableParams = []string{"script_uuid", "status", "comments"}

// @API COC POST /v1/job/scripts/{script_uuid}/action
func ResourceScriptApproval() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScriptApprovalCreate,
		ReadContext:   resourceScriptApprovalRead,
		UpdateContext: resourceScriptApprovalUpdate,
		DeleteContext: resourceScriptApprovalDelete,

		CustomizeDiff: config.FlexibleForceNew(scriptAcceptNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"script_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comments": {
				Type:     schema.TypeString,
				Optional: true,
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

func buildScriptApprovalCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"status":   d.Get("status"),
		"comments": utils.ValueIgnoreEmpty(d.Get("comments")),
	}

	return bodyParams
}

func resourceScriptApprovalCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/job/scripts/{script_uuid}/action"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	scriptUUID := d.Get("script_uuid").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{script_uuid}", scriptUUID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildScriptApprovalCreateOpts(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error approving the COC script (%s): %s", scriptUUID, err)
	}
	d.SetId(scriptUUID)

	return nil
}

func resourceScriptApprovalRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceScriptApprovalUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceScriptApprovalDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting script approve resource is not supported. The script approve resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
