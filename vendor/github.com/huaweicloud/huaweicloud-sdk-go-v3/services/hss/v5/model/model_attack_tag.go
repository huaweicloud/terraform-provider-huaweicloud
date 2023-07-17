package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AttackTag 攻击标识，包含如下：   - attack_success : 攻击成功   - attack_attempt : 攻击尝试   - attack_blocked : 攻击被阻断   - abnormal_behavior : 异常行为   - collapsible_host : 主机失陷   - system_vulnerability : 系统脆弱性
type AttackTag struct {
}

func (o AttackTag) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AttackTag struct{}"
	}

	return strings.Join([]string{"AttackTag", string(data)}, " ")
}
