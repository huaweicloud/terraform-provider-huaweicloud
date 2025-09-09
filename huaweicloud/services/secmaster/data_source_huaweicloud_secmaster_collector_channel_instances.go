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

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/collector/channels/instances
func DataSourceCollectorChannelInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCollectorChannelInstancesRead,

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
			"channel_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the channel ID.",
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
				Description: "Specifies the sort key.",
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the sort direction. Supported values are **asc** and **desc**.",
			},
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The collector channel instances list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The channel name.",
						},
						"config_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The collector channel configuration status.",
						},
						"create_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IAM user ID.",
						},
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The node name.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region.",
						},
						"mini_on_online": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether online.",
						},
						"public_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public IP address.",
						},
						"private_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The private IP address.",
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
						"read_write": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The read write record information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"channel_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The channel ID (UUID).",
									},
									"minion_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The minion ID (UUID).",
									},
									"accept_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The accept count.",
									},
									"send_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The send count.",
									},
									"accept_rate": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The accept rate.",
									},
									"send_rate": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The send rate.",
									},
									"heart_beat_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The time when the last heartbeat signal was received.",
									},
									"latest_transmission_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The time of the last transmission.",
									},
									"channel_instance_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of collector channel instances.",
									},
									"heart_beat": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether the node successfully received heartbeat signals.",
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

func buildCollectorChannelInstancesQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("channel_id"); ok {
		queryParams = fmt.Sprintf("%s&channel_id=%v", queryParams, v)
	}
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

func dataSourceCollectorChannelInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/collector/channels/instances"
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
	requestPath += buildCollectorChannelInstancesQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		requestResp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster collector channel instances: %s", err)
		}

		requestRespBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}

		recordsResp := utils.PathSearch("records", requestRespBody, make([]interface{}, 0)).([]interface{})
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
		d.Set("records", flattenCollectorChannelInstancesRecords(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCollectorChannelInstancesRecords(recordsResp []interface{}) []interface{} {
	if len(recordsResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(recordsResp))
	for _, v := range recordsResp {
		rst = append(rst, map[string]interface{}{
			"channel_name":       utils.PathSearch("channel_name", v, nil),
			"config_status":      utils.PathSearch("config_status", v, nil),
			"create_by":          utils.PathSearch("create_by", v, nil),
			"node_name":          utils.PathSearch("node_name", v, nil),
			"region":             utils.PathSearch("region", v, nil),
			"mini_on_online":     utils.PathSearch("mini_on_online", v, nil),
			"public_ip_address":  utils.PathSearch("public_ip_address", v, nil),
			"private_ip_address": utils.PathSearch("private_ip_address", v, nil),
			"monitor":            flattenCollectorChannelInstancesMonitor(v),
			"read_write":         flattenCollectorChannelInstancesReadWrite(v),
		})
	}

	return rst
}

func flattenCollectorChannelInstancesMonitor(respBody interface{}) []interface{} {
	monitorResp := utils.PathSearch("monitor", respBody, nil)
	if monitorResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"mini_on_online":  utils.PathSearch("mini_on_online", monitorResp, nil),
			"memory_count":    utils.PathSearch("memory_count", monitorResp, nil),
			"memory_usage":    utils.PathSearch("memory_usage", monitorResp, nil),
			"memory_free":     utils.PathSearch("memory_free", monitorResp, nil),
			"memory_shared":   utils.PathSearch("memory_shared", monitorResp, nil),
			"memory_cache":    utils.PathSearch("memory_cache", monitorResp, nil),
			"cpu_usage":       utils.PathSearch("cpu_usage", monitorResp, nil),
			"cpu_idle":        utils.PathSearch("cpu_idle", monitorResp, nil),
			"up_pps":          utils.PathSearch("up_pps", monitorResp, nil),
			"down_pps":        utils.PathSearch("down_pps", monitorResp, nil),
			"write_rate":      utils.PathSearch("write_rate", monitorResp, nil),
			"read_rate":       utils.PathSearch("read_rate", monitorResp, nil),
			"disk_count":      utils.PathSearch("disk_count", monitorResp, nil),
			"disk_usage":      utils.PathSearch("disk_usage", monitorResp, nil),
			"heart_beat_time": utils.PathSearch("heart_beat_time", monitorResp, nil),
			"health_status":   utils.PathSearch("health_status", monitorResp, nil),
			"heart_beat":      utils.PathSearch("heart_beat", monitorResp, nil),
		},
	}
}

func flattenCollectorChannelInstancesReadWrite(respBody interface{}) []interface{} {
	readWriteResp := utils.PathSearch("read_write", respBody, nil)
	if readWriteResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"channel_id":               utils.PathSearch("channel_id", readWriteResp, nil),
			"minion_id":                utils.PathSearch("minion_id", readWriteResp, nil),
			"accept_count":             utils.PathSearch("accept_count", readWriteResp, nil),
			"send_count":               utils.PathSearch("send_count", readWriteResp, nil),
			"accept_rate":              utils.PathSearch("accept_rate", readWriteResp, nil),
			"send_rate":                utils.PathSearch("send_rate", readWriteResp, nil),
			"heart_beat_time":          utils.PathSearch("heart_beat_time", readWriteResp, nil),
			"latest_transmission_time": utils.PathSearch("latest_transmission_time", readWriteResp, nil),
			"channel_instance_count":   utils.PathSearch("channel_instance_count", readWriteResp, nil),
			"heart_beat":               utils.PathSearch("heart_beat", readWriteResp, nil),
		},
	}
}
