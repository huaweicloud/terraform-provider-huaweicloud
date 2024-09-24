package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// KafkaMessageDiagnosisReportInfoEntity 消息积压诊断报告信息实体类
type KafkaMessageDiagnosisReportInfoEntity struct {

	// 诊断报告ID
	ReportId string `json:"report_id"`

	// 消息积压诊断任务状态。 - diagnosing：诊断中 - failed：诊断失败 - deleted：手动删除 - finished：诊断完成 - normal：诊断结果正常 - abnormal：诊断结果异常
	Status KafkaMessageDiagnosisReportInfoEntityStatus `json:"status"`

	// 诊断任务开始时间
	BeginTime string `json:"begin_time"`

	// 诊断任务结束时间
	EndTime *string `json:"end_time,omitempty"`

	// 该次诊断任务诊断的消费组名称
	GroupName string `json:"group_name"`

	// 该次诊断任务诊断的topic名称
	TopicName string `json:"topic_name"`

	// 该次诊断任务发现的存在消息堆积的分区数
	AccumulatedPartitions int32 `json:"accumulated_partitions"`
}

func (o KafkaMessageDiagnosisReportInfoEntity) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KafkaMessageDiagnosisReportInfoEntity struct{}"
	}

	return strings.Join([]string{"KafkaMessageDiagnosisReportInfoEntity", string(data)}, " ")
}

type KafkaMessageDiagnosisReportInfoEntityStatus struct {
	value string
}

type KafkaMessageDiagnosisReportInfoEntityStatusEnum struct {
	DIAGNOSING KafkaMessageDiagnosisReportInfoEntityStatus
	FAILED     KafkaMessageDiagnosisReportInfoEntityStatus
	DELETED    KafkaMessageDiagnosisReportInfoEntityStatus
	FINISHED   KafkaMessageDiagnosisReportInfoEntityStatus
	NORMAL     KafkaMessageDiagnosisReportInfoEntityStatus
	ABNORMAL   KafkaMessageDiagnosisReportInfoEntityStatus
}

func GetKafkaMessageDiagnosisReportInfoEntityStatusEnum() KafkaMessageDiagnosisReportInfoEntityStatusEnum {
	return KafkaMessageDiagnosisReportInfoEntityStatusEnum{
		DIAGNOSING: KafkaMessageDiagnosisReportInfoEntityStatus{
			value: "diagnosing",
		},
		FAILED: KafkaMessageDiagnosisReportInfoEntityStatus{
			value: "failed",
		},
		DELETED: KafkaMessageDiagnosisReportInfoEntityStatus{
			value: "deleted",
		},
		FINISHED: KafkaMessageDiagnosisReportInfoEntityStatus{
			value: "finished",
		},
		NORMAL: KafkaMessageDiagnosisReportInfoEntityStatus{
			value: "normal",
		},
		ABNORMAL: KafkaMessageDiagnosisReportInfoEntityStatus{
			value: "abnormal",
		},
	}
}

func (c KafkaMessageDiagnosisReportInfoEntityStatus) Value() string {
	return c.value
}

func (c KafkaMessageDiagnosisReportInfoEntityStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *KafkaMessageDiagnosisReportInfoEntityStatus) UnmarshalJSON(b []byte) error {
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
