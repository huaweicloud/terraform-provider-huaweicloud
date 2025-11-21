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

// @API Workspace GET /v1/{project_id}/app-server-groups
func DataSourceAppServerGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppServerGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the resource.`,
			},
			"server_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the server group.`,
			},
			"server_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the server group.`,
			},
			"app_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of application group.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The enterprise project ID.`,
			},
			"is_secondary_server_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Whether it is a secondary server group.`,
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The tag value to filter server groups.`,
			},

			// Attribute
			"server_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of server groups.`,
				Elem:        serverGroupSchema(),
			},
		},
	}
}

func serverGroupSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique ID of the server group.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the server group.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the server group.`,
			},
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The image ID used to create servers in this group.`,
			},
			"os_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the operating system.`,
			},
			"product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The product ID.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subnet ID for the network interface.`,
			},
			"system_disk_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the system disk.`,
			},
			"system_disk_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The size of the system disk.`,
			},
			"is_vdi": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether it is VDI single-session mode.`,
			},
			"extra_session_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The paid session type.`,
			},
			"extra_session_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of paid sessions.`,
			},
			"app_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The application type of app server groups that matched filter parameters.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the server group.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The last update time of the server group.`,
			},
			"storage_mount_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The NAS storage directory mounting policy on APS.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise project ID.`,
			},
			"primary_server_group_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of primary server group IDs.`,
			},
			"secondary_server_group_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of secondary server group IDs.`,
			},
			"server_group_status": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the server group is enabled.`,
			},
			"site_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The site type.`,
			},
			"site_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The site ID.`,
			},
			"app_server_flavor_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total number of server configurations.`,
			},
			"app_server_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total number of servers.`,
			},
			"app_group_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total number of associated application groups.`,
			},
			"image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The image name of the group.`,
			},
			"subnet_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subnet name of the group.`,
			},
			"ou_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The default organization name of the group.`,
			},
			"product_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The product specification information.`,
				Elem:        productInfoSchema(),
			},
			"scaling_policy": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The auto-scaling policy.`,
				Elem:        scalingPolicySchema(),
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The tag information of app server groups that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The key of the tag.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The value of the tag.`,
						},
					},
				},
			},
		},
	}
}

func productInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the product.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the flavor.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the product.`,
			},
			"architecture": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The architecture of the product.`,
			},
			"cpu": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The CPU information.`,
			},
			"cpu_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The CPU description.`,
			},
			"memory": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The memory size in MB.`,
			},
			"is_gpu": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether it is GPU type.`,
			},
			"system_disk_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the system disk.`,
			},
			"system_disk_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The size of the system disk.`,
			},
			"gpu_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The GPU description.`,
			},
			"descriptions": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The product description.`,
			},
			"charge_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The charging mode.`,
			},
			"contain_data_disk": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the package includes data disk.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the resource.`,
			},
			"cloud_service_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the cloud service.`,
			},
			"volume_product_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the volume product.`,
			},
			"sessions": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum number of sessions supported by the package.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the product package in sales mode.`,
			},
			"cond_operation_az": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the product package in the availability zone.`,
			},
			"sub_product_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of sub products.`,
			},
			"domain_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of domain IDs.`,
			},
			"package_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the package.`,
			},
			"expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The expiration time of the product package.`,
			},
			"support_gpu_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The GPU type supported by the product package.`,
			},
		},
	}
}

func scalingPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable the policy.`,
			},
			"max_scaling_amount": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum scaling amount.`,
			},
			"single_expansion_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of instances to add in a single scaling operation.`,
			},
			"scaling_policy_by_session": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The session-based scaling policy.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"session_usage_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The total session usage threshold of the group.`,
						},
						"shrink_after_session_idle_minutes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The release time for instances without session connections.`,
						},
					},
				},
			},
		},
	}
}

func buildAppServerGroupsQueryParams(d *schema.ResourceData) string {
	var params []string

	if v, ok := d.GetOk("server_group_name"); ok {
		params = append(params, fmt.Sprintf("server_group_name=%v", v))
	}
	if v, ok := d.GetOk("server_group_id"); ok {
		params = append(params, fmt.Sprintf("server_group_id=%v", v))
	}
	if v, ok := d.GetOk("app_type"); ok {
		params = append(params, fmt.Sprintf("app_type=%v", v))
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		params = append(params, fmt.Sprintf("enterprise_project_id=%v", v))
	}
	if v, ok := d.GetOk("is_secondary_server_group"); ok {
		params = append(params, fmt.Sprintf("is_secondary_server_group=%v", v))
	}
	if v, ok := d.GetOk("tags"); ok {
		params = append(params, fmt.Sprintf("tags=%v", v))
	}

	if len(params) < 1 {
		return ""
	}
	return "&" + strings.Join(params, "&")
}

func listAppServerGroups(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/app-server-groups?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildAppServerGroupsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		items := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, items...)
		if len(items) < limit {
			break
		}

		offset += len(items)
	}

	return result, nil
}

func flattenAppServerGroups(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(items))
	for i, item := range items {
		result[i] = map[string]interface{}{
			"id":                 utils.PathSearch("id", item, nil),
			"name":               utils.PathSearch("name", item, nil),
			"description":        utils.PathSearch("description", item, nil),
			"image_id":           utils.PathSearch("image_id", item, nil),
			"os_type":            utils.PathSearch("os_type", item, nil),
			"product_id":         utils.PathSearch("product_id", item, nil),
			"subnet_id":          utils.PathSearch("subnet_id", item, nil),
			"system_disk_type":   utils.PathSearch("system_disk_type", item, nil),
			"system_disk_size":   utils.PathSearch("system_disk_size", item, nil),
			"is_vdi":             utils.PathSearch("is_vdi", item, nil),
			"extra_session_type": utils.PathSearch("extra_session_type", item, nil),
			"extra_session_size": utils.PathSearch("extra_session_size", item, nil),
			"app_type":           utils.PathSearch("app_type", item, nil),
			"create_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
				utils.PathSearch("create_time", item, "").(string))/1000, false),
			"update_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
				utils.PathSearch("update_time", item, "").(string))/1000, false),
			"storage_mount_policy":       utils.PathSearch("storage_mount_policy", item, nil),
			"enterprise_project_id":      utils.PathSearch("enterprise_project_id", item, nil),
			"primary_server_group_ids":   utils.PathSearch("primary_server_group_ids", item, nil),
			"secondary_server_group_ids": utils.PathSearch("secondary_server_group_ids", item, nil),
			"server_group_status":        utils.PathSearch("server_group_status", item, nil),
			"site_type":                  utils.PathSearch("site_type", item, nil),
			"site_id":                    utils.PathSearch("site_id", item, nil),
			"app_server_flavor_count":    utils.PathSearch("app_server_flavor_count", item, nil),
			"app_server_count":           utils.PathSearch("app_server_count", item, nil),
			"app_group_count":            utils.PathSearch("app_group_count", item, nil),
			"image_name":                 utils.PathSearch("image_name", item, nil),
			"subnet_name":                utils.PathSearch("subnet_name", item, nil),
			"ou_name":                    utils.PathSearch("ou_name", item, nil),
			"product_info": flattenProductInfo(utils.PathSearch("product_info", item,
				make(map[string]interface{})).(map[string]interface{})),
			"scaling_policy": flattenScalingPolicy(utils.PathSearch("scaling_policy", item,
				make(map[string]interface{})).(map[string]interface{})),
			"tags": flattenAppServerGroupsTags(utils.PathSearch("tags", item,
				make([]interface{}, 0)).([]interface{})),
		}
	}

	return result
}

func flattenProductInfo(rawData map[string]interface{}) []map[string]interface{} {
	if len(rawData) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"product_id":          utils.PathSearch("product_id", rawData, nil),
			"flavor_id":           utils.PathSearch("flavor_id", rawData, nil),
			"type":                utils.PathSearch("type", rawData, nil),
			"architecture":        utils.PathSearch("architecture", rawData, nil),
			"cpu":                 utils.PathSearch("cpu", rawData, nil),
			"cpu_desc":            utils.PathSearch("cpu_desc", rawData, nil),
			"memory":              utils.PathSearch("memory", rawData, nil),
			"is_gpu":              utils.PathSearch("is_gpu", rawData, nil),
			"system_disk_type":    utils.PathSearch("system_disk_type", rawData, nil),
			"system_disk_size":    utils.PathSearch("system_disk_size", rawData, nil),
			"gpu_desc":            utils.PathSearch("gpu_desc", rawData, nil),
			"descriptions":        utils.PathSearch("descriptions", rawData, nil),
			"charge_mode":         utils.PathSearch("charge_mode", rawData, nil),
			"contain_data_disk":   utils.PathSearch("contain_data_disk", rawData, nil),
			"resource_type":       utils.PathSearch("resource_type", rawData, nil),
			"cloud_service_type":  utils.PathSearch("cloud_service_type", rawData, nil),
			"volume_product_type": utils.PathSearch("volume_product_type", rawData, nil),
			"sessions":            utils.PathSearch("sessions", rawData, nil),
			"status":              utils.PathSearch("status", rawData, nil),
			"cond_operation_az":   utils.PathSearch("cond_operation_az", rawData, nil),
			"sub_product_list":    utils.PathSearch("sub_product_list", rawData, make([]interface{}, 0)).([]interface{}),
			"domain_ids":          utils.PathSearch("domain_ids", rawData, make([]interface{}, 0)).([]interface{}),
			"package_type":        utils.PathSearch("package_type", rawData, nil),
			"expire_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
				utils.PathSearch("expire_time", rawData, "").(string))/1000, false),
			"support_gpu_type": utils.PathSearch("support_gpu_type", rawData, nil),
		},
	}
}

func flattenScalingPolicy(rawData map[string]interface{}) []map[string]interface{} {
	if len(rawData) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"enable":                 utils.PathSearch("enable", rawData, nil),
			"max_scaling_amount":     utils.PathSearch("max_scaling_amount", rawData, nil),
			"single_expansion_count": utils.PathSearch("single_expansion_count", rawData, nil),
			"scaling_policy_by_session": flattenScalingPolicyBySession(utils.PathSearch("scaling_policy_by_session",
				rawData, make(map[string]interface{})).(map[string]interface{})),
		},
	}
}

func flattenScalingPolicyBySession(rawData map[string]interface{}) []map[string]interface{} {
	if len(rawData) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"session_usage_threshold":           utils.PathSearch("session_usage_threshold", rawData, nil),
			"shrink_after_session_idle_minutes": utils.PathSearch("shrink_after_session_idle_minutes", rawData, nil),
		},
	}
}

func flattenAppServerGroupsTags(rawData []interface{}) []map[string]interface{} {
	if len(rawData) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(rawData))
	for i, v := range rawData {
		tagMap := v.(map[string]interface{})
		result[i] = map[string]interface{}{
			"key":   utils.PathSearch("key", tagMap, nil),
			"value": utils.PathSearch("value", tagMap, nil),
		}
	}

	return result
}

func dataSourceAppServerGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	items, err := listAppServerGroups(client, d)
	if err != nil {
		return diag.Errorf("error querying server groups: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("server_groups", flattenAppServerGroups(items)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
