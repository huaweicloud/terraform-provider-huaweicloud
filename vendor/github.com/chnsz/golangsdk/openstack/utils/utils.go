package utils

import "fmt"

func DeleteNotPassParams(params *map[string]interface{}, not_pass_params []string) {
	for _, i := range not_pass_params {
		delete(*params, i)
	}
}

// GenerateEpsIDQuery is used to generate a request URL with enterprise_project_id
func GenerateEpsIDQuery(epsID string) string {
	if len(epsID) == 0 {
		return ""
	}
	return fmt.Sprintf("?enterprise_project_id=%s", epsID)
}
