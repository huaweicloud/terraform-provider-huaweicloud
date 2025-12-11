package mrs

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var clusterNodeBatchShrinkNonUpdatableParams = []string{
	"cluster_id",
	"node_group_name",
	"node_count",
	"resource_ids",
}

// @API MRS POST /v2/{project_id}/clusters/{cluster_id}/shrink
// @API MRS GET /v1.1/{project_id}/cluster_infos/{cluster_id}
func ResourceClusterNodeBatchShrink() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterNodeBatchShrinkCreate,
		ReadContext:   resourceClusterNodeBatchShrinkRead,
		UpdateContext: resourceClusterNodeBatchShrinkUpdate,
		DeleteContext: resourceClusterNodeBatchShrinkDelete,

		CustomizeDiff: config.FlexibleForceNew(clusterNodeBatchShrinkNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the MRS cluster is located.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster to be shrunk.`,
			},
			"node_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the node group to which the nodes to be shrunk belong.`,
			},
			"node_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   `The number of nodes to be deleted from the node group.`,
				ConflictsWith: []string{"resource_ids"},
			},
			"resource_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The ID list of resource nodes to be deleted.`,
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildClusterNodeBatchShrinkBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"node_group_name": d.Get("node_group_name").(string),
		"count":           utils.ValueIgnoreEmpty(d.Get("node_count")),
		"resource_ids":    utils.ValueIgnoreEmpty(d.Get("resource_ids")),
	}
}

func resourceClusterNodeBatchShrinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v2/{project_id}/clusters/{cluster_id}/shrink"
		clusterId = d.Get("cluster_id").(string)
	)

	// Lock the resource to prevent concurrent.
	config.MutexKV.Lock(clusterId)
	defer config.MutexKV.Unlock(clusterId)

	client, err := cfg.NewServiceClient("mrs", region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	shrinkPath := client.Endpoint + httpUrl
	shrinkPath = strings.ReplaceAll(shrinkPath, "{project_id}", client.ProjectID)
	shrinkPath = strings.ReplaceAll(shrinkPath, "{cluster_id}", clusterId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
		JSONBody: utils.RemoveNil(buildClusterNodeBatchShrinkBodyParams(d)),
	}

	_, err = client.Request("POST", shrinkPath, &createOpt)
	if err != nil {
		return diag.Errorf("error shrinking nodes for the cluster (%s): %s", clusterId, err)
	}

	err = waitForClusterStatusCompleted(ctx, client, clusterId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the cluster shrink to complete: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate resource ID: %s", err)
	}
	d.SetId(randUUID)

	return nil
}

func resourceClusterNodeBatchShrinkRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterNodeBatchShrinkUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterNodeBatchShrinkDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for batch shrinking the nodes of the cluster. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
