package notebook

import (
	"fmt"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type CreateOpts struct {
	Name        string         `json:"name" required:"true"`
	Flavor      string         `json:"flavor" required:"true"`
	ImageId     string         `json:"image_id" required:"true"`
	Volume      VolumeReq      `json:"volume" required:"true"`
	Description string         `json:"description,omitempty"`
	Duration    *int           `json:"duration,omitempty"`
	Endpoints   []EndpointsReq `json:"endpoints,omitempty"`
	Feature     string         `json:"feature,omitempty"`
	PoolId      string         `json:"pool_id,omitempty"`
	WorkspaceId string         `json:"workspace_id,omitempty"`
}

type EndpointsReq struct {
	AllowedAccessIps []string `json:"allowed_access_ips,omitempty"`
	Service          string   `json:"service,omitempty"`
	KeyPairNames     []string `json:"key_pair_names,omitempty"`
}

type VolumeReq struct {
	Category  string `json:"category" required:"true"`
	Ownership string `json:"ownership" required:"true"`
	Capacity  *int   `json:"capacity,omitempty"`
	Uri       string `json:"uri,omitempty"`
}

type ListOpts struct {
	Feature     string `q:"feature,omitempty"`
	Limit       int    `q:"limit,omitempty"`
	Name        string `q:"name,omitempty"`
	Offset      int    `q:"offset,omitempty"`
	Owner       string `q:"owner,omitempty"`
	SortDir     string `q:"sort_dir,omitempty"`
	SortKey     string `q:"sort_key,omitempty"`
	Status      string `q:"status,omitempty"`
	WorkspaceId string `q:"workspaceId,omitempty"`
}

type UpdateOpts struct {
	Description    *string        `json:"description,omitempty"`
	Endpoints      []EndpointsReq `json:"endpoints,omitempty"`
	Flavor         string         `json:"flavor,omitempty"`
	ImageId        string         `json:"image_id,omitempty"`
	Name           string         `json:"name,omitempty"`
	StorageNewSize *int           `json:"storage_new_size,omitempty"`
}

type ListImageOpts struct {
	Limit       int    `q:"limit,omitempty"`
	Name        string `q:"name,omitempty"`
	Namespace   string `q:"namespace,omitempty"`
	Offset      int    `q:"offset,omitempty"`
	SortDir     string `q:"sort_dir,omitempty"`
	SortKey     string `q:"sort_key,omitempty"`
	Type        string `q:"type,omitempty"`
	WorkspaceId string `q:"workspace_id,omitempty"`
}

type MountStorageOpts struct {
	Category  string `json:"category"`
	MountPath string `json:"mount_path"`
	Uri       string `json:"uri"`
}

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*Notebook, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst Notebook
	_, err = c.Post(createURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func Delete(c *golangsdk.ServiceClient, id string) (*Notebook, error) {
	url := deleteURL(c, id)
	var rst Notebook
	_, err := c.DeleteWithResponse(url, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func List(c *golangsdk.ServiceClient, opts ListOpts) (*ListNotebooks, error) {
	url := listURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst ListNotebooks
	_, err = c.Get(url, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func Get(c *golangsdk.ServiceClient, id string) (*Notebook, error) {
	var rst Notebook
	_, err := c.Get(getURL(c, id), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Update(c *golangsdk.ServiceClient, id string, opts UpdateOpts) (*Notebook, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst Notebook
	_, err = c.Put(updateURL(c, id), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func UpdateLease(c *golangsdk.ServiceClient, id string, duration int) (*Lease, error) {
	url := updateLeaseURL(c, id)
	url += fmt.Sprintf("?duration=%d", duration)
	var rst Lease
	emptyBody := make(map[string]string)
	_, err := c.Patch(url, emptyBody, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func Start(c *golangsdk.ServiceClient, id string, duration int) (*Notebook, error) {
	url := startURL(c, id)
	url += fmt.Sprintf("?duration=%d", duration)

	var rst Notebook
	emptyBody := make(map[string]string)
	_, err := c.Post(url, emptyBody, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func Stop(c *golangsdk.ServiceClient, id string) (*Notebook, error) {
	url := stopURL(c, id)
	opts := make(map[string]interface{})
	var rst Notebook

	_, err := c.Post(url, opts, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func ListImages(c *golangsdk.ServiceClient, opts ListImageOpts) (*pagination.Pager, error) {
	url := imagesURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	page := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := ImagePage{pagination.OffsetPageBase{PageResult: r}}
		return p
	})

	return &page, nil
}

func ListSwitchableFlavors(c *golangsdk.ServiceClient, id string) (*flavorResp, error) {
	url := switchableFlavorsURL(c, id)

	var rst flavorResp
	_, err := c.Get(url, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func Mount(c *golangsdk.ServiceClient, id string, opts MountStorageOpts) (*MountStorage, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst MountStorage
	_, err = c.Post(mountURL(c, id), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func GetMount(c *golangsdk.ServiceClient, id, storageId string) (*MountStorage, error) {
	var rst MountStorage
	_, err := c.Get(mountDetailURL(c, id, storageId), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func DeleteMount(c *golangsdk.ServiceClient, id, storageId string) (*MountStorage, error) {
	url := mountDetailURL(c, id, storageId)
	var rst MountStorage
	_, err := c.DeleteWithResponse(url, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func ListMounts(c *golangsdk.ServiceClient, id string) (*MountStorageListResp, error) {
	var rst MountStorageListResp
	_, err := c.Get(mountURL(c, id), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}
