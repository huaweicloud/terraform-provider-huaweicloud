package policies

/*
Package policies enables management and retrieval of backup policies
VBS service.

Example to List Policies

	listopts := policies.ListOpts{PolicyID:"ef5c5859-a8f4-48cb-869a-dbbc466cd6b6",Frequency:0}
	list,err := policies.List(client,listopts)
	if err != nil {
		panic(err)
	}


Example to Create a Policy

	creatopts := policies.CreateOpts{PolicyName: "Demo_policy", ScheduledPolicy: policies.ScheduledPolicy{StartTime: "12:00", Status: "ON", Frequency: 1, RententionNum: 12, RemainFirstBackup: "Y",}, Tags:[]policies.Tag{{Key:"key",Value:"value"}}}
	create, err := policies.Create(client, creatopts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Policy

	delete := policies.Delete(client,"c776176d-eb3c-4f48-902a-da3a12c4fea9")
	if delete.Err != nil {
		panic(delete.Err)
	}

Example to Associate a resource to a Policy

	assopts := policies.AssociateOpts{PolicyID:"5b549fad-c4e5-4d7e-83b9-eea366f27017",Resources:[]policies.AssociateResource{{ResourceID:"bdec76de-3cca-46b4-8b71-a333467a1b79",ResourceType:"volume"},{ResourceID:"286b8b84-6640-4f6f-acde-2a58e490f371",ResourceType:"volume"}}}
	associate,err := policies.Associate(client,assopts).ExtractResource()

Example to disassociate a resource to a Policy
	disassopts := policies.DisassociateOpts{Resources: []policies.DisassociateResource{{ResourceID: "bdec76de-3cca-46b4-8b71-a333467a1b79"},{ResourceID: "286b8b84-6640-4f6f-acde-2a58e490f371"}}}
	disassociate,err := policies.Disassociate(client,"5b549fad-c4e5-4d7e-83b9-eea366f27017",disassopts).ExtractResource()

*/
