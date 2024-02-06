package cbr

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/cbr/v3/backups"
	"github.com/chnsz/golangsdk/openstack/cbr/v3/members"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API CBR PUT /v3/{project_id}/backups/{backup_id}/members/{member_id}
// @API CBR GET /v3/{project_id}/backups/{backup_id}
func ResourceBackupShareAccepter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBackupShareAccepterCreate,
		ReadContext:   resourceBackupShareAccepterRead,
		DeleteContext: resourceBackupShareAccepterDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the backup will be stored.",
			},
			"backup_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the shared source backup.",
			},
			"vault_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the vault which the backup will be stored.",
			},
			"source_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the project to which the source backup belongs.",
			},
		},
	}
}

func resourceBackupShareAccepterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CbrV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	var (
		backupId = d.Get("backup_id").(string)
		opts     = members.UpdateOpts{
			BackupId: backupId,
			MemberId: client.ProjectID,
			Status:   "accepted",
			VaultId:  d.Get("vault_id").(string),
		}
	)

	_, err = members.Update(client, opts)
	if err != nil {
		return diag.Errorf("error modifying backup share (%s) status: %s", backupId, err)
	}
	// Use the backup ID as the ID of this resource.
	// Note: One backup can only manage one resource.
	d.SetId(backupId)

	return resourceBackupShareAccepterRead(ctx, d, meta)
}

func resourceBackupShareAccepterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CbrV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	resp, err := backups.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CBR backup share accepter")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("backup_id", resp.ID),
		d.Set("vault_id", resp.VaultId),
		d.Set("source_project_id", resp.ProjectId),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving resource fields of the CBR backup share accepter: %s", err)
	}
	return nil
}

func resourceBackupShareAccepterDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CbrV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	var (
		backupId = d.Id()
		opts     = members.UpdateOpts{
			BackupId: backupId,
			MemberId: client.ProjectID,
			Status:   "rejected",
		}
	)
	_, err = members.Update(client, opts)
	if err != nil {
		return diag.Errorf("error denying backup share: %s", err)
	}
	return nil
}
