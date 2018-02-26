package backendmember

import (
	//"fmt"
	"github.com/huaweicloud/golangsdk"
	//"github.com/huaweicloud/golangsdk/pagination"
)

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type AddOptsBuilder interface {
	ToBackendAddMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type AddOpts struct {
	// server_id
	ServerId string `json:"server_id" required:"true"`
	// The load balancer on which to provision this listener.
	Address string `json:"address" required:"true"`
}

// ToBackendAddMap casts a CreateOpts struct to a map.
func (opts AddOpts) ToBackendAddMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Add is an operation which provisions a new Listeners based on the
// configuration defined in the AddOpts struct. Once the request is
// validated and progress has started on the provisioning process, a
// AddResult will be returned.
//
// Users with an admin role can create Listeners on behalf of other tenants by
// specifying a TenantID attribute different than their own.
func Add(c *golangsdk.ServiceClient, listener_id string, opts AddOptsBuilder) (r AddResult) {
	b, err := opts.ToBackendAddMap()
	// API takes an array of these...
	a := make([]map[string]interface{}, 1)
	a[0] = b
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(addURL(c, listener_id), a, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// RemoveOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Remove operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type RemoveOptsBuilder interface {
	ToBackendRemoveMap() (map[string]interface{}, error)
}

type LoadBalancerID struct {
	// backend member id to remove
	ID string `json:"id" required:"true"`
}

// RemoveOpts is the common options struct used in this package's Remove
// operation.
type RemoveOpts struct {
	RemoveMember []LoadBalancerID `json:"removeMember" required:"true"`
}

// ToBackendCreateMap casts a CreateOpts struct to a map.
func (opts RemoveOpts) ToBackendRemoveMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Remove will permanently remove a particular backend based on its unique ID.
func Remove(c *golangsdk.ServiceClient, listener_id string, id string) (r RemoveResult) {
	/*lbid := LoadBalancerID{
		ID: id,
	}
	lbids := []LoadBalancerID{lbid}
	removeOpts := RemoveOpts{
		removeMember: lbids,
	}
	fmt.Printf("removeOpts=%+v.\n", removeOpts)
	b, err := removeOpts.ToBackendRemoveMap() */
	lbid := make(map[string]interface{})
	lbid["id"] = id
	lbids := make([]map[string]interface{}, 1)
	lbids[0] = lbid
	b := make(map[string]interface{})
	b["removeMember"] = lbids
	//fmt.Printf("b=%+v.\n", b)
	/* if err != nil {
		r.Err = err
		return
	} */
	_, r.Err = c.Post(removeURL(c, listener_id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get retrieves a particular Health Monitor based on its unique ID.
func Get(c *golangsdk.ServiceClient, listener_id, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, listener_id, id), &r.Body, nil)
	return
}
