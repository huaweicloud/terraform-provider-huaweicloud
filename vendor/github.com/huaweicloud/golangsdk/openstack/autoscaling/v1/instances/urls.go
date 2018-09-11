package instances

import (
	"github.com/huaweicloud/golangsdk"
)

const resourcePath = "scaling_group_instance"

//getURL will build the querystring by which can be able to query all the instances
//of group
func listURL(client *golangsdk.ServiceClient, groupID string) string {
	return client.ServiceURL(resourcePath, groupID, "list")
}

//deleteURL will build the query url by which can be able to delete an instance from
//the group
func deleteURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(resourcePath, id)
}

//batchURL will build the query url by which can be able to batch add or delete
//instances
func batchURL(client *golangsdk.ServiceClient, groupID string) string {
	return client.ServiceURL(resourcePath, groupID, "action")
}
