package alarmrule

import (
	"log"

	"github.com/huaweicloud/golangsdk"
)

type CreateOptsBuilder interface {
	ToAlarmRuleCreateMap() (map[string]interface{}, error)
}

type DimensionOpts struct {
	Name  string `json:"name" required:"true"`
	Value string `json:"value" required:"true"`
}

type MetricOpts struct {
	Namespace  string          `json:"namespace" required:"true"`
	MetricName string          `json:"metric_name" required:"true"`
	Dimensions []DimensionOpts `json:"dimensions" required:"true"`
}

type ConditionOpts struct {
	Period             int    `json:"period" required:"true"`
	Filter             string `json:"filter" required:"true"`
	ComparisonOperator string `json:"comparison_operator" required:"true"`
	// The Value ranges from 0 to MAX_VALUE
	Value int    `json:"value"`
	Unit  string `json:"unit,omitempty"`
	Count int    `json:"count" required:"true"`
}

type ActionOpts struct {
	Type             string   `json:"type" required:"true"`
	NotificationList []string `json:"notificationList" required:"true"`
}

type CreateOpts struct {
	AlarmName               string        `json:"alarm_name" required:"true"`
	AlarmDescription        string        `json:"alarm_description,omitempty"`
	AlarmType               string        `json:"alarm_type,omitempty"`
	AlarmLevel              int           `json:"alarm_level,omitempty"`
	Metric                  MetricOpts    `json:"metric" required:"true"`
	Condition               ConditionOpts `json:"condition" required:"true"`
	AlarmActions            []ActionOpts  `json:"alarm_actions,omitempty"`
	InsufficientdataActions []ActionOpts  `json:"insufficientdata_actions,omitempty"`
	OkActions               []ActionOpts  `json:"ok_actions,omitempty"`
	AlarmEnabled            bool          `json:"alarm_enabled"`
	AlarmActionEnabled      bool          `json:"alarm_action_enabled"`
}

func (opts CreateOpts) ToAlarmRuleCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToAlarmRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	log.Printf("[DEBUG] create AlarmRule url:%q, body=%#v, opt=%#v", rootURL(c), b, opts)
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{201}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

type UpdateOptsBuilder interface {
	ToAlarmRuleUpdateMap() (map[string]interface{}, error)
}

type UpdateOpts struct {
	AlarmEnabled bool `json:"alarm_enabled"`
}

func (opts UpdateOpts) ToAlarmRuleUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func Update(c *golangsdk.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToAlarmRuleUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(actionURL(c, id), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{204}}
	_, r.Err = c.Delete(resourceURL(c, id), reqOpt)
	return
}
