package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// StartLogsRequest Request Object
type StartLogsRequest struct {

	// 指定开启日志的集群ID。
	ClusterId string `json:"cluster_id"`

	// action支持base_log_collect和real_time_log_collect两种，base就是之前历史的能力，real_time为实时采集能力，默认不传就是base，兼容之前的逻辑
	Action *StartLogsRequestAction `json:"action,omitempty"`

	Body *StartLogsReq `json:"body,omitempty"`
}

func (o StartLogsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartLogsRequest struct{}"
	}

	return strings.Join([]string{"StartLogsRequest", string(data)}, " ")
}

type StartLogsRequestAction struct {
	value string
}

type StartLogsRequestActionEnum struct {
	BASE_LOG_COLLECT      StartLogsRequestAction
	REAL_TIME_LOG_COLLECT StartLogsRequestAction
}

func GetStartLogsRequestActionEnum() StartLogsRequestActionEnum {
	return StartLogsRequestActionEnum{
		BASE_LOG_COLLECT: StartLogsRequestAction{
			value: "base_log_collect",
		},
		REAL_TIME_LOG_COLLECT: StartLogsRequestAction{
			value: "real_time_log_collect",
		},
	}
}

func (c StartLogsRequestAction) Value() string {
	return c.value
}

func (c StartLogsRequestAction) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *StartLogsRequestAction) UnmarshalJSON(b []byte) error {
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
