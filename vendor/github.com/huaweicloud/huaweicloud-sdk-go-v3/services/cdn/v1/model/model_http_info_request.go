package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type HttpInfoRequest struct {
	Https *HttpInfoRequestBody `json:"https"`
}

func (o HttpInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HttpInfoRequest struct{}"
	}

	return strings.Join([]string{"HttpInfoRequest", string(data)}, " ")
}
