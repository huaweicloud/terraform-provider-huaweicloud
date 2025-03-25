package alarmrule

import (
	"github.com/chnsz/golangsdk"
)

type DimensionOpts struct {
	Name  string `json:"name" required:"true"`
	Value string `json:"value,omitempty"`
}

type PolicyOpts struct {
	MetricName string `json:"metric_name" required:"true"`
	// The value can be 0
	Period             int    `json:"period"`
	Filter             string `json:"filter" required:"true"`
	ComparisonOperator string `json:"comparison_operator" required:"true"`
	// The Value ranges from 0 to MAX_VALUE
	Value float64 `json:"value"`
	Unit  string  `json:"unit,omitempty"`
	Count int     `json:"count" required:"true"`
	// The value can be 0
	SuppressDuration int `json:"suppress_duration"`
	Level            int `json:"level,omitempty"`
}

type NotificationOpts struct {
	Type             string   `json:"type" required:"true"`
	NotificationList []string `json:"notification_list" required:"true"`
}

type CreateOpts struct {
	Name                  string             `json:"name" required:"true"`
	Description           string             `json:"description,omitempty"`
	Namespace             string             `json:"namespace" required:"true"`
	ResourceGroupID       string             `json:"resource_group_id,omitempty"`
	Resources             [][]DimensionOpts  `json:"resources" required:"true"`
	Policies              []PolicyOpts       `json:"policies,omitempty"`
	Type                  string             `json:"type" required:"true"`
	AlarmNotifications    []NotificationOpts `json:"alarm_notifications,omitempty"`
	OkNotifications       []NotificationOpts `json:"ok_notifications,omitempty"`
	NotificationBeginTime string             `json:"notification_begin_time,omitempty"`
	NotificationEndTime   string             `json:"notification_end_time,omitempty"`
	EnterpriseProjectID   string             `json:"enterprise_project_id,omitempty"`
	Enabled               bool               `json:"enabled"`
	NotificationEnabled   bool               `json:"notification_enabled"`
	AlarmTemplateID       string             `json:"alarm_template_id,omitempty"`
}

func Create(c *golangsdk.ServiceClient, opts CreateOpts) (r CreateResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{201}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(listURL(c, id), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	})
	return
}

type ActionOpts struct {
	AlarmIDs     []string `json:"alarm_ids" required:"true"`
	AlarmEnabled bool     `json:"alarm_enabled"`
}

func Action(c *golangsdk.ServiceClient, id string, opts ActionOpts) (r ActionResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(actionURL(c), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

type UpdateResourcesOpts struct {
	Resources [][]DimensionOpts `json:"resources" required:"true"`
}

func BatchResources(c *golangsdk.ServiceClient, id, operation string, opts UpdateResourcesOpts) (r BatchResourcesResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(batchResourcesURL(c, id, operation), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

type UpdatePoliciesOpts struct {
	Policies []PolicyOpts `json:"policies" required:"true"`
}

func PoliciesModify(c *golangsdk.ServiceClient, id string, opts UpdatePoliciesOpts) (r PoliciesResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(policiesURL(c, id), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

func GetResources(c *golangsdk.ServiceClient, id string) (r GetResourcesResult) {
	_, r.Err = c.Get(resourcesURL(c, id), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json; charset=UTF-8",
		},
	})
	return
}

type DeleteOpts struct {
	AlarmIDs []string `json:"alarm_ids"`
}

func Delete(c *golangsdk.ServiceClient, opts DeleteOpts) (r DeleteResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(deleteURL(c), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
