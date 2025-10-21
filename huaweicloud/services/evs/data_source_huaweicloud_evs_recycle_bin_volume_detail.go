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

// @API EVS GET /v3/{project_id}/recycle-bin-volumes/{volume_id}
func DataSourceRecycleBinVolumeDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRecycleBinVolumeDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volume": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     recycleBinVolumesSchema(),
			},
		},
	}
}

func dataSourceRecycleBinVolumeDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "evs"
		httpUrl = "v3/{project_id}/recycle-bin-volumes/{volume_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{volume_id}", d.Get("volume_id").(string))
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving EVS recycle bin volume detail: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("volume", flattenRecycleBinVolumeDetail(utils.PathSearch("volume", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRecycleBinVolumeDetail(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":          utils.PathSearch("id", resp, nil),
			"name":        utils.PathSearch("name", resp, nil),
			"description": utils.PathSearch("description", resp, nil),
			"status":      utils.PathSearch("status", resp, nil),
			"attachments": flattenRecycleBinVolumeAttachments(
				utils.PathSearch("attachments", resp, make([]interface{}, 0)).([]interface{})),
			"multiattach":            utils.PathSearch("multiattach", resp, nil),
			"size":                   utils.PathSearch("size", resp, nil),
			"metadata":               utils.PathSearch("metadata", resp, nil),
			"bootable":               utils.PathSearch("bootable", resp, nil),
			"tags":                   utils.PathSearch("tags", resp, nil),
			"availability_zone":      utils.PathSearch("availability_zone", resp, nil),
			"created_at":             utils.PathSearch("created_at", resp, nil),
			"service_type":           utils.PathSearch("service_type", resp, nil),
			"updated_at":             utils.PathSearch("updated_at", resp, nil),
			"volume_type":            utils.PathSearch("volume_type", resp, nil),
			"enterprise_project_id":  utils.PathSearch("enterprise_project_id", resp, nil),
			"plan_delete_at":         utils.PathSearch("plan_delete_at", resp, nil),
			"pre_deleted_at":         utils.PathSearch("pre_deleted_at", resp, nil),
			"dedicated_storage_id":   utils.PathSearch("dedicated_storage_id", resp, nil),
			"dedicated_storage_name": utils.PathSearch("dedicated_storage_name", resp, nil),
			"volume_image_metadata":  utils.PathSearch("volume_image_metadata", resp, nil),
			"wwn":                    utils.PathSearch("wwn", resp, nil),
		},
	}
}

func flattenRecycleBinVolumeAttachments(resp []interface{}) []interface{} {
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
