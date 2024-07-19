package cce

import (
	"context"
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

func resourcePoolNodesAddDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting nodes add resource is not supported. The nodes add resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
