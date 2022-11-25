package routetables

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient, instanceId string) string {
	return client.ServiceURL("enterprise-router", instanceId, "route-tables")
}

func resourceURL(client *golangsdk.ServiceClient, instanceId, routeTableId string) string {
	return client.ServiceURL("enterprise-router", instanceId, "route-tables", routeTableId)
}
