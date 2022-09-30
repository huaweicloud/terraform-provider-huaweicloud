package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SetRdsBackupCnfReq struct {

	// 自动创建快照的名称前缀，需要用户自己手动输入。只能包含1~32位小写字母、数字、中划线或者下划线，并且以小写字母开头。
	Prefix string `json:"prefix"`

	// 每天创建快照的时刻，只支持整点，后面需加上时区，格式为“HH:mm z”，“HH:mm”表示整点时间，“z”表示时区。比如“00:00 GMT+08:00”、“01:00 GMT+08:00”等。
	Period string `json:"period"`

	// 自定义设置快照保留的天数，范围是1～90。系统在半点时刻会自动删除超过保留天数的快照。
	Keepday int32 `json:"keepday"`

	// 是否开启自动创建快照策略。 - true：表示开启自动创建快照策略。 - false：表示关闭自动创建快照策略。
	Enable string `json:"enable"`

	// 表示关闭自动创建快照策略时，是否需要清除所有自动创建的快照。默认为“false”，表示不会删除之前已自动创建的快照。设置为true，表示在关闭自动创建快照策略的同时，删除所有已创建的快照。
	DeleteAuto *string `json:"deleteAuto,omitempty"`
}

func (o SetRdsBackupCnfReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetRdsBackupCnfReq struct{}"
	}

	return strings.Join([]string{"SetRdsBackupCnfReq", string(data)}, " ")
}
