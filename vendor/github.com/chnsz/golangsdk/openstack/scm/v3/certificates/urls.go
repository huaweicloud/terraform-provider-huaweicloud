package certificates

import "github.com/chnsz/golangsdk"

const (
	rootPath     = "scm"
	resourcePath = "certificates"
	importPath   = "import"
	pushPath     = "push"
	exportPath   = "export"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func importURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath, importPath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func pushURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id, pushPath)
}

func exportURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id, exportPath)
}
