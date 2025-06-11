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
			"object_type": {
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

func buildQueryParams(d *schema.ResourceData) []string {
	queryParams := []string{}

	if v, ok := d.GetOk("external_project_id"); ok {
		queryParams = append(queryParams, "external_project_id="+v.(string))
	}
	if v, ok := d.GetOk("region_id"); ok {
		queryParams = append(queryParams, "region_id="+v.(string))
	}
	if v, ok := d.GetOk("cloud_type"); ok {
		queryParams = append(queryParams, "cloud_type="+v.(string))
	}
	if v, ok := d.GetOk("object_type"); ok {
		queryParams = append(queryParams, "object_type="+v.(string))
	}
	if v, ok := d.GetOk("protect_type"); ok {
		queryParams = append(queryParams, "protect_type="+v.(string))
	}
	if v, ok := d.GetOk("vault_id"); ok {
		queryParams = append(queryParams, "vault_id="+v.(string))
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

	queryParams := buildQueryParams(d)
	if len(queryParams) > 0 {
		requestPath = requestPath + "?" + strings.Join(queryParams, "&")
	}

	allVaults := make([]interface{}, 0)
	limit := 1000
	offset := 0

	for {
		paginationPath := requestPath
		if strings.Contains(paginationPath, "?") {
			paginationPath += fmt.Sprintf("&limit=%d&offset=%d", limit, offset)
		} else {
			paginationPath += fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
		}

		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		resp, err := client.Request("GET", paginationPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error querying CBR external vaults: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		vaults := utils.PathSearch("vaults", respBody, []interface{}{}).([]interface{})
		allVaults = append(allVaults, vaults...)

		if len(vaults) == 0 {
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
			"id":                      utils.PathSearch("id", vault, ""),
			"name":                    utils.PathSearch("name", vault, ""),
			"description":             utils.PathSearch("description", vault, ""),
			"provider_id":             utils.PathSearch("provider_id", vault, ""),
			"project_id":              utils.PathSearch("project_id", vault, ""),
			"enterprise_project_id":   utils.PathSearch("enterprise_project_id", vault, ""),
			"created_at":              utils.PathSearch("created_at", vault, ""),
			"auto_bind":               utils.PathSearch("auto_bind", vault, false),
			"user_id":                 utils.PathSearch("user_id", vault, ""),
			"auto_expand":             utils.PathSearch("auto_expand", vault, false),
			"smn_notify":              utils.PathSearch("smn_notify", vault, false),
			"threshold":               utils.PathSearch("threshold", vault, 0),
			"sys_lock_source_service": utils.PathSearch("sys_lock_source_service", vault, ""),
			"locked":                  utils.PathSearch("locked", vault, false),
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

	tagsList := tags.([]interface{})
	if len(tagsList) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(tagsList))
	for i, tag := range tagsList {
		result[i] = map[string]interface{}{
			"key":   utils.PathSearch("key", tag, ""),
			"value": utils.PathSearch("value", tag, ""),
		}
	}

	return result
}

func flattenExternalVaultResources(resources interface{}) []map[string]interface{} {
	if resources == nil {
		return nil
	}

	resourcesList := resources.([]interface{})
	if len(resourcesList) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(resourcesList))
	for i, res := range resourcesList {
		result[i] = map[string]interface{}{
			"id":             utils.PathSearch("id", res, ""),
			"type":           utils.PathSearch("type", res, ""),
			"name":           utils.PathSearch("name", res, ""),
			"protect_status": utils.PathSearch("protect_status", res, ""),
			"size":           utils.PathSearch("size", res, 0),
			"backup_size":    utils.PathSearch("backup_size", res, 0),
			"backup_count":   utils.PathSearch("backup_count", res, 0),
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
	excludeVolumesList := make([]string, len(excludeVolumes))
	for i, v := range excludeVolumes {
		excludeVolumesList[i] = v.(string)
	}

	result := map[string]interface{}{
		"exclude_volumes": excludeVolumesList,
	}

	return []map[string]interface{}{result}
}

func flattenExternalVaultBilling(billing interface{}) []map[string]interface{} {
	if billing == nil {
		return nil
	}

	result := map[string]interface{}{
		"allocated":        utils.PathSearch("allocated", billing, 0),
		"used":             utils.PathSearch("used", billing, 0),
		"size":             utils.PathSearch("size", billing, 0),
		"status":           utils.PathSearch("status", billing, ""),
		"charging_mode":    utils.PathSearch("charging_mode", billing, ""),
		"cloud_type":       utils.PathSearch("cloud_type", billing, ""),
		"consistent_level": utils.PathSearch("consistent_level", billing, ""),
		"protect_type":     utils.PathSearch("protect_type", billing, ""),
		"object_type":      utils.PathSearch("object_type", billing, ""),
		"spec_code":        utils.PathSearch("spec_code", billing, ""),
		"order_id":         utils.PathSearch("order_id", billing, ""),
		"product_id":       utils.PathSearch("product_id", billing, ""),
		"storage_unit":     utils.PathSearch("storage_unit", billing, ""),
		"frozen_scene":     utils.PathSearch("frozen_scene", billing, ""),
		"is_multi_az":      utils.PathSearch("is_multi_az", billing, false),
	}

	return []map[string]interface{}{result}
}
