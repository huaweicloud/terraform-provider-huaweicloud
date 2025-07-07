package css

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

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

var cssClusterSchema = map[string]*schema.Schema{
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
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
		Default:  "elasticsearch",
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
		Type:             schema.TypeString,
		Optional:         true,
		RequiredWith:     []string{"vpc_id", "subnet_id", "security_group_id"},
		Computed:         true,
		DiffSuppressFunc: utils.SuppressStringSepratedByCommaDiffs,
		Description:      "schema: Required",
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
					Type:             schema.TypeString,
					Optional:         true,
					DiffSuppressFunc: utils.SuppressStringSepratedByCommaDiffs,
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
					Type:             schema.TypeString,
					DiffSuppressFunc: utils.SuppressStringSepratedByCommaDiffs,
					Optional:         true,
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
}

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

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(clusterNonUpdatableParams, cssClusterSchema),
			config.MergeDefaultTags(),
		),

		Schema: cssClusterSchema,
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
	v1client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	v2client, err := conf.CssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V2 client: %s", err)
	}

	createClusterHttpUrl := "v2.0/{project_id}/clusters"
	createClusterPath := v2client.Endpoint + createClusterHttpUrl
	createClusterPath = strings.ReplaceAll(createClusterPath, "{project_id}", v2client.ProjectID)

	createClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	bodyParams, paramErr := buildClusterCreateParameters(d, conf)
	if paramErr != nil {
		return diag.FromErr(paramErr)
	}
	createClusterOpt.JSONBody = utils.RemoveNil(bodyParams)
	createClusterResp, err := v2client.Request("POST", createClusterPath, &createClusterOpt)
	if err != nil {
		return diag.Errorf("error creating CSS cluster, err: %s", err)
	}

	createClusterRespBody, err := utils.FlattenResponse(createClusterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterId := utils.PathSearch("cluster.id", createClusterRespBody, "").(string)
	orderId := utils.PathSearch("orderId", createClusterRespBody, "").(string)

	if orderId == "" {
		if clusterId == "" {
			return diag.Errorf("error creating CSS cluster: id is not found in API response,%#v", createClusterRespBody)
		}
		d.SetId(clusterId)
	} else {
		bssClient, err := conf.BssV2Client(conf.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}

		// 1. If charging mode is PrePaid, wait for the order to be completed.
		err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		// 2. get the resource ID, must be after order success
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resourceId)
	}

	createResultErr := checkClusterCreateResult(ctx, v1client, d.Id(), d.Timeout(schema.TimeoutCreate))
	if createResultErr != nil {
		return diag.FromErr(createResultErr)
	}

	return resourceCssClusterRead(ctx, d, meta)
}

func buildClusterCreateParameters(d *schema.ResourceData, conf *config.Config) (map[string]interface{}, error) {
	cluster := map[string]interface{}{
		"name": d.Get("name").(string),
		"datastore": map[string]interface{}{
			"type":    d.Get("engine_type"),
			"version": d.Get("engine_version"),
		},
		"enterprise_project_id": utils.ValueIgnoreEmpty(conf.GetEnterpriseProjectID(d)),
		"tags":                  utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
		"backupStrategy":        buildCssClusterCreateBackupStrategy(d.Get("backup_strategy").([]interface{})),
	}

	roles := make([]interface{}, 0)
	if ess, ok := d.GetOk("ess_node_config"); ok {
		essNode := ess.([]interface{})[0].(map[string]interface{})
		roles = append(roles, buildCreateClusterRole(essNode, InstanceTypeEss))

		if v, ok := d.GetOk("master_node_config"); ok {
			node := v.([]interface{})[0].(map[string]interface{})
			roles = append(roles, buildCreateClusterRole(node, InstanceTypeEssMaster))
		}
		if v, ok := d.GetOk("client_node_config"); ok {
			node := v.([]interface{})[0].(map[string]interface{})
			roles = append(roles, buildCreateClusterRole(node, InstanceTypeEssClient))
		}
		if v, ok := d.GetOk("cold_node_config"); ok {
			node := v.([]interface{})[0].(map[string]interface{})
			roles = append(roles, buildCreateClusterRole(node, InstanceTypeEssCold))
		}

		cluster["availability_zone"] = utils.StringIgnoreEmpty(d.Get("availability_zone").(string))
		cluster["nics"] = map[string]interface{}{
			"vpcId":           d.Get("vpc_id"),
			"netId":           d.Get("subnet_id"),
			"securityGroupId": d.Get("security_group_id"),
		}
		cluster["roles"] = roles
	} else {
		// Compatible with previous version
		cluster["availability_zone"] = utils.StringIgnoreEmpty(d.Get("node_config.0.availability_zone").(string))
		cluster["nics"] = map[string]interface{}{
			"vpcId":           d.Get("node_config.0.network_info.0.vpc_id"),
			"netId":           d.Get("node_config.0.network_info.0.subnet_id"),
			"securityGroupId": d.Get("node_config.0.network_info.0.security_group_id"),
		}

		// add ess role
		cluster["roles"] = append(roles, map[string]interface{}{
			"flavorRef":   d.Get("node_config.0.flavor"),
			"type":        InstanceTypeEss,
			"instanceNum": d.Get("expect_node_num"),
			"volume": map[string]interface{}{
				"size":        d.Get("node_config.0.volume.0.size"),
				"volume_type": d.Get("node_config.0.volume.0.volume_type"),
			},
		})
	}

	securityMode := d.Get("security_mode").(bool)
	if securityMode {
		adminPassword := d.Get("password").(string)
		if adminPassword == "" {
			return nil, fmt.Errorf("administrator password is required in security mode")
		}
		cluster["authorityEnable"] = true
		cluster["adminPwd"] = adminPassword
		cluster["httpsEnable"] = d.Get("https_enabled")
	}

	if _, ok := d.GetOk("vpcep_endpoint"); ok {
		loadBalance := map[string]interface{}{
			"endpointWithDnsName": d.Get("vpcep_endpoint.0.endpoint_with_dns_name"),
		}

		vpcPermissions := utils.ExpandToStringList(d.Get("vpcep_endpoint.0.whitelist").([]interface{}))
		if len(vpcPermissions) > 0 {
			loadBalance["vpcPermissions"] = vpcPermissions
		}
		cluster["loadBalance"] = loadBalance
	}

	if _, ok := d.GetOk("kibana_public_access"); ok {
		whitelist, ok := d.GetOk("kibana_public_access.0.whitelist")
		cluster["publicKibanaReq"] = map[string]interface{}{
			"eipSize": d.Get("kibana_public_access.0.bandwidth"),
			"elbWhiteList": map[string]interface{}{
				"enableWhiteList": ok,
				"whiteList":       whitelist,
			},
		}
	}

	if _, ok := d.GetOk("public_access"); ok {
		whitelist, ok := d.GetOk("public_access.0.whitelist")
		cluster["publicIPReq"] = map[string]interface{}{
			"publicBindType": "auto_assign",
			"eip": map[string]interface{}{
				"bandWidth": map[string]interface{}{
					"size": d.Get("public_access.0.bandwidth"),
				},
			},
			"elbWhiteListReq": map[string]interface{}{
				"enableWhiteList": ok,
				"whiteList":       whitelist,
			},
		}
	}

	if payModel, ok := d.GetOk("period_unit"); ok && d.Get("charging_mode").(string) != "postPaid" {
		payInfo := map[string]interface{}{
			"period":    d.Get("period"),
			"isAutoPay": 1,
		}

		if payModel == "month" {
			payInfo["payModel"] = 2
		} else {
			payInfo["payModel"] = 3
		}

		if d.Get("auto_renew").(string) == "true" {
			payInfo["isAutoRenew"] = 1
		}
		cluster["payInfo"] = payInfo
	}

	bodyParams := map[string]interface{}{
		"cluster": cluster,
	}
	return bodyParams, nil
}

func buildCreateClusterRole(node map[string]interface{}, nodeType string) map[string]interface{} {
	roleBoby := map[string]interface{}{
		"flavorRef":   node["flavor"],
		"type":        nodeType,
		"instanceNum": node["instance_number"],
	}

	// Ess node and cold node support local disk. The volume value is empty.
	// Master node and client node do not support local volume. The volume value is required.
	if volumes := node["volume"].([]interface{}); len(volumes) > 0 {
		volume := volumes[0].(map[string]interface{})
		roleBoby["volume"] = map[string]interface{}{
			"size":        volume["size"],
			"volume_type": volume["volume_type"],
		}
	}
	return roleBoby
}

func buildCssClusterCreateBackupStrategy(backupRaw []interface{}) map[string]interface{} {
	if len(backupRaw) == 0 {
		return nil
	}
	raw := backupRaw[0].(map[string]interface{})
	return map[string]interface{}{
		"prefix":   raw["prefix"],
		"period":   raw["start_time"],
		"keepday":  raw["keep_days"],
		"bucket":   utils.ValueIgnoreEmpty(raw["bucket"]),
		"basePath": utils.ValueIgnoreEmpty(raw["backup_path"]),
		"agency":   utils.ValueIgnoreEmpty(raw["agency"]),
	}
}

func resourceCssClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	v1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	clusterDetail, err := getClusterDetails(v1Client, d.Id())
	if err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015"), "error getting CSS cluster")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", clusterDetail, nil)),
		d.Set("engine_type", utils.PathSearch("datastore.type", clusterDetail, nil)),
		d.Set("engine_version", utils.PathSearch("datastore.version", clusterDetail, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterpriseProjectId", clusterDetail, nil)),
		d.Set("vpc_id", utils.PathSearch("vpcId", clusterDetail, nil)),
		d.Set("subnet_id", utils.PathSearch("subnetId", clusterDetail, nil)),
		d.Set("security_group_id", utils.PathSearch("securityGroupId", clusterDetail, nil)),
		d.Set("nodes", flattenClusterNodes(utils.PathSearch("instances", clusterDetail, make([]interface{}, 0)).([]interface{}))),
		setNodeConfigsAndAzToState(clusterDetail, d),
		d.Set("vpcep_ip", utils.PathSearch("vpcepIp", clusterDetail, nil)),
		d.Set("kibana_public_access", flattenKibana(utils.PathSearch("publicKibanaResp", clusterDetail, nil))),
		d.Set("public_access", flattenPublicAccess(clusterDetail)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", clusterDetail, nil))),
		d.Set("created", utils.PathSearch("created", clusterDetail, nil)),
		d.Set("endpoint", utils.PathSearch("endpoint", clusterDetail, nil)),
		d.Set("status", utils.PathSearch("status", clusterDetail, nil)),
		d.Set("security_mode", utils.PathSearch("authorityEnable", clusterDetail, false)),
		d.Set("https_enabled", utils.PathSearch("httpsEnable", clusterDetail, nil)),
		d.Set("created_at", utils.PathSearch("created", clusterDetail, nil)),
		d.Set("updated_at", utils.PathSearch("updated", clusterDetail, nil)),
		d.Set("bandwidth_resource_id", utils.PathSearch("bandwidthResourceId", clusterDetail, nil)),
		d.Set("is_period", utils.PathSearch("period", clusterDetail, nil)),
		d.Set("backup_available", utils.PathSearch("backupAvailable", clusterDetail, nil)),
		d.Set("disk_encrypted", utils.PathSearch("diskEncrypted", clusterDetail, nil)),
	)

	getVpcepConnectionRespBody, err := getVpcepConnection(d, v1Client)
	if err != nil {
		return diag.Errorf("error extracting vpcep connection: %s", err)
	}
	mErr = multierror.Append(mErr,
		d.Set("vpcep_endpoint_id", utils.PathSearch("connections|[0].id", getVpcepConnectionRespBody, nil)),
	)

	getAutoCreatePolicyRespBody, err := getAutoCreateBackupStrategy(d, v1Client)
	if err != nil {
		return diag.Errorf("error extracting cluster: backup_strategy, err: %s", err)
	}
	policyEnable := utils.PathSearch("enable", getAutoCreatePolicyRespBody, nil)
	if policyEnable == "true" {
		strategy := []map[string]interface{}{
			{
				"prefix":      utils.PathSearch("prefix", getAutoCreatePolicyRespBody, nil),
				"start_time":  utils.PathSearch("period", getAutoCreatePolicyRespBody, nil),
				"keep_days":   int(utils.PathSearch("keepday", getAutoCreatePolicyRespBody, float64(0)).(float64)),
				"bucket":      utils.PathSearch("bucket", getAutoCreatePolicyRespBody, nil),
				"backup_path": utils.PathSearch("basePath", getAutoCreatePolicyRespBody, nil),
				"agency":      utils.PathSearch("agency", getAutoCreatePolicyRespBody, nil),
			},
		}
		mErr = multierror.Append(mErr, d.Set("backup_strategy", strategy))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func getClusterDetails(client *golangsdk.ServiceClient, clusterID string) (interface{}, error) {
	getClusterDetailsHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}"
	getClusterDetailsPath := client.Endpoint + getClusterDetailsHttpUrl
	getClusterDetailsPath = strings.ReplaceAll(getClusterDetailsPath, "{project_id}", client.ProjectID)
	getClusterDetailsPath = strings.ReplaceAll(getClusterDetailsPath, "{cluster_id}", clusterID)

	getClusterDetailsPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getClusterDetailsResp, err := client.Request("GET", getClusterDetailsPath, &getClusterDetailsPathOpt)
	if err != nil {
		return getClusterDetailsResp, err
	}

	return utils.FlattenResponse(getClusterDetailsResp)
}

func flattenClusterNodes(nodes []interface{}) []interface{} {
	if len(nodes) == 0 {
		return nil
	}

	rst := make([]interface{}, len(nodes))
	for i, v := range nodes {
		rst[i] = map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"type":              utils.PathSearch("type", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"availability_zone": utils.PathSearch("azCode", v, nil),
			"status":            utils.PathSearch("status", v, nil),
			"spec_code":         utils.PathSearch("specCode", v, nil),
			"ip":                utils.PathSearch("ip", v, nil),
			"resource_id":       utils.PathSearch("resourceId", v, nil),
		}
	}
	return rst
}

func getAutoCreateBackupStrategy(d *schema.ResourceData, client *golangsdk.ServiceClient) (interface{}, error) {
	getAutoCreatePolicyHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/policy"
	getAutoCreatePolicyPath := client.Endpoint + getAutoCreatePolicyHttpUrl
	getAutoCreatePolicyPath = strings.ReplaceAll(getAutoCreatePolicyPath, "{project_id}", client.ProjectID)
	getAutoCreatePolicyPath = strings.ReplaceAll(getAutoCreatePolicyPath, "{cluster_id}", d.Id())

	getAutoCreatePolicyPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getClusterBackupStrategyResp, err := client.Request("GET", getAutoCreatePolicyPath, &getAutoCreatePolicyPathOpt)
	if err != nil {
		return getClusterBackupStrategyResp, err
	}

	return utils.FlattenResponse(getClusterBackupStrategyResp)
}

func flattenKibana(publicKibana interface{}) []interface{} {
	elbWhiteListResp := utils.PathSearch("elbWhiteListResp", publicKibana, nil)
	if publicKibana == nil || elbWhiteListResp == nil {
		return nil
	}

	result := map[string]interface{}{
		"bandwidth":         int(utils.PathSearch("eipSize", publicKibana, float64(0)).(float64)),
		"whitelist_enabled": utils.PathSearch("enableWhiteList", elbWhiteListResp, nil),
		"whitelist":         utils.PathSearch("whiteList", elbWhiteListResp, nil),
		"public_ip":         utils.PathSearch("publicKibanaIp", publicKibana, nil),
	}
	return []interface{}{result}
}

func flattenPublicAccess(clusterDetail interface{}) []interface{} {
	elbWhiteList := utils.PathSearch("elbWhiteList", clusterDetail, nil)
	bandwidth := int(utils.PathSearch("bandwidthSize", clusterDetail, float64(0)).(float64))
	publicIp := utils.PathSearch("publicIp", clusterDetail, nil)
	if elbWhiteList == nil || publicIp == nil {
		return nil
	}

	result := map[string]interface{}{
		"bandwidth":         bandwidth,
		"whitelist_enabled": utils.PathSearch("enableWhiteList", elbWhiteList, nil),
		"whitelist":         utils.PathSearch("whiteList", elbWhiteList, nil),
		"public_ip":         publicIp,
	}
	return []interface{}{result}
}

func getVpcepConnection(d *schema.ResourceData, client *golangsdk.ServiceClient) (interface{}, error) {
	getVpcepConnectionHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/vpcepservice/connections"
	getVpcepConnectionPath := client.Endpoint + getVpcepConnectionHttpUrl
	getVpcepConnectionPath = strings.ReplaceAll(getVpcepConnectionPath, "{project_id}", client.ProjectID)
	getVpcepConnectionPath = strings.ReplaceAll(getVpcepConnectionPath, "{cluster_id}", d.Id())

	getVpcepConnectionPathOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getVpcepConnectionResp, err := client.Request("GET", getVpcepConnectionPath, &getVpcepConnectionPathOpt)
	if err != nil {
		// CSS.5182 : The VPC endpoint service is not enabled.
		if hasErrorCode(err, "CSS.5182") {
			return getVpcepConnectionResp, nil
		}
		return getVpcepConnectionResp, err
	}

	return utils.FlattenResponse(getVpcepConnectionResp)
}

func hasErrorCode(err error, expectCode string) bool {
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var response interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &response); jsonErr == nil {
			errorCode, parseErr := jmespath.Search("errCode", response)
			if parseErr != nil {
				log.Printf("[WARN] failed to parse errCode from response body: %s", parseErr)
			}

			if errorCode == expectCode {
				return true
			}
		}
	}

	return false
}

func setNodeConfigsAndAzToState(clusterDetail interface{}, d *schema.ResourceData) error {
	instances := utils.PathSearch("instances", clusterDetail, make([]interface{}, 0)).([]interface{})
	if len(instances) == 0 {
		return nil
	}
	nodeConfigMap, azArray := getNodeConfigMapAndAzArray(instances)
	azArray = utils.RemoveDuplicateElem(azArray)
	az := strings.Join(azArray, ",")
	mErr := multierror.Append(
		d.Set("availability_zone", az),
	)
	for k, v := range nodeConfigMap {
		switch k {
		case InstanceTypeEss:
			// old version nodeConfig, NO volume return, so get from state
			volumeSize := utils.PathSearch("volume|[0]|size", v, 0).(int)
			volumeType := utils.PathSearch("volume|[0]|volume_type", v, "").(string)
			nodeConfig := map[string]interface{}{
				"flavor":            v["flavor"],
				"availability_zone": az,
				"volume": []interface{}{map[string]interface{}{
					"size":        volumeSize,
					"volume_type": volumeType,
				}},
				"network_info": []interface{}{map[string]interface{}{
					"vpc_id":            utils.PathSearch("vpcId", clusterDetail, nil),
					"subnet_id":         utils.PathSearch("subnetId", clusterDetail, nil),
					"security_group_id": utils.PathSearch("securityGroupId", clusterDetail, nil),
				}},
			}

			// NO volume return, so get from state
			if volumeSize == 0 && volumeType == "" {
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
			volumeSize := utils.PathSearch("volume|[0]|size", v, 0).(int)
			volumeType := utils.PathSearch("volume|[0]|volume_type", v, "").(string)
			if volumeSize == 0 && volumeType == "" {
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
			volumeSize := utils.PathSearch("volume|[0]|size", v, 0).(int)
			volumeType := utils.PathSearch("volume|[0]|volume_type", v, "").(string)
			if volumeSize == 0 && volumeType == "" {
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
			volumeSize := utils.PathSearch("volume|[0]|size", v, 0).(int)
			volumeType := utils.PathSearch("volume|[0]|volume_type", v, "").(string)
			if volumeSize == 0 && volumeType == "" {
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

func getNodeConfigMapAndAzArray(instances []interface{}) (map[string]map[string]interface{}, []string) {
	nodeConfigMap := make(map[string]map[string]interface{})
	var azArray []string
	for _, v := range instances {
		azArray = append(azArray, utils.PathSearch("azCode", v, "").(string))

		nodeType := utils.PathSearch("type", v, "").(string)
		if node, ok := nodeConfigMap[nodeType]; ok {
			node["instance_number"] = node["instance_number"].(int) + 1
		} else {
			nodeConfigMap[nodeType] = map[string]interface{}{
				"flavor":          utils.PathSearch("specCode", v, nil),
				"instance_number": 1,
				"volume": []interface{}{map[string]interface{}{
					"size":        int(utils.PathSearch("volume.size", v, float64(0)).(float64)),
					"volume_type": utils.PathSearch("volume.type", v, nil),
				}},
				"shrink_node_ids": nil,
			}
		}
	}

	return nodeConfigMap, azArray
}

func resourceCssClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	clusterId := d.Id()

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
		err := updateNodeConfig(ctx, d, client, cfg)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update backup strategy
	if d.HasChange("backup_strategy") {
		err = updateBackupStrategy(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(client, d, "css-cluster", clusterId)
		if tagErr != nil {
			return diag.Errorf("error updating tags of CSS cluster:%s, err:%s", clusterId, tagErr)
		}
	}

	// update vpc endpoint
	if d.HasChange("vpcep_endpoint") {
		err = updateVpcepEndpoint(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update security mode
	if d.HasChange("security_mode") {
		err = updateSafeMode(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update in safe mode
	if d.Get("security_mode").(bool) {
		err = updateInSafeMode(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update security group ID
	if d.HasChange("security_group_id") {
		err = updateSecurityGroup(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("charging_mode", "auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}

		err = updateChangingModeOrAutoRenew(ctx, d, client, bssClient)
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

func updateInSafeMode(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	// reset admin pasword
	if d.HasChange("password") {
		err := updateAdminPassword(ctx, d, client)
		if err != nil {
			return err
		}
	}

	// update kibana
	if d.HasChange("kibana_public_access") {
		err := updateKibanaPublicAccess(ctx, d, client)
		if err != nil {
			return err
		}
	}

	// update public_access
	if d.HasChange("public_access") {
		err := updatePublicAccess(ctx, d, client)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateNodeConfig(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, conf *config.Config) error {
	addMasterNode := isAddNode(d, "master_node_config")
	addClientNode := isAddNode(d, "client_node_config")

	flavorList, err := getFlavorList(client)
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

func updateBackupStrategy(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	rawList := d.Get("backup_strategy").([]interface{})

	if len(rawList) == 0 {
		autoPolicyBodyParams := map[string]interface{}{
			"prefix":  "snapshot",
			"period":  "00:00 GMT+08:00",
			"keepday": 7,
			"enable":  "false",
		}
		err := createAutoCreatePolicy(client, d.Id(), autoPolicyBodyParams)
		if err != nil {
			return fmt.Errorf("error updating backup strategy: %s", err)
		}
	} else {
		raw := rawList[0].(map[string]interface{})
		if d.HasChanges("backup_strategy.0.bucket", "backup_strategy.0.backup_path", "backup_strategy.0.agency") {
			// If obs is specified, update basic configurations
			err := updateSnapshotSetting(client, d.Id(), raw)
			if err != nil {
				return fmt.Errorf("error modifying basic configurations of a cluster snapshot: %s", err)
			}
		}

		// check backup strategy, if the policy was disabled, we should enable it
		getAutoCreatePolicyRespBody, err := getAutoCreateBackupStrategy(d, client)
		if err != nil {
			return fmt.Errorf("error extracting cluster backup_strategy, err: %s", err)
		}
		policyEnable := utils.PathSearch("enable", getAutoCreatePolicyRespBody, "").(string)
		if policyEnable == "false" && raw["bucket"] == nil {
			// If obs is not specified,  create  basic configurations automatically
			err := startSnapshotAutoSetting(client, d.Id())
			if err != nil {
				return fmt.Errorf("error enable snapshot function: %s", err)
			}
		}

		// update policy
		if d.HasChanges("backup_strategy.0.prefix", "backup_strategy.0.start_time",
			"backup_strategy.0.keep_days") {
			autoPolicyBodyParams := map[string]interface{}{
				"prefix":  raw["prefix"],
				"period":  raw["start_time"],
				"keepday": raw["keep_days"],
				"enable":  "true",
			}
			err := createAutoCreatePolicy(client, d.Id(), autoPolicyBodyParams)
			if err != nil {
				return fmt.Errorf("error updating backup strategy: %s", err)
			}
		}
	}
	return nil
}

func createAutoCreatePolicy(client *golangsdk.ServiceClient, clusterId string, bodyParams map[string]interface{}) error {
	createAutoCreatePolicyHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/policy"
	createAutoCreatePolicyPath := client.Endpoint + createAutoCreatePolicyHttpUrl
	createAutoCreatePolicyPath = strings.ReplaceAll(createAutoCreatePolicyPath, "{project_id}", client.ProjectID)
	createAutoCreatePolicyPath = strings.ReplaceAll(createAutoCreatePolicyPath, "{cluster_id}", clusterId)

	createAutoCreatePolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         bodyParams,
	}

	_, err := client.Request("POST", createAutoCreatePolicyPath, &createAutoCreatePolicyOpt)
	if err != nil {
		return err
	}

	return nil
}

func updateSnapshotSetting(client *golangsdk.ServiceClient, clusterId string, raw map[string]interface{}) error {
	updateSnapshotSettingHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/setting"
	updateSnapshotSettingPath := client.Endpoint + updateSnapshotSettingHttpUrl
	updateSnapshotSettingPath = strings.ReplaceAll(updateSnapshotSettingPath, "{project_id}", client.ProjectID)
	updateSnapshotSettingPath = strings.ReplaceAll(updateSnapshotSettingPath, "{cluster_id}", clusterId)

	updateSnapshotSettingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"bucket":   raw["bucket"],
			"basePath": raw["backup_path"],
			"agency":   raw["agency"],
		},
	}

	_, err := client.Request("POST", updateSnapshotSettingPath, &updateSnapshotSettingOpt)
	if err != nil {
		return err
	}

	return nil
}

func startSnapshotAutoSetting(client *golangsdk.ServiceClient, clusterId string) error {
	startSnapshotAutoSettingHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/index_snapshot/auto_setting"
	startSnapshotAutoSettingPath := client.Endpoint + startSnapshotAutoSettingHttpUrl
	startSnapshotAutoSettingPath = strings.ReplaceAll(startSnapshotAutoSettingPath, "{project_id}", client.ProjectID)
	startSnapshotAutoSettingPath = strings.ReplaceAll(startSnapshotAutoSettingPath, "{cluster_id}", clusterId)

	startSnapshotAutoSettingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err := client.Request("POST", startSnapshotAutoSettingPath, &startSnapshotAutoSettingOpt)
	if err != nil {
		return err
	}

	return nil
}

func updateVpcepEndpoint(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	o, n := d.GetChange("vpcep_endpoint")
	oValue := o.([]interface{})
	nValue := n.([]interface{})
	clusterId := d.Id()
	switch len(nValue) - len(oValue) {
	case -1: // delete vpc endpoint
		err := stopVpcep(client, clusterId)
		if err != nil {
			return fmt.Errorf("error deleting the VPC endpoint of CSS cluster: %s, err: %s", clusterId, err)
		}
		err = checkClusterOperationResult(ctx, client, clusterId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

	case 1: // start vpc endpoint
		startVpcepBodyParams := map[string]interface{}{
			"endpointWithDnsName": d.Get("vpcep_endpoint.0.endpoint_with_dns_name"),
		}
		err := startVpcep(client, clusterId, startVpcepBodyParams)
		if err != nil {
			return fmt.Errorf("error creating the VPC endpoint of CSS cluster: %s, err: %s", clusterId, err)
		}
		err = checkClusterOperationResult(ctx, client, clusterId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

	case 0: // update vpc endpoint
		// update whitelist
		if d.HasChange("vpcep_endpoint.0.whitelist") {
			updateVpcepWhitelistBodyParams := map[string]interface{}{
				"vpcPermissions": d.Get("vpcep_endpoint.0.whitelist"),
			}
			err := updateVpcepWhitelist(client, clusterId, updateVpcepWhitelistBodyParams)
			if err != nil {
				return fmt.Errorf("error updating the VPC endpoint whitelist of CSS cluster: %s, err: %s", clusterId, err)
			}
		}
	}
	return nil
}

func stopVpcep(client *golangsdk.ServiceClient, clusterId string) error {
	stopVpcepHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/vpcepservice/close"
	stopVpcepPath := client.Endpoint + stopVpcepHttpUrl
	stopVpcepPath = strings.ReplaceAll(stopVpcepPath, "{project_id}", client.ProjectID)
	stopVpcepPath = strings.ReplaceAll(stopVpcepPath, "{cluster_id}", clusterId)

	stopVpcepOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err := client.Request("PUT", stopVpcepPath, &stopVpcepOpt)
	if err != nil {
		return err
	}

	return nil
}

func startVpcep(client *golangsdk.ServiceClient, clusterId string, bodyParams map[string]interface{}) error {
	startVpcepHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/vpcepservice/open"
	startVpcepPath := client.Endpoint + startVpcepHttpUrl
	startVpcepPath = strings.ReplaceAll(startVpcepPath, "{project_id}", client.ProjectID)
	startVpcepPath = strings.ReplaceAll(startVpcepPath, "{cluster_id}", clusterId)

	startVpcepOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         bodyParams,
	}

	_, err := client.Request("POST", startVpcepPath, &startVpcepOpt)
	if err != nil {
		return err
	}

	return nil
}

func updateVpcepWhitelist(client *golangsdk.ServiceClient, clusterId string, bodyParams map[string]interface{}) error {
	updateVpcepWhitelistHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/vpcepservice/permissions"
	updateVpcepWhitelistPath := client.Endpoint + updateVpcepWhitelistHttpUrl
	updateVpcepWhitelistPath = strings.ReplaceAll(updateVpcepWhitelistPath, "{project_id}", client.ProjectID)
	updateVpcepWhitelistPath = strings.ReplaceAll(updateVpcepWhitelistPath, "{cluster_id}", clusterId)

	updateVpcepWhitelistOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         bodyParams,
	}

	_, err := client.Request("POST", updateVpcepWhitelistPath, &updateVpcepWhitelistOpt)
	if err != nil {
		return err
	}

	return nil
}

func updateKibanaPublicAccess(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	o, n := d.GetChange("kibana_public_access")
	oValue := o.([]interface{})
	nValue := n.([]interface{})
	clusterId := d.Id()

	switch len(nValue) - len(oValue) {
	case -1: // delete kibana_public_access
		err := closeKibanaPublic(client, clusterId, oValue)
		if err != nil {
			return fmt.Errorf("error diabling the kibana public access of CSS cluster: %s, err: %s", clusterId, err)
		}
		err = checkClusterOperationResult(ctx, client, clusterId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

	case 1: // enable kibana_public_access
		err := startKibanaPublic(client, clusterId, d)
		if err != nil {
			return fmt.Errorf("error enabling the kibana public access of CSS cluster: %s, err: %s", clusterId, err)
		}
		err = checkClusterOperationResult(ctx, client, clusterId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

	case 0:
		// update bandwidth
		if d.HasChange("kibana_public_access.0.bandwidth") {
			err := updateKibanaPublicBandwidth(client, clusterId, d)
			if err != nil {
				return fmt.Errorf("error modifing bandwidth of the kibana public access of CSS cluster: %s, err: %s", clusterId, err)
			}
		}

		// update whitelist
		if d.HasChanges("kibana_public_access.0.whitelist", "kibana_public_access.0.whitelist_enabled") {
			// disable whitelist
			if !d.Get("kibana_public_access.0.whitelist_enabled").(bool) {
				err := stopKibanaPublicWhitelis(client, clusterId)
				if err != nil {
					return fmt.Errorf("error disabing the whitelist of the kibana public access of CSS cluster: %s, err: %s", clusterId, err)
				}
			} else {
				err := updateKibanaPublicWhitelist(client, clusterId, d)
				if err != nil {
					return fmt.Errorf("error modifing whitelist of the kibana public access of CSS cluster: %s, err: %s", clusterId, err)
				}
			}
		}
	}

	return nil
}

func updateKibanaPublicWhitelist(client *golangsdk.ServiceClient, clusterId string, d *schema.ResourceData) error {
	updateKibanaPublicWhitelistHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/publickibana/whitelist/update"
	updateKibanaPublicWhitelistPath := client.Endpoint + updateKibanaPublicWhitelistHttpUrl
	updateKibanaPublicWhitelistPath = strings.ReplaceAll(updateKibanaPublicWhitelistPath, "{project_id}", client.ProjectID)
	updateKibanaPublicWhitelistPath = strings.ReplaceAll(updateKibanaPublicWhitelistPath, "{cluster_id}", clusterId)

	updateKibanaPublicWhitelistOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"whiteList": d.Get("kibana_public_access.0.whitelist"),
		},
	}

	_, err := client.Request("POST", updateKibanaPublicWhitelistPath, &updateKibanaPublicWhitelistOpt)
	if err != nil {
		return err
	}

	return nil
}

func stopKibanaPublicWhitelis(client *golangsdk.ServiceClient, clusterId string) error {
	stopKibanaPublicWhitelisHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/publickibana/bandwidth"
	stopKibanaPublicWhitelisPath := client.Endpoint + stopKibanaPublicWhitelisHttpUrl
	stopKibanaPublicWhitelisPath = strings.ReplaceAll(stopKibanaPublicWhitelisPath, "{project_id}", client.ProjectID)
	stopKibanaPublicWhitelisPath = strings.ReplaceAll(stopKibanaPublicWhitelisPath, "{cluster_id}", clusterId)

	stopKibanaPublicWhitelisOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err := client.Request("PUT", stopKibanaPublicWhitelisPath, &stopKibanaPublicWhitelisOpt)
	if err != nil {
		return err
	}

	return nil
}

func updateKibanaPublicBandwidth(client *golangsdk.ServiceClient, clusterId string, d *schema.ResourceData) error {
	updateKibanaPublicBandwidthHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/publickibana/bandwidth"
	updateKibanaPublicBandwidthPath := client.Endpoint + updateKibanaPublicBandwidthHttpUrl
	updateKibanaPublicBandwidthPath = strings.ReplaceAll(updateKibanaPublicBandwidthPath, "{project_id}", client.ProjectID)
	updateKibanaPublicBandwidthPath = strings.ReplaceAll(updateKibanaPublicBandwidthPath, "{cluster_id}", clusterId)

	updateKibanaPublicBandwidthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"bandWidth": map[string]interface{}{
				"size": d.Get("kibana_public_access.0.bandwidth"),
			},
			"isAutoPay": 1,
		},
	}

	_, err := client.Request("POST", updateKibanaPublicBandwidthPath, &updateKibanaPublicBandwidthOpt)
	if err != nil {
		return err
	}

	return nil
}

func startKibanaPublic(client *golangsdk.ServiceClient, clusterId string, d *schema.ResourceData) error {
	startKibanaPublicHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/publickibana/open"
	startKibanaPublicPath := client.Endpoint + startKibanaPublicHttpUrl
	startKibanaPublicPath = strings.ReplaceAll(startKibanaPublicPath, "{project_id}", client.ProjectID)
	startKibanaPublicPath = strings.ReplaceAll(startKibanaPublicPath, "{cluster_id}", clusterId)

	startKibanaPublicOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"eipSize": d.Get("kibana_public_access.0.bandwidth"),
			"elbWhiteList": map[string]interface{}{
				"enableWhiteList": d.Get("kibana_public_access.0.whitelist_enabled"),
				"whiteList":       d.Get("kibana_public_access.0.whitelist"),
			},
		},
	}

	_, err := client.Request("POST", startKibanaPublicPath, &startKibanaPublicOpt)
	if err != nil {
		return err
	}

	return nil
}

func closeKibanaPublic(client *golangsdk.ServiceClient, clusterId string, oValue []interface{}) error {
	closeKibanaPublicHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/publickibana/close"
	closeKibanaPublicPath := client.Endpoint + closeKibanaPublicHttpUrl
	closeKibanaPublicPath = strings.ReplaceAll(closeKibanaPublicPath, "{project_id}", client.ProjectID)
	closeKibanaPublicPath = strings.ReplaceAll(closeKibanaPublicPath, "{cluster_id}", clusterId)

	closeKibanaPublicOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"eipSize": utils.PathSearch("bandwidth", oValue[0], 0),
			"elbWhiteList": map[string]interface{}{
				"enableWhiteList": utils.PathSearch("whitelist_enabled", oValue[0], false),
				"whiteList":       utils.PathSearch("whitelist", oValue[0], nil),
			},
		},
	}

	_, err := client.Request("PUT", closeKibanaPublicPath, &closeKibanaPublicOpt)
	if err != nil {
		return err
	}

	return nil
}

func updatePublicAccess(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	o, n := d.GetChange("public_access")
	oValue := o.([]interface{})
	nValue := n.([]interface{})
	clusterId := d.Id()

	switch len(nValue) - len(oValue) {
	case -1: // delete public_access
		err := closePublicAccess(client, clusterId, oValue)
		if err != nil {
			return fmt.Errorf("error diabling public access of CSS cluster: %s, err: %s", clusterId, err)
		}
		err = checkClusterOperationResult(ctx, client, clusterId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

	case 1:
		// enable public_access
		err := startPublicAccess(client, clusterId, d)
		if err != nil {
			return fmt.Errorf("error enabling public access of CSS cluster: %s, err: %s", clusterId, err)
		}

		err = checkClusterOperationResult(ctx, client, clusterId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

		if whitelist, ok := d.GetOk("public_access.0.whitelist"); ok {
			err := updatePublicWhitelist(client, clusterId, whitelist)
			if err != nil {
				return fmt.Errorf("error updating whitelist of public access of CSS cluster: %s, err: %s", clusterId, err)
			}
		}

	case 0:
		// disable whitelist
		if d.HasChanges("public_access.0.whitelist", "public_access.0.whitelist_enabled") {
			if !d.Get("public_access.0.whitelist_enabled").(bool) {
				err := stopPublicWhitelist(client, clusterId)
				if err != nil {
					return fmt.Errorf("error disabling whitelist of public access of CSS cluster: %s, err: %s", clusterId, err)
				}
			} else {
				err := updatePublicWhitelist(client, clusterId, d.Get("public_access.0.whitelist"))
				if err != nil {
					return fmt.Errorf("error updating whitelist of public access of CSS cluster: %s, err: %s", clusterId, err)
				}
			}
		}

		// update bandwidth
		if d.HasChange("public_access.0.bandwidth") {
			err := updatePublicBandWidth(client, clusterId, d)
			if err != nil {
				return fmt.Errorf("error disabling the whitelist of the kibana public access of CSS cluster: %s, err: %s", clusterId, err)
			}
		}
	}

	return nil
}

func updatePublicBandWidth(client *golangsdk.ServiceClient, clusterId string, d *schema.ResourceData) error {
	updatePublicBandWidthHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/public/bandwidth"
	updatePublicBandWidthPath := client.Endpoint + updatePublicBandWidthHttpUrl
	updatePublicBandWidthPath = strings.ReplaceAll(updatePublicBandWidthPath, "{project_id}", client.ProjectID)
	updatePublicBandWidthPath = strings.ReplaceAll(updatePublicBandWidthPath, "{cluster_id}", clusterId)

	updatePublicBandWidthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"bandWidth": map[string]interface{}{
				"size": d.Get("public_access.0.bandwidth"),
			},
			"isAutoPay": 1,
		},
	}

	_, err := client.Request("POST", updatePublicBandWidthPath, &updatePublicBandWidthOpt)
	if err != nil {
		return err
	}

	return nil
}

func stopPublicWhitelist(client *golangsdk.ServiceClient, clusterId string) error {
	stopPublicWhitelistHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/public/whitelist/close"
	stopPublicWhitelistPath := client.Endpoint + stopPublicWhitelistHttpUrl
	stopPublicWhitelistPath = strings.ReplaceAll(stopPublicWhitelistPath, "{project_id}", client.ProjectID)
	stopPublicWhitelistPath = strings.ReplaceAll(stopPublicWhitelistPath, "{cluster_id}", clusterId)

	stopPublicWhitelistOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err := client.Request("PUT", stopPublicWhitelistPath, &stopPublicWhitelistOpt)
	if err != nil {
		return err
	}

	return nil
}

func updatePublicWhitelist(client *golangsdk.ServiceClient, clusterId string, whitelist interface{}) error {
	updatePublicWhitelistHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/public/whitelist/update"
	updatePublicWhitelistPath := client.Endpoint + updatePublicWhitelistHttpUrl
	updatePublicWhitelistPath = strings.ReplaceAll(updatePublicWhitelistPath, "{project_id}", client.ProjectID)
	updatePublicWhitelistPath = strings.ReplaceAll(updatePublicWhitelistPath, "{cluster_id}", clusterId)

	updatePublicWhitelistOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"whiteList": whitelist,
		},
	}

	_, err := client.Request("POST", updatePublicWhitelistPath, &updatePublicWhitelistOpt)
	if err != nil {
		return err
	}

	return nil
}

func startPublicAccess(client *golangsdk.ServiceClient, clusterId string, d *schema.ResourceData) error {
	startPublicAccessHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/public/open"
	startPublicAccessPath := client.Endpoint + startPublicAccessHttpUrl
	startPublicAccessPath = strings.ReplaceAll(startPublicAccessPath, "{project_id}", client.ProjectID)
	startPublicAccessPath = strings.ReplaceAll(startPublicAccessPath, "{cluster_id}", clusterId)

	startPublicAccessOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"eip": map[string]interface{}{
				"bandWidth": map[string]interface{}{
					"size": d.Get("public_access.0.bandwidth"),
				},
			},
			"isAutoPay": 1,
		},
	}

	_, err := client.Request("POST", startPublicAccessPath, &startPublicAccessOpt)
	if err != nil {
		return err
	}

	return nil
}

func closePublicAccess(client *golangsdk.ServiceClient, clusterId string, oValue []interface{}) error {
	closePublicAccessHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/public/close"
	closePublicAccessPath := client.Endpoint + closePublicAccessHttpUrl
	closePublicAccessPath = strings.ReplaceAll(closePublicAccessPath, "{project_id}", client.ProjectID)
	closePublicAccessPath = strings.ReplaceAll(closePublicAccessPath, "{cluster_id}", clusterId)

	closePublicAccessOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"eip": map[string]interface{}{
				"bandWidth": map[string]interface{}{
					"size": utils.PathSearch("bandwidth", oValue[0], 0),
				},
			},
		},
	}

	_, err := client.Request("PUT", closePublicAccessPath, &closePublicAccessOpt)
	if err != nil {
		return err
	}

	return nil
}

func resourceCssClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterId := d.Id()
	client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	deleteClusterHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}"
	deleteClusterPath := client.Endpoint + deleteClusterHttpUrl
	deleteClusterPath = strings.ReplaceAll(deleteClusterPath, "{project_id}", client.ProjectID)
	deleteClusterPath = strings.ReplaceAll(deleteClusterPath, "{cluster_id}", clusterId)

	deleteClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteClusterPath, &deleteClusterOpt)
	if err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		return common.CheckDeletedDiag(d, common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015"), "error deleting the CSS cluster")
	}

	err = checkClusterDeleteResult(ctx, client, clusterId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("failed to check the result of deletion %s", err)
	}
	return nil
}

func checkClusterCreateResult(ctx context.Context, client *golangsdk.ServiceClient, clusterId string,
	timeout time.Duration) error {
	createStateConf := &resource.StateChangeConf{
		Pending: []string{ClusterStatusInProcess},
		Target:  []string{ClusterStatusAvailable},
		Refresh: func() (interface{}, string, error) {
			resp, err := getClusterDetails(client, clusterId)
			if err != nil {
				return nil, "failed", err
			}
			return resp, utils.PathSearch("status", resp, nil).(string), err
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

func checkClusterDeleteResult(ctx context.Context, client *golangsdk.ServiceClient, clusterId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Done"},
		Refresh: func() (interface{}, string, error) {
			_, err := getClusterDetails(client, clusterId)
			if err != nil {
				err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return true, "Done", nil
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

func checkClusterOperationResult(ctx context.Context, client *golangsdk.ServiceClient, clusterId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Done"},
		Refresh: func() (interface{}, string, error) {
			resp, err := getClusterDetails(client, clusterId)
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

func checkCssClusterIsReady(detail interface{}) bool {
	status := utils.PathSearch("status", detail, "").(string)
	actions := utils.PathSearch("actions", detail, make([]interface{}, 0)).([]interface{})
	instances := utils.PathSearch("instances", detail, make([]interface{}, 0)).([]interface{})
	if status != ClusterStatusAvailable {
		return false
	}

	// actions --- the behaviors on a cluster
	if len(actions) > 0 {
		return false
	}

	if len(instances) == 0 {
		return false
	}
	for _, v := range instances {
		status := utils.PathSearch("status", v, "").(string)
		if status != ClusterStatusAvailable {
			return false
		}
	}
	return true
}

func updateFlavor(ctx context.Context, d *schema.ResourceData,
	flavorsResp interface{}, conf *config.Config, addMasterNode, addClientNode bool) error {
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
	resp interface{}, conf *config.Config) error {
	region := conf.GetRegion(d)
	cssV1Client, err := conf.CssV1Client(region)
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

	err = checkClusterOperationResult(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	return nil
}

func flattenFlavorId(nodeType string, d *schema.ResourceData, resp interface{}) (string, error) {
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

func getFlavorList(client *golangsdk.ServiceClient) (interface{}, error) {
	getEsFlavorHttpUrl := "v1.0/{project_id}/es-flavors"
	getEsFlavorPath := client.Endpoint + getEsFlavorHttpUrl
	getEsFlavorPath = strings.ReplaceAll(getEsFlavorPath, "{project_id}", client.ProjectID)

	getEsFlavorOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getEsFlavorResp, err := client.Request("GET", getEsFlavorPath, &getEsFlavorOpt)
	if err != nil {
		return getEsFlavorResp, err
	}

	return utils.FlattenResponse(getEsFlavorResp)
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

func updateExtendInstanceStorage(ctx context.Context, d *schema.ResourceData, conf *config.Config, bodyParams map[string]interface{}) error {
	region := conf.GetRegion(d)
	cssV1Client, err := conf.CssV1Client(region)
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

	err = checkClusterOperationResult(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
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

	err = checkClusterOperationResult(ctx, cssV1Client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	return nil
}

func updateSafeMode(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
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

	err = checkClusterOperationResult(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate))
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

func updateAdminPassword(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
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

	err = checkClusterOperationResult(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func updateSecurityGroup(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
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

	err = checkClusterOperationResult(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate))
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

	if len(shrinkNodeIds) == 0 {
		bodyParams := buildUpdateShrinkNodeByTypeBodyParams(nodeType, shrinkNodeSize)
		err := updateShrinkInstanceNodeByType(ctx, d, client, bodyParams)
		if err != nil {
			return err
		}
	} else {
		if shrinkNodeSize != len(shrinkNodeIds) {
			return fmt.Errorf("instance_number changing number is inconsistent with the length"+
				" of shrink_node_ids, node type: %s", nodeType)
		}
		err := updateShrinkInstanceNodeById(ctx, d, client, shrinkNodeIds)
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

func updateShrinkInstanceNodeByType(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	bodyParams map[string]interface{}) error {
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

	err = checkClusterOperationResult(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func updateShrinkInstanceNodeById(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, nodeIds []interface{}) error {
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

	err = checkClusterOperationResult(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	return nil
}

type ResponseError struct {
	ErrorCode string `json:"errCode"`
	ErrorMsg  string `json:"externalMessage"`
}
