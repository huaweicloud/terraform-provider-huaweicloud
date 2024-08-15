// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product WAF
// ---------------------------------------------------------------

package waf

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF PUT /v1/{project_id}/waf/policy/{policy_id}/{rule_type}/{rule_id}/status
// @API WAF POST /v1/{project_id}/waf/policy/{policy_id}/geoip
// @API WAF PUT /v1/{project_id}/waf/policy/{policy_id}/geoip/{rule_id}
// @API WAF DELETE /v1/{project_id}/waf/policy/{policy_id}/geoip/{rule_id}
// @API WAF GET /v1/{project_id}/waf/policy/{policy_id}/geoip/{rule_id}
func ResourceRuleGeolocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRuleGeolocationCreate,
		UpdateContext: resourceRuleGeolocationUpdate,
		ReadContext:   resourceRuleGeolocationRead,
		DeleteContext: resourceRuleGeolocationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWAFRuleImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the policy ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of WAF geolocation access control rule.`,
			},
			"geolocation": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the locations that can be configured in the geolocation access control rule.`,
			},
			"action": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the protective action of WAF geolocation access control rule.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the enterprise project ID of WAF geolocation access control rule.`,
			},
			"status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: `Specifies the status of WAF geolocation access control rule.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of WAF geolocation access control rule.`,
			},
		},
	}
}

func resourceRuleGeolocationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/policy/{policy_id}/geoip"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF Client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{policy_id}", d.Get("policy_id").(string))
	createPath += buildQueryParams(d, cfg)
	createOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody:         buildCreateBodyParams(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating WAF geolocation access control rule: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("id", createRespBody)
	if err != nil {
		return diag.Errorf("error creating WAF geolocation access control rule: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourceRuleGeolocationRead(ctx, d, meta)
}

func buildCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"geoip":       d.Get("geolocation"),
		"white":       d.Get("action"),
		"status":      d.Get("status"),
		"description": d.Get("description"),
	}
}

func resourceRuleGeolocationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/{project_id}/waf/policy/{policy_id}/geoip/{rule_id}"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF Client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{policy_id}", d.Get("policy_id").(string))
	getPath = strings.ReplaceAll(getPath, "{rule_id}", d.Id())
	getPath += buildQueryParams(d, cfg)
	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		// If the rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving WAF geolocation access control rule")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("policy_id", utils.PathSearch("policyid", getRespBody, nil)),
		d.Set("geolocation", utils.PathSearch("geoip", getRespBody, nil)),
		d.Set("action", utils.PathSearch("white", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceRuleGeolocationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF Client: %s", err)
	}

	updateRuleGeolocationChanges := []string{
		"name",
		"geolocation",
		"action",
		"description",
	}
	if d.HasChanges(updateRuleGeolocationChanges...) {
		updatePath := client.Endpoint + "v1/{project_id}/waf/policy/{policy_id}/geoip/{rule_id}"
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{policy_id}", d.Get("policy_id").(string))
		updatePath = strings.ReplaceAll(updatePath, "{rule_id}", d.Id())
		updatePath += buildQueryParams(d, cfg)
		updateOpt := golangsdk.RequestOpts{
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=utf8",
			},
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateBodyParams(d)),
		}

		_, err := client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating WAF geolocation access control rule: %s", err)
		}
	}

	if d.HasChange("status") {
		if err := updateRuleStatus(client, d, cfg, "geoip"); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceRuleGeolocationRead(ctx, d, meta)
}

func buildUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"geoip":       d.Get("geolocation"),
		"white":       d.Get("action"),
		"description": d.Get("description"),
	}
}

func resourceRuleGeolocationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/policy/{policy_id}/geoip/{rule_id}"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF Client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{policy_id}", d.Get("policy_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{rule_id}", d.Id())
	deletePath += buildQueryParams(d, cfg)
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// If the rule does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting WAF geolocation access control rule")
	}

	return nil
}
