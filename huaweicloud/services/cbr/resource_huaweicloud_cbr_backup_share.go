package cbr

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cbr/v3/members"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API CBR POST /v3/{project_id}/backups/{backup_id}/members
// @API CBR GET /v3/{project_id}/backups/{backup_id}/members
// @API CBR DELETE /v3/{project_id}/backups/{backup_id}/members/{member_id}
func ResourceBackupShare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBackupShareCreate,
		ReadContext:   resourceBackupShareRead,
		UpdateContext: resourceBackupShareUpdate,
		DeleteContext: resourceBackupShareDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceBackupShareImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the shared backup is located.",
			},
			"backup_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The backup ID.",
			},
			"members": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dest_project_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the project with which the backup is shared.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the backup shared member record.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The backup shared status.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the backup shared member.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The latest update time of the backup shared member.",
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the image registered with the shared backup copy.",
						},
						"vault_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the vault where the shared backup is stored.",
						},
					},
				},
				Description: "The list of shared members configuration.",
			},
		},
	}
}

// Collect destination project IDs of all shared members.
func buildBackupSharedMembersInfo(memberList *schema.Set) []string {
	if memberList.Len() < 1 {
		return nil
	}
	result := make([]string, memberList.Len())
	for i, val := range memberList.List() {
		member := val.(map[string]interface{})
		result[i] = member["dest_project_id"].(string)
	}
	return result
}

func resourceBackupShareCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.CbrV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	var (
		backupId      = d.Get("backup_id").(string)
		sharedMembers = d.Get("members").(*schema.Set)
		opts          = members.CreateOpts{
			BackupId: backupId,
			Members:  buildBackupSharedMembersInfo(sharedMembers),
		}
	)

	_, err = members.Create(client, opts)
	if err != nil {
		return diag.Errorf("failed while adding shared member to specified backup (%s): %s", backupId, err)
	}
	// Use the backup ID as the ID of this resource.
	// Note: One backup can only manage one resource.
	d.SetId(backupId)

	return resourceBackupShareRead(ctx, d, meta)
}

func flattenBackupSharedMembersInfo(memberList []members.Member) []map[string]interface{} {
	if len(memberList) < 1 {
		return nil
	}
	result := make([]map[string]interface{}, len(memberList))
	for i, member := range memberList {
		result[i] = map[string]interface{}{
			"dest_project_id": member.DestProjectId,
			"id":              member.ID,
			"status":          member.Status,
			"created_at":      member.CreatedAt,
			"updated_at":      member.UpdatedAt,
			"image_id":        member.ImageId,
			"vault_id":        member.VaultId,
		}
	}
	return result
}

func resourceBackupShareRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CbrV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	var (
		backupId = d.Id()
		opts     = members.ListOpts{
			BackupId: backupId,
		}
	)
	// Query the shared member list under the specified backup.
	memberList, err := members.List(client, opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CBR backup share")
	}
	if len(memberList) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "CBR backup share")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("members", flattenBackupSharedMembersInfo(memberList)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving CBR backup shared: %s", err)
	}
	return nil
}

func resourceBackupShareUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.CbrV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	var (
		backupId = d.Id()
		mErr     *multierror.Error

		oldRaws, newRaws = d.GetChange("members")
		addSet           = newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
		rmSet            = oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
	)

	for _, val := range rmSet.List() {
		member := val.(map[string]interface{})
		err = members.Delete(client, backupId, member["dest_project_id"].(string))
		if err != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("error delete shared member (%s) from a specified CBR backup (%s): %s",
				member["id"].(string), backupId, err))
		}
	}

	if addSet.Len() > 0 {
		opts := members.CreateOpts{
			BackupId: backupId,
			Members:  buildBackupSharedMembersInfo(addSet),
		}
		_, err = members.Create(client, opts)
		if err != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("failed while adding shared member to specified backup (%s): %s",
				backupId, err))
		}
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return resourceBackupShareRead(ctx, d, meta)
}

func resourceBackupShareDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.CbrV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	var (
		backupId = d.Id()
		mErr     *multierror.Error

		listOpts = members.ListOpts{
			BackupId: backupId,
		}
	)
	memberList, err := members.List(client, listOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, parseBackupError(err), "CBR backup share")
	}

	for _, member := range memberList {
		err = members.Delete(client, backupId, member.DestProjectId)
		if err != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("error delete shared member (%s) from a specified CBR backup (%s): %s",
				member.ID, backupId, err))
		}
	}
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBackupShareImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	// The resource ID, which is also the backup ID.
	return []*schema.ResourceData{d}, d.Set("backup_id", d.Id())
}
