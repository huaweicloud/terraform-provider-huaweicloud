package loggroups

import "github.com/huaweicloud/golangsdk"

// CreateOptsBuilder is used for creating log group parameters.
type CreateOptsBuilder interface {
	ToLogGroupsCreateMap() (map[string]interface{}, error)
}

// CreateOpts is a struct that contains all the parameters.
type CreateOpts struct {
	// Specifies the log group name.
	LogGroupName string `json:"log_group_name" required:"true"`

	// Specifies the log expiration time.
	TTL int `json:"ttl_in_days,omitempty"`
}

// ToLogGroupsCreateMap is used for type convert
func (ops CreateOpts) ToLogGroupsCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(ops, "")
}

// Create a log group with given parameters.
func Create(client *golangsdk.ServiceClient, ops CreateOptsBuilder) (r CreateResult) {
	b, err := ops.ToLogGroupsCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{201},
	})

	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToLogGroupsUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contain options for updating an existing Group.
// For more information about the parameters, see the LogGroup object.
type UpdateOpts struct {
	// Specifies the log expiration time.
	TTL int `json:"ttl_in_days,omitempty"`
}

// ToLogGroupsUpdateMap is used for type convert
func (ops UpdateOpts) ToLogGroupsUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(ops, "")
}

// update a log group with given parameters by id.
func Update(client *golangsdk.ServiceClient, ops UpdateOptsBuilder, id string) (r UpdateResult) {
	b, err := ops.ToLogGroupsUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(updateURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})

	return
}

// Delete a log group by id
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	opts := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=utf8"},
	}
	_, r.Err = client.Delete(deleteURL(client, id), &opts)
	return
}

// Get log group list
func List(client *golangsdk.ServiceClient) (r ListResults) {
	opts := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=utf8"},
	}
	_, r.Err = client.Get(listURL(client), &r.Body, &opts)
	return
}
