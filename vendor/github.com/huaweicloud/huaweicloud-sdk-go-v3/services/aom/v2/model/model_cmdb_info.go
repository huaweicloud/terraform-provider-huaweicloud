package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CmdbInfo struct {

	// 应用id。
	AppId *string `json:"app_id,omitempty"`

	// 节点信息列表。
	NodeIds *[]NodeInfo `json:"node_ids,omitempty"`
}

func (o CmdbInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CmdbInfo struct{}"
	}

	return strings.Join([]string{"CmdbInfo", string(data)}, " ")
}
