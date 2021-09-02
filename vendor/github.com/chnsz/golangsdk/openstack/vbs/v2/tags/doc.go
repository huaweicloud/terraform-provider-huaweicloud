package tags

/*
Package tags enables management and retrieval of policy's tags
VBS service.

Example to Add tag a Policy

	createopts := tags.CreateOpts{
						Tag:tags.Tag{Key:"test",Value:"demo"}
				 }

	create,err := tags.Create(client,"ed8b9f73-4415-494d-a54e-5f3373bc353d",createopts).Extract()

Example to Get tag of a Policy

   get,err := tags.Get(client,"ed8b9f73-4415-494d-a54e-5f3373bc353d").Extract()
	fmt.Println(get,err)

Example to delete a tag of a Policy

    delete := tags.Delete(client,"5b549fad-c4e5-4d7e-83b9-eea366f27017","ECS")
	fmt.Println(delete)

Example to ListResources policy details based on tags

   queryOpts := tags.ListOpts{
					Action:"filter",
					NotAnyTags:[]tags.RespTags{{Key:"newKey",Values:[]string{"volumeSphere"}}}
				}

   query,err :=tags.ListResources(client,queryOpts).ExtractResources()
   fmt.Println(query,err)

Example to perform batch tag actions on a backup policy

	actionopts := tags.BatchOpts{
						Action:"update",
						RespTags:[]tags.Tag{{Key:"k22",Value:"v22"}}
				}

    action := tags.BatchAction(client,"ed8b9f73-4415-494d-a54e-5f3373bc353d",actionopts)
*/
