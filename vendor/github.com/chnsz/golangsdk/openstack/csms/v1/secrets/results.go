package secrets

type SecretRst struct {
	Secret Secret `json:"secret"`
}

type Secret struct {
	ID                  string        `json:"id"`
	Name                string        `json:"name"`
	State               string        `json:"state"`
	KmsKeyID            string        `json:"kms_key_id"`
	Description         string        `json:"description"`
	CreateTime          int           `json:"create_time"`
	UpdateTime          int           `json:"update_time"`
	ScheduledDeleteTime int           `json:"scheduled_delete_time"`
	SecretType          string        `json:"secret_type"`
	EnterpriseProjectID string        `json:"enterprise_project_id"`
	EventSubscriptions  []interface{} `json:"event_subscriptions"`
}

type Version struct {
	VersionMetadata VersionMetadata `json:"version_metadata"`
	SecretBinary    string          `json:"secret_binary"`
	SecretString    string          `json:"secret_string"`
}

// VersionMetadata 凭据版本被标记的状态。
type VersionMetadata struct {
	ID            string   `json:"id"`
	CreateTime    int      `json:"create_time"`
	KmsKeyID      string   `json:"kms_key_id"`
	SecretName    string   `json:"secret_name"`
	VersionStages []string `json:"version_stages"`
	ExpireTime    int      `json:"expire_time"`
}
