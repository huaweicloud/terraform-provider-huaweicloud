package auditlog

import "github.com/chnsz/golangsdk"

// POST /v3/{project_id}/instance/{instance_id}/audit-log/switch
func updateURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instance", instanceId, "audit-log", "switch")
}

// GET /v3/{project_id}/instance/{instance_id}/audit-log/switch-status
func getURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instance", instanceId, "audit-log", "switch-status")
}
