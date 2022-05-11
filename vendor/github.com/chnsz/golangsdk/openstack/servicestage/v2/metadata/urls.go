package metadata

import "github.com/chnsz/golangsdk"

const rootPath = "cas/metadata"

func runtimeURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, "runtimes")
}

func flavorURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, "flavors")
}
