package cloudvolumes

import (
	"github.com/chnsz/golangsdk"
)

type QoSModifyOpts struct {
	IopsAndThroughputOpts IopsAndThroughputOpts `json:"qos_modify" required:"true"`
}

type IopsAndThroughputOpts struct {
	Iops       int `json:"iops" required:"true"`
	Throughput int `json:"throughput,omitempty"`
}

func (opts QoSModifyOpts) ToVolumeUpdateQoSMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

type UpdateQoSOptsBuilder interface {
	ToVolumeUpdateQoSMap() (map[string]interface{}, error)
}

func ModifyQoS(client *golangsdk.ServiceClient, id string, opts UpdateQoSOptsBuilder) (r JobResult) {
	b, err := opts.ToVolumeUpdateQoSMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(qoSURL(client, id), b, &r.Body, nil)
	return
}
