package servers

import (
	"net/http"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/common"
)

// GetServer get server detail
func GetServer(client *golangsdk.ServiceClient, serverID string) (r GetResult) {
	GetUrl := getURL(client, serverID)

	var resp *http.Response
	resp, r.Err = client.Get(GetUrl, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	if r.Err != nil {
		return
	}
	defer resp.Body.Close()

	return
}

// UpdateOptsBuilder to handle Delete request's body
type UpdateOptsBuilder interface {
	ToUpdateBodyMap() (map[string]interface{}, error)
}

// UpdateInstance to update one server
type UpdateInstance struct {
	UpdateServer UpdateOpts `json:"server"`
}

// UpdateOpts to update one server options
type UpdateOpts struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ToUpdateBodyMap converts UpdateInstance structures to map[string]interface{}
func (opts UpdateInstance) ToUpdateBodyMap() (map[string]interface{}, error) {
	body, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return body, nil
}

// UpdateServer updates one server in UpdateOptsBuilder
func UpdateServer(client *golangsdk.ServiceClient, opts UpdateOptsBuilder, serverID string) (r UpdateResult) {
	body, err := opts.ToUpdateBodyMap()
	if err != nil {
		r.Err = err
		return
	}

	updateURL := updateURL(client, serverID)

	var resp *http.Response
	resp, r.Err = client.Put(updateURL, body, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	if r.Err != nil {
		return
	}
	defer resp.Body.Close()

	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToServerCreateMap() (map[string]interface{}, error)
}

// ToServerCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToServerCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"server": b}, nil
}

type CreateOpts struct {
	common.ResourceOpts
	Coverage common.Coverage `json:"coverage" required:"true"`
}

// CreateServer requests a server to be provisioned to the user in the current tenant with response entity.
func CreateServer(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (createResult CreateCloudServerResponse, err error) {
	var r CreateResult
	reqBody, err := opts.ToServerCreateMap()
	if err != nil {
		return
	}

	createURL := createURL(client)
	var resp *http.Response
	resp, err = client.Post(createURL, reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusAccepted},
	})
	if err != nil {
		return
	}

	job, errJob := r.ExtractJob()
	if errJob != nil {
		err = errJob
		return
	}
	server, errServer := r.ExtractServer()
	if errServer != nil {
		err = errServer
		return
	}
	createResult = CreateCloudServerResponse{Job: job,
		ServerIDs: server}

	defer resp.Body.Close()
	return
}

// DeleteOptsBuilder to handle Delete request's body
type DeleteOptsBuilder interface {
	ToDeleteBodyMap() (map[string]interface{}, error)
}

// DeleteOpts to delete all Servers in server array
type DeleteOpts struct {
	Servers []cloudservers.Server `json:"servers"`
}

// ToDeleteBodyMap converts DeleteOpts structures to map[string]interface{}
func (opts DeleteOpts) ToDeleteBodyMap() (map[string]interface{}, error) {
	body, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return body, nil
}

// DeleteServers deletes all server in DeleteOptsBuilder
func DeleteServers(client *golangsdk.ServiceClient, opts DeleteOptsBuilder) (r DeleteResult) {
	body, err := opts.ToDeleteBodyMap()
	if err != nil {
		r.Err = err
		return
	}
	deleteURL := deleteAllServersURL(client)

	var resp *http.Response
	resp, r.Err = client.Post(deleteURL, body, nil, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusNoContent},
	})
	if r.Err != nil {
		return
	}
	defer resp.Body.Close()

	return
}
