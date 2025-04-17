package evs

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EVS GET /v2/{project_id}/cloudsnapshots/{snapshot_id}
// @API EVS POST /v2/{project_id}/cloudsnapshots/{snapshot_id}/rollback
// ResourceSnapshotRollBack is a definition of the one-time action resource that used to manage snapshot rollback.
func ResourceSnapshotRollBack() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSnapshotRollBackCreate,
		ReadContext:   resourceSnapshotRollBackRead,
		DeleteContext: resourceSnapshotRollBackDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSnapshotRollBackCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                     = meta.(*config.Config)
		region                  = cfg.GetRegion(d)
		rollbackSnapshotHttpUrl = "v2/{project_id}/cloudsnapshots/{snapshot_id}/rollback"
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

	// The API response is the same as the request body, and we do not need to retrieve any content from the response,
	// so ignore it here.
	_, err = rollbackSnapshotClient.Request("POST", rollbackSnapshotPath, &rollbackSnapshotOpt)
	if err != nil {
		return diag.Errorf("failed to rollback snapshot to volume: %s", err)
	}

	d.SetId(d.Get("snapshot_id").(string))

	return resourceSnapshotRollBackRead(ctx, d, meta)
}

func resourceSnapshotRollBackRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceSnapshotRollBackDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a one-time action resource.
	return nil
}
