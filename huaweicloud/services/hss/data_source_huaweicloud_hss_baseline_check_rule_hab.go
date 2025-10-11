package hss

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

// @API HSS POST /v5/{project_id}/baseline/check-rule/handle-affect-baseline
func DataSourceBaselineCheckRuleHAB() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBaselineCheckRuleHABRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			// Request Body Parameters.
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the operation performed during the baseline check.`,
			},
			"handle_status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the handle status of the baseline check rule.`,
			},
			"check_rule_list": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        checkRuleListSchema(),
				Description: `Specifies the list of baseline check rules to be handled.`,
			},
			"host_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the host ID.`,
			},
			// Request Query Parameters.
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			// Attribute Parameters.
			"total_rule_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total number of affected items.`,
			},
			"rule_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of check items affected by the operation.`,
			},
			"host_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of hosts affected by the operation.`,
			},
			"data_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataListAttributeSchema(),
				Description: `The detailed information about the operation impact scope.`,
			},
		},
	}
}

func dataListAttributeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"host_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The host ID.`,
			},
			"host_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The host name.`,
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The public IP address of the host.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The private IP address of the host.`,
			},
			"asset_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The asset value of the host.`,
			},
			"check_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the baseline check.`,
			},
			"standard": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The standard type of the host.`,
			},
			"tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The check type of check items in the baseline check.`,
			},
			"check_rule_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of check items in the baseline check.`,
			},
		},
	}
}

func checkRuleListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"check_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the baseline check.`,
			},
			"check_rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the baseline check rule.`,
			},
			"standard": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The standard type of the host.`,
			},
		},
	}
}

func buildBaselineCheckRuleHABQueryParams(epsId string) string {
	rst := "?limit=200"
	if epsId != "" {
		rst = fmt.Sprintf("%s&enterprise_project_id=%s", rst, epsId)
	}

	return rst
}

func buildCheckRuleHABCheckRuleListRequestBody(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"check_name":    utils.ValueIgnoreEmpty(rawMap["check_name"]),
			"check_rule_id": utils.ValueIgnoreEmpty(rawMap["check_rule_id"]),
			"standard":      utils.ValueIgnoreEmpty(rawMap["standard"]),
		})
	}

	return rst
}

func buildBaselineCheckRuleHABRequestBody(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"action":          d.Get("action"),
		"handle_status":   d.Get("handle_status"),
		"check_rule_list": buildCheckRuleHABCheckRuleListRequestBody(d.Get("check_rule_list").([]interface{})),
		"host_id":         utils.ValueIgnoreEmpty(d.Get("host_id")),
	}
}

func dataSourceBaselineCheckRuleHABRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		epsId        = cfg.GetEnterpriseProjectID(d)
		product      = "hss"
		httpUrl      = "v5/{project_id}/baseline/check-rule/handle-affect-baseline"
		offset       = 0
		allData      = make([]interface{}, 0)
		totalRuleNum = 0
		ruleNum      = 0
		hostNum      = 0
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildBaselineCheckRuleHABQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildBaselineCheckRuleHABRequestBody(d)),
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%d", requestPath, offset)
		resp, err := client.Request("POST", requestPathWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving HSS baseline check rule HAB: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalRuleNum = int(utils.PathSearch("total_rule_num", respBody, float64(0)).(float64))
		ruleNum = int(utils.PathSearch("rule_num", respBody, float64(0)).(float64))
		hostNum = int(utils.PathSearch("host_num", respBody, float64(0)).(float64))

		dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataList) == 0 {
			break
		}

		offset += len(dataList)
		allData = append(allData, dataList...)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_rule_num", totalRuleNum),
		d.Set("rule_num", ruleNum),
		d.Set("host_num", hostNum),
		d.Set("data_list", flattenBaselineCheckRuleHABDataList(allData)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBaselineCheckRuleHABDataList(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		result = append(result, map[string]interface{}{
			"host_id":         utils.PathSearch("host_id", v, nil),
			"host_name":       utils.PathSearch("host_name", v, nil),
			"public_ip":       utils.PathSearch("public_ip", v, nil),
			"private_ip":      utils.PathSearch("private_ip", v, nil),
			"asset_value":     utils.PathSearch("asset_value", v, nil),
			"check_type":      utils.PathSearch("check_type", v, nil),
			"standard":        utils.PathSearch("standard", v, nil),
			"tag":             utils.PathSearch("tag", v, nil),
			"check_rule_name": utils.PathSearch("check_rule_name", v, nil),
		})
	}

	return result
}
