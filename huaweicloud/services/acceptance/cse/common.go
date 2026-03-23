package cse

import "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

// getAcceptanceEpsId returns the acceptance enterprise project ID.
// If the acceptance enterprise project ID is not set, return "0" (default enterprise project ID).
// Notes: This function is only available for enterprise project granted users.
func getAcceptanceEpsId() string {
	if acceptance.HW_ENTERPRISE_PROJECT_ID_TEST == "" {
		return "0"
	}
	return acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
}
