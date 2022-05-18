package engines

import "github.com/chnsz/golangsdk"

// CreateOpts is the structure required by the Create method to create a new engine.
type CreateOpts struct {
	// The name of the dedicated microservice engine.
	// The value can contain 3 to 24 characters, including letters, digits, and hyphens (-).
	// It must start with a letter and cannot end with a hyphen.
	Name string `json:"name" required:"true"`
	// The charging mode of the dedicated microservice engine.
	Payment string `json:"payment" required:"true"`
	// The flavor of the dedicated microservice engine.
	//   cse.s1.small2: High availability 100 instance engine.
	//   cse.s1.medium2: High availability 200 instance engine.
	//   cse.s1.large2: High availability 500 instance engine.
	//   cse.s1.xlarge2: High availability 2000 instance engine.
	Flavor string `json:"flavor" required:"true"`
	// List of available zones for the current region.
	AvailabilityZones []string `json:"azList" required:"true"`
	// The authentication method for the dedicated microservice engine.
	// The "RBAC" is security authentication, and the "NONE" is no authentication.
	AuthType string `json:"authType" required:"true"`
	// The VPC name.
	VpcName string `json:"vpc" required:"true"`
	// The network ID of the subnet.
	NetworkId string `json:"networkId" required:"true"`
	// The subnet CIDR.
	SubnetCidr string `json:"subnetCidr" required:"true"`
	// The deployment type of the dedicated microservice engine. The fixed value is "CSE2".
	SpecType string `json:"specType" required:"true"`
	// The description of the dedicated microservice engine. The value can contain up to 255 characters.
	Description string `json:"description,omitempty"`
	// The VPC ID.
	VpcId string `json:"vpcId,omitempty"`
	// The public access for the dedicated microservice engine.
	PublicIpId string `json:"publicIpId,omitempty"`
	// The dedicated microservice engine must be passed when security authentication is selected,
	// including the authentication information of the engine.
	AuthCred *AuthCred `json:"auth_cred,omitempty"`
	// The additional parameters of the dedicated microservice engine.
	Inputs map[string]interface{} `json:"inputs,omitempty"`
	// The enterprise project ID to which the dedicated microservice engine.
	EnterpriseProjectId string `json:"-"`
}

// AuthCred is an object that specifies the configuration of the security authentication.
type AuthCred struct {
	// The root account password that needs to be specified when selecting security authentication.
	Password string `json:"pwd" required:"true"`
}

func buildMoreHeaderUsingEpsId(epsId string) map[string]string {
	moreHeader := map[string]string{
		"Content-Type": "application/json",
		"X-Language":   "en-us",
	}

	if epsId != "" {
		moreHeader["X-Enterprise-Project-ID"] = epsId
	}

	return moreHeader
}

// Create is a method to create a microservice engine using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*RequestResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r RequestResp
	_, err = c.Post(rootURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: buildMoreHeaderUsingEpsId(opts.EnterpriseProjectId),
	})
	return &r, err
}

// Get is a method to retrieves a particular configuration based on its unique ID and enterprise project ID.
func Get(c *golangsdk.ServiceClient, engineId, epsId string) (*Engine, error) {
	var r Engine
	_, err := c.Get(resourceURL(c, engineId), &r, &golangsdk.RequestOpts{
		MoreHeaders: buildMoreHeaderUsingEpsId(epsId),
	})
	return &r, err
}

// Delete is a method to remove an existing engine using its unique ID and enterprise project ID.
func Delete(c *golangsdk.ServiceClient, engineId, epsId string) (*RequestResp, error) {
	var r RequestResp
	_, err := c.Delete(resourceURL(c, engineId), &golangsdk.RequestOpts{
		JSONResponse: &r,
		MoreHeaders:  buildMoreHeaderUsingEpsId(epsId),
	})
	return &r, err
}

// GetJob is a method to obtain the job detail using engine ID and job ID.
func GetJob(c *golangsdk.ServiceClient, engineId, jobId string) (*Job, error) {
	var r Job
	_, err := c.Get(jobURL(c, engineId, jobId), &r, nil)
	return &r, err
}
