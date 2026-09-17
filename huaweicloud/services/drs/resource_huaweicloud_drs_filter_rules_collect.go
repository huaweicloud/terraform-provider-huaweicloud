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

var drsFilterRulesCollectNonUpdatableParams = []string{
	"job_id",
	"filter_type",
	"is_all_rules",
	"limit",
	"offset",
}

// @API DRS POST /v5/{project_id}/jobs/{job_id}/filter-rules
func ResourceDrsFilterRulesCollect() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDrsFilterRulesCollectCreate,
		ReadContext:   resourceDrsFilterRulesCollectRead,
		UpdateContext: resourceDrsFilterRulesCollectUpdate,
		DeleteContext: resourceDrsFilterRulesCollectDelete,

		CustomizeDiff: config.FlexibleForceNew(drsFilterRulesCollectNonUpdatableParams),

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
				Description: "Specifies the ID of the DRS job to collect filter rules from.",
			},
			"filter_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"common", "config"}, false),
				Description:  "Specifies the type of the data filter rules to collect.",
			},
			"is_all_rules": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies whether to query all filter rules (including synchronized and unsent ones).",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the number of records to return.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the offset of the records to query.",
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

func resourceDrsFilterRulesCollectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{job_id}", d.Get("job_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateDrsFilterRulesCollectBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error collecting DRS filter rules: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error collecting DRS filter rules: ID is not found in API response")
	}
	d.SetId(id)

	return resourceDrsFilterRulesCollectRead(ctx, d, meta)
}

func buildCreateDrsFilterRulesCollectBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"filter_type":  d.Get("filter_type"),
		"is_all_rules": utils.ValueIgnoreEmpty(d.Get("is_all_rules")),
		"limit":        utils.ValueIgnoreEmpty(d.Get("limit")),
		"offset":       utils.ValueIgnoreEmpty(d.Get("offset")),
	}
	return bodyParams
}

func resourceDrsFilterRulesCollectRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrsFilterRulesCollectUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrsFilterRulesCollectDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DRS filter rules collect resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
