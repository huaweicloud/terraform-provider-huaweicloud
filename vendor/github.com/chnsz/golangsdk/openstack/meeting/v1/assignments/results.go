package assignments

// Administrator is the structure that represents the detail of the administrator account.
type Administrator struct {
	// User account.
	Account string `json:"account"`
	// User ID.
	ID string `json:"id"`
	// User name.
	Name string `json:"name"`
	// Administrator type.
	//   0: normal administrator.
	//   1: default administrator.
	AdminType int `json:"adminType"`
	// Country or region.
	Country string `json:"country"`
	// Department details.
	Department Department `json:"dept"`
	// Email address.
	Email string `json:"emial"`
	// Phone number.
	Phone string `json:"phone"`
}

// Department is the structure that represents the detail of the department.
type Department struct {
	// Corporation ID.
	CorpId string `json:"corpId"`
	// Department code.
	Code string `json:"deptCode"`
	// Department name.
	Name string `json:"deptName"`
	// Department path.
	NamePath string `json:"deptNamePath"`
	// Parent department code.
	ParentCode string `json:"parentDeptCode"`
}

// ErrResponse is the structure that represents the request error.
type ErrResponse struct {
	// Error code.
	Code string `json:"error_code"`
	// Error message.
	Message string `json:"error_msg"`
}
