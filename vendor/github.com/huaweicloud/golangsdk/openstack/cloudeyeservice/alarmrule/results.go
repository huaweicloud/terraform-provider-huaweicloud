package alarmrule

import (
	"fmt"
	"log"

	"github.com/huaweicloud/golangsdk"
)

type CreateResponse struct {
	AlarmID string `json:"alarm_id"`
}

type CreateResult struct {
	golangsdk.Result
}

func (c CreateResult) Extract() (*CreateResponse, error) {
	r := &CreateResponse{}
	return r, c.ExtractInto(r)
}

type DimensionInfo struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type MetricInfo struct {
	Namespace  string          `json:"namespace"`
	MetricName string          `json:"metric_name"`
	Dimensions []DimensionInfo `json:"dimensions"`
}

type ConditionInfo struct {
	Period             int    `json:"period"`
	Filter             string `json:"filter"`
	ComparisonOperator string `json:"comparison_operator"`
	Value              int    `json:"value"`
	Unit               string `json:"unit"`
	Count              int    `json:"count"`
}

type ActionInfo struct {
	Type             string   `json:"type"`
	NotificationList []string `json:"notification_list"`
}

type AlarmRule struct {
	AlarmName               string        `json:"alarm_name"`
	AlarmDescription        string        `json:"alarm_description"`
	Metric                  MetricInfo    `json:"metric"`
	Condition               ConditionInfo `json:"condition"`
	AlarmActions            []ActionInfo  `json:"alarm_actions"`
	InsufficientdataActions []ActionInfo  `json:"insufficientdata_actions"`
	OkActions               []ActionInfo  `json:"ok_actions"`
	AlarmEnabled            bool          `json:"alarm_enabled"`
	AlarmActionEnabled      bool          `json:"alarm_action_enabled"`
	UpdateTime              int64         `json:"update_time"`
	AlarmState              string        `json:"alarm_state"`
}

type realActionInfo struct {
	Type             string   `json:"type"`
	NotificationList []string `json:"notificationList"`
}

func copyActionInfo(src []realActionInfo) []ActionInfo {
	if len(src) == 0 {
		return nil
	}

	dest := make([]ActionInfo, len(src), len(src))
	for i, s := range src {
		d := &dest[i]
		d.Type = s.Type
		d.NotificationList = s.NotificationList
	}
	log.Printf("[DEBUG] copyActionOpts:: src = %#v, dest = %#v", src, dest)
	return dest
}

type alarmRule struct {
	AlarmName               string           `json:"alarm_name"`
	AlarmDescription        string           `json:"alarm_description"`
	Metric                  MetricInfo       `json:"metric"`
	Condition               ConditionInfo    `json:"condition"`
	AlarmActions            []realActionInfo `json:"alarm_actions"`
	InsufficientdataActions []realActionInfo `json:"insufficientdata_actions"`
	OkActions               []realActionInfo `json:"ok_actions"`
	AlarmEnabled            bool             `json:"alarm_enabled"`
	AlarmActionEnabled      bool             `json:"alarm_action_enabled"`
	ID                      string           `json:"alarm_id"`
	UpdateTime              int64            `json:"update_time"`
	AlarmState              string           `json:"alarm_state"`
}

type GetResult struct {
	golangsdk.Result
}

func (g GetResult) Extract() (*AlarmRule, error) {
	var r struct {
		MetricAlarms []alarmRule `json:"metric_alarms"`
	}
	err := g.ExtractInto(&r)
	if err != nil {
		return nil, err
	}
	if len(r.MetricAlarms) != 1 {
		return nil, fmt.Errorf("get %d alarm rules", len(r.MetricAlarms))
	}
	ar0 := r.MetricAlarms[0]
	ar := AlarmRule{
		AlarmName:               ar0.AlarmName,
		AlarmDescription:        ar0.AlarmDescription,
		Metric:                  ar0.Metric,
		Condition:               ar0.Condition,
		AlarmActions:            copyActionInfo(ar0.AlarmActions),
		InsufficientdataActions: copyActionInfo(ar0.InsufficientdataActions),
		OkActions:               copyActionInfo(ar0.OkActions),
		AlarmEnabled:            ar0.AlarmEnabled,
		AlarmActionEnabled:      ar0.AlarmActionEnabled,
		UpdateTime:              ar0.UpdateTime,
		AlarmState:              ar0.AlarmState,
	}
	return &ar, nil
}

type UpdateResult struct {
	golangsdk.ErrResult
}

type DeleteResult struct {
	golangsdk.ErrResult
}
