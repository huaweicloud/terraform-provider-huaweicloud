package associations

import "github.com/chnsz/golangsdk"

func enableURL(client *golangsdk.ServiceClient, instanceId, routeTableId string) string {
	return client.ServiceURL("enterprise-router", instanceId, "route-tables", routeTableId, "associate")
}

func queryURL(client *golangsdk.ServiceClient, instanceId, routeTableId string) string {
	return client.ServiceURL("enterprise-router", instanceId, "route-tables", routeTableId, "associations")
}

func disableURL(client *golangsdk.ServiceClient, instanceId, routeTableId string) string {
	return client.ServiceURL("enterprise-router", instanceId, "route-tables", routeTableId, "disassociate")
}
