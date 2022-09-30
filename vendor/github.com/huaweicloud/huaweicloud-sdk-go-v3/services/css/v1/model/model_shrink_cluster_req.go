package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ShrinkClusterReq struct {

	// 需要缩容的节点类型和数量集合。
	Shrink []ShrinkNodeReq `json:"shrink"`
}

func (o ShrinkClusterReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShrinkClusterReq struct{}"
	}

	return strings.Join([]string{"ShrinkClusterReq", string(data)}, " ")
}
