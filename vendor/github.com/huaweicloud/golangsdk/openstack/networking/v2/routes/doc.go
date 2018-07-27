/*
Package routes enables management and retrieval of Routes
Route service.

Example to List Routes

	listroute:=routes.ListOpts{VPC_ID:"93e94d8e-31a6-4c22-bdf7-8b23c7b67329"}
	out,err:=routes.List(client,listroute)
	fmt.Println(out[0].RouteID)


Example to Create a Route

	route:=routes.CreateOpts{
		Type:"peering",
		NextHop:"d2dea4ba-e988-4e9c-8162-652e74b2560c",
		Destination:"192.168.0.0/16",
		VPC_ID:"3127e30b-5f8e-42d1-a3cc-fdadf412c5bf"}
	outroute,err:=routes.Create(client,route).Extract()
	fmt.Println(outroute)


Example to Delete a Route

	out:=routes.Delete(client,"39a07dcb-f30e-41c1-97ac-182c8f0d43c1")
		fmt.Println(out)
*/
package routes
