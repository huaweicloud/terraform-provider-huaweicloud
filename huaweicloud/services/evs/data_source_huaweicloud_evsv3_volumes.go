package evs

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

// @API EVS GET /v3/{project_id}/volumes/detail
func DataSourceV3Volumes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV3VolumesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metadata": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
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
							},
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
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attached_at": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"attachment_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"device": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"server_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"volume_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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
						"bootable": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"multiattach": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"volume_image_metadata": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"iops": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"frozened": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"total_val": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"volume_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"throughput": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"frozened": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"total_val": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"volume_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"snapshot_policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildV3VolumesQueryParams(d *schema.ResourceData) string {
	rst := ""

	if v, ok := d.GetOk("name"); ok {
		rst += fmt.Sprintf("&name=%v", v)
	}
	if v, ok := d.GetOk("sort_key"); ok {
		rst += fmt.Sprintf("&sort_key=%v", v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%v", v)
	}
	if v, ok := d.GetOk("status"); ok {
		rst += fmt.Sprintf("&status=%v", v)
	}
	if v, ok := d.GetOk("metadata"); ok {
		rst += fmt.Sprintf("&metadata=%v", v)
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		rst += fmt.Sprintf("&availability_zone=%v", v)
	}
	if rst != "" {
		rst = "?" + rst[1:]
	}

	return rst
}

func buildV3VolumesRequestPathWithOffset(queryParams string, offset int) string {
	if offset == 0 {
		// Ignore the offset of the first page.
		return queryParams
	}

	if queryParams == "" {
		// No query conditions were added.
		return fmt.Sprintf("?offset=%d", offset)
	}

	return fmt.Sprintf("%s&offset=%d", queryParams, offset)
}

func dataSourceV3VolumesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v3/{project_id}/volumes/detail"
		product     = "evs"
		queryParams = buildV3VolumesQueryParams(d)
		offset      = 0
		allVolumes  []interface{}
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := requestPath + buildV3VolumesRequestPathWithOffset(queryParams, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving EVS v3 volumes: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		volumesResp := utils.PathSearch("volumes", respBody, make([]interface{}, 0)).([]interface{})
		if len(volumesResp) == 0 {
			break
		}

		allVolumes = append(allVolumes, volumesResp...)
		offset += len(volumesResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("volumes", flattenV3Volumes(allVolumes)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenV3Volumes(allVolumes []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(allVolumes))
	for _, v := range allVolumes {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"links":             flattenV3VolumesLinks(v),
			"name":              utils.PathSearch("name", v, nil),
			"status":            utils.PathSearch("status", v, nil),
			"attachments":       flattenV3VolumesAttachments(v),
			"availability_zone": utils.PathSearch("availability_zone", v, nil),
			"snapshot_id":       utils.PathSearch("snapshot_id", v, nil),
			"description":       utils.PathSearch("description", v, nil),
			"bootable":          utils.PathSearch("bootable", v, nil),
			"created_at":        utils.PathSearch("created_at", v, nil),
			"volume_type":       utils.PathSearch("volume_type", v, nil),
			"metadata": utils.ExpandToStringMap(utils.PathSearch("metadata", v,
				make(map[string]interface{})).(map[string]interface{})),
			"size":        utils.PathSearch("size", v, nil),
			"updated_at":  utils.PathSearch("updated_at", v, nil),
			"multiattach": utils.PathSearch("multiattach", v, nil),
			"volume_image_metadata": utils.ExpandToStringMap(utils.PathSearch("volume_image_metadata", v,
				make(map[string]interface{})).(map[string]interface{})),
			"iops":               flattenV3VolumesIops(v),
			"throughput":         flattenV3VolumesThroughput(v),
			"snapshot_policy_id": utils.PathSearch("snapshot_policy_id", v, nil),
		})
	}

	return rst
}

func flattenV3VolumesLinks(respBody interface{}) []interface{} {
	linksResp := utils.PathSearch("links", respBody, make([]interface{}, 0)).([]interface{})
	if len(linksResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(linksResp))
	for _, v := range linksResp {
		rst = append(rst, map[string]interface{}{
			"href": utils.PathSearch("href", v, nil),
			"rel":  utils.PathSearch("rel", v, nil),
		})
	}

	return rst
}

func flattenV3VolumesAttachments(respBody interface{}) []interface{} {
	attachmentsResp := utils.PathSearch("attachments", respBody, make([]interface{}, 0)).([]interface{})
	rst := make([]interface{}, 0, len(attachmentsResp))
	for _, v := range attachmentsResp {
		rst = append(rst, map[string]interface{}{
			"attached_at":   utils.PathSearch("attached_at", v, nil),
			"attachment_id": utils.PathSearch("attachment_id", v, nil),
			"device":        utils.PathSearch("device", v, nil),
			"host_name":     utils.PathSearch("host_name", v, nil),
			"id":            utils.PathSearch("id", v, nil),
			"server_id":     utils.PathSearch("server_id", v, nil),
			"volume_id":     utils.PathSearch("volume_id", v, nil),
		})
	}

	return rst
}

func flattenV3VolumesIops(respBody interface{}) []interface{} {
	iopsResp := utils.PathSearch("iops", respBody, nil)
	if iopsResp == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"frozened":  utils.PathSearch("frozened", iopsResp, nil),
		"id":        utils.PathSearch("id", iopsResp, nil),
		"total_val": utils.PathSearch("total_val", iopsResp, nil),
		"volume_id": utils.PathSearch("volume_id", iopsResp, nil),
	}}
}

func flattenV3VolumesThroughput(respBody interface{}) []interface{} {
	throughputResp := utils.PathSearch("throughput", respBody, nil)
	if throughputResp == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"frozened":  utils.PathSearch("frozened", throughputResp, nil),
		"id":        utils.PathSearch("id", throughputResp, nil),
		"total_val": utils.PathSearch("total_val", throughputResp, nil),
		"volume_id": utils.PathSearch("volume_id", throughputResp, nil),
	}}
}
