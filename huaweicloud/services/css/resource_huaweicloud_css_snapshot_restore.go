package css

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

var snapshotRestoreNonUpdatableParams = []string{
	"source_cluster_id",
	"target_cluster_id",
	"snapshot_id",
	"indices",
	"rename_pattern",
	"rename_replacement",
}

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/{snapshot_id}/restore
// @API CSS GET /v1.0/{project_id}/clusters
func ResourceSnapshotRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSnapshotRestoreCreate,
		ReadContext:   resourceSnapshotRestoreRead,
		UpdateContext: resourceSnapshotRestoreUpdate,
		DeleteContext: resourceSnapshotRestoreDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(snapshotRestoreNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the source cluster ID.`,
			},
			"snapshot_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the snapshot to be restored.`,
			},
			"target_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the target cluster ID.`,
			},
			"indices": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Name of an index to be restored.`,
			},
			"rename_pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Rule for defining the indexes to be restored.`,
			},
			"rename_replacement": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Rule for renaming an index.`,
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

func resourceSnapshotRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	targetClusterId := d.Get("target_cluster_id").(string)
	snapshotId := d.Get("snapshot_id").(string)
	clusterId := d.Get("source_cluster_id").(string)
	restoreSnapshotHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/{snapshot_id}/restore"
	restoreSnapshotPath := client.Endpoint + restoreSnapshotHttpUrl
	restoreSnapshotPath = strings.ReplaceAll(restoreSnapshotPath, "{project_id}", client.ProjectID)
	restoreSnapshotPath = strings.ReplaceAll(restoreSnapshotPath, "{cluster_id}", clusterId)
	restoreSnapshotPath = strings.ReplaceAll(restoreSnapshotPath, "{snapshot_id}", snapshotId)

	restoreSnapshotOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	restoreSnapshotOpt.JSONBody = utils.RemoveNil(buildCreateSnapshotRestoreBodyParams(d))

	_, err = client.Request("POST", restoreSnapshotPath, &restoreSnapshotOpt)
	if err != nil {
		return diag.Errorf("error restore CSS cluster snapshot extend, cluster_id: %s, error: %s", d.Id(), err)
	}

	err = checkClusterOperationResult(ctx, client, targetClusterId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	return resourceSnapshotRestoreRead(ctx, d, meta)
}

func buildCreateSnapshotRestoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"targetCluster":     d.Get("target_cluster_id"),
		"indices":           utils.ValueIgnoreEmpty(d.Get("indices")),
		"renamePattern":     utils.ValueIgnoreEmpty(d.Get("rename_pattern")),
		"renameReplacement": utils.ValueIgnoreEmpty(d.Get("rename_replacement")),
	}
	return bodyParams
}

func resourceSnapshotRestoreRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSnapshotRestoreUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSnapshotRestoreDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting restoration record is not supported. The restoration record is only removed from the state," +
		" but it remains in the cloud. And the instance doesn't return to the state before restoration."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
