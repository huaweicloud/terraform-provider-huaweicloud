package cfw

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

// @API CFW POST /v1/{project_id}/advanced-ips-rule
func ResourceAdvancedIpsRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdvancedIpsRuleCreate,
		ReadContext:   resourceAdvancedIpsRuleRead,
		UpdateContext: resourceAdvancedIpsRuleUpdate,
		DeleteContext: resourceAdvancedIpsRuleDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"action",
			"ips_rule_id",
			"ips_rule_type",
			"object_id",
			"param",
			"status",
			"fw_instance_id",
			"enterprise_project_id",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"action": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ips_rule_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ips_rule_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"object_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"param": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"fw_instance_id": {
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

func buildAdvancedIpsRuleQueryParams(d *schema.ResourceData) string {
	rst := ""

	if v, ok := d.GetOk("fw_instance_id"); ok {
		rst += fmt.Sprintf("&fw_instance_id=%s", v.(string))
	}

	if v, ok := d.GetOk("enterprise_project_id"); ok {
		rst += fmt.Sprintf("&enterprise_project_id=%s", v.(string))
	}

	if rst != "" {
		rst = "?" + rst[1:]
	}

	return rst
}

func buildAdvancedIpsRuleCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"action":        d.Get("action"),
		"ips_rule_id":   d.Get("ips_rule_id"),
		"ips_rule_type": d.Get("ips_rule_type"),
		"object_id":     d.Get("object_id"),
		"param":         d.Get("param"),
		"status":        d.Get("status"),
	}
}

func resourceAdvancedIpsRuleCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/advanced-ips-rule"
	)

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAdvancedIpsRuleQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildAdvancedIpsRuleCreateBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating CFW advanced IPS rule: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error parsing response: %s", err)
	}

	id := utils.PathSearch("data.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CFW advanced IPS rule: ID is not found in API response")
	}
	d.SetId(id)

	return nil
}

func resourceAdvancedIpsRuleRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceAdvancedIpsRuleUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceAdvancedIpsRuleDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to create advanced IPS rule. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
