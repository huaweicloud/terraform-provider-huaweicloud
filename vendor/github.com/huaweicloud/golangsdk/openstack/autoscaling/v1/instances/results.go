package instances

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

//Instance is a struct which represents all the infromation of a instance
type Instance struct {
	ID                string `json:"instance_id"`
	Name              string `json:"instance_name"`
	GroupID           string `json:"scaling_group_id"`
	GroupName         string `json:"scaling_group_name"`
	LifeCycleStatus   string `json:"life_cycle_state"`
	HealthStatus      string `json:"health_status"`
	ConfigurationName string `json:"scaling_configuration_name"`
	ConfigurationID   string `json:"scaling_configuration_id"`
	CreateTime        string `json:"create_time"`
}

//InstancePage is a struct which can do the page function
type InstancePage struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a instances number equal to 0.
func (r InstancePage) IsEmpty() (bool, error) {
	groups, err := r.Extract()
	return len(groups) == 0, err
}

func (r InstancePage) Extract() ([]Instance, error) {
	var instances []Instance
	err := r.Result.ExtractIntoSlicePtr(&instances, "scaling_group_instances")
	return instances, err
}

//DeleteResult is a struct which contains the result of deletion instance
type DeleteResult struct {
	golangsdk.ErrResult
}

//BatchResult is a struct which contains the result of batch operations
type BatchResult struct {
	golangsdk.ErrResult
}
