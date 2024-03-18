package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ReduceVolumeRequestBody struct {
	ReduceVolume *ReduceVolumeObject `json:"reduce_volume"`
}

func (o ReduceVolumeRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReduceVolumeRequestBody struct{}"
	}

	return strings.Join([]string{"ReduceVolumeRequestBody", string(data)}, " ")
}
