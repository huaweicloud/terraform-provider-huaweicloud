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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var clusterNodeBatchExpandNonUpdatableParams = []string{
	"cluster_id",
	"node_group_name",
	"node_count",
	"skip_bootstrap_scripts",
	"scale_without_start",
}

// @API MRS POST /v2/{project_id}/clusters/{cluster_id}/expand
// @API MRS GET /v1.1/{project_id}/cluster_infos/{cluster_id}
func ResourceClusterNodeBatchExpand() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterNodeBatchExpandCreate,
		ReadContext:   resourceClusterNodeBatchExpandRead,
		UpdateContext: resourceClusterNodeBatchExpandUpdate,
		DeleteContext: resourceClusterNodeBatchExpandDelete,

		CustomizeDiff: config.FlexibleForceNew(clusterNodeBatchExpandNonUpdatableParams),

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
				Description: `The ID of the cluster to which the nodes to be expanded belong.`,
			},
			"node_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the node group to which the nodes to be expanded belong.`,
			},
			"node_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The number of nodes to be expanded.`,
			},
			"skip_bootstrap_scripts": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to skip bootstrap scripts.`,
			},
			"scale_without_start": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to start the components on the node after it has been expanded.`,
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

func buildClusterNodeBatchExpandBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"node_group_name":     d.Get("node_group_name"),
		"count":               d.Get("node_count"),
		"scale_without_start": d.Get("scale_without_start"),
	}

	skipBootstrapScripts := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "skip_bootstrap_scripts")
	if skipBootstrapScripts != nil {
		bodyParams["skip_bootstrap_scripts"] = skipBootstrapScripts.(bool)
	}

	return bodyParams
}

func resourceClusterNodeBatchExpandCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v2/{project_id}/clusters/{cluster_id}/expand"
		clusterId = d.Get("cluster_id").(string)
	)

	// Lock the resource to prevent concurrent.
	config.MutexKV.Lock(clusterId)
	defer config.MutexKV.Unlock(clusterId)

	client, err := cfg.NewServiceClient("mrs", region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	expandPath := client.Endpoint + httpUrl
	expandPath = strings.ReplaceAll(expandPath, "{project_id}", client.ProjectID)
	expandPath = strings.ReplaceAll(expandPath, "{cluster_id}", clusterId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
		JSONBody: buildClusterNodeBatchExpandBodyParams(d),
	}

	resp, err := client.Request("POST", expandPath, &createOpt)
	if err != nil {
		return diag.Errorf("error expanding nodes for the cluster (%s): %s", clusterId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate resource ID: %s", err)
	}
	d.SetId(randUUID)

	orderId := utils.PathSearch("order_id", respBody, "").(string)
	if orderId != "" {
		bssClient, err := cfg.NewServiceClient("bssv2", cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS client: %s", err)
		}

		// The expand nodes interface does not support auto-payment, so the CBC side interface is called to pay for the order.
		if err := cbc.PaySubscriptionOrder(bssClient, orderId); err != nil {
			return diag.Errorf("error paying for expansion order (%s) of cluster (%s): %s", orderId, clusterId, err)
		}

		if err := common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.FromErr(err)
		}
	}

	// The order for scaling nodes is not a primary resource, therefore the `common.WaitOrderResourceComplete` method cannot be used.
	err = waitForClusterStatusCompleted(ctx, client, clusterId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the cluster expansion to complete: %s", err)
	}

	return nil
}

func resourceClusterNodeBatchExpandRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterNodeBatchExpandUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterNodeBatchExpandDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for batch expanding the nodes of the cluster. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
