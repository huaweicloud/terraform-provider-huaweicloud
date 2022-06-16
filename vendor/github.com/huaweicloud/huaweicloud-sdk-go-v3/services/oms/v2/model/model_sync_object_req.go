package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 同步事件请求体
type SyncObjectReq struct {

	// 待同步对象的列表,其中待同步对象最大数量为10,列表中object_key为URL编码处理后的结果
	ObjectKeys []string `json:"object_keys"`
}

func (o SyncObjectReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SyncObjectReq struct{}"
	}

	return strings.Join([]string{"SyncObjectReq", string(data)}, " ")
}
