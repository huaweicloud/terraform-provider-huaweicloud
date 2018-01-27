package policies

import (
	"github.com/huawei-clouds/golangsdk"
)

const resourcePath = "scaling_policy"

//createURL will build the rest query url of creation
//the create url is endpoint/scaling_policy
func createURL(client *golangsdk.ServiceClientExtension) string {
	return client.ServiceURL(client.ProjectID, resourcePath)
}

//deleteURL will build the url of deletion
//its pattern is endpoint/scaling_policy/<policy-id>
func deleteURL(client *golangsdk.ServiceClientExtension, id string) string {
	return client.ServiceURL(client.ProjectID, resourcePath, id)
}

//getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClientExtension, id string) string {
	return client.ServiceURL(client.ProjectID, resourcePath, id)
}

func updateURL(c *golangsdk.ServiceClientExtension, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id)
}
