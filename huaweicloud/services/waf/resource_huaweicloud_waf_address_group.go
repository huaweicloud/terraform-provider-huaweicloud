// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product WAF
// ---------------------------------------------------------------

package waf

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF GET /v1/{project_id}/waf/ip-group/{id}
// @API WAF PUT /v1/{project_id}/waf/ip-group/{id}
// @API WAF DELETE /v1/{project_id}/waf/ip-group/{id}
// @API WAF POST /v1/{project_id}/waf/ip-groups
func ResourceWafAddressGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAddressGroupCreate,
		UpdateContext: resourceAddressGroupUpdate,
		ReadContext:   resourceAddressGroupRead,
		DeleteContext: resourceAddressGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWAFImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the address group.`,
			},
			"ip_addresses": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `Specifies the IP addresses or IP address ranges.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the enterprise project ID of WAF address group.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the description of the address group.`,
			},
			"rules": {
				Type:        schema.TypeList,
				Elem:        rulesSchema(),
				Computed:    true,
				Description: `The list of rules that use the IP address group.`,
			},
		},
	}
}

func rulesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of rule.`,
			},
			"rule_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of rule.`,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of policy.`,
			},
			"policy_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of policy.`,
			},
		},
	}
	return &sc
}

func resourceAddressGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/ip-groups"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildWAFAddressGroupQueryParams(d, cfg)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildWAFAddressGroupBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating WAF address group: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	groupId := utils.PathSearch("id", respBody, "").(string)
	if groupId == "" {
		return diag.Errorf("error creating WAF address group: ID is not found in API response")
	}
	d.SetId(groupId)

	return resourceAddressGroupRead(ctx, d, meta)
}

func buildWAFAddressGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	ipAddresses := d.Get("ip_addresses").([]interface{})
	addresses := make([]string, 0, len(ipAddresses))
	for _, v := range ipAddresses {
		addresses = append(addresses, v.(string))
	}

	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"ips":         strings.Join(addresses, ","),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func buildWAFAddressGroupQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	epsId := cfg.GetEnterpriseProjectID(d)
	if epsId == "" {
		return ""
	}
	return fmt.Sprintf("?enterprise_project_id=%s", epsId)
}

func resourceAddressGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/{project_id}/waf/ip-group/{id}"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{id}", d.Id())
	requestPath += buildWAFAddressGroupQueryParams(d, cfg)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		// If the address group does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving WAF address group")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("ip_addresses", flattenAddressGroupResponseBodyIpAddresses(respBody)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("rules", flattenAddressGroupResponseBodyRules(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAddressGroupResponseBodyIpAddresses(resp interface{}) []string {
	if resp == nil {
		return nil
	}
	ips := utils.PathSearch("ips", resp, "").(string)
	return strings.Split(ips, ",")
}

func flattenAddressGroupResponseBodyRules(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("rules", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"rule_id":     utils.PathSearch("rule_id", v, nil),
			"rule_name":   utils.PathSearch("rule_name", v, nil),
			"policy_id":   utils.PathSearch("policy_id", v, nil),
			"policy_name": utils.PathSearch("policy_name", v, nil),
		})
	}
	return rst
}

func resourceAddressGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateWAFAddressGroupChanges := []string{
		"name",
		"ip_addresses",
		"description",
	}

	if d.HasChanges(updateWAFAddressGroupChanges...) {
		var (
			httpUrl = "v1/{project_id}/waf/ip-group/{id}"
			product = "waf"
		)
		client, err := cfg.NewServiceClient(product, region)
		if err != nil {
			return diag.Errorf("error creating WAF client: %s", err)
		}

		requestPath := client.Endpoint + httpUrl
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestPath = strings.ReplaceAll(requestPath, "{id}", d.Id())
		requestPath += buildWAFAddressGroupQueryParams(d, cfg)
		requestOpt := golangsdk.RequestOpts{
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=utf8",
			},
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildWAFAddressGroupBodyParams(d)),
		}

		_, err = client.Request("PUT", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating WAF address group: %s", err)
		}
	}
	return resourceAddressGroupRead(ctx, d, meta)
}

func resourceAddressGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/ip-group/{id}"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{id}", d.Id())
	requestPath += buildWAFAddressGroupQueryParams(d, cfg)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		// If the address group does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting WAF address group")
	}

	return nil
}

// resourceWAFImportState use to import an id with format <id> or <id>/<enterprise_project_id>
func resourceWAFImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	if !strings.Contains(d.Id(), "/") {
		return []*schema.ResourceData{d}, nil
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <id>/<enterprise_project_id>")
	}
	d.SetId(parts[0])
	return []*schema.ResourceData{d}, d.Set("enterprise_project_id", parts[1])
}
