/*
Package shares enables management and retrieval of shares
Share service.
Example to List Shares
	listshares := shares.ListOpts{}
	allshares, err := shares.List(client,listshares)
	if err != nil {
		panic(err)
	}
	fmt.Println(allshares)

Example to Create a share
	share:=shares.CreateOpts{
			Name:"sfs-test",
			ShareProto:"NFS",
			Size:1,
		}
	outshare,err:=shares.Create(client,share).Extract()
	if err != nil {
				panic(err)
			}
	fmt.Println(outshare)

Example to Update a share
	updateshare:=shares.UpdateOpts{DisplayName:"sfs-test-1",DisplayDescription:"test sfs"}
	out,err:=shares.Update(client,"6149e448-dcac-4691-96d9-041e09ef617f",updateshare).Extract()
	if err != nil {
			panic(err)
		}
	fmt.Println(out)

Example to Delete a share
	out:=shares.Delete(client,"6149e448-dcac-4691-96d9-041e09ef617f")
	fmt.Println(out)
	if err != nil {
		panic(err)
	}

Example to Get share
	getshare,err:=shares.Get(client, "6149e448-dcac-4691-96d9-041e09ef617f").Extract()
	fmt.Println(getshare)
	if err != nil {
			panic(err)
		}
Example to Allow Access
createSet:=shares.GrantAccessOpts{AccessLevel:"rw",AccessTo:"5232f396-d6cc-4a81-8de3-afd7a7ecdfd8",AccessType:"cert"}
	access,err:=shares.GrantAccess(client,"dff2df5f-00e7-4517-ac32-1d0ab8dc0d68",createSet).Extract()
	fmt.Println(access)
	if err != nil {
			panic(err)
		}

Example to Deny Access
	deleteSet := shares.DeleteAccessOpts{AccessID:"fc32500f-fa78-4f06-8caf-06ad7fb9726c"}
	remove:=shares.DeleteAccess(client,"1b8facf8-b822-4349-a033-e078b2a84b7f",deleteSet)
	fmt.Println(remove)
	if err != nil {
			panic(err)
		}

Example to Get Access Rule Detail
	rule_list,err:= shares.ListAccessRights(client,"42381b5b-f8cb-445e-9465-89a718e071a7").ExtractAccessRights()
	if err != nil {
			panic(err)
		}
	fmt.Println(rule_list)

Example to Get Mount Location Details
	mount, err := shares.GetExportLocations(client, "dff2df5f-00e7-4517-ac32-1d0ab8dc0d68").ExtractExportLocations()
	fmt.Println(mount)
	if err != nil {
			panic(err)
		}

Example to Extend share
	extendsfs:=shares.ExpandOpts{OSExtend: shares.OSExtendOpts{NewSize: 512}}
	shares.Expand(client,"45a3af18-8ab0-405c-9ead-06c51a415f79",extendsfs)

Example to Shrink share
	shrinksfs:=shares.ShrinkOpts{OSShrink: shares.OSShrinkOpts{NewSize: 8}}
	shares.Shrink(client,"45a3af18-8ab0-405c-9ead-06c51a415f79",shrinksfs)
*/
package shares
