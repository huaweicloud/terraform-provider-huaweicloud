package dws

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS POST /v1.0/{project_id}/snapshots/{snapshot_id}/linked-copy
// @API DWS GET /v1.0/{project_id}/snapshots/{snapshot_id}
// @API DWS DELETE /v1.0/{project_id}/snapshots/{snapshot_id}
func ResourceSnapshotCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSnapshotCopyCreate,
		ReadContext:   resourceSnapshotCopyRead,
		DeleteContext: resourceSnapshotCopyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSnapshotCopyImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"snapshot_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the ID of the automated snapshot to be copied.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the copy snapshot.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the description of the copy snapshot.`,
			},
		},
	}
}

func resourceSnapshotCopyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	snapshotId := d.Get("snapshot_id").(string)
	// For the same automatic snapshot, you cannot perform the copy snapshot and restore cluster operations at the same time.
	config.MutexKV.Lock(snapshotId)
	config.MutexKV.Unlock(snapshotId)

	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	httpUrl := "v1.0/{project_id}/snapshots/{snapshot_id}/linked-copy"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{snapshot_id}", d.Get("snapshot_id").(string))

	opt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildSnapshotCopyBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating DWS snapshot copy: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("snapshot_id", respBody, nil)
	if err != nil {
		return diag.Errorf("error copying DWS automated snapshot: ID is not found in API response")
	}

	d.SetId(id.(string))

	return resourceSnapshotCopyRead(ctx, d, meta)
}

func buildSnapshotCopyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"backup_name": d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceSnapshotCopyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	respBody, err := GetSnapshotById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving copied snapshot")
	}
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("snapshot.name", respBody, nil)),
		d.Set("description", utils.PathSearch("snapshot.description", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSnapshotCopyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	if err = deleteSnapshotById(client, d.Id()); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting copied snapshot")
	}
	return nil
}

func resourceSnapshotCopyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<snapshot_id>/<id>', but got '%s'",
			importedId)
	}

	mErr := multierror.Append(nil, d.Set("snapshot_id", parts[0]))
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
