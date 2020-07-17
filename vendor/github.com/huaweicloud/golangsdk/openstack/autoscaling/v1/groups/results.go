package groups

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

//CreateGroupResult is a struct retured by CreateGroup request
type CreateResult struct {
	golangsdk.Result
}

//Extract the create group result as a string type.
func (r CreateResult) Extract() (string, error) {
	var a struct {
		GroupID string `json:"scaling_group_id"`
	}
	err := r.Result.ExtractInto(&a)
	return a.GroupID, err
}

//DeleteGroupResult contains the body of the deleting group request
type DeleteResult struct {
	golangsdk.ErrResult
}

//GetGroupResult contains the body of getting detailed group request
type GetResult struct {
	golangsdk.Result
}

//Extract method will parse the result body into Group struct
func (r GetResult) Extract() (Group, error) {
	var g Group
	err := r.Result.ExtractIntoStructPtr(&g, "scaling_group")
	return g, err
}

//Group represents the struct of one autoscaling group
type Group struct {
	Name                      string          `json:"scaling_group_name"`
	ID                        string          `json:"scaling_group_id"`
	Status                    string          `json:"scaling_group_status"`
	ConfigurationID           string          `json:"scaling_configuration_id"`
	ConfigurationName         string          `json:"scaling_configuration_name"`
	ActualInstanceNumber      int             `json:"current_instance_number"`
	DesireInstanceNumber      int             `json:"desire_instance_number"`
	MinInstanceNumber         int             `json:"min_instance_number"`
	MaxInstanceNumber         int             `json:"max_instance_number"`
	CoolDownTime              int             `json:"cool_down_time"`
	LBListenerID              string          `json:"lb_listener_id"`
	LBaaSListeners            []LBaaSListener `json:"lbaas_listeners"`
	AvailableZones            []string        `json:"available_zones"`
	Networks                  []Network       `json:"networks"`
	SecurityGroups            []SecurityGroup `json:"security_groups"`
	CreateTime                string          `json:"create_time"`
	VpcID                     string          `json:"vpc_id"`
	Detail                    string          `json:"detail"`
	IsScaling                 bool            `json:"is_scaling"`
	HealthPeriodicAuditMethod string          `json:"health_periodic_audit_method"`
	HealthPeriodicAuditTime   int             `json:"health_periodic_audit_time"`
	HealthPeriodicAuditGrace  int             `json:"health_periodic_audit_grace_period"`
	InstanceTerminatePolicy   string          `json:"instance_terminate_policy"`
	Notifications             []string        `json:"notifications"`
	DeletePublicip            bool            `json:"delete_publicip"`
	CloudLocationID           string          `json:"cloud_location_id"`
}

type Network struct {
	ID string `json:"id"`
}

type SecurityGroup struct {
	ID string `json:"id"`
}

type LBaaSListener struct {
	ListenerID   string `json:"listener_id"`
	PoolID       string `json:"pool_id"`
	ProtocolPort int    `json:"protocol_port"`
	Weight       int    `json:"weight"`
}

type GroupPage struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a ListResult contains no Volumes.
func (r GroupPage) IsEmpty() (bool, error) {
	groups, err := r.Extract()
	return len(groups) == 0, err
}

func (r GroupPage) Extract() ([]Group, error) {
	var gs []Group
	err := r.Result.ExtractIntoSlicePtr(&gs, "scaling_groups")
	return gs, err
}

//UpdateResult is a struct from which can get the result of udpate method
type UpdateResult struct {
	golangsdk.Result
}

//Extract will deserialize the result to group id with string
func (r UpdateResult) Extract() (string, error) {
	var a struct {
		ID string `json:"scaling_group_id"`
	}
	err := r.Result.ExtractInto(&a)
	return a.ID, err
}

//this is the action result which is the result of enable or disable operations
type ActionResult struct {
	golangsdk.ErrResult
}
