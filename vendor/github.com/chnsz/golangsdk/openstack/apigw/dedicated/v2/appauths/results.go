package appauths

import "github.com/chnsz/golangsdk/pagination"

type createResp struct {
	// App authorization records.
	Auths []Authorization `json:"auths"`
}

// Authorization is the struct that represents the application authorization records.
type Authorization struct {
	// The application ID.
	ApiId string `json:"api_id"`
	// The authorization result.
	AuthResult AuthResult `json:"auth_result"`
	// The authorization time.
	AuthTime string `json:"auth_time"`
	// The authorization record ID.
	ID string `json:"id"`
	// The ID of the application allowed to access the API.
	AppId string `json:"app_id"`
	// The authorizer.
	// + PROVIDER: API provider
	// + CONSUMER: API user
	AuthRole string `json:"auth_role"`
	// Authorization channel type.
	// + NORMAL: normal channel
	// + GREEN: green channel
	// The default value is NORMAL.
	AuthTunnel string `json:"auth_tunnel"`
	// Whitelist for the green channel.
	AuthWhitelist []string `json:"auth_whitelist"`
	// Blacklist for the green channel.
	AuthBlacklist []string `json:"auth_blacklist"`
	// Access parameters.
	VisitParams string `json:"visit_params"`
}

// AuthResult is the structure that represents the authorization details.
type AuthResult struct {
	// Authorization result.
	// + SUCCESS
	// + SKIPPED
	// + FAILED
	Status string `json:"status"`
	// Error message.
	ErrorMsg string `json:"error_msg"`
	// Error code.
	ErrorCode string `json:"error_code"`
	// Name of the API for which authorization fails.
	ApiName string `json:"api_name"`
	// Name of the app that fails to be authorized.
	AppName string `json:"app_name"`
}

// AuthorizedPage is a page structure that represents each page information.
type AuthorizedPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a AuthorizedPage struct is empty.
func (b AuthorizedPage) IsEmpty() (bool, error) {
	arr, err := ExtractAuthorizedApis(b)
	return len(arr) == 0, err
}

// ExtractAuthorizedApis is a method to extract the list of authorized APIs for specified application.
func ExtractAuthorizedApis(r pagination.Page) ([]ApiAuthInfo, error) {
	var s []ApiAuthInfo
	err := r.(AuthorizedPage).Result.ExtractIntoSlicePtr(&s, "auths")
	return s, err
}

// ApiAuthInfo is the structure that represents the authorized API information.
type ApiAuthInfo struct {
	// Authorization record ID.
	ID string `json:"id"`
	// The API ID.
	ApiId string `json:"api_id"`
	// The API name.
	ApiName string `json:"api_name"`
	// Name of the API group to which the API belongs.
	GroupName string `json:"group_name"`
	// API type.
	ApiType int `json:"api_type"`
	// API description.
	ApiRemark string `json:"api_remark"`
	// ID of the environment in which an app has been authorized to call the API.
	EnvId string `json:"env_id"`
	// Authorizer.
	AuthRole string `json:"auth_role"`
	// Authorization time.
	AuthTime string `json:"auth_time"`
	// Application name.
	AppName string `json:"app_name"`
	// Application description.
	AppRemark string `json:"app_remark"`
	// Application type.
	AppType string `json:"app_type"`
	// Creator of the app.
	// + USER: The app is created by a tenant.
	// + MARKET: The app is allocated by KooGallery.
	AppCreator string `json:"app_creator"`
	// API publication record ID.
	PublishId string `json:"publish_id"`
	// ID of the API group to which the API belongs.
	GroupId string `json:"group_id"`
	// Authorization channel type.
	// + NORMAL: normal channel
	// + GREEN: green channel
	AuthTunnel string `json:"auth_tunnel"`
	// Whitelist for the green channel.
	AuthWhitelist []string `json:"auth_whitelist"`
	// Blacklist for the green channel.
	AuthBlacklist []string `json:"auth_blacklist"`
	// Access parameters.
	VisitParams string `json:"visit_params"`
	// ROMA application type.
	// + subscription: subscription application
	// + integration: integration application
	RomaAppType string `json:"roma_app_type"`
	// Name of the environment in which the app has been authorized to call the API.
	EnvName string `json:"env_name"`
	// Application ID.
	AppId string `json:"app_id"`
}

// UnauthorizedPage is a page structure that represents each page information.
type UnauthorizedPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a UnauthorizedPage struct is empty.
func (b UnauthorizedPage) IsEmpty() (bool, error) {
	arr, err := ExtractUnauthorizedApis(b)
	return len(arr) == 0, err
}

// ExtractUnauthorizedApis is a method to extract the list of unauthorized APIs for specified application.
func ExtractUnauthorizedApis(r pagination.Page) ([]ApiOutlineInfo, error) {
	var s []ApiOutlineInfo
	err := r.(UnauthorizedPage).Result.ExtractIntoSlicePtr(&s, "apis")
	return s, err
}

// ApiOutlineInfo is the structure that represents the unauthorized API information.
type ApiOutlineInfo struct {
	// API authentication mode.
	AuthType string `json:"auth_type"`
	// Name of the environment in which the API has been published.
	RunEnvName string `json:"run_env_name"`
	// Name of the API group to which the API belongs.
	GroupName string `json:"group_name"`
	// Publication record ID.
	PublishId string `json:"publish_id"`
	// ID of the API group to which the API belongs.
	GroupId string `json:"group_id"`
	// API name.
	Name string `json:"name"`
	// API description.
	Remark string `json:"remark"`
	// ID of the environment in which the API has been published.
	RunEnvId string `json:"run_env_id"`
	// Application ID.
	ID string `json:"id"`
	// API request address.
	ReqUri string `json:"req_uri"`
}
