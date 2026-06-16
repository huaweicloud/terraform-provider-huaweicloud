package secmaster

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/collector/nodes
func DataSourceCollectorNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCollectorNodesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"health_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address": {
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
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     collectorNodesRecordsSchema(),
			},
		},
	}
}

func collectorNodesRecordsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"channel_instance_refer_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"monitor": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     collectorNodesMonitorSchema(),
			},
			"node_expansion": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     collectorNodesExpansionSchema(),
			},
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"specification": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vpc_endpoint_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_endpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func collectorNodesMonitorSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cpu_idle": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_usage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_count": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_usage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"down_pps": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"heart_beat": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"heart_beat_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory_cache": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memory_count": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memory_free": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memory_shared": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memory_usage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mini_on_online": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"read_rate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"up_pps": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"write_rate": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func collectorNodesExpansionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"custom_label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_center": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintainer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_plane": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCollectorNodesQueryParams(d *schema.ResourceData, offset int) string {
	rst := fmt.Sprintf("?limit=500&offset=%d", offset)

	if v, ok := d.GetOk("health_status"); ok {
		rst += fmt.Sprintf("&health_status=%v", v)
	}

	if v, ok := d.GetOk("node_id"); ok {
		rst += fmt.Sprintf("&node_id=%v", v)
	}

	if v, ok := d.GetOk("node_name"); ok {
		rst += fmt.Sprintf("&node_name=%v", v)
	}

	if v, ok := d.GetOk("ip_address"); ok {
		rst += fmt.Sprintf("&ip_address=%v", v)
	}

	if v, ok := d.GetOk("sort_key"); ok {
		rst += fmt.Sprintf("&sort_key=%v", v)
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%v", v)
	}

	return rst
}

func dataSourceCollectorNodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v1/{project_id}/workspaces/{workspace_id}/collector/nodes"
		product    = "secmaster"
		offset     = 0
		allRecords = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := requestPath + buildCollectorNodesQueryParams(d, offset)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster collector nodes: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		recordsResp := utils.PathSearch("records", respBody, make([]interface{}, 0)).([]interface{})
		if len(recordsResp) == 0 {
			break
		}

		allRecords = append(allRecords, recordsResp...)
		offset += len(recordsResp)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenCollectorNodesRecords(allRecords)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCollectorNodesRecords(nodes []interface{}) []interface{} {
	if len(nodes) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(nodes))
	for _, v := range nodes {
		rst = append(rst, map[string]interface{}{
			"channel_instance_refer_count": utils.PathSearch(
				"channel_instance_refer_count", v, nil),
			"create_by":    utils.PathSearch("create_by", v, nil),
			"custom_label": utils.PathSearch("custom_label", v, nil),
			"description":  utils.PathSearch("description", v, nil),
			"device_type":  utils.PathSearch("device_type", v, nil),
			"ip_address":   utils.PathSearch("ip_address", v, nil),
			"monitor": flattenRecordsMonitor(
				utils.PathSearch("monitor", v, nil)),
			"node_expansion": flattenRecordsExpansion(
				utils.PathSearch("node_expansion", v, nil)),
			"node_id":              utils.PathSearch("node_id", v, nil),
			"node_name":            utils.PathSearch("node_name", v, nil),
			"os_type":              utils.PathSearch("os_type", v, nil),
			"private_ip_address":   utils.PathSearch("private_ip_address", v, nil),
			"project_id":           utils.PathSearch("project_id", v, nil),
			"region":               utils.PathSearch("region", v, nil),
			"specification":        utils.PathSearch("specification", v, nil),
			"update_time":          utils.PathSearch("update_time", v, nil),
			"vpc_endpoint_address": utils.PathSearch("vpc_endpoint_address", v, nil),
			"vpc_endpoint_id":      utils.PathSearch("vpc_endpoint_id", v, nil),
			"vpc_id":               utils.PathSearch("vpc_id", v, nil),
			"workspace_id":         utils.PathSearch("workspace_id", v, nil),
		})
	}

	return rst
}

func flattenRecordsMonitor(raw interface{}) []interface{} {
	if raw == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"cpu_idle":        utils.PathSearch("cpu_idle", raw, nil),
			"cpu_usage":       utils.PathSearch("cpu_usage", raw, nil),
			"disk_count":      utils.PathSearch("disk_count", raw, nil),
			"disk_usage":      utils.PathSearch("disk_usage", raw, nil),
			"down_pps":        utils.PathSearch("down_pps", raw, nil),
			"health_status":   utils.PathSearch("health_status", raw, nil),
			"heart_beat":      utils.PathSearch("heart_beat", raw, nil),
			"heart_beat_time": utils.PathSearch("heart_beat_time", raw, nil),
			"memory_cache":    utils.PathSearch("memory_cache", raw, nil),
			"memory_count":    utils.PathSearch("memory_count", raw, nil),
			"memory_free":     utils.PathSearch("memory_free", raw, nil),
			"memory_shared":   utils.PathSearch("memory_shared", raw, nil),
			"memory_usage":    utils.PathSearch("memory_usage", raw, nil),
			"mini_on_online":  utils.PathSearch("mini_on_online", raw, nil),
			"read_rate":       utils.PathSearch("read_rate", raw, nil),
			"up_pps":          utils.PathSearch("up_pps", raw, nil),
			"write_rate":      utils.PathSearch("write_rate", raw, nil),
		},
	}
}

func flattenRecordsExpansion(raw interface{}) []interface{} {
	if raw == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"custom_label":  utils.PathSearch("custom_label", raw, nil),
			"data_center":   utils.PathSearch("data_center", raw, nil),
			"description":   utils.PathSearch("description", raw, nil),
			"maintainer":    utils.PathSearch("maintainer", raw, nil),
			"network_plane": utils.PathSearch("network_plane", raw, nil),
			"node_id":       utils.PathSearch("node_id", raw, nil),
		},
	}
}
