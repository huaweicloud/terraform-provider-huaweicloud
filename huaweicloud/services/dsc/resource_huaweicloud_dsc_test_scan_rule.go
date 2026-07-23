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

var testScanRuleNonUpdatableParams = []string{
	"category",
	"data",
	"effective_mode",
	"location",
	"rule_content",
	"rule_id",
	"rule_name",
}

// @API DSC POST /v1/{project_id}/scan-rules/test
func ResourceTestScanRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTestScanRuleCreate,
		ReadContext:   resourceTestScanRuleRead,
		UpdateContext: resourceTestScanRuleUpdate,
		DeleteContext: resourceTestScanRuleDelete,

		CustomizeDiff: config.FlexibleForceNew(testScanRuleNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"effective_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_content": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"is_match": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"match_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildTestScanRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"category":       utils.ValueIgnoreEmpty(d.Get("category")),
		"data":           utils.ValueIgnoreEmpty(d.Get("data")),
		"effective_mode": utils.ValueIgnoreEmpty(d.Get("effective_mode")),
		"location":       utils.ValueIgnoreEmpty(d.Get("location")),
		"rule_content":   utils.ValueIgnoreEmpty(d.Get("rule_content").([]interface{})),
		"rule_id":        utils.ValueIgnoreEmpty(d.Get("rule_id")),
		"rule_name":      utils.ValueIgnoreEmpty(d.Get("rule_name")),
	}
}

func resourceTestScanRuleCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/scan-rules/test"
		product = "dsc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildTestScanRuleBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error testing DSC scan rule: %s", err)
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
		d.Set("is_match", utils.PathSearch("is_match", respBody, nil)),
		d.Set("match_group", utils.PathSearch("match_group", respBody, nil)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting test DSC scan rule fields: %s", mErr)
	}

	return nil
}

func resourceTestScanRuleRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceTestScanRuleUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceTestScanRuleDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to test a DSC scan rule. Deleting this resource will
    not undo the test action, but will only remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
