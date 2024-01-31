package policyexecutelogs

import (
	"github.com/chnsz/golangsdk"
)

func listURL(c *golangsdk.ServiceClient, policyID string) string {
	return c.ServiceURL("scaling_policy_execute_log", policyID)
}
