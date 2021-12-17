package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 命令参数body
type CommandBody struct {
	// 命令名称，分为：START、STOP、DELETE、SYNC、UPLOAD_LOG、RSET_LOG_ACL

	CommandName string `json:"command_name"`
	// 命令执行结果  success代表执行命令成功  fail代表命令执行失败

	Result string `json:"result"`
	// JSON格式的命令执行结果，只用于保存数据库，没有其他作用

	ResultDetail *interface{} `json:"result_detail"`
}

func (o CommandBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CommandBody struct{}"
	}

	return strings.Join([]string{"CommandBody", string(data)}, " ")
}
