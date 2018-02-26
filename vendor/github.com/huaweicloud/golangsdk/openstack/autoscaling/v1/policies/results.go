package policies

import (
	"github.com/huaweicloud/golangsdk"
)

//Create Result is a struct which represents the create result of policy
type CreateResult struct {
	golangsdk.Result
}

//Extract of CreateResult will deserialize the result of Creation
func (r CreateResult) Extract() (string, error) {
	var a struct {
		ID string `json:"scaling_policy_id"`
	}
	err := r.Result.ExtractInto(&a)
	return a.ID, err
}

//DeleteResult is a struct which represents the delete result.
type DeleteResult struct {
	golangsdk.ErrResult
}

//Policy is a struct that represents the result of get policy
type Policy struct {
	ID             string         `json:"scaling_group_id"`
	Name           string         `json:"scaling_policy_name"`
	Status         string         `json:"policy_status"`
	Type           string         `json:"scaling_policy_type"`
	AlarmID        string         `json:"alarm_id"`
	SchedulePolicy SchedulePolicy `json:"scheduled_policy"`
	Action         Action         `json:"scaling_policy_action"`
	CoolDownTime   int            `json:"cool_down_time"`
	CreateTime     string         `json:"create_time"`
}

type SchedulePolicy struct {
	LaunchTime      string `json:"launch_time"`
	RecurrenceType  string `json:"recurrence_type"`
	RecurrenceValue string `json:"recurrence_value"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
}

type Action struct {
	Operation   string `json:"operation"`
	InstanceNum int    `json:"instance_number"`
}

//GetResult is a struct which represents the get result
type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (Policy, error) {
	var p Policy
	err := r.Result.ExtractIntoStructPtr(&p, "scaling_policy")
	return p, err
}

//UpdateResult is a struct from which can get the result of udpate method
type UpdateResult struct {
	golangsdk.Result
}

//Extract will deserialize the result to group id with string
func (r UpdateResult) Extract() (string, error) {
	var a struct {
		ID string `json:"scaling_policy_id"`
	}
	err := r.Result.ExtractInto(&a)
	return a.ID, err
}
