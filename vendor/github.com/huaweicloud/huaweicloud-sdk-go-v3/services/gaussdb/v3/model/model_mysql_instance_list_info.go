package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type MysqlInstanceListInfo struct {
	// 实例ID。

	Id string `json:"id"`
	// 创建的实例名称。

	Name string `json:"name"`
	// 实例状态。

	Status *string `json:"status,omitempty"`
	// 实例写内网IP地址列表。弹性云服务器创建成功后该值存在，其他情况下为空字符串。

	PrivateIps *[]string `json:"private_ips,omitempty"`
	// 实例外网IP地址列表。

	PublicIps *[]string `json:"public_ips,omitempty"`
	// 数据库端口号。

	Port *string `json:"port,omitempty"`
	// 实例类型，取值为“Cluster”。

	Type *string `json:"type,omitempty"`
	// 实例所在区域。

	Region *string `json:"region,omitempty"`

	Datastore *MysqlDatastore `json:"datastore,omitempty"`
	// 创建时间，格式为\"yyyy-mm-ddThh:mm:ssZ\"。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。说明：创建时返回值为空，数据库实例创建成功后该值不为空。

	Created *string `json:"created,omitempty"`
	// 更新时间，格式与\"created\"字段对应格式完全相同。说明：创建时返回值为空，数据库实例创建成功后该值不为空。

	Updated *string `json:"updated,omitempty"`
	// 默认用户名。

	DbUserName *string `json:"db_user_name,omitempty"`
	// 虚拟私有云ID。

	VpcId *string `json:"vpc_id,omitempty"`
	// 子网的网络ID信息。

	SubnetId *string `json:"subnet_id,omitempty"`
	// 安全组ID。

	SecurityGroupId *string `json:"security_group_id,omitempty"`
	// 规格码。

	FlavorRef *string `json:"flavor_ref,omitempty"`

	FlavorInfo *MysqlFlavorInfo `json:"flavor_info,omitempty"`

	Volume *MysqlVolumeInfo `json:"volume,omitempty"`

	BackupStrategy *MysqlBackupStrategy `json:"backup_strategy,omitempty"`
	// 企业项目ID。

	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`
	// 时区。

	TimeZone *string `json:"time_zone,omitempty"`

	ChargeInfo *MysqlChargeInfo `json:"charge_info,omitempty"`
	// 专属资源池ID，只有数据库实例属于专属资源池才会返回该参数。

	DedicatedResourceId *string `json:"dedicated_resource_id,omitempty"`
	// 标签列表。

	Tags *[]InstanceTagItem `json:"tags,omitempty"`
}

func (o MysqlInstanceListInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlInstanceListInfo struct{}"
	}

	return strings.Join([]string{"MysqlInstanceListInfo", string(data)}, " ")
}
