package dsc

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

// @API DSC GET /v1/{project_id}/security-policies/dbss-oem-info
func DataSourceDbssOemInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDbssOemInfoRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ins_info_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     oemInsInfoSchema(),
			},
		},
	}
}

func oemInsInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ins_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ins_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"related_datasource_policy_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     relatedDatasourceInfoSchema(),
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func relatedDatasourceInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datasource_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datasource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datasource_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ddm_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     ddmPolicySchema(),
			},
			"gde_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gdePolicySchema(),
			},
			"proxy_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sdm_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     sdmPolicySchema(),
			},
		},
	}
}

func ddmPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"columns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     columnSchema(),
			},
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func gdePolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"action": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"alg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"columns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     columnSchema(),
			},
			"table": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func sdmPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"columns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     columnSchema(),
			},
			"do_mask": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"do_move": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func columnSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"mask": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDbssOemInfoQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?type=%s", d.Get("type").(string))
}

func dataSourceDbssOemInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/security-policies/dbss-oem-info"
		product = "dsc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildDbssOemInfoQueryParams(d)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC DBSS OEM info: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("ins_info_list", flattenOemInsInfoList(utils.PathSearch(
			"ins_info_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOemInsInfoList(insInfoList []interface{}) []interface{} {
	if len(insInfoList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(insInfoList))
	for _, v := range insInfoList {
		rst = append(rst, utils.RemoveNil(map[string]interface{}{
			"ins_id":   utils.PathSearch("ins_id", v, nil),
			"ins_name": utils.PathSearch("ins_name", v, nil),
			"related_datasource_policy_list": flattenRelatedDatasourcePolicyList(utils.PathSearch(
				"related_datasource_policy_list", v, make([]interface{}, 0)).([]interface{})),
			"subnet_id": utils.PathSearch("subnet_id", v, nil),
			"vpc_id":    utils.PathSearch("vpc_id", v, nil),
		}))
	}

	return rst
}

func flattenRelatedDatasourcePolicyList(policyList []interface{}) []interface{} {
	if len(policyList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policyList))
	for _, v := range policyList {
		rst = append(rst, utils.RemoveNil(map[string]interface{}{
			"datasource_address": utils.PathSearch("datasource_address", v, nil),
			"datasource_id":      utils.PathSearch("datasource_id", v, nil),
			"datasource_port":    utils.PathSearch("datasource_port", v, nil),
			"ddm_policies": flattenDdmPolicies(utils.PathSearch(
				"ddm_policies", v, make([]interface{}, 0)).([]interface{})),
			"gde_policies": flattenGdePolicies(utils.PathSearch(
				"gde_policies", v, make([]interface{}, 0)).([]interface{})),
			"proxy_port": utils.PathSearch("proxy_port", v, nil),
			"sdm_policies": flattenSdmPolicies(utils.PathSearch(
				"sdm_policies", v, make([]interface{}, 0)).([]interface{})),
		}))
	}

	return rst
}

func flattenDdmPolicies(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, utils.RemoveNil(map[string]interface{}{
			"columns":   flattenColumns(utils.PathSearch("columns", v, make([]interface{}, 0)).([]interface{})),
			"namespace": utils.PathSearch("namespace", v, nil),
			"table":     utils.PathSearch("table", v, nil),
		}))
	}

	return rst
}

func flattenGdePolicies(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, utils.RemoveNil(map[string]interface{}{
			"action":  utils.PathSearch("action", v, nil),
			"alg":     utils.PathSearch("alg", v, nil),
			"columns": flattenColumns(utils.PathSearch("columns", v, make([]interface{}, 0)).([]interface{})),
			"table":   utils.PathSearch("table", v, nil),
		}))
	}

	return rst
}

func flattenSdmPolicies(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, utils.RemoveNil(map[string]interface{}{
			"columns":   flattenColumns(utils.PathSearch("columns", v, make([]interface{}, 0)).([]interface{})),
			"do_mask":   utils.PathSearch("do_mask", v, nil),
			"do_move":   utils.PathSearch("do_move", v, nil),
			"namespace": utils.PathSearch("namespace", v, nil),
			"table":     utils.PathSearch("table", v, nil),
		}))
	}

	return rst
}

func flattenColumns(columns []interface{}) []interface{} {
	if len(columns) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(columns))
	for _, v := range columns {
		rst = append(rst, utils.RemoveNil(map[string]interface{}{
			"mask": utils.PathSearch("mask", v, nil),
			"name": utils.PathSearch("name", v, nil),
		}))
	}

	return rst
}
