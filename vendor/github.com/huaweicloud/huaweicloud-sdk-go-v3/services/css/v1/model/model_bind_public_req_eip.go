package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BindPublicReqEip 弹性IP信息。
type BindPublicReqEip struct {
	BandWidth *BindPublicReqEipBandWidth `json:"bandWidth"`
}

func (o BindPublicReqEip) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BindPublicReqEip struct{}"
	}

	return strings.Join([]string{"BindPublicReqEip", string(data)}, " ")
}
