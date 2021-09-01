package icagents

import (
	"github.com/chnsz/golangsdk"
)

type CreateOptsBuilder interface {
	ToIcAgentIstallMap() (map[string]interface{}, error)
}

// Installing parameters for ICagent in cce cluster
type InstallParam struct {
	// The ID of Cluster
	ClusterId string `json:"clusterId" required:"true"`
	// Namespace for agent
	NameSpace string `json:"nameSpace" required:"true"`
}

// ToIcAgentIstallMap builds a create request body from InstallParam.
func (installParam InstallParam) ToIcAgentIstallMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(installParam, "")
}

// Create accepts a CreateOpts struct and uses the values to intall ic agent in cluster.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToIcAgentIstallMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}
