package workspace

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

// @API Workspace GET /v2/{project_id}/desktop-pools
func DataSourceDesktopPools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDesktopPoolsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the desktop pools.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the desktop pool.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the desktop pool.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The enterprise project ID to which the desktop pool belongs",
			},
			"in_maintenance_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the desktop pool is in maintenance mode.",
			},
			"desktop_pools": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        desktopPoolSchema(),
				Description: "The list of desktop pools.",
			},
		},
	}
}

func desktopPoolSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the desktop pool.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the desktop pool.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the desktop pool.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the desktop pool.",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the desktop pool.",
			},
			"charging_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The charging mode of the desktop pool.",
			},
			"desktop_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of desktops in the pool.",
			},
			"desktop_used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of used desktops in the pool.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The availability zone of the desktop pool.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The subnet ID of the desktop pool.",
			},
			"product": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        desktopPoolProductSchema(),
				Description: "The product information of the desktop pool.",
			},
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The image ID used by the desktop pool.",
			},
			"image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The image name used by the desktop pool.",
			},
			"image_os_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The OS type of the image.",
			},
			"image_os_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The OS version of the image.",
			},
			"image_os_platform": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The OS platform of the image.",
			},
			"image_product_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The product code of the image.",
			},
			"root_volume": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        desktopPoolDesktopVolumeInfoSchema(),
				Description: "The root volume information of the desktop pool.",
			},
			"data_volumes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        desktopPoolDesktopVolumeInfoSchema(),
				Description: "The data volumes information of the desktop pool.",
			},
			"security_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        desktopPoolSecurityGroupSchema(),
				Description: "The security groups of the desktop pool.",
			},
			"disconnected_retention_period": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The disconnected retention period in minutes.",
			},
			"enable_autoscale": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether auto scaling is enabled.",
			},
			"autoscale_policy": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        desktopPoolAutoscalePolicySchema(),
				Description: "The auto scaling policy of the desktop pool.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the desktop pool.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The enterprise project ID.",
			},
			"in_maintenance_mode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the desktop pool is in maintenance mode.",
			},
			"desktop_name_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The desktop name policy ID.",
			},
		},
	}
}

func desktopPoolProductSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The product ID.",
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The flavor ID.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The product type.",
			},
			"cpu": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CPU specification.",
			},
			"memory": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The memory specification.",
			},
			"descriptions": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The product description.",
			},
			"charge_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The charging mode.",
			},
		},
	}
}

func desktopPoolDesktopVolumeInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The volume ID.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The volume type.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The volume size in GB.",
			},
			"resource_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource specification code.",
			},
		},
	}
}

func desktopPoolSecurityGroupSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The security group ID.",
			},
		},
	}
}

func desktopPoolAutoscalePolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"autoscale_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The auto scaling type.",
			},
			"max_auto_created": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of auto-created desktops.",
			},
			"min_idle": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The minimum number of idle desktops.",
			},
			"once_auto_created": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of desktops will create in one auto scaling operation.",
			},
		},
	}
}

func buildListDesktopPoolsParams(d *schema.ResourceData) string {
	params := ""
	if v, ok := d.GetOk("name"); ok {
		params = fmt.Sprintf("%s&name=%v", params, v)
	}
	if v, ok := d.GetOk("type"); ok {
		params = fmt.Sprintf("%s&type=%v", params, v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		params = fmt.Sprintf("%s&enterprise_project_id=%v", params, v)
	}
	if v, ok := d.GetOk("in_maintenance_mode"); ok {
		params = fmt.Sprintf("%s&in_maintenance_mode=%v", params, v)
	}
	return params
}

func listDesktopPools(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/desktop-pools?limit={limit}"
		offset  = 0
		limit   = 1000
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildListDesktopPoolsParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		desktopPools := utils.PathSearch("desktop_pools", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, desktopPools...)
		if len(desktopPools) < limit {
			return result, nil
		}
		offset += len(desktopPools)
	}
}

func flattenDesktopPoolsProduct(product interface{}) []interface{} {
	if product == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"product_id":   utils.PathSearch("product_id", product, nil),
			"flavor_id":    utils.PathSearch("flavor_id", product, nil),
			"type":         utils.PathSearch("type", product, nil),
			"cpu":          utils.PathSearch("cpu", product, nil),
			"memory":       utils.PathSearch("memory", product, nil),
			"descriptions": utils.PathSearch("descriptions", product, nil),
			"charge_mode":  utils.PathSearch("charge_mode", product, nil),
		},
	}
}

func flattenDesktopPoolsDesktopRootVolume(volume interface{}) []interface{} {
	if volume == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":                 utils.PathSearch("id", volume, nil),
			"type":               utils.PathSearch("type", volume, nil),
			"size":               utils.PathSearch("size", volume, nil),
			"resource_spec_code": utils.PathSearch("resource_spec_code", volume, nil),
		},
	}
}

func flattenDesktopPoolsDesktopDataVolumes(volumes []interface{}) []interface{} {
	if len(volumes) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(volumes))
	for _, item := range volumes {
		result = append(result, map[string]interface{}{
			"id":                 utils.PathSearch("id", item, nil),
			"type":               utils.PathSearch("type", item, nil),
			"size":               utils.PathSearch("size", item, nil),
			"resource_spec_code": utils.PathSearch("resource_spec_code", item, nil),
		})
	}

	return result
}

func flattenDesktopPoolsDesktopSecurityGroups(groups []interface{}) []interface{} {
	if len(groups) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(groups))
	for _, item := range groups {
		result = append(result, map[string]interface{}{
			"id": utils.PathSearch("id", item, nil),
		})
	}

	return result
}

func flattenDesktopPoolsAutoscalePolicy(policy interface{}) []interface{} {
	if policy == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"autoscale_type":    utils.PathSearch("autoscale_type", policy, nil),
			"max_auto_created":  utils.PathSearch("max_auto_created", policy, nil),
			"min_idle":          utils.PathSearch("min_idle", policy, nil),
			"once_auto_created": utils.PathSearch("once_auto_created", policy, nil),
		},
	}
}

func flattenDesktopPools(pools []interface{}) []interface{} {
	if len(pools) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(pools))
	for _, item := range pools {
		result = append(result, map[string]interface{}{
			"id":                 utils.PathSearch("id", item, nil),
			"name":               utils.PathSearch("name", item, nil),
			"type":               utils.PathSearch("type", item, nil),
			"description":        utils.PathSearch("description", item, nil),
			"created_time":       utils.PathSearch("created_time", item, nil),
			"charging_mode":      utils.PathSearch("charging_mode", item, nil),
			"desktop_count":      utils.PathSearch("desktop_count", item, nil),
			"desktop_used":       utils.PathSearch("desktop_used", item, nil),
			"availability_zone":  utils.PathSearch("availability_zone", item, nil),
			"subnet_id":          utils.PathSearch("subnet_id", item, nil),
			"product":            flattenDesktopPoolsProduct(utils.PathSearch("product", item, nil)),
			"image_id":           utils.PathSearch("image_id", item, nil),
			"image_name":         utils.PathSearch("image_name", item, nil),
			"image_os_type":      utils.PathSearch("image_os_type", item, nil),
			"image_os_version":   utils.PathSearch("image_os_version", item, nil),
			"image_os_platform":  utils.PathSearch("image_os_platform", item, nil),
			"image_product_code": utils.PathSearch("image_product_code", item, nil),
			"root_volume":        flattenDesktopPoolsDesktopRootVolume(utils.PathSearch("root_volume", item, nil)),
			"data_volumes": flattenDesktopPoolsDesktopDataVolumes(
				utils.PathSearch("data_volumes", item, make([]interface{}, 0)).([]interface{})),
			"security_groups": flattenDesktopPoolsDesktopSecurityGroups(
				utils.PathSearch("security_groups", item, make([]interface{}, 0)).([]interface{})),
			"disconnected_retention_period": utils.PathSearch("disconnected_retention_period", item, nil),
			"enable_autoscale":              utils.PathSearch("enable_autoscale", item, nil),
			"autoscale_policy":              flattenDesktopPoolsAutoscalePolicy(utils.PathSearch("autoscale_policy", item, nil)),
			"status":                        utils.PathSearch("status", item, nil),
			"enterprise_project_id":         utils.PathSearch("enterprise_project_id", item, nil),
			"in_maintenance_mode":           utils.PathSearch("in_maintenance_mode", item, nil),
			"desktop_name_policy_id":        utils.PathSearch("desktop_name_policy_id", item, nil),
		})
	}

	return result
}

func dataSourceDesktopPoolsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	desktopPools, err := listDesktopPools(client, d)
	if err != nil {
		return diag.Errorf("error querying Workspace desktop pools: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("desktop_pools", flattenDesktopPools(desktopPools)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
