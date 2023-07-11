package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateYmlsReqEditModify 配置文件操作。 - modify: 修改参数配置。 - delete: 删除参数配置。 - reset: 重置参数配置。
type UpdateYmlsReqEditModify struct {

	// 参数配置列表。值为需要修改的json数据。
	ElasticsearchYml *interface{} `json:"elasticsearch.yml"`
}

func (o UpdateYmlsReqEditModify) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateYmlsReqEditModify struct{}"
	}

	return strings.Join([]string{"UpdateYmlsReqEditModify", string(data)}, " ")
}
