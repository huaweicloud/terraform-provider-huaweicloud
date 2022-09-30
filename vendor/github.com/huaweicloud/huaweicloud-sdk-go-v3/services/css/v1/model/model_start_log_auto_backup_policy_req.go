package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StartLogAutoBackupPolicyReq struct {

	// 备份开始时间。格式：格林威治标准时间。
	Period string `json:"period"`
}

func (o StartLogAutoBackupPolicyReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartLogAutoBackupPolicyReq struct{}"
	}

	return strings.Join([]string{"StartLogAutoBackupPolicyReq", string(data)}, " ")
}
