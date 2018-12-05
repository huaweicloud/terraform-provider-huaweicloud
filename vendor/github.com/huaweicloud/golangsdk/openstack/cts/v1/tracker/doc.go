/*
Package tracker provides operation records for cloud service resources.

Example to List Tracker
	listTracker := tracker.ListOpts{}
	allTracker, err := tracker.List(client,listTracker)
	if err != nil {
		panic(err)
	}
	fmt.Println(allTracker)



Example to Create a Tracker
	createTracker:=tracker.CreateOpts{
		BucketName: "obs-e51d",
		FilePrefixName: "mytracker",
		SimpleMessageNotification:tracker.SimpleMessageNotification{
			IsSupportSMN: true,
			TopicID: "urn:smn:eu-de:626ce20e52a346c090b09cffc3e038e5:c2c-topic",
			IsSendAllKeyOperation: false,
			Operations: []string{"login"},
			NeedNotifyUserList: []string{"user1","user2"},
	}}
	out,err:=tracker.Create(client, createTracker).Extract()
	fmt.Println(out)
	fmt.Println(err)


Example to Update a Tracker
	updateTracker:=tracker.UpdateOpts{
		BucketName : "ciros-img",
		FilePrefixName : "mytracker",
		Status : "disabled",
		SimpleMessageNotification:tracker.SimpleMessageNotification{
			IsSupportSMN: false,
			TopicID: "urn:smn:eu-de:626ce20e52a346c090b09cffc3e038e5:c2c-topic",
			IsSendAllKeyOperation:false,
			Operations: []string{"delete","create","login"},
			NeedNotifyUserList:[]string{"user1","user2"},
		},
		   }
	out,err:=tracker.Update(client, updateTracker).Extract()
	fmt.Println(out)


Example to Delete a Tracker
	out:= tracker.Delete(client).ExtractErr()
	fmt.Println(out)


*/
package tracker
