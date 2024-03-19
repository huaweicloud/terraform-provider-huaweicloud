package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ManualVulScanRequestInfo 手工检测漏洞接口请求体
type ManualVulScanRequestInfo struct {

	// 操作类型,包含如下：   -linux_vul : linux漏洞   -windows_vul : windows漏洞   -web_cms : Web-CMS漏洞   -app_vul : 应用漏洞   -urgent_vul : 应急漏洞
	ManualScanType *[]ManualVulScanRequestInfoManualScanType `json:"manual_scan_type,omitempty"`

	// 是否是批量操作,为true时扫描所有支持的主机
	BatchFlag *bool `json:"batch_flag,omitempty"`

	// 扫描主机的范围，包含如下：   -all_host : 扫描全部主机,此类型不需要填写agent_id_list   -specific_host : 扫描指定主机
	RangeType *string `json:"range_type,omitempty"`

	// 主机列表
	AgentIdList *[]string `json:"agent_id_list,omitempty"`

	// 扫描的应急漏洞id列表，若为空则扫描所有应急漏洞 包含如下： \"URGENT-CVE-2023-46604   Apache ActiveMQ远程代码执行漏洞\" \"URGENT-HSSVD-2020-1109  Elasticsearch 未授权访问漏洞\" \"URGENT-CVE-2022-26134   Atlassian Confluence OGNL 远程代码执行漏洞（CVE-2022-26134）\" \"URGENT-CVE-2023-22515   Atlassian Confluence Data Center and Server 权限提升漏洞(CVE-2023-22515)\" \"URGENT-CVE-2023-22518   Atlassian Confluence Data Center & Server 授权机制不恰当漏洞(CVE-2023-22518)\" \"URGENT-CVE-2023-28432   MinIO 信息泄露漏洞（CVE-2023-28432）\" \"URGENT-CVE-2023-37582   Apache RocketMQ 远程代码执行漏洞(CVE-2023-37582)\" \"URGENT-CVE-2023-33246   Apache RocketMQ 远程代码执行漏洞(CVE-2023-33246)\" \"URGENT-CNVD-2023-02709  禅道项目管理系统远程命令执行漏洞(CNVD-2023-02709)\" \"URGENT-CVE-2022-36804   Atlassian Bitbucket Server 和 Data Center 命令注入漏洞(CVE-2022-36804)\" \"URGENT-CVE-2022-22965   Spring Framework JDK >= 9 远程代码执行漏洞\" \"URGENT-CVE-2022-25845   fastjson <1.2.83 远程代码执行漏洞\" \"URGENT-CVE-2019-14439   Jackson-databind远程命令执行漏洞（CVE-2019-14439）\" \"URGENT-CVE-2020-13933   Apache Shiro身份验证绕过漏洞（CVE-2020-13933）\" \"URGENT-CVE-2020-26217   XStream < 1.4.14 远程代码执行漏洞（CVE-2020-26217）\" \"URGENT-CVE-2021-4034    Linux Polkit 权限提升漏洞预警（CVE-2021-4034）\" \"URGENT-CVE-2021-44228   Apache Log4j2 远程代码执行漏洞（CVE-2021-44228、CVE-2021-45046）\" \"URGENT-CVE-2022-0847    Dirty Pipe - Linux 内核本地提权漏洞（CVE-2022-0847）\"
	UrgentVulIdList *[]string `json:"urgent_vul_id_list,omitempty"`
}

func (o ManualVulScanRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ManualVulScanRequestInfo struct{}"
	}

	return strings.Join([]string{"ManualVulScanRequestInfo", string(data)}, " ")
}

type ManualVulScanRequestInfoManualScanType struct {
	value string
}

type ManualVulScanRequestInfoManualScanTypeEnum struct {
	LINUX_VUL   ManualVulScanRequestInfoManualScanType
	WINDOWS_VUL ManualVulScanRequestInfoManualScanType
	WEB_CMS     ManualVulScanRequestInfoManualScanType
	APP_VUL     ManualVulScanRequestInfoManualScanType
	URGENT_VUL  ManualVulScanRequestInfoManualScanType
}

func GetManualVulScanRequestInfoManualScanTypeEnum() ManualVulScanRequestInfoManualScanTypeEnum {
	return ManualVulScanRequestInfoManualScanTypeEnum{
		LINUX_VUL: ManualVulScanRequestInfoManualScanType{
			value: "linux_vul",
		},
		WINDOWS_VUL: ManualVulScanRequestInfoManualScanType{
			value: "windows_vul",
		},
		WEB_CMS: ManualVulScanRequestInfoManualScanType{
			value: "web_cms",
		},
		APP_VUL: ManualVulScanRequestInfoManualScanType{
			value: "app_vul",
		},
		URGENT_VUL: ManualVulScanRequestInfoManualScanType{
			value: "urgent_vul",
		},
	}
}

func (c ManualVulScanRequestInfoManualScanType) Value() string {
	return c.value
}

func (c ManualVulScanRequestInfoManualScanType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ManualVulScanRequestInfoManualScanType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
