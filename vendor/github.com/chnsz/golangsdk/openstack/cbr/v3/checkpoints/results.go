package checkpoints

// Checkpoint is the structure that represents the checkpoint and corresponding backup resources detail.
type Checkpoint struct {
	// The restore point ID.
	ID string `json:"id"`
	// The project ID.
	ProjectId string `json:"project_id"`
	// The creation time, the format is 'HHHH-MM-DDThh:mm:dd.000000'.
	CreatedAt string `json:"created_at"`
	// The status information.
	// The valid values are as follows:
	// + available
	// + deleting
	// + protecting
	// + deleted
	// + error-deleting
	// + error
	Status string `json:"status"`
	// The vault information.
	Vault VaultInfo `json:"vault"`
	// The extended information.
	ExtraInfo CheckpointExtraInfo `json:"extra_info"`
}

// VaultInfo is the structure that represents the backup storage information.
type VaultInfo struct {
	// The vault ID.
	ID string `json:"id"`
	// The vault name.
	Name string `json:"name"`
	// The backup objects.
	Resources []CheckpointResource `json:"resources"`
	// The resources skipped during backup.
	SkippedResources []SkippedResource `json:"skipped_resources"`
}

// CheckpointResource is the structure that represents the backup resources detail.
type CheckpointResource struct {
	// Extra information of the resource.
	ExtraInfo string `json:"extra_info"`
	// ID of the resource to be backed up.
	ID string `json:"id"`
	// Name of the resource to be backed up.
	Name string `json:"name"`
	// Protected status.
	// The valid values are as follows:
	// + available
	// + error
	// + protecting
	// + restoring
	// + removing
	ProtectStatus string `json:"protect_status"`
	// Allocated capacity for the associated resource, in GB.
	ResourceSize string `json:"resource_size"`
	// The type of the resource to be backed up, which can be:
	// + OS::Nova::Server
	// + OS::Cinder::Volume
	// + OS::Ironic::BareMetalServer
	// + OS::Native::Server
	// + OS::Sfs::Turbo
	// + OS::Workspace::DesktopV2
	Type string `json:"type"`
	// Backup size.
	BackupSize string `json:"backup_size"`
	// Number of backups.
	BackupCount string `json:"backup_count"`
}

// CheckpointResource is the structure that represents the resources list which are not backed up.
type SkippedResource struct {
	// Resource ID.
	ID string `json:"id"`
	// Resource type.
	Type string `json:"type"`
	// Resource name.
	Name string `json:"name"`
	// Error code.
	Code string `json:"code"`
	// Reason for the skipping. For example, the resource is being backed up.
	Reason string `json:"reason"`
}

// CheckpointExtraInfo is the structure that represents the extra information of the checkpoint.
type CheckpointExtraInfo struct {
	// Backup name.
	Name string `json:"name"`
	// Backup description.
	Description string `json:"description"`
	// Number of days that backups can be retained.
	RetentionDuration int `json:"retention_duration"`
}
