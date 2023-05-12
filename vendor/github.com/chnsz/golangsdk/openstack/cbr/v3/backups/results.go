package backups

type getResp struct {
	// The backup detail.
	Backup BackupResp `json:"backup"`
}

// BackupResp is the structure that represents the backup detail.
type BackupResp struct {
	// The restore point ID
	CheckpointId string `json:"checkpoint_id"`
	// The creation time of the backup.
	CreatedAt string `json:"created_at"`
	// The backup description.
	Description string `json:"description"`
	// The expiration time of the backup.
	ExpiredAt string `json:"expired_at"`
	// The extended information.
	ExtendInfo BackupExtendInfo `json:"extend_info"`
	// The backup ID.
	ID string `json:"id"`
	// The backup type.
	ImageType string `json:"image_type"`
	// The backup name.
	Name string `json:"name"`
	// The parent backup ID.
	ParentId string `json:"parent_id"`
	// The project ID to which the backup belongs.
	ProjectId string `json:"project_id"`
	// Backup time.
	ProtectedAt string `json:"protected_at"`
	// The availability zone where the backup resource is located.
	ResourceAz string `json:"resource_az"`
	// The backup resource ID.
	ResourceId string `json:"resource_id"`
	// The backup resource name.
	ResourceName string `json:"resource_name"`
	// The backup resource size, in GB.
	ResourceSize int `json:"resource_size"`
	// The backup resource type.
	ResourceType string `json:"resource_type"`
	// The backup status.
	Status string `json:"status"`
	// The latest update time of the backup.
	UpdatedAt string `json:"updated_at"`
	// The vault to which the backup resource belongs.
	VaultId string `json:"vault_id"`
	// The replication records.
	ReplicationRecords []ReplicationRecord `json:"replication_record"`
	// The enterprise project to which the backup resource belongs.
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// The provider ID.
	ProviderId string `json:"provider_id"`
	// The backup list of the child resources.
	Children []BackupResp `json:"children"`
}

// BackupExtendInfo is an object that represents the extended information of the backup.
type BackupExtendInfo struct {
	// Whether the backup is automatically generated.
	AutoTrigger bool `json:"auto_trigger"`
	// Whether the backup is a system disk backup.
	Bootable bool `json:"bootable"`
	// Whether the backup is an incremental backup.
	Incremental bool `json:"incremental"`
	// Snapshot ID of the disk backup.
	SnapshotId string `json:"snapshot_id"`
	// Whether to allow lazyloading for fast restoration.
	SupportLld bool `json:"support_lld"`
	// The restoration mode.
	SupportRestoreMode string `json:"supported_restore_mode"`
	// The ID list of images created using backups.
	OsImagesData []ImageData `json:"os_image_data"`
	// Whether the VM backup data contains system disk data.
	ContainSystemDisk bool `json:"contain_system_disk"`
	// Whether the backup is encrypted.
	Encrypted bool `json:"encrypted"`
	// Whether the disk is a system disk.
	SystemDisk bool `json:"system_disk"`
}

// ImageData is an object that represents the backup image detail.
type ImageData struct {
	// Backup image ID.
	ImageId string `json:"image_id"`
}

// ReplicationRecord is an object that represents the replication record detail.
type ReplicationRecord struct {
	// The creation time of the replication.
	CreatedAt string `json:"created_at"`
	// The ID of the destination backup used for replication.
	DestinationBackupId string `json:"destination_backup_id"`
	// The record ID of the destination backup used for replication.
	DestinationCheckpointId string `json:"destination_checkpoint_id"`
	// The ID of the replication destination project.
	DestinationProjectId string `json:"destination_project_id"`
	// The replication destination region.
	DestinationRegion string `json:"destination_region"`
	// The destination vault ID.
	DestinationVaultId string `json:"destination_vault_id"`
	// The additional information of the replication.
	ExtraInfo ReplicationRecordExtraInfo `json:"extra_info"`
	// The replication record ID.
	ID string `json:"id"`
	// The ID of the source backup used for replication.
	SourceBackupId string `json:"source_backup_id"`
	// The ID of the source backup record used for replication.
	SourceCheckpointId string `json:"source_checkpoint_id"`
	// The ID of the replication source project.
	SourceProjectId string `json:"source_project_id"`
	// The replication source region.
	SourceRegion string `json:"source_region"`
	// The replication status.
	Status string `json:"status"`
	// The ID of the vault where the backup resides.
	VaultId string `json:"vault_id"`
}

// ReplicationRecordExtraInfo is an object that represents the additional information of the replication.
type ReplicationRecordExtraInfo struct {
	// The replication progress.
	Progress int `json:"progress"`
	// The error code.
	FailCode string `json:"fail_code"`
	// The error cause.
	FailReason string `json:"fail_reason"`
	// Whether replication is automatically scheduled.
	AutoTrigger bool `json:"auto_trigger"`
	// The destination vault ID.
	DestinationVaultId string `json:"destination_vault_id"`
}
