package cloudvolumes

import (
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToVolumeCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains options for creating a Volume. This object is passed to
// the cloudvolumes.Create function.
type CreateOpts struct {
	Volume     VolumeOpts     `json:"volume" required:"true"`
	ChargeInfo *BssParam      `json:"bssParam,omitempty"`
	Scheduler  *SchedulerOpts `json:"OS-SCH-HNT:scheduler_hints,omitempty"`
	ServerID   string         `json:"server_id,omitempty"`
}

type BssParam struct {
	// Specifies the billing mode. The default value is postPaid.
	//   prePaid: indicates the yearly/monthly billing mode.
	//   postPaid: indicates the pay-per-use billing mode.
	ChargingMode string `json:"chargingMode" required:"true"`
	// Specifies the unit of the subscription term.
	// This parameter is valid and mandatory only when chargingMode is set to prePaid.
	//   month: indicates that the unit is month.
	//   year: indicates that the unit is year.
	PeriodType string `json:"periodType,omitempty"`
	// Specifies the subscription term. This parameter is valid and mandatory only when chargingMode is set to prePaid.
	//   When periodType is set to month, the parameter value ranges from 1 to 9.
	//   When periodType is set to year, the parameter value must be set to 1.
	PeriodNum int `json:"periodNum,omitempty"`
	// Specifies whether to pay immediately. This parameter is valid only when chargingMode is set to prePaid. The default value is false.
	//   false: indicates not to pay immediately after an order is created.
	//   true: indicates to pay immediately after an order is created. The system will automatically deduct fees from the account balance.
	IsAutoPay string `json:"isAutoPay,omitempty"`
	// Specifies whether to automatically renew the subscription.
	// This parameter is valid only when chargingMode is set to prePaid. The default value is false.
	//   false: indicates not to automatically renew the subscription.
	//   true: indicates to automatically renew the subscription. The automatic renewal term is the same as the subscription term.
	IsAutoRenew string `json:"isAutoRenew,omitempty"`
}

// VolumeOpts contains options for creating a Volume.
type VolumeOpts struct {
	// The availability zone
	AvailabilityZone string `json:"availability_zone" required:"true"`
	// The associated volume type
	VolumeType string `json:"volume_type" required:"true"`
	// The volume name
	Name string `json:"name,omitempty"`
	// The volume description
	Description string `json:"description,omitempty"`
	// The size of the volume, in GB
	Size int `json:"size,omitempty"`
	// The number to be created in a batch
	Count int `json:"count,omitempty"`
	// The backup_id
	BackupID string `json:"backup_id,omitempty"`
	// the ID of the existing volume snapshot
	SnapshotID string `json:"snapshot_id,omitempty"`
	// the ID of the image in IMS
	ImageID string `json:"imageRef,omitempty"`
	// Shared disk
	Multiattach bool `json:"multiattach,omitempty"`
	// One or more metadata key and value pairs to associate with the volume
	Metadata map[string]string `json:"metadata,omitempty"`
	// One or more tag key and value pairs to associate with the volume
	Tags map[string]string `json:"tags,omitempty"`
	// the enterprise project id
	EnterpriseProjectID string `json:"enterprise_project_id,omitempty"`
	// The iops of evs volume. Only required when volume_type is `GPSSD2` or `ESSD2`
	IOPS int `json:"iops,omitempty"`
	// The throughput of evs volume. Only required when volume_type is `GPSSD2`
	Throughput int `json:"throughput,omitempty"`
}

// SchedulerOpts contains the scheduler hints
type SchedulerOpts struct {
	StorageID string `json:"dedicated_storage_id,omitempty"`
}

// ToVolumeCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToVolumeCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new Volume based on the values in CreateOpts.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r JobResult) {
	b, err := opts.ToVolumeCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	// the version of create API is v2.1
	newClient := *client
	baseURL := newClient.ResourceBaseURL()
	newClient.ResourceBase = strings.Replace(baseURL, "/v2/", "/v2.1/", 1)

	_, r.Err = newClient.Post(createURL(&newClient), b, &r.Body, nil)
	return
}

// ExtendOptsBuilder allows extensions to add additional parameters to the
// ExtendSize request.
type ExtendOptsBuilder interface {
	ToVolumeExtendMap() (map[string]interface{}, error)
}

// ExtendOpts contains options for extending the size of an existing Volume.
// This object is passed to the cloudvolumes.ExtendSize function.
type ExtendOpts struct {
	SizeOpts   ExtendSizeOpts    `json:"os-extend" required:"true"`
	ChargeInfo *ExtendChargeOpts `json:"bssParam,omitempty"`
}

// ExtendSizeOpts contains the new size of the volume, in GB.
type ExtendSizeOpts struct {
	NewSize int `json:"new_size" required:"true"`
}

// ExtendChargeOpts contains the charging parameters of the volume
type ExtendChargeOpts struct {
	IsAutoPay string `json:"isAutoPay,omitempty"`
}

// ToVolumeExtendMap assembles a request body based on the contents of an
// ExtendOpts.
func (opts ExtendOpts) ToVolumeExtendMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// ExtendSize will extend the size of the volume based on the provided information.
// This operation does not return a response body.
func ExtendSize(client *golangsdk.ServiceClient, id string, opts ExtendOptsBuilder) (r JobResult) {
	b, err := opts.ToVolumeExtendMap()
	if err != nil {
		r.Err = err
		return
	}
	// the version of extend API is v2.1
	newClient := *client
	baseURL := newClient.ResourceBaseURL()
	newClient.ResourceBase = strings.Replace(baseURL, "/v2/", "/v2.1/", 1)

	_, r.Err = newClient.Post(actionURL(&newClient, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})
	return
}

type RetypeOptsBuilder interface {
	ToVolumeRetypeMap() (map[string]interface{}, error)
}

type RetypeOpts struct {
	BssParam *BssParamOpts `json:"bssParam,omitempty"`
	OSRetype OSRetypeOpts  `json:"os-retype" required:"true"`
}

type BssParamOpts struct {
	IsAutoPay string `json:"isAutoPay,omitempty"`
}

type OSRetypeOpts struct {
	NewType    string `json:"new_type" required:"true"`
	Iops       int    `json:"iops,omitempty"`
	Throughput int    `json:"throughput,omitempty"`
}

func (opts RetypeOpts) ToVolumeRetypeMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func UpdateVolumeType(client *golangsdk.ServiceClient, id string, opts RetypeOptsBuilder) (r JobResult) {
	b, err := opts.ToVolumeRetypeMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(retypeURL(client, id), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToVolumeUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contain options for updating an existing Volume. This object is passed
// to the cloudvolumes.Update function.
type UpdateOpts struct {
	Name        string  `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ToVolumeUpdateMap assembles a request body based on the contents of an
// UpdateOpts.
func (opts UpdateOpts) ToVolumeUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "volume")
}

// Update will update the Volume with provided information. To extract the updated
// Volume from the response, call the Extract method on the UpdateResult.
func Update(client *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToVolumeUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// DeleteOptsBuilder is an interface by which can be able to build the query string
// of volume deletion.
type DeleteOptsBuilder interface {
	ToVolumeDeleteQuery() (string, error)
}

// DeleteOpts contain options for deleting an existing Volume. This object is passed
// to the cloudvolumes.Delete function.
type DeleteOpts struct {
	// Specifies to delete all snapshots associated with the EVS disk.
	Cascade bool `q:"cascade"`
}

// ToVolumeDeleteQuery assembles a request body based on the contents of an
// DeleteOpts.
func (opts DeleteOpts) ToVolumeDeleteQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// Delete will delete the existing Volume with the provided ID
func Delete(client *golangsdk.ServiceClient, id string, opts DeleteOptsBuilder) (r DeleteResult) {
	url := resourceURL(client, id)
	if opts != nil {
		q, err := opts.ToVolumeDeleteQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += q
	}
	_, r.Err = client.Delete(url, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get retrieves the Volume with the provided ID. To extract the Volume object
// from the response, call the Extract method on the GetResult.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToVolumeListQuery() (string, error)
}

// ListOpts holds options for listing Volumes. It is passed to the volumes.List
// function.
type ListOpts struct {
	// Name will filter by the specified volume name.
	Name string `q:"name"`
	// Status will filter by the specified status.
	Status string `q:"status"`
	// Metadata will filter results based on specified metadata.
	Metadata string `q:"metadata"`
	// Specifies the disk ID.
	ID string `q:"id"`
	// Specifies the disk IDs. The parameter value is in the ids=["id1","id2",...,"idx"] format.
	// In the response, the ids value contains valid disk IDs only. Invalid disk IDs will be ignored.
	// Details about a maximum of 60 disks can be queried.
	// If parameters id and ids are both specified in the request, id will be ignored.
	IDs string `q:"ids"`
	// Specifies the AZ.
	AvailabilityZone string `q:"availability_zone"`
	// Specifies the ID of the DSS storage pool. All disks in the DSS storage pool can be filtered out.
	// Only precise match is supported.
	DedicatedStorageID string `q:"dedicated_storage_id"`
	// Specifies the name of the DSS storage pool. All disks in the DSS storage pool can be filtered out.
	// Fuzzy match is supported.
	DedicatedStorageName string `q:"dedicated_storage_name"`
	// Specifies the enterprise project ID for filtering. If input parameter all_granted_eps exists, disks in all
	// enterprise projects that are within the permission scope will be queried.
	EnterpriseProjectID string `q:"enterprise_project_id"`
	// Specifies whether the disk is shareable.
	//   true: specifies a shared disk.
	//   false: specifies a non-shared disk.
	Multiattach bool `q:"multiattach"`
	// Specifies the service type. Currently, the supported services are EVS, DSS, and DESS.
	ServiceType string `q:"service_type"`
	// Specifies the server ID.
	// This parameter is used to filter all the EVS disks that have been attached to this server.
	ServerID string `q:"server_id"`
	// Specifies the keyword based on which the returned results are sorted.
	// The value can be id, status, size, or created_at, and the default value is created_at.
	SortKey string `q:"sort_key"`
	// Specifies the result sorting order. The default value is desc.
	//   desc: indicates the descending order.
	//   asc: indicates the ascending order.
	SortDir string `q:"sort_dir"`
	// Specifies the disk type ID.
	// You can obtain the disk type ID in Querying EVS Disk Types.
	// That is, the id value in the volume_types parameter description table.
	VolumeTypeID string `q:"volume_type_id"`
	// Requests a page size of items.
	Limit int `q:"limit"`
	// Used in conjunction with limit to return a slice of items.
	Offset int `q:"offset"`
	// The ID of the last-seen item.
	Marker string `q:"marker"`
}

// ToVolumeListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToVolumeListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List returns Volumes optionally limited by the conditions provided in ListOpts.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToVolumeListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := VolumePage{pagination.OffsetPageBase{PageResult: r}}
		return p
	})
}

func ListPage(client *golangsdk.ServiceClient, opts ListOptsBuilder) ([]Volume, error) {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToVolumeListQuery()
		if err != nil {
			return nil, err
		}
		url += query
	}

	var rst golangsdk.Result
	_, err := client.Get(url, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})

	if err != nil {
		return nil, err
	}

	var r struct {
		Volumes []Volume
	}
	err = rst.ExtractInto(&r)
	return r.Volumes, err
}
