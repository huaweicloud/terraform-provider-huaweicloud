package terminals

import "github.com/chnsz/golangsdk/pagination"

// TerminalBindingResp is the structure that represents the bind information between MAC addresses and desktops.
type TerminalBindingResp struct {
	// Line info.
	Line int `json:"line"`
	// Terminal MAC address.
	MAC string `json:"mac"`
	// Desktop name.
	DesktopName string `json:"desktop_name"`
	// Description.
	Description string `json:"description"`
	// Bind ID.
	ID string `json:"id"`
	// SID.
	SID string `json:"sid"`
	// The validation result code.
	ValidationResultCode string `json:"validation_result_code"`
}

// TerminalBindingPage is a single page maximum result representing a query by offset page.
type TerminalBindingPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a TerminalBindingPage struct is empty.
func (b TerminalBindingPage) IsEmpty() (bool, error) {
	arr, err := ExtractTerminalBindings(b)
	return len(arr) == 0, err
}

// ExtractTerminalBindings is a method to extract the list of terminal binding desktops information.
func ExtractTerminalBindings(r pagination.Page) ([]TerminalBindingResp, error) {
	var s []TerminalBindingResp
	err := r.(TerminalBindingPage).Result.ExtractIntoSlicePtr(&s, "bind_list")
	return s, err
}

type deleteResp struct {
	// Delete result list.
	ResultList []DeleteResult `json:"result_list"`
}

// DeleteResult is the structure that represents the request response of the terminal binding.
type DeleteResult struct {
	// Bind ID.
	ID string `json:"id"`
	// Delete result code.
	DeleteResultCode string `json:"delete_result_code"`
	// Delete result message.
	DeleteResultMsg string `json:"delete_result_msg"`
}
