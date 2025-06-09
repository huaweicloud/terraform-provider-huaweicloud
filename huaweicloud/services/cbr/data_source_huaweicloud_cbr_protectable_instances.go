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

// @API CBR GET /v3/{project_id}/protectables/{protectable_type}/instances
func DataSourceProtectableInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProtectableInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"protectable_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the object type.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource name.`,
			},
			"server_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the server ID.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource status.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The protectable instances.`,
				Elem:        dataProtectableInstancesSchema(),
			},
		},
	}
}

func dataProtectableInstancesSchema() *schema.Resource {
	sc := &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Define fields `children` and `detail` as strings in JSON format.
			"children": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The children resources.`,
			},
			"detail": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource detail.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource name.`,
			},
			"protectable": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The backup information.`,
				Elem:        dataProtectableSchema(),
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The size of the resource, in GB.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource status.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the resource to be backed up.`,
			},
		},
	}

	return sc
}

func dataProtectableSchema() *schema.Resource {
	sc := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The error code for unsupported backup.`,
			},
			"reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The reason why backup is not supported.`,
			},
			"result": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether backup is supported.`,
			},
			"vault": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The associated vault.`,
				Elem:        instancesProtectableVaultSchema(),
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The reason why the resource cannot be backed up.`,
			},
		},
	}

	return sc
}

func instancesProtectableVaultSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"billing": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The operation information.`,
				Elem:        dataBillingSchema(),
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user-defined vault description.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The vault ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The vault name.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The project ID.`,
			},
			"provider_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the vault resource type.`,
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The resources.`,
				Elem:        dataResourcesSchema(),
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The vault tags.`,
				Elem:        dataVaultTagsSchema(),
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise project ID.`,
			},
			"auto_bind": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether automatic association is enabled.`,
			},
			"bind_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The association rules.`,
				Elem:        dataBindRulesSchema(),
			},
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user ID.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"auto_expand": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable auto capacity expansion for the vault.`,
			},
			"smn_notify": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `The SMN notification switch for the vault.`,
			},
			"threshold": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The vault capacity threshold.`,
			},
			"sys_lock_source_service": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The identity of the SMB service.`,
			},
			"locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the vault is locked.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version.`,
			},
		},
	}
}

func dataBindRulesSchema() *schema.Resource {
	sc := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The tags using to filter automatically associated resources.`,
				Elem:        dataBindRulesTagsSchema(),
			},
		},
	}

	return sc
}

func dataBindRulesTagsSchema() *schema.Resource {
	sc := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The tag key.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The tag value.`,
			},
		},
	}

	return sc
}

func dataVaultTagsSchema() *schema.Resource {
	sc := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The tag key.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The tag value.`,
			},
		},
	}

	return sc
}

func dataResourcesSchema() *schema.Resource {
	sc := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"extra_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The additional information of the resource.`,
				Elem:        dataResourceExtraInfoSchema(),
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the resource to be backed up.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the resource to be backed up.`,
			},
			"protect_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The protection status.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The allocated capacity for the associated resources, in GB.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the resource to be backed up.`,
			},
			"backup_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The backup size.`,
			},
			"backup_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of backups.`,
			},
		},
	}

	return sc
}

func dataResourceExtraInfoSchema() *schema.Resource {
	sc := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"exclude_volumes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The ID of the disk that will not be backed up.`,
			},
		},
	}

	return sc
}

func dataBillingSchema() *schema.Resource {
	sc := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allocated": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The allocated capacity, in GB.`,
			},
			"charging_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The charging mode.`,
			},
			"cloud_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cloud type.`,
			},
			"consistent_level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backup specifications.`,
			},
			"object_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The object type.`,
			},
			"order_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The order ID.`,
			},
			"product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The product ID.`,
			},
			"protect_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The protection type.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The capacity, in GB.`,
			},
			"spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The specification code.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The vault status.`,
			},
			"storage_unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the bucket for the vault.`,
			},
			"used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The used capacity, in MB.`,
			},
			"frozen_scene": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The scenario when an account is frozen.`,
			},
			"is_multi_az": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `The multi-AZ attribute of a vault.`,
			},
		},
	}

	return sc
}

func buildProtectableInstancesQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := "?limit=50"

	if v, ok := d.GetOk("resource_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}

	if v, ok := d.GetOk("server_id"); ok {
		res = fmt.Sprintf("%s&server_id=%v", res, v)
	}

	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}

	if epsID := cfg.GetEnterpriseProjectID(d); epsID != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsID)
	}

	return res
}

// There is a problem with the paging of the API.
// When querying server type resources, the query results are the same when the offset is `0` and `1`.
// However, this paging effect is different when querying disk type resources, so the logic cannot be unified.
func dataSourceProtectableInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		httpUrl         = "v3/{project_id}/protectables/{protectable_type}/instances"
		product         = "cbr"
		protectableType = d.Get("protectable_type").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{protectable_type}", protectableType)
	requestPath += buildProtectableInstancesQueryParams(d, cfg)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CBR protectable instances: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", flattenProtectableInstancesAttribute(respBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenProtectableInstancesAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	instances := utils.PathSearch("instances", respBody, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, 0, len(instances))
	for _, v := range instances {
		result = append(result, map[string]interface{}{
			"children":    utils.JsonToString(utils.PathSearch("children", v, nil)),
			"detail":      utils.JsonToString(utils.PathSearch("detail", v, nil)),
			"id":          utils.PathSearch("id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"protectable": flattenProtectableAttribute(utils.PathSearch("protectable", v, nil)),
			"size":        utils.PathSearch("size", v, nil),
			"status":      utils.PathSearch("status", v, nil),
			"type":        utils.PathSearch("type", v, nil),
		})
	}
	return result
}

func flattenProtectableAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstMap := map[string]interface{}{
		"code":    utils.PathSearch("code", respBody, nil),
		"reason":  utils.PathSearch("reason", respBody, nil),
		"result":  utils.PathSearch("result", respBody, nil),
		"vault":   flattenVaultAttribute(utils.PathSearch("vault", respBody, nil)),
		"message": utils.PathSearch("message", respBody, nil),
	}

	return []interface{}{rstMap}
}

func flattenVaultAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstMap := map[string]interface{}{
		"billing":                 flattenBillingAttribute(utils.PathSearch("billing", respBody, nil)),
		"description":             utils.PathSearch("description", respBody, nil),
		"id":                      utils.PathSearch("id", respBody, nil),
		"name":                    utils.PathSearch("name", respBody, nil),
		"project_id":              utils.PathSearch("project_id", respBody, nil),
		"provider_id":             utils.PathSearch("provider_id", respBody, nil),
		"resources":               flattenResourcesAttribute(utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})),
		"tags":                    flattenTagsAttribute(utils.PathSearch("tags", respBody, make([]interface{}, 0)).([]interface{})),
		"enterprise_project_id":   utils.PathSearch("enterprise_project_id", respBody, nil),
		"auto_bind":               utils.PathSearch("auto_bind", respBody, nil),
		"bind_rules":              flattenBindRulesAttribute(utils.PathSearch("bind_rules", respBody, nil)),
		"user_id":                 utils.PathSearch("user_id", respBody, nil),
		"created_at":              utils.PathSearch("created_at", respBody, nil),
		"auto_expand":             utils.PathSearch("auto_expand", respBody, nil),
		"smn_notify":              utils.PathSearch("smn_notify", respBody, nil),
		"threshold":               utils.PathSearch("threshold", respBody, nil),
		"sys_lock_source_service": utils.PathSearch("sys_lock_source_service", respBody, nil),
		"locked":                  utils.PathSearch("locked", respBody, nil),
		"updated_at":              utils.PathSearch("updated_at", respBody, nil),
		"version":                 utils.PathSearch("version", respBody, nil),
	}

	return []interface{}{rstMap}
}

func flattenBillingAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstMap := map[string]interface{}{
		"allocated":        utils.PathSearch("allocated", respBody, nil),
		"charging_mode":    utils.PathSearch("charging_mode", respBody, nil),
		"cloud_type":       utils.PathSearch("cloud_type", respBody, nil),
		"consistent_level": utils.PathSearch("consistent_level", respBody, nil),
		"object_type":      utils.PathSearch("object_type", respBody, nil),
		"order_id":         utils.PathSearch("order_id", respBody, nil),
		"product_id":       utils.PathSearch("product_id", respBody, nil),
		"protect_type":     utils.PathSearch("protect_type", respBody, nil),
		"size":             utils.PathSearch("size", respBody, nil),
		"spec_code":        utils.PathSearch("spec_code", respBody, nil),
		"status":           utils.PathSearch("status", respBody, nil),
		"storage_unit":     utils.PathSearch("storage_unit", respBody, nil),
		"used":             utils.PathSearch("used", respBody, nil),
		"frozen_scene":     utils.PathSearch("frozen_scene", respBody, nil),
		"is_multi_az":      utils.PathSearch("is_multi_az", respBody, nil),
	}

	return []interface{}{rstMap}
}

func flattenResourcesAttribute(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"extra_info":     flattenResourceExtraInfoAttribute(utils.PathSearch("extra_info", v, nil)),
			"id":             utils.PathSearch("id", v, nil),
			"name":           utils.PathSearch("name", v, nil),
			"protect_status": utils.PathSearch("protect_status", v, nil),
			"size":           utils.PathSearch("size", v, nil),
			"type":           utils.PathSearch("type", v, nil),
			"backup_size":    utils.PathSearch("backup_size", v, nil),
			"backup_count":   utils.PathSearch("backup_count", v, nil),
		})
	}

	return rst
}

func flattenResourceExtraInfoAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	excludeVolumes := utils.PathSearch("exclude_volumes", respBody, nil)
	if excludeVolumes == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"exclude_volumes": excludeVolumes,
	}}
}

func flattenTagsAttribute(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}

	return rst
}

func flattenBindRulesAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstMap := map[string]interface{}{
		"tags": flattenBindRulesTags(utils.PathSearch("tags", respBody, make([]interface{}, 0)).([]interface{})),
	}

	return []interface{}{rstMap}
}

func flattenBindRulesTags(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}

	return rst
}
