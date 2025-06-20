package modelarts

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v2ReourcePoolNonUpdatableParams = []string{
	"spec",
	"spec.*.type",
	"spec.*.network",
	"spec.*.network.name",
	"spec.*.network.vpc_id",
	"spec.*.network.subnet_id",
	"spec.*.user_login",
	"spec.*.user_login.key_pair_name",
	"spec.*.user_login.password",
	"spec.*.clusters",
	"spec.*.clusters.name",
	"spec.*.clusters.provider_id",
}

const postPaid = "0"
const prePaid = "1"

// @API ModelArts POST /v2/{project_id}/pools
// @API ModelArts DELETE /v2/{project_id}/pools/{id}
// @API ModelArts GET /v2/{project_id}/pools/{id}
// @API ModelArts PATCH /v2/{project_id}/pools/{id}
func Resourcev2ResourcePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2ResourcePoolCreate,
		ReadContext:   resourceV2ResourcePoolRead,
		UpdateContext: resourceV2ResourcePoolUpdate,
		DeleteContext: resourceV2ResourcePoolDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(v2ReourcePoolNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"metadata": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"labels": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  `The labels of the resource pool, in JSON format.`,
						},
						"annotations": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  `The annotations of the resource pool, in JSON format.`,
						},
						// Internal attributes.
						"labels_origin": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The script configuration value of this change is also the original value used for the new value
next time in the list build method and reorder. The corresponding parameter name is 'metadata.labels'.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
					},
				},
				Description: `The metadata of the resource pool.`,
			},
			"spec": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resources": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        resourcePoolv2SpecResourceSchema(),
							Description: `The list of resource specifications in the resource pool.`,
						},
						"scope": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of job types supported by the resource pool.`,
						},
						"network": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `The name of the network.`,
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `The ID of the VPC.`,
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `The ID of the subnet.`,
									},
								},
							},
							Description: `The network of the resource pool.`,
						},
						"user_login": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{Schema: map[string]*schema.Schema{
								"key_pair_name": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: `The name of the key pair.`,
								},
								"password": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: `The password of the resource pool.`,
								},
							}},
							Description: `The user login information of the privileged pool.`,
						},
						"clusters": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{Schema: map[string]*schema.Schema{
								"name": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: `The name of the cluster.`,
								},
								"provider_id": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: `The provider ID of the cluster.`,
								},
							}},
							Description: `The type of the cluster.`,
						},
					},
				},
				Description: `The specification of the resource pool.`,
			},
		},
	}
}

func resourcePoolv2SpecResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The flavor of the resource pool.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The count of the resource pool.`,
			},
			"max_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The max number of resources of the corresponding flavors.`,
			},
			"node_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of resource pool nodes.`,
			},
			"taints": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"key": {
						Type:        schema.TypeString,
						Required:    true,
						Description: `The key of the taint.`,
					},
					"effect": {
						Type:        schema.TypeString,
						Required:    true,
						Description: `The effect of the taint.`,
					},
					"value": {
						Type:        schema.TypeString,
						Optional:    true,
						Computed:    true,
						Description: `The value of the taint.`,
					},
				}},
				Description: `The taint list of the resource pool.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The key/value pairs labels of resource pool.`,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The key of the tag.`,
						},
						"value": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The value of the tag.`,
						},
					},
				},
				Description: `The tags of resource pool.`,
			},
			"network": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The ID of the VPC.`,
						},
						"subnet": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The ID of the subnet.`,
						},
						"security_groups": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The ID list of the security group.`,
						},
					},
				},
				Description: `The network of the privileged pool.`,
			},
			"extend_params": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(_, o, n string, _ *schema.ResourceData) bool {
					// The current SuppressMapDiffs method just only supports object type sub-parameters, and does not
					// support list type sub-parameters.
					return utils.ContainsAllKeyValues(utils.TryMapValueAnalysis(o), utils.TryMapValueAnalysis(n))
				},
				Description: `The extend params of the resource pool, in JSON format.`,
			},
			"creating_step": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"step": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The creation step of the resource pool nodes.`,
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the resource pool nodes.`,
						},
					},
				},
				Description: `The creation step configuration of the resource pool nodes.`,
			},
			"root_volume": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the root volume.`,
						},
						"size": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The size of the root volume.`,
						},
					},
				},
				Description: `The root volume of the resource pool nodes.`,
			},
			"data_volumes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the data volume.`,
						},
						"size": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The size of the data volume.`,
						},
						"extend_params": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							DiffSuppressFunc: func(_, o, n string, _ *schema.ResourceData) bool {
								// The current SuppressMapDiffs method just only supports object type sub-parameters, and does not
								// support list type sub-parameters.
								return utils.ContainsAllKeyValues(utils.TryMapValueAnalysis(o), utils.TryMapValueAnalysis(n))
							},
							Description: `The extend parameters of the data volume, in JSON format.`,
						},
						"count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `The count of the current data volume configuration.`,
						},
					},
				},
				Description: `The data volumes of the resource pool nodes.`,
			},
			"volume_group_configs": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_group": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the volume group.`,
						},
						"docker_thin_pool": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `The percentage of container volumes to data volumes on resource pool nodes.`,
						},
						"lvm_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lv_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The LVM write mode.`,
									},
									"path": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `The volume mount path.`,
									},
								},
							},
							Description: `The configuration of the LVM management.`,
						},
						"types": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The storage types of the volume group.`,
						},
					},
				},
				Description: `The extend configurations of the volume groups.`,
			},
			"os": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        modelartsResourcePoolResourcesOsSchema(),
				Description: `The image information for the specified OS.`,
			},
			"driver": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        modelartsResourcePoolResourcesDriverSchema(),
				Description: `The driver information of the resource pool nodes.`,
			},
			// Internal parameters.
			"azs": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"az": {
							Type:     schema.TypeString,
							Required: true,
							Description: utils.SchemaDesc(
								`The AZ name.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"count": {
							Type:     schema.TypeInt,
							Required: true,
							Description: utils.SchemaDesc(
								`The number of nodes in the AZ.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
					},
				},
				Description: utils.SchemaDesc(
					`The AZ list of the resource pool nodes.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			// Internal attribute(s).
			"volume_group_configs_origin": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_group": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The script configuration value of this change is also the original value used for the new value
next time in the list build method and reorder. The corresponding parameter name is 'volume_group_configs.volume_group'.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"docker_thin_pool": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: utils.SchemaDesc(
								`The script configuration value of this change is also the original value used for the new value
next time in the list build method and reorder. The corresponding parameter name is 'volume_group_configs.docker_thin_pool'.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"lvm_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lv_type": {
										Type:     schema.TypeString,
										Computed: true,
										Description: utils.SchemaDesc(
											`The script configuration value of this change is also the original value used for the new value
next time in the list build method and reorder. The corresponding parameter name is 'volume_group_configs.lvm_config.lv_type'.`,
											utils.SchemaDescInput{
												Internal: true,
											},
										),
									},
									"path": {
										Type:     schema.TypeString,
										Computed: true,
										Description: utils.SchemaDesc(
											`The script configuration value of this change is also the original value used for the new value
next time in the list build method and reorder. The corresponding parameter name is 'volume_group_configs.lvm_config.path'.`,
											utils.SchemaDescInput{
												Internal: true,
											},
										),
									},
								},
							},
							Description: utils.SchemaDesc(
								`The script configuration value of this change is also the original value used for the new value
next time in the list build method and reorder. The corresponding parameter name is 'volume_group_configs.lvm_config'.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Description: utils.SchemaDesc(
								`The script configuration value of this change is also the original value used for the new value
next time in the list build method and reorder. The corresponding parameter name is 'volume_group_configs.types'.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
					},
				},
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for the new value
next time in the list build method and reorder. The corresponding parameter name is 'volume_group_configs'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func resourceV2ResourcePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/pools"
	)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCreateV2ResourcePoolBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Modelarts resource pool: %s", err)
	}

	respBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	poolName := utils.PathSearch("metadata.name", respBody, nil)
	if poolName == nil {
		return diag.Errorf("unable to find the resource pool name in the API response")
	}
	d.SetId(poolName.(string))

	if getChargingMode(d) == prePaid {
		// wait 30 seconds so that the resource pool can be queried
		// TODO: 测试是否需要等待30s
		// time.Sleep(30 * time.Second)
		orderId := utils.PathSearch(`metadata.annotations."os.modelarts/order.id"`, respBody, nil)
		if orderId == nil {
			return diag.Errorf("unable to find the order ID in the API response")
		}
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = common.WaitOrderResourceComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = waitForResourcePoolStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts V2 resource pool (%s) creation to complete: %s", d.Id(), err)
	}

	if err = waitForDriverStatusCompleted(ctx, cfg, region, d); err != nil {
		return diag.Errorf("error waiting for the Modelarts V2 resource pool (%s) driver status to become running: %s", d.Id(), err)
	}

	return resourceV2ResourcePoolRead(ctx, d, meta)
}

func getChargingMode(d *schema.ResourceData) string {
	annotations, ok := d.GetOk("metadata.0.annotations")
	if ok {
		return utils.PathSearch("os.modelarts/billing.mode", utils.StringToJson(annotations.(string)), postPaid).(string)
	}
	return postPaid
}

func buildCreateV2ResourcePoolBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"apiVersion": "v2",
		"kind":       "Pool",
		"metadata":   buildCreateV2ResourcePoolMetaData(d),
		"spec":       buildCreateV2ResourcePoolSpec(d),
	}
	return bodyParams
}

func buildCreateV2ResourcePoolMetaData(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"labels":      utils.StringToJson(d.Get("metadata.0.labels").(string)),
		"annotations": utils.StringToJson(d.Get("metadata.0.annotations").(string)),
	}
}

func buildCreateV2ResourcePoolSpec(d *schema.ResourceData) map[string]interface{} {
	specifications := d.Get("spec").([]interface{})
	if len(specifications) < 1 || specifications[0] == nil {
		return nil
	}

	specification := specifications[0]
	params := map[string]interface{}{
		// Currently only support `Dedicate`.
		"type":      "Dedicate",
		"resources": buildV2ResourcePoolSpecResources(d),
		"scope": utils.ValueIgnoreEmpty(utils.ExpandToStringListBySet(utils.PathSearch("scope", specification,
			schema.NewSet(schema.HashString, nil)).(*schema.Set))),
		"network": buildCreateV2ResourcePoolSpecNetwork(utils.PathSearch("network", specification,
			make([]interface{}, 0)).([]interface{})),
		"userLogin": buildCreateV2ResourcePoolSpecUserLogin(utils.PathSearch("user_login", specification,
			make([]interface{}, 0)).([]interface{})),
		"clusters": buildCreateV2ResourcePoolSpecClusters(utils.PathSearch("clusters", specification,
			make([]interface{}, 0)).([]interface{})),
	}
	return params
}

func buildV2ResourcePoolSpecResources(d *schema.ResourceData) []map[string]interface{} {
	oldResourcesVal, newResourcesVal := d.GetChange("spec.0.resources")
	oldResources := oldResourcesVal.([]interface{})
	newResources := newResourcesVal.([]interface{})

	result := make([]map[string]interface{}, len(newResources))
	for i, v := range newResources {
		result[i] = map[string]interface{}{
			"flavor":   utils.PathSearch("flavor", v, nil),
			"count":    utils.PathSearch("count", v, nil),
			"maxCount": utils.ValueIgnoreEmpty(utils.PathSearch("max_count", v, nil)),
			"azs": buildV2ResourcePoolResourcesAzs(utils.PathSearch("azs", v,
				schema.NewSet(schema.HashString, nil)).(*schema.Set)),
			"nodePool": utils.ValueIgnoreEmpty(utils.PathSearch("node_pool", v, nil)),
			"taints": buildV2ResourcePoolResourcesTaints(utils.PathSearch("taints", v,
				schema.NewSet(schema.HashString, nil)).(*schema.Set)),
			"labels":  utils.ValueIgnoreEmpty(utils.PathSearch("labels", v, nil)),
			"tags":    buildV2ResourcePoolResourcesTags(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
			"network": buildV2ResourcePoolSpecResourcesNetwork(utils.PathSearch("network", v, make([]interface{}, 0)).([]interface{})),
			"extendParams": buildV2ResourcePoolResourcesExtendParams(
				utils.PathSearch(fmt.Sprintf("[%d].extend_params", i), oldResources, "{}").(string),
				utils.PathSearch("extend_params", v, "{}").(string),
			),
			"creatingStep": buildResourcePoolResourcesCreatingStep(
				utils.PathSearch("creating_step", v, make([]interface{}, 0)).([]interface{})),
			"rootVolume": buildV2ResourcePoolResourcesRootVolume(utils.PathSearch("root_volume", v,
				make([]interface{}, 0)).([]interface{})),
			"dataVolumes": buildV2ResourcePoolResourcesDataVolumes(
				utils.PathSearch(fmt.Sprintf("[%d].data_volumes", i), oldResources, make([]interface{}, 0)).([]interface{}),
				utils.PathSearch("data_volumes", v, make([]interface{}, 0)).([]interface{}),
			),
			"volumeGroupConfigs": buildV2ResourcePoolResourcesVolumeGroupConfigs(
				utils.PathSearch(fmt.Sprintf("[%d].volume_group_configs_origin", i), oldResources, make([]interface{}, 0)).([]interface{}),
				utils.PathSearch("volume_group_configs", v, schema.NewSet(schema.HashString, nil)).(*schema.Set).List(),
			),
			"os":     buildV2ResourcePoolResourcesOsInfo(utils.PathSearch("os", v, make([]interface{}, 0)).([]interface{})),
			"driver": buildV2ResourcePoolResourcesDriver(utils.PathSearch("driver", v, make([]interface{}, 0)).([]interface{})),
		}
	}
	return result
}

func buildV2ResourcePoolResourcesAzs(azs *schema.Set) []map[string]interface{} {
	if azs.Len() < 1 {
		return nil
	}

	result := make([]map[string]interface{}, azs.Len())
	for i, az := range azs.List() {
		result[i] = map[string]interface{}{
			"az":    utils.PathSearch("az", az, nil),
			"count": utils.PathSearch("count", az, nil),
		}
	}

	return result
}

func buildV2ResourcePoolResourcesTaints(taints *schema.Set) []map[string]interface{} {
	if taints.Len() < 1 {
		return nil
	}

	result := make([]map[string]interface{}, taints.Len())
	for i, taint := range taints.List() {
		result[i] = map[string]interface{}{
			"key":    utils.PathSearch("key", taint, nil),
			"value":  utils.PathSearch("value", taint, nil),
			"effect": utils.ValueIgnoreEmpty(utils.PathSearch("effect", taint, nil)),
		}
	}

	return result
}

func buildV2ResourcePoolResourcesTags(tags []interface{}) []interface{} {
	if len(tags) < 1 {
		return nil
	}
	result := make([]interface{}, len(tags))
	for i, tag := range tags {
		result[i] = map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.PathSearch("value", tag, nil),
		}
	}
	return result
}

func buildV2ResourcePoolSpecResourcesNetwork(resourceNetworks []interface{}) map[string]interface{} {
	if len(resourceNetworks) < 1 || resourceNetworks[0] == nil {
		return nil
	}

	return utils.RemoveNil(map[string]interface{}{
		"vpc":    utils.ValueIgnoreEmpty(utils.PathSearch("[0].vpc", resourceNetworks, nil)),
		"subnet": utils.ValueIgnoreEmpty(utils.PathSearch("[0].subnet", resourceNetworks, nil)),
		"securityGroups": utils.ValueIgnoreEmpty(utils.ExpandToStringListBySet(utils.PathSearch("[0].security_groups",
			resourceNetworks, schema.NewSet(schema.HashString, nil)).(*schema.Set))),
	})
}

func buildCreateV2ResourcePoolSpecNetwork(networks []interface{}) map[string]interface{} {
	if len(networks) < 1 || networks[0] == nil {
		return nil
	}

	return utils.RemoveNil(map[string]interface{}{
		"name":     utils.ValueIgnoreEmpty(utils.PathSearch("[0].name", networks, nil)),
		"vpcId":    utils.ValueIgnoreEmpty(utils.PathSearch("[0].vpc_id", networks, nil)),
		"subnetId": utils.ValueIgnoreEmpty(utils.PathSearch("[0].subnet_id", networks, nil)),
	})
}

func buildCreateV2ResourcePoolSpecUserLogin(userLogin []interface{}) map[string]interface{} {
	if len(userLogin) < 1 || userLogin[0] == nil {
		return nil
	}
	return utils.RemoveNil(map[string]interface{}{
		"keyPairName": utils.ValueIgnoreEmpty(utils.PathSearch("[0].key_pair_name", userLogin, nil)),
		"password":    utils.ValueIgnoreEmpty(utils.PathSearch("[0].password", userLogin, nil)),
	})
}

func buildCreateV2ResourcePoolSpecClusters(clusters []interface{}) []interface{} {
	if len(clusters) < 1 {
		return nil
	}

	result := make([]interface{}, len(clusters))
	for i, v := range clusters {
		result[i] = map[string]interface{}{
			"name":       utils.ValueIgnoreEmpty(utils.PathSearch("name", v, nil)),
			"providerId": utils.ValueIgnoreEmpty(utils.PathSearch("provider_id", v, nil)),
		}
	}
	return result
}

func buildV2ResourcePoolResourcesExtendParams(oldExtendParams, newExtendParams string) map[string]interface{} {
	extendParams := utils.TryMapValueAnalysis(utils.StringToJson(oldExtendParams))
	if objExtendParams := utils.TryMapValueAnalysis(utils.StringToJson(newExtendParams)); len(objExtendParams) > 0 {
		for k, v := range objExtendParams {
			extendParams[k] = v
		}
	}
	return extendParams
}

func buildV2ResourcePoolResourcesRootVolume(rootVolumes []interface{}) map[string]interface{} {
	if len(rootVolumes) < 1 || rootVolumes[0] == nil {
		return nil
	}

	return utils.RemoveNil(map[string]interface{}{
		"volumeType": utils.PathSearch("[0].volume_type", rootVolumes, nil),
		"size":       utils.PathSearch("[0].size", rootVolumes, nil),
	})
}

func buildV2ResourcePoolResourcesDataVolumes(oldDataVolumes, newDataVolumes []interface{}) []map[string]interface{} {
	if len(newDataVolumes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(newDataVolumes))
	for i, dataVolume := range newDataVolumes {
		result = append(result, map[string]interface{}{
			"volumeType": utils.PathSearch("volume_type", dataVolume, nil),
			"size":       utils.PathSearch("size", dataVolume, nil),
			"extendParams": buildV2ResourcePoolResourcesExtendParams(
				utils.PathSearch(fmt.Sprintf("[%d].extend_params", i), oldDataVolumes, "").(string),
				utils.PathSearch("extend_params", dataVolume, "").(string),
			),
			"count": utils.ValueIgnoreEmpty(utils.PathSearch("count", dataVolume, nil)),
		})
	}

	return result
}

func buildV2ResourcePoolResourcesVolumeGroupConfigs(oldVolumeGroupConfigs, newVolumeGroupConfigs []interface{}) []map[string]interface{} {
	if len(oldVolumeGroupConfigs) < 1 {
		result := make([]map[string]interface{}, 0, len(newVolumeGroupConfigs))
		for _, volumeGroupConfig := range newVolumeGroupConfigs {
			result = append(result, map[string]interface{}{
				"volumeGroup":    utils.PathSearch("volume_group", volumeGroupConfig, nil),
				"dockerThinPool": utils.ValueIgnoreEmpty(utils.PathSearch("docker_thin_pool", volumeGroupConfig, nil)),
				"lvmConfig": buildResourceVolumeGroupConfigsLvmConfig(utils.PathSearch("lvm_config", volumeGroupConfig,
					make([]interface{}, 0)).([]interface{})),
				"types": utils.ValueIgnoreEmpty(utils.ExpandToStringListBySet(utils.PathSearch("types", volumeGroupConfig,
					schema.NewSet(schema.HashString, nil)).(*schema.Set))),
			})
		}
		return result
	}

	result := make([]map[string]interface{}, len(oldVolumeGroupConfigs))
	for _, volumeGroupConfig := range oldVolumeGroupConfigs {
		newVolumeGroupConfig := utils.PathSearch(fmt.Sprintf("[?volume_group=='%s']|[0]",
			utils.PathSearch("volume_group", volumeGroupConfig, "").(string)), newVolumeGroupConfigs, make(map[string]interface{}))

		elem := map[string]interface{}{
			// Required parameter.
			"volumeGroup": utils.PathSearch("volume_group", volumeGroupConfig, nil),
		}

		if dockerThinPool := utils.PathSearch("docker_thin_pool", newVolumeGroupConfig, 0).(int); dockerThinPool != 0 {
			elem["dockerThinPool"] = dockerThinPool
		} else {
			elem["dockerThinPool"] = utils.ValueIgnoreEmpty(utils.PathSearch("docker_thin_pool", volumeGroupConfig, nil))
		}

		if lvmConfigs := utils.PathSearch("lvm_config", newVolumeGroupConfig, make([]interface{}, 0)).([]interface{}); len(lvmConfigs) > 0 {
			elem["lvmConfig"] = buildResourceVolumeGroupConfigsLvmConfig(lvmConfigs)
		} else {
			elem["lvmConfig"] = utils.ValueIgnoreEmpty(buildResourceVolumeGroupConfigsLvmConfig(
				utils.PathSearch("lvm_config", volumeGroupConfig, make([]interface{}, 0)).([]interface{})))
		}

		if types := utils.PathSearch("types", newVolumeGroupConfig, schema.NewSet(schema.HashString, nil)).(*schema.Set); types.Len() > 0 {
			elem["types"] = utils.ValueIgnoreEmpty(types.List())
		} else {
			elem["types"] = utils.ValueIgnoreEmpty(utils.ExpandToStringListBySet(utils.PathSearch("types", volumeGroupConfig,
				schema.NewSet(schema.HashString, nil)).(*schema.Set)))
		}

		result = append(result, elem)
	}
	return result
}

func buildV2ResourcePoolResourcesOsInfo(osInfos []interface{}) map[string]interface{} {
	if len(osInfos) < 1 || osInfos[0] == nil {
		return nil
	}

	osInfo := osInfos[0]
	return map[string]interface{}{
		"name":      utils.ValueIgnoreEmpty(utils.PathSearch("name", osInfo, nil)),
		"imageId":   utils.ValueIgnoreEmpty(utils.PathSearch("image_id", osInfo, nil)),
		"imageType": utils.ValueIgnoreEmpty(utils.PathSearch("image_type", osInfo, nil)),
	}
}

func buildV2ResourcePoolResourcesDriver(drivers []interface{}) map[string]interface{} {
	if len(drivers) < 1 || drivers[0] == nil {
		return nil
	}

	driver := drivers[0]
	return map[string]interface{}{
		"version": utils.ValueIgnoreEmpty(utils.PathSearch("version", driver, nil)),
	}
}

func waitForResourcePoolStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			resourcePool, err := queryResourcePool(cfg, cfg.GetRegion(d), d)
			if err != nil {
				return nil, "ERROR", err
			}

			if utils.PathSearch("status.resources.abnormal", resourcePool, nil) != nil {
				return nil, "ERROR", fmt.Errorf("error creating resource pool: the resource pool is abnormal")
			}

			status := utils.PathSearch(`status.phase`, resourcePool, "").(string)
			if status == "Error" {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `status.phase`)
			}

			creating := utils.PathSearch("status.resources.creating", resourcePool, make([]interface{}, 0)).([]interface{})
			deleting := utils.PathSearch("status.resources.creating", resourcePool, make([]interface{}, 0)).([]interface{})
			creationFaild := utils.PathSearch("status.resources.creationFaild", resourcePool, make([]interface{}, 0)).([]interface{})
			// Check whether the resource pool is running and the node expansion and contraction is successful.
			if status == "Running" && len(creating) == 0 && len(deleting) == 0 && len(creationFaild) == 0 {
				return resourcePool, "COMPLETED", nil
			}

			return resourcePool, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        30 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceV2ResourcePoolRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	resourcePool, err := queryResourcePool(cfg, region, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Modelarts V2 resource pool")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("metadata", flattenV2ResourcePoolMetadataInfo(d, utils.PathSearch("metadata", resourcePool, nil))),
		d.Set("spec", flattenV2ResourcePoolSpec(utils.PathSearch("spec", resourcePool, nil))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenV2ResourcePoolMetadataInfo(d *schema.ResourceData, metadata interface{}) []map[string]interface{} {
	if metadata == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"labels":        utils.PathSearch("labels", metadata, nil),
			"annotations":   utils.PathSearch("metadata.0.annotations", metadata, nil),
			"labels_origin": d.Get("metadata.0.labels_origin"),
		},
	}
}

func flattenV2ResourcePoolSpec(spec interface{}) []map[string]interface{} {
	if spec == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":       utils.PathSearch("type", spec, nil),
			"resources":  flattenV2ResourcePoolSpecResources(utils.PathSearch("resources", spec, make([]interface{}, 0)).([]interface{})),
			"scope":      utils.PathSearch("scope", spec, nil),
			"network":    flattenV2ResourcePoolSpecNetwork(utils.PathSearch("network", spec, nil)),
			"user_login": flattenV2ResourcePoolSpecUserLogin(utils.PathSearch("user_login", spec, nil)),
			"clusters":   flattenV2ResourcePoolClusters(utils.PathSearch("clusters", spec, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenV2ResourcePoolSpecResources(resources []interface{}) []map[string]interface{} {
	if len(resources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(resources))
	for _, v := range resources {
		result = append(result, map[string]interface{}{
			"flavor":        utils.PathSearch("flavor", v, nil),
			"count":         utils.PathSearch("count", v, nil),
			"max_count":     utils.PathSearch("maxCount", v, nil),
			"azs":           flattenV2ResourcePoolResourceAzs(utils.PathSearch("azs", v, make([]interface{}, 0)).([]interface{})),
			"node_pool":     utils.PathSearch("nodePool", v, nil),
			"taints":        flattenV2ResourcePoolResourceTaints(utils.PathSearch("taints", v, make([]interface{}, 0)).([]interface{})),
			"labels":        utils.PathSearch("labels", v, nil),
			"tags":          flattenResourcePoolSpecResourcesTags(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
			"network":       flattenResourcePoolSpecResourcesNetwork(utils.PathSearch("network", v, nil)),
			"extend_params": utils.JsonToString(utils.PathSearch("extendParams", v, nil)),
			"creating_step": flattenResourcePoolResourcesCreatingStep(utils.PathSearch("creatingStep", v, nil)),
			"root_volume":   flattenResourcePoolResourcesRootVolume(utils.PathSearch("rootVolume", v, nil)),
			"data_volumes": flattenResourcePoolResourcesDataVolumes(utils.PathSearch("dataVolumes",
				v, make([]interface{}, 0)).([]interface{})),
			"volume_group_configs": flattenResourcePoolResourcesVolumeGroupConfigs(utils.PathSearch("volumeGroupConfigs",
				v, make([]interface{}, 0)).([]interface{})),
			"volume_group_configs_origin": flattenResourcePoolResourcesVolumeGroupConfigs(utils.PathSearch("volumeGroupConfigs",
				v, make([]interface{}, 0)).([]interface{})),
			"os":     flattenResourcePoolResourcesOsInfo(utils.PathSearch("os", v, nil)),
			"driver": flattenResourcePoolResourcesDriver(utils.PathSearch("driver", v, nil))})
	}

	return result
}

func flattenResourcePoolSpecResourcesTags(tags []interface{}) []interface{} {
	if len(tags) == 0 {
		return nil
	}

	result := make([]interface{}, len(tags))
	for i, tag := range tags {
		result[i] = map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.PathSearch("value", tag, nil),
		}
	}
	return result
}

func flattenResourcePoolSpecResourcesNetwork(network interface{}) []map[string]interface{} {
	if network == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"vpc":             utils.PathSearch("vpc", network, nil),
			"subnet":          utils.PathSearch("subnet", network, nil),
			"security_groups": utils.PathSearch("securityGroups", network, nil),
		},
	}
}

func flattenV2ResourcePoolSpecNetwork(network interface{}) []map[string]interface{} {
	if network == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":      utils.PathSearch("name", network, nil),
			"vpc_id":    utils.PathSearch("vpcId", network, nil),
			"subnet_id": utils.PathSearch("subnetId", network, nil),
		},
	}
}

func flattenV2ResourcePoolSpecUserLogin(userLogin interface{}) []map[string]interface{} {
	if userLogin == nil {
		return nil
	}

	keyPairName := utils.PathSearch("keyPairName", userLogin, "").(string)
	if keyPairName != "" {
		return []map[string]interface{}{
			{
				"key_pair_name": keyPairName,
			},
		}
	}
	return nil
}

func resourceV2ResourcePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	if !d.HasChanges("metadata", "spec") {
		return nil
	}

	updateResourcePoolClient, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resp, err := updateResourcePool(updateResourcePoolClient, d, d.Id())
	if err != nil {
		return diag.Errorf("error updating Modelarts V2 resource pool: %s", err)
	}

	// Only when expanding prepaid type nodes, we need to determine the order status.
	// The result of the getChargingMode method is the charging mode of the nodes or resource pool.
	//     + When expanding the nodes, it indicates the charging mode of the nodes.
	if getChargingMode(d) == prePaid && d.HasChange("spec.0.resources") {
		// Whenever any count in resouces changes, the order status needs to be determined.
		oldRaw, newRaw := d.GetChange("spec.0.resources")
		if isAnyNodeScalling(oldRaw.([]interface{}), newRaw.([]interface{})) {
			updateRespBody, err := utils.FlattenResponse(resp)
			if err != nil {
				return diag.FromErr(err)
			}

			orderId := utils.PathSearch(`metadata.annotations."os.modelarts/order.id"`, updateRespBody, nil)
			if orderId == nil {
				return diag.Errorf("error updating Modelarts resource pool: order ID is not found in API response")
			}

			bssClient, err := cfg.BssV2Client(region)
			if err != nil {
				return diag.Errorf("error creating BSS v2 client: %s", err)
			}

			err = common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}

			_, err = common.WaitOrderAllResourceComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	err = waitForResourcePoolStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts V2 resource pool (%s) update to complete: %s", d.Id(), err)
	}

	return resourceV2ResourcePoolRead(ctx, d, meta)
}

func isAnyNodeScalling(oldResource, newResource []interface{}) bool {
	for i, v := range newResource {
		oldCount := utils.PathSearch(fmt.Sprintf("[%d].count", i), oldResource, 0).(int)
		newCount := utils.PathSearch("count", v, 0).(int)
		if oldCount != newCount {
			return true
		}
	}
	return false
}

func updateResourcePool(client *golangsdk.ServiceClient, d *schema.ResourceData, resourcePoolId string) (*http.Response, error) {
	updatehttpUrl := "v2/{project_id}/pools/{pool_name}"
	updatePath := client.Endpoint + updatehttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{pool_name}", resourcePoolId)

	updateResourcePoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/merge-patch+json"},
		JSONBody:         utils.RemoveNil(buildUpdateV2ResourcePoolBodyParams(d)),
	}

	return client.Request("PATCH", updatePath, &updateResourcePoolOpt)
}

func buildUpdateV2ResourcePoolBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": buildUpdateV2ResourcePoolMetaData(d),
		"spec":     buildUpdateV2ResourcePoolSpec(d),
	}
	return bodyParams
}

func buildUpdateV2ResourcePoolMetaData(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"labels":      utils.StringToJson(d.Get("metadata.0.labels").(string)).(map[string]interface{}),
		"annotations": buildUpdateV2ResourcePoolMetaDataAnnotations(d),
	}
	return params
}

func buildUpdateV2ResourcePoolMetaDataAnnotations(d *schema.ResourceData) map[string]interface{} {
	params := make(map[string]interface{})
	if annotations, ok := d.GetOk("metadata.0.annotations"); ok {
		params = utils.StringToJson(annotations.(string)).(map[string]interface{})
	}
	// If the node pools are not increased, delete the billing mode related parameters.
	oldRaw, newRaw := d.GetChange("spec.0.resources")
	if !isAnyNodePoolCountIncrease(oldRaw.([]interface{}), newRaw.([]interface{})) {
		delete(params, "os.modelarts/billing.mode")
		delete(params, "os.modelarts/period.num")
		delete(params, "os.modelarts/period.type")
		delete(params, "os.modelarts/auto.renew")
		delete(params, "os.modelarts/promotion.info")
		delete(params, "os.modelarts/service.console.url")
		delete(params, "os.modelarts/flavor.resource.ids")
		delete(params, "os.modelarts/order.id")
		delete(params, "os.modelarts/auto.pay")
	}
	return params
}

func isAnyNodePoolCountIncrease(oldResource, newRawResource []interface{}) bool {
	for i, v := range newRawResource {
		oldCount := utils.PathSearch(fmt.Sprintf("[%d].count", i), oldResource, 0).(int)
		newCount := utils.PathSearch("count", v, 0).(int)
		if newCount > oldCount {
			return true
		}
	}
	return false
}

func buildUpdateV2ResourcePoolSpec(d *schema.ResourceData) map[string]interface{} {
	specifications := d.Get("spec").([]interface{})
	if len(specifications) < 1 || specifications[0] == nil {
		return nil
	}

	specification := specifications[0]
	params := map[string]interface{}{
		"resources": buildV2ResourcePoolSpecResources(d),
		"scope": utils.ValueIgnoreEmpty(utils.ExpandToStringListBySet(utils.PathSearch("[0].scope", specification,
			schema.NewSet(schema.HashString, nil)).(*schema.Set))),
	}
	return params
}

func resourceV2ResourcePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		resourcePoolId = d.Id()
	)

	resourcePool, err := queryResourcePool(cfg, region, d)
	if err != nil {
		return diag.Errorf("error querying Modelarts V2 resource pool: %s", err)
	}

	chargingMode := utils.PathSearch(`metadata.annotations."os.modelarts/billing.mode"`, resourcePool, postPaid).(string)
	// 节点中是否有包周期节点，如有，执行批量退订， 没有时再退订资源池
	if chargingMode == prePaid {
		mainResourceId := utils.PathSearch(`metadata.0.labels."os.modelarts/resource.id"`, resourcePool, "").(string)
		if mainResourceId == "" {
			return diag.Errorf("error getting main resource ID from the resource pool(%s)", resourcePoolId)
		}

		if err := common.UnsubscribePrePaidResource(d, cfg, []string{mainResourceId}); err != nil {
			return diag.Errorf("error unsubscribing Modelarts V2 resource pool: %s", err)
		}
	} else {
		err := deleteResourcePool(cfg, d, region)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = deleteResourcePoolWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts V2 resource pool (%s) deletion to complete: %s", resourcePoolId, err)
	}
	return nil
}
