package function

import (
	"io/ioutil"
	"net/http"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

//Create function
type CreateOptsBuilder interface {
	ToCreateFunctionMap() (map[string]interface{}, error)
}

//funcCode struct
type FunctionCodeOpts struct {
	File string `json:"file" required:"true"`
	Link string `json:"-"`
}

//function struct
type CreateOpts struct {
	FuncName      string           `json:"func_name" required:"true"`
	Package       string           `json:"package" required:"true"`
	CodeType      string           `json:"code_type" required:"true"`
	CodeUrl       string           `json:"code_url,omitempty"`
	Description   string           `json:"description,omitempty"`
	CodeFilename  string           `json:"code_filename,omitempty"`
	Handler       string           `json:"handler" required:"true"`
	MemorySize    int              `json:"memory_size" required:"true"`
	Runtime       string           `json:"runtime" required:"true"`
	Timeout       int              `json:"timeout" required:"true"`
	UserData      string           `json:"user_data,omitempty"`
	Xrole         string           `json:"xrole,omitempty"`
	AppXrole      string           `json:"app_xrole,omitempty"`
	DependencyPkg string           `json:"dependency_pkg,omitempty"`
	FuncCode      FunctionCodeOpts `json:"func_code" required:"true"`
}

func (opts CreateOpts) ToCreateFunctionMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

//create funtion
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	f, err := opts.ToCreateFunctionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), f, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

//functions list struct
type ListOpts struct {
	Marker   string `q:"marker"`
	MaxItems string `q:"maxitems"`
}

func (opts ListOpts) ToMetricsListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

type ListOptsBuilder interface {
	ToMetricsListQuery() (string, error)
}

//functions list
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToMetricsListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return FunctionPage{pagination.SinglePageBase(r)}
	})
}

//Querying the Metadata Information of a Function
func GetMetadata(c *golangsdk.ServiceClient, functionUrn string) (r GetResult) {
	_, r.Err = c.Get(getMetadataURL(c, functionUrn), &r.Body, nil)
	return
}

//Querying the Code of a Function
func GetCode(c *golangsdk.ServiceClient, functionUrn string) (r GetResult) {
	_, r.Err = c.Get(getCodeURL(c, functionUrn), &r.Body, nil)
	return
}

//Deleting a Function or Function Version
func Delete(c *golangsdk.ServiceClient, functionUrn string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, functionUrn), nil)
	return
}

type UpdateOptsBuilder interface {
	ToUpdateMap() (map[string]interface{}, error)
}

//Function struct for update
type UpdateCodeOpts struct {
	CodeType      string           `json:"code_type" required:"true"`
	CodeUrl       string           `json:"code_url,omitempty"`
	DependencyPkg string           `json:"dependency_pkg,omitempty"`
	FuncCode      FunctionCodeOpts `json:"func_code,omitempty"`
}

func (opts UpdateCodeOpts) ToUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

//Modifying the Code of a Function
func UpdateCode(c *golangsdk.ServiceClient, functionUrn string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(updateCodeURL(c, functionUrn), b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

//Metadata struct for update
type UpdateMetadataOpts struct {
	Runtime       string `json:"runtime" required:"true"`
	CodeType      string `json:"code_type" required:"true"`
	CodeUrl       string `json:"code_url,omitempty"`
	Description   string `json:"description,omitempty"`
	MemorySize    int    `json:"memory_size" required:"true"`
	Handler       string `json:"handler" required:"true"`
	Timeout       int    `json:"timeout" required:"true"`
	UserData      string `json:"user_data,omitempty"`
	DependencyPkg string `json:"dependency_pkg,omitempty"`
	Xrole         string `json:"xrole,omitempty"`
	AppXrole      string `json:"app_xrole,omitempty"`
}

func (opts UpdateMetadataOpts) ToUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

//Modifying the Metadata Information of a Function
func UpdateMetadata(c *golangsdk.ServiceClient, functionUrn string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(updateMetadataURL(c, functionUrn), b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

//verstion struct
type CreateVersionOpts struct {
	Digest      string `json:"digest,omitempty"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
}

func (opts CreateVersionOpts) ToCreateFunctionMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

//Publishing a Function Version
func CreateVersion(c *golangsdk.ServiceClient, opts CreateOptsBuilder, functionUrn string) (r CreateResult) {
	b, err := opts.ToCreateFunctionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createVersionURL(c, functionUrn), b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200, 201}})
	return
}

//Querying the Alias Information of a Function Version
func ListVersions(c *golangsdk.ServiceClient, opts ListOptsBuilder, functionUrn string) pagination.Pager {
	url := listVersionURL(c, functionUrn)
	if opts != nil {
		query, err := opts.ToMetricsListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return FunctionPage{pagination.SinglePageBase(r)}
	})
}

//Alias struct
type CreateAliasOpts struct {
	Name    string `json:"name" required:"true"`
	Version string `json:"version" required:"true"`
}

func (opts CreateAliasOpts) ToCreateFunctionMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

//Creating an Alias for a Function Version
func CreateAlias(c *golangsdk.ServiceClient, opts CreateOptsBuilder, functionUrn string) (r CreateResult) {
	b, err := opts.ToCreateFunctionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createAliasURL(c, functionUrn), b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

//Alias struct for update
type UpdateAliasOpts struct {
	Version     string `json:"version" required:"true"`
	Description string `json:"description,omitempty"`
}

func (opts UpdateAliasOpts) ToUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

//Modifying the Alias Information of a Function Version
func UpdateAlias(c *golangsdk.ServiceClient, functionUrn, aliasName string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(updateAliasURL(c, functionUrn, aliasName), b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

//Deleting an Alias of a Function Version
func DeleteAlias(c *golangsdk.ServiceClient, functionUrn, aliasName string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteAliasURL(c, functionUrn, aliasName), &golangsdk.RequestOpts{OkCodes: []int{204}})
	return
}

//Querying the Alias Information of a Function Version
func GetAlias(c *golangsdk.ServiceClient, functionUrn, aliasName string) (r GetResult) {
	_, r.Err = c.Get(getAliasURL(c, functionUrn, aliasName), &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

//Querying the Aliases of a Function's All Versions
func ListAlias(c *golangsdk.ServiceClient, functionUrn string) pagination.Pager {
	return pagination.NewPager(c, listAliasURL(c, functionUrn), func(r pagination.PageResult) pagination.Page {
		return FunctionPage{pagination.SinglePageBase(r)}
	})
}

//Executing a Function Synchronously
func Invoke(c *golangsdk.ServiceClient, m map[string]interface{}, functionUrn string) (r CreateResult) {
	var resp *http.Response
	resp, r.Err = c.Post(invokeURL(c, functionUrn), m, nil, &golangsdk.RequestOpts{
		OkCodes:      []int{200},
		JSONResponse: nil,
	})
	if resp != nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		r.Body = string(body)
	}
	return
}

//Executing a Function Asynchronously
func AsyncInvoke(c *golangsdk.ServiceClient, m map[string]interface{}, functionUrn string) (r CreateResult) {
	_, r.Err = c.Post(asyncInvokeURL(c, functionUrn), m, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{202}})
	return
}
