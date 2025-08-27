package aad

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD POST /v2/aad/policies/waf/blackwhite-list
// @API AAD DELETE /v2/aad/policies/waf/blackwhite-list
// @API AAD GET /v2/aad/policies/waf/blackwhite-list
func ResourcePolicyBlackWhiteRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyBlackWhiteRuleCreate,
		UpdateContext: resourcePolicyBlackWhiteRuleUpdate,
		ReadContext:   resourcePolicyBlackWhiteRuleRead,
		DeleteContext: resourcePolicyBlackWhiteRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyBlackWhiteRuleImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{"domain_name", "ip", "overseas_type", "type"}),

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the domain name.`,
			},
			"overseas_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the protection area.`,
			},
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the IP address or IP segment.`,
			},
			"type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the rule type.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain ID.`,
			},
		},
	}
}

func buildCreatePolicyBlackWhiteRuleBodyParams(domainName string, overseasType int, ip string,
	ruleType int) map[string]interface{} {
	return map[string]interface{}{
		"domain_name":   domainName,
		"ips":           []interface{}{ip},
		"overseas_type": overseasType,
		"type":          ruleType,
	}
}

func buildPolicyBlackWhiteRuleQueryParams(domainName string, overseasType int) string {
	return fmt.Sprintf("?domain_name=%v&overseas_type=%v", domainName, overseasType)
}

func GetPolicyBlackWhiteRule(client *golangsdk.ServiceClient, domainName string, overseasType int, ip string,
	ruleType int) (interface{}, error) {
	requestPath := client.Endpoint + "v2/aad/policies/waf/blackwhite-list"
	requestPath += buildPolicyBlackWhiteRuleQueryParams(domainName, overseasType)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	expression := fmt.Sprintf("[black, white][] | [?ip == '%s' && type == `%d`] | [0]", ip, ruleType)
	ruleResp := utils.PathSearch(expression, respBody, nil)

	if ruleResp == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return ruleResp, nil
}

func resourcePolicyBlackWhiteRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		product      = "aad"
		httpUrl      = "v2/aad/policies/waf/blackwhite-list"
		domainName   = d.Get("domain_name").(string)
		overseasType = d.Get("overseas_type").(int)
		ip           = d.Get("ip").(string)
		ruleType     = d.Get("type").(int)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: buildCreatePolicyBlackWhiteRuleBodyParams(domainName, overseasType, ip, ruleType),
	}

	_, err = client.Request("POST", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error creating AAD policy black white rule: %s", err)
	}

	ruleResp, err := GetPolicyBlackWhiteRule(client, domainName, overseasType, ip, ruleType)
	if err != nil {
		return diag.Errorf("error creating AAD policy black white rule: target rule not found in query API response: %s", err)
	}

	ruleID := utils.PathSearch("id", ruleResp, "").(string)
	if ruleID == "" {
		return diag.Errorf("error creating AAD policy black white rule: ID is not found in query API response")
	}

	d.SetId(ruleID)

	return resourcePolicyBlackWhiteRuleRead(ctx, d, meta)
}

func resourcePolicyBlackWhiteRuleUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourcePolicyBlackWhiteRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		product      = "aad"
		domainName   = d.Get("domain_name").(string)
		overseasType = d.Get("overseas_type").(int)
		ip           = d.Get("ip").(string)
		ruleType     = d.Get("type").(int)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	ruleResp, err := GetPolicyBlackWhiteRule(client, domainName, overseasType, ip, ruleType)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AAD policy black white rule")
	}

	// Backfill the value of ID. Valid in import operation scenarios.
	ruleID := utils.PathSearch("id", ruleResp, "").(string)
	if ruleID != "" {
		d.SetId(ruleID)
	}

	mErr := multierror.Append(
		d.Set("type", utils.PathSearch("type", ruleResp, nil)),
		d.Set("ip", utils.PathSearch("ip", ruleResp, nil)),
		d.Set("domain_id", utils.PathSearch("domain_id", ruleResp, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDeletePolicyBlackWhiteRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"ids":           []interface{}{d.Id()},
		"domain_name":   d.Get("domain_name"),
		"overseas_type": d.Get("overseas_type"),
	}
}

func resourcePolicyBlackWhiteRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v2/aad/policies/waf/blackwhite-list"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: buildDeletePolicyBlackWhiteRuleBodyParams(d),
	}

	_, err = client.Request("DELETE", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error deleting AAD policy black white rule: %s", err)
	}

	return nil
}

func resourcePolicyBlackWhiteRuleImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 4 {
		return nil, errors.New("invalid format specified for import ID, must be <domain_name>/<overseas_type>/<type>/<ip>")
	}

	mErr := multierror.Append(
		d.Set("domain_name", parts[0]),
		d.Set("overseas_type", convertStringtoInt(parts[1])),
		d.Set("ip", parts[2]),
		d.Set("type", convertStringtoInt(parts[3])),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}

func convertStringtoInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("[ERROR] convert the string %s to int failed.", s)
	}
	return i
}
