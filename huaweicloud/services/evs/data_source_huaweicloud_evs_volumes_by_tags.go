package evs

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

// @API EVS POST /v2/{project_id}/cloudvolumes/resource_instances/action
func DataSourceEvsVolumesByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvsVolumesByTagsRead,

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
			"tags": {
				Type:     schema.TypeList,
				Required: true,
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
				Elem:     dataSourceEvsVolumesByTagsResourceSchema(),
			},
		},
	}
}

func dataSourceEvsVolumesByTagsResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
				Elem:     dataSourceEvsVolumesByTagsResourceDetailSchema(),
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
		},
	}
}

func dataSourceEvsVolumesByTagsResourceDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceEvsVolumesByTagsLinksSchema(),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceEvsVolumesByTagsAttachmentSchema(),
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_vol_tenant_attr_tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_image_metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"volume_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bootable": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceEvsVolumesByTagsMetadataSchema(),
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"multiattach": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dedicated_storage_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dedicated_storage_name": {
				Type:     schema.TypeString,
				Computed: true,
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
			"wwn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceEvsVolumesByTagsLinksSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"href": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rel": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceEvsVolumesByTagsMetadataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"__system__cmkid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"__system__encrypted": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"full_clone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hw_passthrough": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceEvsVolumesByTagsAttachmentSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attached_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildListEvsVolumesByTagsBody(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"action": d.Get("action").(string),
	}
	if v, ok := d.GetOk("tags"); ok {
		body["tags"] = expandTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("matches"); ok {
		if list, ok := v.([]interface{}); ok && len(list) > 0 {
			body["matches"] = expandMatches(list)
		}
	}
	return body
}

func expandTags(rawTags []interface{}) []map[string]interface{} {
	tags := make([]map[string]interface{}, 0, len(rawTags))
	for _, t := range rawTags {
		tag := t.(map[string]interface{})
		tags = append(tags, map[string]interface{}{
			"key":    tag["key"],
			"values": tag["values"],
		})
	}
	return tags
}

func expandMatches(rawMatches []interface{}) []map[string]interface{} {
	if len(rawMatches) == 0 {
		return nil
	}
	matches := make([]map[string]interface{}, 0, len(rawMatches))
	for _, m := range rawMatches {
		match, ok := m.(map[string]interface{})
		if !ok {
			continue
		}
		matches = append(matches, map[string]interface{}{
			"key":   match["key"],
			"value": match["value"],
		})
	}
	return matches
}

func flattenAllEvsVolumesByTags(resp []interface{}) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(resp))
	for _, resource := range resp {
		results = append(results, map[string]interface{}{
			"resource_id":     utils.PathSearch("resource_id", resource, nil),
			"resource_name":   utils.PathSearch("resource_name", resource, nil),
			"tags":            flattenEvsTags(utils.PathSearch("tags", resource, make([]interface{}, 0)).([]interface{})),
			"resource_detail": flattenEvsResourceDetail(utils.PathSearch("resource_detail", resource, nil)),
		})
	}
	return results
}

func flattenEvsTags(tags []interface{}) []map[string]interface{} {
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

func flattenEvsResourceDetail(volumeMap interface{}) []map[string]interface{} {
	if volumeMap == nil {
		return nil
	}
	resourceDetailLinks := utils.PathSearch("links", volumeMap, make([]interface{}, 0)).([]interface{})
	resourceDetailAttachments := utils.PathSearch("attachments", volumeMap, make([]interface{}, 0)).([]interface{})
	resourceDetailMetadata := utils.PathSearch("metadata", volumeMap, make(map[string]interface{})).(map[string]interface{})
	resourceDetailTags := utils.PathSearch("tags", volumeMap, make([]interface{}, 0)).([]interface{})

	return []map[string]interface{}{
		{
			"id":                           utils.PathSearch("id", volumeMap, nil),
			"links":                        flattenEvsLinks(resourceDetailLinks),
			"name":                         utils.PathSearch("name", volumeMap, nil),
			"status":                       utils.PathSearch("status", volumeMap, nil),
			"attachments":                  flattenEvsAttachments(resourceDetailAttachments),
			"availability_zone":            utils.PathSearch("availability_zone", volumeMap, nil),
			"snapshot_id":                  utils.PathSearch("snapshot_id", volumeMap, nil),
			"description":                  utils.PathSearch("description", volumeMap, nil),
			"created_at":                   utils.PathSearch("created_at", volumeMap, nil),
			"os_vol_tenant_attr_tenant_id": utils.PathSearch("os-vol-tenant-attr:tenant_id", volumeMap, nil),
			"volume_image_metadata":        utils.PathSearch("volume_image_metadata", volumeMap, map[string]interface{}{}),
			"volume_type":                  utils.PathSearch("volume_type", volumeMap, nil),
			"size":                         utils.PathSearch("size", volumeMap, nil),
			"bootable":                     utils.PathSearch("bootable", volumeMap, nil),
			"metadata":                     flattenEvsMetadata(resourceDetailMetadata),
			"updated_at":                   utils.PathSearch("updated_at", volumeMap, nil),
			"service_type":                 utils.PathSearch("service_type", volumeMap, nil),
			"multiattach":                  utils.PathSearch("multiattach", volumeMap, nil),
			"dedicated_storage_id":         utils.PathSearch("dedicated_storage_id", volumeMap, nil),
			"dedicated_storage_name":       utils.PathSearch("dedicated_storage_name", volumeMap, nil),
			"tags":                         flattenEvsTags(resourceDetailTags),
			"wwn":                          utils.PathSearch("wwn", volumeMap, nil),
			"enterprise_project_id":        utils.PathSearch("enterprise_project_id", volumeMap, nil),
		},
	}
}

func flattenEvsLinks(links []interface{}) []map[string]interface{} {
	if links == nil {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(links))
	for _, link := range links {
		rst = append(rst, map[string]interface{}{
			"href": utils.PathSearch("href", link, nil),
			"rel":  utils.PathSearch("rel", link, nil),
		})
	}
	return rst
}

func flattenEvsMetadata(metadata map[string]interface{}) []map[string]interface{} {
	if metadata == nil {
		return nil
	}

	return []map[string]interface{}{{
		"__system__cmkid":     utils.PathSearch("__system__cmkid", metadata, nil),
		"__system__encrypted": utils.PathSearch("__system__encrypted", metadata, nil),
		"full_clone":          utils.PathSearch("full_clone", metadata, nil),
		"hw_passthrough":      utils.PathSearch("hw:passthrough", metadata, nil),
	}}
}

func flattenEvsAttachments(attachments []interface{}) []map[string]interface{} {
	if attachments == nil {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(attachments))
	for _, att := range attachments {
		rst = append(rst, map[string]interface{}{
			"server_id":     utils.PathSearch("server_id", att, nil),
			"attachment_id": utils.PathSearch("attachment_id", att, nil),
			"attached_at":   utils.PathSearch("attached_at", att, nil),
			"volume_id":     utils.PathSearch("volume_id", att, nil),
			"device":        utils.PathSearch("device", att, nil),
			"id":            utils.PathSearch("id", att, nil),
			"host_name":     utils.PathSearch("host_name", att, nil),
		})
	}
	return rst
}

func dataSourceEvsVolumesByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/cloudvolumes/resource_instances/action"
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	var allVolumes []interface{}
	offset := 0
	totalCount := 0

	for {
		requestBody := buildListEvsVolumesByTagsBody(d)
		if requestBody["action"] == "filter" {
			requestBody["limit"] = 1000
			requestBody["offset"] = offset
		}
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         requestBody,
		}
		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error querying EVS volumes by tags: %s", err)
		}
		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}
		volumes := utils.PathSearch("resources", respBody, []interface{}{}).([]interface{})
		totalCount = int(utils.PathSearch("total_count", respBody, float64(0)).(float64))
		if len(volumes) == 0 {
			break
		}
		allVolumes = append(allVolumes, volumes...)
		offset += len(volumes)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_count", totalCount),
		d.Set("resources", flattenAllEvsVolumesByTags(allVolumes)),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving data source fields of the EVS volumes by tags: %s", mErr)
	}
	return nil
}
