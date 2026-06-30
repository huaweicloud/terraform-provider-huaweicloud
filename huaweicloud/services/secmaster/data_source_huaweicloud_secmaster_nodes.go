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

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/nodes
func DataSourceNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNodesRead,

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
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_name": {
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
				Elem:     nodesRecordsSchema(),
			},
		},
	}
}

func nodesRecordsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"create_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     nodesDescriptionSchema(),
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
				Elem:     nodesMonitorSchema(),
			},
			"node_expansion": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     nodesNodeExpansionSchema(),
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
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"specification": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
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
			"vpcep_service_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func nodesDescriptionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"error_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func nodesMonitorSchema() *schema.Resource {
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

func nodesNodeExpansionSchema() *schema.Resource {
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

func buildNodesQueryParams(d *schema.ResourceData, offset int) string {
	rst := ""

	if v, ok := d.GetOk("node_id"); ok {
		rst += fmt.Sprintf("&node_id=%v", v)
	}

	if v, ok := d.GetOk("node_name"); ok {
		rst += fmt.Sprintf("&node_name=%v", v)
	}

	if v, ok := d.GetOk("sort_key"); ok {
		rst += fmt.Sprintf("&sort_key=%v", v)
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%v", v)
	}

	if offset > 0 {
		rst += fmt.Sprintf("&offset=%d", offset)
	}

	if rst != "" {
		rst = "?" + rst[1:]
	}

	return rst
}

func dataSourceNodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v1/{project_id}/workspaces/{workspace_id}/nodes"
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
		currentPath := requestPath + buildNodesQueryParams(d, offset)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster nodes: %s", err)
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
		d.Set("records", flattenNodesRecords(allRecords)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenNodesRecords(nodes []interface{}) []interface{} {
	if len(nodes) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(nodes))
	for _, v := range nodes {
		rst = append(rst, map[string]interface{}{
			"create_by":            utils.PathSearch("create_by", v, nil),
			"create_time":          utils.PathSearch("create_time", v, nil),
			"description":          flattenNodesDescription(utils.PathSearch("description", v, nil)),
			"device_type":          utils.PathSearch("device_type", v, nil),
			"ip_address":           utils.PathSearch("ip_address", v, nil),
			"monitor":              flattenNodesMonitor(utils.PathSearch("monitor", v, nil)),
			"node_expansion":       flattenNodesNodeExpansion(utils.PathSearch("node_expansion", v, nil)),
			"node_id":              utils.PathSearch("node_id", v, nil),
			"node_name":            utils.PathSearch("node_name", v, nil),
			"os_type":              utils.PathSearch("os_type", v, nil),
			"private_ip_address":   utils.PathSearch("private_ip_address", v, nil),
			"region":               utils.PathSearch("region", v, nil),
			"specification":        utils.PathSearch("specification", v, nil),
			"subnet_id":            utils.PathSearch("subnet_id", v, nil),
			"update_time":          utils.PathSearch("update_time", v, nil),
			"vpc_endpoint_address": utils.PathSearch("vpc_endpoint_address", v, nil),
			"vpc_endpoint_id":      utils.PathSearch("vpc_endpoint_id", v, nil),
			"vpc_id":               utils.PathSearch("vpc_id", v, nil),
			"vpcep_service_ip":     utils.PathSearch("vpcep_service_ip", v, nil),
		})
	}

	return rst
}

func flattenNodesDescription(dataResp interface{}) []interface{} {
	if dataResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"error_code": utils.PathSearch("error_code", dataResp, nil),
			"error_msg":  utils.PathSearch("error_msg", dataResp, nil),
		},
	}
}

func flattenNodesMonitor(dataResp interface{}) []interface{} {
	if dataResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"cpu_idle":        utils.PathSearch("cpu_idle", dataResp, nil),
			"cpu_usage":       utils.PathSearch("cpu_usage", dataResp, nil),
			"disk_count":      utils.PathSearch("disk_count", dataResp, nil),
			"disk_usage":      utils.PathSearch("disk_usage", dataResp, nil),
			"down_pps":        utils.PathSearch("down_pps", dataResp, nil),
			"health_status":   utils.PathSearch("health_status", dataResp, nil),
			"heart_beat":      utils.PathSearch("heart_beat", dataResp, nil),
			"heart_beat_time": utils.PathSearch("heart_beat_time", dataResp, nil),
			"memory_cache":    utils.PathSearch("memory_cache", dataResp, nil),
			"memory_count":    utils.PathSearch("memory_count", dataResp, nil),
			"memory_free":     utils.PathSearch("memory_free", dataResp, nil),
			"memory_shared":   utils.PathSearch("memory_shared", dataResp, nil),
			"memory_usage":    utils.PathSearch("memory_usage", dataResp, nil),
			"mini_on_online":  utils.PathSearch("mini_on_online", dataResp, nil),
			"read_rate":       utils.PathSearch("read_rate", dataResp, nil),
			"up_pps":          utils.PathSearch("up_pps", dataResp, nil),
			"write_rate":      utils.PathSearch("write_rate", dataResp, nil),
		},
	}
}

func flattenNodesNodeExpansion(dataResp interface{}) []interface{} {
	if dataResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"custom_label":  utils.PathSearch("custom_label", dataResp, nil),
			"data_center":   utils.PathSearch("data_center", dataResp, nil),
			"description":   utils.PathSearch("description", dataResp, nil),
			"maintainer":    utils.PathSearch("maintainer", dataResp, nil),
			"network_plane": utils.PathSearch("network_plane", dataResp, nil),
			"node_id":       utils.PathSearch("node_id", dataResp, nil),
		},
	}
}
