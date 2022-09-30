package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 集群对象。
type CreateClusterBody struct {
	Instance *CreateClusterInstanceBody `json:"instance"`

	Datastore *CreateClusterDatastoreBody `json:"datastore"`

	// 集群名称。4～32个字符，只能包含数字、字母、中划线和下划线，且必须以字母开头。
	Name string `json:"name"`

	// 集群实例个数，取值范围为1~32。
	InstanceNum int32 `json:"instanceNum"`

	BackupStrategy *CreateClusterBackupStrategyBody `json:"backupStrategy,omitempty"`

	// 设置是否进行通信加密。取值范围为true或false。默认关闭通信加密功能。当httpsEnable设置为true时，authorityEnable字段需要设置为true。  - true：表示集群进行通信加密。 - false：表示集群不进行通信加密。  >此参数只有6.5.4及之后版本支持。
	HttpsEnable *bool `json:"httpsEnable,omitempty"`

	// 是否开启认证，取值范围为true或false。默认关闭认证功能。  - true：表示集群开启认证。 - false：表示集群不开启认证。  >此参数只有6.5.4及之后版本支持。
	AuthorityEnable *bool `json:"authorityEnable,omitempty"`

	// 安全模式下集群管理员admin的密码，只有在创建集群时authorityEnable设置为true时需要设置此参数。   - 管理员密码需要满足规则：    - 可输入的字符串长度为8-32个字符。    - 密码至少包含大写字母，小写字母，数字和特殊字符中的三类，其中可输入的特殊字符为：~!@#$%^&*()-_=+\\\\|[{}];:,<.>/?。   - 安全集群的密码会进行弱口令校验，建议设置安全性高的密码。
	AdminPwd *string `json:"adminPwd,omitempty"`

	// 企业项目ID。创建集群时，给集群绑定企业项目ID。最大长度36个字符，带\"-\"连字符的UUID格式，或者是字符串\"0\"。\"0\"表示默认企业项目。  关于企业项目ID的获取及企业项目特性的详细信息，请参见[[《企业管理服务用户指南》](https://support.huaweicloud.com/usermanual-em/zh-cn_topic_0123692049.html)](tag:hc)[[《企业管理服务用户指南》](https://support.huaweicloud.com/intl/zh-cn/usermanual-em/zh-cn_topic_0123692049.html)](tag:hk)。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 集群标签。 关于标签特性的详细信息，请参见[[《标签管理服务介绍》](https://support.huaweicloud.com/productdesc-tms/zh-cn_topic_0071335169.html)](tag:hc)[[《标签管理服务介绍》](https://support.huaweicloud.com/intl/zh-cn/productdesc-tms/zh-cn_topic_0071335169.html)](tag:hk)。
	Tags *[]CreateClusterTagsBody `json:"tags,omitempty"`

	PayInfo *PayInfoBody `json:"payInfo,omitempty"`
}

func (o CreateClusterBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClusterBody struct{}"
	}

	return strings.Join([]string{"CreateClusterBody", string(data)}, " ")
}
