package clusters

type RequestResp struct {
	// Cluster ID.
	ClusterId string `json:"cluster_id"`
}

type Cluster struct {
	// Cluster datastore.
	Datastore Datastore `json:"datastore"`
	// The current status list of the cluster:
	//   CREATING
	//   GROWING
	Actions []string `json:"actions"`
	// Whether the openTSDB is enabled.
	OpenTSDBEnabled bool `json:"enable_openTSDB"`
	// Whether the Lemon is enabled.
	LemonEnabled bool `json:"enable_lemon"`
	// The name of the CloudTable cluster.
	Name string `json:"cluster_name"`
	// The number of RegionServers.
	CUNum string `json:"cu_num"`
	// The number of TSD nodes.
	TSDNum string `json:"tsd_num"`
	// The number of Lemon nodes.
	LemonNum string `json:"lemon_num"`
	// Cluster bottom storage type:
	//   OBS
	//   HDFS
	StorageType string `json:"storage_type"`
	// Cluster storage quota.
	StorageQuota string `json:"storage_quota"`
	// Storage space currently in use.
	StorageUsed string `json:"used_storage_size"`
	// Whether the IAM auth is enabled.
	IAMAuthEnabled bool `json:"auth_mode"`
	// The time when the disk was updated.
	UpdatedAt string `json:"updated"`
	// The time when the disk was created.
	CreateAt string `json:"created"`
	// Cluster ID.
	ID string `json:"cluster_id"`
	// Cluster status.
	//   100 Creating
	//   200 Running
	//   300 Abnormal
	//   303 Creation failed
	//   400 Deleted
	//   800 Frezon
	Status string `json:"status"`
	// Intranet OpenTSDB connection access address.
	OpenTSDBLink string `json:"openTSDB_link"`
	// OpenTSDB public network endpoint address.
	TSDPublicEndpoint string `json:"tsd_public_endpoint"`
	// Intranet Lemon connection access address.
	LemonLink string `json:"lemon_link"`
	// Intranet ZooKeeper connection access address.
	ZookeeperLink string `json:"zookeeper_link"`
	// HBase connection access address on the public network.
	HbasePublicEndpoint string `json:"hbase_public_endpoint"`
	// Whether the cluster is frozen.
	//   false
	//   true
	IsFrozen string `json:"is_frozen"`
	// The VPC where the cluster is located.
	VpcId string `json:"vpc_id"`
	// The ID of the network where the CloudTable cluster is located.
	SubnetId string `json:"subnet_id"`
	// The ID of the security group to which the CloudTable belongs.
	SecurityGroupId string `json:"security_group_id"`
	// The ID of the availability zone where the cluster is located.
	AvailabilityZone string `json:"availability_zone"`
}
