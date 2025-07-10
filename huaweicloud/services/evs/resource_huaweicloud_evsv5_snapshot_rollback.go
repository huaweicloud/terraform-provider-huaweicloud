package evs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var snapshotRollBackNonUpdatableParams = []string{
	"snapshot_id",
	"volume_ids",
	"name",
}

// @API EVS POST /v5/{project_id}/cloudsnapshots/{snapshot_id}/rollback
func ResourceV5SnapshotRollBack() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5SnapshotRollBackCreate,
		ReadContext:   resourceV5SnapshotRollBackRead,
		UpdateContext: resourceV5SnapshotRollBackUpdate,
		DeleteContext: resourceV5SnapshotRollBackDelete,

		CustomizeDiff: config.FlexibleForceNew(snapshotRollBackNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceV5SnapshotRollBackCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                     = meta.(*config.Config)
		region                  = cfg.GetRegion(d)
		rollbackSnapshotHttpUrl = "v5/{project_id}/snapshots/{snapshot_id}/rollback"
		rollbackSnapshotProduct = "evs"
	)
	rollbackSnapshotClient, err := cfg.NewServiceClient(rollbackSnapshotProduct, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	rollbackSnapshotPath := rollbackSnapshotClient.Endpoint + rollbackSnapshotHttpUrl
	rollbackSnapshotPath = strings.ReplaceAll(rollbackSnapshotPath, "{project_id}",
		rollbackSnapshotClient.ProjectID)
	rollbackSnapshotPath = strings.ReplaceAll(rollbackSnapshotPath, "{snapshot_id}",
		d.Get("snapshot_id").(string))

	rollbackSnapshotBodyParams := map[string]interface{}{
		"rollback": map[string]interface{}{
			"volume_id": d.Get("volume_id"),
			"name":      utils.ValueIgnoreEmpty(d.Get("name")),
		},
	}

	rollbackSnapshotOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(rollbackSnapshotBodyParams),
	}

	_, err = rollbackSnapshotClient.Request("POST", rollbackSnapshotPath, &rollbackSnapshotOpt)
	if err != nil {
		return diag.Errorf("failed to rollback snapshot to volume: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return nil
}

func resourceV5SnapshotRollBackRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceV5SnapshotRollBackUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceV5SnapshotRollBackDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to rollback snapshot.
Deleting this resource will not reset the rollbacked snapshot, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
