package partitions

import (
	"reflect"

	"github.com/chnsz/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json"},
}

// ListOpts allows the filtering of list data using given parameters.
type ListOpts struct {
	Name string `json:"name"`
}

// List returns collection of partitions.
func List(client *golangsdk.ServiceClient, clusterID string, opts ListOpts) ([]Partitions, error) {
	var r ListResult
	_, r.Err = client.Get(rootURL(client, clusterID), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})

	allPartitions, err := r.ExtractPartitions()

	if err != nil {
		return nil, err
	}

	return FilterPartitions(allPartitions, opts), nil
}

func FilterPartitions(partitions []Partitions, opts ListOpts) []Partitions {
	var refinedPartitions []Partitions
	var matched bool

	m := map[string]FilterStruct{}

	if opts.Name != "" {
		m["Name"] = FilterStruct{Value: opts.Name, Driller: []string{"Metadata"}}
	}

	if len(m) > 0 && len(partitions) > 0 {
		for _, partition := range partitions {
			matched = true

			for key, value := range m {
				if sVal := GetStructNestedField(&partition, key, value.Driller); !(sVal == value.Value) {
					matched = false
				}
			}
			if matched {
				refinedPartitions = append(refinedPartitions, partition)
			}
		}
	} else {
		refinedPartitions = partitions
	}
	return refinedPartitions
}

func GetStructNestedField(v *Partitions, field string, structDriller []string) string {
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
type CreateOpts struct {
	// API type, fixed value Partition
	Kind string `json:"kind" required:"true"`
	// API version, fixed value v3
	ApiVersion string `json:"apiversion" required:"true"`
	// Metadata required to create a Partition
	Metadata CreateMetaData `json:"metadata"`
	// specifications to create a Partition
	Spec Spec `json:"spec" required:"true"`
}

// Metadata required to create a Partition
type CreateMetaData struct {
	// Partition name
	Name string `json:"name,omitempty"`
	// Partition tag, key value pair format
	Labels map[string]string `json:"labels,omitempty"`
	// Partition annotation, key value pair format
	Annotations map[string]string `json:"annotations,omitempty"`
}

// Create accepts a CreateOpts struct and uses the values to create a new
// logical Partition. When it is created, the Partition does not have an internal
// interface
type CreateOptsBuilder interface {
	ToPartitionCreateMap() (map[string]interface{}, error)
}

// ToPartitionCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToPartitionCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// logical Partition.
func Create(c *golangsdk.ServiceClient, clusterid string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPartitionCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{201}}
	_, r.Err = c.Post(rootURL(c, clusterid), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular partitions based on its unique ID and cluster ID.
func Get(c *golangsdk.ServiceClient, clusterid, partitionName string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, clusterid, partitionName), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToPartitionUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a new partition
type UpdateOpts struct {
	Metadata UpdateMetadata `json:"metadata,omitempty"`
}

type UpdateMetadata struct {
	ContainerNetwork []ContainerNetwork `json:"containerNetwork,omitempty"`
}

// ToPartitionUpdateMap builds an update body based on UpdateOpts.
func (opts UpdateOpts) ToPartitionUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update allows partitions to be updated.
func Update(c *golangsdk.ServiceClient, clusterid, partitionName string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToPartitionUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, clusterid, partitionName), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular partition based on its unique Name and cluster ID.
func Delete(c *golangsdk.ServiceClient, clusterid, partitionName string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, clusterid, partitionName), &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}
