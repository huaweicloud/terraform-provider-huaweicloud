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
	ccenodes "github.com/chnsz/golangsdk/openstack/cce/v3/nodes"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE PUT /api/v3/projects/{project_id}/clusters/{cluster_id}/nodes/operation/remove
var NodesRemoveNonUpdatableParams = []string{"cluster_id", "nodes", "nodes.*.uid"}

func ResourceNodesRemove() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodesRemoveCreate,
		ReadContext:   resourceNodesRemoveRead,
		UpdateContext: resourceNodesRemoveUpdate,
		DeleteContext: resourceNodesRemoveDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(NodesRemoveNonUpdatableParams),

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
			"nodes": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
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

func buildNodesRemoveCreateOpts(d *schema.ResourceData) map[string]interface{} {
	result := map[string]interface{}{
		"kind":       "RemoveNodesTask",
		"apiVersion": "v3",
		"spec":       buildNodesRemoveSpec(d),
	}
	return result
}

func buildNodesRemoveSpec(d *schema.ResourceData) map[string]interface{} {
	result := map[string]interface{}{
		"nodes": buildNodesRemoveSpecNodes(d),
	}
	return result
}

func buildNodesRemoveSpecNodes(d *schema.ResourceData) []map[string]interface{} {
	nodesRaw := d.Get("nodes").([]interface{})
	if len(nodesRaw) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(nodesRaw))

	for i, v := range nodesRaw {
		nodes := v.(map[string]interface{})
		result[i] = map[string]interface{}{
			"uid": nodes["uid"],
		}
	}
	return result
}

func resourceNodesRemoveCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		createNodesRemoveHttpUrl = "api/v3/projects/{project_id}/clusters/{cluster_id}/nodes/operation/remove"
	)

	createNodesRemovePath := client.Endpoint + createNodesRemoveHttpUrl
	createNodesRemovePath = strings.ReplaceAll(createNodesRemovePath, "{project_id}", client.ProjectID)
	createNodesRemovePath = strings.ReplaceAll(createNodesRemovePath, "{cluster_id}", d.Get("cluster_id").(string))

	createNodesRemoveOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createNodesRemoveOpt.JSONBody = utils.RemoveNil(buildNodesRemoveCreateOpts(d))
	createNodesRemoveResp, err := client.Request("PUT",
		createNodesRemovePath, &createNodesRemoveOpt)
	if err != nil {
		return diag.Errorf("error removing nodes: %s", err)
	}

	createNodesRemoveRespBody, err := utils.FlattenResponse(createNodesRemoveResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobID := utils.PathSearch("status.jobID", createNodesRemoveRespBody, "")
	if jobID == "" {
		return diag.Errorf("error removing nodes: status.jobID is not found in API response")
	}

	stateJob := &resource.StateChangeConf{
		Pending:      []string{"Initializing", "Running"},
		Target:       []string{"Success"},
		Refresh:      waitForJobStatus(client, jobID.(string)),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        90 * time.Second,
		PollInterval: 15 * time.Second,
	}

	v, err := stateJob.WaitForStateContext(ctx)
	if err != nil {
		if job, ok := v.(*ccenodes.Job); ok {
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

	return nil
}

func resourceNodesRemoveRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceNodesRemoveUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceNodesRemoveDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting nodes remove resource is not supported. The nodes add resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
