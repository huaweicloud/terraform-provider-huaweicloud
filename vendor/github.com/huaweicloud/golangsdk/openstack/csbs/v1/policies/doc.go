/*
Package backup policies enables management and retrieval of
backup servers periodically.

Example to List Backup Policies
	listpolicies := policies.ListOpts{}
	allpolicies, err := policies.List(client,listpolicies)
	if err != nil {
		panic(err)
	}
	fmt.Println(allpolicies)



Example to Create a Backup Policy
	policy:=policies.CreateOpts{
				Name : "c2c-policy",
				Description : "My plan",
				ProviderId : "fc4d5750-22e7-4798-8a46-f48f62c4c1da",
				Parameters : policies.PolicyParam{
				Common:map[string]interface{}{},
				},
				ScheduledOperations : []policies.ScheduledOperation{ {
					Name:  "my-backup",
					Description: "My backup policy",
					Enabled: true,
					OperationDefinition: policies.OperationDefinition{
						MaxBackups: 5,
					},
					Trigger: policies.Trigger{
						Properties : policies.TriggerProperties{
							Pattern : "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n",
						},
					},
					OperationType: "backup",
				}},
				Resources : []policies.Resource{{
					Id:  "9422f270-6fcf-4ba2-9319-a007f2f63a8e",
					Type: "OS::Nova::Server",
					Name: "resource4"
				}},
		}
	out,err:=policies.Create(client,policy).Extract()
	fmt.Println(out)
	fmt.Println(err)


Example to Update a Backup Policy
	updatepolicy:=policies.UpdateOpts{
								Name:"my-plan-c2c-update",
								Parameters : policies.PolicyParamUpdate{
									Common:map[string]interface{}{},
								},
								ScheduledOperations:[]policies.ScheduledOperationToUpdate{{
									Id:"b70c712d-f48b-43f7-9a0f-3bab86d59149",
									Name:"my-backup-policy",
									Description:"My backup policy",
									Enabled:true,
									OperationDefinition:policies.OperationDefinition{
										RetentionDurationDays:-1,
										MaxBackups:20,
									},
									Trigger:policies.Trigger{
										Properties:policies.TriggerProperties{
											Pattern:"BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"}}
										}
									}
								}
							}
		out,err:=policies.Update(client,"5af626d2-19b9-4dc4-8e95-ddba008318b3",updatepolicy).Extract()
		fmt.Println(out)

Example to Delete a Backup Policy
	out:=policies.Delete(client,"16d4bf9e-85b2-41e2-a482-e48ace2ad726")
	fmt.Println(out)
	if err != nil {
		panic(err)
	}

Example to Get Backup Policy
	result:=policies.Get(client,"5af626d2-19b9-4dc4-8e95-ddba008318b3")
	out,err:=result.Extract()
	fmt.Println(out)

*/
package policies
