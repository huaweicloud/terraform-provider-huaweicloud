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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE GET /api/v3/projects/{project_id}/clusters/{cluster_id}
// @API CCE POST /api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools/{nodepool_id}/operation/scale
// @API CCE GET /api/v3/projects/{project_id}/jobs/{job_id}
var nodePoolScaleNonUpdatableParams = []string{"cluster_id", "nodepool_id", "desired_node_count", "scale_groups", "scalable_checking"}

func ResourceNodePoolScale() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodePoolScaleCreate,
		ReadContext:   resourceNodePoolScaleRead,
		UpdateContext: resourceNodePoolScaleUpdate,
		DeleteContext: resourceNodePoolScaleDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(nodePoolScaleNonUpdatableParams),

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
			"desired_node_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"scale_groups": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"scalable_checking": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// charge info: charging_mode, period_unit, period, auto_renew
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenew(nil),

			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildNodePoolScaleCreateOpts(d *schema.ResourceData) map[string]interface{} {
	result := map[string]interface{}{
		"kind":       "NodePool",
		"apiVersion": "v3",
		"spec":       buildNodePoolScaleSpecOpts(d),
	}
	return result
}

func buildNodePoolScaleSpecOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"desiredNodeCount": d.Get("desired_node_count"),
		"scaleGroups":      d.Get("scale_groups"),
		"options":          buildNodePoolScaleOptionsOpts(d),
	}

	return bodyParams
}

func buildNodePoolScaleOptionsOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"scalableChecking":      utils.ValueIgnoreEmpty(d.Get("scalable_checking")),
		"billingConfigOverride": buildNodePoolScaleBillingConfigOverrideOpts(d),
	}

	return bodyParams
}

func buildNodePoolScaleBillingConfigOverrideOpts(d *schema.ResourceData) map[string]interface{} {
	var bodyParams map[string]interface{}

	if d.Get("charging_mode").(string) == "prePaid" {
		bodyParams = map[string]interface{}{
			"billingMode": 1,
			"extendParam": buildNodePoolScaleExtendParamOpts(d),
		}
	}

	if d.Get("charging_mode").(string) == "postPaid" {
		bodyParams = map[string]interface{}{
			"billingMode": 0,
		}
	}

	return bodyParams
}

func buildNodePoolScaleExtendParamOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"periodNum":   d.Get("period"),
		"periodType":  d.Get("period_unit"),
		"isAutoRenew": utils.StringToBool(d.Get("auto_renew")),
		"isAutoPay":   true,
	}

	return bodyParams
}

func resourceNodePoolScaleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.CceV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 client: %s", err)
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		if err := common.ValidatePrePaidChargeInfo(d); err != nil {
			return diag.FromErr(err)
		}
	}

	clusterID := d.Get("cluster_id").(string)
	nodePoolID := d.Get("nodepool_id").(string)

	// Wait for the cce cluster to become available

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
		createNodePoolScaleHttpUrl = "api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools/{nodepool_id}/operation/scale"
	)

	createNodePoolScalePath := client.Endpoint + createNodePoolScaleHttpUrl
	createNodePoolScalePath = strings.ReplaceAll(createNodePoolScalePath, "{project_id}", client.ProjectID)
	createNodePoolScalePath = strings.ReplaceAll(createNodePoolScalePath, "{cluster_id}", clusterID)
	createNodePoolScalePath = strings.ReplaceAll(createNodePoolScalePath, "{nodepool_id}", nodePoolID)

	createNodePoolScaleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createNodePoolScaleOpt.JSONBody = utils.RemoveNil(buildNodePoolScaleCreateOpts(d))
	_, err = client.Request("POST", createNodePoolScalePath, &createNodePoolScaleOpt)
	if err != nil {
		return diag.Errorf("error scaling CCE node pool: %s", err)
	}

	d.SetId(nodePoolID)

	stateConf := &resource.StateChangeConf{
		// The statuses of pending phase includes "Synchronizing" and "Synchronized".
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      nodePoolStateRefreshFunc(client, clusterID, nodePoolID, []string{""}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        90 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for CCE node pool to become available: %s", err)
	}

	return resourceNodePoolScaleRead(ctx, d, meta)
}

func resourceNodePoolScaleRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceNodePoolScaleUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceNodePoolScaleDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting node pool scale resource is not supported. The node pool scale resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
