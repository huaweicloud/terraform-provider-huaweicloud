package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SetOffSiteBackupPolicyRequestBody struct {
	PolicyPara *OffSiteBackupPolicy `json:"policy_para"`
}

func (o SetOffSiteBackupPolicyRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetOffSiteBackupPolicyRequestBody struct{}"
	}

	return strings.Join([]string{"SetOffSiteBackupPolicyRequestBody", string(data)}, " ")
}
