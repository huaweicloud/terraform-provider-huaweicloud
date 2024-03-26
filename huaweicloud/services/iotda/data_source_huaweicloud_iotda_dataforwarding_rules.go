package iotda

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

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

func dataSourceDataForwardingRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	var (
		allRoutingRules []model.RoutingRule
		limit           = int32(50)
		offset          int32
	)

	for {
		listOpts := model.ListRoutingRulesRequest{
			RuleName: utils.StringIgnoreEmpty(d.Get("name").(string)),
			Resource: utils.StringIgnoreEmpty(d.Get("resource").(string)),
			Event:    utils.StringIgnoreEmpty(d.Get("trigger").(string)),
			AppType:  utils.StringIgnoreEmpty(d.Get("app_type").(string)),
			AppId:    utils.StringIgnoreEmpty(d.Get("space_id").(string)),
			Limit:    utils.Int32(limit),
			Offset:   &offset,
		}

		active, ok := d.GetOk("enabled")
		if ok {
			listOpts.Active = utils.StringToBool(active)
		}

		listResp, listErr := client.ListRoutingRules(&listOpts)
		if listErr != nil {
			return diag.Errorf("error querying IoTDA dataforwarding rules: %s", listErr)
		}

		if len(*listResp.Rules) == 0 {
			break
		}
		allRoutingRules = append(allRoutingRules, *listResp.Rules...)
		offset += limit
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuId)

	targetRoutingRules := filterListDataForwardingRules(allRoutingRules, d)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenDataForwardingRules(targetRoutingRules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterListDataForwardingRules(routingRules []model.RoutingRule, d *schema.ResourceData) []model.RoutingRule {
	if len(routingRules) == 0 {
		return nil
	}

	rst := make([]model.RoutingRule, 0, len(routingRules))
	for _, v := range routingRules {
		if ruleId, ok := d.GetOk("rule_id"); ok &&
			fmt.Sprint(ruleId) != utils.StringValue(v.RuleId) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenDataForwardingRules(routingRules []model.RoutingRule) []interface{} {
	if len(routingRules) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(routingRules))
	for _, v := range routingRules {
		rst = append(rst, map[string]interface{}{
			"id":          v.RuleId,
			"name":        v.RuleName,
			"description": v.Description,
			"resource":    v.Subject.Resource,
			"trigger":     v.Subject.Event,
			"app_type":    v.AppType,
			"space_id":    v.AppId,
			"enabled":     v.Active,
			"select":      v.Select,
			"where":       v.Where,
		})
	}

	return rst
}
