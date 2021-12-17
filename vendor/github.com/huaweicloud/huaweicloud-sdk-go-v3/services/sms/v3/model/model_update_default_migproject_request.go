package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type UpdateDefaultMigprojectRequest struct {
	// 迁移项目ID

	MigProjectId string `json:"mig_project_id"`
}

func (o UpdateDefaultMigprojectRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDefaultMigprojectRequest struct{}"
	}

	return strings.Join([]string{"UpdateDefaultMigprojectRequest", string(data)}, " ")
}
