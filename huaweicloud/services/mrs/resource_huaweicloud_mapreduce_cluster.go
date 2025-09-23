package mrs

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/mrs/v1/cluster"
	clusterv2 "github.com/chnsz/golangsdk/openstack/mrs/v2/clusters"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	typeAnalysis = "ANALYSIS"
	typeStream   = "STREAMING"
	typeHybrid   = "MIXED"
	typeCustom   = "CUSTOM"

	masterGroup        = "master_node_default_group"
	analysisCoreGroup  = "core_node_analysis_group"
	streamingCoreGroup = "core_node_streaming_group"
	analysisTaskGroup  = "task_node_analysis_group"
	streamingTaskGroup = "task_node_streaming_group"
	customNodeGroup    = "Core"

	mrsHostDefaultPageNum  = 1
	mrsHostDefaultPageSize = 100
)

type stateRefresh struct {
	Pending      []string
	Target       []string
	Timeout      time.Duration
	Delay        time.Duration
	PollInterval time.Duration
}

// @API MRS GET /v1.1/{project_id}/clusters/{cluster_id}/hosts
// @API MRS DELETE /v1.1/{project_id}/clusters/{cluster_id}
// @API MRS POST /v2/{project_id}/clusters
// @API MRS PUT /v2/{project_id}/clusters/{cluster_id}/cluster-name
// @API MRS POST /v1.1/{project_id}/{resourceType}/{id}/tags/action
// @API MRS GET /v1.1/{project_id}/cluster_infos/{cluster_id}
// @API MRS PUT /v1.1/{project_id}/cluster_infos/{cluster_id}
// @API EIP GET /v1/{project_id}/publicips
// @API VPC GET /v1/{project_id}/vpcs
// @API VPC GET /v1/{project_id}/subnets
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrate
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources/filter
func ResourceMRSClusterV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMRSClusterV2Create,
		ReadContext:   resourceMRSClusterV2Read,
		UpdateContext: resourceMRSClusterV2Update,
		DeleteContext: resourceMRSClusterV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

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
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"eip_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
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
			"template_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"master_nodes": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     nodeGroupSchemaResource("master_node_default_group", false),
			},
			"analysis_core_nodes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     nodeGroupSchemaResource("core_node_analysis_group", true),
			},
			"streaming_core_nodes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     nodeGroupSchemaResource("core_node_streaming_group", true),
			},
			"analysis_task_nodes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     nodeGroupSchemaResource("task_node_analysis_group", true),
			},
			"streaming_task_nodes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     nodeGroupSchemaResource("task_node_streaming_group", true),
			},
			"custom_nodes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     nodeGroupSchemaResource("", false),
			},
			"component_configs": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     componentConfigsSchemaResource(),
			},
			"tags": common.TagsSchema(),
			"external_datasources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     externalDatasourceSchema(),
				ForceNew: true,
			},
			"bootstrap_scripts": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     bootstrapScriptsSchema(),
				ForceNew: true,
			},
			"smn_notify": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     smnNotifySchema(),
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
			// Attributes.
			"total_node_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"master_node_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
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

/*
when custom node,the groupName should been empty
*/
func nodeGroupSchemaResource(groupName string, nodeScalable bool) *schema.Resource {
	nodeResource := schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Required: true,
				ForceNew: true,
			},
			"assigned_roles": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if clusterType := d.Get("type").(string); clusterType != typeCustom {
						return true
					}
					return false
				},
			},
			"host_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}

	// when custom type,should support user specified name
	if groupName == "" {
		nodeResource.Schema["group_name"] = &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		}
	}

	if nodeScalable {
		nodeResource.Schema["node_number"] = &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		}
	} else {
		nodeResource.Schema["node_number"] = &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
		}
	}

	// All task node groups do not support `prepaid`.
	if groupName != analysisTaskGroup && groupName != streamingTaskGroup {
		nodeResource.Schema["charging_mode"] = common.SchemaChargingMode(nil)
		nodeResource.Schema["period_unit"] = &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"month", "year",
			}, false),
		}
		nodeResource.Schema["period"] = &schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 9),
		}
		nodeResource.Schema["auto_renew"] = common.SchemaAutoRenewUpdatable(nil)
	}

	return &nodeResource
}

func componentConfigsSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"configs": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"config_file_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func externalDatasourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"component_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_connection_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func bootstrapScriptsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Name of a bootstrap action script.`,
			},
			"uri": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Path of a bootstrap action script. Set this parameter to an OBS bucket path or a local VM path.`,
			},
			"nodes": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				ForceNew:    true,
				Description: `Name of the node group where the bootstrap action script is executed.`,
			},
			"fail_action": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The action after the bootstrap action script fails to be executed.`,
			},
			"parameters": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Bootstrap action script parameters.`,
			},
			"active_master": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Whether the bootstrap action script runs only on active master nodes.`,
			},
			"before_component_start": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Whether the bootstrap action script is executed before component start.`,
			},
			"execute_need_sudo_root": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Whether the bootstrap action script involves root user operations.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The execution time of one bootstrap action script, in RFC-3339 format.`,
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of one bootstrap action script.`,
			},
		},
	}
}

func smnNotifySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"topic_urn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The Uniform Resource Name (URN) of the topic.`,
			},
			"subscription_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The subscription rule name.`,
			},
		},
	}
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

// buildMrsClusterNodeGroups is a method which to build a node group list with all node group arguments.
func buildMrsClusterNodeGroups(d *schema.ResourceData) []clusterv2.NodeGroupOpts {
	var (
		groupOpts      []clusterv2.NodeGroupOpts
		nodeGroupTypes = map[string]string{
			"master_nodes":         masterGroup,
			"analysis_core_nodes":  analysisCoreGroup,
			"analysis_task_nodes":  analysisTaskGroup,
			"streaming_core_nodes": streamingCoreGroup,
			"streaming_task_nodes": streamingTaskGroup,
			"custom_nodes":         customNodeGroup,
		}
	)
	for k, v := range nodeGroupTypes {
		if optsRaw, ok := d.GetOk(k); ok {
			opts := buildNodeGroupOpts(d, optsRaw.([]interface{}), v)
			groupOpts = append(groupOpts, opts...)
		}
	}
	return groupOpts
}

func buildNodeGroupOpts(d *schema.ResourceData, optsRaw []interface{}, defaultName string) []clusterv2.NodeGroupOpts {
	var result []clusterv2.NodeGroupOpts
	for i := 0; i < len(optsRaw); i++ {
		var nodeGroup = clusterv2.NodeGroupOpts{}
		opts := optsRaw[i].(map[string]interface{})

		nodeGroup.GroupName = defaultName
		customName := opts["group_name"]
		if customName != nil {
			nodeGroup.GroupName = customName.(string)
		}

		nodeGroup.NodeSize = opts["flavor"].(string)
		nodeGroup.NodeNum = opts["node_number"].(int)
		nodeGroup.RootVolume = &clusterv2.Volume{
			Type: opts["root_volume_type"].(string),
			Size: opts["root_volume_size"].(int),
		}
		volumeCount := opts["data_volume_count"].(int)
		if volumeCount != 0 {
			nodeGroup.DataVolume = &clusterv2.Volume{
				Type: opts["data_volume_type"].(string),
				Size: opts["data_volume_size"].(int),
			}
		} else {
			// According to the API rules, when the data disk is not mounted, the parameters in the structure still
			// need to be filled in (but not used), fill in the system disk data here.
			nodeGroup.DataVolume = &clusterv2.Volume{
				Type: opts["root_volume_type"].(string),
				Size: opts["root_volume_size"].(int),
			}
		}

		nodeType := nodeGroup.GroupName
		// All task node groups do not support `prepaid`.
		if nodeType != "task_node_analysis_group" && nodeType != "task_node_streaming_group" {
			nodeGroup.ChargeInfo = bulidNodeGroupChargeInfo(opts, d)
		}

		nodeGroup.DataVolumeCount = golangsdk.IntToPointer(volumeCount)
		// This parameter is mandatory when the cluster type is CUSTOM. Specifies the roles deployed in a node group.
		if clusterType := d.Get("type").(string); clusterType == typeCustom {
			for _, v := range opts["assigned_roles"].([]interface{}) {
				nodeGroup.AssignedRoles = append(nodeGroup.AssignedRoles, v.(string))
			}
		}

		result = append(result, nodeGroup)
	}
	return result
}

func bulidClusterChargeMode(d *schema.ResourceData) *clusterv2.ChargeInfo {
	autoRenew, _ := strconv.ParseBool(d.Get("auto_renew").(string))
	return &clusterv2.ChargeInfo{
		ChargeMode:  "prePaid",
		PeriodNum:   d.Get("period").(int),
		PeriodType:  d.Get("period_unit").(string),
		IsAutoRenew: utils.Bool(autoRenew),
		IsAutoPay:   utils.Bool(true),
	}
}

func bulidNodeGroupChargeInfo(nodeGroup map[string]interface{}, d *schema.ResourceData) *clusterv2.ChargeInfo {
	if nodeGroup["charging_mode"].(string) != "prePaid" && d.Get("charging_mode").(string) != "prePaid" {
		return nil
	}

	// If mode group does not specify charge information, it will inherit the cluster.
	nodeChargeInfo := bulidClusterChargeMode(d)
	if period, ok := nodeGroup["period"].(int); ok && period != 0 {
		nodeChargeInfo.PeriodNum = period
	}

	if periodUnit, ok := nodeGroup["period_unit"].(string); ok && periodUnit != "" {
		nodeChargeInfo.PeriodType = periodUnit
	}

	if autoRenew, ok := nodeGroup["auto_renew"].(string); ok && autoRenew != "" {
		isAutoRenew, _ := strconv.ParseBool(autoRenew)
		nodeChargeInfo.IsAutoRenew = utils.Bool(isAutoRenew)
	}

	return nodeChargeInfo
}

func buildComponentConfigOpts(d *schema.ResourceData) []clusterv2.ComponentConfigOpts {
	v, ok := d.GetOk("component_configs")
	if !ok {
		return nil
	}

	optsRaw := v.([]interface{})
	var result = make([]clusterv2.ComponentConfigOpts, len(optsRaw))
	for i, v := range optsRaw {
		opts := v.(map[string]interface{})
		configOptsRaw := opts["configs"].([]interface{})

		var configOpts = make([]clusterv2.ConfigOpts, len(configOptsRaw))
		for j, item := range configOptsRaw {
			opt := item.(map[string]interface{})
			configOpts[j] = clusterv2.ConfigOpts{
				Key:            opt["key"].(string),
				Value:          opt["value"].(string),
				ConfigFileName: opt["config_file_name"].(string),
			}
		}

		result[i] = clusterv2.ComponentConfigOpts{
			Name:    opts["name"].(string),
			Configs: configOpts,
		}
	}

	return result
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

func waitForMrsClusterStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, id string, refresh stateRefresh) error {
	stateConf := &resource.StateChangeConf{
		Pending:      refresh.Pending,
		Target:       refresh.Target,
		Refresh:      clusterV2StateRefreshFunc(client, id),
		Timeout:      refresh.Timeout,
		Delay:        refresh.Delay,
		PollInterval: refresh.PollInterval,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		// the system will recyle the cluster when creating failed
		return err
	}
	return nil
}

func resourceMRSClusterV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	mrsV1Client, err := cfg.MrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating MRS V1 client: %s", err)
	}
	mrsV2Client, err := cfg.MrsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating MRS V2 client: %s", err)
	}

	vpcId := d.Get("vpc_id").(string)
	vpcResp, err := vpc.GetVpcById(cfg, region, vpcId)
	if err != nil {
		return diag.Errorf("unable to find the vpc (%s) on the server: %s", vpcId, err)
	}
	subnetId := d.Get("subnet_id").(string)
	subnetResp, err := vpc.GetVpcSubnetById(cfg, region, subnetId)
	if err != nil {
		return diag.Errorf("unable to find the subnet (%s) on the server: %s", subnetId, err)
	}

	networkingClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating networking client: %s", err)
	}

	epsID := "all_granted_eps"
	eipId, publicIp, err := queryEipInfo(networkingClient, d.Get("eip_id").(string), d.Get("public_ip").(string), epsID)
	if err != nil {
		return diag.Errorf("unable to find the eip_id=%s,public_ip=%s on the server: %s", d.Get("eip_id").(string),
			d.Get("public_ip").(string), err)
	}

	createOpts := &clusterv2.CreateOpts{
		Region:               region,
		AvailabilityZone:     d.Get("availability_zone").(string),
		ClusterVersion:       d.Get("version").(string),
		ClusterName:          d.Get("name").(string),
		ClusterType:          d.Get("type").(string),
		ManagerAdminPassword: d.Get("manager_admin_pass").(string),
		VpcName:              vpcResp.Name,
		VpcId:                vpcId,
		SubnetId:             subnetId,
		SubnetName:           subnetResp.Name,
		EipId:                eipId,
		EipAddress:           publicIp,
		Components:           strings.Join(utils.ExpandToStringListBySet(d.Get("component_list").(*schema.Set)), ","),
		EnterpriseProjectId:  cfg.GetEnterpriseProjectID(d),
		LogCollection:        buildLogCollection(d),
		NodeGroups:           buildMrsClusterNodeGroups(d),
		SafeMode:             buildMrsSafeMode(d),
		SecurityGroupsIds:    strings.Join(utils.ExpandToStringListBySet(d.Get("security_group_ids").(*schema.Set)), ","),
		ComponentConfigs:     buildComponentConfigOpts(d),
		TemplateId:           d.Get("template_id").(string),
		Tags:                 utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		ExternalDatasources:  buildClusterExternalDatasources(d.Get("external_datasources")),
		BootstrapScripts:     buildBootstrapScripts(d.Get("bootstrap_scripts").(*schema.Set)),
		SMNNotifyConfig:      buildSMNNotify(d),
	}
	if v, ok := d.GetOk("node_key_pair"); ok {
		createOpts.NodeKeypair = v.(string)
		createOpts.LoginMode = "KEYPAIR"
	} else {
		createOpts.NodeRootPassword = d.Get("node_admin_pass").(string)
		createOpts.LoginMode = "PASSWORD"
	}

	// add charge info
	chargingMode := d.Get("charging_mode").(string)
	if chargingMode == "prePaid" {
		createOpts.ChargeInfo = bulidClusterChargeMode(d)
		resp, err := clusterv2.Create(mrsV2Client, createOpts).Extract()
		if err != nil {
			return diag.Errorf("error creating Cluster: %s", err)
		}

		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, resp.OrdeId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, resp.OrdeId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(resourceId)
	} else {
		resp, err := clusterv2.Create(mrsV2Client, createOpts).Extract()
		if err != nil {
			return diag.Errorf("error creating Cluster: %s", err)
		}
		d.SetId(resp.ID)
		refresh := stateRefresh{
			Pending:      []string{"starting"},
			Target:       []string{"running"},
			Timeout:      d.Timeout(schema.TimeoutCreate),
			Delay:        480 * time.Second,
			PollInterval: 15 * time.Second,
		}
		// After request send, check the cluster state and wait for it become running.
		if err = waitForMrsClusterStateCompleted(ctx, mrsV1Client, d.Id(), refresh); err != nil {
			d.SetId("")
			return diag.Errorf("error waiting for MapReduce cluster (%s) to become ready: %s", d.Id(), err)
		}
	}

	return resourceMRSClusterV2Read(ctx, d, meta)
}

func buildSMNNotify(d *schema.ResourceData) *clusterv2.SMNNotifyConfigOpts {
	if v, ok := d.GetOk("smn_notify"); ok {
		smnNotifyConfigs := v.([]interface{})
		smnNotifyConfig := smnNotifyConfigs[0].(map[string]interface{})
		return &clusterv2.SMNNotifyConfigOpts{
			TopicURN:         smnNotifyConfig["topic_urn"].(string),
			SubscriptionName: smnNotifyConfig["subscription_name"].(string),
		}
	}
	return nil
}

func buildClusterExternalDatasources(rawParams interface{}) []clusterv2.ExternalDatasource {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]clusterv2.ExternalDatasource, 0)
		for _, raw := range rawArray {
			item := raw.(map[string]interface{})
			param := clusterv2.ExternalDatasource{
				ComponentName: item["component_name"].(string),
				RoleType:      item["role_type"].(string),
				SourceType:    item["source_type"].(string),
				ConnectorId:   item["data_connection_id"].(string),
			}
			params = append(params, param)
		}
		return params
	}
	return nil
}

func buildBootstrapScripts(rawParams *schema.Set) []clusterv2.ScriptOpts {
	if rawParams.Len() < 1 {
		return nil
	}

	params := make([]clusterv2.ScriptOpts, 0)
	for _, raw := range rawParams.List() {
		if item, ok := raw.(map[string]interface{}); ok {
			param := clusterv2.ScriptOpts{
				Name:                 item["name"].(string),
				URI:                  item["uri"].(string),
				Parameters:           item["parameters"].(string),
				Nodes:                utils.ExpandToStringList(item["nodes"].([]interface{})),
				ActiveMaster:         utils.Bool(item["active_master"].(bool)),
				BeforeComponentStart: utils.Bool(item["before_component_start"].(bool)),
				FailAction:           item["fail_action"].(string),
				ExecuteNeedSudoRoot:  utils.Bool(item["execute_need_sudo_root"].(bool)),
			}
			params = append(params, param)
		}
	}
	return params
}

func queryEipInfo(client *golangsdk.ServiceClient, eipId, publicIp, epsID string) (eipID string, publicIP string, err error) {
	if eipId == "" && publicIp == "" {
		return "", "", nil
	}

	listOpts := eips.ListOpts{
		EnterpriseProjectId: epsID,
	}
	if eipId != "" {
		listOpts.Id = []string{eipId}
	}
	if publicIp != "" {
		listOpts.PublicIp = []string{publicIp}
	}

	pages, err := eips.List(client, listOpts).AllPages()
	if err != nil {
		return "", "", err
	}

	allEips, err := eips.ExtractPublicIPs(pages)
	if err != nil {
		return "", "", fmt.Errorf("unable to retrieve eips: %s ", err)
	}

	if len(allEips) < 1 {
		return "", "", fmt.Errorf("unable to retrieve eips")
	}

	return allEips[0].ID, allEips[0].PublicAddress, nil
}

func setMrsClsuterType(d *schema.ResourceData, resp *cluster.Cluster) error {
	// The returned ClusterType is an 'Int' type, with a value of 0 to 2,
	// which respectively represent:'ANALYSIS','STREAMING' ,'MIXED' and 'CUSTOM'.
	clusterType := []string{"ANALYSIS", "STREAMING", "MIXED", "CUSTOM"}
	if resp.ClusterType >= len(clusterType) || resp.ClusterType < 0 {
		return fmt.Errorf("the cluster type of the response is '%d', not in the map", resp.ClusterType)
	}
	return d.Set("type", clusterType[resp.ClusterType])
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
	return d.Set("create_time", utils.FormatTimeStampRFC3339(createTime, false))
}

func setMrsClsuterUpdateTimestamp(d *schema.ResourceData, resp *cluster.Cluster) error {
	updateTime, err := strconv.ParseInt(resp.Updateat, 10, 64)
	if err != nil {
		return err
	}
	return d.Set("update_time", utils.FormatTimeStampRFC3339(updateTime, false))
}

func setMrsClsuterChargingTimestamp(d *schema.ResourceData, resp *cluster.Cluster) error {
	chargingStartTime, err := strconv.ParseInt(resp.Chargingstarttime, 10, 64)
	if err != nil {
		return err
	}
	return d.Set("charging_start_time", utils.FormatTimeStampRFC3339(chargingStartTime, false))
}

func setMrsClsuterComponentList(d *schema.ResourceData, resp *cluster.Cluster) error {
	result := make([]interface{}, len(resp.Componentlist))
	for i, attachment := range resp.Componentlist {
		result[i] = attachment.Componentname
	}
	return d.Set("component_list", result)
}

func setMrsClusterNodeGroups(d *schema.ResourceData, mrsV1Client *golangsdk.ServiceClient,
	resp *cluster.Cluster) error {
	var groupMapDecl = map[string]string{
		masterGroup:        "master_nodes",
		analysisCoreGroup:  "analysis_core_nodes",
		streamingCoreGroup: "streaming_core_nodes",
		analysisTaskGroup:  "analysis_task_nodes",
		streamingTaskGroup: "streaming_task_nodes",
	}

	clustHostMap, err := queryMrsClusterHosts(d, mrsV1Client)
	if err != nil {
		return err
	}
	var values = make(map[string][]map[string]interface{})

	for _, node := range resp.NodeGroups {
		value, ok := groupMapDecl[node.GroupName]
		isCustomNode := false
		if !ok {
			value = "custom_nodes"
			isCustomNode = true
		}
		groupMap := map[string]interface{}{
			"node_number":      node.NodeNum,
			"flavor":           node.NodeSize,
			"root_volume_type": node.RootVolumeType,
			"root_volume_size": node.RootVolumeSize,
		}
		hostIps := parseHostIps(node.GroupName, isCustomNode, clustHostMap)
		if len(hostIps) == 0 {
			log.Printf("[WARN] One nodeGroup lost host_ips information by some internal error,nodeGroup= %+v", node)
		}
		groupMap["host_ips"] = hostIps

		if isCustomNode {
			groupMap["group_name"] = node.GroupName
		}

		groupMap["assigned_roles"] = node.AssignedRoles
		if node.DataVolumeCount != 0 {
			groupMap["data_volume_type"] = node.DataVolumeType
			groupMap["data_volume_size"] = node.DataVolumeSize
			groupMap["data_volume_count"] = node.DataVolumeCount
		}

		if value != "analysis_task_nodes" && value != "streaming_task_nodes" {
			nodeGroup, ok := d.GetOk(value)
			if ok {
				groupMap["charging_mode"] = utils.PathSearch("[0].charging_mode", nodeGroup, "")
				groupMap["period"] = utils.PathSearch("[0].period", nodeGroup, 1)
				groupMap["period_unit"] = utils.PathSearch("[0].period_unit", nodeGroup, "")
				groupMap["auto_renew"] = utils.PathSearch("[0].auto_renew", nodeGroup, "")
			}
		}

		log.Printf("[DEBUG] node group '%s' is : %+v", value, groupMap)
		values[value] = append(values[value], groupMap)
	}

	for k, v := range values {
		// lintignore:R001
		if err := d.Set(k, v); err != nil {
			return fmt.Errorf("set nodeGroup= %s error", k)
		}
	}

	return nil
}

func parseHostIps(groupName string, isCustomNode bool, clustHostMap map[string][]string) []string {
	var rt []string
	if !isCustomNode {
		if hostIps, ok := clustHostMap[groupName]; ok {
			rt = append(rt, hostIps...)
		}
	} else {
		key := fmt.Sprintf("%s-%s", customNodeGroup, groupName)
		if hostIps, ok := clustHostMap[key]; ok {
			rt = append(rt, hostIps...)
		}
	}

	return rt
}

func flattenTags(tagsString string) map[string]string {
	// the format of tagsRaw: "d=d2,aa=aa"
	result := make(map[string]string)
	if len(tagsString) > 0 {
		tagsArray := strings.Split(tagsString, ",")
		for _, item := range tagsArray {
			tag := strings.SplitN(item, "=", 2)
			if len(tag) == 2 {
				result[tag[0]] = tag[1]
			}
		}
	}
	return result
}

func flattenBootstrapScripts(bootstrapScripts []cluster.BootStrapScript) []map[string]interface{} {
	if size := len(bootstrapScripts); size > 0 {
		result := make([]map[string]interface{}, 0, size)
		for _, v := range bootstrapScripts {
			result = append(result, map[string]interface{}{
				"name":                   v.Name,
				"uri":                    v.URI,
				"nodes":                  v.Nodes,
				"fail_action":            v.FailAction,
				"parameters":             v.Parameters,
				"active_master":          v.ActiveMaster,
				"before_component_start": v.BeforeComponentStart,
				"execute_need_sudo_root": v.ExecuteNeedSudoRoot,
				"start_time":             utils.FormatTimeStampRFC3339(int64(v.StartTime), false),
				"state":                  v.State,
			})
		}
		return result
	}
	return nil
}

func getMrsClusterFromServer(client *golangsdk.ServiceClient, clusterID string) (*cluster.Cluster, error) {
	resp, err := cluster.Get(client, clusterID).Extract()
	if err != nil {
		return nil, err
	}

	if resp.Clusterstate == "terminated" {
		log.Printf("[WARN] Retrieved Cluster %s, but it was terminated, abort it", clusterID)
		return nil, golangsdk.ErrDefault404{}
	}
	return resp, nil
}

func resourceMRSClusterV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.MrsV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	clusterID := d.Id()
	resp, err := getMrsClusterFromServer(client, clusterID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting MapReduce cluster")
	}

	log.Printf("[DEBUG] Retrieved Cluster %s: %#v", clusterID, resp)
	mErr := multierror.Append(
		d.Set("region", resp.Datacenter),
		d.Set("availability_zone", resp.AvailabilityZone),
		d.Set("name", resp.Clustername),
		d.Set("version", resp.Clusterversion),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("node_key_pair", resp.Nodepubliccertname),
		d.Set("vpc_id", resp.Vpcid),
		d.Set("subnet_id", resp.Subnetid),
		d.Set("master_node_ip", resp.Masternodeip),
		d.Set("private_ip", resp.Privateipfirst),
		d.Set("status", resp.Clusterstate),
		d.Set("public_ip", resp.EipAddress),
		d.Set("eip_id", resp.EipId),
		setMrsClsuterType(d, resp),
		setMrsClsuterComponentList(d, resp),
		d.Set("safe_mode", resp.Safemode == 1),
		d.Set("log_collection", resp.LogCollection == 1),
		d.Set("security_group_ids", strings.Split(resp.Securitygroupsid, ",")),
		setMrsClsuterTotalNodeNumber(d, resp),
		setMrsClsuterCreateTimestamp(d, resp),
		setMrsClsuterUpdateTimestamp(d, resp),
		setMrsClsuterChargingTimestamp(d, resp),
		setMrsClusterNodeGroups(d, client, resp),
		d.Set("tags", flattenTags(resp.Tags)),
		d.Set("bootstrap_scripts", flattenBootstrapScripts(resp.BootstrapScripts)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

// resizeMRSClusterCoreNodes is a method which used to resize core node for each cluster type.
// The resizeCount is a number of the group size changing, nagetive means scale in group.
func resizeMRSClusterCoreNodes(ctx context.Context, client *golangsdk.ServiceClient, id, groupType string, resizeCount int) error {
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
		return fmt.Errorf("error resizing core node")
	}
	refresh := stateRefresh{
		Pending:      []string{"scaling-out", "scaling-in"},
		Target:       []string{"running"},
		Delay:        2 * time.Minute,
		Timeout:      1 * time.Hour,
		PollInterval: 15 * time.Second,
	}
	if err = waitForMrsClusterStateCompleted(ctx, client, id, refresh); err != nil {
		return fmt.Errorf("error waiting for MRS cluster resize to be complated: %s", err)
	}
	return nil
}

// resizeMRSClusterTaskNodes is a method which use to scale out/in the (analysis/streaming) nodes.
func resizeMRSClusterTaskNodes(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	groupType, nodeType string) error {
	oldRaws, newRaws := d.GetChange(nodeType)
	oldList := oldRaws.([]interface{})
	newList := newRaws.([]interface{})
	resizeCount := getNodeResizeNumber(oldList, newList)

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
	_, err := cluster.Update(client, d.Id(), opts).Extract()
	if err != nil {
		return fmt.Errorf("error resizing task node")
	}
	refresh := stateRefresh{
		Pending:      []string{"scaling-out", "scaling-in"},
		Target:       []string{"running"},
		Delay:        2 * time.Minute,
		Timeout:      1 * time.Hour,
		PollInterval: 15 * time.Second,
	}
	if err = waitForMrsClusterStateCompleted(ctx, client, d.Id(), refresh); err != nil {
		return fmt.Errorf("error waiting for MRS cluster resize to be complated: %s", err)
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

// calculate the number of the custom group resize option. Dont support add new nodeGroup
func parseCustomNodeResize(oldList, newList []interface{}) map[string]int {
	var rst = make(map[string]int)

	for newIndex := 0; newIndex < len(oldList); newIndex++ {
		newNode := newList[newIndex].(map[string]interface{})

		groupName := newNode["group_name"].(string)
		newSize := newNode["node_number"].(int)

		oldNode := oldList[newIndex].(map[string]interface{})
		oldSize := oldNode["node_number"].(int)

		// Distinguish scale out and scale in by positive and negative
		rst[groupName] = newSize - oldSize
	}
	return rst
}

func updateMRSClusterNodes(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	clusterType := d.Get("type").(string)
	if clusterType == typeAnalysis || clusterType == typeHybrid {
		if d.HasChange("analysis_core_nodes") {
			oldRaws, newRaws := d.GetChange("analysis_core_nodes")
			num := getNodeResizeNumber(oldRaws.([]interface{}), newRaws.([]interface{}))
			err := resizeMRSClusterCoreNodes(ctx, client, d.Id(), analysisCoreGroup, num)
			if err != nil {
				return err
			}
		}
		if d.HasChange("analysis_task_nodes") {
			err := resizeMRSClusterTaskNodes(ctx, client, d, analysisTaskGroup, "analysis_task_nodes")
			if err != nil {
				return err
			}
		}
	}
	if clusterType == typeStream || clusterType == typeHybrid {
		if d.HasChange("streaming_core_nodes") {
			oldRaws, newRaws := d.GetChange("streaming_core_nodes")
			num := getNodeResizeNumber(oldRaws.([]interface{}), newRaws.([]interface{}))
			err := resizeMRSClusterCoreNodes(ctx, client, d.Id(), streamingCoreGroup, num)
			if err != nil {
				return err
			}
		}
		if d.HasChange("streaming_task_nodes") {
			err := resizeMRSClusterTaskNodes(ctx, client, d, streamingTaskGroup, "streaming_task_nodes")
			if err != nil {
				return err
			}
		}
	}

	if clusterType == typeCustom {
		if d.HasChange("custom_nodes") {
			oldRaws, newRaws := d.GetChange("custom_nodes")
			scaleMap := parseCustomNodeResize(oldRaws.([]interface{}), newRaws.([]interface{}))
			for k, num := range scaleMap {
				err := resizeMRSClusterCoreNodes(ctx, client, d.Id(), k, num)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func updateClusterName(d *schema.ResourceData, cfg *config.Config) error {
	region := cfg.GetRegion(d)

	var (
		updateNameHttpUrl = "v2/{project_id}/clusters/{id}/cluster-name"
		updateNameProduct = "mrs"
	)
	updateNameClient, err := cfg.NewServiceClient(updateNameProduct, region)
	if err != nil {
		return fmt.Errorf("error creating MRS Client: %s", err)
	}

	updateNamePath := updateNameClient.Endpoint + updateNameHttpUrl
	updateNamePath = strings.ReplaceAll(updateNamePath, "{project_id}", updateNameClient.ProjectID)
	updateNamePath = strings.ReplaceAll(updateNamePath, "{id}", d.Id())

	updateNameOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	updateNameOpt.JSONBody = utils.RemoveNil(buildUpdateNameBodyParams(d))
	_, err = updateNameClient.Request("PUT", updateNamePath, &updateNameOpt)
	if err != nil {
		return fmt.Errorf("error updating name of MRS cluster: %s", err)
	}

	return nil
}

func buildUpdateNameBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"cluster_name": d.Get("name"),
	}
	return bodyParams
}

func resourceMRSClusterV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	clusterId := d.Id()
	client, err := cfg.MrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	if d.HasChange("name") {
		err = updateClusterName(d, cfg)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		tagErr := updateResourceTagsWithSleep(client, d, "clusters", clusterId)
		if tagErr != nil {
			return diag.Errorf("error updating tags of MRS cluster:%s, err:%s", clusterId, tagErr)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   clusterId,
			ResourceType: "clusters",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	// lintignore:R019
	if d.HasChanges("analysis_core_nodes", "streaming_core_nodes", "analysis_task_nodes",
		"streaming_task_nodes", "custom_nodes") {
		err = updateMRSClusterNodes(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}

		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), clusterId); err != nil {
			return diag.Errorf("error updating the auto-renew of the cluster (%s): %s", clusterId, err)
		}
	}

	return resourceMRSClusterV2Read(ctx, d, meta)
}

func resourceMRSClusterV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.MrsV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	// if charging mode is pre-paid, unsubscribe the order.
	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		if err := common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
			return diag.Errorf("error unsubscribing MRS cluster: %s", err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"running", "terminating"},
			Target:       []string{"terminated", "DELETED"},
			Refresh:      clusterV2StateRefreshFunc(client, d.Id()),
			Timeout:      d.Timeout(schema.TimeoutDelete),
			Delay:        15 * time.Second,
			PollInterval: 10 * time.Second,
		}

		_, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error deleting MRS cluster: %s", err)
		}

		return nil
	}

	err = cluster.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting MRS cluster: %s", err)
	}
	refresh := stateRefresh{
		Pending:      []string{"running", "terminating"},
		Target:       []string{"terminated", "DELETED"},
		Delay:        45 * time.Second,
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 10 * time.Second,
	}
	if err = waitForMrsClusterStateCompleted(ctx, client, d.Id(), refresh); err != nil {
		d.SetId("")
		return diag.Errorf("error waiting for MRS cluster (%s) to be terminated: %s", d.Id(), err)
	}

	return nil
}

/*
When the host type is core, the map key format: {type}-{groupName},
parse from the host name : {clustId}_{groupName}xxxx[-000x]
*/
func queryMrsClusterHosts(d *schema.ResourceData, mrsV1Client *golangsdk.ServiceClient) (map[string][]string, error) {
	clusterId := d.Id()

	hostOpts := cluster.HostOpts{
		CurrentPage: mrsHostDefaultPageNum,
		PageSize:    mrsHostDefaultPageSize,
	}

	resp, err := cluster.ListHosts(mrsV1Client, clusterId, hostOpts)
	if err != nil {
		return nil, fmt.Errorf("query mapreduce cluster host failed: %s", err)
	}
	log.Printf("[DEBUG] Get mapreduce cluster host list response: %#v", resp)
	hostsMap := make(map[string][]string)
	if len(resp.Hosts) > 0 {
		for _, item := range resp.Hosts {
			switch hostType := item.Type; hostType {
			case "Master":
				hostsMap[masterGroup] = append(hostsMap[masterGroup], item.Ip)
			case "Streaming_Core":
				hostsMap[streamingCoreGroup] = append(hostsMap[streamingCoreGroup], item.Ip)
			case "Streaming_Task":
				hostsMap[streamingTaskGroup] = append(hostsMap[streamingTaskGroup], item.Ip)
			case "Analysis_Core":
				hostsMap[analysisCoreGroup] = append(hostsMap[analysisCoreGroup], item.Ip)
			case "Analysis_Task":
				hostsMap[analysisTaskGroup] = append(hostsMap[analysisTaskGroup], item.Ip)
			default:
				reg := regexp.MustCompile(fmt.Sprintf(`%s_(\w*)\w{4}(-\d{4})?`, clusterId))
				allSubMatchStr := reg.FindAllStringSubmatch(item.Name, -1)

				if len(allSubMatchStr) < 1 {
					return nil, fmt.Errorf("parse host info failed. host=%v", item)
				}

				key := fmt.Sprintf("%s-%s", hostType, allSubMatchStr[0][1])
				hostsMap[key] = append(hostsMap[key], item.Ip)
			}
		}
	}

	return hostsMap, nil
}

func updateResourceTagsWithSleep(conn *golangsdk.ServiceClient, d *schema.ResourceData, resourceType, id string) error {
	if d.HasChange("tags") {
		oRaw, nRaw := d.GetChange("tags")
		oMap := oRaw.(map[string]interface{})
		nMap := nRaw.(map[string]interface{})

		// remove old tags
		if len(oMap) > 0 {
			removedTags := utils.ExpandResourceTags(oMap)
			err := tags.Delete(conn, resourceType, id, removedTags).ExtractErr()
			if err != nil {
				return err
			}
			// lintignore:R018
			time.Sleep(5 * time.Second)
		}

		// set new tags
		if len(nMap) > 0 {
			taglist := utils.ExpandResourceTags(nMap)
			err := tags.Create(conn, resourceType, id, taglist).ExtractErr()
			if err != nil {
				return err
			}
			// lintignore:R018
			time.Sleep(5 * time.Second)
		}
	}

	return nil
}
