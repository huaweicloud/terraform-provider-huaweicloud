// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DWS
// ---------------------------------------------------------------

package dws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	deleteNotExistMsg            = "logical cluster is not existed"
	deleteFirstLogicalClusterMsg = "the first logical cluster can't be deleted"
)

// @API DWS POST /v2/{project_id}/clusters/{cluster_id}/logical-clusters
// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/logical-clusters
// @API DWS DELETE /v2/{project_id}/clusters/{cluster_id}/logical-clusters/{logical_cluster_id}
func ResourceLogicalCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogicalClusterCreate,
		ReadContext:   resourceLogicalClusterRead,
		DeleteContext: resourceLogicalClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLogicalClusterImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the DWS cluster ID.`,
			},
			"logical_cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the logical cluster name.`,
			},
			"cluster_rings": {
				Type:        schema.TypeSet,
				Elem:        logicalClusterRingsSchema(),
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the DWS cluster ring list information.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The DWS logical cluster status.`,
			},
			"first_logical_cluster": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether it is the first logical cluster. The first logical cluster cannot be deleted.`,
			},
			"edit_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether editing is allowed.`,
			},
			"restart_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to allow restart.`,
			},
			"delete_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether deletion is allowed.`,
			},
		},
	}
}

func logicalClusterRingsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"ring_hosts": {
				Type:        schema.TypeSet,
				Elem:        logicalRingHostsSchema(),
				Required:    true,
				ForceNew:    true,
				Description: `Indicates the cluster host ring information.`,
			},
		},
	}
	return &sc
}

func logicalRingHostsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"host_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the host name.`,
			},
			"back_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the backend IP address.`,
			},
			"cpu_cores": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the number of CPU cores.`,
			},
			"memory": {
				Type:        schema.TypeFloat,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the host memory.`,
			},
			"disk_size": {
				Type:        schema.TypeFloat,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the host disk size.`,
			},
		},
	}
	return &sc
}

func resourceLogicalClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v2/{project_id}/clusters/{cluster_id}/logical-clusters"
		product     = "dws"
		clusterName = d.Get("logical_cluster_name").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", d.Get("cluster_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateLogicalClusterBodyParams(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DWS logical cluster: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	errCode := utils.PathSearch("error_code", createRespBody, "").(string)
	if errCode != "DWS.0000" {
		errMsg := utils.PathSearch("error_msg", createRespBody, "").(string)
		return diag.Errorf("error creating DWS logical cluster: error code: %s, error message: %s", errCode, errMsg)
	}

	clusterRespBody, err := waitingForStateCompleted(ctx, client, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for DWS logical cluster (%s) creation to complete: %s", clusterName, err)
	}

	id, err := jmespath.Search("logical_cluster_id", clusterRespBody)
	if err != nil || id == nil {
		return diag.Errorf("error creating DWS logical cluster: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourceLogicalClusterRead(ctx, d, meta)
}

func buildCreateLogicalClusterBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"logical_cluster": map[string]interface{}{
			"logical_cluster_name": d.Get("logical_cluster_name"),
			"cluster_rings":        buildLogicalClusterRingsRequestBody(d.Get("cluster_rings")),
		},
	}
}

func buildLogicalClusterRingsRequestBody(rawParams interface{}) []map[string]interface{} {
	if rawSet, ok := rawParams.(*schema.Set); ok && rawSet.Len() > 0 {
		rst := make([]map[string]interface{}, 0, rawSet.Len())
		for _, v := range rawSet.List() {
			raw, isMap := v.(map[string]interface{})
			if !isMap {
				continue
			}

			rst = append(rst, map[string]interface{}{
				"ring_hosts": buildLogicalRingHostsRequestBody(raw["ring_hosts"]),
			})
		}
		return rst
	}
	return nil
}

func buildLogicalRingHostsRequestBody(rawParams interface{}) []map[string]interface{} {
	if rawSet, ok := rawParams.(*schema.Set); ok && rawSet.Len() > 0 {
		rst := make([]map[string]interface{}, 0, rawSet.Len())
		for _, v := range rawSet.List() {
			raw, isMap := v.(map[string]interface{})
			if !isMap {
				continue
			}

			rst = append(rst, map[string]interface{}{
				"host_name": raw["host_name"],
				"back_ip":   raw["back_ip"],
				"cpu_cores": raw["cpu_cores"],
				"memory":    raw["memory"],
				"disk_size": raw["disk_size"],
			})
		}
		return rst
	}
	return nil
}

func waitingForStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) (interface{}, error) {
	clusterName := d.Get("logical_cluster_name").(string)
	expression := fmt.Sprintf("logical_clusters[?logical_cluster_name=='%s']|[0]", clusterName)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			clusterRespBody, err := readLogicalClusters(client, d)
			if err != nil {
				return nil, "ERROR", err
			}

			cluster := utils.PathSearch(expression, clusterRespBody, nil)
			if cluster == nil {
				return nil, "ERROR", golangsdk.ErrDefault404{}
			}

			completed := utils.PathSearch("action_info.completed", cluster, false).(bool)
			result := utils.PathSearch("action_info.result", cluster, "").(string)
			if completed && result == "success" {
				return cluster, "COMPLETED", nil
			}

			if completed && result == "failed" {
				return cluster, "ERROR", fmt.Errorf("the DWS logical cluster (%s) is failed", clusterName)
			}
			return cluster, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}

func resourceLogicalClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr                     *multierror.Error
		cfg                      = meta.(*config.Config)
		region                   = cfg.GetRegion(d)
		getLogicalClusterProduct = "dws"
	)
	client, err := cfg.NewServiceClient(getLogicalClusterProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	clusterRespBody, err := readLogicalClusters(client, d)
	// The list API response status code is `404` when the cluster does not exist.
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DWS logical cluster")
	}

	expression := fmt.Sprintf("logical_clusters[?logical_cluster_id=='%s']|[0]", d.Id())
	cluster := utils.PathSearch(expression, clusterRespBody, nil)
	if cluster == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("logical_cluster_name", utils.PathSearch("logical_cluster_name", cluster, nil)),
		d.Set("cluster_rings", flattenResponseBodyClusterRings(cluster)),
		d.Set("status", utils.PathSearch("status", cluster, nil)),
		d.Set("first_logical_cluster", utils.PathSearch("first_logical_cluster", cluster, nil)),
		d.Set("edit_enable", utils.PathSearch("edit_enable", cluster, nil)),
		d.Set("restart_enable", utils.PathSearch("restart_enable", cluster, nil)),
		d.Set("delete_enable", utils.PathSearch("delete_enable", cluster, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenResponseBodyClusterRings(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("cluster_rings", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"ring_hosts": flattenRingHosts(v),
		}
	}
	return rst
}

// waitingForDeleteStateEnable This method is used to wait for operable status before deleting.
// Deleting operations can only be performed when `delete_enable` is true.
func waitingForDeleteStateEnable(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	expression := fmt.Sprintf("logical_clusters[?logical_cluster_id=='%s']|[0]", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED", "FIRST_LOGICAL_CLUSTER"},
		Refresh: func() (interface{}, string, error) {
			clusterRespBody, err := readLogicalClusters(client, d)
			if err != nil {
				return nil, "ERROR", err
			}

			cluster := utils.PathSearch(expression, clusterRespBody, nil)
			if cluster == nil {
				return nil, "ERROR", golangsdk.ErrDefault404{}
			}

			enable := utils.PathSearch("delete_enable", cluster, false).(bool)
			if enable {
				return enable, "COMPLETED", nil
			}

			// When `first_logical_cluster` is true, field `delete_enable` will always be false.
			isFirstCluster := utils.PathSearch("first_logical_cluster", cluster, false).(bool)
			if isFirstCluster {
				return enable, "FIRST_LOGICAL_CLUSTER", nil
			}

			return enable, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceLogicalClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/logical-clusters/{logical_cluster_id}"
		product = "dws"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	if err := waitingForDeleteStateEnable(ctx, client, d, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for DWS logical cluster (%s) to become delete enable: %s", d.Id(), err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", d.Get("cluster_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{logical_cluster_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting DWS logical cluster: %s", err)
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// The successful API call does not mean that the deletion is successful. Also need to judge and parse the response body.
	if diagResponse := checkLogicalClusterDeleteResponse(d, deleteRespBody); diagResponse != nil {
		return diagResponse
	}

	err = waitingForStateDeleted(ctx, client, d, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for DWS logical cluster (%s) deletion to complete: %s", d.Id(), err)
	}
	return nil
}

// waitingForStateDeleted This method is used to wait for delete to complete.
func waitingForStateDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	expression := fmt.Sprintf("logical_clusters[?logical_cluster_id=='%s']|[0]", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			clusterRespBody, err := readLogicalClusters(client, d)
			if err != nil {
				return nil, "ERROR", err
			}

			cluster := utils.PathSearch(expression, clusterRespBody, nil)
			if cluster == nil {
				obj := map[string]string{"code": "COMPLETED"}
				return obj, "COMPLETED", nil
			}
			return cluster, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

// checkLogicalClusterDeleteResponse Check whether the DWS logical cluster delete API response body contains error msg.
// Not exist response body:     `{"error_code": "DWS.9999","error_msg": "logical cluster is not existed"}`
// First cluster response body: `{"error_code": "DWS.9999","error_msg": "the first logical cluster can't be deleted"}`
// Success delete response body:`{"error_code": "DWS.0000","error_msg": null}`
func checkLogicalClusterDeleteResponse(d *schema.ResourceData, respBody interface{}) diag.Diagnostics {
	errCode := utils.PathSearch("error_code", respBody, "").(string)
	if errCode == "DWS.0000" {
		return nil
	}

	errMsg := utils.PathSearch("error_msg", respBody, "").(string)
	if errMsg == deleteNotExistMsg {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	if errMsg == deleteFirstLogicalClusterMsg {
		errMsg = "The first logical cluster can't be deleted. The project is only removed from the state," +
			" but it remains in the cloud."
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  errMsg,
			},
		}
	}
	return diag.Errorf("error deleting DWS logical cluster: error code: %s, error message: %s", errCode, errMsg)
}

func resourceLogicalClusterImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	partLength := len(parts)

	if partLength == 2 {
		d.SetId(parts[1])
		return []*schema.ResourceData{d}, d.Set("cluster_id", parts[0])
	}
	return nil, fmt.Errorf("invalid format specified for import ID, must be <cluster_id>/<id>")
}
