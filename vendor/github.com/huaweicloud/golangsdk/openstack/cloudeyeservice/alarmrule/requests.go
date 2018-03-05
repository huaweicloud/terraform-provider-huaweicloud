package alarmrule

import (
	"fmt"
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
	Value              int    `json:"value" required:"true"`
	Unit               string `json:"unit,omitempty"`
	Count              int    `json:"count" required:"true"`
}

type ActionOpts struct {
	Type             string   `json:"type" required:"true"`
	NotificationList []string `json:"notification_list" required:"true"`
}

type CreateOpts struct {
	AlarmName               string        `json:"alarm_name" required:"true"`
	AlarmDescription        string        `json:"alarm_description,omitempty"`
	Metric                  MetricOpts    `json:"metric" required:"true"`
	Condition               ConditionOpts `json:"condition" required:"true"`
	AlarmActions            []ActionOpts  `json:"alarm_actions,omitempty"`
	InsufficientdataActions []ActionOpts  `json:"insufficientdata_actions,omitempty"`
	OkActions               []ActionOpts  `json:"ok_actions,omitempty"`
	AlarmEnabled            bool          `json:"alarm_enabled"`
	AlarmActionEnabled      bool          `json:"alarm_action_enabled"`
}

func (opts CreateOpts) ToAlarmRuleCreateMap() (map[string]interface{}, error) {
	return nil, fmt.Errorf("no need implement")
}

type realActionOpts struct {
	Type             string   `json:"type" required:"true"`
	NotificationList []string `json:"notificationList" required:"true"`
}

type createOpts struct {
	AlarmName               string           `json:"alarm_name" required:"true"`
	AlarmDescription        string           `json:"alarm_description,omitempty"`
	Metric                  MetricOpts       `json:"metric" required:"true"`
	Condition               ConditionOpts    `json:"condition" required:"true"`
	AlarmActions            []realActionOpts `json:"alarm_actions,omitempty"`
	InsufficientdataActions []realActionOpts `json:"insufficientdata_actions,omitempty"`
	OkActions               []realActionOpts `json:"ok_actions,omitempty"`
	AlarmEnabled            bool             `json:"alarm_enabled"`
	AlarmActionEnabled      bool             `json:"alarm_action_enabled"`
}

func (opts createOpts) ToAlarmRuleCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func copyActionOpts(src []ActionOpts) []realActionOpts {
	if len(src) == 0 {
		return nil
	}

	dest := make([]realActionOpts, len(src), len(src))
	for i, s := range src {
		d := &dest[i]
		d.Type = s.Type
		d.NotificationList = s.NotificationList
	}
	log.Printf("[DEBUG] copyActionOpts:: src = %#v, dest = %#v", src, dest)
	return dest
}

func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	opt := opts.(CreateOpts)
	opts1 := createOpts{
		AlarmName:               opt.AlarmName,
		AlarmDescription:        opt.AlarmDescription,
		Metric:                  opt.Metric,
		Condition:               opt.Condition,
		AlarmActions:            copyActionOpts(opt.AlarmActions),
		InsufficientdataActions: copyActionOpts(opt.InsufficientdataActions),
		OkActions:               copyActionOpts(opt.OkActions),
		AlarmEnabled:            opt.AlarmEnabled,
		AlarmActionEnabled:      opt.AlarmActionEnabled,
	}
	b, err := opts1.ToAlarmRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	log.Printf("[DEBUG] create AlarmRule url:%q, body=%#v, opt=%#v", rootURL(c), b, opts1)
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
