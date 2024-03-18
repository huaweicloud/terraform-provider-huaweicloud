package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ServerlessInfoResponse struct {

	// Serverless型实例的算力范围最小值。取值范围：0.5 ~ 8，单位：RCU。
	MinComputeUnit string `json:"min_compute_unit"`

	// Serverless型实例的算力范围最大值。取值范围：0.5 ~ 8，单位：RCU。
	MaxComputeUnit string `json:"max_compute_unit"`
}

func (o ServerlessInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ServerlessInfoResponse struct{}"
	}

	return strings.Join([]string{"ServerlessInfoResponse", string(data)}, " ")
}
