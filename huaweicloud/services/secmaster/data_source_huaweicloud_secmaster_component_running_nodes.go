package secmaster

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

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/components/{component_id}/running-nodes
func DataSourceComponentRunningNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComponentRunningNodesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the workspace ID.",
			},
			"component_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the component ID.",
			},
			"node_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the node ID.",
			},
			"node_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the node name.",
			},
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the attribute fields for sorting.",
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the sorting order. Supported values are **ASC** and **DESC**.",
			},
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The component running nodes list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"component_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The component ID.",
						},
						"component_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The component name.",
						},
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The node ID.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The creation time (timestamp in milliseconds).",
						},
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The node name.",
						},
						"specification": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The specification.",
						},
						"config_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The node configuration status.",
						},
						"fail_deploy_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The deployment failure message.",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address.",
						},
						"private_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The private IP address.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region.",
						},
						"vpc_endpoint_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC endpoint ID.",
						},
						"vpc_endpoint_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC endpoint address.",
						},
						"monitor": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The monitor information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mini_on_online": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether online.",
									},
									"memory_count": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The number of physical memory modules.",
									},
									"memory_usage": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The amount of physical memory used.",
									},
									"memory_free": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The amount of currently free physical memory.",
									},
									"memory_shared": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The total amount of memory shared by multiple processes.",
									},
									"memory_cache": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The memory size of cached data.",
									},
									"cpu_usage": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The current CPU usage rate.",
									},
									"cpu_idle": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The percentage of CPU idle time.",
									},
									"up_pps": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The number of upload data packets per second.",
									},
									"down_pps": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The number of download data packets per second.",
									},
									"write_rate": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The disk write rate.",
									},
									"read_rate": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The disk read rate.",
									},
									"disk_count": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The number of disk devices in the system.",
									},
									"disk_usage": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The current disk space usage.",
									},
									"heart_beat_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The time when the last heartbeat signal was received, in ISO 8601 format.",
									},
									"health_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The health status of the node.",
									},
									"heart_beat": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether the node successfully received heartbeat signals.",
									},
								},
							},
						},
						"node_expansion": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The node expansion information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The node ID.",
									},
									"data_center": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The data center.",
									},
									"custom_label": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The custom label.",
									},
									"network_plane": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The network plane.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The description information.",
									},
									"maintainer": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The maintainer.",
									},
								},
							},
						},
						"node_apply_fail_enum": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The node application success or failure status and reason.",
						},
						"list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The component configuration parameter list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"configuration_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The configuration ID.",
									},
									"component_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The component ID.",
									},
									"node_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The node ID.",
									},
									"file_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The file name.",
									},
									"file_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The file path.",
									},
									"file_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The file type.",
									},
									"param": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The parameter.",
									},
									"version": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The version.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The configuration type.",
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

func buildComponentRunningNodesQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("node_id"); ok {
		queryParams = fmt.Sprintf("%s&node_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("node_name"); ok {
		queryParams = fmt.Sprintf("%s&node_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_key"); ok {
		queryParams = fmt.Sprintf("%s&sort_key=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams = fmt.Sprintf("%s&sort_dir=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceComponentRunningNodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/components/{component_id}/running-nodes"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{component_id}", d.Get("component_id").(string))
	requestPath += buildComponentRunningNodesQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		requestResp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster component running nodes: %s", err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}

		recordsResp := utils.PathSearch("records", respBody, make([]interface{}, 0)).([]interface{})
		if len(recordsResp) == 0 {
			break
		}

		result = append(result, recordsResp...)
		offset += len(recordsResp)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenComponentRunningNodesRecords(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenComponentRunningNodesRecords(recordsResp []interface{}) []interface{} {
	if len(recordsResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(recordsResp))
	for _, v := range recordsResp {
		rst = append(rst, map[string]interface{}{
			"component_id":         utils.PathSearch("component_id", v, nil),
			"component_name":       utils.PathSearch("component_name", v, nil),
			"node_id":              utils.PathSearch("node_id", v, nil),
			"create_time":          utils.PathSearch("create_time", v, nil),
			"node_name":            utils.PathSearch("node_name", v, nil),
			"specification":        utils.PathSearch("specification", v, nil),
			"config_status":        utils.PathSearch("config_status", v, nil),
			"fail_deploy_message":  utils.PathSearch("fail_deploy_message", v, nil),
			"ip_address":           utils.PathSearch("ip_address", v, nil),
			"private_ip_address":   utils.PathSearch("private_ip_address", v, nil),
			"region":               utils.PathSearch("region", v, nil),
			"vpc_endpoint_id":      utils.PathSearch("vpc_endpoint_id", v, nil),
			"vpc_endpoint_address": utils.PathSearch("vpc_endpoint_address", v, nil),
			"monitor":              flattenComponentRunningNodesMonitor(utils.PathSearch("monitor", v, nil)),
			"node_expansion":       flattenComponentRunningNodesNodeExpansion(utils.PathSearch("node_expansion", v, nil)),
			"node_apply_fail_enum": utils.PathSearch("node_apply_fail_enum", v, nil),
			"list":                 flattenComponentRunningNodesList(utils.PathSearch("list", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenComponentRunningNodesMonitor(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"mini_on_online":  utils.PathSearch("mini_on_online", resp, nil),
			"memory_count":    utils.PathSearch("memory_count", resp, nil),
			"memory_usage":    utils.PathSearch("memory_usage", resp, nil),
			"memory_free":     utils.PathSearch("memory_free", resp, nil),
			"memory_shared":   utils.PathSearch("memory_shared", resp, nil),
			"memory_cache":    utils.PathSearch("memory_cache", resp, nil),
			"cpu_usage":       utils.PathSearch("cpu_usage", resp, nil),
			"cpu_idle":        utils.PathSearch("cpu_idle", resp, nil),
			"up_pps":          utils.PathSearch("up_pps", resp, nil),
			"down_pps":        utils.PathSearch("down_pps", resp, nil),
			"write_rate":      utils.PathSearch("write_rate", resp, nil),
			"read_rate":       utils.PathSearch("read_rate", resp, nil),
			"disk_count":      utils.PathSearch("disk_count", resp, nil),
			"disk_usage":      utils.PathSearch("disk_usage", resp, nil),
			"heart_beat_time": utils.PathSearch("heart_beat_time", resp, nil),
			"health_status":   utils.PathSearch("health_status", resp, nil),
			"heart_beat":      utils.PathSearch("heart_beat", resp, nil),
		},
	}
}

func flattenComponentRunningNodesNodeExpansion(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"node_id":       utils.PathSearch("node_id", resp, nil),
			"data_center":   utils.PathSearch("data_center", resp, nil),
			"custom_label":  utils.PathSearch("custom_label", resp, nil),
			"network_plane": utils.PathSearch("network_plane", resp, nil),
			"description":   utils.PathSearch("description", resp, nil),
			"maintainer":    utils.PathSearch("maintainer", resp, nil),
		},
	}
}

func flattenComponentRunningNodesList(listResp []interface{}) []interface{} {
	if len(listResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(listResp))
	for _, v := range listResp {
		rst = append(rst, map[string]interface{}{
			"configuration_id": utils.PathSearch("configuration_id", v, nil),
			"component_id":     utils.PathSearch("component_id", v, nil),
			"node_id":          utils.PathSearch("node_id", v, nil),
			"file_name":        utils.PathSearch("file_name", v, nil),
			"file_path":        utils.PathSearch("file_path", v, nil),
			"file_type":        utils.PathSearch("file_type", v, nil),
			"param":            utils.PathSearch("param", v, nil),
			"version":          utils.PathSearch("version", v, nil),
			"type":             utils.PathSearch("type", v, nil),
		})
	}

	return rst
}
