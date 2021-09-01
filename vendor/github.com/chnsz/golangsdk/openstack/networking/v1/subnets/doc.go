/*
Package Subnets enables management and retrieval of Subnets

Example to List Vpcs

	listOpts := subnets.ListOpts{}
	allSubnets, err := subnets.List(subnetClient, listOpts)
	if err != nil {
		panic(err)
	}

	for _, subnet := range allSubnets {
		fmt.Printf("%+v\n", subnet)
	}

Example to Create a Vpc

	createOpts := subnets.CreateOpts{
		Name:          "test_subnets",
		CIDR:          "192.168.0.0/16"
		GatewayIP:	   "192.168.0.1"
		PRIMARY_DNS:   "8.8.8.8"
		SECONDARY_DNS: "8.8.4.4"
		AvailabilityZone:"eu-de-02"
		VPC_ID:"3b9740a0-b44d-48f0-84ee-42eb166e54f7"

	}
	vpc, err := subnets.Create(subnetClient, createOpts).Extract()

	if err != nil {
		panic(err)
	}

Example to Update a Vpc

	subnetID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"

	updateOpts := subnets.UpdateOpts{
		Name:          "testsubnet",
	}

	subnet, err := subnets.Update(subnetClient, subnetID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Vpc

	subnetID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"

	err := subnets.Delete(subnetClient, subnetID).ExtractErr()

	if err != nil {
		panic(err)
	}
*/
package subnets
