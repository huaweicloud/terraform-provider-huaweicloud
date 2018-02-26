package users

import "github.com/huaweicloud/golangsdk"

const (
	tenantPath = "tenants"
	userPath   = "users"
	rolePath   = "roles"
)

func ResourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(userPath, id)
}

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(userPath)
}

func listRolesURL(c *golangsdk.ServiceClient, tenantID, userID string) string {
	return c.ServiceURL(tenantPath, tenantID, userPath, userID, rolePath)
}
