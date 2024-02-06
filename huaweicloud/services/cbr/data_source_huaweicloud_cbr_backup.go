package cbr

import (
	"context"
	"reflect"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/cbr/v3/backups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API CBR GET /v3/{project_id}/backups/{backup_id}
func DataSourceBackup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBackupRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The region where the CBR backup is located.`,
			},
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The backup ID.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The parent backup ID.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backup type.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backup name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backup description.`,
			},
			"checkpoint_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The restore point ID.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backup resource ID.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backup resource type.`,
			},
			"resource_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The backup resource size, in GB.`,
			},
			"resource_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backup resource name.`,
			},
			"resource_az": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The availability zone where the backup resource is located.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise project to which the backup resource belongs.`,
			},
			"vault_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The vault to which the backup resource belongs.`,
			},
			"replication_records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: schemaReplicationRecordDetail(),
				},
				Description: `The replication records.`,
			},
			"children": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: schemaChildrenBackupDetail(),
				},
				Description: `The backup list of the sub-backup resources.`,
			},
			"extend_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: schemaBackupExtendDetail(),
				},
				Description: `The extended information.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backup status.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the backup.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the backup.`,
			},
			"expired_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The expiration time of the backup.`,
			},
		},
	}
}

func schemaChildrenBackupDetail() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The sub-backup ID.`,
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The sub-backup name.`,
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The sub-backup description.`,
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The sub-backup type.`,
		},
		"checkpoint_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The restore point ID of the sub-backup resource.`,
		},
		"resource_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The sub-backup resource ID.`,
		},
		"resource_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The sub-backup resource type.`,
		},
		"resource_size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: `The sub-backup resource size, in GB.`,
		},
		"resource_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The sub-backup resource name.`,
		},
		"resource_az": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The availability zone where the backup sub-backup resource is located.`,
		},
		"enterprise_project_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The enterprise project to which the backup sub-backup resource belongs.`,
		},
		"vault_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The vault to which the backup sub-backup resource belongs.`,
		},
		"replication_records": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: schemaReplicationRecordDetail(),
			},
			Description: `The replication records.`,
		},
		"extend_info": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: schemaBackupExtendDetail(),
			},
			Description: `The extended information.`,
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The sub-backup status.`,
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The creation time of the sub-backup.`,
		},
		"updated_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The latest update time of the sub-backup.`,
		},
		"expired_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The expiration time of the sub-backup.`,
		},
	}
}

func schemaReplicationRecordDetail() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The replication record ID.`,
		},
		"destination_backup_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The ID of the destination backup used for replication.`,
		},
		"destination_checkpoint_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The record ID of the destination backup used for replication.`,
		},
		"destination_project_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The ID of the replication destination project.`,
		},
		"destination_region": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The replication destination region.`,
		},
		"destination_vault_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The destination vault ID.`,
		},
		"source_backup_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The ID of the source backup used for replication.`,
		},
		"source_checkpoint_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The ID of the source backup record used for replication.`,
		},
		"source_project_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The ID of the replication source project.`,
		},
		"source_region": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The replication source region.`,
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The replication status.`,
		},
		"vault_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The ID of the vault where the backup resides.`,
		},
		"extra_info": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"progress": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: `The replication progress.`,
					},
					"fail_code": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: `The error code.`,
					},
					"fail_reason": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: `The error cause.`,
					},
					"auto_trigger": {
						Type:        schema.TypeBool,
						Computed:    true,
						Description: `Whether replication is automatically scheduled.`,
					},
					"destination_vault_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: `The destination vault ID.`,
					},
				},
			},
			Description: `The additional information of the replication.`,
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The creation time of the replication.`,
		},
	}
}

func schemaBackupExtendDetail() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"auto_trigger": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: `Whether the backup is automatically generated.`,
		},
		"bootable": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: `Whether the backup is a system disk backup.`,
		},
		"incremental": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: `Whether the backup is an incremental backup.`,
		},
		"snapshot_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `Snapshot ID of the disk backup.`,
		},
		"support_lld": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: `Whether to allow lazyloading for fast restoration.`,
		},
		"supported_restore_mode": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The restoration mode.`,
		},
		"os_registry_images": {
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: `The ID list of images created using backups.`,
		},
		"contain_system_disk": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: `Whether the VM backup data contains system disk data.`,
		},
		"encrypted": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: `Whether the backup is encrypted.`,
		},
		"is_system_disk": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: `Whether the disk is a system disk.`,
		},
	}
}

func flattenReplicationRecordExtraInfo(extraInfo backups.ReplicationRecordExtraInfo) []map[string]interface{} {
	if reflect.DeepEqual(extraInfo, backups.ReplicationRecordExtraInfo{}) {
		return nil
	}

	return []map[string]interface{}{
		{
			"progress":             extraInfo.Progress,
			"fail_code":            extraInfo.FailCode,
			"fail_reason":          extraInfo.FailReason,
			"auto_trigger":         extraInfo.AutoTrigger,
			"destination_vault_id": extraInfo.DestinationVaultId,
		},
	}
}

func flattenReplicationRecords(records []backups.ReplicationRecord) []map[string]interface{} {
	if len(records) < 1 {
		return nil
	}
	result := make([]map[string]interface{}, len(records))
	for i, record := range records {
		result[i] = map[string]interface{}{
			"id":                        record.ID,
			"destination_backup_id":     record.DestinationBackupId,
			"destination_checkpoint_id": record.DestinationCheckpointId,
			"destination_project_id":    record.DestinationProjectId,
			"destination_region":        record.DestinationRegion,
			"destination_vault_id":      record.DestinationVaultId,
			"source_backup_id":          record.SourceBackupId,
			"source_checkpoint_id":      record.SourceCheckpointId,
			"source_project_id":         record.SourceProjectId,
			"source_region":             record.SourceRegion,
			"status":                    record.Status,
			"vault_id":                  record.VaultId,
			"extra_info":                flattenReplicationRecordExtraInfo(record.ExtraInfo),
		}
	}
	return result
}

func flattenChildBackups(children []backups.BackupResp) []map[string]interface{} {
	if len(children) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(children))
	for i, child := range children {
		result[i] = map[string]interface{}{
			"id":                    child.ID,
			"name":                  child.Name,
			"description":           child.Description,
			"type":                  child.ImageType,
			"resource_id":           child.ResourceId,
			"resource_type":         child.ResourceType,
			"resource_size":         child.ResourceSize,
			"resource_name":         child.ResourceName,
			"resource_az":           child.ResourceAz,
			"enterprise_project_id": child.EnterpriseProjectId,
			"vault_id":              child.VaultId,
			"replication_records":   flattenReplicationRecords(child.ReplicationRecords),
			"extend_info":           flattenExtendInfo(child.ExtendInfo),
			"status":                child.Status,
			"created_at":            child.CreatedAt,
			"updated_at":            child.UpdatedAt,
			"expired_at":            child.ExpiredAt,
		}
	}
	return result
}

func flattenRegistryImages(images []backups.ImageData) []string {
	if len(images) < 1 {
		return nil
	}
	result := make([]string, len(images))
	for i, image := range images {
		result[i] = image.ImageId
	}
	return result
}

func flattenExtendInfo(extendInfo backups.BackupExtendInfo) []map[string]interface{} {
	if reflect.DeepEqual(extendInfo, backups.BackupExtendInfo{}) {
		return nil
	}

	return []map[string]interface{}{
		{
			"auto_trigger":           extendInfo.AutoTrigger,
			"bootable":               extendInfo.Bootable,
			"incremental":            extendInfo.Incremental,
			"snapshot_id":            extendInfo.SnapshotId,
			"support_lld":            extendInfo.SupportLld,
			"supported_restore_mode": extendInfo.SupportRestoreMode,
			"os_registry_images":     flattenRegistryImages(extendInfo.OsImagesData),
			"contain_system_disk":    extendInfo.ContainSystemDisk,
			"encrypted":              extendInfo.Encrypted,
			"is_system_disk":         extendInfo.SystemDisk,
		},
	}
}

func dataSourceBackupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CbrV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	backupId := d.Get("id").(string)
	resp, err := backups.Get(client, backupId)
	if err != nil {
		return diag.Errorf("error querying backup detail: %s", err)
	}
	d.SetId(backupId)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("parent_id", resp.ParentId),
		d.Set("type", resp.ImageType),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("checkpoint_id", resp.CheckpointId),
		d.Set("resource_id", resp.ResourceId),
		d.Set("resource_type", resp.ResourceType),
		d.Set("resource_size", resp.ResourceSize),
		d.Set("resource_name", resp.ResourceName),
		d.Set("resource_az", resp.ResourceAz),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("vault_id", resp.VaultId),
		d.Set("replication_records", flattenReplicationRecords(resp.ReplicationRecords)),
		d.Set("children", flattenChildBackups(resp.Children)),
		d.Set("extend_info", flattenExtendInfo(resp.ExtendInfo)),
		d.Set("status", resp.Status),
		d.Set("created_at", resp.CreatedAt),
		d.Set("updated_at", resp.UpdatedAt),
		d.Set("expired_at", resp.ExpiredAt),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving backup data-source fields: %s", err)
	}
	return nil
}
