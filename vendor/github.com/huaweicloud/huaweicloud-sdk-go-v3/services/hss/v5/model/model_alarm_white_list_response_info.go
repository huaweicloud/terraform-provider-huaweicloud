package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 告警白名单详情
type AlarmWhiteListResponseInfo struct {

	// SHA256
	Hash *string `json:"hash,omitempty"`

	// 描述信息
	Description *string `json:"description,omitempty"`

	// 事件类型，包含如下:   - 1001 : 恶意软件   - 1010 : Rootkit   - 1011 : 勒索软件   - 1015 : Webshell   - 1017 : 反弹Shell   - 2001 : 一般漏洞利用   - 2047 : Redis漏洞利用   - 2048 : Hadoop漏洞利用   - 2049 : MySQL漏洞利用   - 3002 : 文件提权   - 3003 : 进程提权   - 3004 : 关键文件变更   - 3005 : 文件/目录变更   - 3007 : 进程异常行为   - 3015 : 高危命令执行   - 3018 : 异常Shell   - 3027 : Crontab可疑任务   - 4002 : 暴力破解   - 4004 : 异常登录   - 4006 : 非法系统账号
	EventType *int32 `json:"event_type,omitempty"`

	// 更新时间，毫秒
	UpdateTime *int64 `json:"update_time,omitempty"`
}

func (o AlarmWhiteListResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AlarmWhiteListResponseInfo struct{}"
	}

	return strings.Join([]string{"AlarmWhiteListResponseInfo", string(data)}, " ")
}
