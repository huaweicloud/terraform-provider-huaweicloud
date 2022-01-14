package secrets

import "github.com/chnsz/golangsdk"

// rootURL /v1/{project_id}/secrets
func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("secrets")
}

// resourceURL /v1/{project_id}/secrets/{secret_name}
func resourceURL(c *golangsdk.ServiceClient, secretName string) string {
	return c.ServiceURL("secrets", secretName)
}

// resourceVersionURL /v1/{project_id}/secrets/{secret_name}/versions/{version_id}
func resourceVersionURL(c *golangsdk.ServiceClient, secretName string, versionID string) string {
	return c.ServiceURL("secrets", secretName, "versions", versionID)
}

// versionRootURL /v1/{project_id}/secrets/{secret_name}/versions
func versionRootURL(c *golangsdk.ServiceClient, secretName string) string {
	return c.ServiceURL("secrets", secretName, "versions")
}
