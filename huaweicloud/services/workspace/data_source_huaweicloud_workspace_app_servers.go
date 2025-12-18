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

// @API Workspace GET /v1/{project_id}/app-servers
func DataSourceAppServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppServersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the resource.`,
			},
			"server_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the server group.`,
			},
			"server_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the server.`,
			},
			"machine_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The machine name of the server.`,
			},
			"ip_addr": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The IP address of the server.`,
			},
			"server_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the server.`,
			},
			"maintain_status": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether the server is in maintenance status.`,
			},
			"scaling_auto_create": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether the server is created by auto-scaling.`,
			},

			// Attributes
			"servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of servers.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the server.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the server.`,
						},
						"machine_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The machine name of the server group that matched filter parameters.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the server.`,
						},
						"server_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the server group that matched filter parameters.`,
						},
						"server_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the server group.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the server.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the server.`,
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time of the server.`,
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the image.`,
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The availability zone of the server.`,
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain of the server.`,
						},
						"ou_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The organization unit name.`,
						},
						"sid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The SID of the instance.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the instance.`,
						},
						"os_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version of the operating system.`,
						},
						"os_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the operating system.`,
						},
						"order_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the order.`,
						},
						"maintain_status": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `The maintenance status of the server group that matched filter parameters.`,
						},
						"scaling_auto_create": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the server is created by auto scaling.`,
						},
						"job_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the last executed job.`,
						},
						"job_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the job.`,
						},
						"job_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the job.`,
						},
						"job_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The execution time of the last job.`,
						},
						"resource_pool_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the resource pool.`,
						},
						"resource_pool_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the resource pool.`,
						},
						"host_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the dedicated host.`,
						},
						"session_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of sessions.`,
						},
						"vm_status": {
							Type:     schema.TypeString,
							Computed: true,
							Description: `The steady state of a server, the stable state in which a certain operation
                                         is completed.`,
						},
						"task_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The task status of the server.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the enterprise project.`,
						},
						"flavor": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        workspaceAppServersFlavorSchema(),
							Description: `The flavor information of the server.`,
						},
						"product_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The product information of the server.`,
							Elem: &schema.Resource{
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
										Description: `Whether the flavor is GPU type.`,
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
							},
						},
						"metadata": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The metadata of the server.`,
						},
						"freeze": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The freeze information of the server.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"effect": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The effect of the freeze operation.`,
									},
									"scene": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The scene of the service status update.`,
									},
								},
							},
						},
						"host_address": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The network information of the server.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addr": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The IP address.`,
									},
									"version": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The IP address version.`,
									},
									"mac_addr": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The MAC address.`,
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The IP address allocation type.`,
									},
									"port_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The port ID of the IP address.`,
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the VPC.`,
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the subnet.`,
									},
									"tenant_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the tenant.`,
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The tags of the server.`,
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
				},
			},
		},
	}
}

func workspaceAppServersFlavorSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the flavor.`,
			},
			"links": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The quick link information for relevant tags corresponding to server specifications.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rel": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The shortcut link tag name.`,
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The corresponding shortcut link.`,
						},
					},
				},
			},
		},
	}
}

func buildAppServersQueryParams(d *schema.ResourceData) string {
	res := ""

	if serverGroupId, ok := d.GetOk("server_group_id"); ok {
		res = fmt.Sprintf("%s&server_group_id=%v", res, serverGroupId)
	}
	if serverName, ok := d.GetOk("server_name"); ok {
		res = fmt.Sprintf("%s&server_name=%v", res, serverName)
	}
	if machineName, ok := d.GetOk("machine_name"); ok {
		res = fmt.Sprintf("%s&machine_name=%v", res, machineName)
	}
	if ipAddr, ok := d.GetOk("ip_addr"); ok {
		res = fmt.Sprintf("%s&ip_addr=%v", res, ipAddr)
	}
	if serverId, ok := d.GetOk("server_id"); ok {
		res = fmt.Sprintf("%s&server_id=%v", res, serverId)
	}
	if maintainStatus, ok := d.GetOk("maintain_status"); ok {
		res = fmt.Sprintf("%s&maintain_status=%v", res, maintainStatus)
	}
	if scalingAutoCreate, ok := d.GetOk("scaling_auto_create"); ok {
		res = fmt.Sprintf("%s&scaling_auto_create=%v", res, scalingAutoCreate)
	}

	return res
}

func listAppServers(client *golangsdk.ServiceClient, queryParams ...string) ([]interface{}, error) {
	var (
		httpUrl  = "v1/{project_id}/app-servers?limit={limit}"
		limit    = 100
		offset   = 0
		result   = make([]interface{}, 0)
		respBody interface{}
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	if len(queryParams) > 0 && queryParams[0] != "" {
		listPath += queryParams[0]
	}

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

		respBody, err = utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		servers := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, servers...)
		if len(servers) < limit {
			break
		}

		offset += len(servers)
	}

	return result, nil
}

func flattenAppServersFlavorLinks(links []interface{}) []map[string]interface{} {
	if len(links) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(links))
	for _, link := range links {
		result = append(result, map[string]interface{}{
			"href": utils.PathSearch("href", link, nil),
			"rel":  utils.PathSearch("rel", link, nil),
		})
	}

	return result
}

func flattenAppServersFlavor(flavor map[string]interface{}) []map[string]interface{} {
	if len(flavor) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"id": utils.PathSearch("id", flavor, nil),
			"links": flattenAppServersFlavorLinks(utils.PathSearch("links", flavor,
				make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenAppServersProductInfo(productInfo map[string]interface{}) []map[string]interface{} {
	if len(productInfo) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"product_id":          utils.PathSearch("product_id", productInfo, nil),
			"flavor_id":           utils.PathSearch("flavor_id", productInfo, nil),
			"type":                utils.PathSearch("type", productInfo, nil),
			"architecture":        utils.PathSearch("architecture", productInfo, nil),
			"cpu":                 utils.PathSearch("cpu", productInfo, nil),
			"cpu_desc":            utils.PathSearch("cpu_desc", productInfo, nil),
			"memory":              utils.PathSearch("memory", productInfo, nil),
			"is_gpu":              utils.PathSearch("is_gpu", productInfo, nil),
			"system_disk_type":    utils.PathSearch("system_disk_type", productInfo, nil),
			"system_disk_size":    utils.PathSearch("system_disk_size", productInfo, nil),
			"gpu_desc":            utils.PathSearch("gpu_desc", productInfo, nil),
			"descriptions":        utils.PathSearch("descriptions", productInfo, nil),
			"charge_mode":         utils.PathSearch("charge_mode", productInfo, nil),
			"contain_data_disk":   utils.PathSearch("contain_data_disk", productInfo, nil),
			"resource_type":       utils.PathSearch("resource_type", productInfo, nil),
			"cloud_service_type":  utils.PathSearch("cloud_service_type", productInfo, nil),
			"volume_product_type": utils.PathSearch("volume_product_type", productInfo, nil),
			"sessions":            utils.PathSearch("sessions", productInfo, nil),
			"status":              utils.PathSearch("status", productInfo, nil),
			"cond_operation_az":   utils.PathSearch("cond_operation_az", productInfo, nil),
			"sub_product_list":    utils.PathSearch("sub_product_list", productInfo, nil),
			"domain_ids":          utils.PathSearch("domain_ids", productInfo, nil),
			"package_type":        utils.PathSearch("package_type", productInfo, nil),
			"expire_time":         utils.PathSearch("expire_time", productInfo, nil),
			"support_gpu_type":    utils.PathSearch("support_gpu_type", productInfo, nil),
		},
	}
}

func flattenV1AppServersFreeze(freeze []interface{}) []map[string]interface{} {
	if len(freeze) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"effect": utils.PathSearch("effect", freeze, nil),
			"scene":  utils.PathSearch("scene", freeze, nil),
		},
	}
}

func flattenV1AppServersHostAddress(addressList []interface{}) []map[string]interface{} {
	if len(addressList) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(addressList))
	for _, addr := range addressList {
		result = append(result, map[string]interface{}{
			"addr":        utils.PathSearch("addr", addr, nil),
			"version":     utils.PathSearch("version", addr, nil),
			"mac_addr":    utils.PathSearch(`"OS-EXT-IPS-MAC:mac_addr"`, addr, nil),
			"type":        utils.PathSearch(`"OS-EXT-IPS:type"`, addr, nil),
			"port_id":     utils.PathSearch(`"OS-EXT-IPS:port_id"`, addr, nil),
			"vpc_id":      utils.PathSearch("vpc_id", addr, nil),
			"subnet_id":   utils.PathSearch("subnet_id", addr, nil),
			"tenant_type": utils.PathSearch("tenant_type", addr, nil),
		})
	}

	return result
}

func flattenV1AppServersTags(tagsList []interface{}) []map[string]interface{} {
	if len(tagsList) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(tagsList))
	for _, t := range tagsList {
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", t, nil),
			"value": utils.PathSearch("value", t, nil),
		})
	}

	return result
}

func flattenAppServers(servers []interface{}) []map[string]interface{} {
	if len(servers) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(servers))
	for _, server := range servers {
		result = append(result, map[string]interface{}{
			"id":                    utils.PathSearch("id", server, nil),
			"name":                  utils.PathSearch("name", server, nil),
			"machine_name":          utils.PathSearch("machine_name", server, nil),
			"description":           utils.PathSearch("description", server, nil),
			"server_group_id":       utils.PathSearch("server_group_id", server, nil),
			"server_group_name":     utils.PathSearch("server_group_name", server, nil),
			"status":                utils.PathSearch("status", server, nil),
			"create_time":           utils.PathSearch("create_time", server, nil),
			"update_time":           utils.PathSearch("update_time", server, nil),
			"image_id":              utils.PathSearch("image_id", server, nil),
			"availability_zone":     utils.PathSearch("availability_zone", server, nil),
			"domain":                utils.PathSearch("domain", server, nil),
			"ou_name":               utils.PathSearch("ou_name", server, nil),
			"sid":                   utils.PathSearch("sid", server, nil),
			"instance_id":           utils.PathSearch("instance_id", server, nil),
			"os_version":            utils.PathSearch("os_version", server, nil),
			"os_type":               utils.PathSearch("os_type", server, nil),
			"order_id":              utils.PathSearch("order_id", server, nil),
			"maintain_status":       utils.PathSearch("maintain_status", server, nil),
			"scaling_auto_create":   utils.PathSearch("scaling_auto_create", server, nil),
			"job_id":                utils.PathSearch("job_id", server, nil),
			"job_type":              utils.PathSearch("job_type", server, nil),
			"job_status":            utils.PathSearch("job_status", server, nil),
			"job_time":              utils.PathSearch("job_time", server, nil),
			"resource_pool_id":      utils.PathSearch("resource_pool_id", server, nil),
			"resource_pool_type":    utils.PathSearch("resource_pool_type", server, nil),
			"host_id":               utils.PathSearch("host_id", server, nil),
			"session_count":         utils.PathSearch("session_count", server, nil),
			"vm_status":             utils.PathSearch("vm_status", server, nil),
			"task_status":           utils.PathSearch("task_status", server, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", server, nil),
			"metadata":              utils.PathSearch("metadata", server, nil),
			"flavor": flattenAppServersFlavor(utils.PathSearch("flavor", server,
				make(map[string]interface{})).(map[string]interface{})),
			"product_info": flattenAppServersProductInfo(utils.PathSearch("product_info", server,
				make(map[string]interface{})).(map[string]interface{})),
			"freeze": flattenV1AppServersFreeze(utils.PathSearch("freeze", server,
				make([]interface{}, 0)).([]interface{})),
			"host_address": flattenV1AppServersHostAddress(utils.PathSearch("host_address", server,
				make([]interface{}, 0)).([]interface{})),
			"tags": flattenV1AppServersTags(utils.PathSearch("tags", server,
				make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func dataSourceAppServersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	servers, err := listAppServers(client, buildAppServersQueryParams(d))
	if err != nil {
		return diag.Errorf("error getting app servers: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("servers", flattenAppServers(servers)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
