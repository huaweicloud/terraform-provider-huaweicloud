package cbr

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CBR POST /v3/{project_id}/vaults/resource_instances/action
func DataSourceVaultsByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVaultsByTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"without_any_tag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cloud_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"tags_any": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"not_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"not_tags_any": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"sys_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"object_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataResourceDetailSchema(),
			},
		},
	}
}

func dataResourceDetailSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
				},
			},
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_detail": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vault": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataResourceDetailVaultSchema(),
						},
					},
				},
			},
			"sys_tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
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

func dataResourceDetailVaultSchema() *schema.Resource {
	sc := schema.Resource{
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
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataResourceDetailVaultResourcesSchema(),
			},
			"provider_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_bind": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"bind_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
				},
			},
			"auto_expand": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"smn_notify": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"threshold": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"billing": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataResourceDetailVaultBillingSchema(),
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
				},
			},
			"backup_name_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"demand_billing": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cbc_delete_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"frozen": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sys_lock_source_service": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"supplier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cross_account": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cross_account_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataResourceDetailVaultResourcesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_protect": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backup_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backup_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protect_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extra_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retention_duration": {
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
						"exclude_volumes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
	return &sc
}

func dataResourceDetailVaultBillingSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"allocated": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"consistent_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"frozen_scene": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protect_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_multi_az": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_double_az": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"partner_bp_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func buildListVaultsByTagsBody(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{}
	if v, ok := d.GetOk("action"); ok {
		body["action"] = v.(string)
	}
	if d.Get("without_any_tag").(bool) {
		body["without_any_tag"] = true
	}
	if v, ok := d.GetOk("cloud_type"); ok {
		body["cloud_type"] = v.(string)
	}
	if v, ok := d.GetOk("object_type"); ok {
		body["object_type"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		body["tags"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("tags_any"); ok {
		body["tags_any"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("not_tags"); ok {
		body["not_tags"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("not_tags_any"); ok {
		body["not_tags_any"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("sys_tags"); ok {
		body["sys_tags"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("matches"); ok {
		body["matches"] = expandMatches(v.([]interface{}))
	}
	return body
}

func expandTags(rawTags []interface{}) []map[string]interface{} {
	tags := make([]map[string]interface{}, len(rawTags))
	for i, raw := range rawTags {
		tag := raw.(map[string]interface{})
		tags[i] = map[string]interface{}{
			"key":    tag["key"].(string),
			"values": tag["values"].([]interface{}),
		}
	}
	return tags
}

func expandMatches(rawTags []interface{}) []map[string]interface{} {
	tags := make([]map[string]interface{}, len(rawTags))
	for i, raw := range rawTags {
		tag := raw.(map[string]interface{})
		tags[i] = map[string]interface{}{
			"key":   tag["key"].(string),
			"value": tag["value"].(string),
		}
	}
	return tags
}

func flattenAllVaultsByTags(resp []interface{}) []map[string]interface{} {
	if len(resp) == 0 {
		return nil
	}
	results := make([]map[string]interface{}, 0, len(resp))
	for _, resource := range resp {
		resourceMap := map[string]interface{}{
			"resource_id":     utils.PathSearch("resource_id", resource, nil),
			"resource_name":   utils.PathSearch("resource_name", resource, nil),
			"tags":            flattenVaultTags(utils.PathSearch("tags", resource, make([]interface{}, 0)).([]interface{})),
			"resource_detail": flattenResourceDetail(utils.PathSearch("resource_detail.vault", resource, nil)),
		}
		results = append(results, resourceMap)
	}
	return results
}

func flattenResourceDetail(res interface{}) []interface{} {
	if res == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"vault": flattenVault(res),
		},
	}
}

func flattenVault(vaultMap interface{}) []interface{} {
	if vaultMap == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":                    utils.PathSearch("id", vaultMap, nil),
			"name":                  utils.PathSearch("name", vaultMap, nil),
			"provider_id":           utils.PathSearch("provider_id", vaultMap, nil),
			"created_at":            utils.PathSearch("created_at", vaultMap, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", vaultMap, nil),
			"auto_bind":             utils.PathSearch("auto_bind", vaultMap, nil),
			"auto_expand":           utils.PathSearch("auto_expand", vaultMap, nil),
			"smn_notify":            utils.PathSearch("smn_notify", vaultMap, nil),
			"threshold":             utils.PathSearch("threshold", vaultMap, nil),
			"bind_rules":            flattenVaultBindRules(utils.PathSearch("bind_rules", vaultMap, nil)),
			"resources":             flattenVaultRes(utils.PathSearch("resources", vaultMap, make([]interface{}, 0)).([]interface{})),
			"billing":               flattenVaultBilling(utils.PathSearch("billing", vaultMap, nil)),
			"tags":                  flattenVaultTags(utils.PathSearch("tags", vaultMap, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenVaultRes(resources []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(resources))

	for _, res := range resources {
		rst = append(rst, map[string]interface{}{
			"id":             utils.PathSearch("id", res, nil),
			"type":           utils.PathSearch("type", res, nil),
			"name":           utils.PathSearch("name", res, nil),
			"auto_protect":   utils.PathSearch("auto_protect", res, nil),
			"size":           utils.PathSearch("size", res, nil),
			"backup_size":    utils.PathSearch("backup_size", res, nil),
			"backup_count":   utils.PathSearch("backup_count", res, nil),
			"protect_status": utils.PathSearch("protect_status", res, nil),
			"extra_info":     flattenVaultResourcesExtraInfo(res),
		})
	}
	return rst
}

func flattenVaultResourcesExtraInfo(tags interface{}) []map[string]interface{} {
	if tags == nil {
		return nil
	}
	tagMap := utils.PathSearch("extra_info", tags, make(map[string]interface{}, 0)).(map[string]interface{})
	rst := make([]map[string]interface{}, 0, len(tagMap))
	for key, value := range tagMap {
		rst = append(rst, map[string]interface{}{
			"key":   key,
			"value": value,
		})
	}
	return rst
}

func flattenVaultBilling(billing interface{}) []interface{} {
	if billing == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"allocated":        utils.PathSearch("allocated", billing, nil),
			"cloud_type":       utils.PathSearch("cloud_type", billing, nil),
			"consistent_level": utils.PathSearch("consistent_level", billing, nil),
			"charging_mode":    utils.PathSearch("charging_mode", billing, nil),
			"order_id":         utils.PathSearch("order_id", billing, nil),
			"product_id":       utils.PathSearch("product_id", billing, nil),
			"protect_type":     utils.PathSearch("protect_type", billing, nil),
			"object_type":      utils.PathSearch("object_type", billing, nil),
			"spec_code":        utils.PathSearch("spec_code", billing, nil),
			"used":             utils.PathSearch("used", billing, 0),
			"status":           utils.PathSearch("status", billing, nil),
			"size":             utils.PathSearch("size", billing, 0),
			"is_multi_az":      utils.PathSearch("is_multi_az", billing, nil),
			"is_double_az":     utils.PathSearch("is_double_az", billing, nil),
			"storage_unit":     utils.PathSearch("storage_unit", billing, nil),
			"partner_bp_id":    utils.PathSearch("partner_bp_id", billing, nil),
		},
	}
}

func flattenVaultBindRules(rules interface{}) []interface{} {
	if rules == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"key":   utils.PathSearch("key", rules, nil),
			"value": utils.PathSearch("value", rules, nil),
		},
	}
}

func flattenVaultTags(tags []interface{}) []map[string]interface{} {
	if tags == nil {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(tags))
	for _, tag := range tags {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.PathSearch("value", tag, nil),
		})
	}
	return rst
}

func dataSourceVaultsByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/vault/resource_instances/action"
		product = "cbr"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	var allVaults []interface{}
	offset := 0
	totalCount := 0

	for {
		requestBody := buildListVaultsByTagsBody(d)
		if requestBody["action"] == "filter" {
			requestBody["limit"] = "1000"
			requestBody["offset"] = strconv.Itoa(offset)
		}

		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         requestBody,
		}
		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error querying CBR vaults by tags: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		vaults := utils.PathSearch("resources", respBody, []interface{}{}).([]interface{})
		totalCount = int(utils.PathSearch("total_count", respBody, float64(0)).(float64))
		if len(vaults) == 0 {
			break
		}
		allVaults = append(allVaults, vaults...)
		offset += len(vaults)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_count", totalCount),
		d.Set("resources", flattenAllVaultsByTags(allVaults)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving data source fields of the CBR vault by tags: %s", mErr)
	}
	return nil
}
