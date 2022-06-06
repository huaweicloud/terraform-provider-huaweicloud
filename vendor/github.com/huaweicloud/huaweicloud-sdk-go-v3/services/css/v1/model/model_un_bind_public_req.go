package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UnBindPublicReq struct {
	Eip *BindPublicReqEipReq `json:"eip"`
}

func (o UnBindPublicReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UnBindPublicReq struct{}"
	}

	return strings.Join([]string{"UnBindPublicReq", string(data)}, " ")
}
