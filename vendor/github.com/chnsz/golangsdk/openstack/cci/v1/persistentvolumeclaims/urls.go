package persistentvolumeclaims

import "github.com/chnsz/golangsdk"

const rootPath = "namespaces"

func rootURL(client *golangsdk.ServiceClient, ns string) string {
	return client.ServiceURL(rootPath, ns, "extended-persistentvolumeclaims")
}

func resourceURL(client *golangsdk.ServiceClient, ns, name string) string {
	return client.ServiceURL(rootPath, ns, "persistentvolumeclaims", name)
}
