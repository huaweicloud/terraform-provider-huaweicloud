package softwareconfig

import (
	"reflect"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the software config attributes you want to see returned. Marker and Limit are used for pagination.
type ListOpts struct {
	Id     string
	Name   string
	Marker string `q:"marker"`
	Limit  int    `q:"limit"`
}

// List returns collection of
// Software Config. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those Software Config that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]SoftwareConfig, error) {
	q, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return nil, err
	}
	u := rootURL(c) + q.String()
	pages, err := pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return SoftwareConfigPage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()

	allConfigs, err := ExtractSoftwareConfigs(pages)
	if err != nil {
		return nil, err
	}

	return FilterSoftwareConfig(allConfigs, opts)
}

func FilterSoftwareConfig(config []SoftwareConfig, opts ListOpts) ([]SoftwareConfig, error) {

	var refinedSoftwareConfig []SoftwareConfig
	var matched bool
	m := map[string]interface{}{}

	if opts.Id != "" {
		m["Id"] = opts.Id
	}
	if opts.Name != "" {
		m["Name"] = opts.Name
	}

	if len(m) > 0 && len(config) > 0 {
		for _, config := range config {
			matched = true

			for key, value := range m {
				if sVal := getStructField(&config, key); !(sVal == value) {
					matched = false
				}
			}

			if matched {
				refinedSoftwareConfig = append(refinedSoftwareConfig, config)
			}
		}
	} else {
		refinedSoftwareConfig = config
	}
	return refinedSoftwareConfig, nil
}

func getStructField(v *SoftwareConfig, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToSoftwareConfigCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new Software Config. There are
// no required values.
type CreateOpts struct {
	// Specifies the script used for defining the configuration.
	Config string `json:"config,omitempty"`
	//Specifies the name of the software configuration group.
	Group string `json:"group,omitempty"`
	//Specifies the name of the software configuration.
	Name string `json:"name" required:"true"`
	//Specifies the software configuration input.
	Inputs []map[string]interface{} `json:"inputs,omitempty"`
	//Specifies the software configuration output.
	Outputs []map[string]interface{} `json:"outputs,omitempty"`
	//Specifies options used by a software configuration management tool.
	Options map[string]interface{} `json:"options,omitempty"`
}

// ToSoftwareConfigCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToSoftwareConfigCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and uses the values to create a new Software config
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSoftwareConfigCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular software config based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// Delete will permanently delete a particular Software Config based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
