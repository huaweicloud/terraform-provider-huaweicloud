package addons

import (
	"reflect"

	"github.com/chnsz/golangsdk"
)

// ListOpts allows the filtering of list data using given parameters.
type ListOpts struct {
	AddonTemplateName string `json:"addonTemplateName"`
	Uid               string `json:"uid"`
	Version           string `json:"version"`
	Status            string `json:"status"`
}

// List returns collection of addons.
func List(client *golangsdk.ServiceClient, clusterID string, opts ListOpts) ([]Addon, error) {
	var r ListResult
	_, r.Err = client.Get(resourceListURL(client, clusterID), &r.Body, nil)

	allAddons, err := r.ExtractAddon()

	if err != nil {
		return nil, err
	}

	return FilterAddons(allAddons, opts), nil
}

func FilterAddons(addons []Addon, opts ListOpts) []Addon {

	var refinedAddons []Addon
	var matched bool

	m := map[string]FilterStruct{}

	if opts.AddonTemplateName != "" {
		m["AddonTemplateName"] = FilterStruct{Value: opts.AddonTemplateName, Driller: []string{"Spec"}}
	}
	if opts.Version != "" {
		m["Version"] = FilterStruct{Value: opts.Version, Driller: []string{"Spec"}}
	}
	if opts.Uid != "" {
		m["Id"] = FilterStruct{Value: opts.Uid, Driller: []string{"Metadata"}}
	}
	if opts.Status != "" {
		m["Status"] = FilterStruct{Value: opts.Status, Driller: []string{"Status"}}
	}

	if len(m) > 0 && len(addons) > 0 {
		for _, addon := range addons {
			matched = true

			for key, value := range m {
				if sVal := GetStructNestedField(&addon, key, value.Driller); sVal != value.Value {
					matched = false
					break
				}
			}
			if matched {
				refinedAddons = append(refinedAddons, addon)
			}
		}
	}

	return refinedAddons
}

func GetStructNestedField(v *Addon, field string, structDriller []string) string {
	r := reflect.ValueOf(v)
	for _, drillField := range structDriller {
		f := reflect.Indirect(r).FieldByName(drillField).Interface()
		r = reflect.ValueOf(f)
	}
	f1 := reflect.Indirect(r).FieldByName(field)
	return string(f1.String())
}

type FilterStruct struct {
	Value   string
	Driller []string
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToAddonCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new addon
type CreateOpts struct {
	// API type, fixed value Addon
	Kind string `json:"kind" required:"true"`
	// API version, fixed value v3
	ApiVersion string `json:"apiVersion" required:"true"`
	// Metadata required to create an addon
	Metadata CreateMetadata `json:"metadata" required:"true"`
	// specifications to create an addon
	Spec RequestSpec `json:"spec" required:"true"`
}

type CreateMetadata struct {
	Anno Annotations `json:"annotations" required:"true"`
}

type Annotations struct {
	AddonInstallType string `json:"addon.install/type" required:"true"`
}

//Specifications to create an addon
type RequestSpec struct {
	// For the addon version.
	Version string `json:"version" required:"true"`
	// Cluster ID.
	ClusterID string `json:"clusterID" required:"true"`
	// Addon Template Name.
	AddonTemplateName string `json:"addonTemplateName" required:"true"`
	// Addon Parameters
	Values Values `json:"values" required:"true"`
}

type Values struct {
	Basic  map[string]interface{} `json:"basic" required:"true"`
	Custom map[string]interface{} `json:"custom,omitempty"`
	Flavor map[string]interface{} `json:"flavor,omitempty"`
}

// ToAddonCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToAddonCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// addon.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder, cluster_id string) (r CreateResult) {
	b, err := opts.ToAddonCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{201}}
	_, r.Err = c.Post(rootURL(c, cluster_id), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular addon based on its unique ID.
func Get(c *golangsdk.ServiceClient, id, cluster_id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id, cluster_id), &r.Body, nil)
	return
}

// Delete will permanently delete a particular addon based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id, cluster_id string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Delete(resourceURL(c, id, cluster_id), reqOpt)
	return
}
