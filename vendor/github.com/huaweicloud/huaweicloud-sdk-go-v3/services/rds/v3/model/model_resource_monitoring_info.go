package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ResourceMonitoringInfo 资源监控信息
type ResourceMonitoringInfo struct {

	// 实例id
	Id string `json:"id"`

	// 实例名称
	Name string `json:"name"`

	InstanceState *InstanceState `json:"instance_state"`

	// 实例类型
	Type ResourceMonitoringInfoType `json:"type"`

	// cpu大小
	Cpu string `json:"cpu"`

	// 内存大小（单位：GB）
	Mem string `json:"mem"`

	// 引擎名称
	EngineName ResourceMonitoringInfoEngineName `json:"engine_name"`

	// 引擎版本
	EngineVersion string `json:"engine_version"`

	// cpu使用率
	CpuUsage float64 `json:"cpu_usage"`

	// 内存使用率
	MemoryUsage float64 `json:"memory_usage"`

	// 磁盘使用率
	DiskUsage float64 `json:"disk_usage"`

	// tps
	Tps float64 `json:"tps"`

	// qps
	Qps *float64 `json:"qps,omitempty"`

	// iops
	Iops float64 `json:"iops"`

	// 活跃连接数
	ActiveConnections int32 `json:"active_connections"`

	// 慢SQL
	SlowSql *float64 `json:"slow_sql,omitempty"`

	// 只读实例资源监控指标
	ReadonlyInstanceResourceMonitoringInfo []ResourceMonitoringInfo `json:"readonly_instance_resource_monitoring_info"`
}

func (o ResourceMonitoringInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResourceMonitoringInfo struct{}"
	}

	return strings.Join([]string{"ResourceMonitoringInfo", string(data)}, " ")
}

type ResourceMonitoringInfoType struct {
	value string
}

type ResourceMonitoringInfoTypeEnum struct {
	SINGLE     ResourceMonitoringInfoType
	HA         ResourceMonitoringInfoType
	REPLICA    ResourceMonitoringInfoType
	ENTERPRISE ResourceMonitoringInfoType
}

func GetResourceMonitoringInfoTypeEnum() ResourceMonitoringInfoTypeEnum {
	return ResourceMonitoringInfoTypeEnum{
		SINGLE: ResourceMonitoringInfoType{
			value: "Single",
		},
		HA: ResourceMonitoringInfoType{
			value: "Ha",
		},
		REPLICA: ResourceMonitoringInfoType{
			value: "Replica",
		},
		ENTERPRISE: ResourceMonitoringInfoType{
			value: "Enterprise",
		},
	}
}

func (c ResourceMonitoringInfoType) Value() string {
	return c.value
}

func (c ResourceMonitoringInfoType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ResourceMonitoringInfoType) UnmarshalJSON(b []byte) error {
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

type ResourceMonitoringInfoEngineName struct {
	value string
}

type ResourceMonitoringInfoEngineNameEnum struct {
	MYSQL      ResourceMonitoringInfoEngineName
	POSTGRESQL ResourceMonitoringInfoEngineName
	SQLSERVER  ResourceMonitoringInfoEngineName
}

func GetResourceMonitoringInfoEngineNameEnum() ResourceMonitoringInfoEngineNameEnum {
	return ResourceMonitoringInfoEngineNameEnum{
		MYSQL: ResourceMonitoringInfoEngineName{
			value: "mysql",
		},
		POSTGRESQL: ResourceMonitoringInfoEngineName{
			value: "postgresql",
		},
		SQLSERVER: ResourceMonitoringInfoEngineName{
			value: "sqlserver",
		},
	}
}

func (c ResourceMonitoringInfoEngineName) Value() string {
	return c.value
}

func (c ResourceMonitoringInfoEngineName) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ResourceMonitoringInfoEngineName) UnmarshalJSON(b []byte) error {
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
