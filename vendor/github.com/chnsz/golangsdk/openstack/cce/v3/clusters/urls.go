package clusters

import "github.com/chnsz/golangsdk"

const (
	rootPath      = "clusters"
	certPath      = "clustercert"
	masterIpPath  = "mastereip"
	operationPath = "operation"
)

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id)
}

func certificateURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id, certPath)
}

func masterIpURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id, masterIpPath)
}

func operationURL(c *golangsdk.ServiceClient, id, action string) string {
	return c.ServiceURL(rootPath, id, operationPath, action)
}
