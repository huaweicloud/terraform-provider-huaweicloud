package agency

import "github.com/chnsz/golangsdk"

const (
	rootPath     = "OS-AGENCY"
	resourcePath = "agencies"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func roleURL(c *golangsdk.ServiceClient, resource, resourceID, agencyID, roleID string) string {
	return c.ServiceURL(rootPath, resource, resourceID, resourcePath, agencyID, "roles", roleID)
}

func listRolesURL(c *golangsdk.ServiceClient, resource, resourceID, agencyID string) string {
	return c.ServiceURL(rootPath, resource, resourceID, resourcePath, agencyID, "roles")
}

func inheritedURL(client *golangsdk.ServiceClient, domainID, agencyID, roleID string) string {
	return client.ServiceURL("OS-INHERIT", "domains", domainID, resourcePath, agencyID,
		"roles", roleID, "inherited_to_projects")
}

func listInheritedURL(client *golangsdk.ServiceClient, domainID, agencyID string) string {
	return client.ServiceURL("OS-INHERIT", "domains", domainID, resourcePath, agencyID, "roles", "inherited_to_projects")
}
