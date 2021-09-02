/*
Package backup enables management and retrieval of
back up resources.

Example to List Backup
	listbackup := backup.ListOpts{ID: "7b99acfd-18c3-4f26-9d39-b4ebd2ea3e12"}
	allbackups, err := backup.List(client,listbackup)
	if err != nil {
		panic(err)
	}
	fmt.Println(allbackups)



Example to Create a Backup
	createBackup:=backup.CreateOpts{BackupName: "c2c-backup", Description: "mybackup"}
	out,err:=backup.Create(client,"fc4d5750-22e7-4798-8a46-f48f62c4c1da", "f8ddc472-cf00-4384-851e-5f2a68c33762",
							createBackup).Extract()
	fmt.Println(out)
	fmt.Println(err)

Example to Query if resources can be backed up
	createQuery:=backup.ResourceBackupCapOpts{CheckProtectable:[]backup.ResourceCapQueryParams{{ResourceId: "069e678a-f1d1-4a38-880b-459bde82fcc6",
											ResourceType: "OS::Nova::Server"}}}
	out,err:=backup.QueryResourceBackupCapability(client,"fc4d5750-22e7-4798-8a46-f48f62c4c1da",
		createQuery).ExtractQueryResponse()
	fmt.Println(out)
	fmt.Println(err)


Example to Delete a Backup
	out:=backup.Delete(client,"fc4d5750-22e7-4798-8a46-f48f62c4c1da")
	fmt.Println(out)
	if err != nil {
		panic(err)
	}

Example to Get Backup
	result:=backup.Get(client,"7b99acfd-18c3-4f26-9d39-b4ebd2ea3e12")
	out,err:=result.ExtractBackup()
	fmt.Println(out)

*/
package backup
