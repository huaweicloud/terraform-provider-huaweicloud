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

// @API WAF POST /v1/{project_id}/waf/rule/ip-reputation
func ResourceWafBatchCreateIpReputationRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafBatchCreateIpReputationRulesCreate,
		ReadContext:   resourceWafBatchCreateIpReputationRulesRead,
		UpdateContext: resourceWafBatchCreateIpReputationRulesUpdate,
		DeleteContext: resourceWafBatchCreateIpReputationRulesDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"name",
			"type",
			"tags",
			"action",
			"action.*.category",
			"policy_ids",
			"description",
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
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"action": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     buildBatchCreateIpReputationRulesActionSchema(),
			},
			"policy_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Type:     schema.TypeString,
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

func buildBatchCreateIpReputationRulesActionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func buildBatchCreateIpReputationRulesQueryParams(epsId string) string {
	if epsId == "" {
		return ""
	}
	return fmt.Sprintf("?enterprise_project_id=%s", epsId)
}

// Empty array is allowed. SO the response type is interface{}
func buildBatchCreateIpReputationRulesActionBodyParam(d *schema.ResourceData) map[string]interface{} {
	rawArray, ok := d.Get("action").([]interface{})
	if !ok || len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"category": rawMap["category"],
	}
}

func buildBatchCreateIpReputationRulesBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"type":        d.Get("type"),
		"tags":        d.Get("tags"),
		"action":      buildBatchCreateIpReputationRulesActionBodyParam(d),
		"policy_ids":  d.Get("policy_ids"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func resourceWafBatchCreateIpReputationRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/rule/ip-reputation"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath += buildBatchCreateIpReputationRulesQueryParams(epsId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildBatchCreateIpReputationRulesBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error batch creating WAF IP reputation rules: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	return resourceWafBatchCreateIpReputationRulesRead(ctx, d, meta)
}

func resourceWafBatchCreateIpReputationRulesRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is an action resource.
	return nil
}

func resourceWafBatchCreateIpReputationRulesUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is an action resource.
	return nil
}

func resourceWafBatchCreateIpReputationRulesDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch create WAF IP reputation rules. Deleting this resource
    will not remove the created rules, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
