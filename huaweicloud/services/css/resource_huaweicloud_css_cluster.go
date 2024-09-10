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

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	cssv1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/css/v1"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/css/v1/model"
	cssv2model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/css/v2/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	// Instance type. The options are ess, ess-cold, ess-master, and ess-client.
	InstanceTypeEss       = "ess"
	InstanceTypeEssCold   = "ess-cold"
	InstanceTypeEssMaster = "ess-master"
	InstanceTypeEssClient = "ess-client"

	ClusterStatusInProcess   = "100" // The operation, such as instance creation, is in progress.
	ClusterStatusAvailable   = "200"
	ClusterStatusUnavailable = "303"
)

var clusterNonUpdatableParams = []string{"engine_version", "availability_zone"}

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/role_extend
// @API CSS POST /v1.0/{project_id}/clusters
// @API CSS POST /v1.0/{project_id}/{resource_type}/{cluster_id}/tags
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/public/open
// @API CSS PUT /v1.0/{project_id}/clusters/{cluster_id}/public/whitelist/close
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/publickibana/bandwidth
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/publickibana/whitelist/update
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/setting
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/vpcepservice/open
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/vpcepservice/permissions
// @API CSS DELETE /v1.0/{project_id}/{resource_type}/{cluster_id}/tags/{key}
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/auto_setting
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/public/whitelist/update
// @API CSS DELETE /v1.0/{project_id}/clusters/{cluster_id}
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}
// @API CSS POST /v2.0/{project_id}/clusters
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/publickibana/open
// @API CSS PUT /v1.0/{project_id}/clusters/{cluster_id}/publickibana/whitelist/close
// @API CSS PUT /v1.0/{project_id}/clusters/{cluster_id}/vpcepservice/close
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/vpcepservice/connections
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/policy
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/policy
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/public/bandwidth
// @API CSS PUT /v1.0/{project_id}/clusters/{cluster_id}/public/close
// @API CSS PUT /v1.0/{project_id}/clusters/{cluster_id}/publickibana/close
// @API CSS GET /v1.0/{project_id}/es-flavors
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/{types}/flavor
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/type/{type}/independent
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/mode/change
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/password/reset
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/sg/change
// @API CSS POST /v1.0/{project_id}/cluster/{cluster_id}/period
// @API CSS POST /v1.0/extend/{project_id}/clusters/{cluster_id}/role/shrink
// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/node/offline
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
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

		CustomizeDiff: config.FlexibleForceNew(clusterNonUpdatableParams),

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
			},
			"security_mode": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"https_enabled": {
				Type:         schema.TypeBool,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				RequiredWith: []string{"security_mode"},
			},
			"ess_node_config": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ExactlyOneOf:  []string{"node_config", "ess_node_config"},
				ConflictsWith: []string{"expect_node_num"},
				Computed:      true,
				Elem:          essOrColdNodeSchema(),
				Description:   "schema: Required",
			},
			"master_node_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     masterOrClientNodeSchema(),
			},
			"client_node_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     masterOrClientNodeSchema(),
			},
			"cold_node_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     essOrColdNodeSchema(),
			},
			"availability_zone": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"vpc_id", "subnet_id", "security_group_id"},
				Computed:     true,
				Description:  "schema: Required",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "schema: Required",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "schema: Required",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "schema: Required",
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
				Computed: true,
			},
			// charging_mode, period_unit and period only support changing post-paid to pre-paid billing mode.
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid",
				}, false),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"period"},
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"period_unit"},
			},
			"auto_renew": common.SchemaAutoRenewUpdatable(nil),
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
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_period": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"backup_available": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disk_encrypted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "schema: Deprecated; use created_at instead",
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func essOrColdNodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},

			"instance_number": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"volume": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
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
			"shrink_node_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func masterOrClientNodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},

			"instance_number": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"volume": {
				Type:     schema.TypeList,
				Required: true,
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
						},
					},
				},
			},
			"shrink_node_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceCssClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.HcCssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}
	cssV2Client, err := conf.HcCssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V2 client: %s", err)
	}

	createClusterOpts, paramErr := buildClusterCreateParameters(d, conf)
	if paramErr != nil {
		return diag.FromErr(paramErr)
	}

	r, err := cssV2Client.CreateCluster(createClusterOpts)
	if err != nil {
		return diag.Errorf("error creating CSS cluster, err: %s", err)
	}

	if (r.Cluster == nil || r.Cluster.Id == nil) && r.OrderId == nil {
		return diag.Errorf("error creating CSS cluster: id is not found in API response,%#v", r)
	}

	if r.OrderId == nil {
		if r.Cluster == nil || r.Cluster.Id == nil {
			return diag.Errorf("error creating CSS cluster: id is not found in API response,%#v", r)
		}
		d.SetId(*r.Cluster.Id)
	} else {
		bssClient, err := conf.BssV2Client(conf.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}

		// 1. If charging mode is PrePaid, wait for the order to be completed.
		err = common.WaitOrderComplete(ctx, bssClient, *r.OrderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		// 2. get the resource ID, must be after order success
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, *r.OrderId, d.Timeout(schema.TimeoutCreate))
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

func buildClusterCreateParameters(d *schema.ResourceData, conf *config.Config) (*cssv2model.CreateClusterRequest, error) {
	createOpts := cssv2model.CreateClusterBody{
		Name: d.Get("name").(string),
		Datastore: &cssv2model.CreateClusterDatastoreBody{
			Type:    d.Get("engine_type").(string),
			Version: d.Get("engine_version").(string),
		},
		EnterpriseProjectId: utils.StringIgnoreEmpty(conf.GetEnterpriseProjectID(d)),
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
			return nil, fmt.Errorf("administrator password is required in security mode")
		}
		createOpts.AuthorityEnable = utils.Bool(true)
		createOpts.AdminPwd = utils.String(adminPassword)

		createOpts.HttpsEnable = utils.Bool(d.Get("https_enabled").(bool))
	}

	if _, ok := d.GetOk("vpcep_endpoint"); ok {
		createOpts.LoadBalance = &cssv2model.CreateClusterLoadBalance{
			EndpointWithDnsName: d.Get("vpcep_endpoint.0.endpoint_with_dns_name").(bool),
		}

		vpcPermissions := utils.ExpandToStringList(d.Get("vpcep_endpoint.0.whitelist").([]interface{}))
		if len(vpcPermissions) > 0 {
			createOpts.LoadBalance.VpcPermissions = &vpcPermissions
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

	if payModel, ok := d.GetOk("period_unit"); ok && d.Get("charging_mode").(string) != "postPaid" {
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
	clusterRolesBody := cssv2model.CreateClusterRolesBody{
		FlavorRef:   node["flavor"].(string),
		Type:        nodeType,
		InstanceNum: int32(node["instance_number"].(int)),
	}

	// Ess node and cold node support local disk. The volume value is empty.
	// Master node and client node do not support local volume. The volume value is required.
	if volumes := node["volume"].([]interface{}); len(volumes) > 0 {
		volume := volumes[0].(map[string]interface{})
		clusterRolesBody.Volume = &cssv2model.CreateClusterInstanceVolumeBody{
			Size:       int32(volume["size"].(int)),
			VolumeType: volume["volume_type"].(string),
		}
	}
	return clusterRolesBody
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

func resourceCssClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.HcCssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	clusterDetail, err := cssV1Client.ShowClusterDetail(&model.ShowClusterDetailRequest{ClusterId: d.Id()})
	if err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		err = ConvertExpectedHwSdkErrInto404Err(err, http.StatusForbidden, "CSS.0015", "")
		return common.CheckDeletedDiag(d, err, "CSS cluster")
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
		d.Set("https_enabled", clusterDetail.HttpsEnable),
		setClusterBackupStrategy(d, cssV1Client),
		d.Set("created_at", clusterDetail.Created),
		d.Set("updated_at", clusterDetail.Updated),
		d.Set("bandwidth_resource_id", clusterDetail.BandwidthResourceId),
		d.Set("is_period", clusterDetail.Period),
		d.Set("backup_available", clusterDetail.BackupAvailable),
		d.Set("disk_encrypted", clusterDetail.DiskEncrypted),
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
			"ip":                v.Ip,
			"resource_id":       v.ResourceId,
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

func setClusterBackupStrategy(d *schema.ResourceData, client *cssv1.CssClient) error {
	policy, err := client.ShowAutoCreatePolicy(&model.ShowAutoCreatePolicyRequest{ClusterId: d.Id()})
	if err != nil {
		return fmt.Errorf("error extracting cluster: backup_strategy, err: %s", err)
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

func flattenPublicAccess(resp *model.ElbWhiteListResp, bandwidth *int32, publicIp *string) []interface{} {
	if resp == nil || publicIp == nil {
		return nil
	}

	result := map[string]interface{}{
		"bandwidth":         int(*bandwidth),
		"whitelist_enabled": resp.EnableWhiteList,
		"whitelist":         resp.WhiteList,
		"public_ip":         publicIp,
	}
	return []interface{}{result}
}

func setVpcEndpointIdToState(d *schema.ResourceData, cssV1Client *cssv1.CssClient) error {
	resp, err := cssV1Client.ShowVpcepConnection(&model.ShowVpcepConnectionRequest{ClusterId: d.Id()})
	if err != nil {
		if err, ok := err.(*sdkerr.ServiceResponseError); ok {
			errCode := err.ErrorCode
			// CSS.5182 : The VPC endpoint service is not enabled.
			if errCode == "" {
				var apiError ResponseError
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
				"volume": []interface{}{map[string]interface{}{
					"size":        v.Volume.Size,
					"volume_type": v.Volume.Type,
				}},
				"shrink_node_ids": nil,
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
			volumeSize := utils.PathSearch("volume|[0]|size", v, 0).(*int32)
			volumeType := utils.PathSearch("volume|[0]|volume_type", v, "").(*string)
			nodeConfig := map[string]interface{}{
				"flavor":            v["flavor"],
				"availability_zone": az,
				"volume": []interface{}{map[string]interface{}{
					"size":        *volumeSize,
					"volume_type": *volumeType,
				}},
				"network_info": []interface{}{map[string]interface{}{
					"vpc_id":            detail.VpcId,
					"subnet_id":         detail.SubnetId,
					"security_group_id": detail.SecurityGroupId,
				}},
			}

			// NO volume return, so get from state
			if *volumeSize == 0 && *volumeType == "" {
				v["volume"] = []interface{}{map[string]interface{}{
					"size":        d.Get("ess_node_config.0.volume.0.size").(int),
					"volume_type": d.Get("ess_node_config.0.volume.0.volume_type").(string),
				}}
			}
			if ids, ok := d.GetOk("ess_node_config.0.shrink_node_ids"); ok {
				v["shrink_node_ids"] = ids.([]interface{})
			}
			mErr = multierror.Append(mErr,
				d.Set("node_config", []interface{}{nodeConfig}),
				d.Set("ess_node_config", []interface{}{v}),
				d.Set("expect_node_num", v["instance_number"]),
			)
		case InstanceTypeEssMaster:
			volumeSize := utils.PathSearch("volume|[0]|size", v, 0).(*int32)
			volumeType := utils.PathSearch("volume|[0]|volume_type", v, "").(*string)
			if *volumeSize == 0 && *volumeType == "" {
				v["volume"] = []interface{}{map[string]interface{}{
					"size":        d.Get("master_node_config.0.volume.0.size").(int),
					"volume_type": d.Get("master_node_config.0.volume.0.volume_type").(string),
				}}
			}
			if ids, ok := d.GetOk("master_node_config.0.shrink_node_ids"); ok {
				v["shrink_node_ids"] = ids.([]interface{})
			}
			mErr = multierror.Append(mErr,
				d.Set("master_node_config", []interface{}{v}),
			)
		case InstanceTypeEssClient:
			volumeSize := utils.PathSearch("volume|[0]|size", v, 0).(*int32)
			volumeType := utils.PathSearch("volume|[0]|volume_type", v, "").(*string)
			if *volumeSize == 0 && *volumeType == "" {
				v["volume"] = []interface{}{map[string]interface{}{
					"size":        d.Get("client_node_config.0.volume.0.size").(int),
					"volume_type": d.Get("client_node_config.0.volume.0.volume_type").(string),
				}}
			}
			if ids, ok := d.GetOk("client_node_config.0.shrink_node_ids"); ok {
				v["shrink_node_ids"] = ids.([]interface{})
			}
			mErr = multierror.Append(mErr,
				d.Set("client_node_config", []interface{}{v}),
			)
		case InstanceTypeEssCold:
			volumeSize := utils.PathSearch("volume|[0]|size", v, 0).(*int32)
			volumeType := utils.PathSearch("volume|[0]|volume_type", v, "").(*string)
			if *volumeSize == 0 && *volumeType == "" {
				v["volume"] = []interface{}{map[string]interface{}{
					"size":        d.Get("cold_node_config.0.volume.0.size").(int),
					"volume_type": d.Get("cold_node_config.0.volume.0.volume_type").(string),
				}}
			}
			if ids, ok := d.GetOk("cold_node_config.0.shrink_node_ids"); ok {
				v["shrink_node_ids"] = ids.([]interface{})
			}
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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	clusterId := d.Id()
	cssV1Client, err := cfg.HcCssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	client, err := cfg.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	nodeConfigChanges := []string{
		"ess_node_config",
		"master_node_config",
		"client_node_config",
		"cold_node_config",
	}

	if d.HasChanges(nodeConfigChanges...) {
		err := updateNodeConfig(ctx, d, cssV1Client, cfg)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update backup strategy
	if d.HasChange("backup_strategy") {
		err = updateBackupStrategy(d, cssV1Client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		oRaw, nRaw := d.GetChange("tags")
		err = updateCssTags(cssV1Client, clusterId, oRaw.(map[string]interface{}), nRaw.(map[string]interface{}))
		if err != nil {
			return diag.Errorf("error updating tags of CSS cluster: %s, err: %s", clusterId, err)
		}
	}

	// update vpc endpoint
	if d.HasChange("vpcep_endpoint") {
		err = updateVpcepEndpoint(ctx, d, cssV1Client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update security mode
	if d.HasChange("security_mode") {
		err = updateSafeMode(ctx, d, cssV1Client, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update in safe mode
	if d.Get("security_mode").(bool) {
		err = updateInSafeMode(ctx, d, cssV1Client, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update security group ID
	if d.HasChange("security_group_id") {
		err = updateSecurityGroup(ctx, d, cssV1Client, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("charging_mode", "auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}

		err = updateChangingModeOrAutoRenew(ctx, d, cssV1Client, bssClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   clusterId,
			ResourceType: "css-cluster",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceCssClusterRead(ctx, d, meta)
}

func updateInSafeMode(ctx context.Context, d *schema.ResourceData,
	cssV1Client *cssv1.CssClient, client *golangsdk.ServiceClient) error {
	// reset admin pasword
	if d.HasChange("password") {
		err := updateAdminPassword(ctx, d, cssV1Client, client)
		if err != nil {
			return err
		}
	}

	// update kibana
	if d.HasChange("kibana_public_access") {
		err := updateKibanaPublicAccess(ctx, d, cssV1Client)
		if err != nil {
			return err
		}
	}

	// update public_access
	if d.HasChange("public_access") {
		err := updatePublicAccess(ctx, d, cssV1Client)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateNodeConfig(ctx context.Context, d *schema.ResourceData,
	cssV1Client *cssv1.CssClient, conf *config.Config) error {
	addMasterNode := isAddNode(d, "master_node_config")
	addClientNode := isAddNode(d, "client_node_config")

	flavorList, err := getFlavorList(cssV1Client)
	if err != nil {
		return err
	}

	if addMasterNode {
		flavorId, err := flattenFlavorId(InstanceTypeEssMaster, d, flavorList)
		if err != nil {
			return err
		}

		bodyParams := buildAddMasterNodeParams(d, flavorId)
		err = addMastersOrClients(ctx, d, conf, InstanceTypeEssMaster, bodyParams)
		if err != nil {
			return err
		}
	} else if d.HasChange("master_node_config.0.volume") {
		return fmt.Errorf("ess-master node volume not supports to be updated")
	}

	if addClientNode {
		flavorId, err := flattenFlavorId(InstanceTypeEssClient, d, flavorList)
		if err != nil {
			return err
		}

		bodyParams := buildAddClientNodeParams(d, flavorId)
		err = addMastersOrClients(ctx, d, conf, InstanceTypeEssClient, bodyParams)
		if err != nil {
			return err
		}
	} else if d.HasChange("client_node_config.0.volume") {
		return fmt.Errorf("ess-client node volume not supports to be updated")
	}

	// update flavor
	flavorChanges := []string{
		"ess_node_config.0.flavor",
		"master_node_config.0.flavor",
		"client_node_config.0.flavor",
		"cold_node_config.0.flavor",
	}
	if d.HasChanges(flavorChanges...) {
		err = updateFlavor(ctx, d, flavorList, conf, addMasterNode, addClientNode)
		if err != nil {
			return err
		}
	}

	// extend instance number
	instanceNumChanges := []string{
		"ess_node_config.0.instance_number",
		"master_node_config.0.instance_number",
		"client_node_config.0.instance_number",
		"cold_node_config.0.instance_number",
		"expect_node_num",
	}
	if d.HasChanges(instanceNumChanges...) {
		err = extendInstanceNumber(ctx, d, conf, addMasterNode, addClientNode)
		if err != nil {
			return err
		}
	}

	// extend volume size
	instanceVolumeSizeChanges := []string{
		"ess_node_config.0.volume.0.size",
		"cold_node_config.0.volume.0.size",
		"node_config.0.volume.0.size",
	}
	if d.HasChanges(instanceVolumeSizeChanges...) {
		err = extendVolumeSize(ctx, d, conf)
		if err != nil {
			return err
		}
	}

	return nil
}

func isAddNode(d *schema.ResourceData, node string) bool {
	oldRaws, newRaws := d.GetChange(node)
	return len(oldRaws.([]interface{})) == 0 && len(newRaws.([]interface{})) == 1
}

func updateBackupStrategy(d *schema.ResourceData, cssV1Client *cssv1.CssClient) error {
	rawList := d.Get("backup_strategy").([]interface{})

	if len(rawList) == 0 {
		// stop auto backup strategy
		_, err := cssV1Client.CreateAutoCreatePolicy(&model.CreateAutoCreatePolicyRequest{
			ClusterId: d.Id(),
			Body: &model.SetRdsBackupCnfReq{
				Prefix:  utils.String("snapshot"),
				Period:  utils.String("00:00 GMT+08:00"),
				Keepday: utils.Int32(7),
				Enable:  "false",
			},
		})

		if err != nil {
			return fmt.Errorf("error updating backup strategy: %s", err)
		}
	} else {
		raw := rawList[0].(map[string]interface{})

		if d.HasChanges("backup_strategy.0.bucket", "backup_strategy.0.backup_path",
			"backup_strategy.0.agency") {
			// If obs is specified, update basic configurations
			_, err := cssV1Client.UpdateSnapshotSetting(&model.UpdateSnapshotSettingRequest{
				ClusterId: d.Id(),
				Body: &model.UpdateSnapshotSettingReq{
					Bucket:   raw["bucket"].(string),
					BasePath: raw["backup_path"].(string),
					Agency:   raw["agency"].(string),
				},
			})
			if err != nil {
				return fmt.Errorf("error modifying basic configurations of a cluster snapshot: %s", err)
			}
		}

		// check backup strategy, if the policy was disabled, we should enable it
		policy, err := cssV1Client.ShowAutoCreatePolicy(&model.ShowAutoCreatePolicyRequest{ClusterId: d.Id()})
		if err != nil {
			return fmt.Errorf("error extracting cluster backup_strategy, err: %s", err)
		}

		if utils.StringValue(policy.Enable) == "false" && raw["bucket"] == nil {
			// If obs is not specified,  create  basic configurations automatically
			_, err = cssV1Client.StartAutoSetting(&model.StartAutoSettingRequest{ClusterId: d.Id()})
			if err != nil {
				return fmt.Errorf("error enable snapshot function: %s", err)
			}
		}

		// update policy
		if d.HasChanges("backup_strategy.0.prefix", "backup_strategy.0.start_time",
			"backup_strategy.0.keep_days") {
			opts := &model.CreateAutoCreatePolicyRequest{
				ClusterId: d.Id(),
				Body: &model.SetRdsBackupCnfReq{
					Prefix:  utils.String(raw["prefix"].(string)),
					Period:  utils.String(raw["start_time"].(string)),
					Keepday: utils.Int32(int32(raw["keep_days"].(int))),
					Enable:  "true",
				},
			}
			_, err = cssV1Client.CreateAutoCreatePolicy(opts)
			if err != nil {
				return fmt.Errorf("error updating backup strategy: %s", err)
			}
		}
	}
	return nil
}

func updateVpcepEndpoint(ctx context.Context, d *schema.ResourceData, cssV1Client *cssv1.CssClient) error {
	o, n := d.GetChange("vpcep_endpoint")
	oValue := o.([]interface{})
	nValue := n.([]interface{})
	switch len(nValue) - len(oValue) {
	case -1: // delete vpc endpoint
		_, err := cssV1Client.StopVpecp(&model.StopVpecpRequest{ClusterId: d.Id()})
		if err != nil {
			return fmt.Errorf("error deleting the VPC endpoint of CSS cluster: %s, err: %s", d.Id(), err)
		}
		err = checkClusterOperationCompleted(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

	case 1: // start vpc endpoint
		_, err := cssV1Client.StartVpecp(&model.StartVpecpRequest{
			ClusterId: d.Id(),
			Body: &model.StartVpecpReq{
				EndpointWithDnsName: utils.Bool(d.Get("vpcep_endpoint.0.endpoint_with_dns_name").(bool)),
			},
		})
		if err != nil {
			return fmt.Errorf("error creating the VPC endpoint of CSS cluster: %s, err: %s", d.Id(), err)
		}
		err = checkClusterOperationCompleted(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

	case 0: // update vpc endpoint
		// update whitelist
		if d.HasChange("vpcep_endpoint.0.whitelist") {
			_, err := cssV1Client.UpdateVpcepWhitelist(&model.UpdateVpcepWhitelistRequest{
				ClusterId: d.Id(),
				Body: &model.UpdateVpcepWhitelistReq{
					VpcPermissions: utils.ExpandToStringList(d.Get("vpcep_endpoint.0.whitelist").([]interface{})),
				},
			})
			if err != nil {
				return fmt.Errorf("error updating the VPC endpoint whitelist of CSS cluster: %s, err: %s", d.Id(), err)
			}
		}
	}
	return nil
}

func updateKibanaPublicAccess(ctx context.Context, d *schema.ResourceData, cssV1Client *cssv1.CssClient) error {
	o, n := d.GetChange("kibana_public_access")
	oValue := o.([]interface{})
	nValue := n.([]interface{})

	switch len(nValue) - len(oValue) {
	case -1: // delete kibana_public_access
		_, err := cssV1Client.UpdateCloseKibana(&model.UpdateCloseKibanaRequest{
			ClusterId: d.Id(),
			Body: &model.CloseKibanaPublicReq{
				EipSize: utils.Int32(int32(utils.PathSearch("bandwidth", oValue[0], 0).(int))),
				ElbWhiteList: &model.StartKibanaPublicReqElbWhitelist{
					EnableWhiteList: utils.PathSearch("whitelist_enabled", oValue[0], false).(bool),
					WhiteList:       utils.PathSearch("whitelist", oValue[0], "").(string),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("error diabling the kibana public access of CSS cluster: %s, err: %s", d.Id(), err)
		}
		err = checkClusterOperationCompleted(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

	case 1: // enable kibana_public_access
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
			return fmt.Errorf("error enabling the kibana public access of CSS cluster: %s, err: %s", d.Id(), err)
		}
		err = checkClusterOperationCompleted(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

	case 0:
		// update bandwidth
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
				return fmt.Errorf("error modifing bandwidth of the kibana public access of CSS cluster: %s, err: %s", d.Id(), err)
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
					return fmt.Errorf("error disabing the whitelist of the kibana public access of CSS cluster: %s, err: %s", d.Id(), err)
				}
			} else {
				_, err := cssV1Client.UpdatePublicKibanaWhitelist(&model.UpdatePublicKibanaWhitelistRequest{
					ClusterId: d.Id(),
					Body: &model.UpdatePublicKibanaWhitelistReq{
						WhiteList: d.Get("kibana_public_access.0.whitelist").(string),
					},
				})
				if err != nil {
					return fmt.Errorf("error modifing whitelist of the kibana public access of CSS cluster: %s, err: %s", d.Id(), err)
				}
			}
		}
	}

	return nil
}

func updatePublicAccess(ctx context.Context, d *schema.ResourceData, cssV1Client *cssv1.CssClient) error {
	o, n := d.GetChange("public_access")
	oValue := o.([]interface{})
	nValue := n.([]interface{})

	switch len(nValue) - len(oValue) {
	case -1: // delete public_access
		_, err := cssV1Client.UpdateUnbindPublic(&model.UpdateUnbindPublicRequest{
			ClusterId: d.Id(),
			Body: &model.UnBindPublicReq{
				Eip: &model.UnBindPublicReqEipReq{
					BandWidth: &model.BindPublicReqEipBandWidth{
						Size: int32(utils.PathSearch("bandwidth", oValue[0], 0).(int)),
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("error diabling public access of CSS cluster: %s, err: %s", d.Id(), err)
		}
		err = checkClusterOperationCompleted(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

	case 1:
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
			return fmt.Errorf("error enabling public access of CSS cluster: %s, err: %s", d.Id(), err)
		}

		err = checkClusterOperationCompleted(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

		if whitelist, ok := d.GetOk("public_access.0.whitelist"); ok {
			_, err := cssV1Client.StartPublicWhitelist(&model.StartPublicWhitelistRequest{
				ClusterId: d.Id(),
				Body: &model.StartPublicWhitelistReq{
					WhiteList: whitelist.(string),
				},
			})
			if err != nil {
				return fmt.Errorf("error updating whitelist of public access of CSS cluster: %s, err: %s", d.Id(), err)
			}
		}

	case 0:
		// disable whitelist
		if d.HasChanges("public_access.0.whitelist", "public_access.0.whitelist_enabled") {
			if !d.Get("public_access.0.whitelist_enabled").(bool) {
				_, err := cssV1Client.StopPublicWhitelist(&model.StopPublicWhitelistRequest{ClusterId: d.Id()})
				if err != nil {
					return fmt.Errorf("error disabling whitelist of public access of CSS cluster: %s, err: %s", d.Id(), err)
				}
			} else {
				_, err := cssV1Client.StartPublicWhitelist(&model.StartPublicWhitelistRequest{
					ClusterId: d.Id(),
					Body: &model.StartPublicWhitelistReq{
						WhiteList: d.Get("public_access.0.whitelist").(string),
					},
				})
				if err != nil {
					return fmt.Errorf("error updating whitelist of public access of CSS cluster: %s, err: %s", d.Id(), err)
				}
			}
		}

		// update bandwidth
		if d.HasChange("public_access.0.bandwidth") {
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
				return fmt.Errorf("error disabling the whitelist of the kibana public access of CSS cluster: %s, err: %s", d.Id(), err)
			}
		}
	}

	return nil
}

func resourceCssClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.HcCssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	_, err = cssV1Client.DeleteCluster(&model.DeleteClusterRequest{ClusterId: d.Id()})
	if err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		err = ConvertExpectedHwSdkErrInto404Err(err, http.StatusForbidden, "CSS.0015", "")
		return common.CheckDeletedDiag(d, err, "error deleting the CSS cluster")
	}

	err = checkClusterDeleteResult(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("failed to check the result of deletion %s", err)
	}
	d.SetId("")
	return nil
}

func checkClusterCreateResult(ctx context.Context, cssV1Client *cssv1.CssClient, clusterId string,
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

func checkClusterDeleteResult(ctx context.Context, cssV1Client *cssv1.CssClient, clusterId string,
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
						var apiError ResponseError
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

func checkClusterOperationCompleted(ctx context.Context, cssV1Client *cssv1.CssClient, clusterId string,
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

	// actions --- the behaviors on a cluster
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

func updateCssTags(cssV1Client *cssv1.CssClient, id string, old, new map[string]interface{}) error {
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

func updateFlavor(ctx context.Context, d *schema.ResourceData,
	flavorsResp map[string]interface{}, conf *config.Config, addMasterNode, addClientNode bool) error {
	if d.HasChange("ess_node_config.0.flavor") {
		err := updateFlavorByType(ctx, InstanceTypeEss, d, flavorsResp, conf)
		if err != nil {
			return err
		}
	}
	if d.HasChange("master_node_config.0.flavor") {
		if !addMasterNode {
			err := updateFlavorByType(ctx, InstanceTypeEssMaster, d, flavorsResp, conf)
			if err != nil {
				return err
			}
		}
	}
	if d.HasChange("client_node_config.0.flavor") {
		if !addClientNode {
			err := updateFlavorByType(ctx, InstanceTypeEssClient, d, flavorsResp, conf)
			if err != nil {
				return err
			}
		}
	}
	if d.HasChange("cold_node_config.0.flavor") {
		err := updateFlavorByType(ctx, InstanceTypeEssCold, d, flavorsResp, conf)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateFlavorByType(ctx context.Context, nodeType string, d *schema.ResourceData,
	resp map[string]interface{}, conf *config.Config) error {
	region := conf.GetRegion(d)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return fmt.Errorf("error creating CSS V1 client: %s", err)
	}
	hcCssV1Client, err := conf.HcCssV1Client(region)
	if err != nil {
		return fmt.Errorf("error creating CSS V1 client: %s", err)
	}

	flavorId, err := flattenFlavorId(nodeType, d, resp)
	if err != nil {
		return err
	}

	updateFlavorHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/{types}/flavor"
	updateFlavorPath := cssV1Client.Endpoint + updateFlavorHttpUrl
	updateFlavorPath = strings.ReplaceAll(updateFlavorPath, "{project_id}", cssV1Client.ProjectID)
	updateFlavorPath = strings.ReplaceAll(updateFlavorPath, "{cluster_id}", d.Id())
	updateFlavorPath = strings.ReplaceAll(updateFlavorPath, "{types}", nodeType)

	updateFlavorOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	updateFlavorOpt.JSONBody = map[string]interface{}{
		"needCheckReplica": true,
		"newFlavorId":      flavorId,
		"isAutoPay":        1,
	}
	updateFlavorResp, err := cssV1Client.Request("POST", updateFlavorPath, &updateFlavorOpt)
	if err != nil {
		return fmt.Errorf("error updating CSS cluster flavor, cluster_id: %s, error: %s", d.Id(), err)
	}

	updateFlavorRespBody, err := utils.FlattenResponse(updateFlavorResp)
	if err != nil {
		return fmt.Errorf("error retrieving CSS cluster updating flavor response: %s", err)
	}

	orderId := utils.PathSearch("orderId", updateFlavorRespBody, "").(string)
	if orderId != "" {
		bssClient, err := conf.BssV2Client(region)
		if err != nil {
			return fmt.Errorf("error creating BSS v2 client: %s", err)
		}

		// If charging mode is PrePaid, wait for the order to be completed.
		err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	err = checkClusterOperationCompleted(ctx, hcCssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	return nil
}

func flattenFlavorId(nodeType string, d *schema.ResourceData, resp map[string]interface{}) (string, error) {
	version := d.Get("engine_version").(string)
	var flavorName string
	switch nodeType {
	case InstanceTypeEss:
		flavorName = d.Get("ess_node_config.0.flavor").(string)
	case InstanceTypeEssCold:
		flavorName = d.Get("cold_node_config.0.flavor").(string)
	case InstanceTypeEssMaster:
		flavorName = d.Get("master_node_config.0.flavor").(string)
	case InstanceTypeEssClient:
		flavorName = d.Get("master_node_config.0.flavor").(string)
	}
	findFlavorIdChar := fmt.Sprintf(
		`versions|[?type=='%s'&&version=='%s']|[0]|flavors|[?name=='%s']|[0]|flavor_id`,
		nodeType, version, flavorName)
	flavorId := utils.PathSearch(findFlavorIdChar, resp, "")
	if flavorId == "" {
		return "", fmt.Errorf("unable to find the ID of flavor(type: %s, version: %s, name: %s)",
			nodeType, version, flavorName)
	}
	return flavorId.(string), nil
}

func getFlavorList(cssV1Client *cssv1.CssClient) (map[string]interface{}, error) {
	flavorsResp, err := cssV1Client.ListFlavors(&model.ListFlavorsRequest{})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve CSS flavors: %s ", err)
	}

	getFlavorsRespJson, err := json.Marshal(flavorsResp)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(getFlavorsRespJson, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func extendInstanceNumber(ctx context.Context, d *schema.ResourceData, conf *config.Config,
	addMasterNode, addClientNode bool) error {
	if d.HasChange("ess_node_config.0.instance_number") {
		oldv, newv := d.GetChange("ess_node_config.0.instance_number")
		oldNum := oldv.(int)
		newNum := newv.(int)
		if newNum < oldNum {
			if d.Get("is_period").(bool) {
				return fmt.Errorf("instance shrinking operation is only supported in the PostPaid charging mode")
			}
			nodeIds := d.Get("ess_node_config.0.shrink_node_ids").([]interface{})
			err := updateShrinkInstance(ctx, d, conf, oldNum-newNum, InstanceTypeEss, nodeIds)
			if err != nil {
				return err
			}
		}

		if newNum > oldNum {
			bodyParams := buildUpdateExtendInstanceStorageBodyParams(InstanceTypeEss, newNum-oldNum, 0)

			err := updateExtendInstanceStorage(ctx, d, conf, bodyParams)
			if err != nil {
				return err
			}
		}
	}
	if d.HasChange("master_node_config.0.instance_number") {
		if !addMasterNode {
			oldv, newv := d.GetChange("master_node_config.0.instance_number")
			oldNum := oldv.(int)
			newNum := newv.(int)
			if newNum < oldNum {
				if d.Get("is_period").(bool) {
					return fmt.Errorf("instance shrinking operation is only supported in the PostPaid charging mode")
				}
				nodeIds := d.Get("master_node_config.0.shrink_node_ids").([]interface{})
				err := updateShrinkInstance(ctx, d, conf, oldNum-newNum, InstanceTypeEssMaster, nodeIds)
				if err != nil {
					return err
				}
			}

			if newNum > oldNum {
				bodyParams := buildUpdateExtendInstanceStorageBodyParams(InstanceTypeEssMaster, newNum-oldNum, 0)

				err := updateExtendInstanceStorage(ctx, d, conf, bodyParams)
				if err != nil {
					return err
				}
			}
		}
	}
	if d.HasChange("client_node_config.0.instance_number") {
		if !addClientNode {
			oldv, newv := d.GetChange("client_node_config.0.instance_number")
			oldNum := oldv.(int)
			newNum := newv.(int)
			if newNum < oldNum {
				if d.Get("is_period").(bool) {
					return fmt.Errorf("instance shrinking operation is only supported in the PostPaid charging mode")
				}
				nodeIds := d.Get("client_node_config.0.shrink_node_ids").([]interface{})
				err := updateShrinkInstance(ctx, d, conf, oldNum-newNum, InstanceTypeEssClient, nodeIds)
				if err != nil {
					return err
				}
			}

			if newNum > oldNum {
				bodyParams := buildUpdateExtendInstanceStorageBodyParams(InstanceTypeEssClient, newNum-oldNum, 0)

				err := updateExtendInstanceStorage(ctx, d, conf, bodyParams)
				if err != nil {
					return err
				}
			}
		}
	}
	if d.HasChange("cold_node_config.0.instance_number") {
		oldv, newv := d.GetChange("cold_node_config.0.instance_number")
		oldNum := oldv.(int)
		newNum := newv.(int)
		if newNum < oldNum {
			if d.Get("is_period").(bool) {
				return fmt.Errorf("instance shrinking operation is only supported in the PostPaid charging mode")
			}
			nodeIds := d.Get("cold_node_config.0.shrink_node_ids").([]interface{})
			err := updateShrinkInstance(ctx, d, conf, oldNum-newNum, InstanceTypeEssCold, nodeIds)
			if err != nil {
				return err
			}
		}

		if newNum > oldNum {
			bodyParams := buildUpdateExtendInstanceStorageBodyParams(InstanceTypeEssCold, newNum-oldNum, 0)

			err := updateExtendInstanceStorage(ctx, d, conf, bodyParams)
			if err != nil {
				return err
			}
		}
	}
	if d.HasChange("expect_node_num") {
		oldv, newv := d.GetChange("expect_node_num")
		nodesize := newv.(int) - oldv.(int)
		if nodesize < 0 {
			return fmt.Errorf("instance_number only supports to be extended")
		}
		bodyParams := buildUpdateExtendInstanceStorageBodyParams(InstanceTypeEss, nodesize, 0)

		err := updateExtendInstanceStorage(ctx, d, conf, bodyParams)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateExtendInstanceStorage(ctx context.Context, d *schema.ResourceData,
	conf *config.Config, bodyParams map[string]interface{}) error {
	region := conf.GetRegion(d)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return fmt.Errorf("error creating CSS V1 client: %s", err)
	}
	hcCssV1Client, err := conf.HcCssV1Client(region)
	if err != nil {
		return fmt.Errorf("error creating CSS V1 client: %s", err)
	}

	updateExtendHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/role_extend"
	updateExtendPath := cssV1Client.Endpoint + updateExtendHttpUrl
	updateExtendPath = strings.ReplaceAll(updateExtendPath, "{project_id}", cssV1Client.ProjectID)
	updateExtendPath = strings.ReplaceAll(updateExtendPath, "{cluster_id}", d.Id())

	updateExtendOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	updateExtendOpt.JSONBody = bodyParams
	updateExtendResp, err := cssV1Client.Request("POST", updateExtendPath, &updateExtendOpt)
	if err != nil {
		return fmt.Errorf("error updating CSS cluster extend, cluster_id: %s, error: %s", d.Id(), err)
	}

	updateExtendRespBody, err := utils.FlattenResponse(updateExtendResp)
	if err != nil {
		return fmt.Errorf("error retrieving CSS cluster updating extend response: %s", err)
	}

	orderId := utils.PathSearch("orderId", updateExtendRespBody, "").(string)
	if orderId != "" {
		bssClient, err := conf.BssV2Client(region)
		if err != nil {
			return fmt.Errorf("error creating BSS v2 client: %s", err)
		}

		// If charging mode is PrePaid, wait for the order to be completed.
		err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	err = checkClusterOperationCompleted(ctx, hcCssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	return nil
}

func buildUpdateExtendInstanceStorageBodyParams(nodeType string, nodesize, disksize int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"grow": []map[string]interface{}{
			{
				"type":     nodeType,
				"nodesize": nodesize,
				"disksize": disksize,
			},
		},
		"isAutoPay": 1,
	}
	return bodyParams
}

func extendVolumeSize(ctx context.Context, d *schema.ResourceData, conf *config.Config) error {
	if d.HasChange("ess_node_config.0.volume.0.size") {
		oldDisksize, newDisksize := d.GetChange("ess_node_config.0.volume.0.size")
		disksize := newDisksize.(int) - oldDisksize.(int)
		if disksize < 0 {
			return fmt.Errorf("volume size only supports to be extended")
		}

		bodyParams := buildUpdateExtendInstanceStorageBodyParams(InstanceTypeEss, 0, disksize)

		err := updateExtendInstanceStorage(ctx, d, conf, bodyParams)
		if err != nil {
			return err
		}
	}
	if d.HasChange("cold_node_config.0.volume.0.size") {
		oldDisksize, newDisksize := d.GetChange("cold_node_config.0.volume.0.size")
		disksize := newDisksize.(int) - oldDisksize.(int)
		if disksize < 0 {
			return fmt.Errorf("volume size only supports to be extended")
		}

		bodyParams := buildUpdateExtendInstanceStorageBodyParams(InstanceTypeEssCold, 0, disksize)

		err := updateExtendInstanceStorage(ctx, d, conf, bodyParams)
		if err != nil {
			return err
		}
	}
	if d.HasChange("node_config.0.volume.0.size") {
		oldDisksize, newDisksize := d.GetChange("node_config.0.volume.0.size")
		disksize := newDisksize.(int) - oldDisksize.(int)
		if disksize < 0 {
			return fmt.Errorf("volume size only supports to be extended")
		}

		bodyParams := buildUpdateExtendInstanceStorageBodyParams(InstanceTypeEss, 0, disksize)

		err := updateExtendInstanceStorage(ctx, d, conf, bodyParams)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildAddMasterNodeParams(d *schema.ResourceData, flavorId string) map[string]interface{} {
	// to do: is_auto_pay param need to add after api bug fixed
	bodyParams := map[string]interface{}{
		"type": map[string]interface{}{
			"flavor_ref":  flavorId,
			"node_size":   d.Get("master_node_config.0.instance_number"),
			"volume_type": d.Get("master_node_config.0.volume.0.volume_type"),
		},
	}

	return bodyParams
}

func buildAddClientNodeParams(d *schema.ResourceData, flavorId string) map[string]interface{} {
	// to do: is_auto_pay param need to add after api bug fixed
	bodyParams := map[string]interface{}{
		"type": map[string]interface{}{
			"flavor_ref":  flavorId,
			"node_size":   d.Get("client_node_config.0.instance_number"),
			"volume_type": d.Get("client_node_config.0.volume.0.volume_type"),
		},
	}
	return bodyParams
}

func addMastersOrClients(ctx context.Context, d *schema.ResourceData,
	conf *config.Config, nodeType string, bodyParams map[string]interface{}) error {
	region := conf.GetRegion(d)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return fmt.Errorf("error creating CSS V1 client: %s", err)
	}
	hcCssV1Client, err := conf.HcCssV1Client(region)
	if err != nil {
		return fmt.Errorf("error creating CSS V1 client: %s", err)
	}

	addNodeHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/type/{type}/independent"
	addNodePath := cssV1Client.Endpoint + addNodeHttpUrl
	addNodePath = strings.ReplaceAll(addNodePath, "{project_id}", cssV1Client.ProjectID)
	addNodePath = strings.ReplaceAll(addNodePath, "{cluster_id}", d.Id())
	addNodePath = strings.ReplaceAll(addNodePath, "{type}", nodeType)

	addNodeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	addNodeOpt.JSONBody = bodyParams
	addNodeResp, err := cssV1Client.Request("POST", addNodePath, &addNodeOpt)
	if err != nil {
		return fmt.Errorf("error add CSS cluster %s node, cluster_id: %s, error: %s", nodeType, d.Id(), err)
	}

	addNodeRespBody, err := utils.FlattenResponse(addNodeResp)
	if err != nil {
		return fmt.Errorf("error retrieving CSS cluster updating extend response: %s", err)
	}

	orderId := utils.PathSearch("orderId", addNodeRespBody, "").(string)
	if orderId != "" {
		bssClient, err := conf.BssV2Client(region)
		if err != nil {
			return fmt.Errorf("error creating BSS v2 client: %s", err)
		}

		// If charging mode is PrePaid, wait for the order to be completed.
		err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	err = checkClusterOperationCompleted(ctx, hcCssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	return nil
}

func updateSafeMode(ctx context.Context, d *schema.ResourceData,
	cssV1Client *cssv1.CssClient, client *golangsdk.ServiceClient) error {
	adminPassword := d.Get("password").(string)
	if d.Get("security_mode").(bool) && adminPassword == "" {
		return fmt.Errorf("administrator password is required when security mode changes")
	}
	updateSafeModeHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/mode/change"
	updateSafeModePath := client.Endpoint + updateSafeModeHttpUrl
	updateSafeModePath = strings.ReplaceAll(updateSafeModePath, "{project_id}", client.ProjectID)
	updateSafeModePath = strings.ReplaceAll(updateSafeModePath, "{cluster_id}", d.Id())

	updateSafeModeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	updateSafeModeOpt.JSONBody = buildUpdateSafeModeParams(d, adminPassword)

	_, err := client.Request("POST", updateSafeModePath, &updateSafeModeOpt)
	if err != nil {
		return fmt.Errorf("error updating CSS cluster security mode, cluster_id: %s, error: %s", d.Id(), err)
	}

	err = checkClusterOperationCompleted(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func buildUpdateSafeModeParams(d *schema.ResourceData, adminPassword string) map[string]interface{} {
	body := map[string]interface{}{
		"authorityEnable": d.Get("security_mode").(bool),
		"httpsEnable":     false,
	}
	if d.Get("security_mode").(bool) {
		body["adminPwd"] = adminPassword
		body["httpsEnable"] = d.Get("https_enabled").(bool)
	}

	return body
}

func updateAdminPassword(ctx context.Context, d *schema.ResourceData,
	cssV1Client *cssv1.CssClient, client *golangsdk.ServiceClient) error {
	adminPassword := d.Get("password").(string)
	if adminPassword == "" {
		return fmt.Errorf("administrator password is required when security mode changes")
	}
	updatePasswordHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/password/reset"
	updatePasswordPath := client.Endpoint + updatePasswordHttpUrl
	updatePasswordPath = strings.ReplaceAll(updatePasswordPath, "{project_id}", client.ProjectID)
	updatePasswordPath = strings.ReplaceAll(updatePasswordPath, "{cluster_id}", d.Id())

	updatePasswordOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	updatePasswordOpt.JSONBody = map[string]interface{}{
		"newpassword": adminPassword,
	}

	_, err := client.Request("POST", updatePasswordPath, &updatePasswordOpt)
	if err != nil {
		return fmt.Errorf("error resetting CSS cluster administrator password, cluster_id: %s, error: %s", d.Id(), err)
	}

	err = checkClusterOperationCompleted(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func updateSecurityGroup(ctx context.Context, d *schema.ResourceData,
	cssV1Client *cssv1.CssClient, client *golangsdk.ServiceClient) error {
	updateSecurityGroupHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/sg/change"
	updateSecurityGroupPath := client.Endpoint + updateSecurityGroupHttpUrl
	updateSecurityGroupPath = strings.ReplaceAll(updateSecurityGroupPath, "{project_id}", client.ProjectID)
	updateSecurityGroupPath = strings.ReplaceAll(updateSecurityGroupPath, "{cluster_id}", d.Id())

	updateSecurityGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	updateSecurityGroupOpt.JSONBody = map[string]interface{}{
		"security_group_ids": d.Get("security_group_id").(string),
	}

	_, err := client.Request("POST", updateSecurityGroupPath, &updateSecurityGroupOpt)
	if err != nil {
		return fmt.Errorf("error updating CSS cluster security group ID, cluster_id: %s, error: %s", d.Id(), err)
	}

	err = checkClusterOperationCompleted(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func updateShrinkInstance(ctx context.Context, d *schema.ResourceData, conf *config.Config,
	shrinkNodeSize int, nodeType string, shrinkNodeIds []interface{}) error {
	region := conf.GetRegion(d)
	client, err := conf.CssV1Client(region)
	if err != nil {
		return fmt.Errorf("error creating CSS V1 client: %s", err)
	}
	hcClient, err := conf.HcCssV1Client(region)
	if err != nil {
		return fmt.Errorf("error creating CSS V1 client: %s", err)
	}

	if len(shrinkNodeIds) == 0 {
		bodyParams := buildUpdateShrinkNodeByTypeBodyParams(nodeType, shrinkNodeSize)
		err := updateShrinkInstanceNodeByType(ctx, d, hcClient, client, bodyParams)
		if err != nil {
			return err
		}
	} else {
		if shrinkNodeSize != len(shrinkNodeIds) {
			return fmt.Errorf("instance_number changing number is inconsistent with the length"+
				" of shrink_node_ids, node type: %s", nodeType)
		}
		err := updateShrinkInstanceNodeById(ctx, d, hcClient, client, shrinkNodeIds)
		if err != nil {
			return err
		}
	}

	return nil
}

func buildUpdateShrinkNodeByTypeBodyParams(nodeType string, nodesize int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"shrink": []map[string]interface{}{
			{
				"type":           nodeType,
				"reducedNodeNum": nodesize,
			},
		},
	}
	return bodyParams
}

func updateShrinkInstanceNodeByType(ctx context.Context, d *schema.ResourceData, hcClient *cssv1.CssClient,
	client *golangsdk.ServiceClient, bodyParams map[string]interface{}) error {
	shrinkNodeHttpUrl := "v1.0/extend/{project_id}/clusters/{cluster_id}/role/shrink"
	shrinkNodePath := client.Endpoint + shrinkNodeHttpUrl
	shrinkNodePath = strings.ReplaceAll(shrinkNodePath, "{project_id}", client.ProjectID)
	shrinkNodePath = strings.ReplaceAll(shrinkNodePath, "{cluster_id}", d.Id())

	shrinkNodeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	shrinkNodeOpt.JSONBody = bodyParams

	_, err := client.Request("POST", shrinkNodePath, &shrinkNodeOpt)
	if err != nil {
		return fmt.Errorf("error shrinking CSS cluster node by type, cluster_id: %s, error: %s", d.Id(), err)
	}

	err = checkClusterOperationCompleted(ctx, hcClient, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func updateShrinkInstanceNodeById(ctx context.Context, d *schema.ResourceData, hcClient *cssv1.CssClient,
	client *golangsdk.ServiceClient, nodeIds []interface{}) error {
	shrinkNodeHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/node/offline"
	shrinkNodePath := client.Endpoint + shrinkNodeHttpUrl
	shrinkNodePath = strings.ReplaceAll(shrinkNodePath, "{project_id}", client.ProjectID)
	shrinkNodePath = strings.ReplaceAll(shrinkNodePath, "{cluster_id}", d.Id())

	shrinkNodeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	shrinkNodeOpt.JSONBody = map[string]interface{}{
		"migrate_data": true,
		"shrinkNodes":  nodeIds,
	}
	_, err := client.Request("POST", shrinkNodePath, &shrinkNodeOpt)
	if err != nil {
		return fmt.Errorf("error shrinking CSS cluster node by ID, cluster_id: %s, error: %s", d.Id(), err)
	}

	err = checkClusterOperationCompleted(ctx, hcClient, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	return nil
}

type ResponseError struct {
	ErrorCode string `json:"errCode"`
	ErrorMsg  string `json:"externalMessage"`
}
