package routes

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient, routeTableId string) string {
	return client.ServiceURL("enterprise-router/route-tables", routeTableId, "static-routes")
}

func resourceURL(client *golangsdk.ServiceClient, routeId, routeTableId string) string {
	return client.ServiceURL("enterprise-router/route-tables", routeTableId, "static-routes", routeId)
}
