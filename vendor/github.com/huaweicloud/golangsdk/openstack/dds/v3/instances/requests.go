package instances

import (
	"net/http"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type CreateOpts struct {
	Name             string         `json:"name"  required:"true"`
	DataStore        DataStore      `json:"datastore" required:"true"`
	Region           string         `json:"region" required:"true"`
	AvailabilityZone string         `json:"availability_zone" required:"true"`
	VpcId            string         `json:"vpc_id" required:"true"`
	SubnetId         string         `json:"subnet_id" required:"true"`
	SecurityGroupId  string         `json:"security_group_id" required:"true"`
	Password         string         `json:"password" required:"true"`
	DiskEncryptionId string         `json:"disk_encryption_id,omitempty"`
	Ssl              string         `json:"ssl_option,omitempty"`
	Mode             string         `json:"mode" required:"true"`
	Flavor           []Flavor       `json:"flavor" required:"true"`
	BackupStrategy   BackupStrategy `json:"backup_strategy,omitempty"`
}

type DataStore struct {
	Type          string `json:"type" required:"true"`
	Version       string `json:"version" required:"true"`
	StorageEngine string `json:"storage_engine" required:"true"`
}

type Flavor struct {
	Type     string `json:"type" required:"true"`
	Num      int    `json:"num" required:"true"`
	Storage  string `json:"storage,omitempty"`
	Size     int    `json:"size,omitempty"`
	SpecCode string `json:"spec_code" required:"true"`
}

type BackupStrategy struct {
	StartTime string `json:"start_time" required:"true"`
	KeepDays  int    `json:"keep_days,omitempty"`
}

type CreateInstanceBuilder interface {
	ToInstancesCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToInstancesCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Create(client *golangsdk.ServiceClient, opts CreateInstanceBuilder) (r CreateResult) {
	b, err := opts.ToInstancesCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})
	return
}

type DeleteInstance struct {
	InstanceId string `json:"instance_id" required:"true"`
}

type DeleteInstanceBuilder interface {
	ToInstancesDeleteMap() (map[string]interface{}, error)
}

func (opts DeleteInstance) ToInstancesDeleteMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(&opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Delete(client *golangsdk.ServiceClient, instanceId string) (r DeleteInstanceResult) {

	url := deleteURL(client, instanceId)

	_, r.Err = client.Delete(url, &golangsdk.RequestOpts{JSONResponse: &r.Body, MoreHeaders: map[string]string{"Content-Type": "application/json"}})
	return
}

type ListInstanceOpts struct {
	Id            string `q:"id"`
	Name          string `q:"name"`
	Mode          string `q:"mode"`
	DataStoreType string `q:"datastore_type"`
	VpcId         string `q:"vpc_id"`
	SubnetId      string `q:"subnet_id"`
	Offset        int    `q:"offset"`
	Limit         int    `q:"limit"`
}

type ListInstanceBuilder interface {
	ToInstanceListDetailQuery() (string, error)
}

func (opts ListInstanceOpts) ToInstanceListDetailQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, opts ListInstanceBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToInstanceListDetailQuery()

		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageList := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return InstancePage{pagination.SinglePageBase(r)}
	})

	return pageList
}

// UpdateOpt defines the basic information for update APIs
// URI: <Method> base_url/<Action>
// Request body: {<Param>: <Value>}
// the supported value for Method including: "post" and "put"
type UpdateOpt struct {
	Param  string
	Value  interface{}
	Action string
	Method string
}

func Update(client *golangsdk.ServiceClient, instanceId string, opts []UpdateOpt) (r UpdateInstanceResult) {
	for _, optRaw := range opts {
		url := modifyURL(client, instanceId, optRaw.Action)
		body := map[string]interface{}{
			optRaw.Param: optRaw.Value,
		}

		var httpMethod func(string, interface{}, interface{}, *golangsdk.RequestOpts) (*http.Response, error)
		if optRaw.Method == "post" {
			httpMethod = client.Post
		} else {
			httpMethod = client.Put
		}

		_, r.Err = httpMethod(url, body, &r.Body, &golangsdk.RequestOpts{
			OkCodes: []int{200, 202},
		})

		if r.Err != nil {
			break
		}
	}
	return
}
