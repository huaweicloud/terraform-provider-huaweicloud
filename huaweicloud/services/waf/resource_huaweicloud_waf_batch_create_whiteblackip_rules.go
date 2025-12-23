package waf

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF POST /v1/{project_id}/waf/rule/whiteblackip
func ResourceWafBatchCreateWhiteBlackIpRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafBatchCreateWhiteBlackIpRulesCreate,
		ReadContext:   resourceWafBatchCreateWhiteBlackIpRulesRead,
		UpdateContext: resourceWafBatchCreateWhiteBlackIpRulesUpdate,
		DeleteContext: resourceWafBatchCreateWhiteBlackIpRulesDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"name",
			"white",
			"policy_ids",
			"addr",
			"description",
			"ip_group_id",
			"time_mode",
			"start",
			"terminal",
			"enterprise_project_id",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"white": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"policy_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"addr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"time_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"terminal": {
				Type:     schema.TypeInt,
				Optional: true,
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

func buildBatchCreateWhiteBlackIpRulesQueryParams(epsId string) string {
	if epsId == "" {
		return ""
	}

	return fmt.Sprintf("?enterprise_project_id=%s", epsId)
}

func buildBatchCreateWhiteBlackIpRulesBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"white":       d.Get("white"),
		"policy_ids":  d.Get("policy_ids"),
		"addr":        utils.ValueIgnoreEmpty(d.Get("addr")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"ip_group_id": utils.ValueIgnoreEmpty(d.Get("ip_group_id")),
		"time_mode":   utils.ValueIgnoreEmpty(d.Get("time_mode")),
		"start":       utils.ValueIgnoreEmpty(d.Get("start")),
		"terminal":    utils.ValueIgnoreEmpty(d.Get("terminal")),
	}
}

func resourceWafBatchCreateWhiteBlackIpRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/rule/whiteblackip"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath += buildBatchCreateWhiteBlackIpRulesQueryParams(epsId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildBatchCreateWhiteBlackIpRulesBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error batch creating WAF whiteblack IP rules: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	return resourceWafBatchCreateWhiteBlackIpRulesRead(ctx, d, meta)
}

func resourceWafBatchCreateWhiteBlackIpRulesRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is an action resource.
	return nil
}

func resourceWafBatchCreateWhiteBlackIpRulesUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is an action resource.
	return nil
}

func resourceWafBatchCreateWhiteBlackIpRulesDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch create whiteblack IP rules. Deleting this resource
    will not remove the created rules, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
