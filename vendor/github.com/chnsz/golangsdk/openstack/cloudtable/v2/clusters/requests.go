package clusters

import "github.com/chnsz/golangsdk"

// CreateOpts is a struct which will be used to create a new cloudtable cluster.
type CreateOpts struct {
	// Create cluster database parameters.
	Datastore Datastore `json:"datastore" required:"true"`
	// Cluster name
	Name string `json:"name" required:"true"`
	// The instance object of the cluster.
	Instance Instance `json:"instance" required:"true"`
	// Storage type, the valid values are:
	//   ULTRAHIGH
	//   COMMON
	StorageType string `json:"storage_type" required:"true"`
	// The VPC where the cluster is located.
	VpcId string `json:"vpc_id" required:"true"`
	// Whether the IAM auth is enabled.
	IAMAuthEnabled bool `json:"auth_mode,omitempty"`
	// Whether the Lemon is enabled.
	LemonEnabled bool `json:"enable_lemon,omitempty"`
	// Whether the OpenTSDB is enabled.
	OpenTSDBEnabled bool `json:"enable_openTSDB,omitempty"`
	// The size of the stored value.
	StorageSize int `json:"storage_size,omitempty"`
}

// Datastore is an object specifying the cluster storage.
type Datastore struct {
	// Cluster database type.
	Type string `json:"type" required:"true"`
	// Controller version number, default to '1.0.6'.
	Version string `json:"version" required:"true"`
}

// Instance is an object specifying the information of the nodes number, az and network.
type Instance struct {
	// The ID of the availability zone where the cluster is located.
	AvailabilityZone string `json:"availability_zone" required:"true"`
	// The number of computing unit nodes in the CloudTable cluster must be at least 2.
	CUNum int `json:"cu_num" required:"true"`
	// Information about the network where the cluster is located.
	Networks []Network `json:"nics" required:"true"`
	// The number of Lemon nodes in the CloudTable cluster.
	LemonNum int `json:"lemon_num,omitempty"`
	// The number of TSD nodes in the CloudTable cluster must be at least 2.
	TSDNum int `json:"tsd_num,omitempty"`
}

// Network is an object specifying the cluster network.
type Network struct {
	// The ID of the network where the CloudTable cluster is located.
	SubnetId string `json:"net_id" required:"true"`
	// The ID of the security group to which the CloudTable belongs.
	SecurityGroupId string `json:"security_group_id" required:"true"`
}

// Create is a method to create a new cluster.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*RequestResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "cluster")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(rootURL(c), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
	})
	if err == nil {
		var r RequestResp
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// Get is a method to obtain the detail of the cluster.
func Get(c *golangsdk.ServiceClient, clusterId string) (*Cluster, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, clusterId), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
	})
	if err == nil {
		var r Cluster
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// Delete is a method to remove an existing cluster by ID.
func Delete(c *golangsdk.ServiceClient, clusterId string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult
	_, r.Err = c.Delete(resourceURL(c, clusterId), &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
	})
	return &r
}
