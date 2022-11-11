package css

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	v1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/css/v1"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/css/v1/model"
	cssv2model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/css/v2/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

const (
	// Instance type. The options are ess, ess-cold, ess-master, and ess-client.
	InstanceTypeEss       = "ess"
	InstanceTypeEssCold   = "ess-cold"
	InstanceTypeEssMaster = "ess-master"
	InstanceTypeEssClient = "ess-client"

	ClusterStatusInProcess   = "100" //The operation, such as instance creation, is in progress.
	ClusterStatusAvailable   = "200"
	ClusterStatusUnavailable = "303"
)

func ResourceCssCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCssClusterCreate,
		ReadContext:   resourceCssClusterRead,
		UpdateContext: resourceCssClusterUpdate,
		DeleteContext: resourceCssClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"engine_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "elasticsearch",
				ValidateFunc: validation.StringInSlice([]string{"elasticsearch", "logstash"}, false),
			},

			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"security_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
				ForceNew:  true,
			},

			"ess_node_config": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				MaxItems:      1,
				ExactlyOneOf:  []string{"node_config", "ess_node_config"},
				ConflictsWith: []string{"expect_node_num"},
				Computed:      true,
				Elem:          cssNodeSchema(1, 200, true),
			},

			"master_node_config": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     cssNodeSchema(3, 10, false),
			},

			"client_node_config": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     cssNodeSchema(1, 32, false),
			},

			"cold_node_config": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     cssNodeSchema(1, 32, true),
			},

			"availability_zone": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"vpc_id", "subnet_id", "security_group_id"},
				Computed:     true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"backup_strategy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Required: true,
						},

						"keep_days": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  7,
						},

						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "snapshot",
						},

						"bucket": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"backup_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"agency": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"public_access": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				RequiredWith: []string{"password"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"whitelist_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"whitelist": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"vpcep_endpoint": { // none query API
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_with_dns_name": { // auto create private domain name
							Type:     schema.TypeBool,
							Required: true,
						},

						"whitelist": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"kibana_public_access": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				RequiredWith: []string{"password"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth": { // Mbit/s
							Type:     schema.TypeInt,
							Required: true,
						},

						"whitelist_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"whitelist": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": common.TagsSchema(),

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),

			"expect_node_num": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "please use ess_node_config.instance_number instead",
				Computed:   true,
			},

			"node_config": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "please use ess_node_config instead",
				ConflictsWith: []string{"master_node_config", "client_node_config", "cold_node_config", "availability_zone"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flavor": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"volume": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:     schema.TypeInt,
										Required: true,
									},

									"volume_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},

						"network_info": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_group_id": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},

									"subnet_id": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},

									"vpc_id": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},

						"availability_zone": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"spec_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vpcep_endpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vpcep_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func cssNodeSchema(min, max int, canExtendsVolume bool) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_number": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(min, max),
			},

			"volume": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:         schema.TypeInt,
							Required:     true,
							ForceNew:     !canExtendsVolume,
							ValidateFunc: validation.IntDivisibleBy(10),
						},
						"volume_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{"COMMON", "HIGH", "ULTRAHIGH"},
								false),
						},
					},
				},
			},
		},
	}
}

func resourceCssClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	cssV1Client, err := config.HcCssV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating CSS V1 client: %s", err)
	}
	cssV2Client, err := config.HcCssV2Client(region)
	if err != nil {
		return diag.Errorf("Error creating CSS V2 client: %s", err)
	}

	createClusterOpts, paramErr := buildClusterCreateParameters(d, config)
	if paramErr != nil {
		return diag.FromErr(paramErr)
	}

	r, err := cssV2Client.CreateCluster(createClusterOpts)
	if err != nil {
		return diag.Errorf("error creating CSS cluster, err=%s", err)
	}

	if (r.Cluster == nil || r.Cluster.Id == nil) && r.OrdeId == nil {
		return diag.Errorf("error creating CSS cluster: id is not found in API response,%#v", r)
	}

	if r.OrdeId == nil {
		if r.Cluster == nil || r.Cluster.Id == nil {
			return diag.Errorf("error creating CSS cluster: id is not found in API response,%#v", r)
		}
		d.SetId(*r.Cluster.Id)

	} else {
		bssClient, err := config.BssV2Client(config.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}

		// 1. If charging mode is PrePaid, wait for the order to be completed.
		err = common.WaitOrderComplete(ctx, bssClient, *r.OrdeId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		// 2. get the resource ID, must be after order success
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, *r.OrdeId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resourceId)
	}

	createResultErr := checkClusterCreateResult(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutCreate))
	if createResultErr != nil {
		return diag.FromErr(createResultErr)
	}

	return resourceCssClusterRead(ctx, d, meta)
}

func buildClusterCreateParameters(d *schema.ResourceData, config *config.Config) (*cssv2model.CreateClusterRequest, error) {
	createOpts := cssv2model.CreateClusterBody{
		Name: d.Get("name").(string),
		Datastore: &cssv2model.CreateClusterDatastoreBody{
			Type:    d.Get("engine_type").(string),
			Version: d.Get("engine_version").(string),
		},
		EnterpriseProjectId: utils.StringIgnoreEmpty(config.GetEnterpriseProjectID(d)),
		Tags:                buildCssTags(d.Get("tags").(map[string]interface{})),
		BackupStrategy:      resourceCssClusterCreateBackupStrategy(d.Get("backup_strategy").([]interface{})),
	}

	if ess, ok := d.GetOk("ess_node_config"); ok {
		essNode := ess.([]interface{})[0].(map[string]interface{})
		createOpts.Roles = append(createOpts.Roles, buildCreateClusterRole(essNode, InstanceTypeEss))

		if v, ok := d.GetOk("master_node_config"); ok {
			node := v.([]interface{})[0].(map[string]interface{})
			createOpts.Roles = append(createOpts.Roles, buildCreateClusterRole(node, InstanceTypeEssMaster))
		}
		if v, ok := d.GetOk("client_node_config"); ok {
			node := v.([]interface{})[0].(map[string]interface{})
			createOpts.Roles = append(createOpts.Roles, buildCreateClusterRole(node, InstanceTypeEssClient))
		}
		if v, ok := d.GetOk("cold_node_config"); ok {
			node := v.([]interface{})[0].(map[string]interface{})
			createOpts.Roles = append(createOpts.Roles, buildCreateClusterRole(node, InstanceTypeEssCold))
		}

		createOpts.AvailabilityZone = utils.StringIgnoreEmpty(d.Get("availability_zone").(string))
		createOpts.Nics = &cssv2model.CreateClusterInstanceNicsBody{
			VpcId:           d.Get("vpc_id").(string),
			NetId:           d.Get("subnet_id").(string),
			SecurityGroupId: d.Get("security_group_id").(string),
		}

	} else {
		// Compatible with previous version
		createOpts.AvailabilityZone = utils.StringIgnoreEmpty(d.Get("node_config.0.availability_zone").(string))
		createOpts.Nics = &cssv2model.CreateClusterInstanceNicsBody{
			VpcId:           d.Get("node_config.0.network_info.0.vpc_id").(string),
			NetId:           d.Get("node_config.0.network_info.0.subnet_id").(string),
			SecurityGroupId: d.Get("node_config.0.network_info.0.security_group_id").(string),
		}

		// add ess role
		createOpts.Roles = append(createOpts.Roles, cssv2model.CreateClusterRolesBody{
			FlavorRef:   d.Get("node_config.0.flavor").(string),
			Type:        InstanceTypeEss,
			InstanceNum: int32(d.Get("expect_node_num").(int)),
			Volume: &cssv2model.CreateClusterInstanceVolumeBody{
				Size:       int32(d.Get("node_config.0.volume.0.size").(int)),
				VolumeType: d.Get("node_config.0.volume.0.volume_type").(string),
			},
		})
	}

	securityMode := d.Get("security_mode").(bool)
	if securityMode {
		adminPassword := d.Get("password").(string)
		if adminPassword == "" {
			return nil, fmtp.Errorf("administrator password is required in security mode")
		}
		createOpts.HttpsEnable = utils.Bool(true)
		createOpts.AuthorityEnable = utils.Bool(true)
		createOpts.AdminPwd = utils.String(adminPassword)
	}

	if _, ok := d.GetOk("vpcep_endpoint"); ok {
		createOpts.LoadBalance = &cssv2model.CreateClusterLoadBalance{
			EndpointWithDnsName: d.Get("vpcep_endpoint.0.endpoint_with_dns_name").(bool),
		}

		vpcPermissions := utils.ExpandToStringList(d.Get("vpcep_endpoint.0.whitelist").([]interface{}))
		if len(vpcPermissions) > 0 {
			createOpts.LoadBalance.VpcPermisssions = &vpcPermissions
		}
	}

	if _, ok := d.GetOk("kibana_public_access"); ok {
		whitelist, ok := d.GetOk("kibana_public_access.0.whitelist")
		createOpts.PublicKibanaReq = &cssv2model.CreateClusterPublicKibanaReq{
			EipSize: int32(d.Get("kibana_public_access.0.bandwidth").(int)),
			ElbWhiteList: &cssv2model.CreateClusterPublicKibanaElbWhiteList{
				EnableWhiteList: ok,
				WhiteList:       whitelist.(string),
			},
		}
	}

	if _, ok := d.GetOk("public_access"); ok {
		whitelist, ok := d.GetOk("public_access.0.whitelist")
		createOpts.PublicIPReq = &cssv2model.CreateClusterPublicIpReq{
			PublicBindType: "auto_assign",
			Eip: &cssv2model.CreateClusterPublicEip{
				BandWidth: &cssv2model.CreateClusterPublicEipSize{
					Size: int32(d.Get("public_access.0.bandwidth").(int)),
				},
			},
			ElbWhiteListReq: &cssv2model.CreateClusterElbWhiteList{
				EnableWhiteList: ok,
				WhiteList:       utils.StringIgnoreEmpty(whitelist.(string)),
			},
		}
	}

	if payModel, ok := d.GetOk("period_unit"); ok || d.Get("charging_mode").(string) == "prePaid" {
		createOpts.PayInfo = &cssv2model.PayInfoBody{
			Period:    int32(d.Get("period").(int)),
			IsAutoPay: utils.Int32(1),
		}

		if payModel == "month" {
			createOpts.PayInfo.PayModel = 2
		} else {
			createOpts.PayInfo.PayModel = 3
		}

		if d.Get("auto_renew").(string) == "true" {
			createOpts.PayInfo.IsAutoRenew = utils.Int32(1)
		}
	}

	return &cssv2model.CreateClusterRequest{Body: &cssv2model.CreateClusterReq{Cluster: &createOpts}}, nil
}

func buildCreateClusterRole(node map[string]interface{}, nodeType string) cssv2model.CreateClusterRolesBody {
	volume := node["volume"].([]interface{})[0].(map[string]interface{})

	return cssv2model.CreateClusterRolesBody{
		FlavorRef:   node["flavor"].(string),
		Type:        nodeType,
		InstanceNum: int32(node["instance_number"].(int)),
		Volume: &cssv2model.CreateClusterInstanceVolumeBody{
			Size:       int32(volume["size"].(int)),
			VolumeType: volume["volume_type"].(string),
		},
	}
}

func buildCssTags(tagmap map[string]interface{}) *[]cssv2model.CreateClusterTagsBody {
	var taglist []cssv2model.CreateClusterTagsBody

	for k, v := range tagmap {
		tag := cssv2model.CreateClusterTagsBody{
			Key:   k,
			Value: v.(string),
		}
		taglist = append(taglist, tag)
	}

	return &taglist
}

func resourceCssClusterCreateBackupStrategy(backupRaw []interface{}) *cssv2model.CreateClusterBackupStrategyBody {
	if len(backupRaw) == 0 {
		return nil
	}
	raw := backupRaw[0].(map[string]interface{})
	opts := cssv2model.CreateClusterBackupStrategyBody{
		Prefix:   raw["prefix"].(string),
		Period:   raw["start_time"].(string),
		Keepday:  int32(raw["keep_days"].(int)),
		Bucket:   utils.StringIgnoreEmpty(raw["bucket"].(string)),
		BasePath: utils.StringIgnoreEmpty(raw["backup_path"].(string)),
		Agency:   utils.StringIgnoreEmpty(raw["agency"].(string)),
	}
	return &opts
}

func resourceCssClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	cssV1Client, err := config.HcCssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	clusterDetail, err := cssV1Client.ShowClusterDetail(&model.ShowClusterDetailRequest{ClusterId: d.Id()})
	if err != nil {
		return diag.Errorf("query cluster detail failed, cluster_id=%s, err=%s", d.Id(), err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", clusterDetail.Name),
		d.Set("engine_type", clusterDetail.Datastore.Type),
		d.Set("engine_version", clusterDetail.Datastore.Version),
		d.Set("enterprise_project_id", clusterDetail.EnterpriseProjectId),
		d.Set("vpc_id", clusterDetail.VpcId),
		d.Set("subnet_id", clusterDetail.SubnetId),
		d.Set("security_group_id", clusterDetail.SecurityGroupId),
		d.Set("nodes", flattenClusterNodes(clusterDetail.Instances)),
		setNodeConfigsAndAzToState(d, clusterDetail),
		setVpcEndpointIdToState(d, cssV1Client),
		d.Set("vpcep_ip", clusterDetail.VpcepIp),
		d.Set("kibana_public_access", flattenKibana(clusterDetail.PublicKibanaResp)),
		d.Set("public_access", flattenPublicAccess(clusterDetail.ElbWhiteList, clusterDetail.BandwidthSize,
			clusterDetail.PublicIp)),
		d.Set("tags", flattenTags(clusterDetail.Tags)),
		d.Set("created", clusterDetail.Created),
		d.Set("endpoint", clusterDetail.Endpoint),
		d.Set("status", clusterDetail.Status),
		d.Set("security_mode", flattenSecurity(clusterDetail.AuthorityEnable)),
		setClusterBackupStrategy(d, cssV1Client),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenClusterNodes(s *[]model.ClusterDetailInstances) []interface{} {
	if s == nil {
		return make([]interface{}, 0)
	}

	rst := make([]interface{}, len(*s))
	for i, v := range *s {
		rst[i] = map[string]interface{}{
			"id":                v.Id,
			"type":              v.Type,
			"name":              v.Name,
			"availability_zone": v.AzCode,
			"status":            v.Status,
			"spec_code":         v.SpecCode,
		}
	}
	return rst
}

func flattenSecurity(authorityEnable *bool) bool {
	if authorityEnable == nil {
		return false
	}
	return *authorityEnable
}

func setClusterBackupStrategy(d *schema.ResourceData, client *v1.CssClient) error {
	policy, err := client.ShowAutoCreatePolicy(&model.ShowAutoCreatePolicyRequest{ClusterId: d.Id()})
	if err != nil {
		return fmt.Errorf("error extracting Cluster:backup_strategy, err: %s", err)
	}

	var strategy []map[string]interface{}
	if utils.StringValue(policy.Enable) == "true" {
		strategy = []map[string]interface{}{
			{
				"prefix":      policy.Prefix,
				"start_time":  policy.Period,
				"keep_days":   policy.Keepday,
				"bucket":      policy.Bucket,
				"backup_path": policy.BasePath,
				"agency":      policy.Agency,
			},
		}
	}
	return d.Set("backup_strategy", strategy)
}

func flattenTags(tags *[]model.ClusterDetailTags) map[string]string {
	if tags == nil {
		return nil
	}

	result := make(map[string]string)
	for _, val := range *tags {
		result[*val.Key] = utils.StringValue(val.Value)
	}
	return result
}

func flattenKibana(publicKibana *model.PublicKibanaRespBody) []interface{} {
	if publicKibana == nil || publicKibana.ElbWhiteListResp == nil {
		return nil
	}

	result := map[string]interface{}{
		"bandwidth":         int(*publicKibana.EipSize),
		"whitelist_enabled": publicKibana.ElbWhiteListResp.EnableWhiteList,
		"whitelist":         publicKibana.ElbWhiteListResp.WhiteList,
		"public_ip":         publicKibana.PublicKibanaIp,
	}
	return []interface{}{result}
}

func flattenPublicAccess(resp *model.ElbWhiteListResp, bandwidth *int32, public_ip *string) []interface{} {
	if resp == nil || public_ip == nil {
		return nil
	}

	result := map[string]interface{}{
		"bandwidth":         int(*bandwidth),
		"whitelist_enabled": resp.EnableWhiteList,
		"whitelist":         resp.WhiteList,
		"public_ip":         public_ip,
	}
	return []interface{}{result}
}

func setVpcEndpointIdToState(d *schema.ResourceData, cssV1Client *v1.CssClient) error {
	resp, err := cssV1Client.ShowVpcepConnection(&model.ShowVpcepConnectionRequest{ClusterId: d.Id()})
	if err != nil {
		if err, ok := err.(*sdkerr.ServiceResponseError); ok {
			errCode := err.ErrorCode
			// CSS.5182 : The VPC endpoint service is not enabled.
			if errCode == "" {
				var apiError CssError
				pErr := json.Unmarshal([]byte(err.ErrorMessage), &apiError)
				if pErr == nil && apiError.ErrorCode == "CSS.5182" {
					return nil
				}
			}

		}
		return err
	}

	var endpointId string
	if resp.Connections != nil && len(*resp.Connections) != 0 {
		connects := *resp.Connections
		endpointId = utils.StringValue(connects[0].Id)
	}

	return d.Set("vpcep_endpoint_id", endpointId)
}

func setNodeConfigsAndAzToState(d *schema.ResourceData, detail *model.ShowClusterDetailResponse) error {
	if detail.Instances == nil || len(*detail.Instances) == 0 {
		return nil
	}
	nodeConfigMap := make(map[string]map[string]interface{})
	var azArray []string
	for _, v := range *detail.Instances {
		azArray = append(azArray, utils.StringValue(v.AzCode))

		nodeType := utils.StringValue(v.Type)
		if node, ok := nodeConfigMap[nodeType]; ok {
			node["instance_number"] = node["instance_number"].(int) + 1
		} else {
			nodeConfigMap[nodeType] = map[string]interface{}{
				"flavor":          v.SpecCode,
				"instance_number": 1,
			}
		}
	}
	azArray = utils.RemoveDuplicateElem(azArray)
	az := strings.Join(azArray, ",")
	mErr := multierror.Append(
		d.Set("availability_zone", az),
	)
	for k, v := range nodeConfigMap {
		switch k {
		case InstanceTypeEss:
			// old version nodeConfig, NO volume return, so get from state
			nodeConfig := map[string]interface{}{
				"flavor":            v["flavor"],
				"availability_zone": az,
				"volume": []interface{}{map[string]interface{}{
					"size":        d.Get("node_config.0.volume.0.size").(int),
					"volume_type": d.Get("node_config.0.volume.0.volume_type").(string),
				}},
				"network_info": []interface{}{map[string]interface{}{
					"vpc_id":            detail.VpcId,
					"subnet_id":         detail.SubnetId,
					"security_group_id": detail.SecurityGroupId,
				}},
			}

			// NO volume return, so get from state
			v["volume"] = []interface{}{map[string]interface{}{
				"size":        d.Get("ess_node_config.0.volume.0.size").(int),
				"volume_type": d.Get("ess_node_config.0.volume.0.volume_type").(string),
			}}
			mErr = multierror.Append(mErr,
				d.Set("node_config", []interface{}{nodeConfig}),
				d.Set("ess_node_config", []interface{}{v}),
				d.Set("expect_node_num", v["instance_number"]),
			)
		case InstanceTypeEssMaster:
			v["volume"] = []interface{}{map[string]interface{}{
				"size":        d.Get("master_node_config.0.volume.0.size").(int),
				"volume_type": d.Get("master_node_config.0.volume.0.volume_type").(string),
			}}
			mErr = multierror.Append(mErr,
				d.Set("master_node_config", []interface{}{v}),
			)
		case InstanceTypeEssClient:
			v["volume"] = []interface{}{map[string]interface{}{
				"size":        d.Get("client_node_config.0.volume.0.size").(int),
				"volume_type": d.Get("client_node_config.0.volume.0.volume_type").(string),
			}}
			mErr = multierror.Append(mErr,
				d.Set("client_node_config", []interface{}{v}),
			)
		case InstanceTypeEssCold:
			v["volume"] = []interface{}{map[string]interface{}{
				"size":        d.Get("cold_node_config.0.volume.0.size").(int),
				"volume_type": d.Get("cold_node_config.0.volume.0.volume_type").(string),
			}}
			mErr = multierror.Append(mErr,
				d.Set("cold_node_config", []interface{}{v}),
			)
		default:
			log.Printf("[ERROR] Does not support to set the %s node config to state", k)
		}

	}
	return mErr.ErrorOrNil()
}

func resourceCssClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	cssV1Client, err := config.HcCssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	// extend cluster
	if d.HasChanges("ess_node_config", "master_node_config", "client_node_config",
		"cold_node_config", "expect_node_num") {
		opts, err := buildCssClusterV1ExtendClusterParameters(d)
		if err != nil {
			return diag.Errorf("error building the request body of api(extend_cluster), err=%s", err)
		}
		_, err = cssV1Client.UpdateExtendInstanceStorage(opts)
		if err != nil {
			return diag.Errorf("extend CSS cluster instance storage failed, cluster_id=%s, error=%s", d.Id(), err)
		}

		err = checkClusterOperationCompleted(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	//update backup strategy
	if d.HasChange("backup_strategy") {
		value, ok := d.GetOk("backup_strategy")
		if !ok {
			// stop auto backup strategy
			_, err := cssV1Client.CreateAutoCreatePolicy(&model.CreateAutoCreatePolicyRequest{
				ClusterId: d.Id(),
				Body: &model.SetRdsBackupCnfReq{
					Prefix:  "snapshot",
					Period:  "00:00 GMT+08:00",
					Keepday: 7,
					Enable:  "false",
				},
			})

			if err != nil {
				return diag.Errorf("error updating backup strategy: %s", err)
			}
		} else {
			rawList := value.([]interface{})
			if len(rawList) == 1 {
				raw := rawList[0].(map[string]interface{})

				if d.HasChanges("backup_strategy.0.bucket", "backup_strategy.0.backup_path",
					"backup_strategy.0.agency") {
					// If obs is specified, update basic configurations
					_, err = cssV1Client.UpdateSnapshotSetting(&model.UpdateSnapshotSettingRequest{
						ClusterId: d.Id(),
						Body: &model.UpdateSnapshotSettingReq{
							Bucket:   raw["bucket"].(string),
							BasePath: raw["backup_path"].(string),
							Agency:   raw["agency"].(string),
						},
					})
					if err != nil {
						return diag.Errorf("error Modifying Basic Configurations of a Cluster Snapshot: %s", err)
					}
				}

				// check backup strategy, if the policy was disabled, we should enable it
				policy, err := cssV1Client.ShowAutoCreatePolicy(&model.ShowAutoCreatePolicyRequest{ClusterId: d.Id()})
				if err != nil {
					return diag.Errorf("Error extracting Cluster backup_strategy, err: %s", err)
				}

				if utils.StringValue(policy.Enable) == "false" && raw["bucket"] == nil {
					// If obs is not specified,  create  basic configurations automatically
					_, err = cssV1Client.StartAutoSetting(&model.StartAutoSettingRequest{ClusterId: d.Id()})
					if err != nil {
						return diag.Errorf("error enable snapshot function: %s", err)
					}
				}

				// update policy
				if d.HasChanges("backup_strategy.0.prefix", "backup_strategy.0.start_time",
					"backup_strategy.0.keep_days") {
					opts := &model.CreateAutoCreatePolicyRequest{
						ClusterId: d.Id(),
						Body: &model.SetRdsBackupCnfReq{
							Prefix:  raw["prefix"].(string),
							Period:  raw["start_time"].(string),
							Keepday: int32(raw["keep_days"].(int)),
							Enable:  "true",
						},
					}
					_, err = cssV1Client.CreateAutoCreatePolicy(opts)
					if err != nil {
						return diag.Errorf("error updating backup strategy: %s", err)
					}
				}
			}
		}
	}

	if d.HasChange("tags") {
		oRaw, nRaw := d.GetChange("tags")
		err = updateCssTags(cssV1Client, d.Id(), oRaw.(map[string]interface{}), nRaw.(map[string]interface{}))
		if err != nil {
			return diag.Errorf("error updating tags of CSS cluster= %s, err:%s", d.Id(), err)
		}
	}

	// update vpc endpoint
	if d.HasChange("vpcep_endpoint") {
		o, n := d.GetChange("vpcep_endpoint")
		oValue := o.([]interface{})
		nValue := n.([]interface{})
		if len(nValue) == 0 { // delete vpc endpoint
			_, err := cssV1Client.StopVpecp(&model.StopVpecpRequest{ClusterId: d.Id()})
			if err != nil {
				return diag.Errorf("error deleting the VPC endpoint of CSS cluster= %s, err: %s", d.Id(), err)
			}
		} else if len(oValue) == 0 { // start vpc endpoint
			_, err := cssV1Client.StartVpecp(&model.StartVpecpRequest{
				ClusterId: d.Id(),
				Body: &model.StartVpecpReq{
					EndpointWithDnsName: d.Get("vpcep_endpoint.0.endpoint_with_dns_name").(bool),
				},
			})
			if err != nil {
				return diag.Errorf("error creating the VPC endpoint of CSS cluster= %s, err: %s", d.Id(), err)
			}
		}

		err = checkClusterOperationCompleted(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}

		// update whitelist
		if len(nValue) > 0 && d.HasChange("vpcep_endpoint.0.whitelist") {
			_, err = cssV1Client.UpdateVpcepWhitelist(&model.UpdateVpcepWhitelistRequest{
				ClusterId: d.Id(),
				Body: &model.UpdateVpcepWhitelistReq{
					VpcPermissions: utils.ExpandToStringList(d.Get("vpcep_endpoint.0.whitelist").([]interface{})),
				},
			})
			if err != nil {
				return diag.Errorf("error updating the VPC endpoint whitelist of CSS cluster= %s, err: %s", d.Id(), err)
			}
		}
	}

	// update kibana
	if d.Get("security_mode").(bool) && d.HasChange("kibana_public_access") {
		o, n := d.GetChange("kibana_public_access")
		oValue := o.([]interface{})
		nValue := n.([]interface{})
		if len(nValue) == 0 { // delete kibana_public_access
			_, err := cssV1Client.UpdateCloseKibana(&model.UpdateCloseKibanaRequest{ClusterId: d.Id()})
			if err != nil {
				return diag.Errorf("error diabling the kibana public access of CSS cluster= %s, err: %s", d.Id(), err)
			}
			err = checkClusterOperationCompleted(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		} else if len(oValue) == 0 {
			// enable kibana_public_access
			_, err := cssV1Client.StartKibanaPublic(&model.StartKibanaPublicRequest{
				ClusterId: d.Id(),
				Body: &model.StartKibanaPublicReq{
					EipSize: int32(d.Get("kibana_public_access.0.bandwidth").(int)),
					ElbWhiteList: &model.StartKibanaPublicReqElbWhitelist{
						EnableWhiteList: d.Get("kibana_public_access.0.whitelist_enabled").(bool),
						WhiteList:       d.Get("kibana_public_access.0.whitelist").(string),
					},
				},
			})
			if err != nil {
				return diag.Errorf("error enabling the kibana public access of CSS cluster= %s, err: %s", d.Id(), err)
			}
			err = checkClusterOperationCompleted(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			//update bandwidth
			if d.HasChange("kibana_public_access.0.bandwidth") {
				_, err := cssV1Client.UpdateAlterKibana(&model.UpdateAlterKibanaRequest{
					ClusterId: d.Id(),
					Body: &model.UpdatePublicKibanaBandwidthReq{
						BandWidth: &model.UpdatePublicKibanaBandwidthReqBandWidth{
							Size: int32(d.Get("kibana_public_access.0.bandwidth").(int)),
						},
						IsAutoPay: utils.Int32(1),
					},
				})
				if err != nil {
					return diag.Errorf("error modifing bandwidth of the kibana public access of CSS cluster= %s, err: %s", d.Id(), err)
				}
			}

			// update whitelist
			if d.HasChanges("kibana_public_access.0.whitelist", "kibana_public_access.0.whitelist_enabled") {
				// disable whitelist
				if !d.Get("kibana_public_access.0.whitelist_enabled").(bool) {
					_, err := cssV1Client.StopPublicKibanaWhitelist(&model.StopPublicKibanaWhitelistRequest{
						ClusterId: d.Id(),
					})
					if err != nil {
						return diag.Errorf("error disabing the whitelist of the kibana public access of CSS cluster= %s, err: %s", d.Id(), err)
					}
				} else {
					_, err := cssV1Client.UpdatePublicKibanaWhitelist(&model.UpdatePublicKibanaWhitelistRequest{
						ClusterId: d.Id(),
						Body: &model.UpdatePublicKibanaWhitelistReq{
							WhiteList: d.Get("kibana_public_access.0.whitelist").(string),
						},
					})
					if err != nil {
						return diag.Errorf("error modifing whitelist of the kibana public access of CSS cluster= %s, err: %s", d.Id(), err)
					}
				}
			}
		}
	}

	// update public_access
	if d.Get("security_mode").(bool) && d.HasChange("public_access") {
		o, n := d.GetChange("public_access")
		oValue := o.([]interface{})
		nValue := n.([]interface{})
		if len(nValue) == 0 { // delete public_access
			_, err := cssV1Client.UpdateUnbindPublic(&model.UpdateUnbindPublicRequest{ClusterId: d.Id()})
			if err != nil {
				return diag.Errorf("error diabling public access of CSS cluster= %s, err: %s", d.Id(), err)
			}
			err = checkClusterOperationCompleted(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		} else if len(oValue) == 0 {
			// enable public_access
			_, err := cssV1Client.CreateBindPublic(&model.CreateBindPublicRequest{
				ClusterId: d.Id(),
				Body: &model.BindPublicReq{
					Eip: &model.BindPublicReqEip{
						BandWidth: &model.BindPublicReqEipBandWidth{
							Size: int32(d.Get("public_access.0.bandwidth").(int)),
						},
					},
					IsAutoPay: utils.Int32(1),
				},
			})

			if err != nil {
				return diag.Errorf("error enabling public access of CSS cluster= %s, err: %s", d.Id(), err)
			}

			err = checkClusterOperationCompleted(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if len(nValue) > 0 && d.HasChanges("public_access.0.whitelist", "public_access.0.whitelist_enabled") {
			// disable whitelist
			if !d.Get("kibana_public_access.0.whitelist_enabled").(bool) {
				_, err := cssV1Client.StopPublicWhitelist(&model.StopPublicWhitelistRequest{ClusterId: d.Id()})
				if err != nil {
					return diag.Errorf("error disabling whitelist of public access of CSS cluster= %s, err: %s", d.Id(), err)
				}
			} else {
				_, err := cssV1Client.StartPublicWhitelist(&model.StartPublicWhitelistRequest{
					ClusterId: d.Id(),
					Body: &model.StartPublicWhitelistReq{
						WhiteList: d.Get("kibana_public_access.0.whitelist").(string),
					},
				})
				if err != nil {
					return diag.Errorf("error modifing whitelist of public access of CSS cluster= %s, err: %s", d.Id(), err)
				}
			}
		}

		// update bandwidth
		if len(nValue) > 0 && d.HasChange("public_access.0.bandwidth") {
			_, err := cssV1Client.UpdatePublicBandWidth(&model.UpdatePublicBandWidthRequest{
				ClusterId: d.Id(),
				Body: &model.BindPublicReqEipReq{
					BandWidth: &model.BindPublicReqEipBandWidth{
						Size: int32(d.Get("public_access.0.bandwidth").(int)),
					},
					IsAutoPay: utils.Int32(1),
				},
			})
			if err != nil {
				return diag.Errorf("error disabling the whitelist of the kibana public access of CSS cluster= %s, err: %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := config.BssV2Client(config.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
			return diag.Errorf("error updating the auto-renew of the cluster (%s): %s", d.Id(), err)
		}
	}

	return resourceCssClusterRead(ctx, d, meta)
}

func buildCssClusterV1ExtendClusterParameters(d *schema.ResourceData) (*model.UpdateExtendInstanceStorageRequest, error) {
	var grow = make([]model.RoleExtendGrowReq, 0, 4)

	if d.HasChange("ess_node_config") {
		oldv, newv := d.GetChange("ess_node_config.0.instance_number")
		nodesize := newv.(int) - oldv.(int)
		if nodesize < 0 {
			return nil, fmt.Errorf("instance_number only supports to be extended")
		}

		oldDisksize, newDisksize := d.GetChange("ess_node_config.0.volume.0.size")
		disksize := newDisksize.(int) - oldDisksize.(int)
		if disksize < 0 {
			return nil, fmt.Errorf("volume size only supports to be extended")
		}

		grow = append(grow, model.RoleExtendGrowReq{
			Type:     InstanceTypeEss,
			Nodesize: int32(nodesize),
			Disksize: int32(disksize),
		})
	}

	if d.HasChange("cold_node_config") {
		oldv, newv := d.GetChange("cold_node_config.0.instance_number")
		nodesize := newv.(int) - oldv.(int)
		if nodesize < 0 {
			return nil, fmt.Errorf("instance_number only supports to be extended")
		}

		oldDisksize, newDisksize := d.GetChange("cold_node_config.0.volume.0.size")
		disksize := newDisksize.(int) - oldDisksize.(int)
		if disksize < 0 {
			return nil, fmt.Errorf("volume size only supports to be extended")
		}

		grow = append(grow, model.RoleExtendGrowReq{
			Type:     InstanceTypeEssCold,
			Nodesize: int32(nodesize),
			Disksize: int32(disksize),
		})
	}

	if d.HasChange("master_node_config") {
		oldv, newv := d.GetChange("master_node_config.0.instance_number")
		nodesize := newv.(int) - oldv.(int)
		if nodesize < 0 {
			return nil, fmt.Errorf("instance_number only supports to be extended")
		}

		grow = append(grow, model.RoleExtendGrowReq{
			Type:     InstanceTypeEssMaster,
			Nodesize: int32(nodesize),
		})
	}

	if d.HasChange("client_node_config") {
		oldv, newv := d.GetChange("client_node_config.0.instance_number")
		nodesize := newv.(int) - oldv.(int)
		if nodesize < 0 {
			return nil, fmt.Errorf("instance_number only supports to be extended")
		}

		grow = append(grow, model.RoleExtendGrowReq{
			Type:     InstanceTypeEssClient,
			Nodesize: int32(nodesize),
		})
	}

	if d.HasChanges("node_config.0.volume.0.size", "expect_node_num") {
		oldv, newv := d.GetChange("expect_node_num")
		nodesize := newv.(int) - oldv.(int)
		if nodesize < 0 {
			return nil, fmt.Errorf("expect_node_num only supports to be extended")
		}

		oldDisksize, newDisksize := d.GetChange("node_config.0.volume.0.size")
		disksize := newDisksize.(int) - oldDisksize.(int)
		if disksize < 0 {
			return nil, fmt.Errorf("volume size only supports to be extended")
		}

		grow = append(grow, model.RoleExtendGrowReq{
			Type:     InstanceTypeEss,
			Nodesize: 1,
			Disksize: int32(disksize),
		})
	}

	return &model.UpdateExtendInstanceStorageRequest{
		ClusterId: d.Id(),
		Body: &model.RoleExtendReq{
			Grow: grow,
		},
	}, nil
}

func resourceCssClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	cssV1Client, err := config.HcCssV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating CSS V1 client: %s", err)
	}

	_, err = cssV1Client.DeleteCluster(&model.DeleteClusterRequest{ClusterId: d.Id()})
	if err != nil {
		return diag.Errorf("delete CSS Cluster %s failed, error= %s", d.Id(), err)
	}

	err = checkClusterDeleteResult(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("failed to check the result of deletion %s", err)
	}
	d.SetId("")
	return nil
}

func checkClusterCreateResult(ctx context.Context, cssV1Client *v1.CssClient, clusterId string,
	timeout time.Duration) error {
	createStateConf := &resource.StateChangeConf{
		Pending: []string{ClusterStatusInProcess},
		Target:  []string{ClusterStatusAvailable},
		Refresh: func() (interface{}, string, error) {
			resp, err := cssV1Client.ShowClusterDetail(&model.ShowClusterDetailRequest{ClusterId: clusterId})
			if err != nil {
				return nil, "failed", err
			}
			return resp, *resp.Status, err
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CSS (%s) to be created: %s", clusterId, err)
	}
	return nil
}

func checkClusterDeleteResult(ctx context.Context, cssV1Client *v1.CssClient, clusterId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Done"},
		Refresh: func() (interface{}, string, error) {
			_, err := cssV1Client.ShowClusterDetail(&model.ShowClusterDetailRequest{ClusterId: clusterId})
			if err != nil {
				if err, ok := err.(*sdkerr.ServiceResponseError); ok {
					if err.StatusCode == http.StatusNotFound {
						return true, "Done", nil
					}

					if err.StatusCode == 403 {
						var apiError CssError
						pErr := json.Unmarshal([]byte(err.ErrorMessage), &apiError)
						if pErr == nil && apiError.ErrorCode == "CSS.0015" {
							return true, "Done", nil
						}
					}
				}
				return nil, "ERROR", err
			}
			return true, "Pending", nil
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CSS (%s) to be delete: %s", clusterId, err)
	}
	return nil
}

func checkClusterOperationCompleted(ctx context.Context, cssV1Client *v1.CssClient, clusterId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Done"},
		Refresh: func() (interface{}, string, error) {
			resp, err := cssV1Client.ShowClusterDetail(&model.ShowClusterDetailRequest{ClusterId: clusterId})
			if err != nil {
				return nil, "failed", err
			}

			if checkCssClusterIsReady(resp) {
				return resp, "Done", nil
			}
			return resp, "Pending", nil
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CSS (%s) to be extend: %s", clusterId, err)
	}
	return nil
}

func checkCssClusterIsReady(detail *model.ShowClusterDetailResponse) bool {
	if utils.StringValue(detail.Status) != ClusterStatusAvailable {
		return false
	}

	//actions --- the behaviors on a cluster
	if detail.Actions != nil && len(*detail.Actions) > 0 {
		return false
	}

	if detail.Instances == nil {
		return false
	}
	for _, v := range *detail.Instances {
		if utils.StringValue(v.Status) != ClusterStatusAvailable {
			return false
		}
	}
	return true
}

func updateCssTags(cssV1Client *v1.CssClient, id string, old, new map[string]interface{}) error {
	// remove old tags
	for k := range old {
		_, err := cssV1Client.DeleteClustersTags(&model.DeleteClustersTagsRequest{
			ResourceType: "css-cluster",
			ClusterId:    id,
			Key:          k,
		})
		if err != nil {
			return err
		}
	}

	// set new tags
	for k, v := range new {
		_, err := cssV1Client.CreateClustersTags(&model.CreateClustersTagsRequest{
			ResourceType: "css-cluster",
			ClusterId:    id,
			Body: &model.TagReq{
				Tag: &model.Tag{
					Key:   k,
					Value: v.(string),
				},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

type CssError struct {
	ErrorCode string `json:"errCode"`
	ErrorMsg  string `json:"externalMessage"`
}
