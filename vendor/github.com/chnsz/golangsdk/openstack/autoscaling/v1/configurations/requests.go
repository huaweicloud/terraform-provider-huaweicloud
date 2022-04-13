package configurations

import (
	"encoding/base64"
	"log"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
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
	InstanceID  string            `json:"instance_id,omitempty"`
	FlavorRef   string            `json:"flavorRef,omitempty"`
	ImageRef    string            `json:"imageRef,omitempty"`
	Disk        []DiskOpts        `json:"disk,omitempty"`
	SSHKey      string            `json:"key_name,omitempty"`
	Personality []PersonalityOpts `json:"personality,omitempty"`
	PubicIP     *PublicIpOpts     `json:"public_ip,omitempty"`
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
	Eip EipOpts `json:"eip" required:"true"`
}

type EipOpts struct {
	Type      string        `json:"ip_type" required:"true"`
	Bandwidth BandwidthOpts `json:"bandwidth" required:"true"`
}

type BandwidthOpts struct {
	ShareType    string `json:"share_type" required:"true"`
	Size         int    `json:"size,omitempty"`
	ChargingMode string `json:"charging_mode,omitempty"`
	ID           string `json:"id,omitempty"`
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
