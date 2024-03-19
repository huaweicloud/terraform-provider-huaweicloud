package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TempContentInfo struct {

	// 报文id或者事务id或者插件id
	ContentId *int32 `json:"content_id,omitempty"`

	// 内容
	Content *[]Content `json:"content,omitempty"`

	// 索引
	Index *int32 `json:"index,omitempty"`

	// 数据指令内容
	Data *interface{} `json:"data,omitempty"`

	// 数据指令类型（0：默认请求卡片；1：数据指令；201：循环指令；202：条件指令；301：集合点[；203：vu百分比控制器；204：吞吐量控制器；302：插件请求](tag:hws,hws_hk)）
	DataType *int32 `json:"data_type,omitempty"`
}

func (o TempContentInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TempContentInfo struct{}"
	}

	return strings.Join([]string{"TempContentInfo", string(data)}, " ")
}
