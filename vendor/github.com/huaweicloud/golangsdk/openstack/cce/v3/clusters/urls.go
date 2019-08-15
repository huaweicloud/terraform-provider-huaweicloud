package clusters

import "github.com/huaweicloud/golangsdk"

const (
	rootPath = "clusters"
	certPath = "clustercert"
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
