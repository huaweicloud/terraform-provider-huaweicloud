package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeDetailResponse Response Object
type UpgradeDetailResponse struct {

	// 下发执行接口次数。
	TotalSize *int32 `json:"totalSize,omitempty"`

	DetailList     *[]GetUpgradeDetailInfo `json:"detailList,omitempty"`
	HttpStatusCode int                     `json:"-"`
}

func (o UpgradeDetailResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeDetailResponse struct{}"
	}

	return strings.Join([]string{"UpgradeDetailResponse", string(data)}, " ")
}
