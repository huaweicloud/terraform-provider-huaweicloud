package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AccessAreaFilter 区域访问控制。   > - 使用该功能需要提交工单开通区域访问控制功能。   > - CDN会定期更新IP地址库，部分未在地址库的IP将无法识别到所属位置。如果CDN无法识别用户所在位置，将采取放行策略，返回对应的资源给用户。
type AccessAreaFilter struct {

	// 规则类型，黑、白名单二选一。   - black: 黑名单，如果匹配到黑名单规则，则黑名单所选区域内的用户将无法访问当前资源，返回403状态码。   - white: 白名单，白名单所选区域以外的用户均无法访问当前资源，返回403状态码。
	Type *string `json:"type,omitempty"`

	// 生效类型。   - all: 所有文件，所有文件均遵循配置的规则。   - file_directory: 目录路径，指定目录路径的资源遵循配置的规则。   - file_path: 全路径，指定路径的资源遵循配置的规则。
	ContentType *string `json:"content_type,omitempty"`

	// 生效规则。当content_type为all时，为空或不传。 当content_type为file_directory时，输入要求以“/”作为首字符，多个目录以“,”进行分隔，如/test/folder01,/test/folder02，并且输入的目录路径总数不超过100个。 当content_type为file_path时，输入要求以“/”或“\\*”作为首字符，支持配置通配符“\\*”，通配符不能连续出现且不能超过两个。多个路径以“,”进行分割，如/test/a.txt,/test/b.txt，并且输出的总数不能超过100个。   > - 不允许配置两条完全一样的白名单或黑名单规则。   > - 仅允许配置一条生效类型为“所有文件”的规则。
	ContentValue *string `json:"content_value,omitempty"`

	// 配置规则适用的区域，多个区域以“,”进行分隔，支持的区域如：CN_IN：中国大陆，AF：阿富汗，IE：爱尔兰，EG：埃及，AU：澳大利亚等。具体的位置编码参见《附录-地理位置编码》查询。
	Area *string `json:"area,omitempty"`

	// 例外IP，配置指定IP不执行当前规则。
	ExceptionIp *string `json:"exception_ip,omitempty"`
}

func (o AccessAreaFilter) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AccessAreaFilter struct{}"
	}

	return strings.Join([]string{"AccessAreaFilter", string(data)}, " ")
}
