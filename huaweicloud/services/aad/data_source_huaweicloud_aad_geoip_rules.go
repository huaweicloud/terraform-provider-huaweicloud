package aad

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD GET /v2/aad/policies/waf/geoip-rule
func DataSourceGeoIpRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeoIpRulesRead,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the domain name to query.`,
			},
			"overseas_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the protection region.`,
			},
			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of Geo IP rules.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the rule.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the rule.`,
						},
						"geoip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The geographical location code.`,
						},
						"overseas_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The protection region.`,
						},
						"timestamp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The creation timestamp of the rule.`,
						},
						"white": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The protection action.`,
						},
					},
				},
			},
		},
	}
}

func buildGeoIpRulesQueryParams(d *schema.ResourceData) string {
	domainName := d.Get("domain_name").(string)
	oversightType := d.Get("overseas_type").(string)

	return fmt.Sprintf("?domain_name=%s&overseas_type=%s", domainName, oversightType)
}

func dataSourceGeoIpRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v2/aad/policies/waf/geoip-rule"
	)

	client, err := cfg.NewServiceClient("aad", region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += buildGeoIpRulesQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving AAD Geo IP rules: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("items", flattenGeoIpRules(utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGeoIpRules(rules []interface{}) []map[string]interface{} {
	if len(rules) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		result = append(result, map[string]interface{}{
			"id":            utils.PathSearch("id", rule, nil),
			"name":          utils.PathSearch("name", rule, nil),
			"geoip":         utils.PathSearch("geoip", rule, nil),
			"overseas_type": utils.PathSearch("overseas_type", rule, nil),
			"timestamp":     utils.PathSearch("timestamp", rule, nil),
			"white":         utils.PathSearch("white", rule, nil),
		})
	}

	return result
}
