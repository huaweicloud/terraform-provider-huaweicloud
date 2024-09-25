package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type Traces struct {

	// 标识事件对应的云服务资源ID。
	ResourceId *string `json:"resource_id,omitempty"`

	// 标识查询事件列表对应的事件名称。由0-9,a-z,A-Z,'-','.','_',组成，长度为1～64个字符，且以首字符必须为字母。
	TraceName *string `json:"trace_name,omitempty"`

	// 标识事件等级，目前有三种：正常（normal），警告（warning），事故（incident）。
	TraceRating *TracesTraceRating `json:"trace_rating,omitempty"`

	// 标识事件发生源头类型，管理类事件主要包括API调用（ApiCall），Console页面调用（ConsoleAction）和系统间调用（SystemAction）。 数据类事件主要包括ObsSDK，ObsAPI。
	TraceType *string `json:"trace_type,omitempty"`

	// 标识事件对应接口请求内容，即资源操作请求体。
	Request *string `json:"request,omitempty"`

	// 记录用户请求的响应，标识事件对应接口响应内容，即资源操作结果返回体。
	Response *string `json:"response,omitempty"`

	// 记录用户请求的响应，标识事件对应接口返回的HTTP状态码。
	Code *string `json:"code,omitempty"`

	// 标识事件对应的云服务接口版本。
	ApiVersion *string `json:"api_version,omitempty"`

	// 标识其他云服务为此条事件添加的备注信息。
	Message *string `json:"message,omitempty"`

	// 标识云审计服务记录本次事件的时间戳。
	RecordTime *int64 `json:"record_time,omitempty"`

	// 标识事件的ID，由系统生成的UUID。
	TraceId *string `json:"trace_id,omitempty"`

	// 标识事件产生的时间戳。
	Time *int64 `json:"time,omitempty"`

	User *UserInfo `json:"user,omitempty"`

	// 标识查询事件列表对应的云服务类型。必须为已对接CTS的云服务的英文缩写，且服务类型一般为大写字母。
	ServiceType *string `json:"service_type,omitempty"`

	// 查询事件列表对应的资源类型。
	ResourceType *string `json:"resource_type,omitempty"`

	// 标识触发事件的租户IP。
	SourceIp *string `json:"source_ip,omitempty"`

	// 标识事件对应的资源名称。
	ResourceName *string `json:"resource_name,omitempty"`

	// 记录本次请求的request id
	RequestId *string `json:"request_id,omitempty"`

	// 记录本次请求出错后，问题定位所需要的辅助信息。
	LocationInfo *string `json:"location_info,omitempty"`

	// 云资源的详情页面
	Endpoint *string `json:"endpoint,omitempty"`

	// 云资源的详情页面的访问链接（不含endpoint）
	ResourceUrl *string `json:"resource_url,omitempty"`

	// 标识资源所在的企业项目ID。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 标识资源所在的账号ID。仅在跨租户操作资源时有值。
	ResourceAccountId *string `json:"resource_account_id,omitempty"`
}

func (o Traces) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Traces struct{}"
	}

	return strings.Join([]string{"Traces", string(data)}, " ")
}

type TracesTraceRating struct {
	value string
}

type TracesTraceRatingEnum struct {
	NORMAL   TracesTraceRating
	WARNING  TracesTraceRating
	INCIDENT TracesTraceRating
}

func GetTracesTraceRatingEnum() TracesTraceRatingEnum {
	return TracesTraceRatingEnum{
		NORMAL: TracesTraceRating{
			value: "normal",
		},
		WARNING: TracesTraceRating{
			value: "warning",
		},
		INCIDENT: TracesTraceRating{
			value: "incident",
		},
	}
}

func (c TracesTraceRating) Value() string {
	return c.value
}

func (c TracesTraceRating) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TracesTraceRating) UnmarshalJSON(b []byte) error {
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
