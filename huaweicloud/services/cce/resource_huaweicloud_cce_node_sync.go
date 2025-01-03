package cce

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE GET /api/v3/projects/{project_id}/clusters/{cluster_id}
// @API CCE GET /api/v2/projects/{project_id}/clusters/{cluster_id}/nodes/{node_id}/sync
var nodesSyncNonUpdatableParams = []string{"cluster_id", "node_id"}

func ResourceNodeSync() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodeSyncCreate,
		ReadContext:   resourceNodeSyncRead,
		UpdateContext: resourceNodeSyncUpdate,
		DeleteContext: resourceNodeSyncDelete,

		CustomizeDiff: config.FlexibleForceNew(nodesSyncNonUpdatableParams),

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
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceNodeSyncCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		createNodeSyncHttpUrl = "api/v2/projects/{project_id}/clusters/{cluster_id}/nodes/{node_id}/sync"
	)

	createNodeSyncPath := client.Endpoint + createNodeSyncHttpUrl
	createNodeSyncPath = strings.ReplaceAll(createNodeSyncPath, "{project_id}", client.ProjectID)
	createNodeSyncPath = strings.ReplaceAll(createNodeSyncPath, "{cluster_id}", d.Get("cluster_id").(string))
	createNodeSyncPath = strings.ReplaceAll(createNodeSyncPath, "{node_id}", d.Get("node_id").(string))

	createNodeSyncOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("GET", createNodeSyncPath, &createNodeSyncOpt)
	if err != nil {
		return diag.Errorf("error syncing CCE node: %s", err)
	}

	d.SetId(d.Get("node_id").(string))

	return resourceNodeSyncRead(ctx, d, meta)
}

func resourceNodeSyncRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceNodeSyncUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceNodeSyncDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting node sync resource is not supported. The node sync resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
