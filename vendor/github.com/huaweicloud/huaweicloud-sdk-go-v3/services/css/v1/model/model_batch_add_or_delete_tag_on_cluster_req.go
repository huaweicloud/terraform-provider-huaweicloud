package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BatchAddOrDeleteTagOnClusterReq struct {

	// 操作类型。通过该属性标识当前所需的操作类型。  - create：批量添加标签。 - delete：批量删除标签。
	Action string `json:"action"`

	// 标签列表。
	Tags []Tag `json:"tags"`

	// 系统标签列表。 - key固定为：_sys_enterprise_project_id。 - value是UUID或0,value为0表示默认企业项目。
	SysTags *[]SysTags `json:"sysTags,omitempty"`
}

func (o BatchAddOrDeleteTagOnClusterReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchAddOrDeleteTagOnClusterReq struct{}"
	}

	return strings.Join([]string{"BatchAddOrDeleteTagOnClusterReq", string(data)}, " ")
}
