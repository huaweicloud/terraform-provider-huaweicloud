package users

type User struct {
	// Activation time, UTC timestamp.
	ActiveTime int `json:"activeTime"`
	// Administrator type.
	//   0: default (super) admin
	//   1: Ordinary administrator
	//   2: Non-administrator (that is, an ordinary enterprise member, valid when UserType is 2)
	AdminType int `json:"adminType"`
	// Business registration information.
	// This data is only returned when users query their own information.
	Corp CorpBasicInfo `json:"corp"`
	// The country to which the phone number belongs.
	Country string `json:"country"`
	// Department code.
	DeptCode string `json:"deptCode"`
	// Department name.
	DeptName string `json:"deptName"`
	// Department full name.
	DeptNamePath string `json:"deptNamePath"`
	// Description.
	Description string `json:"desc"`
	// Binding device type information.
	DevType DeviceInfo `json:"devType"`
	// Email.
	Email string `json:"email"`
	// English name.
	EnglishName string `json:"englishName"`
	// User function bits.
	Function UserFunction `json:"function"`
	// Whether to hide the phone number.
	HidePhone bool `json:"hidePhone"`
	// User ID.
	ID string `json:"id"`
	// License.
	//   0: commercial;
	//   1: Free trial.
	License int `json:"license"`
	// Name.
	Name string `json:"name"`
	// Phone number.
	Phone string `json:"phone"`
	// Signature.
	Signature string `json:"signature"`
	// SIP number.
	SipNum string `json:"sipNum"`
	// Address book sorting level, the lower the serial number, the higher the priority.
	SortLevel int `json:"sortLevel"`
	// user status.
	//   0: normal;
	//   1: Disable.
	Status int `json:"status"`
	// Third-party User ID.
	ThirdAccount string `json:"thirdAccount"`
	// Position (Title).
	Title string `json:"title"`
	// HUAWEI CLOUD conference user account.
	UserAccount string `json:"userAccount"`
	// User type.
	// 2: Enterprise member account
	UserType int `json:"userType"`
	// Smart screen unique account.
	VisionAccount string `json:"visionAccount"`
	// Cloud meeting room list.
	VmrList []UserVmr `json:"vmrList"`
}

type CorpBasicInfo struct {
	// Administrator account.
	Account string `json:"account"`
	// business location.
	Address string `json:"address"`
	// Administrator name.
	AdminName string `json:"adminName"`
	// Whether to support automatic account opening.
	AutoUserCreate bool `json:"autoUserCreate"`
	// The country to which the administrator's phone belongs.
	Country string `json:"country"`
	// Administrator email.
	Email string `json:"email"`
	// Whether it has pstn function.
	EnablePstn bool `json:"enablePstn"`
	// Whether to send meeting notices via SMS.
	EnableSMS bool `json:"enableSMS"`
	// Corporation ID.
	Id string `json:"id"`
	// Corporation name.
	Name string `json:"name"`
	// Admin phone number.
	Phone string `json:"phone"`
	// Whether to open cloud disk.
	EnableCloudDisk bool `json:"enableCloudDisk"`
	// Type of corporation.
	CorpType int `json:"corpType"`
}

type DeviceInfo struct {
	// Equipment end product dimensions.
	DeviceSize string `json:"deviceSize"`
	// Terminal model.
	Model string `json:"model"`
	// Terminal equipment purchase channels.
	PurchaseChannel string `json:"purchaseChannel"`
}

type UserVmr struct {
	// The ID of the cloud meeting room.
	// Corresponds to the vmrID in the create conference interface.
	ID string `json:"id"`
	// Fixed meeting ID of the cloud meeting room.
	// Corresponds to the vmrConferenceID of the data returned by the Create Conference API.
	VmrId string `json:"vmrId"`
	// Cloud meeting room name.
	VmrName string `json:"vmrName"`
	// VMR mode.
	//   0: Personal meeting ID
	//   1: Cloud meeting room
	//   2: Webinar
	VmrMode int `json:"vmrMode"`
	// The id of the cloud conference room package. Only cloud conference rooms are returned.
	VmrPkgId string `json:"vmrPkgId"`
	// The participation time of the cloud conference room package.
	// If it is 0, it means unlimited time, and only the cloud conference room is returned.
	VmrPkgLength int `json:"vmrPkgLength"`
	// The name of the cloud conference room package. Only cloud conference rooms are returned.
	VmrPkgName string `json:"vmrPkgName"`
	// The number of concurrent parties in the cloud conference room package.
	// Only the cloud conference room is returned.
	VmrPkgParties int `json:"vmrPkgParties"`
	// Cloud meeting room status.
	//   0: normal
	//   1: disable
	//   2: unassigned
	Status int `json:"status"`
}

type ErrResponse struct {
	// Error code.
	Code string `json:"error_code"`
	// Error message.
	Message string `json:"error_msg"`
	// Request ID.
	RequestId string `json:"request_id"`
}
