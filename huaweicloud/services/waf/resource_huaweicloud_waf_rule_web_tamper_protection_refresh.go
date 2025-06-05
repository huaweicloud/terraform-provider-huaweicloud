package waf

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF POST /v1/{project_id}/waf/policy/{policy_id}/antitamper/{rule_id}/refresh

var refreshCacheNonUpdatableParams = []string{"policy_id", "rule_id", "enterprise_project_id"}

func ResourceRuleWebTamperRefresh() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWebTamperRefreshCreate,
		ReadContext:   resourceWebTamperRefreshRead,
		UpdateContext: resourceWebTamperRefreshUpdate,
		DeleteContext: resourceWebTamperRefreshDelete,

		CustomizeDiff: config.FlexibleForceNew(refreshCacheNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
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

func resourceWebTamperRefreshCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		policyId = d.Get("policy_id").(string)
		ruleId   = d.Get("rule_id").(string)
		epsId    = cfg.GetEnterpriseProjectID(d)
		httpUrl  = "v1/{project_id}/waf/policy/{policy_id}/antitamper/{rule_id}/refresh"
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{policy_id}", policyId)
	createPath = strings.ReplaceAll(createPath, "{rule_id}", ruleId)
	if epsId != "" {
		createPath += fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error updating cache of web tamper protection rule: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	tamperId := utils.PathSearch("id", respBody, "").(string)
	if tamperId == "" {
		return diag.Errorf("unable to find the ID from the API response")
	}

	d.SetId(tamperId)

	return nil
}

func resourceWebTamperRefreshRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a action resource.
	return nil
}

func resourceWebTamperRefreshUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a action resource.
	return nil
}

func resourceWebTamperRefreshDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a action resource.
	return nil
}
