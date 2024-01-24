package workspace

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v2/{project_id}/desktops/detail
func DataSourceDesktops() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDesktopsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"desktop_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fixed_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"desktop_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"user_attached": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"in_maintenance_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"desktops": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     desktopSchema(),
			},
		},
	}
}

func desktopSchema() *schema.Resource {
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"in_maintenance_mode": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"internet_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_addresses": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"root_volume": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     volumeSchema(),
			},
			"data_volume": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     volumeSchema(),
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attach_user_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_group": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"site_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"site_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ou_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ou_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attach_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"join_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_support_internet": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func volumeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
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
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildListDesktopsQueryParams(cfg *config.Config, d *schema.ResourceData) string {
	queryParam := ""
	if v, ok := d.GetOk("desktop_id"); ok {
		queryParam += fmt.Sprintf("&desktop_id=%v", v)
	}

	if v, ok := d.GetOk("name"); ok {
		queryParam += fmt.Sprintf("&computer_name=%v", v)
	}

	if v, ok := d.GetOk("user_name"); ok {
		queryParam += fmt.Sprintf("&user_name=%v", v)
	}

	if v, ok := d.GetOk("fixed_ip"); ok {
		queryParam += fmt.Sprintf("&desktop_ip=%v", v)
	}

	if v, ok := d.GetOk("desktop_type"); ok {
		queryParam += fmt.Sprintf("&desktop_type=%v", v)
	}

	if v, ok := d.GetOk("user_attached"); ok {
		queryParam += fmt.Sprintf("&user_attached=%v", v)
	}

	if v, ok := d.GetOk("image_id"); ok {
		queryParam += fmt.Sprintf("&image_id=%v", v)
	}

	if v, ok := d.GetOk("in_maintenance_mode"); ok {
		queryParam += fmt.Sprintf("&in_maintenance_mode=%v", v)
	}

	if v, ok := d.GetOk("status"); ok {
		queryParam += fmt.Sprintf("&status=%v", v)
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		queryParam += fmt.Sprintf("&subnet_id=%v", v)
	}

	epsID := cfg.GetEnterpriseProjectID(d)
	if epsID != "" {
		queryParam = fmt.Sprintf("%s&enterprise_project_id=%v", queryParam, epsID)
	}

	if queryParam != "" {
		queryParam = "?" + queryParam[1:]
	}

	return queryParam
}

func dataSourceDesktopsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.WorkspaceV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	listDesktopsPath := client.ResourceBaseURL() + "desktops/detail"
	listDesktopsQueryParams := buildListDesktopsQueryParams(cfg, d)
	listDesktopsPath += listDesktopsQueryParams

	resp, err := pagination.ListAllItems(
		client,
		"offset",
		listDesktopsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Workspace desktops")
	}

	listRespJson, err := json.Marshal(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	curJson := utils.PathSearch("desktops", listRespBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("desktops", flattenListDesktops(filterListDesktopByTags(curArray, d))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterListDesktopByTags(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	tagFilter := d.Get("tags").(map[string]interface{})
	if len(tagFilter) == 0 {
		return all
	}

	for _, v := range all {
		tags := utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil))
		tagmap := utils.ExpandToStringMap(tags)
		if !utils.HasMapContains(tagmap, tagFilter) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func flattenListDesktops(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, len(resp))
	for i, v := range resp {
		rst[i] = map[string]interface{}{
			"id":                    utils.PathSearch("desktop_id", v, nil),
			"name":                  utils.PathSearch("computer_name", v, nil),
			"type":                  utils.PathSearch("desktop_type", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"image_id":              utils.PathSearch("metadata.\"metering.image_id\"", v, nil),
			"status":                utils.PathSearch("status", v, nil),
			"in_maintenance_mode":   utils.PathSearch("in_maintenance_mode", v, false).(bool),
			"subnet_id":             utils.PathSearch("subnet_id", v, nil),
			"tags":                  utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil)),
			"internet_mode":         utils.PathSearch("internet_mode", v, nil),
			"ip_addresses":          utils.PathSearch("ip_addresses", v, nil),
			"root_volume":           flattenRootVolume(utils.PathSearch("root_volume", v, nil)),
			"data_volume":           flattenDataVolumes(utils.PathSearch("data_volumes", v, make([]interface{}, 0))),
			"availability_zone":     utils.PathSearch("availability_zone", v, nil),
			"attach_user_infos":     flattenAttachUserInfos(utils.PathSearch("attach_user_infos", v, make([]interface{}, 0))),
			"created_at":            utils.PathSearch("created", v, nil),
			"site_type":             utils.PathSearch("site_type", v, nil),
			"site_name":             utils.PathSearch("site_name", v, nil),
			"product_id":            utils.PathSearch("product_id", v, nil),
			"flavor_id":             utils.PathSearch("product.flavor_id", v, nil),
			"ou_name":               utils.PathSearch("ou_name", v, nil),
			"ou_version":            utils.PathSearch("ou_version", v, nil),
			"attach_state":          utils.PathSearch("attach_state", v, nil),
			"join_domain":           utils.PathSearch("join_domain", v, nil),
			"is_support_internet":   utils.PathSearch("is_support_internet", v, nil),
		}
	}
	return rst
}

func getVolume(volume interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":       utils.PathSearch("type", volume, nil),
		"size":       utils.PathSearch("size", volume, nil),
		"device":     utils.PathSearch("device", volume, nil),
		"id":         utils.PathSearch("id", volume, nil),
		"volume_id":  utils.PathSearch("volume_id", volume, nil),
		"created_at": utils.PathSearch("create_time", volume, nil),
		"name":       utils.PathSearch("display_name", volume, nil),
	}
}

func flattenRootVolume(volume interface{}) []map[string]interface{} {
	if volume == nil {
		return nil
	}

	rootVolume := []map[string]interface{}{getVolume(volume)}
	return rootVolume
}

func flattenDataVolumes(curRaw interface{}) []map[string]interface{} {
	curArray := curRaw.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(curArray))
	for i, volume := range curArray {
		result[i] = getVolume(volume)
	}
	return result
}

func flattenAttachUserInfos(raw interface{}) []map[string]interface{} {
	curArray := raw.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(curArray))
	for i, volume := range curArray {
		result[i] = map[string]interface{}{
			"user_name":  utils.PathSearch("user_name", volume, nil),
			"user_group": utils.PathSearch("user_group", volume, nil),
		}
	}
	return result
}
