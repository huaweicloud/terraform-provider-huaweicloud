package function

import (
	"io/ioutil"
	"net/http"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
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
	FuncName            string            `json:"func_name" required:"true"`
	MemorySize          int               `json:"memory_size" required:"true"`
	Package             string            `json:"package" required:"true"`
	Runtime             string            `json:"runtime" required:"true"`
	Timeout             int               `json:"timeout" required:"true"`
	AppXrole            string            `json:"app_xrole,omitempty"`
	CodeFilename        string            `json:"code_filename,omitempty"`
	CodeType            string            `json:"code_type,omitempty"`
	CodeUrl             string            `json:"code_url,omitempty"`
	CustomImage         *CustomImage      `json:"custom_image,omitempty"`
	Description         string            `json:"description,omitempty"`
	EncryptedUserData   string            `json:"encrypted_user_data,omitempty"`
	EnterpriseProjectID string            `json:"enterprise_project_id,omitempty"`
	FuncCode            *FunctionCodeOpts `json:"func_code,omitempty"`
	Handler             string            `json:"handler,omitempty"`
	Type                string            `json:"type,omitempty"`
	UserData            string            `json:"user_data,omitempty"`
	Xrole               string            `json:"xrole,omitempty"`
}

type CustomImage struct {
	Enabled bool   `json:"enabled" required:"true"`
	Image   string `json:"image" required:"true"`
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
	CodeType     string           `json:"code_type" required:"true"`
	CodeUrl      string           `json:"code_url,omitempty"`
	DependList   []string         `json:"depend_list,omitempty"`
	CodeFileName string           `json:"code_filename,omitempty"`
	FuncCode     FunctionCodeOpts `json:"func_code,omitempty"`
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
	Handler            string       `json:"handler" required:"true"`
	MemorySize         int          `json:"memory_size" required:"true"`
	Timeout            int          `json:"timeout" required:"true"`
	Runtime            string       `json:"runtime" required:"true"`
	Package            string       `json:"package,omitempty"`
	FuncVpc            *FuncVpc     `json:"func_vpc,omitempty"`
	MountConfig        *MountConfig `json:"mount_config,omitempty"`
	CodeUrl            string       `json:"code_url,omitempty"`
	Description        string       `json:"description,omitempty"`
	UserData           string       `json:"user_data,omitempty"`
	EncryptedUserData  string       `json:"encrypted_user_data,omitempty"`
	Xrole              string       `json:"xrole,omitempty"`
	AppXrole           string       `json:"app_xrole,omitempty"`
	InitializerHandler string       `json:"initializer_handler,omitempty"`
	InitializerTimeout int          `json:"initializer_timeout,omitempty"`
	CustomImage        *CustomImage `json:"custom_image,omitempty"`
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

// AsyncInvokeConfigOpts is the structure that used to modify the asynchronous invocation configuration.
type AsyncInvokeConfigOpts struct {
	// The maximum validity period of a message.
	MaxAsyncEventAgeInSeconds int `json:"max_async_event_age_in_seconds,omitempty"`
	// The maximum number of retry attempts to be made if asynchronous invocation fails.
	MaxAsyncRetryAttempts int `json:"max_async_retry_attempts,omitempty"`
	// Asynchronous invocation target.
	DestinationConfig DestinationConfig `json:"destination_config,omitempty"`
	// Whether to enable asynchronous invocation status persistence.
	EnableAsyncStatusLog *bool `json:"enable_async_status_log,omitempty"`
}

// DestinationConfig is the structure that represents the asynchronous invocation target.
type DestinationConfig struct {
	// The target to be invoked when a function is successfully executed.
	OnSuccess DestinationConfigDetails `json:"on_success,omitempty"`
	// The target to be invoked when a function fails to be executed due to a  system error or an internal error.
	OnFailure DestinationConfigDetails `json:"on_failure,omitempty"`
}

// DestinationConfigDetails is the structure that represents the configuration details of the asynchronous invocation.
type DestinationConfigDetails struct {
	// The object type.
	Destination string `json:"destination,omitempty"`
	// The parameters (in JSON format) corresponding to the target service.
	Param string `json:"param,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// UpdateAsyncInvokeConfig is the method that used to enable or modify the asynchronous invocation.
func UpdateAsyncInvokeConfig(c *golangsdk.ServiceClient, functionUrn string,
	opts AsyncInvokeConfigOpts) (*AsyncInvokeConfig, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r AsyncInvokeConfig
	_, err = c.Put(asyncInvokeConfigURL(c, functionUrn), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// GetAsyncInvokeConfig is the method that used to query the configuration details of the asynchronous invocation.
func GetAsyncInvokeConfig(c *golangsdk.ServiceClient, functionUrn string) (*AsyncInvokeConfig, error) {
	var r AsyncInvokeConfig
	_, err := c.Get(asyncInvokeConfigURL(c, functionUrn), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// DeleteAsyncInvokeConfig is the method that used to delete the asynchronous invocation.
func DeleteAsyncInvokeConfig(c *golangsdk.ServiceClient, functionUrn string) error {
	_, err := c.Delete(asyncInvokeConfigURL(c, functionUrn), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// MaxInstanceConfig is the structure used to update the max instance configuration for function.
type MaxInstanceConfig struct {
	// The maximum number of instances of the function.
	Number *int `json:"max_instance_num,omitempty"`
}

func UpdateMaxInstanceNumber(c *golangsdk.ServiceClient, functionUrn string, num int) (*Function, error) {
	config := MaxInstanceConfig{
		Number: &num,
	}
	var r Function
	_, err := c.Put(maxInstanceNumberURL(c, functionUrn), config, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// TagsActionOpts is an structure that used to manage function tags.
type TagsActionOpts struct {
	// The action name.
	Action string `json:"action,omitempty"`
	// Tag list.
	Tags []tags.ResourceTag `json:"tags,omitempty"`
	// System tag list.
	SysTags []tags.ResourceTag `json:"sys_tags,omitempty"`
}

// CreateResourceTags is the method that used to add tags to function using given parameters.
func CreateResourceTags(c *golangsdk.ServiceClient, functionUrn string, opts TagsActionOpts) error {
	_, err := c.Post(tagsActionURL(c, functionUrn, "create"), opts, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
		OkCodes:     []int{204},
	})
	return err
}

// DeleteResourceTags is the method that used to delete tags from function using given parameters.
func DeleteResourceTags(c *golangsdk.ServiceClient, functionUrn string, opts TagsActionOpts) error {
	_, err := c.DeleteWithBody(tagsActionURL(c, functionUrn, "delete"), opts, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
