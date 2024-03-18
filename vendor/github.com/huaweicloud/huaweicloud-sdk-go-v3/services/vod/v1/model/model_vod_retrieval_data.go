package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type VodRetrievalData struct {

	// 低频取回量
	RetrievalWarm *float64 `json:"retrieval_warm,omitempty"`

	// 归档标准取回量
	RetrievalCold *float64 `json:"retrieval_cold,omitempty"`

	// 归档快速取回量
	RetrievalColdSpeed *float64 `json:"retrieval_cold_speed,omitempty"`
}

func (o VodRetrievalData) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VodRetrievalData struct{}"
	}

	return strings.Join([]string{"VodRetrievalData", string(data)}, " ")
}
