package cbr

import (
	"context"
	"fmt"
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

// @API CBR POST /v3/{project_id}/vault/resource_instances/action
func DataSourceVaultsByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVaultsByTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"without_any_tag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
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
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
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
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
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
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
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
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"cloud_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"object_type": {
				Type:     schema.TypeString,
				Optional: true,
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
			"resource_id": {
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
			"resource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sys_tags": {
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
		},
	}
	return &sc
}

func dataResourceDetailVaultSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"billing": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataResourceDetailVaultBillingSchema(),
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataResourceDetailVaultResourcesSchema(),
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
					},
				},
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
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
			"sys_lock_source_service": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataResourceDetailVaultResourcesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"extra_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"exclude_volumes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protect_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
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
			"object_type": {
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
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
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
			"frozen_scene": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_multi_az": {
				Type:     schema.TypeBool,
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
			"sys_tags":        flattenVaultTags(utils.PathSearch("sys_tags", resource, make([]interface{}, 0)).([]interface{})),
		}
		results = append(results, resourceMap)
	}
	return results
}

func flattenResourceDetail(res interface{}) []map[string]interface{} {
	if res == nil {
		return nil
	}
	result := map[string]interface{}{
		"vault": flattenVault(res),
	}
	return []map[string]interface{}{result}
}

func flattenVault(vaultMap interface{}) []map[string]interface{} {
	if vaultMap == nil {
		return nil
	}
	result := map[string]interface{}{
		"id":                      utils.PathSearch("id", vaultMap, nil),
		"name":                    utils.PathSearch("name", vaultMap, nil),
		"description":             utils.PathSearch("description", vaultMap, nil),
		"provider_id":             utils.PathSearch("provider_id", vaultMap, nil),
		"enterprise_project_id":   utils.PathSearch("enterprise_project_id", vaultMap, nil),
		"auto_bind":               utils.PathSearch("auto_bind", vaultMap, nil),
		"user_id":                 utils.PathSearch("user_id", vaultMap, nil),
		"created_at":              utils.PathSearch("created_at", vaultMap, nil),
		"auto_expand":             utils.PathSearch("auto_expand", vaultMap, nil),
		"smn_notify":              utils.PathSearch("smn_notify", vaultMap, nil),
		"threshold":               utils.PathSearch("threshold", vaultMap, nil),
		"sys_lock_source_service": utils.PathSearch("sys_lock_source_service", vaultMap, nil),
		"locked":                  utils.PathSearch("locked", vaultMap, nil),
		"bind_rules":              flattenVaultBindRules(utils.PathSearch("bind_rules", vaultMap, nil)),
		"resources":               flattenVaultRes(utils.PathSearch("resources", vaultMap, make([]interface{}, 0)).([]interface{})),
		"billing":                 flattenVaultBilling(utils.PathSearch("billing", vaultMap, nil)),
		"tags":                    flattenVaultTags(utils.PathSearch("tags", vaultMap, make([]interface{}, 0)).([]interface{})),
	}
	return []map[string]interface{}{result}
}

func flattenVaultRes(resources []interface{}) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(resources))
	for _, res := range resources {
		result := map[string]interface{}{
			"id":             utils.PathSearch("id", res, nil),
			"type":           utils.PathSearch("type", res, nil),
			"name":           utils.PathSearch("name", res, nil),
			"size":           utils.PathSearch("size", res, nil),
			"backup_size":    utils.PathSearch("backup_size", res, nil),
			"backup_count":   utils.PathSearch("backup_count", res, nil),
			"protect_status": utils.PathSearch("protect_status", res, nil),
			"extra_info":     flattenVaultResourcesExtraInfo(res),
		}
		results = append(results, result)
	}
	return results
}

func flattenVaultResourcesExtraInfo(res interface{}) []map[string]interface{} {
	if res == nil {
		return nil
	}
	extraInfo := utils.PathSearch("extra_info", res, nil)
	if extraInfo == nil {
		return nil
	}
	excludeVolumes := utils.PathSearch("exclude_volumes", extraInfo, []interface{}{}).([]interface{})
	result := map[string]interface{}{
		"exclude_volumes": utils.ExpandToStringList(excludeVolumes),
	}
	return []map[string]interface{}{result}
}

func flattenVaultBilling(billing interface{}) []map[string]interface{} {
	if billing == nil {
		return nil
	}
	result := map[string]interface{}{
		"allocated":        utils.PathSearch("allocated", billing, nil),
		"charging_mode":    utils.PathSearch("charging_mode", billing, nil),
		"cloud_type":       utils.PathSearch("cloud_type", billing, nil),
		"consistent_level": utils.PathSearch("consistent_level", billing, nil),
		"object_type":      utils.PathSearch("object_type", billing, nil),
		"order_id":         utils.PathSearch("order_id", billing, nil),
		"product_id":       utils.PathSearch("product_id", billing, nil),
		"protect_type":     utils.PathSearch("protect_type", billing, nil),
		"size":             utils.PathSearch("size", billing, nil),
		"spec_code":        utils.PathSearch("spec_code", billing, nil),
		"status":           utils.PathSearch("status", billing, nil),
		"storage_unit":     utils.PathSearch("storage_unit", billing, nil),
		"used":             utils.PathSearch("used", billing, nil),
		"frozen_scene":     utils.PathSearch("frozen_scene", billing, nil),
		"is_multi_az":      utils.PathSearch("is_multi_az", billing, nil),
	}
	return []map[string]interface{}{result}
}

func flattenVaultBindRules(rules interface{}) []map[string]interface{} {
	if rules == nil {
		return nil
	}
	tags := flattenVaultTags(utils.PathSearch("tags", rules, make([]interface{}, 0)).([]interface{}))
	result := map[string]interface{}{
		"tags": tags,
	}
	return []map[string]interface{}{result}
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

func buildVaultsByTagsQueryParams(epsId string) string {
	queryParams := ""
	if epsId != "" {
		queryParams = fmt.Sprintf("%s?enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceVaultsByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsID   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v3/{project_id}/vault/resource_instances/action"
		product = "cbr"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildVaultsByTagsQueryParams(epsID)

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
