package secmaster

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/nodes
func ResourceDeleteNodes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeleteNodesCreate,
		ReadContext:   resourceDeleteNodesRead,
		UpdateContext: resourceDeleteNodesUpdate,
		DeleteContext: resourceDeleteNodesDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{"workspace_id", "delete_ids"}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delete_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"delete_fail_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     nodeSchema(),
			},
			"delete_success_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     nodeSchema(),
			},
		},
	}
}

func nodeSchema() *schema.Resource {
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
				Elem:     isapErrorRspSchema(),
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
				Elem:     monitorSchema(),
			},
			"node_expansion": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     isapNodeExpansionSchema(),
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

func isapErrorRspSchema() *schema.Resource {
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

func monitorSchema() *schema.Resource {
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

func isapNodeExpansionSchema() *schema.Resource {
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

func resourceDeleteNodesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/nodes"
		workspaceID   = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + createHttpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody: map[string]interface{}{
			"delete_ids": d.Get("delete_ids"),
		},
	}

	resp, err := client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster nodes: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("delete_fail_list", flattenDeleteNodesList(utils.PathSearch("delete_fail_list", respBody, nil))),
		d.Set("delete_success_list", flattenDeleteNodesList(utils.PathSearch("delete_success_list", respBody, nil))),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting SecMaster delete nodes attributes: %s", mErr)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	return resourceDeleteNodesRead(ctx, d, meta)
}

func flattenDeleteNodesList(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	respArray, ok := respBody.([]interface{})
	if !ok || len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"create_by":            utils.PathSearch("create_by", v, nil),
			"create_time":          utils.PathSearch("create_time", v, nil),
			"description":          flattenIsapErrorRsp(utils.PathSearch("description", v, nil)),
			"device_type":          utils.PathSearch("device_type", v, nil),
			"ip_address":           utils.PathSearch("ip_address", v, nil),
			"monitor":              flattenMonitor(utils.PathSearch("monitor", v, nil)),
			"node_expansion":       flattenIsapNodeExpansion(utils.PathSearch("node_expansion", v, nil)),
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

func flattenIsapErrorRsp(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"error_code": utils.PathSearch("error_code", respBody, nil),
			"error_msg":  utils.PathSearch("error_msg", respBody, nil),
		},
	}
}

func flattenMonitor(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"cpu_idle":        utils.PathSearch("cpu_idle", respBody, nil),
			"cpu_usage":       utils.PathSearch("cpu_usage", respBody, nil),
			"disk_count":      utils.PathSearch("disk_count", respBody, nil),
			"disk_usage":      utils.PathSearch("disk_usage", respBody, nil),
			"down_pps":        utils.PathSearch("down_pps", respBody, nil),
			"health_status":   utils.PathSearch("health_status", respBody, nil),
			"heart_beat":      utils.PathSearch("heart_beat", respBody, nil),
			"heart_beat_time": utils.PathSearch("heart_beat_time", respBody, nil),
			"memory_cache":    utils.PathSearch("memory_cache", respBody, nil),
			"memory_count":    utils.PathSearch("memory_count", respBody, nil),
			"memory_free":     utils.PathSearch("memory_free", respBody, nil),
			"memory_shared":   utils.PathSearch("memory_shared", respBody, nil),
			"memory_usage":    utils.PathSearch("memory_usage", respBody, nil),
			"mini_on_online":  utils.PathSearch("mini_on_online", respBody, nil),
			"read_rate":       utils.PathSearch("read_rate", respBody, nil),
			"up_pps":          utils.PathSearch("up_pps", respBody, nil),
			"write_rate":      utils.PathSearch("write_rate", respBody, nil),
		},
	}
}

func flattenIsapNodeExpansion(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"custom_label":  utils.PathSearch("custom_label", respBody, nil),
			"data_center":   utils.PathSearch("data_center", respBody, nil),
			"description":   utils.PathSearch("description", respBody, nil),
			"maintainer":    utils.PathSearch("maintainer", respBody, nil),
			"network_plane": utils.PathSearch("network_plane", respBody, nil),
			"node_id":       utils.PathSearch("node_id", respBody, nil),
		},
	}
}

func resourceDeleteNodesRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDeleteNodesUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDeleteNodesDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to delete SecMaster nodes. Deleting this 
resource will not change the current SecMaster nodes, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
