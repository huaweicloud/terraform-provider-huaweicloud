package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListInstancesResourceMetricsRequest Request Object
type ListInstancesResourceMetricsRequest struct {

	// 引擎类型
	Engine ListInstancesResourceMetricsRequestEngine `json:"engine"`

	// 搜索字段
	SearchField *string `json:"search_field,omitempty"`

	// 索引位置，偏移量
	Offset *string `json:"offset,omitempty"`

	// 查询数据条数
	Limit *string `json:"limit,omitempty"`

	// 排序方式
	Order *ListInstancesResourceMetricsRequestOrder `json:"order,omitempty"`

	// 排序字段
	SortField *ListInstancesResourceMetricsRequestSortField `json:"sort_field,omitempty"`
}

func (o ListInstancesResourceMetricsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstancesResourceMetricsRequest struct{}"
	}

	return strings.Join([]string{"ListInstancesResourceMetricsRequest", string(data)}, " ")
}

type ListInstancesResourceMetricsRequestEngine struct {
	value string
}

type ListInstancesResourceMetricsRequestEngineEnum struct {
	MYSQL      ListInstancesResourceMetricsRequestEngine
	POSTGRESQL ListInstancesResourceMetricsRequestEngine
	SQLSERVER  ListInstancesResourceMetricsRequestEngine
}

func GetListInstancesResourceMetricsRequestEngineEnum() ListInstancesResourceMetricsRequestEngineEnum {
	return ListInstancesResourceMetricsRequestEngineEnum{
		MYSQL: ListInstancesResourceMetricsRequestEngine{
			value: "mysql",
		},
		POSTGRESQL: ListInstancesResourceMetricsRequestEngine{
			value: "postgresql",
		},
		SQLSERVER: ListInstancesResourceMetricsRequestEngine{
			value: "sqlserver",
		},
	}
}

func (c ListInstancesResourceMetricsRequestEngine) Value() string {
	return c.value
}

func (c ListInstancesResourceMetricsRequestEngine) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListInstancesResourceMetricsRequestEngine) UnmarshalJSON(b []byte) error {
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

type ListInstancesResourceMetricsRequestOrder struct {
	value string
}

type ListInstancesResourceMetricsRequestOrderEnum struct {
	DESC ListInstancesResourceMetricsRequestOrder
	ASC  ListInstancesResourceMetricsRequestOrder
}

func GetListInstancesResourceMetricsRequestOrderEnum() ListInstancesResourceMetricsRequestOrderEnum {
	return ListInstancesResourceMetricsRequestOrderEnum{
		DESC: ListInstancesResourceMetricsRequestOrder{
			value: "DESC",
		},
		ASC: ListInstancesResourceMetricsRequestOrder{
			value: "ASC",
		},
	}
}

func (c ListInstancesResourceMetricsRequestOrder) Value() string {
	return c.value
}

func (c ListInstancesResourceMetricsRequestOrder) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListInstancesResourceMetricsRequestOrder) UnmarshalJSON(b []byte) error {
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

type ListInstancesResourceMetricsRequestSortField struct {
	value string
}

type ListInstancesResourceMetricsRequestSortFieldEnum struct {
	INSTANCE_NAME      ListInstancesResourceMetricsRequestSortField
	STATUS             ListInstancesResourceMetricsRequestSortField
	TYPE               ListInstancesResourceMetricsRequestSortField
	CPU_USAGE          ListInstancesResourceMetricsRequestSortField
	MEMORY_USAGE       ListInstancesResourceMetricsRequestSortField
	DISK_USAGE         ListInstancesResourceMetricsRequestSortField
	TPS                ListInstancesResourceMetricsRequestSortField
	QPS                ListInstancesResourceMetricsRequestSortField
	IOPS               ListInstancesResourceMetricsRequestSortField
	ACTIVE_CONNECTIONS ListInstancesResourceMetricsRequestSortField
	SLOW_SQL           ListInstancesResourceMetricsRequestSortField
}

func GetListInstancesResourceMetricsRequestSortFieldEnum() ListInstancesResourceMetricsRequestSortFieldEnum {
	return ListInstancesResourceMetricsRequestSortFieldEnum{
		INSTANCE_NAME: ListInstancesResourceMetricsRequestSortField{
			value: "instance_name",
		},
		STATUS: ListInstancesResourceMetricsRequestSortField{
			value: "status",
		},
		TYPE: ListInstancesResourceMetricsRequestSortField{
			value: "type",
		},
		CPU_USAGE: ListInstancesResourceMetricsRequestSortField{
			value: "cpu_usage",
		},
		MEMORY_USAGE: ListInstancesResourceMetricsRequestSortField{
			value: "memory_usage",
		},
		DISK_USAGE: ListInstancesResourceMetricsRequestSortField{
			value: "disk_usage",
		},
		TPS: ListInstancesResourceMetricsRequestSortField{
			value: "tps",
		},
		QPS: ListInstancesResourceMetricsRequestSortField{
			value: "qps",
		},
		IOPS: ListInstancesResourceMetricsRequestSortField{
			value: "iops",
		},
		ACTIVE_CONNECTIONS: ListInstancesResourceMetricsRequestSortField{
			value: "active_connections",
		},
		SLOW_SQL: ListInstancesResourceMetricsRequestSortField{
			value: "slow_sql",
		},
	}
}

func (c ListInstancesResourceMetricsRequestSortField) Value() string {
	return c.value
}

func (c ListInstancesResourceMetricsRequestSortField) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListInstancesResourceMetricsRequestSortField) UnmarshalJSON(b []byte) error {
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
