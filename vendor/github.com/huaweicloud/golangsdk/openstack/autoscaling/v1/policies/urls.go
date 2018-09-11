package policies

import (
	"github.com/huaweicloud/golangsdk"
)

const resourcePath = "scaling_policy"

//createURL will build the rest query url of creation
//the create url is endpoint/scaling_policy
func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(resourcePath)
}

//deleteURL will build the url of deletion
//its pattern is endpoint/scaling_policy/<policy-id>
func deleteURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(resourcePath, id)
}

//getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(resourcePath, id)
}

func updateURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}
