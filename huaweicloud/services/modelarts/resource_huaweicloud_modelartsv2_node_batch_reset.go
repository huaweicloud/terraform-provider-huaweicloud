package modelarts

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v2NodeBatchResetNonUpdatableParams = []string{
	"pool_id",
	"node_names",
	"rolling_config",
}

// @API ModelArts POST /v2/{project_id}/pools/{pool_name}/nodes/batch-reset
// @API ModelArts GET /v2/{project_id}/jobs/{job_id}
func ResourceV2NodeBatchReset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2NodeBatchResetCreate,
		ReadContext:   resourceV2NodeBatchResetRead,
		UpdateContext: resourceV2NodeBatchResetUpdate,
		DeleteContext: resourceV2NodeBatchResetDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(v2NodeBatchResetNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the resource nodes are located.`,
			},

			// Required parameters.
			"pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource pool name to which the resource nodes belong.`,
			},
			"node_names": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The name list of resource nodes to be reset.`,
			},
			"rolling_config": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `The rolling configuration for the node reset operation.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"strategy": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The rolling strategy for the reset operation.`,
						},
						"max_unavailable": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The maximum number or percentage of nodes that can be reset simultaneously.`,
						},
					},
				},
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					}),
			},
		},
	}
}

func buildResetRollingConfig(rollingConfigs []interface{}) map[string]interface{} {
	if len(rollingConfigs) < 1 {
		return nil
	}

	rollingConfig := rollingConfigs[0]
	return map[string]interface{}{
		"strategy":       utils.PathSearch("strategy", rollingConfig, nil),
		"maxUnavailable": utils.PathSearch("max_unavailable", rollingConfig, nil),
	}
}

func resetV2ResourcePoolNodes(ctx context.Context, client *golangsdk.ServiceClient, poolId string,
	nodeNames []interface{}, rollingConfig map[string]interface{}, timeout time.Duration) error {
	httpUrl := "v2/{project_id}/pools/{pool_name}/nodes/batch-reset"
	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{pool_name}", poolId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"nodeNames":     nodeNames,
			"rollingConfig": rollingConfig,
		},
	}

	requestResp, err := client.Request("POST", actionPath, &opt)
	if err != nil {
		return fmt.Errorf("error executing batch reset operation for the specified nodes (%v): %s", nodeNames, err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return err
	}
	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("unable to find job ID under the resource pool (%s) in API response", poolId)
	}
	err = waitForV2JobCompleted(ctx, client, jobId, timeout)
	if err != nil {
		return fmt.Errorf("error waiting for the job status of resource pool (%s) reset to complete: %s",
			poolId, err)
	}
	return nil
}

func resourceV2NodeBatchResetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		poolId        = d.Get("pool_id").(string)
		nodeNames     = d.Get("node_names").([]interface{})
		rollingConfig = buildResetRollingConfig(d.Get("rolling_config").([]interface{}))
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = resetV2ResourcePoolNodes(ctx, client, poolId, nodeNames, rollingConfig,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	return resourceV2NodeBatchResetRead(ctx, d, meta)
}

func resourceV2NodeBatchResetRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2NodeBatchResetUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2NodeBatchResetDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for batch reset the ModelArts nodes. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate
file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
