package mrs

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/mrs/v1/cluster"
	"github.com/huaweicloud/golangsdk/openstack/mrs/v2/clusters"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/vpcs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const (
	typeAnalysis = "ANALYSIS"
	typeStream   = "STREAMING"
	typeHybrid   = "MIXED"

	masterGroup        = "master_node_default_group"
	analysisCoreGroup  = "core_node_analysis_group"
	streamingCoreGroup = "core_node_streaming_group"
	analysisTaskGroup  = "task_node_analysis_group"
	streamingTaskGroup = "task_node_streaming_group"
)

type stateRefresh struct {
	Pending      []string
	Target       []string
	Timeout      time.Duration
	Delay        time.Duration
	PollInterval time.Duration
}

func ResourceMRSClusterV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceMRSClusterV2Create,
		Read:   resourceMRSClusterV2Read,
		Update: resourceMRSClusterV2Update,
		Delete: resourceMRSClusterV2Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[A-Za-z][A-Za-z0-9_-]{1,63}$"),
					"The name consists of 2 to 64 characters, starting with a letter. "+
						"Only letters, digits, hyphens (-) and underscores (_) are allowed."),
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"component_list": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"manager_admin_pass": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  typeAnalysis,
				ValidateFunc: validation.StringInSlice([]string{
					typeAnalysis, typeStream, typeHybrid,
				}, false),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"log_collection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"node_admin_pass": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ForceNew:     true,
				ExactlyOneOf: []string{"node_key_pair"},
			},
			"node_key_pair": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"safe_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"master_nodes": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     masterNodeSchemaResource(1, 2),
			},
			"analysis_core_nodes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     coreTaskNodeSchemaResource(1, 500),
			},
			"streaming_core_nodes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     coreTaskNodeSchemaResource(1, 500),
			},
			"analysis_task_nodes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     coreTaskNodeSchemaResource(1, 500),
			},
			"streaming_task_nodes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     coreTaskNodeSchemaResource(1, 500),
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"total_node_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"master_node_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internal_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_alternate_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vnc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charging_start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func masterNodeSchemaResource(min, max int) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_number": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(min, max),
			},
			"root_volume_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"root_volume_size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"data_volume_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_volume_size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"data_volume_count": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func coreTaskNodeSchemaResource(min, max int) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_number": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(min, max),
			},
			"root_volume_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"root_volume_size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"data_volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"data_volume_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"data_volume_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

// The 'component_list' type of the request body is string, before body build, it should be conversation from set to string.
func buildMrsComponents(d *schema.ResourceData) string {
	components := d.Get("component_list").(*schema.Set)
	return buildStringBySet(components)
}

// The 'security_group_ids' type of the request body is string, before body build, it should be conversation from set to string.
func buildMrsSecurityGroupIds(d *schema.ResourceData) string {
	secgroupIds := d.Get("security_group_ids").(*schema.Set)
	return buildStringBySet(secgroupIds)
}

// The 'log_collection' type of the request body is int,
func buildLogCollection(d *schema.ResourceData) *int {
	if d.Get("log_collection").(bool) {
		return golangsdk.IntToPointer(1)
	}
	return golangsdk.IntToPointer(0)
}

func buildMrsSafeMode(d *schema.ResourceData) string {
	isSafe := d.Get("safe_mode").(bool)
	if isSafe {
		return "KERBEROS"
	}
	return "SIMPLE"
}

func buildStringBySet(set *schema.Set) string {
	slice := make([]string, set.Len())
	for i, v := range set.List() {
		slice[i] = v.(string)
	}
	return strings.Join(slice, ",")
}

// buildMrsClusterNodeGroups is a method which to build a node group list with all node group arguments.
func buildMrsClusterNodeGroups(d *schema.ResourceData) []clusters.NodeGroupOpts {
	var (
		groupOpts      []clusters.NodeGroupOpts
		nodeGroupTypes = map[string]string{
			"master_nodes":         masterGroup,
			"analysis_core_nodes":  analysisCoreGroup,
			"analysis_task_nodes":  analysisTaskGroup,
			"streaming_core_nodes": streamingCoreGroup,
			"streaming_task_nodes": streamingTaskGroup,
		}
	)
	for k, v := range nodeGroupTypes {
		if optsRaw, ok := d.GetOk(k); ok {
			opts := buildNodeGroupOpts(d, optsRaw.([]interface{}), v)
			groupOpts = append(groupOpts, opts)
		}
	}
	return groupOpts
}

func buildNodeGroupOpts(d *schema.ResourceData, optsRaw []interface{}, name string) clusters.NodeGroupOpts {
	var result = clusters.NodeGroupOpts{}
	if len(optsRaw) > 0 {
		opts := optsRaw[0].(map[string]interface{})
		result.GroupName = name
		result.NodeSize = opts["flavor"].(string)
		result.NodeNum = opts["node_number"].(int)
		result.RootVolume = &clusters.Volume{
			Type: opts["root_volume_type"].(string),
			Size: opts["root_volume_size"].(int),
		}
		volumeCount := opts["data_volume_count"].(int)
		if volumeCount != 0 {
			result.DataVolume = &clusters.Volume{
				Type: opts["data_volume_type"].(string),
				Size: opts["data_volume_size"].(int),
			}
		} else {
			// According to the API rules, when the data disk is not mounted, the parameters in the structure still
			// need to be filled in (but not used), fill in the system disk data here.
			result.DataVolume = &clusters.Volume{
				Type: opts["root_volume_type"].(string),
				Size: opts["root_volume_size"].(int),
			}
		}
		result.DataVolumeCount = golangsdk.IntToPointer(volumeCount)
	}
	return result
}

func getVpcNameById(d *schema.ResourceData, config *config.Config, id string) (string, error) {
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return "", fmtp.Errorf("Error creating Huaweicloud Vpc client: %s", err)
	}
	vpc, err := vpcs.Get(client, id).Extract()
	if err != nil {
		return "", fmtp.Errorf("Error retrieving Huaweicloud Vpc: %s", err)
	}
	return vpc.Name, nil
}

func clusterV2StateRefreshFunc(client *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		clusterGet, err := cluster.Get(client, clusterId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return clusterGet, "DELETED", nil
			}
			return nil, "", err
		}
		return clusterGet, clusterGet.Clusterstate, nil
	}
}

func waitForMrsClusterStateCompleted(client *golangsdk.ServiceClient, id string, refresh stateRefresh) error {
	stateConf := &resource.StateChangeConf{
		Pending:      refresh.Pending,
		Target:       refresh.Target,
		Refresh:      clusterV2StateRefreshFunc(client, id),
		Timeout:      refresh.Timeout,
		Delay:        refresh.Delay,
		PollInterval: refresh.PollInterval,
	}
	_, err := stateConf.WaitForState()
	if err != nil {
		//the system will recyle the cluster when creating failed
		return err
	}
	return nil
}

// addTagsToMrsCluster method is inherited from MRS V1 resources.
func addTagsToMrsCluster(d *schema.ResourceData, config *config.Config) error {
	client, err := config.MrsV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud MRS V1 client: %s", err)
	}
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(client, "clusters", d.Id(), taglist).ExtractErr(); tagErr != nil {
			return fmtp.Errorf("Error setting tags of MRS cluster %s: %s", d.Id(), tagErr)
		}
	}
	return nil
}

func resourceMRSClusterV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)

	vpcId := d.Get("vpc_id").(string)
	vpcName, err := getVpcNameById(d, config, vpcId)
	if err != nil {
		return fmtp.Errorf("Unable to find the vpc (%s) on the server: %s", vpcId, err)
	}

	createOpts := &clusters.CreateOpts{
		Region:               config.GetRegion(d),
		AvailabilityZone:     d.Get("availability_zone").(string),
		ClusterVersion:       d.Get("version").(string),
		ClusterName:          d.Get("name").(string),
		ClusterType:          d.Get("type").(string),
		ManagerAdminPassword: d.Get("manager_admin_pass").(string),
		VpcName:              vpcName,
		SubnetId:             d.Get("subnet_id").(string),
		EipAddress:           d.Get("ip_address").(string),
		Components:           buildMrsComponents(d),
		EnterpriseProjectId:  common.GetEnterpriseProjectID(d, config),
		LogCollection:        buildLogCollection(d),
		NodeGroups:           buildMrsClusterNodeGroups(d),
		SafeMode:             buildMrsSafeMode(d),
		SecurityGroupsIds:    buildMrsSecurityGroupIds(d),
	}
	if v, ok := d.GetOk("node_key_pair"); ok {
		createOpts.NodeKeypair = v.(string)
		createOpts.LoginMode = "KEYPAIR"
	} else {
		createOpts.NodeRootPassword = d.Get("node_admin_pass").(string)
		createOpts.LoginMode = "PASSWORD"
	}

	mrsV2Client, err := config.MrsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud MRS V2 client: %s", err)
	}
	resp, err := clusters.Create(mrsV2Client, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating Cluster: %s", err)
	}
	d.SetId(resp.ID)
	refresh := stateRefresh{
		Pending:      []string{"starting"},
		Target:       []string{"running"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        480 * time.Second,
		PollInterval: 15 * time.Second,
	}
	mrsV1Client, err := config.MrsV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud MRS V1 client: %s", err)
	}
	// After request send, check the cluster state and wait for it become running.
	if err = waitForMrsClusterStateCompleted(mrsV1Client, d.Id(), refresh); err != nil {
		d.SetId("")
		return fmtp.Errorf("Error waiting for MapReduce cluster (%s) to become ready: %s", d.Id(), err)
	}
	// After MapReduce cluster state become running, add some tags to the cluster.
	if err = addTagsToMrsCluster(d, config); err != nil {
		return fmtp.Errorf("Error waiting for MapReduce cluster (%s) to become ready: %s", d.Id(), err)
	}

	return resourceMRSClusterV2Read(d, meta)
}

func setMrsClsuterType(d *schema.ResourceData, resp *cluster.Cluster) error {
	// The returned ClusterType is an 'Int' type, with a value of 0 to 2,
	// which respectively represent:'ANALYSIS','STREAMING' and'MIXED'.
	clusterType := []string{"ANALYSIS", "STREAMING", "MIXED"}
	if resp.ClusterType >= len(clusterType) || resp.ClusterType < 0 {
		return fmtp.Errorf("The cluster type of the response is '%d', not in the map", resp.ClusterType)
	}
	return d.Set("type", clusterType[resp.ClusterType])
}

func setMrsClsuterSafeMode(d *schema.ResourceData, resp *cluster.Cluster) error {
	result := true
	if resp.Safemode == 0 {
		result = false
	}
	return d.Set("safe_mode", result)
}

func setMRSClusterLogCollection(d *schema.ResourceData, resp *cluster.Cluster) error {
	result := true
	if resp.LogCollection == 0 {
		result = false
	}
	return d.Set("log_collection", result)
}

func setMrsClsuterSecurityGroupIds(d *schema.ResourceData, resp *cluster.Cluster) error {
	secGroupsIds := strings.Split(resp.Securitygroupsid, ",")
	return d.Set("security_group_ids", secGroupsIds)
}

func setMrsClsuterTotalNodeNumber(d *schema.ResourceData, resp *cluster.Cluster) error {
	totalNodes, err := strconv.Atoi(resp.Totalnodenum)
	if err != nil {
		return err
	}
	return d.Set("total_node_number", totalNodes)
}

func setMrsClsuterCreateTimestamp(d *schema.ResourceData, resp *cluster.Cluster) error {
	createTime, err := strconv.ParseInt(resp.Createat, 10, 64)
	if err != nil {
		return err
	}
	return d.Set("create_time", utils.FormatTimeStampRFC3339(createTime))
}

func setMrsClsuterUpdateTimestamp(d *schema.ResourceData, resp *cluster.Cluster) error {
	updateTime, err := strconv.ParseInt(resp.Updateat, 10, 64)
	if err != nil {
		return err
	}
	return d.Set("update_time", utils.FormatTimeStampRFC3339(updateTime))
}

func setMrsClsuterChargingTimestamp(d *schema.ResourceData, resp *cluster.Cluster) error {
	chargingStartTime, err := strconv.ParseInt(resp.Chargingstarttime, 10, 64)
	if err != nil {
		return err
	}
	return d.Set("charging_start_time", utils.FormatTimeStampRFC3339(chargingStartTime))
}

func setMrsClsuterComponentList(d *schema.ResourceData, resp *cluster.Cluster) error {
	result := make([]interface{}, len(resp.Componentlist))
	for i, attachment := range resp.Componentlist {
		result[i] = attachment.Componentname
	}
	return d.Set("component_list", result)
}

func setMrsClusterNodeGroups(d *schema.ResourceData, resp *cluster.Cluster) error {
	var groupMap = map[string]string{
		masterGroup:        "master_nodes",
		analysisCoreGroup:  "analysis_core_nodes",
		streamingCoreGroup: "streaming_core_nodes",
		analysisTaskGroup:  "analysis_task_nodes",
		streamingTaskGroup: "streaming_task_nodes",
	}
	for _, node := range resp.NodeGroups {
		value, ok := groupMap[node.GroupName]
		if !ok {
			logp.Printf("[DEBUG] %s is not in the resource data", node.GroupName)
		}
		groupMap := map[string]interface{}{
			"node_number":      node.NodeNum,
			"flavor":           node.NodeSize,
			"root_volume_type": node.RootVolumeType,
			"root_volume_size": node.RootVolumeSize,
		}
		if node.DataVolumeCount != 0 {
			groupMap["data_volume_type"] = node.DataVolumeType
			groupMap["data_volume_size"] = node.DataVolumeSize
			groupMap["data_volume_count"] = node.DataVolumeCount
		}
		logp.Printf("[DEBUG] node group '%s' is : %+v", value, groupMap)
		// lintignore:R001
		if err := d.Set(value, []map[string]interface{}{groupMap}); err != nil {
			return err
		}
	}

	return nil
}

func setClsuterTags(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	resourceTags, err := tags.Get(client, "clusters", d.Id()).Extract()
	if err != nil {
		return fmtp.Errorf("Error Fetching tags of MapReduce cluster form server: %s", err)
	}
	tagmap := utils.TagsToMap(resourceTags.Tags)
	return d.Set("tags", tagmap)
}

func getMrsClusterFromServer(d *schema.ResourceData, client *golangsdk.ServiceClient) (*cluster.Cluster, error) {
	resp, err := cluster.Get(client, d.Id()).Extract()
	if err != nil {
		return nil, common.CheckDeleted(d, err, "MapReduce cluster is not exist")
	}
	if resp.Clusterstate == "terminated" {
		d.SetId("")
		return resp, fmtp.Errorf("Retrieved Cluster %s, but it was terminated, abort it", d.Id())
	}
	return resp, nil
}

func resourceMRSClusterV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.MrsV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud MRS client: %s", err)
	}
	resp, err := getMrsClusterFromServer(d, client)
	if err != nil {
		return fmtp.Errorf("Error getting MapReduce form server: %s", err)
	}

	logp.Printf("[DEBUG] Retrieved Cluster %s: %#v", d.Id(), resp)
	d.SetId(resp.Clusterid)
	mErr := multierror.Append(
		d.Set("region", resp.Datacenter),
		d.Set("availability_zone", resp.AvailabilityZone),
		d.Set("name", resp.Clustername),
		d.Set("version", resp.Clusterversion),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("node_key_pair", resp.Nodepubliccertname),
		d.Set("vpc_id", resp.Vpcid),
		d.Set("subnet_id", resp.Subnetid),
		d.Set("ip_address", resp.Internalip),
		d.Set("master_node_ip", resp.Masternodeip),
		d.Set("external_ip", resp.Externalip),
		d.Set("private_ip", resp.Privateipfirst),
		d.Set("internal_ip", resp.Internalip),
		d.Set("external_alternate_ip", resp.Externalalternateip),
		d.Set("vnc", resp.Vnc),
		d.Set("status", resp.Clusterstate),
		setMrsClsuterType(d, resp),
		setMrsClsuterComponentList(d, resp),
		setMrsClsuterSafeMode(d, resp),
		setMRSClusterLogCollection(d, resp),
		setMrsClsuterSecurityGroupIds(d, resp),
		setMrsClsuterTotalNodeNumber(d, resp),
		setMrsClsuterCreateTimestamp(d, resp),
		setMrsClsuterUpdateTimestamp(d, resp),
		setMrsClsuterChargingTimestamp(d, resp),
		setMrsClsuterCreateTimestamp(d, resp),
		setMrsClusterNodeGroups(d, resp),
		setClsuterTags(d, client),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.Errorf("Error setting vault fields: %s", err)
	}

	return nil
}

// resizeMRSClusterCoreNodes is a method which used to resize core node for each cluster type.
// The resizeCount is a number of the group size changing, nagetive means scale in group.
func resizeMRSClusterCoreNodes(client *golangsdk.ServiceClient, id, groupType string, resizeCount int) error {
	var isScaleOut = "scale_out"
	if resizeCount < 0 {
		isScaleOut = "scale_in"
		resizeCount = -resizeCount
	}
	opts := cluster.UpdateOpts{
		Parameters: cluster.ResizeParameters{
			ScaleType: isScaleOut,
			NodeGroup: &groupType,
			NodeId:    "node_orderadd",
			Instances: strconv.Itoa(resizeCount),
		},
	}
	_, err := cluster.Update(client, id, opts).Extract()
	if err != nil {
		return fmtp.Errorf("Error resizing core node")
	}
	refresh := stateRefresh{
		Pending:      []string{"scaling-out", "scaling-in"},
		Target:       []string{"running"},
		Delay:        2 * time.Minute,
		Timeout:      1 * time.Hour,
		PollInterval: 15 * time.Second,
	}
	if err = waitForMrsClusterStateCompleted(client, id, refresh); err != nil {
		return fmtp.Errorf("Error waiting for Mrs cluster resize to be complated: %s", err)
	}
	return nil
}

// resizeMRSClusterTaskNodes is a method which use to scale out/in the node groups.
func resizeMRSClusterTaskNodes(client *golangsdk.ServiceClient, id, groupType string, oldList, newList []interface{},
	resizeCount int) error {
	var isScaleOut = "scale_out"
	newRaw := newList[0].(map[string]interface{})

	if resizeCount < 0 {
		isScaleOut = "scale_in"
		resizeCount = -resizeCount
	}
	params := cluster.ResizeParameters{
		ScaleType: isScaleOut,
		NodeGroup: &groupType,
		NodeId:    "node_orderadd", // Fixed value of resize request
		Instances: strconv.Itoa(resizeCount),
	}
	if len(oldList) == 0 {
		params.TaskNodeInfo = &cluster.TaskNodeInfo{
			NodeSize:        newRaw["flavor"].(string),
			DataVolumeType:  newRaw["data_volume_type"].(string),
			DataVolumeSize:  newRaw["data_volume_size"].(int),
			DataVolumeCount: newRaw["data_volume_count"].(int),
		}
	}
	opts := cluster.UpdateOpts{
		Parameters: params,
	}
	_, err := cluster.Update(client, id, opts).Extract()
	if err != nil {
		return fmtp.Errorf("Error resizing task node")
	}
	refresh := stateRefresh{
		Pending:      []string{"scaling-out", "scaling-in"},
		Target:       []string{"running"},
		Delay:        2 * time.Minute,
		Timeout:      1 * time.Hour,
		PollInterval: 15 * time.Second,
	}
	if err = waitForMrsClusterStateCompleted(client, id, refresh); err != nil {
		return fmtp.Errorf("Error waiting for Mrs cluster resize to be complated: %s", err)
	}
	return nil
}

// the getNodeResizeNumber is a method which use to calculate the number of the group resize option.
func getNodeResizeNumber(oldList, newList []interface{}) int {
	newNode := newList[0].(map[string]interface{})
	newSize := newNode["node_number"].(int)
	if len(oldList) == 0 {
		return newSize
	}
	oldNode := oldList[0].(map[string]interface{})
	oldSize := oldNode["node_number"].(int)
	// Distinguish scale out and scale in by positive and negative
	return newSize - oldSize
}

func updateMRSClusterNodes(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	clusterType := d.Get("type").(string)
	if clusterType == typeAnalysis || clusterType == typeHybrid {
		if d.HasChange("analysis_core_nodes") {
			oldRaws, newRaws := d.GetChange("analysis_core_nodes")
			num := getNodeResizeNumber(oldRaws.([]interface{}), newRaws.([]interface{}))
			err := resizeMRSClusterCoreNodes(client, d.Id(), analysisCoreGroup, num)
			if err != nil {
				return err
			}
		}
		if d.HasChange("analysis_task_nodes") {
			oldRaws, newRaws := d.GetChange("analysis_task_nodes")
			num := getNodeResizeNumber(oldRaws.([]interface{}), newRaws.([]interface{}))
			err := resizeMRSClusterTaskNodes(client, d.Id(), analysisTaskGroup,
				oldRaws.([]interface{}), newRaws.([]interface{}), num)
			if err != nil {
				return err
			}
		}
	}
	if clusterType == typeStream || clusterType == typeHybrid {
		if d.HasChange("streaming_core_nodes") {
			oldRaws, newRaws := d.GetChange("streaming_core_nodes")
			num := getNodeResizeNumber(oldRaws.([]interface{}), newRaws.([]interface{}))
			err := resizeMRSClusterCoreNodes(client, d.Id(), streamingCoreGroup, num)
			if err != nil {
				return err
			}
		}
		if d.HasChange("streaming_task_nodes") {
			oldRaws, newRaws := d.GetChange("streaming_task_nodes")
			num := getNodeResizeNumber(oldRaws.([]interface{}), newRaws.([]interface{}))
			err := resizeMRSClusterTaskNodes(client, d.Id(), streamingTaskGroup,
				oldRaws.([]interface{}), newRaws.([]interface{}), num)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceMRSClusterV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.MrsV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud MRS client: %s", err)
	}
	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(client, d, "clusters", d.Id())
		if tagErr != nil {
			return fmtp.Errorf("Error updating tags of MRS cluster:%s, err:%s", d.Id(), tagErr)
		}
	}
	// lintignore:R005
	if d.HasChange("analysis_core_nodes") || d.HasChange("streaming_core_nodes") ||
		d.HasChange("analysis_task_nodes") || d.HasChange("streaming_task_nodes") {
		err = updateMRSClusterNodes(d, client)
		if err != nil {
			return err
		}
	}

	return resourceMRSClusterV2Read(d, meta)
}

func resourceMRSClusterV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.MrsV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud MRS client: %s", err)
	}

	err = cluster.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud Cluster: %s", err)
	}
	refresh := stateRefresh{
		Pending:      []string{"running", "terminating"},
		Target:       []string{"terminated", "DELETED"},
		Delay:        45 * time.Second,
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 10 * time.Second,
	}
	if err = waitForMrsClusterStateCompleted(client, d.Id(), refresh); err != nil {
		d.SetId("")
		return fmtp.Errorf("Error waiting for Mrs cluster (%s) to be terminated: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
