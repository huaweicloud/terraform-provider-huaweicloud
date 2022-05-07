package auditlog

type UpdateResp struct {
	Result    string `json:"result"`
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

type AuditLogStatus struct {
	SwitchStatus string `json:"switch_status"`
	ErrorCode    string `json:"error_code"`
	ErrorMsg     string `json:"error_msg"`
}
