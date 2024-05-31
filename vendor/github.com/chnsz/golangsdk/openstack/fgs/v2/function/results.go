package function

import (
	"strconv"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type Function struct {
	Id                  string         `json:"-"`
	FuncId              string         `json:"-"`
	FuncUrn             string         `json:"func_urn"`
	FuncName            string         `json:"func_name"`
	DomainId            string         `json:"domain_id"`
	Namespace           string         `json:"namespace"`
	ProjectName         string         `json:"project_name"`
	Package             string         `json:"package"`
	Runtime             string         `json:"runtime"`
	Timeout             int            `json:"timeout"`
	Handler             string         `json:"handler"`
	MemorySize          int            `json:"memory_size"`
	Cpu                 int            `json:"cpu"`
	CodeType            string         `json:"code_type"`
	CodeUrl             string         `json:"code_url"`
	CodeFileName        string         `json:"code_filename"`
	CodeSize            int64          `json:"code_size"`
	CustomImage         CustomImage    `json:"custom_image"`
	UserData            string         `json:"user_data"`
	EncryptedUserData   string         `json:"encrypted_user_data"`
	Digest              string         `json:"digest"`
	Version             string         `json:"version"`
	ImageName           string         `json:"image_name"`
	Xrole               string         `json:"xrole"`
	AppXrole            string         `json:"app_xrole"`
	Description         string         `json:"description"`
	VersionDescription  string         `json:"version_description"`
	LastmodifiedUtc     int64          `json:"-"`
	LastModified        string         `json:"last_modified"`
	FuncCode            FunctionCode   `json:"func_code"`
	FuncVpc             FuncVpc        `json:"func_vpc"`
	MountConfig         MountConfig    `json:"mount_config,omitempty"`
	Concurrency         int            `json:"-"`
	DependList          []string       `json:"depend_list"` // Deprecated
	DependVersionList   []string       `json:"depend_version_list"`
	StrategyConfig      StrategyConfig `json:"strategy_config"`
	ExtendConfig        string         `json:"extend_config"`
	Dependencies        []*Dependency  `json:"dependencies"`
	InitializerTimeout  int            `json:"initializer_timeout,omitempty"`
	InitializerHandler  string         `json:"initializer_handler,omitempty"`
	EnterpriseProjectID string         `json:"enterprise_project_id"`
	Type                string         `json:"type"`
	// GPU memory.
	// Range: 1024 to 16,384, and the value is a multiple of 1024.
	GPUMemory int `json:"gpu_memory"`
	// GPU type.
	GPUType string `json:"gpu_type"`
	// Ephemeral storage size, the maximum value is 10 GB. Defaults to 512 MB.
	EphemeralStorage int `json:"ephemeral_storage"`
	// Whether to allow a long timeout.
	LongTime bool `json:"long_time"`
	// Log group ID.
	LogGroupId string `json:"log_group_id"`
	// Log stream ID.
	LogStreamId string `json:"log_stream_id"`
	// Network configuration.
	NetworkController NetworkControlConfig `json:"network_controller"`
	// Whether stateful functions are supported.
	IsStatefulFunction bool `json:"is_stateful_function"`
	// Whether to enable dynamic memory allocation.
	EnableDynamicMemory bool `json:"enable_dynamic_memory"`
	// Whether to allow authentication information in the request header.
	EnableAuthInHeader bool `json:"enable_auth_in_header"`
	// Private domain name.
	DomainNames string `json:"domain_names"`
	// The pre-stop handler of the function.
	PreStopHandler string `json:"pre_stop_handler"`
	// Maximum duration the function can be initialized.
	PreStopTimeout int `json:"pre_stop_timeout"`
}

type FuncMount struct {
	Id             string             `json:"id,omitempty"`
	MountType      string             `json:"mount_type" required:"true"`
	MountResource  string             `json:"mount_resource" required:"true"`
	MountSharePath string             `json:"mount_share_path" required:"true"`
	LocalMountPath string             `json:"local_mount_path" required:"true" description:"local file path in function runtime environment"`
	UserGroupId    *int               `json:"-"`
	UserId         *int               `json:"-"`
	Status         string             `json:"status,omitempty"` //ACTIVE或DISABLED，和触发器类似。如果已经存在的配置不可用了processrouter不会挂载。
	ProjectId      string             `json:"-"`
	FuncVersions   []*FunctionVersion `json:"-"`
	SaveType       int                `json:"-"` //仅仅在数据处理时用到，如果需要保存新的映射关系，则将其置为1，如要删除老的，将其置为2
}

// noinspection GoNameStartsWithPackageName
type FunctionVersion struct {
	Id                 string        `json:"-"`
	FuncId             string        `json:"-"`
	Runtime            string        `json:"runtime"`
	Timeout            int           `json:"timeout"`
	Handler            string        `json:"handler"`
	MemorySize         int           `json:"memory_size"`
	Cpu                int           `json:"cpu"`
	CodeType           string        `json:"code_type"`
	CodeUrl            string        `json:"code_url"`
	CodeFileName       string        `json:"code_file_name"`
	CodeSize           int64         `json:"code_size"`
	UserData           string        `json:"user_data"`
	EncryptedUserData  string        `json:"encrypted_user_data"`
	ImageName          string        `json:"image_name"`
	Digest             string        `json:"digest"`
	Version            string        `json:"version"`
	Xrole              string        `json:"xrole"`
	AppXrole           string        `json:"app_xrole"`
	Description        string        `json:"description"`
	VersionDescription string        `json:"version_description"`
	LastModified       int64         `json:"last_modified"`
	Concurrency        int           `json:"concurrency"`
	ExtendConfig       string        `json:"extend_config"`
	Dependencies       []*Dependency `json:"dependencies"`
	FuncBase           *FunctionBase `json:"func_base"`
	FuncVpcId          string        `json:"-"`
	FuncMounts         []*FuncMount
	MountConfig        *MountConfig `json:"mount_config"`
	Vpc                *FuncVpc     `json:"vpc"`
	InitializerTimeout int
	InitializerHandler string `description:"the function initializer handler"`
}

// dependency
type Dependency struct {
	Id           string             `json:"id"`
	Owner        string             `json:"owner"`
	Namespace    string             `json:"-"`
	Link         string             `json:"link"`
	Runtime      string             `json:"runtime"`
	ETag         string             `json:"etag"`
	Size         int64              `json:"size"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	FileName     string             `json:"file_name,omitempty"`
	FuncVersions []*FunctionVersion `json:"-"`
	SaveType     int                `json:"-"`
}

type MountUser struct {
	UserId      int `json:"user_id" required:"true"`
	UserGroupId int `json:"user_group_id" required:"true"`
}

type MountConfig struct {
	MountUser  MountUser   `json:"mount_user" required:"true"`
	FuncMounts []FuncMount `json:"func_mounts" required:"true"`
}

// noinspection GoNameStartsWithPackageName
type FunctionCode struct {
	File string `json:"file"`
	Link string `json:"link"`
}

// noinspection GoNameStartsWithPackageName
type FunctionBase struct {
	Id          string `json:"-"`
	FuncName    string `json:"func_name"`
	DomainId    string `description:"domain id"`
	Namespace   string `json:"namespace"`
	ProjectName string `json:"project_name"`
	Package     string `json:"package"`
}

type FuncVpc struct {
	Id             string   `json:"-"`
	DomainId       string   `json:"-" validate:"regexp=^[a-zA-Z0-9-]+$" description:"domain id"`
	Namespace      string   `json:"-"`
	VpcName        string   `json:"vpc_name,omitempty"`
	VpcId          string   `json:"vpc_id,omitempty"`
	SubnetName     string   `json:"subnet_name,omitempty"`
	SubnetId       string   `json:"subnet_id,omitempty"`
	Cidr           string   `json:"cidr,omitempty"`
	Gateway        string   `json:"gateway,omitempty"`
	SecurityGroups []string `json:"security_groups,omitempty"`
}

type TriggerAccessVpcs struct {
	VpcId   string `json:"vpc_id"`
	VpcName string `json:"trigger_access_vpcs"`
}

type NetworkController struct {
	DisablePublicNetwork bool                 `json:"disable_public_network"`
	TriggerAccessVpcs    []*TriggerAccessVpcs `json:"trigger_access_vpcs"`
}

type commonResult struct {
	golangsdk.Result
}

type GetResult struct {
	commonResult
}

type CreateResult struct {
	commonResult
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type UpdateResult struct {
	commonResult
}

type FunctionPage struct {
	pagination.SinglePageBase
}

func (r commonResult) Extract() (*Function, error) {
	var f Function
	err := r.ExtractInto(&f)
	return &f, err
}

// noinspection GoNameStartsWithPackageName
type FunctionList struct {
	Functions  []Function `json:"functions"`
	NextMarker int        `json:"next_marker"`
}

func ExtractList(r pagination.Page) (FunctionList, error) {
	var s FunctionList
	err := (r.(FunctionPage)).ExtractInto(&s)
	return s, err
}

func (r commonResult) ExtractInvoke() (interface{}, error) {
	return r.Body, r.Err
}

type Versions struct {
	Versions   []Function `json:"versions"`
	NextMarker int        `json:"next_marker"`
}

func ExtractVersionlist(r pagination.Page) (Versions, error) {
	var s Versions
	err := (r.(FunctionPage)).ExtractInto(&s)
	return s, err
}

type AliasResult struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Description  string `json:"description"`
	LastModified string `json:"last_modified"`
	AliasUrn     string `json:"alias_urn"`
}

func (r commonResult) ExtractAlias() (*AliasResult, error) {
	var s AliasResult
	err := r.ExtractInto(&s)
	return &s, err
}

func ExtractAliasList(r pagination.Page) ([]AliasResult, error) {
	var s []AliasResult
	err := (r.(FunctionPage)).ExtractInto(&s)
	return s, err
}

// AsyncInvokeConfig is the structure that represents the asynchronous invocation.
type AsyncInvokeConfig struct {
	// Function URN.
	FunctionUrn string `json:"func_urn"`
	// Maximum validity period of a message. Value range: 60–86,400. Unit: second.
	MaxAsyncEventAgeInSeconds int `json:"max_async_event_age_in_seconds"`
	// Maximum number of retry attempts to be made if asynchronous invocation fails. Default value: 3. Value range: 0–8.
	MaxAsyncRetryAttempts int `json:"max_async_retry_attempts"`
	// Asynchronous invocation target.
	DestinationConfig DestinationConfig `json:"destination_config"`
	// Time when asynchronous execution notification was configured.
	CreatedAt string `json:"created_time"`
	// Time when the asynchronous execution notification settings were last modified.
	UpdatedAt string `json:"last_modified"`
	// Whether to enable asynchronous invocation status persistence.
	EnableAsyncStatusLog bool `json:"enable_async_status_log"`
}

// ReservedInstanceConfigResp is the structure that represents the response of the GetReservedInstanceConfig method.
type ReservedInstanceConfigResp struct {
	// The list of reserved instance policy.
	ReservedInstances []ReservedInstancePolicy `json:"reserved_instances"`
	// The page information.
	PageInfo PageInfoObj `json:"page_info"`
	// Number of function.
	Count int `json:"count"`
}

// ReservedInstancePolicy is the structure that represents the reserved instance policy configuration.
type ReservedInstancePolicy struct {
	// Function URN.
	FunctionUrn string `json:"function_urn"`
	// Limited type, the supported values are version and alias.
	QualifierType string `json:"qualifier_type"`
	// The value of the limited type.
	QualifierName string `json:"qualifier_name"`
	// The number of instance reserved.
	MinCount int `json:"min_count"`
	// Whether to enable the idle mode configuration.
	IdleMode bool `json:"idle_mode"`
	// The auto scaling policy configuration.
	TacticsConfig TacticsConfigObj `json:"tactics_config"`
}

// PageInfoObj is the structure that represents pagination information of the function reserved instance.
type PageInfoObj struct {
	// Next record location.
	NextMarker int `json:"next_marker"`
	// Last record location.
	PreviousMarker int `json:"previous_marker"`
	// Total number of current page.
	CurrentCount int `json:"current_count"`
}

// ReservedInstanceConfigPage represents the response pages of the GetReservedInstanceConfig method.
type ReservedInstanceConfigPage struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if no reserved instance.
func (r ReservedInstanceConfigPage) IsEmpty() (bool, error) {
	resp, err := extractReservedInstanceConfigs(r)
	return len(resp) == 0, err
}

// LastMarker returns the last marker index in reserved instance list.
func (r ReservedInstanceConfigPage) LastMarker() (string, error) {
	resp, err := extractPageInfo(r)
	if err != nil {
		return "", err
	}

	if resp.CurrentCount == 0 {
		return "", nil
	}
	return strconv.Itoa(resp.NextMarker), nil
}

// extractPageInfo is a method which to extract the response of the page information.
func extractPageInfo(r pagination.Page) (*PageInfoObj, error) {
	var s ReservedInstanceConfigResp
	err := r.(ReservedInstanceConfigPage).Result.ExtractInto(&s)

	return &s.PageInfo, err
}

// extractReservedInstanceConfigs is a method which to extract the response to reserved instance configuration list.
func extractReservedInstanceConfigs(p pagination.Page) ([]ReservedInstancePolicy, error) {
	var resp ReservedInstanceConfigResp
	err := p.(ReservedInstanceConfigPage).Result.ExtractInto(&resp)

	return resp.ReservedInstances, err
}
