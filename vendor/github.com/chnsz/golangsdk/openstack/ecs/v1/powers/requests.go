package powers

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
)

// PowerOpts allows batch update of the ECS instances power state through the API.
// Parameter 'Type' supports two methods of 'SOFT' and 'HARD' to shut down and reboot the machine's power.
type PowerOpts struct {
	Servers []ServerInfo `json:"servers" required:"true"`
	Type    string       `json:"type,omitempty"`
}

type ServerInfo struct {
	ID string `json:"id" required:"true"`
}

// PowerOptsBuilder allows extensions to add additional parameters to the PowerAction(power on/off and reboot) request.
type PowerOptsBuilder interface {
	ToPowerActionMap(option string) (map[string]interface{}, error)
}

// ToPowerActionMap assembles a request body based on the contents of a PowerOpts and the power option parameter.
func (opts PowerOpts) ToPowerActionMap(option string) (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, option)
}

// PowerAction uses a option parameter to control the power state of the ECS instance.
// The option only supports 'os-start' (power on), 'os-stop' (power off) and 'reboot'.
func PowerAction(client *golangsdk.ServiceClient, opts PowerOptsBuilder, option string) (r cloudservers.JobResult) {
	reqBody, err := opts.ToPowerActionMap(option)
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(actionURL(client), reqBody, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}
