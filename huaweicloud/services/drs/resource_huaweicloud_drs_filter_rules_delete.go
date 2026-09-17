package drs

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

var drsFilterRulesDeleteNonUpdatableParams = []string{
	"job_id",
	"rule_ids",
	"type",
}

// @API DRS DELETE /v5/{project_id}/jobs/{job_id}/filter-rules
func ResourceDrsFilterRulesDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDrsFilterRulesDeleteCreate,
		ReadContext:   resourceDrsFilterRulesDeleteRead,
		UpdateContext: resourceDrsFilterRulesDeleteUpdate,
		DeleteContext: resourceDrsFilterRulesDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(drsFilterRulesDeleteNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the ID of the DRS job to delete filter rules from.",
			},
			"rule_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Specifies the list of filter rule IDs to delete.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"common", "config"}, false),
				Description:  "Specifies the type of the data filter rules to delete.",
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

func resourceDrsFilterRulesDeleteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/jobs/{job_id}/filter-rules"
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{job_id}", d.Get("job_id").(string))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	deleteOpt.JSONBody = utils.RemoveNil(buildDeleteDrsFilterRulesBodyParams(d))

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting DRS filter rules: %s", err)
	}

	_, err = utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// The API returns an empty response body, use job_id as the resource ID
	d.SetId(d.Get("job_id").(string))

	return resourceDrsFilterRulesDeleteRead(ctx, d, meta)
}

func buildDeleteDrsFilterRulesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"rule_ids": utils.ValueIgnoreEmpty(d.Get("rule_ids").([]interface{})),
		"type":     d.Get("type"),
	}
	return bodyParams
}

func resourceDrsFilterRulesDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrsFilterRulesDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrsFilterRulesDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DRS filter rules delete resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
