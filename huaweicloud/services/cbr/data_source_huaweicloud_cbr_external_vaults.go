package cbr

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

// @API CBR GET /v3/{project_id}/vaults/external
func DataSourceExternalVaults() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceExternalVaultsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"external_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protect_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vault_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vaults": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataExternalVaultSchema(),
			},
		},
	}
}

func dataExternalVaultSchema() *schema.Resource {
	return &schema.Resource{
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
			"provider_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
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
				Elem:     dataExternalVaultBindRulesSchema(),
			},
			"user_id": {
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
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataExternalVaultTagsSchema(),
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataExternalVaultResourcesSchema(),
			},
			"billing": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataExternalVaultBillingSchema(),
			},
		},
	}
}

func dataExternalVaultBindRulesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataExternalVaultTagsSchema(),
			},
		},
	}
}

func dataExternalVaultTagsSchema() *schema.Resource {
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

func dataExternalVaultResourcesSchema() *schema.Resource {
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
			"backup_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backup_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"extra_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataExternalVaultResourcesExtraInfoSchema(),
			},
		},
	}
}

func dataExternalVaultResourcesExtraInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"exclude_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataExternalVaultBillingSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allocated": {
				Type:     schema.TypeInt,
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
			"status": {
				Type:     schema.TypeString,
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
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_unit": {
				Type:     schema.TypeString,
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
}

func buildQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=100"
	queryParams = fmt.Sprintf("%s&region_id=%v", queryParams, d.Get("region_id"))
	queryParams = fmt.Sprintf("%s&external_project_id=%v", queryParams, d.Get("external_project_id"))
	if v, ok := d.GetOk("cloud_type"); ok {
		queryParams = fmt.Sprintf("%s&cloud_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("protect_type"); ok {
		queryParams = fmt.Sprintf("%s&protect_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("vault_id"); ok {
		queryParams = fmt.Sprintf("%s&vault_id=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceExternalVaultsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/vaults/external"
		product = "cbr"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildQueryParams(d)

	allVaults := make([]interface{}, 0)
	offset := 0

	for {
		reqUrl := fmt.Sprintf("%s&offset=%d", requestPath, offset)
		resp, err := client.Request("GET", reqUrl, &golangsdk.RequestOpts{
			KeepResponseBody: true,
		})
		if err != nil {
			return diag.Errorf("error querying CBR external vaults: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		vaults := utils.PathSearch("vaults", respBody, []interface{}{}).([]interface{})
		count := int(utils.PathSearch("count", respBody, float64(0)).(float64))
		allVaults = append(allVaults, vaults...)

		if len(vaults) == 0 || len(allVaults) >= count {
			break
		}

		offset += len(vaults)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("vaults", flattenExternalVaults(allVaults)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting data source fields of the CBR external vaults: %s", mErr)
	}
	return nil
}

func flattenExternalVaults(vaults []interface{}) []map[string]interface{} {
	if len(vaults) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(vaults))
	for i, vault := range vaults {
		result[i] = map[string]interface{}{
			"id":                      utils.PathSearch("id", vault, nil),
			"name":                    utils.PathSearch("name", vault, nil),
			"description":             utils.PathSearch("description", vault, nil),
			"provider_id":             utils.PathSearch("provider_id", vault, nil),
			"project_id":              utils.PathSearch("project_id", vault, nil),
			"enterprise_project_id":   utils.PathSearch("enterprise_project_id", vault, nil),
			"created_at":              utils.PathSearch("created_at", vault, nil),
			"auto_bind":               utils.PathSearch("auto_bind", vault, nil),
			"user_id":                 utils.PathSearch("user_id", vault, nil),
			"auto_expand":             utils.PathSearch("auto_expand", vault, nil),
			"smn_notify":              utils.PathSearch("smn_notify", vault, nil),
			"threshold":               utils.PathSearch("threshold", vault, nil),
			"sys_lock_source_service": utils.PathSearch("sys_lock_source_service", vault, nil),
			"locked":                  utils.PathSearch("locked", vault, nil),
			"bind_rules":              flattenExternalVaultBindRules(utils.PathSearch("bind_rules", vault, nil)),
			"tags":                    flattenExternalVaultTags(utils.PathSearch("tags", vault, nil)),
			"resources":               flattenExternalVaultResources(utils.PathSearch("resources", vault, nil)),
			"billing":                 flattenExternalVaultBilling(utils.PathSearch("billing", vault, nil)),
		}
	}

	return result
}

func flattenExternalVaultBindRules(rules interface{}) []map[string]interface{} {
	if rules == nil {
		return nil
	}

	result := map[string]interface{}{
		"tags": flattenExternalVaultTags(utils.PathSearch("tags", rules, nil)),
	}

	return []map[string]interface{}{result}
}

func flattenExternalVaultTags(tags interface{}) []map[string]interface{} {
	if tags == nil {
		return nil
	}

	tagsList, ok := tags.([]interface{})
	if !ok || len(tagsList) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(tagsList))
	for i, tag := range tagsList {
		result[i] = map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.PathSearch("value", tag, nil),
		}
	}

	return result
}

func flattenExternalVaultResources(resources interface{}) []map[string]interface{} {
	if resources == nil {
		return nil
	}

	resourcesList, ok := resources.([]interface{})
	if !ok || len(resourcesList) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(resourcesList))
	for i, res := range resourcesList {
		result[i] = map[string]interface{}{
			"id":             utils.PathSearch("id", res, nil),
			"type":           utils.PathSearch("type", res, nil),
			"name":           utils.PathSearch("name", res, nil),
			"protect_status": utils.PathSearch("protect_status", res, nil),
			"size":           utils.PathSearch("size", res, nil),
			"backup_size":    utils.PathSearch("backup_size", res, nil),
			"backup_count":   utils.PathSearch("backup_count", res, nil),
			"extra_info":     flattenExternalVaultResourcesExtraInfo(utils.PathSearch("extra_info", res, nil)),
		}
	}

	return result
}

func flattenExternalVaultResourcesExtraInfo(extraInfo interface{}) []map[string]interface{} {
	if extraInfo == nil {
		return nil
	}

	excludeVolumes := utils.PathSearch("exclude_volumes", extraInfo, []interface{}{}).([]interface{})

	result := map[string]interface{}{
		"exclude_volumes": utils.ExpandToStringList(excludeVolumes),
	}

	return []map[string]interface{}{result}
}

func flattenExternalVaultBilling(billing interface{}) []map[string]interface{} {
	if billing == nil {
		return nil
	}

	result := map[string]interface{}{
		"allocated":        utils.PathSearch("allocated", billing, nil),
		"used":             utils.PathSearch("used", billing, nil),
		"size":             utils.PathSearch("size", billing, nil),
		"status":           utils.PathSearch("status", billing, nil),
		"charging_mode":    utils.PathSearch("charging_mode", billing, nil),
		"cloud_type":       utils.PathSearch("cloud_type", billing, nil),
		"consistent_level": utils.PathSearch("consistent_level", billing, nil),
		"protect_type":     utils.PathSearch("protect_type", billing, nil),
		"object_type":      utils.PathSearch("object_type", billing, nil),
		"spec_code":        utils.PathSearch("spec_code", billing, nil),
		"order_id":         utils.PathSearch("order_id", billing, nil),
		"product_id":       utils.PathSearch("product_id", billing, nil),
		"storage_unit":     utils.PathSearch("storage_unit", billing, nil),
		"frozen_scene":     utils.PathSearch("frozen_scene", billing, nil),
		"is_multi_az":      utils.PathSearch("is_multi_az", billing, nil),
	}

	return []map[string]interface{}{result}
}
