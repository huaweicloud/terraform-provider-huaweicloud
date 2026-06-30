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

var v2NodeBatchMigrateNonUpdatableParams = []string{
	"source_pool_id",
	"node_names",
	"source_cluster_id",
	"target_cluster_id",
	"target_pool_id",
	"resource_spec",
	"resource_spec.*.flavor",
}

// @API ModelArts POST /v2/{project_id}/pools/{pool_name}/nodes/batch-migrate
// @API ModelArts GET /v2/{project_id}/jobs/{job_id}
func ResourceV2NodeBatchMigrate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2NodeBatchMigrateCreate,
		ReadContext:   resourceV2NodeBatchMigrateRead,
		UpdateContext: resourceV2NodeBatchMigrateUpdate,
		DeleteContext: resourceV2NodeBatchMigrateDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(v2NodeBatchMigrateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the resource pool nodes are located.`,
			},

			// Required parameters.
			"source_pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the source resource pool to which the resource nodes belong.`,
			},
			"source_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the source cluster from which nodes are migrated.`,
			},
			"target_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the target cluster to which nodes are migrated.`,
			},

			// Optional parameters.
			"node_names": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The name list of nodes to be migrated.`,
			},
			"target_pool_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the target resource pool.`,
			},
			"resource_spec": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The configuration information of the migrated nodes in the target resource pool.`,
				Elem:        nodeBatchMigrateResourceSpecSchema(),
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

func nodeBatchMigrateResourceSpecSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource specification name.`,
			},
		},
	}
}

func buildNodeBatchMigrateResourceSpec(resourceSpecs []interface{}) map[string]interface{} {
	if len(resourceSpecs) < 1 {
		return nil
	}

	spec := resourceSpecs[0]
	return map[string]interface{}{
		"flavor": utils.PathSearch("flavor", spec, nil),
	}
}

func buildNodeBatchMigrateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"migrateNodeNames": utils.ValueIgnoreEmpty(d.Get("node_names")),
		"fromClusterName":  d.Get("source_cluster_id"),
		"toClusterName":    d.Get("target_cluster_id"),
		"toPoolName":       utils.ValueIgnoreEmpty(d.Get("target_pool_id")),
		"resourceSpec":     buildNodeBatchMigrateResourceSpec(d.Get("resource_spec").([]interface{})),
	}
}

func migrateV2ResourcePoolNodes(ctx context.Context, client *golangsdk.ServiceClient, poolId string,
	bodyParams map[string]interface{}, timeout time.Duration) error {
	var (
		httpUrl = "v2/{project_id}/pools/{pool_name}/nodes/batch-migrate"
	)

	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{pool_name}", poolId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: bodyParams,
	}

	requestResp, err := client.Request("POST", actionPath, &opt)
	if err != nil {
		return fmt.Errorf("error executing batch migrate operation for the resource pool (%s): %s", poolId, err)
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
		return fmt.Errorf("error waiting for the job status of resource pool (%s) batch migration to complete: %s",
			poolId, err)
	}
	return nil
}

func resourceV2NodeBatchMigrateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		poolId = d.Get("source_pool_id").(string)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = migrateV2ResourcePoolNodes(ctx, client, poolId,
		buildNodeBatchMigrateBodyParams(d), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	return resourceV2NodeBatchMigrateRead(ctx, d, meta)
}

func resourceV2NodeBatchMigrateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2NodeBatchMigrateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2NodeBatchMigrateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for batch migrating the ModelArts nodes. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate
file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
