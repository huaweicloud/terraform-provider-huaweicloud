/*
Package node pools enables management and retrieval of node pools
CCE service.

Example to List node pools

    clusterID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"

	listNodePools := nodepools.ListOpts{}
    allNodePools, err := nodepools.List(client,clusterID).ExtractNodePool(listNodePools)

	if err != nil {
		panic(err)
	}

	for _, nodepool := range allNodePools {
		fmt.Printf("%+v\n", nodepool)
	}

Example to Create a node pool

    createOpts := nodepools.CreateOpts{Kind:"NodePool",
	   					   ApiVersion:"v3",
	   					   Metadata:nodepools.CreateMetaData{Name:"nodepool_1"},
	   					   Spec: nodepools.CreateSpec{
								NodeTemplate: :nodes.Spec{Flavor:"s1.medium",
	   					                   Az:"az1.dc1",
	   					                   Login:nodes.LoginSpec{"myKeypair"},
	   					                   Count:1,
	   					                   RootVolume:nodes.VolumeSpec{Size:10,VolumeType:"SATA"},
	   					                   DataVolumes:[]nodes.VolumeSpec{{Size:10,VolumeType:"SATA"}},
								},
								Autoscaling: nodepools.AutoscalingSpec{Enable:true,MinNodeCount:1,MaxNodeCount:3},
								InitialNodeCount:1,
							}
		}
 	nodepool, err := nodepools.Create(client,clusterID,createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a node pool

    clusterID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"

	nodeID := "3c8e5957-649f-477b-9e5b-f1f75b21c011"

	updateOpts := nodepools.UpdateOpts{
		Kind:       "NodePool",
		ApiVersion: "v3",
		Metadata: nodepools.UpdateMetaData{Name:"nodepool_1"},
		Spec: nodepools.UpdateSpec{
			InitialNodeCount: d.Get("initial_node_count").(int),
			Autoscaling: nodepools.AutoscalingSpec{Enable:true,MinNodeCount:2,MaxNodeCount:9},
		},
	}
	nodepool,err := nodepools.Update(client,clusterID,nodeID,updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a node pool

	clusterID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"

	nodePoolID := "3c8e5957-649f-477b-9e5b-f1f75b21c011"

	err := nodepools.Delete(client,clusterID,nodePoolID).Extract()
	if err != nil {
		panic(err)
	}
*/

package nodepools
