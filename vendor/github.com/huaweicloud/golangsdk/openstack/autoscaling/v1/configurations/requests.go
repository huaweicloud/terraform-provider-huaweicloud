package configurations

import (
	"encoding/base64"
	"log"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type CreateOptsBuilder interface {
	ToConfigurationCreateMap() (map[string]interface{}, error)
}

type CreateOpts struct {
	Name           string             `json:"scaling_configuration_name" required:"true"`
	InstanceConfig InstanceConfigOpts `json:"instance_config" required:"true"`
}

func (opts CreateOpts) ToConfigurationCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] ToConfigurationCreateMap b is: %#v", b)
	log.Printf("[DEBUG] ToConfigurationCreateMap opts is: %#v", opts)
	publicIp := opts.InstanceConfig.PubicIp
	log.Printf("[DEBUG] ToConfigurationCreateMap publicIp is: %#v", publicIp)

	if publicIp != (PublicIpOpts{}) {
		public_ip := map[string]interface{}{}
		eip := map[string]interface{}{}
		bandwidth := map[string]interface{}{}
		eip_raw := publicIp.Eip
		log.Printf("[DEBUG] ToConfigurationCreateMap eip_raw is: %#v", eip_raw)
		if eip_raw != (EipOpts{}) {
			if eip_raw.IpType != "" {
				eip["ip_type"] = eip_raw.IpType
			}
			bandWidth := eip_raw.Bandwidth
			if bandWidth != (BandwidthOpts{}) {
				if bandWidth.Size > 0 {
					bandwidth["size"] = bandWidth.Size
				}
				if bandWidth.ChargingMode != "" {
					bandwidth["charging_mode"] = bandWidth.ChargingMode
				}
				if bandWidth.ShareType != "" {
					bandwidth["share_type"] = bandWidth.ShareType
				}
				eip["bandwidth"] = bandwidth
			}
			public_ip["eip"] = eip
		}
		b["instance_config"].(map[string]interface{})["public_ip"] = public_ip
	}
	if opts.InstanceConfig.UserData != nil {
		var userData string
		if _, err := base64.StdEncoding.DecodeString(string(opts.InstanceConfig.UserData)); err != nil {
			userData = base64.StdEncoding.EncodeToString(opts.InstanceConfig.UserData)
		} else {
			userData = string(opts.InstanceConfig.UserData)
		}
		b["instance_config"].(map[string]interface{})["user_data"] = &userData
	}
	log.Printf("[DEBUG] ToConfigurationCreateMap b is: %#v", b)
	return b, nil
}

//InstanceConfigOpts is an inner struct of CreateOpts
type InstanceConfigOpts struct {
	ID          string            `json:"instance_id,omitempty"`
	FlavorRef   string            `json:"flavorRef,omitempty"`
	ImageRef    string            `json:"imageRef,omitempty"`
	Disk        []DiskOpts        `json:"disk,omitempty"`
	SSHKey      string            `json:"key_name,omitempty"`
	Personality []PersonalityOpts `json:"personality,omitempty"`
	PubicIp     PublicIpOpts      `json:"-"`
	// UserData contains configuration information or scripts to use upon launch.
	// Create will base64-encode it for you, if it isn't already.
	UserData []byte                 `json:"-"`
	Metadata map[string]interface{} `json:"metadata,omitempty"` //TODO not sure the type
}

//DiskOpts is an inner struct of InstanceConfigOpts
type DiskOpts struct {
	Size       int               `json:"size" required:"true"`
	VolumeType string            `json:"volume_type" required:"true"`
	DiskType   string            `json:"disk_type" required:"true"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type PersonalityOpts struct {
	Path    string `json:"path" required:"true"`
	Content string `json:"content" required:"true"`
}

type PublicIpOpts struct {
	Eip EipOpts `json:"-"`
}

type EipOpts struct {
	IpType    string
	Bandwidth BandwidthOpts `json:"-"`
}

type BandwidthOpts struct {
	Size         int
	ShareType    string
	ChargingMode string
}

//Create is a method by which can be able to access to create a configuration
//of autoscaling
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToConfigurationCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

//Get is a method by which can be able to access to get a configuration of
//autoscaling detailed information
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

//Delete
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

type ListOptsBuilder interface {
	ToConfigurationListQuery() (string, error)
}

type ListOpts struct {
	Name    string `q:"scaling_configuration_name"`
	ImageID string `q:"image_id"`
}

func (opts ListOpts) ToConfigurationListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

//List is method that can be able to list all configurations of autoscaling service
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToConfigurationListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ConfigurationPage{pagination.SinglePageBase(r)}
	})
}
