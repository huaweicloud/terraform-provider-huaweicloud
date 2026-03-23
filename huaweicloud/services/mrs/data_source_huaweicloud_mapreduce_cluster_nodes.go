package mrs

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

// @API MRS GET /v2/{project_id}/clusters/{cluster_id}/nodes
func DataSourceClusterNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceClusterNodesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the nodes are located.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster.`,
			},
			"node_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the node group to which the node belongs.`,
			},
			"node_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the node.`,
			},
			"query_node_detail": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to query node detail.`,
			},
			"query_ecs_detail": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to query ECS detail.`,
			},
			"internal_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The internal IP address of the node.`,
			},
			"nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of nodes that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the node.`,
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource ID of the node.`,
						},
						"node_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the node group to which the node belongs.`,
						},
						"node_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the node.`,
						},
						"charging_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The billing type of the node.`,
						},
						"deployment_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The deployment type of the node.`,
						},
						"server_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The server information of the node.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the server.`,
									},
									"server_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the server.`,
									},
									"server_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the server.`,
									},
									"data_volumes": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The data disks of the server.`,
										Elem:        clusterNodesVolumeSchema(),
									},
									"root_volume": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The system disk configuration of the server.`,
										Elem:        clusterNodesVolumeSchema(),
									},
									"cpu_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The CPU type of the server.`,
									},
									"cpu": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The CPU size of the server.`,
									},
									"mem": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The memory size of the server, in MB.`,
									},
									"internal_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The internal IP address of the server.`,
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The key/value pairs associated with the node.`,
						},
						"node_detail": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The monitoring information of the node.`,
							Elem:        clusterNodesNodeDetailSchema(),
						},
						"node_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the node.`,
						},
						"component_infos": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of components deployed on the node.`,
							Elem:        clusterNodesComponentInfoSchema(),
						},
					},
				},
			},
		},
	}
}

func clusterNodesVolumeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The disk type.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The disk size, in GB.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The disk count.`,
			},
		},
	}
}

func clusterNodesNodeDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"running_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The running status.`,
			},
			"cpu_usage": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The CPU usage.`,
			},
			"memory_usage": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The memory usage.`,
			},
			"disk_usage": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The disk usage.`,
			},
			"total_memory": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The total memory, in MB.`,
			},
			"available_memory": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The available memory, in MB.`,
			},
			"total_hard_disk_space": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The total hard disk space, in GB.`,
			},
			"available_hard_disk_space": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The available hard disk space, in GB.`,
			},
			"network_read": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The network read speed, in Byte/s.`,
			},
			"network_write": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The network write speed, in Byte/s.`,
			},
		},
	}
}

func clusterNodesComponentInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The component ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The component name.`,
			},
			"instance_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The component instance group name.`,
			},
			"running_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The component running status.`,
			},
			"ha_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The HA status.`,
			},
			"config_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The config status.`,
			},
			"role_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The role name.`,
			},
			"role_short_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The role short name.`,
			},
			"role_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The role type.`,
			},
			"service_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The service name.`,
			},
			"pair_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The pair name.`,
			},
			"relation_pairs": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The relation pairs.`,
			},
			"support_decom": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Whether Decom is supported.`,
			},
			"support_reinstall": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Whether reinstall is supported.`,
			},
			"support_collect_stack_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Whether stack info collection is supported.`,
			},
		},
	}
}

func listClusterNodes(client *golangsdk.ServiceClient, d *schema.ResourceData, clusterId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/nodes"
		// `offset` means the page number, starts from 1.
		offset = 1
		limit  = 100
		result = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{cluster_id}", clusterId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	listPath = fmt.Sprintf("%s?limit=%d%s", listPath, limit, buildListClusterNodesQueryParams(d))
	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		nodes := utils.PathSearch("nodes", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, nodes...)
		if len(nodes) < limit {
			break
		}

		offset++
	}

	return result, nil
}

func buildListClusterNodesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("node_group"); ok {
		res = fmt.Sprintf("%s&node_group=%v", res, v)
	}
	if v, ok := d.GetOk("node_name"); ok {
		res = fmt.Sprintf("%s&node_name=%v", res, v)
	}
	if v, ok := d.GetOk("query_node_detail"); ok {
		res = fmt.Sprintf("%s&query_node_detail=%v", res, v)
	}
	if v, ok := d.GetOk("query_ecs_detail"); ok {
		res = fmt.Sprintf("%s&query_ecs_detail=%v", res, v)
	}
	if v, ok := d.GetOk("internal_ip"); ok {
		res = fmt.Sprintf("%s&internal_ip=%v", res, v)
	}

	return res
}

func resourceClusterNodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient("mrs", region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	nodes, err := listClusterNodes(client, d, clusterId)
	if err != nil {
		return diag.Errorf("error retrieving nodes under the specified cluster (%s): %s", clusterId, err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("nodes", flattenClusterNodes(nodes)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenClusterNodes(nodes []interface{}) []interface{} {
	if len(nodes) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(nodes))
	for _, v := range nodes {
		result = append(result, map[string]interface{}{
			"node_name":       utils.PathSearch("node_name", v, nil),
			"resource_id":     utils.PathSearch("resource_id", v, nil),
			"node_group_name": utils.PathSearch("node_group_name", v, nil),
			"node_type":       utils.PathSearch("node_type", v, nil),
			"charging_mode":   parseClusterNodeChargingMode(utils.PathSearch("billing_type", v, "").(string)),
			"deployment_type": utils.PathSearch("deployment_type", v, nil),
			"server_info":     flattenClusterNodeServerInfo(utils.PathSearch("server_info", v, nil)),
			"tags":            utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil)),
			"node_detail":     flattenClusterNodeDetail(utils.PathSearch("node_detail", v, nil)),
			"node_status":     utils.PathSearch("node_status", v, nil),
			"component_infos": flattenClusterNodeComponentInfos(utils.PathSearch("component_infos",
				v, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func parseClusterNodeChargingMode(chargingMode string) string {
	switch chargingMode {
	case "on-period":
		return "prePaid"
	case "on-quantity":
		return "postPaid"
	}

	return chargingMode
}

func flattenClusterNodeServerInfo(serverInfo interface{}) []interface{} {
	if serverInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"server_id":   utils.PathSearch("server_id", serverInfo, nil),
			"server_name": utils.PathSearch("server_name", serverInfo, nil),
			"server_type": utils.PathSearch("server_type", serverInfo, nil),
			"data_volumes": flattenClusterNodesServerInfoVolumes(utils.PathSearch("data_volumes",
				serverInfo, make([]interface{}, 0)).([]interface{})),
			"root_volume": flattenClusterNodesServerInfoRootVolume(utils.PathSearch("root_volume", serverInfo, nil)),
			"cpu_type":    utils.PathSearch("cpu_type", serverInfo, nil),
			"cpu":         utils.PathSearch("cpu", serverInfo, nil),
			"mem":         utils.PathSearch("mem", serverInfo, nil),
			"internal_ip": utils.PathSearch("internal_ip", serverInfo, nil),
		},
	}
}

func flattenClusterNodesServerInfoVolumes(volumes []interface{}) []interface{} {
	if len(volumes) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(volumes))
	for _, v := range volumes {
		result = append(result, map[string]interface{}{
			"type":  utils.PathSearch("type", v, nil),
			"size":  utils.PathSearch("size", v, nil),
			"count": utils.PathSearch("count", v, nil),
		})
	}
	return result
}

func flattenClusterNodesServerInfoRootVolume(rootVolume interface{}) []interface{} {
	if rootVolume == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"type":  utils.PathSearch("type", rootVolume, nil),
			"size":  utils.PathSearch("size", rootVolume, nil),
			"count": utils.PathSearch("count", rootVolume, nil),
		},
	}
}

func flattenClusterNodeDetail(nodeDetail interface{}) []interface{} {
	if nodeDetail == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"running_status":            utils.PathSearch("running_status", nodeDetail, nil),
			"cpu_usage":                 utils.PathSearch("cpu_usage", nodeDetail, nil),
			"memory_usage":              utils.PathSearch("memory_usage", nodeDetail, nil),
			"disk_usage":                utils.PathSearch("disk_usage", nodeDetail, nil),
			"total_memory":              utils.PathSearch("total_memory", nodeDetail, nil),
			"available_memory":          utils.PathSearch("available_memory", nodeDetail, nil),
			"total_hard_disk_space":     utils.PathSearch("total_hard_disk_space", nodeDetail, nil),
			"available_hard_disk_space": utils.PathSearch("available_hard_disk_space", nodeDetail, nil),
			"network_read":              utils.PathSearch("network_read", nodeDetail, nil),
			"network_write":             utils.PathSearch("network_write", nodeDetail, nil),
		},
	}
}

func flattenClusterNodeComponentInfos(components []interface{}) []interface{} {
	if len(components) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(components))
	for _, v := range components {
		result = append(result, map[string]interface{}{
			"id":                         utils.PathSearch("id", v, nil),
			"name":                       utils.PathSearch("name", v, nil),
			"instance_group_name":        utils.PathSearch("instance_group_name", v, nil),
			"running_status":             utils.PathSearch("running_status", v, nil),
			"ha_status":                  utils.PathSearch("ha_status", v, nil),
			"config_status":              utils.PathSearch("config_status", v, nil),
			"role_name":                  utils.PathSearch("role_name", v, nil),
			"role_short_name":            utils.PathSearch("role_short_name", v, nil),
			"role_type":                  utils.PathSearch("role_type", v, nil),
			"service_name":               utils.PathSearch("service_name", v, nil),
			"pair_name":                  utils.PathSearch("pair_name", v, nil),
			"relation_pairs":             utils.PathSearch("relation_pairs", v, nil),
			"support_decom":              utils.PathSearch("support_decom", v, nil),
			"support_reinstall":          utils.PathSearch("support_reinstall", v, nil),
			"support_collect_stack_info": utils.PathSearch("support_collect_stack_info", v, nil),
		})
	}

	return result
}
