package flavors

import (
	"github.com/huaweicloud/golangsdk"
)

func getURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id)
}

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("flavors", "detail")
}

func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("flavors")
}

func deleteURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id)
}

func accessURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id, "os-flavor-access")
}

func accessActionURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id, "action")
}

func extraSpecsListURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id, "os-extra_specs")
}

func extraSpecsGetURL(client *golangsdk.ServiceClient, id, key string) string {
	return client.ServiceURL("flavors", id, "os-extra_specs", key)
}

func extraSpecsCreateURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id, "os-extra_specs")
}

func extraSpecUpdateURL(client *golangsdk.ServiceClient, id, key string) string {
	return client.ServiceURL("flavors", id, "os-extra_specs", key)
}

func extraSpecDeleteURL(client *golangsdk.ServiceClient, id, key string) string {
	return client.ServiceURL("flavors", id, "os-extra_specs", key)
}
