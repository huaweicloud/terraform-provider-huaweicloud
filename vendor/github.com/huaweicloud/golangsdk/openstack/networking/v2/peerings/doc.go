/*
Package peerings enables management and retrieval of vpc peering connections

Example to List a Vpc Peering Connections
	   listOpts:=peerings.ListOpts{}

		peering,err :=peerings.List(client,sub).AllPages()

		peerings,err:=peerings.ExtractPeerings(peering)


		if err != nil{
			fmt.Println(err)
		}

Example to Get a Vpc Peering Connection

       	peeringID := "6bbacb0f-9f94-4fe8-a6b6-1818bdccb2a3"


		peering,err :=peerings.Get(client,peeringID).Extract()


		if err != nil{
			fmt.Println(err)
		}



Example to Accept a Vpc Peering Connection Request
 // Note:- The TenantId should be of accepter

    peeringID := "6bbacb0f-9f94-4fe8-a6b6-1818bdccb2a3"

    peering,err:=peerings.Accept(client,peeringID).ExtractResult()

	if err != nil{
		fmt.Println(err)
	}


Example to Reject a Vpc Peering Connection Request
 // Note:- The TenantId should be of accepter
    peeringID := "6bbacb0f-9f94-4fe8-a6b6-1818bdccb2a3"

    peering,err:=peerings.Reject(client,peeringID).ExtractResult()

	if err != nil{
		fmt.Println(err)
	}

Example to Create a Vpc Peering Connection

	RequestVpcInfo:=peerings.VpcInfo{VpcId:"3127e30b-5f8e-42d1-a3cc-fdadf412c5bf"}
	AcceptVpcInfo:=peerings.VpcInfo{"c6efbdb7-dca4-4178-b3ec-692f125c1e25","17fbda95add24720a4038ba4b1c705ed"}

	opt:=peerings.CreateOpts{"C2C_test",RequestVpcInfo,AcceptVpcInfo}

	peering,err:=peerings.Create(client,opt).Extract()

	if err != nil{
		fmt.Println(err)
	}


Example to Update a VpcPeeringConnection

	peeringID := "6bbacb0f-9f94-4fe8-a6b6-1818bdccb2a3"

	updateOpts:=peerings.UpdateOpts{"C2C_tes1"}

	peering,err:=peerings.Update(client,peeringID,updateOpts).Extract()

	if err != nil{
		fmt.Println(err)
	}


Example to Delete a VpcPeeringConnection

	peeringID := "6bbacb0f-9f94-4fe8-a6b6-1818bdccb2a3"
	err := peerings.Delete(client,"6bbacb0f-9f94-4fe8-a6b6-1818bdccb2a3")
	if err != nil {
		panic(err)
	}
*/
package peerings
