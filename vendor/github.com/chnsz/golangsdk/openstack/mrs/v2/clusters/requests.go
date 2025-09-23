package clusters

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

// CreateOpts is a structure representing information of the cluster creation.
type CreateOpts struct {
	// Region of the cluster.
	Region string `json:"region" required:"true"`
	// Availability zone name.
	AvailabilityZone string `json:"availability_zone" required:"true"`
	// Cluster name, which can contain 2 to 64 characters.
	// Only letters, digits, hyphens (-), and underscores (_) are allowed.
	ClusterName string `json:"cluster_name" required:"true"`
	// Cluster type. The options are as follows:
	//   ANALYSIS: analysis cluster
	//   STREAMING: streaming cluster
	//   MIXED: hybrid cluster
	//   CUSTOM: customized cluster, which is supported only by MRS 3.x.
	ClusterType string `json:"cluster_type" required:"true"`
	// Cluster version.
	// Possible values are as follows:
	//   MRS 1.8.10
	//   MRS 1.9.2
	//   MRS 2.1.0
	//   MRS 3.0.2
	ClusterVersion string `json:"cluster_version" required:"true"`
	// Name of the VPC where the subnet locates
	// Perform the following operations to obtain the VPC name from the VPC management console:
	// Log in to the management console.
	// Click Virtual Private Cloud and select Virtual Private Cloud from the left list.
	// On the Virtual Private Cloud page, obtain the VPC name from the list.
	VpcName string `json:"vpc_name" required:"true"`
	VpcId   string `json:"vpc_id,omitempty"`
	// List of component names, which are separated by commas (,). The options are as follows:
	// MRS 3.0.5
	//   ANALYSIS: Hadoop,Spark2x,HBase,Hive,Hue,Loader,Flink,Oozie,ZooKeeper,Ranger,Tez,Impala,Presto,Kudu,Alluxio
	//   STREAMING: Kafka,Storm,Flume,ZooKeeper,Ranger
	//   MIXED: Hadoop,Spark2x,HBase,Hive,Hue,Loader,Flink,Oozie,ZooKeeper,Ranger,Tez,Impala,Presto,Kudu,Alluxio,
	//     Kafka,Storm,Flume
	//   CUSTOM: Hadoop,Spark2x,HBase,Hive,Hue,Loader,Kafka,Storm,Flume,Flink,Oozie,ZooKeeper,Ranger,Tez,Impala,
	//     Presto,ClickHouse，Kudu,Alluxio
	// MRS 1.9.2
	//   ANALYSIS: Presto,Hadoop,Spark,HBase,Opentsdb,Hive,Hue,Loader,Tez,Flink,Alluxio,Ranger
	//   STREAMING: Kafka,KafkaManager,Storm,Flume
	Components string `json:"components" required:"true"`
	// Node login mode.
	// PASSWORD: password-based login. If this value is selected, node_root_password cannot be left blank.
	// KEYPAIR: specifies the key pair used for login. If this value is selected, node_keypair_name cannot be left blank.
	LoginMode string `json:"login_mode" required:"true"`
	// Password of the MRS Manager administrator.
	// The password can contain 8 to 26 charactors.
	// The password must contain lowercase letters, uppercase letters, digits, spaces
	// and the special characters: !?,.:-_{}[]@$^+=/.
	// the password cannot be the username or the username spelled backwards.
	ManagerAdminPassword string `json:"manager_admin_password" required:"true"`
	// Information about the node groups in the cluster. For details about the parameters, see Table 5.
	NodeGroups []NodeGroupOpts `json:"node_groups" required:"true"`
	// Running mode of an MRS cluster
	// SIMPLE: normal cluster. In a normal cluster, Kerberos authentication is disabled, and users can use all functions provided by the cluster.
	// KERBEROS: security cluster. In a security cluster, Kerberos authentication is enabled, and common users cannot use the file management and job management functions of an MRS cluster or view cluster resource usage and the job records of Hadoop and Spark. To use more cluster functions, the users must contact the Manager administrator to assign more permissions.
	SafeMode string `json:"safe_mode" required:"true"`
	// Jobs can be submitted when a cluster is created. Currently, only one job can be created. For details about job parameters, see Table 9.
	AddJobs []JobOpts `json:"add_jobs,omitempty"`
	// Whether to create the default security group of the MR S security group, the default is false.
	// If true, no matter what 'security_groups_id' is set, the cluster will create a default security group.
	AutoCreateSecGroup string `json:"auto_create_default_security_group,omitempty"`
	// Bootstrap action script information. For more parameter description, see Table 8.
	// MRS 1.7.2 or later supports this parameter.
	BootstrapScripts []ScriptOpts `json:"bootstrap_scripts,omitempty"`
	// Charging type information.
	ChargeInfo *ChargeInfo `json:"charge_info,omitempty"`
	// Indicates the enterprise project ID.
	// When creating a cluster, associate the enterprise project ID with the cluster.
	// The default value is 0, indicating the default enterprise project.
	// To obtain the enterprise project ID, see the id value in the enterprise_project field data structure table in section Querying the Enterprise Project List of the Enterprise Management API Reference.
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
	// An EIP bound to an MRS cluster can be used to access MRS Manager. The EIP must have been created and must be in the same region as the cluster.
	EipAddress string `json:"eip_address,omitempty"`
	// ID of the bound EIP. This parameter is mandatory when eip_address is configured. To obtain the EIP ID, log in to the VPC console, choose Network > Elastic IP and Bandwidth > Elastic IP, click the EIP to be bound, and obtain the ID in the Basic Information area.
	EipId string `json:"eip_id,omitempty"`
	// Indicate whether it is a dedicated cloud resource, the default is false.
	IsDecProject *bool `json:"is_dec_project,omitempty"`
	// Specifies whether to collect logs when cluster creation fails:
	//   0: Do not collect.
	//   1: Collect.
	// The default value is 1, indicating that OBS buckets will be created and only used to collect logs that record MRS cluster creation failures.
	LogCollection *int `json:"log_collection,omitempty"`
	// Name of the agency bound to a cluster node by default. The value is fixed to MRS_ECS_DEFAULT_AGENCY.
	// An agency allows ECS or BMS to manage MRS resources. You can configure an agency of the ECS type to automatically obtain the AK/SK to access OBS.
	// The MRS_ECS_DEFAULT_AGENCY agency has the OBS OperateAccess permission of OBS and the CES FullAccess (for users who have enabled fine-grained policies), CES Administrator, and KMS Administrator permissions in the region where the cluster is located.
	MrsEcsDefaultAgency string `json:"mrs_ecs_default_agency,omitempty"`
	// Password of user root for logging in to a cluster node.
	// The password can contain 8 to 26 charactors.
	// The password must contain lowercase letters, uppercase letters, digits, spaces
	// and the special characters: !?,.:-_{}[]@$^+=/.
	// the password cannot be the username or the username spelled backwards.
	NodeRootPassword string `json:"node_root_password,omitempty"`
	// Name of a key pair You can use a key pair to log in to the Master node in the cluster.
	NodeKeypair string `json:"node_keypair_name,omitempty"`
	// Security group ID of the cluster
	// If this parameter is left blank, MRS automatically creates a security group, whose name starts with mrs_{cluster_name}.
	// If this parameter is not left blank, a fixed security group is used to create a cluster. The transferred ID must be the security group ID owned by the current tenant. The security group must include an inbound rule in which all protocols and all ports are allowed and the source is the IP address of the specified node on the management plane.
	SecurityGroupsIds string `json:"security_groups_id,omitempty"`
	// Subnet ID.
	SubnetId string `json:"subnet_id,omitempty"`
	// Subnet name.
	// Required if SubnetID is empty.
	SubnetName string `json:"subnet_name,omitempty"`
	// Specifies the template used for node deployment when the cluster type is CUSTOM.
	// mgmt_control_combined_v2: template for jointly deploying the management and control nodes. The management and control roles are co-deployed on the Master node, and data instances are deployed in the same node group. This deployment mode applies to scenarios where the number of control nodes is less than 100, reducing costs.
	// mgmt_control_separated_v2: The management and control roles are deployed on different master nodes, and data instances are deployed in the same node group. This deployment mode is applicable to a cluster with 100 to 500 nodes and delivers better performance in high-concurrency load scenarios.
	// mgmt_control_data_separated_v2: The management role and control role are deployed on different Master nodes, and data instances are deployed in different node groups. This deployment mode is applicable to a cluster with more than 500 nodes. Components can be deployed separately, which can be used for a larger cluster scale.
	TemplateId string `json:"template_id,omitempty"`
	// Cluster tag For more parameter description, see Table 4.
	// A maximum of 10 tags can be added to a cluster.
	Tags []tags.ResourceTag `json:"tags,omitempty"`
	// the component configurations of MRS cluster.
	ComponentConfigs []ComponentConfigOpts `json:"component_configs,omitempty"`
	// When deploying components such as Hive and Ranger, you can associate data connections and store metadata in associated databases
	ExternalDatasources []ExternalDatasource `json:"external_datasources,omitempty"`
	// The OBS path to which cluster logs are dumped.
	// This parameter is available only for cluster versions that support dumping cluster logs to OBS.
	LogURI string `json:"log_uri,omitempty"`
	// The alarm configuration of the cluster.
	SMNNotifyConfig *SMNNotifyConfigOpts `json:"smn_notify,omitempty"`
}

// SMNNotifyConfigOpts is a structure representing the alarm configuration information.
type SMNNotifyConfigOpts struct {
	// The Uniform Resource Name (URN) of the topic.
	TopicURN string `json:"topic_urn" required:"true"`
	// The subscription rule name.
	SubscriptionName string `json:"subscription_name" required:"true"`
}

type ExternalDatasource struct {
	ConnectorId   string `json:"connector_id,omitempty"`
	ComponentName string `json:"component_name,omitempty"`
	/**
	Component role type. The options are as follows:
		hive_metastore: Hive Metastore role
		hive_data: Hive role
		hbase_data: HBase role
		ranger_data: Ranger role
	**/
	RoleType string `json:"role_type,omitempty"`
	/**
	Data connection type. The options are as follows:
		LOCAL_DB: local metadata
		RDS_POSTGRES: RDS PostgreSQL database
		RDS_MYSQL: RDS MySQL database
		gaussdb-mysql: GaussDB(for MySQL)
	**/
	SourceType string `json:"source_type,omitempty"`
}

// ChargeInfo is a structure representing billing information.
type ChargeInfo struct {
	// Billing mode.
	// The valid values are as follows:
	//   postPaid: indicates the pay-per-use billing mode.
	//   prePaid: indicates the yearly/monthly billing mode.
	ChargeMode string `json:"charge_mode" required:"true"`
	// Specifies the unit of the subscription term.
	// This parameter is valid and mandatory only when chargingMode is set to prePaid.
	//   month: indicates that the unit is month.
	//   year: indicates that the unit is year.
	PeriodType string `json:"period_type,omitempty"`
	// Specifies the subscription term. This parameter is valid and mandatory only when chargingMode is set to prePaid.
	//   When periodType is set to month, the parameter value ranges from 1 to 9.
	//   When periodType is set to year, the parameter value ranges from 1 to 3.
	PeriodNum int `json:"period_num,omitempty"`
	// Specifies whether to pay immediately. This parameter is valid only when chargingMode is set to prePaid. The default value is false.
	//   false: indicates not to pay immediately after an order is created.
	//   true: indicates to pay immediately after an order is created. The system will automatically deduct fees from the account balance.
	IsAutoPay *bool `json:"is_auto_pay,omitempty"`
	// Whether auto renew is enabled, default to false.
	IsAutoRenew *bool `json:"is_auto_renew,omitempty"`
}

// NodeGroupOpts is a structure representing node group.
type NodeGroupOpts struct {
	// Instance specifications of a node.
	NodeSize string `json:"node_size" required:"true"`
	// Specifies the node group name.
	// The rules for configuring node groups are as follows:
	//     master_node_default_group: Master node group, which must be included in all cluster types.
	//     core_node_analysis_group: analysis Core node group, which must be contained in the analysis cluster and
	//         hybrid cluster.
	//     core_node_streaming_group: indicates the streaming Core node group, which must be included in both streaming
	//         and hybrid clusters.
	//     task_node_analysis_group: Analysis Task node group.
	//         This node group can be selected for analysis clusters and hybrid clusters as required.
	//     task_node_streaming_group: streaming Task node group.
	//         This node group can be selected for streaming clusters and hybrid clusters as required.
	//     node_group{x}: node group of the customized cluster. A maximum of nine node groups can be added.
	//         The value can contain a maximum of 64 characters, including letters, digits and underscores (_).
	GroupName string `json:"group_name" required:"true"`
	// Number of nodes.
	NodeNum int `json:"node_num" required:"true"`
	// Specifies the system disk information of the node.
	RootVolume *Volume `json:"root_volume,omitempty"`
	// Data disk information.
	DataVolume *Volume `json:"data_volume,omitempty"`
	// Number of data disks of a node. The value range is 0 to 10.
	DataVolumeCount *int `json:"data_volume_count,omitempty"`
	// Billing type of the node group.
	ChargeInfo *ChargeInfo `json:"charge_info,omitempty"`
	// Autoscaling rule corresponding to the node group.
	AsPolicy *AsPolicy `json:"auto_scaling_policy,omitempty"`
	// This parameter is mandatory when the cluster type is CUSTOM. Specifies the roles deployed in a node group.
	// This parameter is a character string array. Each character string represents a role expression.
	// Role expression definition:
	//     If the role is deployed on all nodes in the node group, set this parameter to <role name>, e.g. DataNode.
	//     If the role is deployed on a specified subscript node in the node group:
	//         <role name>:<index1>,<index2>..., <indexN>, e.g. NameNode:1,2.
	//     Some roles support multi-instance deployment (that is, multiple instances of the same role are deployed on a
	//         node): <role name>[<instance count>], for example, EsNode[9].
	AssignedRoles []string `json:"assigned_roles,omitempty"`
}

// Volume is a structure representing node volume configurations.
type Volume struct {
	// Disk Type. The following disk types are supported:
	//     SATA: common I/O disk
	//     SAS: high I/O disk
	//     SSD: ultra-high I/O disk
	Type string `json:"type" required:"true"`
	// Specifies the data disk size, in GB. The value range is 10 to 32768.
	Size int `json:"size" required:"true"`
}

// AsPolicy is a structure representing auto-scaling policy for task nodes.
type AsPolicy struct {
	// Whether to enable the auto scaling rule.
	Enabled string `json:"auto_scaling_enable" required:"true"`
	// Minimum number of nodes left in the node group. The value range is 0 to 500.
	MinCapacity int `json:"min_capacity" required:"true"`
	// Maximum number of nodes in the node group. The value range is 0 to 500.
	MaxCapacity int `json:"max_capacity" required:"true"`
	// List of the resource plan.
	ResourcesPlans []ResourcesPlan `json:"resources_plans,omitempty"`
	// List of custom scaling automation scripts.
	Rules []Rule `json:"rules,omitempty"`
	// List of auto scaling rules.
	ExecScripts []ScaleScript `json:"exec_scripts,omitempty"`
}

// ResourcesPlan is a structure representing resource plan of the policy.
type ResourcesPlan struct {
	// Cycle type of a resource plan.
	PeriodType string `json:"period_type" required:"true"`
	// Start time of a resource plan.
	// The value is in the format of hour:minute, indicating that the time ranges from 0:00 to 23:59.
	StartTime string `json:"start_time" required:"true"`
	// End time of a resource plan. The value is in the same format as that of start_time.
	// The interval between end_time and start_time must be greater than or equal to 30 minutes.
	EndTime string `json:"end_time" required:"true"`
	// Minimum number of the preserved nodes in a node group in a resource plan. The value range is 0 to 500.
	MinCapacity int `json:"min_capacity" required:"true"`
	// Maximum number of the preserved nodes in a node group in a resource plan. The value range is 0 to 500.
	MaxCapacity int `json:"max_capacity" required:"true"`
}

// Rule is a structure representing configuration of the auto-scaling rule.
type Rule struct {
	// Auto scaling rule adjustment type. The options are scale_out and scale_in.
	AdjustmentType string `json:"adjustment_type" required:"true"`
	// Cluster cooling time after an auto scaling rule is triggered, when no auto scaling operation is performed.
	// The unit is minute.
	CoolDownMinutes int `json:"cool_down_minutes" required:"true"`
	// Unique name of an auto scaling rule. A cluster name can contain only 1 to 64 characters.
	// Only letters, digits, hyphens (-), and underscores (_) are allowed.
	Name string `json:"name" required:"true"`
	// Number of nodes that can be adjusted once. The value range is 1 to 100.
	ScalingAdjustment int `json:"scaling_adjustment" required:"true"`
	// Condition for triggering a rule.
	Trigger Trigger `json:"trigger" required:"true"`
	// Description about an auto scaling rule. It contains a maximum of 1,024 characters.
	Description *string `json:"description,omitempty"`
}

// Trigger is a structure representing the condition for the triggering a rule.
type Trigger struct {
	// Number of consecutive five-minute periods, during which a metric threshold is reached.
	// The value range is 1 to 288.
	EvaluationPeriods int `json:"evaluation_periods" required:"true"`
	// Metric name.
	MetricName string `json:"metric_name" required:"true"`
	// Metric threshold to trigger a rule.
	MetricValue string `json:"metric_value" required:"true"`
	// Metric judgment logic operator. The options are LT, GT, LTOE and GTOE.
	ComparisonOperator string `json:"comparison_operator,omitempty"`
}

// ScriptOpts is a structure representing the bootstrap action script information.
type ScriptOpts struct {
	// Whether to continue executing subsequent scripts and creating a cluster after the bootstrap action script fails to be executed.
	// continue: Continue to execute subsequent scripts.
	// errorout: Stop the action.
	// The default value is errorout, indicating that the action is stopped.
	// NOTE:
	// You are advised to set this parameter to continue in the commissioning phase so that the cluster can continue to be installed and started no matter whether the bootstrap action is successful.
	FailAction string `json:"fail_action" required:"true"`
	// Name of a bootstrap action script. It must be unique in a cluster.
	// The value can contain only digits, letters, spaces, hyphens (-), and underscores (_) and must not start with a space.
	// The value can contain 1 to 64 characters.
	Name string `json:"name" required:"true"`
	// Type of a node where the bootstrap action script is executed. The value can be Master, Core, or Task.
	Nodes []string `json:"nodes" required:"true"`
	// Bootstrap action script parameters.
	Parameters string `json:"parameters,omitempty"`
	// Path of a bootstrap action script. Set this parameter to an OBS bucket path or a local VM path.
	// OBS bucket path: Enter a script path manually. For example, enter the path of the public sample script provided by MRS. Example: s3a://bootstrap/presto/presto-install.sh. If dualroles is installed, the parameter of the presto-install.sh script is dualroles. If worker is installed, the parameter of the presto-install.sh script is worker. Based on the Presto usage habit, you are advised to install dualroles on the active Master nodes and worker on the Core nodes.
	// Local VM path: Enter a script path. The script path must start with a slash (/) and end with .sh.
	URI string `json:"uri" required:"true"`
	// Whether the bootstrap action script runs only on active Master nodes.
	// The default value is false, indicating that the bootstrap action script can run on all Master nodes.
	ActiveMaster *bool `json:"active_master,omitempty"`
	// Time when the bootstrap action script is executed. Currently, the following two options are available: Before component start and After component start
	// The default value is false, indicating that the bootstrap action script is executed after the component is started.
	BeforeComponentStart *bool `json:"before_component_start,omitempty"`
	ExecuteNeedSudoRoot  *bool `json:"execute_need_sudo_root,omitempty"`
}

type ComponentConfigOpts struct {
	// The component name of MRS cluster which has installed.
	Name    string       `json:"component_name" required:"true"`
	Configs []ConfigOpts `json:"configs" required:"true"`
}
type ConfigOpts struct {
	// The configuration item key of component installed.
	Key string `json:"key" required:"true"`
	// The configuration item value of component installed.
	Value string `json:"value" required:"true"`
	// The configuration file name of component installed.
	ConfigFileName string `json:"config_file_name" required:"true"`
}

// JobOpts is a structure representing the job which to execution.
type JobOpts struct {
	// SQL program path. This parameter is needed by Spark Script and Hive Script jobs only, and must meet the following requirements:
	// Contains a maximum of 1,023 characters, excluding special characters such as ;|&><'$. The address cannot be empty or full of spaces.
	// Files can be stored in HDFS or OBS. The path varies depending on the file system.
	// OBS: The path must start with s3a://. Files or programs encrypted by KMS are not supported.
	// HDFS: The path starts with a slash (/).
	// Ends with .sql. sql is case-insensitive.
	HiveScriptPath string `json:"hive_script_path" required:"true"`
	// Job name. It contains 1 to 64 characters. Only letters, digits, hyphens (-), and underscores (_) are allowed.
	// NOTE:
	// Identical job names are allowed but not recommended.
	JobName string `json:"job_name" required:"true"`
	// Job type code
	//   1: MapReduce
	//   2: Spark
	//   3: Hive Script
	//   4: HiveQL (not supported currently)
	//   5: DistCp, importing and exporting data (not supported currently)
	//   6: Spark Script
	//   7: Spark SQL, submitting Spark SQL statements (not supported currently).
	// NOTE:
	// Spark and Hive jobs can be added to only clusters that include Spark and Hive components.
	JobType int `json:"job_type" required:"true"`
	// true: Submit a job during cluster creation.
	// false: Submit a job after the cluster is created.
	// Set this parameter to true in this example.
	SubmitJobOnceClusterRun bool `json:"submit_job_once_cluster_run" required:"true"`
	// Key parameter for program execution. The parameter is specified by the function of the user's program. MRS is only responsible for loading the parameter.
	// The parameter contains a maximum of 2,047 characters, excluding special characters such as ;|&>'<$, and can be left blank.
	Arguments string `json:"arguments,omitempty"`
	// Data import and export
	// import
	// export
	FileAction string `json:"file_action,omitempty"`
	// HiveQL statement
	Hql string `json:"hql,omitempty"`
	// Address for inputting data
	// Files can be stored in HDFS or OBS. The path varies depending on the file system.
	//   OBS: The path must start with s3a://. Files or programs encrypted by KMS are not supported.
	//   HDFS: The path starts with a slash (/).
	// The parameter contains a maximum of 1,023 characters, excluding special characters such as ;|&>'<$, and can be left blank.
	Input string `json:"input,omitempty"`
	// Path of the JAR or SQL file for program execution. The parameter must meet the following requirements:
	// Contains a maximum of 1,023 characters, excluding special characters such as ;|&><'$. The parameter value cannot be empty or full of spaces.
	// Files can be stored in HDFS or OBS. The path varies depending on the file system.
	//   OBS: The path must start with s3a://. Files or programs encrypted by KMS are not supported.
	//   HDFS: The path starts with a slash (/).
	// Spark Script must end with .sql while MapReduce and Spark Jar must end with .jar. sql and jar are case-insensitive.
	JarPath string `json:"jar_path,omitempty"`
	// Path for storing job logs that record job running status.
	// Files can be stored in HDFS or OBS. The path varies depending on the file system.
	//   OBS: The path must start with s3a://.
	//   HDFS: The path starts with a slash (/).
	// The parameter contains a maximum of 1,023 characters, excluding special characters such as ;|&>'<$, and can be left blank.
	JobLog string `json:"job_log,omitempty"`
	// Address for outputting data
	// Files can be stored in HDFS or OBS. The path varies depending on the file system.
	//   OBS: The path must start with s3a://.
	//   HDFS: The path starts with a slash (/).
	// If the specified path does not exist, the system will automatically create it.
	// The parameter contains a maximum of 1,023 characters, excluding special characters such as ;|&>'<$, and can be left blank.
	Output string `json:"output,omitempty"`
	// Whether to delete the cluster after the job execution is complete.
	ShutdownCluster bool `json:"shutdown_cluster,omitempty"`
}

// ScaleScript is a structure representing the auto-scaling rules.
type ScaleScript struct {
	// Unique name of a custom automation script. The value can contain 1 to 64 characters.
	// The value can contain digits, letters, spaces, hyphens (-), and underscores (_) and must not start with a space.
	Name string `json:"name" required:"true"`
	// Path of a custom automation script. Set this parameter to an OBS bucket path or a local VM path.
	//     OBS bucket path: Enter a script path manually. for example, s3a://XXX/scale.sh.
	//     Local VM path: Enter a script path. The script path must start with a slash (/) and end with .sh.
	URI string `json:"uri" required:"true"`
	// Type of a node where the custom automation script is executed. The node type can be Master, Core, or Task.
	Nodes []string `json:"nodes" required:"true"`
	// Time when a script is executed. The following four options are supported:
	//     before_scale_out: before scale-out
	//     before_scale_in: before scale-in
	//     after_scale_out: after scale-out
	//     after_scale_in: after scale-in
	ActionStage string `json:"action_stage" required:"true"`
	// Whether to continue to execute subsequent scripts and create a cluster after the custom automation script fails
	// to be executed.
	//     continue: Continue to execute subsequent scripts.
	//     errorout: Stop the action.
	FailAction string `json:"fail_action" required:"true"`
	// Parameters of a custom automation script. Multiple parameters are separated by space.
	// The following predefined system parameters can be transferred:
	//     ${mrs_scale_node_num}: Number of the nodes to be added or removed.
	//     ${mrs_scale_type}: Scaling type. The value can be scale_out or scale_in.
	//     ${mrs_scale_node_hostnames}: Host names of the nodes to be added or removed.
	//     ${mrs_scale_node_ips}: IP addresses of the nodes to be added or removed.
	//     ${mrs_scale_rule_name}: Name of the rule that triggers auto scaling.
	// Other user-defined parameters are used in the same way as those of common shell scripts.
	Parameters string `json:"parameters,omitempty"`
	// Whether the custom automation script runs only on the active Master node.
	ActiveMaster bool `json:"active_master,omitempty"`
}

// CreateOptsBuilder is an interface which to support request body build of the cluster creation.
type CreateOptsBuilder interface {
	ToClusterCreateMap() (map[string]interface{}, error)
}

// ToClusterCreateMap is a method which to build a request body by the CreateOpts.
func (opts CreateOpts) ToClusterCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is a method to create a new mapreduce cluster.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToClusterCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = client.Post(rootURL(client), b, &r.Body, reqOpt)
	return
}
