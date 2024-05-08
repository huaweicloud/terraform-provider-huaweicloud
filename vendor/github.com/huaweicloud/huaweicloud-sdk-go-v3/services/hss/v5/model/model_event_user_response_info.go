package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// EventUserResponseInfo 用户信息
type EventUserResponseInfo struct {

	// 用户uid
	UserId *int32 `json:"user_id,omitempty"`

	// 用户gid
	UserGid *int32 `json:"user_gid,omitempty"`

	// 用户名称
	UserName *string `json:"user_name,omitempty"`

	// 用户组名称
	UserGroupName *string `json:"user_group_name,omitempty"`

	// 用户home目录
	UserHomeDir *string `json:"user_home_dir,omitempty"`

	// 用户登录ip
	LoginIp *string `json:"login_ip,omitempty"`

	// 服务类型，包含如下:   - system   - mysql   - redis
	ServiceType *string `json:"service_type,omitempty"`

	// 登录服务端口
	ServicePort *int32 `json:"service_port,omitempty"`

	// 登录方式
	LoginMode *int32 `json:"login_mode,omitempty"`

	// 用户最后一次登录时间
	LoginLastTime *int64 `json:"login_last_time,omitempty"`

	// 用户登录失败次数
	LoginFailCount *int32 `json:"login_fail_count,omitempty"`

	// 口令hash
	PwdHash *string `json:"pwd_hash,omitempty"`

	// 匿名化处理后的口令
	PwdWithFuzzing *string `json:"pwd_with_fuzzing,omitempty"`

	// 密码使用的天数
	PwdUsedDays *int32 `json:"pwd_used_days,omitempty"`

	// 口令的最短有效期限
	PwdMinDays *int32 `json:"pwd_min_days,omitempty"`

	// 口令的最长有效期限
	PwdMaxDays *int32 `json:"pwd_max_days,omitempty"`

	// 口令无效时提前告警天数
	PwdWarnLeftDays *int32 `json:"pwd_warn_left_days,omitempty"`
}

func (o EventUserResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EventUserResponseInfo struct{}"
	}

	return strings.Join([]string{"EventUserResponseInfo", string(data)}, " ")
}
