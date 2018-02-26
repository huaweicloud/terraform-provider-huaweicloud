package loadbalancer_elbs

import (
	"github.com/huaweicloud/golangsdk"
	// "github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/elbaas/listeners"
	"github.com/huaweicloud/golangsdk/pagination"
	//"fmt"
)

// LoadBalancer is the primary load balancing configuration object that specifies
// the virtual IP address on which client traffic is received, as well
// as other details such as the load balancing method to be use, protocol, etc.
type LoadBalancer struct {
	// Specifies the IP address used by ELB for providing services.
	VipAddress string `json:"vip_address"`
	// Specifies the time when information about the load balancer was updated.
	UpdateTime string `json:"update_time"`
	// Specifies the time when the load balancer was created.
	CreateTime string `json:"create_time"`
	// Specifies the load balancer ID.
	ID string `json:"id"`
	// Specifies the status of the load balancer.
	// The value can be ACTIVE, PENDING_CREATE, or ERROR.
	Status string `json:"status"`
	// Specifies the bandwidth (Mbit/s).
	Bandwidth int `json:"bandwidth"`
	// Specifies the VPC ID.
	VpcID string `json:"vpc_id"`
	// Specifies the status of the load balancer.
	// Optional values:
	// 0: The load balancer is disabled.
	// 1: The load balancer is running properly.
	// 2: The load balancer is frozen.
	AdminStateUp int `json:"admin_state_up"`
	// Specifies the subnet ID of backend ECSs.
	VipSubnetID string `json:"vip_subnet_id"`
	// Specifies the load balancer type.
	Type string `json:"type"`
	// Specifies the load balancer name.
	Name string `json:"name"`
	// Provides supplementary information about the load balancer.
	Description string `json:"description"`
	// Specifies the security group ID.
	SecurityGroupID string `json:"security_group_id"`
	// Specifies the ID of the availability zone (AZ).
	AZ string `json:"az"`
}

// LoadBalancerPage is the page returned by a pager when traversing over a
// collection of routers.
type LoadBalancerPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of routers has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
/* func (r LoadBalancerPage) NextPageURL() (string, error) {
	return "", nil
} */

// IsEmpty checks whether a LoadBalancerPage struct is empty.
func (p LoadBalancerPage) IsEmpty() (bool, error) {
	is, err := ExtractLoadBalancers(p)
	return len(is) == 0, err
}

// ExtractLoadBalancers accepts a Page struct, specifically a LoadbalancerPage struct,
// and extracts the elements into a slice of LoadBalancer structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractLoadBalancers(r pagination.Page) ([]LoadBalancer, error) {
	var s struct {
		LoadBalancers []LoadBalancer `json:"loadbalancers"`
	}
	err := (r.(LoadBalancerPage)).ExtractInto(&s)
	return s.LoadBalancers, err
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a loadbalancer.
func (r commonResult) Extract() (*LoadBalancer, error) {
	//fmt.Printf("Extracting...\n")
	lb := new(LoadBalancer)
	err := r.ExtractInto(lb)
	if err != nil {
		//fmt.Printf("Error: %s.\n", err.Error())
		return nil, err
	} else {
		//fmt.Printf("Returning extract: %+v.\n", lb)
		return lb, nil
	}
}

type GetStatusesResult struct {
	golangsdk.Result
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	commonResult
}
