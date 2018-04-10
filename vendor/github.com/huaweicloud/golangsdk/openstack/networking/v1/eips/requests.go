package eips

import (
	"github.com/huaweicloud/golangsdk"
)

//ApplyOptsBuilder is an interface by which can build the request body of public ip
//application
type ApplyOptsBuilder interface {
	ToPublicIpApplyMap() (map[string]interface{}, error)
}

//ApplyOpts is a struct which is used to create public ip
type ApplyOpts struct {
	IP        PublicIpOpts  `json:"publicip" required:"true"`
	Bandwidth BandwidthOpts `json:"bandwidth" required:"true"`
}

type PublicIpOpts struct {
	Type    string `json:"type" required:"true"`
	Address string `json:"ip_address,omitempty"`
}

type BandwidthOpts struct {
	Name       string `json:"name" required:"true"`
	Size       int    `json:"size" required:"true"`
	ShareType  string `json:"share_type" required:"true"`
	ChargeMode string `json:"charge_mode,omitempty"`
}

func (opts ApplyOpts) ToPublicIpApplyMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

//Apply is a method by which can access to apply the public ip
func Apply(client *golangsdk.ServiceClient, opts ApplyOptsBuilder) (r ApplyResult) {
	b, err := opts.ToPublicIpApplyMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

//Get is a method by which can get the detailed information of public ip
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)
	return
}

//Delete is a method by which can be able to delete a private ip
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, id), nil)
	return
}

//UpdateOptsBuilder is an interface by which can be able to build the request
//body
type UpdateOptsBuilder interface {
	ToPublicIpUpdateMap() (map[string]interface{}, error)
}

//UpdateOpts is a struct which represents the request body of update method
type UpdateOpts struct {
	PortID string `json:"port_id,omitempty"`
}

func (opts UpdateOpts) ToPublicIpUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "publicip")
}

//Update is a method which can be able to update the port of public ip
func Update(client *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToPublicIpUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
