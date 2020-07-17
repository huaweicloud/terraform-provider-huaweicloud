package groups

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

//CreateGroupBuilder is an interface from which can build the request of creating group
type CreateOptsBuilder interface {
	ToGroupCreateMap() (map[string]interface{}, error)
}

//CreateGroupOps is a struct contains the parameters of creating group
type CreateOpts struct {
	Name                      string              `json:"scaling_group_name" required:"true"`
	ConfigurationID           string              `json:"scaling_configuration_id,omitempty"`
	DesireInstanceNumber      int                 `json:"desire_instance_number,omitempty"`
	MinInstanceNumber         int                 `json:"min_instance_number,omitempty"`
	MaxInstanceNumber         int                 `json:"max_instance_number,omitempty"`
	CoolDownTime              int                 `json:"cool_down_time,omitempty"`
	LBListenerID              string              `json:"lb_listener_id,omitempty"`
	LBaaSListeners            []LBaaSListenerOpts `json:"lbaas_listeners,omitempty"`
	AvailableZones            []string            `json:"available_zones,omitempty"`
	Networks                  []NetworkOpts       `json:"networks" required:"ture"`
	SecurityGroup             []SecurityGroupOpts `json:"security_groups" required:"ture"`
	VpcID                     string              `json:"vpc_id" required:"ture"`
	HealthPeriodicAuditMethod string              `json:"health_periodic_audit_method,omitempty"`
	HealthPeriodicAuditTime   int                 `json:"health_periodic_audit_time,omitempty"`
	HealthPeriodicAuditGrace  int                 `json:"health_periodic_audit_grace_period,omitempty"`
	InstanceTerminatePolicy   string              `json:"instance_terminate_policy,omitempty"`
	Notifications             []string            `json:"notifications,omitempty"`
	IsDeletePublicip          bool                `json:"delete_publicip,omitempty"`
}

type NetworkOpts struct {
	ID string `json:"id,omitempty"`
}

type SecurityGroupOpts struct {
	ID string `json:"id,omitempty"`
}

type LBaaSListenerOpts struct {
	PoolID       string `json:"pool_id" required:"true"`
	ProtocolPort int    `json:"protocol_port" required:"true"`
	Weight       int    `json:"weight,omitempty"`
}

func (opts CreateOpts) ToGroupCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

//CreateGroup is a method of creating group
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

//DeleteGroup is a method of deleting a group by group id
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

//GetGroup is a method of getting the detailed information of the group by id
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

type ListOptsBuilder interface {
	ToGroupListQuery() (string, error)
}

type ListOpts struct {
	Name            string `q:"scaling_group_name"`
	ConfigurationID string `q:"scaling_configuration_id"`
	Status          string `q:"scaling_group_status"`
}

// ToGroupListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToGroupListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, ops ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if ops != nil {
		q, err := ops.ToGroupListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += q
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return GroupPage{pagination.SinglePageBase(r)}
	})
}

//UpdateOptsBuilder is an interface which can build the map paramter of update function
type UpdateOptsBuilder interface {
	ToGroupUpdateMap() (map[string]interface{}, error)
}

//UpdateOpts is a struct which represents the parameters of update function
type UpdateOpts struct {
	Name                      string              `json:"scaling_group_name,omitempty"`
	DesireInstanceNumber      int                 `json:"desire_instance_number"`
	MinInstanceNumber         int                 `json:"min_instance_number"`
	MaxInstanceNumber         int                 `json:"max_instance_number"`
	CoolDownTime              int                 `json:"cool_down_time,omitempty"`
	LBListenerID              string              `json:"lb_listener_id,omitempty"`
	LBaaSListeners            []LBaaSListenerOpts `json:"lbaas_listeners,omitempty"`
	AvailableZones            []string            `json:"available_zones,omitempty"`
	Networks                  []NetworkOpts       `json:"networks,omitempty"`
	SecurityGroup             []SecurityGroupOpts `json:"security_groups,omitempty"`
	HealthPeriodicAuditMethod string              `json:"health_periodic_audit_method,omitempty"`
	HealthPeriodicAuditTime   int                 `json:"health_periodic_audit_time,omitempty"`
	HealthPeriodicAuditGrace  int                 `json:"health_periodic_audit_grace_period,omitempty"`
	InstanceTerminatePolicy   string              `json:"instance_terminate_policy,omitempty"`
	Notifications             []string            `json:"notifications,omitempty"`
	IsDeletePublicip          bool                `json:"delete_publicip,omitempty"`
	ConfigurationID           string              `json:"scaling_configuration_id,omitempty"`
}

func (opts UpdateOpts) ToGroupUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

//Update is a method which can be able to update the group via accessing to the
//autoscaling service with Put method and parameters
func Update(client *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	body, err := opts.ToGroupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updateURL(client, id), body, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

type ActionOptsBuilder interface {
	ToActionMap() (map[string]interface{}, error)
}

type ActionOpts struct {
	Action string `json:"action" required:"true"`
}

func (opts ActionOpts) ToActionMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func doAction(client *golangsdk.ServiceClient, id string, opts ActionOptsBuilder) (r ActionResult) {
	b, err := opts.ToActionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(enableURL(client, id), &b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

//Enable is an operation by which can make the group enable service
func Enable(client *golangsdk.ServiceClient, id string) (r ActionResult) {
	opts := ActionOpts{
		Action: "resume",
	}
	return doAction(client, id, opts)
}

//Disable is an operation by which can be able to pause the group
func Disable(client *golangsdk.ServiceClient, id string) (r ActionResult) {
	opts := ActionOpts{
		Action: "pause",
	}
	return doAction(client, id, opts)
}
