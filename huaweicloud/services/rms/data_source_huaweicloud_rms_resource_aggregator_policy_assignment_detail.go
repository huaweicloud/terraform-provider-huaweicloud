package rms

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CONFIG POST /v1/resource-manager/domains/{domain_id}/aggregators/aggregate-data/policy-assignment/detail
func DataSourceAggregatorPolicyAssignmentDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAggregatorPolicyAssignmentDetailRead,

		Schema: map[string]*schema.Schema{
			"aggregator_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_assignment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_assignment_type": {
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
			"policy_filter": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     aggregatorPolicyAssignmentDetailPolicyFilter(),
			},
			"policy_filter_v2": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     aggregatorPolicyAssignmentDetailPolicyFilterV2(),
			},
			"period": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_definition_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     aggregatorPolicyAssignmentDetailCustomPolicy(),
			},
			"parameters": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     aggregatorPolicyAssignmentDetailTags(),
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func aggregatorPolicyAssignmentDetailPolicyFilter() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func aggregatorPolicyAssignmentDetailPolicyFilterV2() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"region_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tag_key_logic": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     aggregatorPolicyAssignmentDetailPolicyFilterV2Tags(),
			},
			"exclude_tag_key_logic": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"exclude_tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     aggregatorPolicyAssignmentDetailPolicyFilterV2Tags(),
			},
		},
	}
	return &sc
}

func aggregatorPolicyAssignmentDetailPolicyFilterV2Tags() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return &sc
}

func aggregatorPolicyAssignmentDetailCustomPolicy() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"function_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_value": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return &sc
}

func aggregatorPolicyAssignmentDetailTags() *schema.Resource {
	sc := schema.Resource{
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
	return &sc
}

func dataSourceAggregatorPolicyAssignmentDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/resource-manager/domains/{domain_id}/aggregators/aggregate-data/policy-assignment/detail"
		product = "rms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Config client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", cfg.DomainID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAggregatorPolicyAssignmentDetailQueryParams(d),
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving Config resource aggregator policy assignment detail, %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("policy_assignment_type", utils.PathSearch("policy_assignment_type", getRespBody, nil)),
		d.Set("policy_assignment_id", utils.PathSearch("id", getRespBody, nil)),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("policy_filter", flattenAggregatorPolicyAssignmentDetailPolicyFilter(getRespBody)),
		d.Set("policy_filter_v2", flattenAggregatorPolicyAssignmentDetailPolicyFilterV2(getRespBody)),
		d.Set("period", utils.PathSearch("period", getRespBody, nil)),
		d.Set("state", utils.PathSearch("state", getRespBody, nil)),
		d.Set("created", utils.PathSearch("created", getRespBody, nil)),
		d.Set("updated", utils.PathSearch("updated", getRespBody, nil)),
		d.Set("policy_definition_id", utils.PathSearch("policy_definition_id", getRespBody, nil)),
		d.Set("custom_policy", flattenAggregatorPolicyAssignmentDetailCustomPolicy(getRespBody)),
		d.Set("parameters", flattenAggregatorPolicyAssignmentDetailParameters(getRespBody)),
		d.Set("tags", flattenAggregatorPolicyAssignmentDetailTags(getRespBody)),
		d.Set("created_by", utils.PathSearch("created_by", getRespBody, nil)),
		d.Set("target_type", utils.PathSearch("target_type", getRespBody, nil)),
		d.Set("target_id", utils.PathSearch("target_id", getRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildAggregatorPolicyAssignmentDetailQueryParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"aggregator_id":        d.Get("aggregator_id"),
		"account_id":           d.Get("account_id"),
		"policy_assignment_id": d.Get("policy_assignment_id"),
	}
	return bodyParams
}

func flattenAggregatorPolicyAssignmentDetailPolicyFilter(resp interface{}) []interface{} {
	curJson := utils.PathSearch("policy_filter", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"region_id":         utils.PathSearch("region_id", curJson, nil),
			"resource_provider": utils.PathSearch("resource_provider", curJson, nil),
			"resource_type":     utils.PathSearch("resource_type", curJson, nil),
			"resource_id":       utils.PathSearch("resource_id", curJson, nil),
			"tag_key":           utils.PathSearch("tag_key", curJson, nil),
			"tag_value":         utils.PathSearch("tag_value", curJson, nil),
		},
	}
	return rst
}

func flattenAggregatorPolicyAssignmentDetailPolicyFilterV2(resp interface{}) []interface{} {
	curJson := utils.PathSearch("policy_filter_v2", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"region_ids":     utils.PathSearch("region_ids", curJson, nil),
			"resource_types": utils.PathSearch("resource_types", curJson, nil),
			"resource_ids":   utils.PathSearch("resource_ids", curJson, nil),
			"tag_key_logic":  utils.PathSearch("tag_key_logic", curJson, nil),
			"tags": flattenAggregatorPolicyAssignmentDetailPolicyFilterV2Tags(
				utils.PathSearch("tags", curJson, nil)),
			"exclude_tag_key_logic": utils.PathSearch("exclude_tag_key_logic", curJson, nil),
			"exclude_tags": flattenAggregatorPolicyAssignmentDetailPolicyFilterV2Tags(
				utils.PathSearch("exclude_tags", curJson, nil)),
		},
	}
	return rst
}

func flattenAggregatorPolicyAssignmentDetailPolicyFilterV2Tags(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"key":    utils.PathSearch("key", resp, nil),
			"values": utils.PathSearch("values", resp, nil),
		},
	}
	return rst
}

func flattenAggregatorPolicyAssignmentDetailCustomPolicy(resp interface{}) []interface{} {
	curJson := utils.PathSearch("custom_policy", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"function_urn": utils.PathSearch("function_urn", curJson, nil),
			"auth_type":    utils.PathSearch("auth_type", curJson, nil),
			"auth_value":   flattenAggregatorPolicyAssignmentDetailCustomPolicyAuthValue(curJson),
		},
	}
	return rst
}

func flattenAggregatorPolicyAssignmentDetailCustomPolicyAuthValue(resp interface{}) map[string]interface{} {
	curJson := utils.PathSearch("auth_value", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := make(map[string]interface{})
	for k, v := range curJson.(map[string]interface{}) {
		jsonBytes, _ := json.Marshal(v)
		rst[k] = string(jsonBytes)
	}
	return rst
}

func flattenAggregatorPolicyAssignmentDetailParameters(resp interface{}) map[string]interface{} {
	curJson := utils.PathSearch("parameters", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := make(map[string]interface{})
	for k, v := range curJson.(map[string]interface{}) {
		jsonBytes, _ := json.Marshal(v)
		rst[k] = string(jsonBytes)
	}
	return rst
}

func flattenAggregatorPolicyAssignmentDetailTags(resp interface{}) []interface{} {
	curJson := utils.PathSearch("tags", resp, nil)
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return rst
}
