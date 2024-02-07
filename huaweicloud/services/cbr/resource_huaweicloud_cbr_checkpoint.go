package cbr

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cbr/v3/backups"
	"github.com/chnsz/golangsdk/openstack/cbr/v3/checkpoints"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CBR POST /v3/{project_id}/checkpoints
// @API CBR GET /v3/{project_id}/checkpoints/{checkpoint_id}
// @API CBR GET /v3/{project_id}/backups
// @API CBR DELETE /v3/{project_id}/backups/{backup_id}
func ResourceCheckpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCheckpointCreate,
		ReadContext:   resourceCheckpointRead,
		DeleteContext: resourceCheckpointDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the vault and backup resources are located.",
			},
			"vault_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the vault where the checkpoint to create.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the checkpoint.",
			},
			"backups": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The type of the backup resource.",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The ID of backup resource.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The backup ID.",
						},
						"resource_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The backup resource size.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The backup status.",
						},
						"protected_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The backup time.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The latest update time of the backup.",
						},
					},
				},
				Description: "The list of backups configuration.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The description of the checkpoint.",
			},
			"incremental": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether the backups are incremental backups.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the checkpoint.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the checkpoint.",
			},
		},
	}
}

func buildResourceDetails(resources *schema.Set) []checkpoints.Resource {
	if resources.Len() < 1 {
		return nil
	}

	result := make([]checkpoints.Resource, resources.Len())
	for i, val := range resources.List() {
		res := val.(map[string]interface{})
		result[i] = checkpoints.Resource{
			Type: res["type"].(string),
			ID:   res["resource_id"].(string),
		}
	}
	return result
}

func resourceCheckpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.CbrV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	opts := checkpoints.CreateOpts{
		VaultId: d.Get("vault_id").(string),
		Parameters: checkpoints.CheckpointParameter{
			Name:            d.Get("name").(string),
			Description:     d.Get("description").(string),
			Incremental:     utils.Bool(d.Get("incremental").(bool)),
			ResourceDetails: buildResourceDetails(d.Get("backups").(*schema.Set)),
		},
	}
	resp, err := checkpoints.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating CBR checkpoint: %s", err)
	}
	checkpointId := resp.ID
	d.SetId(checkpointId)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      checkpointStateRefreshFunc(client, checkpointId, []string{"available"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        30 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the checkpoint (%s) to become available: %s", checkpointId, err)
	}

	stateConf = &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      backupsStateRefreshFunc(client, checkpointId, []string{"available"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        30 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for all backups to become available: %s", err)
	}

	return resourceCheckpointRead(ctx, d, meta)
}

func flattenBackupResources(backupList []backups.BackupResp) []map[string]interface{} {
	if len(backupList) < 1 {
		return nil
	}
	result := make([]map[string]interface{}, len(backupList))
	for i, backup := range backupList {
		result[i] = map[string]interface{}{
			"type":          backup.ResourceType,
			"resource_id":   backup.ResourceId,
			"id":            backup.ID,
			"resource_size": backup.ResourceSize,
			"status":        backup.Status,
			"protected_at":  backup.ProtectedAt,
			"updated_at":    backup.UpdatedAt,
		}
	}
	return result
}

func resourceCheckpointRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CbrV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	var (
		checkpointId = d.Id()
		opts         = backups.ListOpts{
			CheckpointId: checkpointId,
			Incremental:  d.Get("incremental").(bool), // Ensure the correct results can be queried.
		}
	)
	// Query the checkpoint details: vault_id, status, created_at.
	checkpoint, err := checkpoints.Get(client, checkpointId)
	if err != nil {
		return common.CheckDeletedDiag(d, parseBackupError(err), "CBR checkpoint")
	}
	// Query the backup list under the specified checkpoint.
	backupList, err := backups.List(client, opts)
	if err != nil {
		return common.CheckDeletedDiag(d, parseBackupError(err), "CBR backups")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("vault_id", checkpoint.Vault.ID),
		d.Set("backups", flattenBackupResources(backupList)),
		d.Set("created_at", checkpoint.CreatedAt),
		d.Set("status", checkpoint.Status),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving CBR checkpoint: %s", err)
	}
	return nil
}

func resourceCheckpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.CbrV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	var (
		checkpointId = d.Id()
		mErr         *multierror.Error

		listOpts = backups.ListOpts{
			CheckpointId: checkpointId,
		}
	)
	backupResp, err := backups.List(client, listOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, parseBackupError(err), "CBR backups")
	}

	for _, backup := range backupResp {
		err = backups.Delete(client, backup.ID)
		analysedErr := parseBackupError(err)
		if analysedErr != nil {
			if _, ok := analysedErr.(golangsdk.ErrDefault404); !ok {
				mErr = multierror.Append(mErr, fmt.Errorf("error deleting CBR backup (%s): %s", backup.ID, err))
			}
		}
	}
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      backupsStateRefreshFunc(client, checkpointId, nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        20 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the backups to be deleted: %s", err)
	}
	return nil
}

func checkpointStateRefreshFunc(client *golangsdk.ServiceClient, checkpointId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := checkpoints.Get(client, checkpointId)
		if err != nil {
			analysedErr := parseBackupError(err)
			if _, ok := analysedErr.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return resp, "COMPLETED", nil
			}
			return resp, "ERROR", analysedErr
		}

		checkpointStatus := resp.Status
		// Unexpected status.
		if utils.StrSliceContains([]string{"error"}, checkpointStatus) {
			return resp, "ERROR", fmt.Errorf("unexpect status (%s)", checkpointStatus)
		}

		if utils.StrSliceContains(targets, checkpointStatus) {
			return resp, "COMPLETED", nil
		}
		return resp, "PENDING", nil
	}
}

func backupsStateRefreshFunc(client *golangsdk.ServiceClient, checkpointId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := checkpoints.Get(client, checkpointId)
		if err != nil {
			analysedErr := parseBackupError(err)
			if _, ok := analysedErr.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return resp, "COMPLETED", nil
			}
			return resp, "ERROR", analysedErr
		}

		backupRes := resp.Vault.Resources
		if len(backupRes) < 1 {
			return resp, "ERROR", golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Body: []byte("unable to find the resource backups"),
				},
			}
		}

		for _, resource := range backupRes {
			backupStatus := resource.ProtectStatus
			// Unexpected status.
			if utils.StrSliceContains([]string{"error"}, backupStatus) {
				return resp, "ERROR", fmt.Errorf("unexpect status (%s)", backupStatus)
			}
			if !utils.StrSliceContains(targets, backupStatus) {
				return resp, "PENDING", nil
			}
		}
		return resp, "COMPLETED", nil
	}
}

type backupNotFoundError struct {
	// Error code
	ErrCode string `json:"code"`
	// Error message
	ErrMessage string `json:"message"`
}

func parseBackupError(respErr error) error {
	var apiError struct {
		Unauthorized backupNotFoundError `json:"unauthorized"`
	}

	if errCode, ok := respErr.(golangsdk.ErrDefault401); ok {
		pErr := json.Unmarshal(errCode.Body, &apiError)
		if pErr == nil && (apiError.Unauthorized.ErrCode == "401" &&
			apiError.Unauthorized.ErrMessage == "Malformed request url") {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return respErr
}
