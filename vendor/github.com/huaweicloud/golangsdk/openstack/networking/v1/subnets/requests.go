package subnets

import (
	"reflect"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the floating IP attributes you want to see returned. SortKey allows you to
// sort by a particular network attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.

type ListOpts struct {
	// ID is the unique identifier for the subnet.
	ID string `json:"id"`

	// Name is the human readable name for the subnet. It does not have to be
	// unique.
	Name string `json:"name"`

	//Specifies the network segment on which the subnet resides.
	CIDR string `json:"cidr"`

	// Status indicates whether or not a subnet is currently operational.
	Status string `json:"status"`

	//Specifies the gateway of the subnet.
	GatewayIP string `json:"gateway_ip"`

	//Specifies the IP address of DNS server 1 on the subnet.
	PRIMARY_DNS string `json:"primary_dns"`

	//Specifies the IP address of DNS server 2 on the subnet.
	SECONDARY_DNS string `json:"secondary_dns"`

	//Identifies the availability zone (AZ) to which the subnet belongs.
	AvailabilityZone string `json:"availability_zone"`

	//Specifies the ID of the VPC to which the subnet belongs.
	VPC_ID string `json:"vpc_id"`
}

// List returns collection of
// subnets. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those subnets that are owned by the
// tenant who submits the request, unless an admin user submits the request.

func List(c *golangsdk.ServiceClient, opts ListOpts) ([]Subnet, error) {
	u := rootURL(c)
	pages, err := pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return SubnetPage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()

	allSubnets, err := ExtractSubnets(pages)
	if err != nil {
		return nil, err
	}

	return FilterSubnets(allSubnets, opts)
}

func FilterSubnets(subnets []Subnet, opts ListOpts) ([]Subnet, error) {

	var refinedSubnets []Subnet
	var matched bool
	m := map[string]interface{}{}

	if opts.ID != "" {
		m["ID"] = opts.ID
	}
	if opts.Name != "" {
		m["Name"] = opts.Name
	}
	if opts.CIDR != "" {
		m["CIDR"] = opts.CIDR
	}
	if opts.Status != "" {
		m["Status"] = opts.Status
	}
	if opts.GatewayIP != "" {
		m["GatewayIP"] = opts.GatewayIP
	}
	if opts.PRIMARY_DNS != "" {
		m["PRIMARY_DNS"] = opts.PRIMARY_DNS
	}
	if opts.SECONDARY_DNS != "" {
		m["SECONDARY_DNS"] = opts.SECONDARY_DNS
	}
	if opts.AvailabilityZone != "" {
		m["AvailabilityZone"] = opts.AvailabilityZone
	}
	if opts.VPC_ID != "" {
		m["VPC_ID"] = opts.VPC_ID
	}

	if len(m) > 0 && len(subnets) > 0 {
		for _, subnet := range subnets {
			matched = true

			for key, value := range m {
				if sVal := getStructField(&subnet, key); !(sVal == value) {
					matched = false
				}
			}

			if matched {
				refinedSubnets = append(refinedSubnets, subnet)
			}
		}

	} else {
		refinedSubnets = subnets
	}

	return refinedSubnets, nil
}

func getStructField(v *Subnet, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToSubnetCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new subnets. There are
// no required values.
type CreateOpts struct {
	Name             string         `json:"name" required:"true"`
	CIDR             string         `json:"cidr" required:"true"`
	DnsList          []string       `json:"dnsList,omitempty"`
	GatewayIP        string         `json:"gateway_ip" required:"true"`
	EnableDHCP       bool           `json:"dhcp_enable" no_default:"y"`
	PRIMARY_DNS      string         `json:"primary_dns,omitempty"`
	SECONDARY_DNS    string         `json:"secondary_dns,omitempty"`
	AvailabilityZone string         `json:"availability_zone,omitempty"`
	VPC_ID           string         `json:"vpc_id" required:"true"`
	ExtraDhcpOpts    []ExtraDhcpOpt `json:"extra_dhcp_opts,omitempty"`
}

type ExtraDhcpOpt struct {
	OptName  string `json:"opt_name" required:"true"`
	OptValue string `json:"opt_value,omitempty"`
}

// ToSubnetCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToSubnetCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "subnet")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// logical subnets. When it is created, the subnets does not have an internal
// interface - it is not associated to any subnet.
//
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSubnetCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular subnets based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	//ToSubnetUpdateMap() (map[string]interface{}, error)
	ToSubnetUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains the values used when updating a subnets.
type UpdateOpts struct {
	Name          string         `json:"name,omitempty"`
	EnableDHCP    bool           `json:"dhcp_enable"`
	PRIMARY_DNS   string         `json:"primary_dns,omitempty"`
	SECONDARY_DNS string         `json:"secondary_dns,omitempty"`
	DnsList       []string       `json:"dnsList,omitempty"`
	ExtraDhcpOpts []ExtraDhcpOpt `json:"extra_dhcp_opts,omitempty"`
}

// ToSubnetUpdateMap builds an update body based on UpdateOpts.
func (opts UpdateOpts) ToSubnetUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "subnet")
}

// Update allows subnets to be updated. You can update the name, administrative
// state, and the external gateway.
func Update(c *golangsdk.ServiceClient, vpcid string, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToSubnetUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(updateURL(c, vpcid, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular subnets based on its unique ID.
func Delete(c *golangsdk.ServiceClient, vpcid string, id string) (r DeleteResult) {
	_, r.Err = c.Delete(updateURL(c, vpcid, id), nil)
	return
}
