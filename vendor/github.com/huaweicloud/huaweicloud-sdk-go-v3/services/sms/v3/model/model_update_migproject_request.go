package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateMigprojectRequest struct {
	// 迁移项目id

	MigProjectId string `json:"mig_project_id"`

	Body *MigProject `json:"body,omitempty"`
}

func (o UpdateMigprojectRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateMigprojectRequest struct{}"
	}

	return strings.Join([]string{"UpdateMigprojectRequest", string(data)}, " ")
}
