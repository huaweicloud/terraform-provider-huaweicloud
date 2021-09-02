/*
Package nodes enables management and retrieval of nodes
CCE service.

Example to List nodes

    clusterID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"

	listNodes := nodes.ListOpts{}
    allNodes, err := nodes.List(client,clusterID).ExtractNode(listNodes)

	if err != nil {
		panic(err)
	}

	for _, node := range allNodes {
		fmt.Printf("%+v\n", node)
	}

Example to Create a node

    createOpts := nodes.CreateOpts{Kind:"Node",
	   					   ApiVersion:"v3",
	   					   Metadata:nodes.CreateMetaData{Name:"node_1"},
	   					   Spec:nodes.Spec{Flavor:"s1.medium",
	   					                   Az:"az1.dc1",
	   					                   Login:nodes.LoginSpec{"myKeypair"},
	   					                   Count:1,
	   					                   RootVolume:nodes.VolumeSpec{Size:10,VolumeType:"SATA"},
	   					                   DataVolumes:[]nodes.VolumeSpec{{Size:10,VolumeType:"SATA"}},
								},
		}
 	node,err := nodes.Create(client,clusterID,createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a cluster

    clusterID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"

	nodeID := "3c8e5957-649f-477b-9e5b-f1f75b21c011"

	updateOpts := nodes.UpdateOpts{Metadata:nodes.UpdateMetadata{Name:"node_1"}}
	node,err := nodes.Update(client,clusterID,nodeID,updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a cluster

	clusterID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"

	nodeID := "3c8e5957-649f-477b-9e5b-f1f75b21c011"

	err := nodes.Delete(client,clusterID,nodeID).Extract()
	if err != nil {
		panic(err)
	}
*/

package nodes
