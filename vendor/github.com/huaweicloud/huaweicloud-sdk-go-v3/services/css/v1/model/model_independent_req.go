package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type IndependentReq struct {
	Type *IndependentBodyReq `json:"type"`
}

func (o IndependentReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IndependentReq struct{}"
	}

	return strings.Join([]string{"IndependentReq", string(data)}, " ")
}
