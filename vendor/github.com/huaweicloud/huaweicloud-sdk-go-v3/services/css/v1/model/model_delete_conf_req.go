package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DeleteConfReq struct {

	// 配置文件名称。
	Name string `json:"name"`
}

func (o DeleteConfReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteConfReq struct{}"
	}

	return strings.Join([]string{"DeleteConfReq", string(data)}, " ")
}
