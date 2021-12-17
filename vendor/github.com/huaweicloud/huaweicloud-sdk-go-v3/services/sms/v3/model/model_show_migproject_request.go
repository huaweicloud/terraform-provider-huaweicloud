package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowMigprojectRequest struct {
	// 迁移项目id

	MigProjectId string `json:"mig_project_id"`
}

func (o ShowMigprojectRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowMigprojectRequest struct{}"
	}

	return strings.Join([]string{"ShowMigprojectRequest", string(data)}, " ")
}
