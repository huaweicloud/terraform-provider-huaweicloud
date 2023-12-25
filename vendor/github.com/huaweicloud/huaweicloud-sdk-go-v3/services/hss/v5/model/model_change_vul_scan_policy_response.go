package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeVulScanPolicyResponse Response Object
type ChangeVulScanPolicyResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o ChangeVulScanPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeVulScanPolicyResponse struct{}"
	}

	return strings.Join([]string{"ChangeVulScanPolicyResponse", string(data)}, " ")
}
