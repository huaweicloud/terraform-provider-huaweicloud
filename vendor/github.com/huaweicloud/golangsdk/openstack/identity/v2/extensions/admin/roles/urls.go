package roles

import "github.com/huaweicloud/golangsdk"

const (
	ExtPath  = "OS-KSADM"
	RolePath = "roles"
	UserPath = "users"
)

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(ExtPath, RolePath, id)
}

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(ExtPath, RolePath)
}

func userRoleURL(c *golangsdk.ServiceClient, tenantID, userID, roleID string) string {
	return c.ServiceURL("tenants", tenantID, UserPath, userID, RolePath, ExtPath, roleID)
}
