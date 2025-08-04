package hss

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var containerKubernetesClusterDaemonsetNonUpdatableParams = []string{
	"cluster_id",
	"cluster_name",
	"auto_upgrade",
	"runtime_info",
	"runtime_info.*.runtime_name",
	"runtime_info.*.runtime_path",
	"enterprise_project_id",
}

// @API HSS POST /v5/{project_id}/container/kubernetes/clusters/{cluster_id}/daemonsets
// @API HSS GET /v5/{project_id}/container/kubernetes/clusters/{cluster_id}/daemonsets
// @API HSS PUT /v5/{project_id}/container/kubernetes/clusters/{cluster_id}/daemonsets
// @API HSS DELETE /v5/{project_id}/container/kubernetes/clusters/{cluster_id}/daemonsets
func ResourceContainerKubernetesClusterDaemonset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContainerKubernetesClusterDaemonsetCreate,
		ReadContext:   resourceContainerKubernetesClusterDaemonsetRead,
		UpdateContext: resourceContainerKubernetesClusterDaemonsetUpdate,
		DeleteContext: resourceContainerKubernetesClusterDaemonsetDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(containerKubernetesClusterDaemonsetNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Current, `cluster_name`, `auto_upgrade`, `runtime_info`, `runtime_info.runtime_name` parameters
			// is optional in the API documentation and is required after confirmation with the HSS service.
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auto_upgrade": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"runtime_info": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"runtime_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"runtime_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"schedule_info": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_selector": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"pod_tolerances": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			// The `agent_version`, `charging_mode`, `cce_protection_type`, and `prefer_packet_cycle`
			// parameters can only be used during update API.
			// Due to API issues, parameters `agent_version`, `invoked_service`, `charging_mode`, `cce_protection_type`,
			// `prefer_packet_cycle` and `schedule_info.pod_tolerances` cannot be tested temporarily.
			"agent_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The `invoked_service` parameter can only be used during update API and delete API.
			"invoked_service": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cce_protection_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"prefer_packet_cycle": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"yaml_content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cluster_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ds_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"desired_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"current_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ready_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"installed_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildContainerKubernetesClusterDaemonsetQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func buildContainerKubernetesClusterDaemonsetRuntimeInfoBodyParams(d *schema.ResourceData) []map[string]interface{} {
	runTimeInfoInput := d.Get("runtime_info").([]interface{})
	runTimeInfoBodyParams := make([]map[string]interface{}, 0, len(runTimeInfoInput))
	for _, v := range runTimeInfoInput {
		runTimeInfoBodyParams = append(runTimeInfoBodyParams, map[string]interface{}{
			"runtime_name": utils.PathSearch("runtime_name", v, nil),
			"runtime_path": utils.ValueIgnoreEmpty(utils.PathSearch("runtime_path", v, nil)),
		})
	}

	return runTimeInfoBodyParams
}

func buildContainerKubernetesClusterDaemonsetScheduleInfoBodyParams(d *schema.ResourceData) map[string]interface{} {
	scheduleInfoInput := d.Get("schedule_info").([]interface{})
	if len(scheduleInfoInput) == 0 {
		return nil
	}

	return map[string]interface{}{
		"node_selector":  utils.ValueIgnoreEmpty(utils.PathSearch("node_selector", scheduleInfoInput[0], nil)),
		"pod_tolerances": utils.ValueIgnoreEmpty(utils.PathSearch("pod_tolerances", scheduleInfoInput[0], nil)),
	}
}

func buildCreateContainerKubernetesClusterDaemonsetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"cluster_name":  d.Get("cluster_name"),
		"auto_upgrade":  d.Get("auto_upgrade"),
		"runtime_info":  buildContainerKubernetesClusterDaemonsetRuntimeInfoBodyParams(d),
		"schedule_info": buildContainerKubernetesClusterDaemonsetScheduleInfoBodyParams(d),
	}

	return bodyParams
}

func resourceContainerKubernetesClusterDaemonsetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		product   = "hss"
		clusterID = d.Get("cluster_id").(string)
		epsId     = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/container/kubernetes/clusters/{cluster_id}/daemonsets"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{cluster_id}", clusterID)
	requestPath += buildContainerKubernetesClusterDaemonsetQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateContainerKubernetesClusterDaemonsetBodyParams(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating HSS container kubernetes cluster daemonset: %s", err)
	}

	d.SetId(clusterID)

	if d.Get("agent_version").(string) != "" || d.Get("invoked_service").(string) != "" ||
		d.Get("charging_mode").(string) != "" || d.Get("cce_protection_type").(string) != "" ||
		d.Get("prefer_packet_cycle").(bool) {
		err = updateContainerKubernetesClusterDaemonset(client, d, epsId, region)
		if err != nil {
			return diag.Errorf(
				"error updating HSS container kubernetes cluster daemonset in creation operation: %s", err)
		}
	}

	return resourceContainerKubernetesClusterDaemonsetRead(ctx, d, meta)
}

func resourceContainerKubernetesClusterDaemonsetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	respBody, err := getClusterDaemonset(client, d.Id(), epsId)
	if err != nil {
		return diag.FromErr(err)
	}

	runtimeInfoResp := utils.PathSearch("runtime_info", respBody, nil)
	yamlContentResp := utils.PathSearch("yaml_content", respBody, nil)

	// The query API always returns `200` status code.
	// So use the `runtime_info` and `yaml_content` fields to determine if resource exist.
	if runtimeInfoResp == nil || yamlContentResp == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("runtime_info", flattenContainerKubernetesClusterDaemonsetRuntimeInfo(runtimeInfoResp)),
		d.Set("schedule_info", flattenContainerKubernetesClusterDaemonsetScheduleInfo(
			utils.PathSearch("schedule_info", respBody, nil))),
		d.Set("yaml_content", yamlContentResp),
		d.Set("node_num", utils.PathSearch("node_num", respBody, nil)),
		d.Set("cluster_status", utils.PathSearch("cluster_status", respBody, nil)),
		d.Set("ds_info", flattenContainerKubernetesClusterDaemonsetDsInfo(
			utils.PathSearch("ds_info", respBody, nil))),
		d.Set("installed_status", utils.PathSearch("installed_status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenContainerKubernetesClusterDaemonsetRuntimeInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := make([]interface{}, len(resp.([]interface{})))
	for i, v := range resp.([]interface{}) {
		result[i] = map[string]interface{}{
			"runtime_name": utils.PathSearch("runtime_name", v, nil),
			"runtime_path": utils.PathSearch("runtime_path", v, nil),
		}
	}

	return result
}

func flattenContainerKubernetesClusterDaemonsetScheduleInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"node_selector":  utils.ExpandToStringList(utils.PathSearch("node_selector", resp, make([]interface{}, 0)).([]interface{})),
			"pod_tolerances": utils.ExpandToStringList(utils.PathSearch("pod_tolerances", resp, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenContainerKubernetesClusterDaemonsetDsInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"desired_num": utils.PathSearch("desired_num", resp, nil),
			"current_num": utils.PathSearch("current_num", resp, nil),
			"ready_num":   utils.PathSearch("ready_num", resp, nil),
		},
	}
}

func buildUpdateContainerKubernetesClusterDaemonsetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"agent_version":       utils.ValueIgnoreEmpty(d.Get("agent_version")),
		"cluster_name":        d.Get("cluster_name"),
		"auto_upgrade":        d.Get("auto_upgrade"),
		"runtime_info":        buildContainerKubernetesClusterDaemonsetRuntimeInfoBodyParams(d),
		"schedule_info":       buildContainerKubernetesClusterDaemonsetScheduleInfoBodyParams(d),
		"invoked_service":     utils.ValueIgnoreEmpty(d.Get("invoked_service")),
		"charging_mode":       utils.ValueIgnoreEmpty(d.Get("charging_mode")),
		"cce_protection_type": utils.ValueIgnoreEmpty(d.Get("cce_protection_type")),
		"prefer_packet_cycle": d.Get("prefer_packet_cycle"),
	}

	return bodyParams
}

func updateContainerKubernetesClusterDaemonset(client *golangsdk.ServiceClient, d *schema.ResourceData, epsId, region string) error {
	requestPath := client.Endpoint + "v5/{project_id}/container/kubernetes/clusters/{cluster_id}/daemonsets"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{cluster_id}", d.Id())
	requestPath += buildContainerKubernetesClusterDaemonsetQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
		JSONBody:         utils.RemoveNil(buildUpdateContainerKubernetesClusterDaemonsetBodyParams(d)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)

	return err
}

func resourceContainerKubernetesClusterDaemonsetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	err = updateContainerKubernetesClusterDaemonset(client, d, epsId, region)
	if err != nil {
		return diag.Errorf("error updating HSS container kubernetes cluster daemonset: %s", err)
	}

	return resourceContainerKubernetesClusterDaemonsetRead(ctx, d, meta)
}

func buildDeleteContainerKubernetesClusterDaemonsetQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := ""

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	if v, ok := d.GetOk("invoked_service"); ok {
		queryParams = fmt.Sprintf("%s&invoked_service=%v", queryParams, v)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func resourceContainerKubernetesClusterDaemonsetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		id      = d.Id()
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	deletePath := client.Endpoint + "v5/{project_id}/container/kubernetes/clusters/{cluster_id}/daemonsets"
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", id)
	deletePath += buildDeleteContainerKubernetesClusterDaemonsetQueryParams(d, epsId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting HSS container kubernetes cluster daemonset: %s", err)
	}

	// Delete API always return `200` status code.
	// So it is necessary to use a query API to determine whether the resource has been successfully deleted.
	if err := waitingForContainerKubernetesClusterDaemonsetDeleted(ctx, client, id, epsId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for HSS container kubernetes cluster daemonset (%s) deleted: %s", id, err)
	}

	return nil
}

func getClusterDaemonset(client *golangsdk.ServiceClient, clusterId, epsId string) (interface{}, error) {
	queryPath := client.Endpoint + "v5/{project_id}/container/kubernetes/clusters/{cluster_id}/daemonsets"
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath = strings.ReplaceAll(queryPath, "{cluster_id}", clusterId)
	queryPath += buildContainerKubernetesClusterDaemonsetQueryParams(epsId)
	queryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", queryPath, &queryOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HSS container kubernetes cluster daemonset: %s", err)
	}

	return utils.FlattenResponse(resp)
}

func waitingForContainerKubernetesClusterDaemonsetDeleted(ctx context.Context, client *golangsdk.ServiceClient, clusterId, epsId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := getClusterDaemonset(client, clusterId, epsId)
			if err != nil {
				return nil, "ERROR", err
			}

			if respBody == nil {
				return "success", "COMPLETED", nil
			}

			runtimeInfoResp := utils.PathSearch("runtime_info", respBody, nil)
			yamlContentResp := utils.PathSearch("yaml_content", respBody, nil)
			if runtimeInfoResp == nil && yamlContentResp == nil {
				return "success", "COMPLETED", nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        3 * time.Second,
		PollInterval: 3 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)

	return err
}
