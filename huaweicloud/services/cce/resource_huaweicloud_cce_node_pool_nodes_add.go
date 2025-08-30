package cce

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE GET /api/v3/projects/{project_id}/clusters/{cluster_id}
// @API CCE POST /api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools/{nodepool_id}/nodes/add
// @API CCE GET /api/v3/projects/{project_id}/jobs/{job_id}
var nodesAddNonUpdatableParams = []string{"cluster_id", "nodepool_id", "node_list", "node_list.*.server_id"}

func ResourcePoolNodesAdd() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePoolNodesAddCreate,
		ReadContext:   resourcePoolNodesAddRead,
		UpdateContext: resourcePoolNodesAddUpdate,
		DeleteContext: resourcePoolNodesAddDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(nodesAddNonUpdatableParams),

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
			"nodepool_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					}},
			},
			"remove_nodes_on_delete": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildPoolNodesAddCreateOpts(d *schema.ResourceData) map[string]interface{} {
	result := map[string]interface{}{
		"kind":       "List",
		"apiVersion": "v3",
		"nodeList":   buildPoolNodesAddNodeListOpts(d),
	}
	return result
}

func buildPoolNodesAddNodeListOpts(d *schema.ResourceData) []map[string]interface{} {
	nodeListRaw := d.Get("node_list").([]interface{})
	if len(nodeListRaw) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(nodeListRaw))

	for i, v := range nodeListRaw {
		nodeList := v.(map[string]interface{})
		result[i] = map[string]interface{}{
			"serverID": nodeList["server_id"],
		}
	}
	return result
}

func resourcePoolNodesAddCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.CceV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 client: %s", err)
	}

	// Wait for the cce cluster to become available
	clusterID := d.Get("cluster_id").(string)
	stateCluster := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      clusterStateRefreshFunc(client, clusterID, []string{"Available"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateCluster.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for CCE cluster to become available: %s", err)
	}

	var (
		createPoolNodesAddHttpUrl = "api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools/{nodepool_id}/nodes/add"
	)

	createPoolNodesAddPath := client.Endpoint + createPoolNodesAddHttpUrl
	createPoolNodesAddPath = strings.ReplaceAll(createPoolNodesAddPath, "{project_id}", client.ProjectID)
	createPoolNodesAddPath = strings.ReplaceAll(createPoolNodesAddPath, "{cluster_id}", d.Get("cluster_id").(string))
	createPoolNodesAddPath = strings.ReplaceAll(createPoolNodesAddPath, "{nodepool_id}", d.Get("nodepool_id").(string))

	createPoolNodesAddOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createPoolNodesAddOpt.JSONBody = utils.RemoveNil(buildPoolNodesAddCreateOpts(d))
	createPoolNodesAddResp, err := client.Request("POST",
		createPoolNodesAddPath, &createPoolNodesAddOpt)
	if err != nil {
		return diag.Errorf("error adding nodes to CCE node pool: %s", err)
	}

	createPoolNodesAddRespBody, err := utils.FlattenResponse(createPoolNodesAddResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobID := utils.PathSearch("jobid", createPoolNodesAddRespBody, "")
	if jobID == "" {
		return diag.Errorf("error adding nodes to CCE node pool: jobid is not found in API response")
	}

	stateJob := &resource.StateChangeConf{
		Pending:      []string{"Initializing", "Running"},
		Target:       []string{"Success"},
		Refresh:      waitForJobStatus(client, jobID.(string)),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        120 * time.Second,
		PollInterval: 20 * time.Second,
	}

	v, err := stateJob.WaitForStateContext(ctx)
	if err != nil {
		if job, ok := v.(*nodes.Job); ok {
			return diag.Errorf("error waiting for job (%s) to become success: %s, reason: %s",
				jobID, err, job.Status.Reason)
		}

		return diag.Errorf("error waiting for job (%s) to become success: %s", jobID, err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return resourcePoolNodesAddRead(ctx, d, meta)
}

func resourcePoolNodesAddRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePoolNodesAddUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePoolNodesAddDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.Get("remove_nodes_on_delete").(bool) {
		errorMsg := "The nodes add resource is only removed from the state."
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  errorMsg,
			},
		}
	}

	cfg := meta.(*config.Config)
	client, err := cfg.CceV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 client: %s", err)
	}

	// Wait for the cce cluster to become available
	clusterID := d.Get("cluster_id").(string)
	stateCluster := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      clusterStateRefreshFunc(client, clusterID, []string{"Available"}),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateCluster.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for CCE cluster to become available: %s", err)
	}

	var (
		nodesRemoveHttpUrl = "api/v3/projects/{project_id}/clusters/{cluster_id}/nodes/operation/remove"
	)

	nodesRemovePath := client.Endpoint + nodesRemoveHttpUrl
	nodesRemovePath = strings.ReplaceAll(nodesRemovePath, "{project_id}", client.ProjectID)
	nodesRemovePath = strings.ReplaceAll(nodesRemovePath, "{cluster_id}", clusterID)

	nodesRemoveOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	nodesReomveReqBody, err := buildNodesRemoveOpts(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	nodesRemoveOpt.JSONBody = utils.RemoveNil(nodesReomveReqBody)
	nodesRemoveResp, err := client.Request("PUT",
		nodesRemovePath, &nodesRemoveOpt)
	if err != nil {
		return diag.Errorf("error removing nodes: %s", err)
	}

	nodesRemoveRespBody, err := utils.FlattenResponse(nodesRemoveResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobID := utils.PathSearch("status.jobID", nodesRemoveRespBody, "")
	if jobID == "" {
		return diag.Errorf("error removing nodes: status.jobID is not found in API response")
	}

	stateJob := &resource.StateChangeConf{
		Pending:      []string{"Initializing", "Running"},
		Target:       []string{"Success"},
		Refresh:      waitForJobStatus(client, jobID.(string)),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        90 * time.Second,
		PollInterval: 15 * time.Second,
	}

	v, err := stateJob.WaitForStateContext(ctx)
	if err != nil {
		if job, ok := v.(*nodes.Job); ok {
			return diag.Errorf("error waiting for job (%s) to become success: %s, reason: %s",
				jobID, err, job.Status.Reason)
		}

		return diag.Errorf("error waiting for job (%s) to become success: %s", jobID, err)
	}

	return nil
}

func buildNodesRemoveOpts(client *golangsdk.ServiceClient, d *schema.ResourceData) (map[string]interface{}, error) {
	var (
		listNodesHttpUrl = "api/v3/projects/{project_id}/clusters/{cluster_id}/nodes"
		clusterID        = d.Get("cluster_id").(string)
		nodepoolID       = d.Get("nodepool_id").(string)
	)

	listNodesPath := client.Endpoint + listNodesHttpUrl
	listNodesPath = strings.ReplaceAll(listNodesPath, "{project_id}", client.ProjectID)
	listNodesPath = strings.ReplaceAll(listNodesPath, "{cluster_id}", clusterID)

	listNodesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listNodesResp, err := client.Request("GET",
		listNodesPath, &listNodesOpt)
	if err != nil {
		return nil, fmt.Errorf("error getting nodes: %s", err)
	}

	listNodesRespBody, err := utils.FlattenResponse(listNodesResp)
	if err != nil {
		return nil, err
	}

	jpath := fmt.Sprintf(`items[?metadata.annotations."kubernetes.io/node-pool.id"=='%s']`, nodepoolID)
	nodeListRaw := utils.PathSearch(jpath, listNodesRespBody, []interface{}{}).([]interface{})
	if nodeListRaw == nil {
		return nil, errors.New("unable to get nodes")
	}

	nodeList := make([]map[string]interface{}, len(nodeListRaw))
	for i, v := range nodeListRaw {
		nodeList[i] = map[string]interface{}{
			"uid":       utils.PathSearch("metadata.uid", v, nil),
			"server_id": utils.PathSearch("status.serverId", v, nil),
		}
	}

	serverIDs := utils.PathSearch(`[].server_id`, d.Get("node_list"), []interface{}{}).([]interface{})
	if serverIDs == nil {
		return nil, errors.New("unable to get node server IDs")
	}

	res := map[string]interface{}{
		"kind":       "RemoveNodesTask",
		"apiVersion": "v3",
		"spec": map[string]interface{}{
			"nodes": buildNodesRemoveNodeBody(nodeList, serverIDs),
		},
	}

	return res, nil
}

func buildNodesRemoveNodeBody(nodeList []map[string]interface{}, serverIDs []interface{}) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(serverIDs))

	for _, v := range nodeList {
		if utils.SliceContains(serverIDs, v["server_id"]) {
			res = append(res, map[string]interface{}{
				"uid": v["uid"],
			})
		}
	}

	return res
}
