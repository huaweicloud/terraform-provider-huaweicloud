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

var migrateDomainNonUpdatableParams = []string{
	"enterprise_project_id",
	"target_enterprise_project_id",
	"host_ids",
	"policy_id",
	"certificate_id",
}

// @API WAF POST /v1/{project_id}/composite-waf/hosts/migration
func ResourceMigrateDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMigrateDomainCreate,
		ReadContext:   resourceMigrateDomainRead,
		UpdateContext: resourceMigrateDomainUpdate,
		DeleteContext: resourceMigrateDomainDelete,

		CustomizeDiff: config.FlexibleForceNew(migrateDomainNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"certificate_id": {
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

func buildMigrateDomainBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"host_ids":       utils.ExpandToStringList(d.Get("host_ids").([]interface{})),
		"policy_id":      d.Get("policy_id"),
		"certificate_id": utils.ValueIgnoreEmpty(d.Get("certificate_id")),
	}

	return bodyParams
}

func resourceMigrateDomainCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		epsId       = d.Get("enterprise_project_id").(string)
		targetEpsId = d.Get("target_enterprise_project_id").(string)
		httpUrl     = "v1/{project_id}/composite-waf/hosts/migration"
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath += fmt.Sprintf("?enterprise_project_id=%s", epsId)
	createPath += fmt.Sprintf("&target_enterprise_project_id=%s", targetEpsId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildMigrateDomainBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error migrating the domain: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	policyId := utils.PathSearch("policy_id", respBody, "").(string)
	if policyId == "" {
		return diag.Errorf("unable to find the policy ID from the API response")
	}

	d.SetId(policyId)

	return nil
}

func resourceMigrateDomainRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a action resource.
	return nil
}

func resourceMigrateDomainUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a action resource.
	return nil
}

func resourceMigrateDomainDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a action resource.
	return nil
}
