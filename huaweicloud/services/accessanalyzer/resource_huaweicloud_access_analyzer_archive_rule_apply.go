package accessanalyzer

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API AccessAnalyzer POST /v5/analyzers/{analyzer_id}/archive-rules/{archive_rule_id}/apply
var nonUpdatableParamsArchiveRuleApply = []string{"analyzer_id", "archive_rule_id"}

func ResourceArchiveRuleApply() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArchiveRuleApplyCreate,
		ReadContext:   resourceArchiveRuleApplyRead,
		UpdateContext: resourceArchiveRuleApplyUpdate,
		DeleteContext: resourceArchiveRuleApplyDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsArchiveRuleApply),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"analyzer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"archive_rule_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArchiveRuleApplyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("accessanalyzer", region)
	if err != nil {
		return diag.Errorf("error creating Access Analyzer client: %s", err)
	}

	createArchiveRuleApplyHttpUrl := "v5/analyzers/{analyzer_id}/archive-rules/{archive_rule_id}/apply"
	createArchiveRuleApplyPath := client.Endpoint + createArchiveRuleApplyHttpUrl
	createArchiveRuleApplyPath = strings.ReplaceAll(createArchiveRuleApplyPath, "{analyzer_id}", d.Get("analyzer_id").(string))
	createArchiveRuleApplyPath = strings.ReplaceAll(createArchiveRuleApplyPath, "{archive_rule_id}", d.Get("archive_rule_id").(string))
	createArchiveRuleApplyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("POST", createArchiveRuleApplyPath, &createArchiveRuleApplyOpt)
	if err != nil {
		return diag.Errorf("error applying Access Analyzer archive rule: %s", err)
	}

	d.SetId(d.Get("archive_rule_id").(string))

	return resourceArchiveRuleApplyRead(ctx, d, meta)
}

func resourceArchiveRuleApplyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceArchiveRuleApplyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceArchiveRuleApplyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for archive rule apply resource. Deleting this resource will
			not change the status of the currently archive rule apply resource, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
