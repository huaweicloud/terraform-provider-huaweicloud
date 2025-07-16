package sdrs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SDRS POST /v1/{project_id}/protected-instances/resource_instances/action
func DataSourceSdrsProtectedInstancesByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSdrsProtectedInstancesByTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the operation to be performed. `,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the tags to query resource list which contain all the specified tags.`,
				Elem:        dataTagParamsSchema(),
			},
			"tags_any": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the tags to query resource list which contain any of the specified tags.`,
				Elem:        dataTagParamsSchema(),
			},
			"not_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the tags to query resource list which do not contain all the specified tags.`,
				Elem:        dataTagParamsSchema(),
			},
			"not_tags_any": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the tags to query resource list which do not contain any of the specified tags.`,
				Elem:        dataTagParamsSchema(),
			},
			"matches": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the search field.`,
				Elem:        dataMatchesParamsSchema(),
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of protected instances that match the search criteria.`,
				Elem:        dataResourcesAttributeSchema(),
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total number of protected instances that match the search criteria.`,
			},
		},
	}
}

func dataTagParamsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the tag key.`,
			},
			"values": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Specifies the tag values.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataMatchesParamsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the search field key.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the search value.`,
			},
		},
	}
}

func dataResourcesAttributeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the protected instance.`,
			},
			"resource_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the protected instance.`,
			},
			"resource_detail": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The details of a protected instance.`,
				Elem:        dataResourceDetailAttributeSchema(),
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The tags of the protected instance.`,
				Elem:        dataTagsAttributeSchema(),
			},
		},
	}
}

func dataTagsAttributeSchema() *schema.Resource {
	return &schema.Resource{
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
}

func dataResourceDetailAttributeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the protected instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the protected instance.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the protected instance.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the protected instance.`,
			},
			"source_server": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the production site server.`,
			},
			"target_server": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the disaster recovery site server.`,
			},
			"server_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the protection group.`,
			},
			// Fields `created_at` and `updated_at` are string format in API document, but the actual value format is int.
			// So we configure the value to TypeInt here.
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The creation time of the protected instance.`,
			},
			"updated_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The last update time of the protected instance.`,
			},
			"metadata": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The metadata of the protected instance.`,
				Elem:        dataMetadataAttributeSchema(),
			},
			"attachment": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The attached replication pairs.`,
				Elem:        dataAttachmentAttributeSchema(),
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The tags of the protected instance.`,
				Elem:        dataResourceDetailTagsAttributeSchema(),
			},
			"progress": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The synchronization progress of the protected instance.`,
			},
			"priority_station": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current production site AZ of the protection group.`,
			},
		},
	}
}

func dataResourceDetailTagsAttributeSchema() *schema.Resource {
	return &schema.Resource{
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
}

func dataAttachmentAttributeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"replication": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the replication pair.`,
			},
			"device": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The device name.`,
			},
		},
	}
}

func dataMetadataAttributeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"__system__frozen": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The frozen status of the resource.`,
			},
		},
	}
}

func buildRequestInstancesTagsBodyParams(tags []interface{}) []map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(tags))
	for _, tag := range tags {
		tagMap, ok := tag.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"key":    tagMap["key"],
			"values": tagMap["values"],
		})
	}
	return rst
}

func buildRequestInstancesMatchesBodyParams(matches []interface{}) []map[string]interface{} {
	if len(matches) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(matches))
	for _, match := range matches {
		matchMap, ok := match.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"key":   matchMap["key"],
			"value": matchMap["value"],
		})
	}
	return rst
}

func buildRequestInstancesBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"action":       d.Get("action").(string),
		"tags":         buildRequestInstancesTagsBodyParams(d.Get("tags").([]interface{})),
		"tags_any":     buildRequestInstancesTagsBodyParams(d.Get("tags_any").([]interface{})),
		"not_tags":     buildRequestInstancesTagsBodyParams(d.Get("not_tags").([]interface{})),
		"not_tags_any": buildRequestInstancesTagsBodyParams(d.Get("not_tags_any").([]interface{})),
		"matches":      buildRequestInstancesMatchesBodyParams(d.Get("matches").([]interface{})),
	}
}

func doRequestInstancesByTags(client *golangsdk.ServiceClient, path string, requestOpt golangsdk.RequestOpts) (interface{}, error) {
	resp, err := client.Request("POST", path, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func dataSourceSdrsProtectedInstancesByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "sdrs"
		path       = "v1/{project_id}/protected-instances/resource_instances/action"
		action     = d.Get("action").(string)
		offset     = 0
		totalCount = 0
		allItems   []interface{}
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}

	requestPath := strings.ReplaceAll(client.Endpoint+path, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildRequestInstancesBodyParams(d)),
	}

	if action == "count" {
		respBody, err := doRequestInstancesByTags(client, requestPath, requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SDRS protected instances count: %s", err)
		}

		totalCount = int(utils.PathSearch("total_count", respBody, float64(0)).(float64))
	}

	if action == "filter" {
		for {
			requestOpt.JSONBody.(map[string]interface{})["offset"] = offset
			respBody, err := doRequestInstancesByTags(client, requestPath, requestOpt)
			if err != nil {
				return diag.Errorf("error retrieving SDRS protected instances list: %s", err)
			}

			resources := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
			if len(resources) == 0 {
				break
			}
			allItems = append(allItems, resources...)
			offset += len(resources)
		}
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("resources", flattenDataResources(allItems)),
		d.Set("total_count", totalCount),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataResources(rawArray []interface{}) []interface{} {
	result := make([]interface{}, 0, len(rawArray))
	for _, rawMap := range rawArray {
		result = append(result, map[string]interface{}{
			"resource_id":     utils.PathSearch("resource_id", rawMap, nil),
			"resource_name":   utils.PathSearch("resource_name", rawMap, nil),
			"resource_detail": flattenDataResourceDetail(utils.PathSearch("resource_detail", rawMap, nil)),
			"tags":            flattenDataTags(utils.PathSearch("tags", rawMap, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenDataResourceDetail(rawMap interface{}) []interface{} {
	if rawMap == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":               utils.PathSearch("id", rawMap, nil),
			"name":             utils.PathSearch("name", rawMap, nil),
			"description":      utils.PathSearch("description", rawMap, nil),
			"status":           utils.PathSearch("status", rawMap, nil),
			"source_server":    utils.PathSearch("source_server", rawMap, nil),
			"target_server":    utils.PathSearch("target_server", rawMap, nil),
			"server_group_id":  utils.PathSearch("server_group_id", rawMap, nil),
			"created_at":       utils.PathSearch("created_at", rawMap, nil),
			"updated_at":       utils.PathSearch("updated_at", rawMap, nil),
			"metadata":         flattenDataMetadata(utils.PathSearch("metadata", rawMap, nil)),
			"attachment":       flattenDataAttachments(utils.PathSearch("attachment", rawMap, make([]interface{}, 0)).([]interface{})),
			"tags":             flattenDataTags(utils.PathSearch("tags", rawMap, make([]interface{}, 0)).([]interface{})),
			"progress":         utils.PathSearch("progress", rawMap, nil),
			"priority_station": utils.PathSearch("priority_station", rawMap, nil),
		},
	}
}

func flattenDataAttachments(rawArray []interface{}) []interface{} {
	result := make([]interface{}, 0, len(rawArray))
	for _, rawMap := range rawArray {
		result = append(result, map[string]interface{}{
			"replication": utils.PathSearch("replication", rawMap, nil),
			"device":      utils.PathSearch("device", rawMap, nil),
		})
	}
	return result
}

func flattenDataMetadata(rawMap interface{}) []interface{} {
	if rawMap == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"__system__frozen": utils.PathSearch("__system__frozen", rawMap, nil),
		},
	}
}

func flattenDataTags(rawArray []interface{}) []interface{} {
	result := make([]interface{}, 0, len(rawArray))
	for _, rawMap := range rawArray {
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", rawMap, nil),
			"value": utils.PathSearch("value", rawMap, nil),
		})
	}
	return result
}
