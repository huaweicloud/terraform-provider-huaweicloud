package iotda

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA GET /v5/iot/{project_id}/routing-rule/rules
func DataSourceDataForwardingRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataForwardingRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"trigger": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"space_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"select": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"where": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildDataForwardingRulesQueryParams(d *schema.ResourceData) string {
	rst := ""
	if v, ok := d.GetOk("name"); ok {
		rst += fmt.Sprintf("&rule_name=%v", v)
	}

	if v, ok := d.GetOk("resource"); ok {
		rst += fmt.Sprintf("&resource=%v", v)
	}

	if v, ok := d.GetOk("trigger"); ok {
		rst += fmt.Sprintf("&event=%v", v)
	}

	if v, ok := d.GetOk("app_type"); ok {
		rst += fmt.Sprintf("&app_type=%v", v)
	}

	if v, ok := d.GetOk("space_id"); ok {
		rst += fmt.Sprintf("&app_id=%v", v)
	}

	if v, ok := d.GetOk("enabled"); ok {
		rst += fmt.Sprintf("&active=%v", v)
	}

	return rst
}

func dataSourceDataForwardingRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v5/iot/{project_id}/routing-rule/rules?limit=50"
		product  = "iotda"
		allRules []interface{}
		offset   = 0
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildDataForwardingRulesQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%d", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error querying IoTDA data forwarding rules: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		rules := utils.PathSearch("rules", respBody, make([]interface{}, 0)).([]interface{})
		if len(rules) == 0 {
			break
		}

		allRules = append(allRules, rules...)
		offset += len(rules)
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuId)

	targetRoutingRules := filterListDataForwardingRules(allRules, d)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenDataForwardingRules(targetRoutingRules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterListDataForwardingRules(allRules []interface{}, d *schema.ResourceData) []interface{} {
	if len(allRules) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(allRules))
	for _, v := range allRules {
		if ruleId, ok := d.GetOk("rule_id"); ok &&
			ruleId.(string) != utils.PathSearch("rule_id", v, "").(string) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenDataForwardingRules(allRules []interface{}) []interface{} {
	if len(allRules) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(allRules))
	for _, v := range allRules {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("rule_id", v, nil),
			"name":        utils.PathSearch("rule_name", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"resource":    utils.PathSearch("subject.resource", v, nil),
			"trigger":     utils.PathSearch("subject.event", v, nil),
			"app_type":    utils.PathSearch("app_type", v, nil),
			"space_id":    utils.PathSearch("app_id", v, nil),
			"enabled":     utils.PathSearch("active", v, nil),
			"select":      utils.PathSearch("select", v, nil),
			"where":       utils.PathSearch("where", v, nil),
		})
	}

	return rst
}
