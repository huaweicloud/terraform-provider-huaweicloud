package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type IndependentRequestBody struct {
	Type *IndependentReq `json:"type"`
}

func (o IndependentRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IndependentRequestBody struct{}"
	}

	return strings.Join([]string{"IndependentRequestBody", string(data)}, " ")
}
