package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddNodesToNodePoolResponse Response Object
type AddNodesToNodePoolResponse struct {

	// 提交任务成功后返回的任务ID，用户可以使用该ID对任务执行情况进行查询。
	Jobid          *string `json:"jobid,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o AddNodesToNodePoolResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddNodesToNodePoolResponse struct{}"
	}

	return strings.Join([]string{"AddNodesToNodePoolResponse", string(data)}, " ")
}
