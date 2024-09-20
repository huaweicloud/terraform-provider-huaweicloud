package autoscaling

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

type UpdateBuilder interface {
	ToUpdateMap() (map[string]interface{}, error)
}

type UpdateAutoScalingOpts struct {
	Status           string           `json:"status" required:"true"`
	ScalingStrategy  *ScalingStrategy `json:"scaling_strategy" required:"true"`
	MonitorCycle     int              `json:"monitor_cycle,omitempty"`
	SilenceCycle     int              `json:"silence_cycle,omitempty"`
	EnlargeThreshold int              `json:"enlarge_threshold,omitempty"`
	MaxFlavor        string           `json:"max_flavor,omitempty"`
	ReduceEnabled    bool             `json:"reduce_enabled,omitempty"`
	MaxReadOnlyCount int              `json:"max_read_only_count,omitempty"`
	ReadOnlyWeight   int              `json:"read_only_weight,omitempty"`
}

type ScalingStrategy struct {
	FlavorSwitch   string `json:"flavor_switch" required:"true"`
	ReadOnlySwitch string `json:"read_only_switch" required:"true"`
}

func (opts UpdateAutoScalingOpts) ToUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Update(c *golangsdk.ServiceClient, instanceId string, opts UpdateBuilder) (r UpdateResult) {
	b, err := opts.ToUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Put(updateURL(c, instanceId), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

func Get(c *golangsdk.ServiceClient, instanceId string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, instanceId), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}
