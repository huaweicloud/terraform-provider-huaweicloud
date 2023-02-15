package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ObsBucket struct {

	// 桶名称
	Bucket *string `json:"bucket,omitempty"`

	// 桶的创建时间
	CreationDate *string `json:"creation_date,omitempty"`

	// 授权结果，取值[0,1]，0表示未授权给转码服务，1表示已授权转码服务
	IsAuthorized *int32 `json:"is_authorized,omitempty"`
}

func (o ObsBucket) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ObsBucket struct{}"
	}

	return strings.Join([]string{"ObsBucket", string(data)}, " ")
}
