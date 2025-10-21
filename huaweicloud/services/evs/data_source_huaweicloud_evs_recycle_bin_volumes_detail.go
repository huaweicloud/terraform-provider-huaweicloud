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

// @API EVS GET /v3/{project_id}/recycle-bin-volumes/detail
func DataSourceRecycleBinVolumesDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRecycleBinVolumesDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     recycleBinVolumesSchema(),
			},
		},
	}
}

func recycleBinVolumesSchema() *schema.Resource {
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
			// This field is of type String in the API documentation, but it is actually of type Bool.
			"multiattach": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"bootable": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan_delete_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pre_deleted_at": {
				Type:     schema.TypeString,
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
			// This field is of type String in the API documentation, but it is actually of type Map.
			"volume_image_metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"wwn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildRecycleBinVolumesDetailQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=1000"

	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		queryParams = fmt.Sprintf("%s&availability_zone=%v", queryParams, v)
	}
	if v, ok := d.GetOk("service_type"); ok {
		queryParams = fmt.Sprintf("%s&service_type=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceRecycleBinVolumesDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "evs"
		httpUrl = "v3/{project_id}/recycle-bin-volumes/detail"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildRecycleBinVolumesDetailQueryParams(d)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving EVS recycle bin volumes detail: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		volumesResp := utils.PathSearch("volumes", respBody, make([]interface{}, 0)).([]interface{})
		if len(volumesResp) == 0 {
			break
		}

		result = append(result, volumesResp...)
		offset += len(volumesResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("volumes", flattenRecycleBinVolumesDetail(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRecycleBinVolumesDetail(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":                     utils.PathSearch("id", v, nil),
			"name":                   utils.PathSearch("name", v, nil),
			"description":            utils.PathSearch("description", v, nil),
			"status":                 utils.PathSearch("status", v, nil),
			"attachments":            flattenRecycleBinVolumesAttachments(utils.PathSearch("attachments", v, make([]interface{}, 0)).([]interface{})),
			"multiattach":            utils.PathSearch("multiattach", v, nil),
			"size":                   utils.PathSearch("size", v, nil),
			"metadata":               utils.PathSearch("metadata", v, nil),
			"bootable":               utils.PathSearch("bootable", v, nil),
			"tags":                   utils.PathSearch("tags", v, nil),
			"availability_zone":      utils.PathSearch("availability_zone", v, nil),
			"created_at":             utils.PathSearch("created_at", v, nil),
			"service_type":           utils.PathSearch("service_type", v, nil),
			"updated_at":             utils.PathSearch("updated_at", v, nil),
			"volume_type":            utils.PathSearch("volume_type", v, nil),
			"enterprise_project_id":  utils.PathSearch("enterprise_project_id", v, nil),
			"plan_delete_at":         utils.PathSearch("plan_delete_at", v, nil),
			"pre_deleted_at":         utils.PathSearch("pre_deleted_at", v, nil),
			"dedicated_storage_id":   utils.PathSearch("dedicated_storage_id", v, nil),
			"dedicated_storage_name": utils.PathSearch("dedicated_storage_name", v, nil),
			"volume_image_metadata":  utils.PathSearch("volume_image_metadata", v, nil),
			"wwn":                    utils.PathSearch("wwn", v, nil),
		})
	}

	return rst
}

func flattenRecycleBinVolumesAttachments(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
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
