package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MysqlInstanceResponse struct {
	// 实例ID。

	Id string `json:"id"`
	// 实例名称。用于表示实例的名称，同一租户下，同类型的实例名称可相同。 取值范围：4~64个字符之间，必须以字母开头，不区分大小写，可以包含字母、数字、中划线或者下划线, 不能包含其它的特殊字符。

	Name string `json:"name"`
	// 实例状态。

	Status *string `json:"status,omitempty"`

	Datastore *MysqlDatastore `json:"datastore,omitempty"`
	// 实例类型，仅支持Cluster。

	Mode *string `json:"mode,omitempty"`
	// 参数组ID。

	ConfigurationId *string `json:"configuration_id,omitempty"`
	// 数据库端口信息。

	Port *string `json:"port,omitempty"`

	BackupStrategy *MysqlBackupStrategy `json:"backup_strategy,omitempty"`
	// 企业项目ID。

	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
	// 区域ID，与请求参数相同。

	Region *string `json:"region,omitempty"`
	// 可用区模式，与请求参数相同。

	AvailabilityZoneMode *string `json:"availability_zone_mode,omitempty"`
	// 主可用区ID。

	MasterAvailabilityZone *string `json:"master_availability_zone,omitempty"`
	// 虚拟私有云ID，与请求参数相同。

	VpcId *string `json:"vpc_id,omitempty"`
	// 安全组ID，与请求参数相同。

	SecurityGroupId *string `json:"security_group_id,omitempty"`
	// 子网ID，与请求参数相同。

	SubnetId *string `json:"subnet_id,omitempty"`
	// 规格码，与请求参数相同。

	FlavorRef *string `json:"flavor_ref,omitempty"`

	ChargeInfo *MysqlChargeInfo `json:"charge_info,omitempty"`
}

func (o MysqlInstanceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlInstanceResponse struct{}"
	}

	return strings.Join([]string{"MysqlInstanceResponse", string(data)}, " ")
}
