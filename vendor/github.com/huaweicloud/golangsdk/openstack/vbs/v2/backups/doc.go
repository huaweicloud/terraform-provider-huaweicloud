/*
Package backups enables management and retrieval of Backups
VBS service.

Example to List Backups

	listOpts := backups.ListOpts{}
	allBackups, err := backups.List(vbsClient, listOpts)
	if err != nil {
		panic(err)
	}

	for _, backup := range allBackups {
		fmt.Printf("%+v\n", backup)
	}

Example to Get a Backup

   getbackup,err:=backups.Get(vbsClient, "6149e448-dcac-4691-96d9-041e09ef617f").Extract()
   if err != nil {
         panic(err)
		}

   fmt.Println(getbackup)

Example to Create a Backup

	createOpts := backups.CreateOpts{
		Name:"backup-test",
		VolumeId:"5024a06e-6990-4f12-9dcc-8fe26b01a710",
	}

	jobInfo, err := backups.Create(vbsClient, createOpts).ExtractJobResponse()
	if err != nil {
		panic(err)
	}

    err1 := backups.WaitForJobSuccess(client, int(120), jobInfo.JobID)
	if err1 != nil {
		panic(err1)
	}

	Label := "backup_id"
    entity, err2 := backups.GetJobEntity(client, jobInfo.JobID, Label)
	fmt.Println(entity)
	if err2 != nil {
		panic(err2)
	}

Example to Delete a Backup

	backupID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"
	err := backups.Delete(vbsClient, backupID).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Restore a Backup

	restoreOpts := backups.BackupRestoreOpts{VolumeId:"5024a06e-6990-4f12-9dcc-8fe26b01a710"}

	restore,err := backups.CreateBackupRestore(vbsClient,"87566ed6-72cb-4053-aa6e-6f6216b3d507",backup).ExtractBackupRestore()
	if err != nil {
		panic(err)
	}
*/
package backups
