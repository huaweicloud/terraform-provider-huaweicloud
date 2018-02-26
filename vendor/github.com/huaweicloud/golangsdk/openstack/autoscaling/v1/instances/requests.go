package instances

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

//ListOptsBuilder is an interface by which can be able to build the query string
//of the list function
type ListOptsBuilder interface {
	ToInstancesListQuery() (string, error)
}

type ListOpts struct {
	LifeCycleStatus string `q:"life_cycle_state"`
	HealthStatus    string `q:"health_status"`
}

func (opts ListOpts) ToInstancesListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

//List is a method by which can be able to access the list function that can get
//instances of a group
func List(client *golangsdk.ServiceClient, groupID string, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client, groupID)
	if opts != nil {
		q, err := opts.ToInstancesListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += q
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return InstancePage{pagination.SinglePageBase(r)}
	})
}

//DeleteOptsBuilder is an interface by whick can be able to build the query string
//of instance deletion
type DeleteOptsBuilder interface {
	ToInstanceDeleteQuery() (string, error)
}

type DeleteOpts struct {
	DeleteInstance bool `q:"instance_delete"`
}

func (opts DeleteOpts) ToInstanceDeleteQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

//Delete is a method by which can be able to delete an instance from a group
func Delete(client *golangsdk.ServiceClient, id string, opts DeleteOptsBuilder) (r DeleteResult) {
	url := deleteURL(client, id)
	if opts != nil {
		q, err := opts.ToInstanceDeleteQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += q
	}
	_, r.Err = client.Delete(url, nil)
	return
}

//BatchOptsBuilder is an interface which can build the query body of batch operation
type BatchOptsBuilder interface {
	ToInstanceBatchMap() (map[string]interface{}, error)
}

//BatchOpts is a struct which represents parameters of batch operations
type BatchOpts struct {
	Instances   []string `json:"instances_id" required:"true"`
	IsDeleteEcs string   `json:"instance_delete,omitempty"`
	Action      string   `json:"action,omitempty"`
}

func (opts BatchOpts) ToInstanceBatchMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

//batch is method which can be able to add/delete numbers instances
func batch(client *golangsdk.ServiceClient, groupID string, opts BatchOptsBuilder) (r BatchResult) {
	b, err := opts.ToInstanceBatchMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(batchURL(client, groupID), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

//BatchAdd is a method by which can add numbers of instances into a group
func BatchAdd(client *golangsdk.ServiceClient, groupID string, instances []string) (r BatchResult) {
	var opts = BatchOpts{
		Instances: instances,
		Action:    "ADD",
	}
	return batch(client, groupID, opts)
}

//BatchDelete is a method by which can delete numbers of instances from a group
func BatchDelete(client *golangsdk.ServiceClient, groupID string, instances []string, deleteEcs string) (r BatchResult) {
	var opts = BatchOpts{
		Instances:   instances,
		IsDeleteEcs: deleteEcs,
		Action:      "REMOVE",
	}
	return batch(client, groupID, opts)
}
