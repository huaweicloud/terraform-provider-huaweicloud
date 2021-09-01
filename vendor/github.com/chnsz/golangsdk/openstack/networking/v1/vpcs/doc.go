/*
Package vpcs enables management and retrieval of Vpcs
VPC service.

Example to List Vpcs

	listOpts := vpcs.ListOpts{}
	allVpcs, err := vpcs.List(vpcClient, listOpts)
	if err != nil {
		panic(err)
	}

	for _, vpc := range allVpcs {
		fmt.Printf("%+v\n", vpc)
	}

Example to Create a Vpc

	createOpts := vpcs.CreateOpts{
		Name:         "vpc_1",
		CIDR:         "192.168.0.0/24"

	}

	vpc, err := vpcs.Create(vpcClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a Vpc

	vpcID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"

	updateOpts := vpcs.UpdateOpts{
		Name:         "vpc_2",
		CIDR:         "192.168.0.0/23"
	}

	vpc, err := vpcs.Update(vpcClient, vpcID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Vpc

	vpcID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"
	err := vpcs.Delete(vpcClient, vpcID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package vpcs
