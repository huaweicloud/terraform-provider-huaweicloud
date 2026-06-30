package elb

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB GET /v3/{project_id}/elb/l7policies/rules
func DataSourceAllL7Rules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAllL7RulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rule_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"compare_type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"provisioning_status": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"value": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     allL7RulesRuleSchema(),
			},
		},
	}
}

func allL7RulesRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compare_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"invert": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"conditions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     allL7RulesRuleConditionSchema(),
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func allL7RulesRuleConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAllL7RulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v3/{project_id}/elb/l7policies/rules?limit=2000"
	)

	client, err := cfg.NewServiceClient("elb", region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildListAllL7RulesQueryParams(d, epsId)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving the ELB l7rules: %s", err)
	}

	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	datasourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(datasourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("rules", flattenListAllL7RulesBody(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListAllL7RulesQueryParams(d *schema.ResourceData, epsId string) string {
	res := ""

	if v, ok := d.GetOk("rule_id"); ok {
		ruleIds := v.([]interface{})
		for _, ruleId := range ruleIds {
			res = fmt.Sprintf("%s&id=%v", res, ruleId)
		}
	}

	if v, ok := d.GetOk("type"); ok {
		matchTypes := v.([]interface{})
		for _, matchType := range matchTypes {
			res = fmt.Sprintf("%s&type=%v", res, matchType)
		}
	}

	if v, ok := d.GetOk("compare_type"); ok {
		compareTypes := v.([]interface{})
		for _, compareType := range compareTypes {
			res = fmt.Sprintf("%s&compare_type=%v", res, compareType)
		}
	}

	if v, ok := d.GetOk("provisioning_status"); ok {
		provisionStatus := v.([]interface{})
		for _, status := range provisionStatus {
			res = fmt.Sprintf("%s&provisioning_status=%v", res, status)
		}
	}

	if v, ok := d.GetOk("value"); ok {
		values := v.([]interface{})
		for _, value := range values {
			res = fmt.Sprintf("%s&value=%v", res, value)
		}
	}

	if epsId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsId)
	}

	return res
}

func flattenListAllL7RulesBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("rules", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":                  utils.PathSearch("id", v, nil),
			"type":                utils.PathSearch("type", v, nil),
			"compare_type":        utils.PathSearch("compare_type", v, nil),
			"key":                 utils.PathSearch("key", v, nil),
			"value":               utils.PathSearch("value", v, nil),
			"provisioning_status": utils.PathSearch("provisioning_status", v, nil),
			"invert":              utils.PathSearch("invert", v, nil),
			"conditions":          flattenListAllL7RulesConditionsBody(v),
			"project_id":          utils.PathSearch("project_id", v, nil),
			"created_at":          utils.PathSearch("created_at", v, nil),
			"updated_at":          utils.PathSearch("updated_at", v, nil),
		})
	}
	return res
}

func flattenListAllL7RulesConditionsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("conditions", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return res
}
