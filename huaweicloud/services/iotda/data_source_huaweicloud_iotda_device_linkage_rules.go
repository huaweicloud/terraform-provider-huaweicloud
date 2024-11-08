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

// @API IoTDA GET /v5/iot/{project_id}/rules
func DataSourceDeviceLinkageRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeviceLinkageRulesRead,

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
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"space_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"triggers": {
							Type:     schema.TypeList,
							Elem:     ruleConditionSchema(),
							Computed: true,
						},
						"actions": {
							Type:     schema.TypeList,
							Elem:     ruleActionSchema(),
							Computed: true,
						},
						"effective_period": {
							Type:     schema.TypeList,
							Elem:     ruleTimeRangeSchema(),
							Computed: true,
						},
						"trigger_logic": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func ruleConditionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_data_condition": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"in_values": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"trigger_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_validatiy_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"simple_timer_condition": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repeat_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"repeat_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"daily_timer_condition": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"days_of_week": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"device_linkage_status_condition": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_list": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
	return &sc
}

func ruleActionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_command": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command_body": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"buffer_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"response_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"smn_forwarding": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topic_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topic_urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"device_alarm": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dimension": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
	return &sc
}

func ruleTimeRangeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"days_of_week": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceDeviceLinkageRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	var (
		allRules []model.RuleResponse
		limit    = int32(50)
		offset   int32
	)

	for {
		listOpts := model.ListRulesRequest{
			AppId:    utils.StringIgnoreEmpty(d.Get("space_id").(string)),
			RuleType: utils.StringIgnoreEmpty(d.Get("type").(string)),
			Limit:    utils.Int32(limit),
			Offset:   &offset,
		}

		listResp, listErr := client.ListRules(&listOpts)
		if listErr != nil {
			return diag.Errorf("error querying IoTDA device linkage rules: %s", listErr)
		}

		if listResp == nil || listResp.Rules == nil {
			break
		}

		if len(*listResp.Rules) == 0 {
			break
		}

		allRules = append(allRules, *listResp.Rules...)
		//nolint:gosec
		offset += int32(len(*listResp.Rules))
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuId)

	targetRules := filterListDeviceLinkageRules(allRules, d)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenDeviceLinkageRules(targetRules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterListDeviceLinkageRules(rules []model.RuleResponse, d *schema.ResourceData) []model.RuleResponse {
	if len(rules) == 0 {
		return nil
	}

	rst := make([]model.RuleResponse, 0, len(rules))
	for _, v := range rules {
		if ruleId, ok := d.GetOk("rule_id"); ok &&
			fmt.Sprint(ruleId) != utils.StringValue(v.RuleId) {
			continue
		}

		if ruleName, ok := d.GetOk("name"); ok &&
			fmt.Sprint(ruleName) != v.Name {
			continue
		}

		if status, ok := d.GetOk("status"); ok &&
			fmt.Sprint(status) != utils.StringValue(v.Status) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenDeviceLinkageRules(rules []model.RuleResponse) []interface{} {
	if len(rules) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rules))
	for _, v := range rules {
		actions := v.Actions
		rst = append(rst, map[string]interface{}{
			"id":               v.RuleId,
			"name":             v.Name,
			"description":      v.Description,
			"type":             v.RuleType,
			"status":           v.Status,
			"space_id":         v.AppId,
			"triggers":         flattenLinkageTriggers(v.ConditionGroup),
			"actions":          flattenLinkageActions(&actions),
			"effective_period": flattenLinkageEffectivePeriod(v.ConditionGroup),
			"trigger_logic":    v.ConditionGroup.Logic,
			"updated_at":       v.LastUpdateTime,
		})
	}

	return rst
}
