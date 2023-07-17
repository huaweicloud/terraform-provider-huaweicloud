package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UnBindPublicReqEipReq 弹性IP信息。
type UnBindPublicReqEipReq struct {
	BandWidth *BindPublicReqEipBandWidth `json:"bandWidth,omitempty"`
}

func (o UnBindPublicReqEipReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UnBindPublicReqEipReq struct{}"
	}

	return strings.Join([]string{"UnBindPublicReqEipReq", string(data)}, " ")
}
