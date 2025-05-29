package clusters

import (
	"reflect"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json"},
}

// ListOpts allows the filtering of list data using given parameters.
type ListOpts struct {
	Name                string `json:"name"`
	ID                  string `json:"uuid"`
	Type                string `json:"type"`
	VpcID               string `json:"vpc"`
	Phase               string `json:"phase"`
	EnterpriseProjectID string `json:"enterpriseProjectId"`
}

// List returns collection of clusters.
func List(client *golangsdk.ServiceClient, opts ListOpts) ([]Clusters, error) {
	var r ListResult
	_, r.Err = client.Get(rootURL(client), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})

	allClusters, err := r.ExtractClusters()
	if err != nil {
		return nil, err
	}

	return FilterClusters(allClusters, opts), nil
}

func FilterClusters(clusters []Clusters, opts ListOpts) []Clusters {

	var refinedClusters []Clusters
	var matched bool
	m := map[string]FilterStruct{}

	if opts.Name != "" {
		m["Name"] = FilterStruct{Value: opts.Name, Driller: []string{"Metadata"}}
	}
	if opts.ID != "" {
		m["Id"] = FilterStruct{Value: opts.ID, Driller: []string{"Metadata"}}
	}
	if opts.Type != "" {
		m["Type"] = FilterStruct{Value: opts.Type, Driller: []string{"Spec"}}
	}
	if opts.VpcID != "" {
		m["VpcId"] = FilterStruct{Value: opts.VpcID, Driller: []string{"Spec", "HostNetwork"}}
	}
	if opts.Phase != "" {
		m["Phase"] = FilterStruct{Value: opts.Phase, Driller: []string{"Status"}}
	}
	if opts.EnterpriseProjectID != "" {
		m["enterpriseProjectId"] = FilterStruct{Value: opts.EnterpriseProjectID, Driller: []string{"Spec", "ExtendParam"}}
	}

	if len(m) > 0 && len(clusters) > 0 {
		for _, cluster := range clusters {
			matched = true

			for key, value := range m {
				if sVal := GetStructNestedField(&cluster, key, value.Driller); !(sVal == value.Value) {
					matched = false
				}
			}
			if matched {
				refinedClusters = append(refinedClusters, cluster)
			}
		}

	} else {
		refinedClusters = clusters
	}

	return refinedClusters
}

type FilterStruct struct {
	Value   string
	Driller []string
}

func GetStructNestedField(v *Clusters, field string, structDriller []string) string {
	r := reflect.ValueOf(v)
	for _, drillField := range structDriller {
		f := reflect.Indirect(r).FieldByName(drillField).Interface()
		r = reflect.ValueOf(f)
	}

	if r.Kind() == reflect.Map {
		keys := r.MapKeys()
		for _, k := range keys {
			if k.String() == field {
				f1 := r.MapIndex(k)
				return f1.Interface().(string)
			}
		}
		return ""
	} else {
		f1 := reflect.Indirect(r).FieldByName(field)
		return f1.String()
	}
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToClusterCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new cluster
type CreateOpts struct {
	// API type, fixed value Cluster
	Kind string `json:"kind" required:"true"`
	// API version, fixed value v3
	ApiVersion string `json:"apiversion" required:"true"`
	// Metadata required to create a cluster
	Metadata CreateMetaData `json:"metadata" required:"true"`
	// specifications to create a cluster
	Spec Spec `json:"spec" required:"true"`
}

// Metadata required to create a cluster
type CreateMetaData struct {
	// Cluster unique name
	Name string `json:"name" required:"true"`
	// Cluster tag, key/value pair format
	Labels map[string]string `json:"labels,omitempty"`
	// Cluster annotation, key/value pair format
	Annotations map[string]string `json:"annotations,omitempty"`
	// Cluster alias
	Alias string `json:"alias,omitempty"`
	// Cluster timezone
	Timezone string `json:"timezone,omitempty"`
}

// ToClusterCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToClusterCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// logical cluster.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToClusterCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{201}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular cluster based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

// GetCert retrieves a particular cluster certificate based on its unique ID.
func GetCert(c *golangsdk.ServiceClient, id string, opts GetCertOpts) (r GetCertResult) {
	body, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(certificateURL(c, id), body, &r.Body, nil)
	return
}

type GetCertOpts struct {
	Duration int `json:"duration" required:"true"`
}

// UpdateOpts contains all the values needed to update a new cluster
type UpdateOpts struct {
	Spec     UpdateSpec      `json:"spec" required:"true"`
	Metadata *UpdateMetadata `json:"metadata,omitempty"`
}

type UpdateMetadata struct {
	// Cluster alias
	Alias string `json:"alias"`
}

type UpdateSpec struct {
	// Cluster description
	Description string `json:"description,omitempty"`
	// Custom san list for certificates
	CustomSan []string `json:"customSan,omitempty"`
	//Container network parameters
	ContainerNetwork *UpdateContainerNetworkSpec `json:"containerNetwork,omitempty"`
	// ENI network parameters
	EniNetwork *EniNetworkSpec `json:"eniNetwork,omitempty"`
	// Node network parameters
	HostNetwork *UpdateHostNetworkSpec `json:"hostNetwork,omitempty"`
}

type UpdateContainerNetworkSpec struct {
	// List of container CIDR blocks. In clusters of v1.21 and later, the cidrs field is used.
	// When the cluster network type is vpc-router, you can add multiple container CIDR blocks.
	// In versions earlier than v1.21, if the cidrs field is used, the first CIDR element in the array is used as the container CIDR block.
	Cidrs []CidrSpec `json:"cidrs,omitempty"`
}

type UpdateHostNetworkSpec struct {
	//The ID of the Security Group used to create the node
	SecurityGroup string `json:"SecurityGroup,omitempty"`
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToClusterUpdateMap() (map[string]interface{}, error)
}

// ToClusterUpdateMap builds an update body based on UpdateOpts.
func (opts UpdateOpts) ToClusterUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update allows clusters to update description.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToClusterUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

type DeleteOpts struct {
	ErrorStatus      string `q:"errorStatus"`
	DeleteEfs        string `q:"delete_efs"`
	DeleteENI        string `q:"delete_eni"`
	DeleteEvs        string `q:"delete_evs"`
	DeleteNet        string `q:"delete_net"`
	DeleteObs        string `q:"delete_obs"`
	DeleteSfs        string `q:"delete_sfs"`
	DeleteSfs30      string `q:"delete_sfs30"`
	LtsReclaimPolicy string `q:"lts_reclaim_policy"`
}

type DeleteOptsBuilder interface {
	ToDeleteQuery() (string, error)
}

func (opts DeleteOpts) ToDeleteQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// DeleteWithOpts will permanently delete a particular cluster based on its unique ID,
// and can delete associated resources based on DeleteOpts.
func DeleteWithOpts(c *golangsdk.ServiceClient, id string, opts DeleteOptsBuilder) (r DeleteResult) {
	url := resourceURL(c, id)
	if opts != nil {
		var query string
		query, r.Err = opts.ToDeleteQuery()
		if r.Err != nil {
			return
		}
		url += query
	}
	_, r.Err = c.Delete(url, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

// Delete will permanently delete a particular cluster based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

type UpdateIpOpts struct {
	Action    string `json:"action" required:"true"`
	Spec      IpSpec `json:"spec,omitempty"`
	ElasticIp string `json:"elasticIp"`
}

type IpSpec struct {
	ID string `json:"id" required:"true"`
}

type UpdateIpOptsBuilder interface {
	ToMasterIpUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateIpOpts) ToMasterIpUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "spec")
}

// Update the access information of a specified cluster.
func UpdateMasterIp(c *golangsdk.ServiceClient, id string, opts UpdateIpOptsBuilder) (r UpdateIpResult) {
	b, err := opts.ToMasterIpUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(masterIpURL(c, id), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

func Operation(c *golangsdk.ServiceClient, id, action string) (r OperationResult) {
	_, r.Err = c.Post(operationURL(c, id, action), nil, nil, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

type ResizeOpts struct {
	FavorResize string             `json:"flavorResize" required:"true"`
	ExtendParam *ResizeExtendParam `json:"extendParam,omitempty"`
}

type ResizeExtendParam struct {
	DecMasterFlavor string `json:"decMasterFlavor,omitempty"`
	IsAutoPay       string `json:"isAutoPay,omitempty"`
}

type ResizeResp struct {
	JobID   string `json:"jobID"`
	OrderID string `json:"orderID"`
}

func Resize(c *golangsdk.ServiceClient, id string, opts ResizeOpts) (ResizeResp, error) {
	var r ResizeResp
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return r, err
	}

	_, err = c.Post(operationURL(c, id, "resize"), b, &r, &golangsdk.RequestOpts{
		OkCodes:     []int{201},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return r, err
}

type UpdateTagsOpts struct {
	Tags []tags.ResourceTag `json:"tags" required:"true"`
}

// AddTags will add tags to the cluster.
func AddTags(c *golangsdk.ServiceClient, id string, tagList []tags.ResourceTag) (r UpdateIpResult) {
	opts := UpdateTagsOpts{
		Tags: tagList,
	}
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(tagsURL(c, id, "create"), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

// RemoveTags will remove tags from the cluster.
func RemoveTags(c *golangsdk.ServiceClient, id string, tagList []tags.ResourceTag) (r UpdateIpResult) {
	tagsWithKeys := make([]tags.ResourceTag, len(tagList))
	for i, v := range tagList {
		tagsWithKeys[i] = tags.ResourceTag{
			Key: v.Key,
		}
	}
	opts := UpdateTagsOpts{
		Tags: tagsWithKeys,
	}
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(tagsURL(c, id, "delete"), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}
