package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AttackPhase 攻击阶段，包含如下：   - reconnaissance : 侦查跟踪   - weaponization : 武器构建   - delivery : 载荷投递   - exploit : 漏洞利用   - installation : 安装植入   - command_and_control : 命令与控制   - actions : 目标达成
type AttackPhase struct {
}

func (o AttackPhase) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AttackPhase struct{}"
	}

	return strings.Join([]string{"AttackPhase", string(data)}, " ")
}
