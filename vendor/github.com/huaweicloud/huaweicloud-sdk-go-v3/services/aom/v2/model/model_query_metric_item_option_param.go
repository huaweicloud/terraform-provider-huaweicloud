package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 参数项。
type QueryMetricItemOptionParam struct {

	// 指标维度列表。
	Dimensions *[]Dimension `json:"dimensions,omitempty"`

	// 指标名称。名称长度取值范围为1~255个字符。 取值范围： AOM提供的基础指标， cpuUsage、cpuCoreUsed等 cpuUage：cpu使用率； cpuCoreUsed：cpu内核占用； 用户上报的自定义指标名称。
	MetricName *string `json:"metricName,omitempty"`

	// 指标命名空间。 取值范围 PAAS.CONTAINER：组件指标、实例指标、进程指标和容器指标的命名空间， PAAS.NODE： 主机指标、网络指标、磁盘指标和文件系统指标的命名空间， PAAS.SLA：SLA指标的命名空间， PAAS.AGGR：集群指标的命名空间， CUSTOMMETRICS：默认的自定义指标的命名空间。
	Namespace QueryMetricItemOptionParamNamespace `json:"namespace"`
}

func (o QueryMetricItemOptionParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "QueryMetricItemOptionParam struct{}"
	}

	return strings.Join([]string{"QueryMetricItemOptionParam", string(data)}, " ")
}

type QueryMetricItemOptionParamNamespace struct {
	value string
}

type QueryMetricItemOptionParamNamespaceEnum struct {
	PAAS_CONTAINER QueryMetricItemOptionParamNamespace
	PAAS_NODE      QueryMetricItemOptionParamNamespace
	PAAS_SLA       QueryMetricItemOptionParamNamespace
	PAAS_AGGR      QueryMetricItemOptionParamNamespace
	CUSTOMMETRICS  QueryMetricItemOptionParamNamespace
}

func GetQueryMetricItemOptionParamNamespaceEnum() QueryMetricItemOptionParamNamespaceEnum {
	return QueryMetricItemOptionParamNamespaceEnum{
		PAAS_CONTAINER: QueryMetricItemOptionParamNamespace{
			value: "PAAS.CONTAINER",
		},
		PAAS_NODE: QueryMetricItemOptionParamNamespace{
			value: "PAAS.NODE",
		},
		PAAS_SLA: QueryMetricItemOptionParamNamespace{
			value: "PAAS.SLA",
		},
		PAAS_AGGR: QueryMetricItemOptionParamNamespace{
			value: "PAAS.AGGR",
		},
		CUSTOMMETRICS: QueryMetricItemOptionParamNamespace{
			value: "CUSTOMMETRICS",
		},
	}
}

func (c QueryMetricItemOptionParamNamespace) Value() string {
	return c.value
}

func (c QueryMetricItemOptionParamNamespace) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *QueryMetricItemOptionParamNamespace) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(string)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to string error")
	}
}
