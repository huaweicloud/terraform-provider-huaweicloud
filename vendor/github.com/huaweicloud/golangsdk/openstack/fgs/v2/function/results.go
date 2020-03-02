package function

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type Function struct {
	Id                 string         `json:"-"`
	FuncId             string         `json:"-"`
	FuncUrn            string         `json:"func_urn"`
	FuncName           string         `json:"func_name"`
	DomainId           string         `json:"domain_id"`
	Namespace          string         `json:"namespace"`
	ProjectName        string         `json:"project_name"`
	Package            string         `json:"package"`
	Runtime            string         `json:"runtime"`
	Timeout            int            `json:"timeout"`
	Handler            string         `json:"handler"`
	MemorySize         int            `json:"memory_size"`
	Cpu                int            `json:"cpu"`
	CodeType           string         `json:"code_type"`
	CodeUrl            string         `json:"code_url"`
	CodeFileName       string         `json:"code_filename"`
	CodeSize           int64          `json:"code_size"`
	UserData           string         `json:"user_data"`
	Digest             string         `json:"digest"`
	Version            string         `json:"version"`
	ImageName          string         `json:"image_name"`
	Xrole              string         `json:"xrole"`
	AppXrole           *string        `json:"app_xrole"`
	Description        string         `json:"description"`
	VersionDescription string         `json:"version_description"`
	LastmodifiedUtc    int64          `json:"-"`
	LastModified       string         `json:"last_modified"`
	FuncCode           FunctionCode   `json:"func_code"`
	FuncVpc            *FuncVpc       `json:"func_vpc"`
	MountConfig        *MountConfig   `json:"mount_config,omitempty"`
	Concurrency        int            `json:"-"`
	DependList         []string       `json:"depend_list"`
	StrategyConfig     StrategyConfig `json:"strategy_config"`
	ExtendConfig       string         `json:"extend_config"`
	Dependencies       []*Dependency  `json:"dependencies"`
	InitializerTimeout int            `json:"initializer_timeout,omitempty"`
	InitializerHandler string         `json:"initializer_handler,omitempty"`
}

type FuncMount struct {
	Id             string             `json:"id"`
	MountType      string             `json:"mount_type"`
	MountResource  string             `json:"mount_resource"`
	MountSharePath string             `json:"mount_share_path"`
	LocalMountPath string             `json:"local_mount_path" description:"local file path in function runtime environment"`
	UserGroupId    *int               `json:"-"`
	UserId         *int               `json:"-"`
	Status         string             `json:"status"` //ACTIVE或DISABLED，和触发器类似。如果已经存在的配置不可用了processrouter不会挂载。
	ProjectId      string             `json:"-"`
	FuncVersions   []*FunctionVersion `json:"-"`
	SaveType       int                `json:"-"` //仅仅在数据处理时用到，如果需要保存新的映射关系，则将其置为1，如要删除老的，将其置为2
}

//noinspection GoNameStartsWithPackageName
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
	UserId      *int `json:"user_id"`
	UserGroupId *int `json:"user_group_id"`
}

type MountConfig struct {
	MountUser  *MountUser   `json:"mount_user"`
	FuncMounts []*FuncMount `json:"func_mounts"`
}

//noinspection GoNameStartsWithPackageName
type FunctionCode struct {
	File string `json:"file"`
	Link string `json:"link"`
}

//noinspection GoNameStartsWithPackageName
type FunctionBase struct {
	Id          string `json:"-"`
	FuncName    string `json:"func_name"`
	DomainId    string `description:"domain id"`
	Namespace   string `json:"namespace"`
	ProjectName string `json:"project_name"`
	Package     string `json:"package"`
}

type FuncVpc struct {
	Id         string `json:"-"`
	DomainId   string `json:"-" validate:"regexp=^[a-zA-Z0-9-]+$" description:"domain id"`
	Namespace  string `json:"-"`
	VpcName    string `json:"vpc_name"`
	VpcId      string `json:"vpc_id"`
	SubnetName string `json:"subnet_name"`
	SubnetId   string `json:"subnet_id"`
	Cidr       string `json:"cidr"`
	Gateway    string `json:"gateway"`
}

type StrategyConfig struct {
	Concurrency *int `json:"concurrency"`
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

//noinspection GoNameStartsWithPackageName
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
