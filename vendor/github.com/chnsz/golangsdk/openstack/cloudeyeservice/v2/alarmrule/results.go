package alarmrule

import (
	"fmt"

	"github.com/chnsz/golangsdk"
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

type ResourcesInfo struct {
	ResourceGroupID   string          `json:"resource_group_id"`
	ResourceGroupName string          `json:"resource_group_name"`
	Dimensions        []DimensionOpts `json:"dimensions"`
}

type AlarmRule struct {
	AlarmID               string             `json:"alarm_id"`
	Name                  string             `json:"name"`
	Description           string             `json:"description"`
	Namespace             string             `json:"namespace"`
	Resources             []ResourcesInfo    `json:"resources"`
	Policies              []PolicyOpts       `json:"policies"`
	Type                  string             `json:"type"`
	AlarmNotifications    []NotificationOpts `json:"alarm_notifications"`
	OkNotifications       []NotificationOpts `json:"ok_notifications"`
	NotificationBeginTime string             `json:"notification_begin_time"`
	NotificationEndTime   string             `json:"notification_end_time"`
	EnterpriseProjectID   string             `json:"enterprise_project_id"`
	Enabled               bool               `json:"enabled"`
	NotificationEnabled   bool               `json:"notification_enabled"`
	AlarmTemplateID       string             `json:"alarm_template_id"`
}

type GetResult struct {
	golangsdk.Result
}

func (g GetResult) Extract() (*AlarmRule, error) {
	var r struct {
		Alarms []AlarmRule `json:"alarms"`
	}
	err := g.ExtractInto(&r)
	if err != nil {
		return nil, err
	}
	if len(r.Alarms) != 1 {
		return nil, fmt.Errorf("get %d alarm rules", len(r.Alarms))
	}
	return &(r.Alarms[0]), nil
}

type GetResourcesResult struct {
	golangsdk.Result
}

func (g GetResourcesResult) Extract() (*[][]DimensionOpts, error) {
	var r struct {
		Resources [][]DimensionOpts `json:"resources"`
	}
	err := g.ExtractInto(&r)
	if err != nil {
		return nil, err
	}
	if len(r.Resources) < 1 {
		return nil, fmt.Errorf("get %d alarm resources", len(r.Resources))
	}
	return &(r.Resources), nil
}

type UpdateResult struct {
	golangsdk.ErrResult
}

type ActionResult struct {
	golangsdk.ErrResult
}

type BatchResourcesResult struct {
	golangsdk.ErrResult
}

type PoliciesResult struct {
	golangsdk.ErrResult
}

type DeleteResult struct {
	golangsdk.ErrResult
}
