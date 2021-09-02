package applications

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type commonResult struct {
	golangsdk.Result
}

// CreateResult represents a result of the Create method.
type CreateResult struct {
	commonResult
}

// GetResult represents a result of the Get operation.
type GetResult struct {
	commonResult
}

// UpdateResult represents a result of the Update operation.
type UpdateResult struct {
	commonResult
}

// ResetSecretResult represents a result of the ResetAppSecret operation.
type ResetSecretResult struct {
	commonResult
}

type Application struct {
	// Creator of the application.
	//     USER: The app is created by the API user.
	//     MARKET: The app is allocated by the marketplace.
	Creator string `json:"creator"`
	// Registraion time.
	RegistraionTime string `json:"register_time"`
	// Update time.
	UpdateTime string `json:"update_time"`
	// App key.
	AppKey string `json:"app_key"`
	// App name.
	Name string `json:"name"`
	// Description.
	Description string `json:"remark"`
	// ID.
	Id string `json:"id"`
	// App secret.
	AppSecret string `json:"app_secret"`
	// App status.
	Status int `json:"status"`
	// App type, Currently only supports 'apig'. List method are not support.
	Type string `json:"app_type"`
	// Number of APIs. Only used for List method.
	BindNum int `json:"bind_num"`
}

func (r commonResult) Extract() (*Application, error) {
	var s Application
	err := r.ExtractInto(&s)
	return &s, err
}

// ApplicationPage represents the response pages of the List operation.
type ApplicationPage struct {
	pagination.SinglePageBase
}

func ExtractApplications(r pagination.Page) ([]Application, error) {
	var s []Application
	err := r.(ApplicationPage).Result.ExtractIntoSlicePtr(&s, "apps")
	return s, err
}

// DeleteResult represents a result of the Delete method.
type DeleteResult struct {
	golangsdk.ErrResult
}

type AppCode struct {
	// AppCode value, which contains 64 to 180 characters, starting with a letter, plus sign (+) or slash (/).
	// Only letters and the following special characters are allowed: +-_!@#$%/=
	Code string `json:"app_code"`
	// AppCode ID.
	Id string `json:"id"`
	// App ID.
	AppId string `json:"app_id"`
	// Creation time, in UTC format.
	CreateTime string `json:"create_time"`
}

// CreateCodeResult represents a result of the CreateAppCode method.
type CreateCodeResult struct {
	CodeResult
}

// AutoGenerateCodeResult represents a result of the AutoGenerateAppCode method.
type AutoGenerateCodeResult struct {
	CodeResult
}

// GetCodeResult represents a result of the GetAppCode method.
type GetCodeResult struct {
	CodeResult
}

type CodeResult struct {
	golangsdk.Result
}

func (r CodeResult) Extract() (*AppCode, error) {
	var s AppCode
	err := r.ExtractInto(&s)
	return &s, err
}

// AppCodePage represents the response pages of the ListAppCode operation.
type AppCodePage struct {
	pagination.SinglePageBase
}

func ExtractAppCodes(r pagination.Page) ([]AppCode, error) {
	var s []AppCode
	err := r.(AppCodePage).Result.ExtractIntoSlicePtr(&s, "app_codes")
	return s, err
}
