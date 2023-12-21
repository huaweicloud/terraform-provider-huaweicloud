package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeDbMajorVersionResponse Response Object
type UpgradeDbMajorVersionResponse struct {

	// 任务ID。
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpgradeDbMajorVersionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeDbMajorVersionResponse struct{}"
	}

	return strings.Join([]string{"UpgradeDbMajorVersionResponse", string(data)}, " ")
}
