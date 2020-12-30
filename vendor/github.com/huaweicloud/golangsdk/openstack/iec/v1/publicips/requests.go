package publicips

import (
	"net/http"

	"github.com/huaweicloud/golangsdk"
)

type PublicIPRequest struct {
	// Specifies the type of the elastic IP address. The value can the
	// 5_telcom, 5_union, 5_bgp, or 5_sbgp. China Northeast: 5_telcom and 5_union China
	// South: 5_sbgp China East: 5_sbgp China North: 5_bgp and 5_sbgp The value must be a
	// type supported by the system. The value can be 5_telcom, 5_union, or 5_bgp.
	Type string `json:"type"`

	// Specifies the elastic IP address to be obtained. The value must
	// be a valid IP address in the available IP address segment.
	IpAddress string `json:"ip_address,omitempty"`

	//Value range: 4, 6, respectively, to create ipv4 and ipv6, when not created ipv4 by default
	IPVersion string `json:"ip_version,omitempty"`

	SiteID string `json:"site_id,omitempty"`
}

type BandWidth struct {
	// Specifies the bandwidth name. The value is a string of 1 to 64
	// characters that can contain letters, digits, underscores (_), and hyphens (-). This
	// parameter is mandatory when share_type is set to PER and is optional when share_type
	// is set to WHOLE with an ID specified.
	Name string `json:"name,omitempty"`

	// Specifies the bandwidth size. The value ranges from 1 Mbit/s to
	// 300 Mbit/s. This parameter is mandatory when share_type is set to PER and is optional
	// when share_type is set to WHOLE with an ID specified.
	Size int `json:"size,omitempty"`

	// Specifies the ID of the bandwidth. You can specify an earlier
	// shared bandwidth when applying for an elastic IP address for the bandwidth whose type
	// is set to WHOLE. The bandwidth whose type is set to WHOLE exclusively uses its own
	// ID. The value can be the ID of the bandwidth whose type is set to WHOLE.
	ID string `json:"id,omitempty"`

	// Specifies whether the bandwidth is shared or exclusive. The
	// value can be PER or WHOLE.
	ShareType string `json:"share_type"`

	// Specifies the charging mode (by traffic or by bandwidth). The
	// value can be bandwidth or traffic. If the value is an empty character string or no
	// value is specified, default value bandwidth is used.
	ChargeMode string `json:"charge_mode,omitempty"`
}

type CreateOpts struct {
	// Specifies the elastic IP address objects.
	Publicip PublicIPRequest `json:"publicip" required:"true"`

	// Specifies the bandwidth objects.
	Bandwidth BandWidth `json:"bandwidth"`
}

type CreateOptsBuilder interface {
	ToCreatePublicIPMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToCreatePublicIPMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCreatePublicIPMap()
	if err != nil {
		r.Err = err
		return
	}

	createURL := CreateURL(client)

	var resp *http.Response
	resp, r.Err = client.Post(createURL, b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{http.StatusOK}})
	if r.Err != nil {
		return
	}
	defer resp.Body.Close()

	return
}

func Delete(client *golangsdk.ServiceClient, publicipId string) (r DeleteResult) {
	deleteURL := DeleteURL(client, publicipId)

	var resp *http.Response
	resp, r.Err = client.Delete(deleteURL, &golangsdk.RequestOpts{OkCodes: []int{http.StatusNoContent}})
	if r.Err != nil {
		return
	}
	defer resp.Body.Close()

	return
}

func Get(client *golangsdk.ServiceClient, publicipId string) (r GetResult) {
	getURL := GetURL(client, publicipId)

	var resp *http.Response
	resp, r.Err = client.Get(getURL, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{http.StatusOK}})
	if r.Err != nil {
		return
	}
	defer resp.Body.Close()

	return
}

type UpdateOpts struct {
	// Specifies the port ID.
	PortId string `json:"port_id,omitempty"`

	//Value range: 4, 6, respectively, to create ipv4 and ipv6, when not created ipv4 by default
	IPVersion int `json:"ip_version,omitempty"`
}

type UpdateOptsBuilder interface {
	ToUpdatePublicIPMap() (map[string]interface{}, error)
}

func (opts UpdateOpts) ToUpdatePublicIPMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "publicip")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Update(client *golangsdk.ServiceClient, publicipId string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToUpdatePublicIPMap()
	if err != nil {
		r.Err = err
		return
	}
	updateURL := UpdateURL(client, publicipId)

	var resp *http.Response
	resp, r.Err = client.Put(updateURL, b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{http.StatusOK}})
	if r.Err != nil {
		return
	}
	defer resp.Body.Close()

	return
}
