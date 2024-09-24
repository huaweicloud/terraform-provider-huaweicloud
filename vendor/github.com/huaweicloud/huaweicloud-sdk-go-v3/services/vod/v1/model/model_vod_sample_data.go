package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type VodSampleData struct {

	// 存储空间。  单位：GB。
	Storage *float32 `json:"storage,omitempty"`

	// 低频存储空间。  单位：GB。
	StorageWarm *float32 `json:"storage_warm,omitempty"`

	// 归档存储空间。  单位：GB。
	StorageCold *float32 `json:"storage_cold,omitempty"`

	// 转码时长。  单位：分钟。
	Transcode *int64 `json:"transcode,omitempty"`
}

func (o VodSampleData) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VodSampleData struct{}"
	}

	return strings.Join([]string{"VodSampleData", string(data)}, " ")
}
