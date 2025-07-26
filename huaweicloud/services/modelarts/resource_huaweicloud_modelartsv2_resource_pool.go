package modelarts

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v2ReourcePoolNonUpdatableParams = []string{
	"spec.*.network",
	"spec.*.network.name",
	"spec.*.network.vpc_id",
	"spec.*.network.subnet_id",
	"spec.*.clusters",
	"spec.*.clusters.provider_id",
	"spec.*.user_login",
	"spec.*.user_login.key_pair_name",
	"spec.*.user_login.password",
}

// @API ModelArts POST /v2/{project_id}/pools
// @API ModelArts GET /v2/{project_id}/pools/{pool_name}
// @API ModelArts PATCH /v2/{project_id}/pools/{pool_name}
// @API ModelArts GET /v2/{project_id}/pools/{pool_name}/nodepools/{nodepool_name}/nodes
// @API ModelArts DELETE /v2/{project_id}/pools/{pool_name}
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
			Delete: schema.DefaultTimeout(90 * time.Minute),
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
							Type:             schema.TypeString,
							Required:         true,
							ValidateFunc:     validation.StringIsJSON,
							DiffSuppressFunc: utils.SuppressObjectDiffs(),
							Description:      `The labels of the resource pool, in JSON format.`,
						},
						"annotations": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
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
				Required: true,
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
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The name of the network.`,
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The ID of the VPC.`,
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The ID of the subnet.`,
									},
								},
							},
							Description: `The network of the resource pool.`,
						},
						"user_login": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_pair_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Sensitive:   true,
										Description: `The name of the key pair.`,
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Sensitive:   true,
										Description: `The password of the resource pool.`,
									},
								},
							},
							Description: `The user login information of the privileged pool.`,
						},
						"clusters": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"provider_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The provider ID of the cluster.`,
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the cluster.`,
									},
								},
							},
							Description: `The cluster information of the privileged pool.`,
						},
						// Internal attributes.
						// "node_pools_order_origin": {
						// 	Type:     schema.TypeList,
						// 	Computed: true,
						// 	Elem:     &schema.Schema{Type: schema.TypeString},
						// 	Description: utils.SchemaDesc(
						// 		`The original value of the node pools order.`,
						// 		utils.SchemaDescInput{
						// 			Internal: true,
						// 		},
						// 	),
						// },
					},
				},
				Description: `The specification of the resource pool.`,
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the resource pool.`,
			},
			// Internal attributes(s).
			"resources_order_origin": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_pool": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(`The node pool of the resource pool.`,
								utils.SchemaDescInput{Internal: true},
							),
						},
						"flavor": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The flavor of the resource pool.`,
								utils.SchemaDescInput{Internal: true},
							),
						},
						"creating_step": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The creating step of the resource pool, in JSON format.`,
								utils.SchemaDescInput{Internal: true},
							),
						},
					},
				},
				Description: utils.SchemaDesc(
					`The script configuration value after the last change, according to which the resources are sorted.`,
					utils.SchemaDescInput{Internal: true},
				),
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
				Description: `The taint list of the resource pool.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The key/value pairs labels of resource pool.`,
			},
			"tags": common.TagsSchema(`The key/value pairs to associate with the resource pool nodes.`),
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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsJSON,
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
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringIsJSON,
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
			// Internal parameters.
			"os": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Description: utils.SchemaDesc(
								`The OS name of the image.`,
								utils.SchemaDescInput{Internal: true},
							),
						},
						"image_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Description: utils.SchemaDesc(
								`The ID of the image.`,
								utils.SchemaDescInput{Internal: true},
							),
						},
						"image_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Description: utils.SchemaDesc(
								`The type of the image.`,
								utils.SchemaDescInput{Internal: true},
							),
						},
					},
				},
				Description: utils.SchemaDesc(
					`The image information for the specified OS.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
			"driver": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Description: utils.SchemaDesc(
								`The driver version.`,
								utils.SchemaDescInput{Internal: true},
							),
						},
					},
				},
				Description: utils.SchemaDesc(
					`The driver information of the resource pool nodes.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
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
		},
	}
}

func createV2ResourcePool(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v2/{project_id}/pools"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCreateV2ResourcePoolBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(createResp)
}

func resourceV2ResourcePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	respBody, err := createV2ResourcePool(client, d)
	if err != nil {
		return diag.Errorf("error creating Modelarts resource pool: %s", err)
	}

	resourcePoolName := utils.PathSearch("metadata.name", respBody, "").(string)
	if resourcePoolName == "" {
		return diag.Errorf("unable to find the resource pool name in the API response")
	}
	d.SetId(resourcePoolName)
	d.Set("resources_order_origin", v2RefreshResourcesOrderOrigin(d.GetRawConfig()))

	if getResourcePoolOrNodesBilingMode(d) == billingModePrePaid {
		orderId := utils.PathSearch(`metadata.annotations."os.modelarts/order.id"`, respBody, "").(string)
		if orderId == "" {
			return diag.Errorf("unable to find the order ID in the API response")
		}

		bssClient, err := cfg.NewServiceClient("bssv2", region)
		if err != nil {
			return diag.Errorf("error creating BSS client: %s", err)
		}

		err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = waitForV2ResourcePoolStateCompleted(ctx, client, resourcePoolName, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts V2 resource pool (%s) creation to complete: %s", d.Id(), err)
	}

	err = waitForV2ReourcePoolDriverStatusCompleted(ctx, client, resourcePoolName, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts V2 resource pool (%s) driver status to become running: %s", d.Id(), err)
	}

	return resourceV2ResourcePoolRead(ctx, d, meta)
}

func getResourcePoolOrNodesBilingMode(d *schema.ResourceData) string {
	annotations, ok := d.GetOk("metadata.0.annotations")
	if ok {
		return utils.PathSearch(`"os.modelarts/billing.mode"`, utils.StringToJson(annotations.(string)), billingModePostPaid).(string)
	}
	return billingModePostPaid
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
		"type": "Dedicate",
		"resources": buildCreateV2ResourcePoolSpecResources(utils.PathSearch("resources", specification,
			make([]interface{}, 0)).([]interface{})),
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

func buildCreateV2ResourcePoolSpecResources(resources []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, len(resources))
	for i, v := range resources {
		result[i] = map[string]interface{}{
			// Required parameters.
			"flavor": utils.PathSearch("flavor", v, nil),
			"count":  utils.PathSearch("count", v, nil),
			// Optional parameters.
			"maxCount": utils.ValueIgnoreEmpty(utils.PathSearch("max_count", v, nil)),
			"nodePool": utils.ValueIgnoreEmpty(utils.PathSearch("node_pool", v, nil)),
			"taints": buildCreateV2ResourcePoolResourceTaints(utils.PathSearch("taints", v,
				schema.NewSet(schema.HashString, nil)).(*schema.Set)),
			"labels": utils.ValueIgnoreEmpty(utils.PathSearch("labels", v, nil)),
			"tags": utils.ValueIgnoreEmpty(utils.ExpandResourceTags(utils.PathSearch("tags", v,
				make(map[string]interface{})).(map[string]interface{}))),
			"network":      buildCreateV2ResourcePoolSpecResourceNetwork(utils.PathSearch("network", v, make([]interface{}, 0)).([]interface{})),
			"extendParams": utils.StringToJson(utils.PathSearch("extend_params", v, "").(string)),
			"creatingStep": buildCreateV2ResourcePoolResourceCreatingStep(utils.PathSearch("creating_step", v,
				make([]interface{}, 0)).([]interface{})),
			"rootVolume": buildCreateV2ResourcePoolResourceRootVolume(utils.PathSearch("root_volume", v,
				make([]interface{}, 0)).([]interface{})),
			"dataVolumes": buildCreateV2ResourcePoolResourceDataVolumes(utils.PathSearch("data_volumes", v,
				make([]interface{}, 0)).([]interface{})),
			"volumeGroupConfigs": buildCreateV2ResourcePoolResourceVolumeGroupConfigs(utils.PathSearch("volume_group_configs", v,
				schema.NewSet(schema.HashString, nil)).(*schema.Set).List(),
			),
			// Internal parameters.
			"os": utils.ValueIgnoreEmpty(buildCreateV2ResourcePoolResourceOsInfo(utils.PathSearch("os", v,
				make([]interface{}, 0)).([]interface{}))),
			"driver": utils.ValueIgnoreEmpty(buildCreateV2ResourcePoolResourceDriver(utils.PathSearch("driver", v,
				make([]interface{}, 0)).([]interface{}))),
			"azs": buildCreateV2ResourcePoolResourceAzs(utils.PathSearch("azs", v,
				schema.NewSet(schema.HashString, nil)).(*schema.Set)),
		}
	}
	return result
}

func buildCreateV2ResourcePoolResourceTaints(taints *schema.Set) []map[string]interface{} {
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

func buildCreateV2ResourcePoolSpecResourceNetwork(resourceNetworks []interface{}) map[string]interface{} {
	if len(resourceNetworks) < 1 || resourceNetworks[0] == nil {
		return nil
	}

	// All parameters are as the Computed behavior.
	return utils.RemoveNil(map[string]interface{}{
		"vpc":    utils.ValueIgnoreEmpty(utils.PathSearch("[0].vpc", resourceNetworks, nil)),
		"subnet": utils.ValueIgnoreEmpty(utils.PathSearch("[0].subnet", resourceNetworks, nil)),
		"securityGroups": utils.ValueIgnoreEmpty(utils.ExpandToStringListBySet(utils.PathSearch("[0].security_groups",
			resourceNetworks, schema.NewSet(schema.HashString, nil)).(*schema.Set))),
	})
}

func buildCreateV2ResourcePoolResourceCreatingStep(creatingSteps []interface{}) map[string]interface{} {
	if len(creatingSteps) < 1 {
		return nil
	}

	return map[string]interface{}{
		"type": utils.ValueIgnoreEmpty(utils.PathSearch("type", creatingSteps[0], nil)),
		"step": utils.ValueIgnoreEmpty(utils.PathSearch("step", creatingSteps[0], nil)),
	}
}

func buildCreateV2ResourcePoolResourceRootVolume(rootVolumes []interface{}) map[string]interface{} {
	if len(rootVolumes) < 1 || rootVolumes[0] == nil {
		return nil
	}

	return utils.RemoveNil(map[string]interface{}{
		"volumeType": utils.PathSearch("[0].volume_type", rootVolumes, nil),
		"size":       utils.PathSearch("[0].size", rootVolumes, nil),
	})
}

func buildCreateV2ResourcePoolResourceDataVolumes(dataVolumes []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(dataVolumes))
	for i, dataVolume := range dataVolumes {
		result[i] = map[string]interface{}{
			"volumeType":   utils.PathSearch("volume_type", dataVolume, nil),
			"size":         utils.PathSearch("size", dataVolume, nil),
			"extendParams": utils.StringToJson(utils.PathSearch("extend_params", dataVolume, "").(string)),
			"count":        utils.ValueIgnoreEmpty(utils.PathSearch("count", dataVolume, nil)),
		}
	}

	return result
}

func buildCreateV2ResourcePoolResourceVolumeGroupConfigs(volumeGroupConfigs []interface{}) []map[string]interface{} {
	if len(volumeGroupConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(volumeGroupConfigs))
	for i, volumeGroupConfig := range volumeGroupConfigs {
		result[i] = map[string]interface{}{
			"volumeGroup":    utils.PathSearch("volume_group", volumeGroupConfig, nil),
			"dockerThinPool": utils.ValueIgnoreEmpty(utils.PathSearch("docker_thin_pool", volumeGroupConfig, nil)),
			"lvmConfig": buildCreateV2ResourcePoolResourceVolumeGroupConfigsLvmConfig(utils.PathSearch("lvm_config", volumeGroupConfig,
				make([]interface{}, 0)).([]interface{})),
			"types": utils.ValueIgnoreEmpty(utils.ExpandToStringListBySet(utils.PathSearch("types", volumeGroupConfig,
				schema.NewSet(schema.HashString, nil)).(*schema.Set))),
		}
	}
	return result
}

func buildCreateV2ResourcePoolResourceVolumeGroupConfigsLvmConfig(lvmConfigs []interface{}) map[string]interface{} {
	if len(lvmConfigs) < 1 {
		return nil
	}

	lvmConfig := lvmConfigs[0]
	return map[string]interface{}{
		"lvType": utils.PathSearch("lv_type", lvmConfig, nil),
		"path":   utils.ValueIgnoreEmpty(utils.PathSearch("path", lvmConfig, nil)),
	}
}

func buildCreateV2ResourcePoolResourceOsInfo(osInfos []interface{}) map[string]interface{} {
	if len(osInfos) < 1 || osInfos[0] == nil {
		return nil
	}

	return utils.RemoveNil(map[string]interface{}{
		"name":      utils.ValueIgnoreEmpty(utils.PathSearch("[0].name", osInfos, nil)),
		"imageId":   utils.ValueIgnoreEmpty(utils.PathSearch("[0].image_id", osInfos, nil)),
		"imageType": utils.ValueIgnoreEmpty(utils.PathSearch("[0].image_type", osInfos, nil)),
	})
}

func buildCreateV2ResourcePoolResourceDriver(drivers []interface{}) map[string]interface{} {
	if len(drivers) < 1 || drivers[0] == nil {
		return nil
	}

	return utils.RemoveNil(map[string]interface{}{
		"version": utils.ValueIgnoreEmpty(utils.PathSearch("version", drivers[0], nil)),
	})
}

func buildCreateV2ResourcePoolResourceAzs(azs *schema.Set) []map[string]interface{} {
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

func buildCreateV2ResourcePoolSpecNetwork(networks []interface{}) map[string]interface{} {
	// All parameters are as the optional behavior.
	if len(networks) < 1 || networks[0] == nil {
		return nil
	}

	// All parameters are as the computed behavior.
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

	return map[string]interface{}{
		"keyPairName": utils.ValueIgnoreEmpty(utils.PathSearch("[0].key_pair_name", userLogin, nil)),
		"password":    utils.ValueIgnoreEmpty(utils.PathSearch("[0].password", userLogin, nil)),
	}
}

func buildCreateV2ResourcePoolSpecClusters(clusters []interface{}) []map[string]interface{} {
	if len(clusters) < 1 || clusters[0] == nil {
		return nil
	}

	result := make([]map[string]interface{}, len(clusters))
	for i, v := range clusters {
		result[i] = map[string]interface{}{
			"name":       utils.ValueIgnoreEmpty(utils.PathSearch("name", v, nil)),
			"providerId": utils.ValueIgnoreEmpty(utils.PathSearch("provider_id", v, nil)),
		}
	}
	return result
}

func waitForV2ResourcePoolStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, resourcePoolName string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resourcePool, err := getV2ResourcePoolByName(client, resourcePoolName)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status.phase", resourcePool, "").(string)
			if status == "Error" {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `status.phase`)
			}

			// The `status.resources` list means resource status under the resource pool.
			creating := utils.PathSearch("status.resources.creating", resourcePool, make([]interface{}, 0)).([]interface{})
			deleting := utils.PathSearch("status.resources.deleting", resourcePool, make([]interface{}, 0)).([]interface{})
			creationFailed := utils.PathSearch("status.resources.creationFailed", resourcePool, make([]interface{}, 0)).([]interface{})
			// Check whether the resource pool is running and the node expansion and contraction is successful.
			if status == "Running" && len(creating) == 0 && len(deleting) == 0 && len(creationFailed) == 0 {
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

func waitForV2ReourcePoolDriverStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, resourcePoolName string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      v2ResourcePoolDriverStatusRefreshFunc(client, resourcePoolName),
		Timeout:      timeout,
		PollInterval: 20 * time.Second,
		// In some cases, the following status changes may occur: Upgrading -> Running -> Creating -> Running
		ContinuousTargetOccurence: 2,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func v2ResourcePoolDriverStatusRefreshFunc(client *golangsdk.ServiceClient, resourcePoolName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resourcePool, err := getV2ResourcePoolByName(client, resourcePoolName)
		if err != nil {
			return resourcePool, "ERROR", err
		}

		driverStatuses := utils.PathSearch("status.driver.*.state", resourcePool, make([]interface{}, 0)).([]interface{})
		if len(driverStatuses) == 0 {
			return "No matches found", "COMPLETED", nil
		}

		for _, status := range driverStatuses {
			if utils.StrSliceContains([]string{"Running", "Abnormal"}, status.(string)) {
				return resourcePool, "COMPLETED", nil
			}
		}

		return resourcePool, "PENDING", nil
	}
}

func resourceV2ResourcePoolRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		resourcePoolName = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resourcePool, err := getV2ResourcePoolByName(client, resourcePoolName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Modelarts V2 resource pool")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("metadata", flattenV2ResourcePoolMetadataInfo(d, utils.PathSearch("metadata", resourcePool, nil))),
		d.Set("spec", flattenV2ResourcePoolSpecification(utils.PathSearch("spec", resourcePool, nil), d)),
		// Attributes.
		d.Set("status", utils.PathSearch("status.phase", resourcePool, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func getV2ResourcePoolByName(client *golangsdk.ServiceClient, resourcePoolName string) (interface{}, error) {
	getHttpUrl := "v2/{project_id}/pools/{pool_name}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{pool_name}", resourcePoolName)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func flattenV2ResourcePoolMetadataInfo(d *schema.ResourceData, metadata interface{}) []map[string]interface{} {
	if metadata == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"annotations": d.Get("metadata.0.annotations"),
			"labels":      utils.JsonToString(utils.PathSearch("labels", metadata, nil)),
			// Attribute(s).
			"labels_origin": d.Get("metadata.0.labels_origin"),
		},
	}
}

func flattenV2ResourcePoolSpecification(spec interface{}, d *schema.ResourceData) []map[string]interface{} {
	if spec == nil {
		return nil
	}

	sortedResources := flattenV2ResourcePoolSpecResourcesAndNodePoolsOrderOrigin(
		utils.PathSearch("resources", spec, make([]interface{}, 0)).([]interface{}),
		d.Get("resources_order_origin").([]interface{}),
	)

	return []map[string]interface{}{
		{
			"resources":  flattenV2ResourcePoolSpecResources(sortedResources),
			"scope":      utils.PathSearch("scope", spec, nil),
			"network":    flattenV2ResourcePoolSpecNetwork(utils.PathSearch("network", spec, nil)),
			"user_login": flattenV2ResourcePoolSpecUserLogin(utils.PathSearch("userLogin", spec, nil)),
			"clusters":   flattenV2ResourcePoolClusters(utils.PathSearch("clusters", spec, make([]interface{}, 0)).([]interface{})),
		},
	}
}

// This method is used to sort the `resources`.
func flattenV2ResourcePoolSpecResourcesAndNodePoolsOrderOrigin(resources, resourcesOrderOrigin []interface{}) []interface{} {
	if len(resourcesOrderOrigin) == 0 {
		return resources
	}

	sortedResources := make([]interface{}, 0)
	// According to the `resources_order_origin` to sort the `resources.
	for _, v := range resourcesOrderOrigin {
		_, index := v2FindResourceByFlavorAndNodePoolAndCreatingStep(
			resources,
			utils.PathSearch("flavor", v, "").(string),
			utils.PathSearch("node_pool", v, "").(string),
			utils.PathSearch("creating_step", v, "").(string),
		)
		if index == -1 {
			continue
		}
		sortedResources = append(sortedResources, resources[index])
		resources = append(resources[:index], resources[index+1:]...)
	}

	sortedResources = append(sortedResources, resources...)
	return sortedResources
}

func flattenV2ResourcePoolSpecResources(resources []interface{}) []interface{} {
	if len(resources) < 1 {
		return nil
	}

	result := make([]interface{}, len(resources))
	for i, v := range resources {
		result[i] = map[string]interface{}{
			"flavor":        utils.PathSearch("flavor", v, nil),
			"count":         utils.PathSearch("count", v, nil),
			"max_count":     utils.PathSearch("maxCount", v, nil),
			"azs":           flattenV2ResourcePoolResourceAzs(utils.PathSearch("azs", v, make([]interface{}, 0)).([]interface{})),
			"node_pool":     utils.PathSearch("nodePool", v, nil),
			"taints":        flattenV2ResourcePoolResourceTaints(utils.PathSearch("taints", v, make([]interface{}, 0)).([]interface{})),
			"labels":        utils.PathSearch("labels", v, nil),
			"tags":          utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil)),
			"network":       flattenV2ResourcePoolSpecResourcesNetwork(utils.PathSearch("network", v, nil)),
			"extend_params": utils.JsonToString(utils.PathSearch("extendParams", v, nil)),
			"creating_step": flattenV2ResourcePoolResourcesCreatingStep(utils.PathSearch("creatingStep", v, nil)),
			"root_volume":   flattenV2ResourcePoolResourcesRootVolume(utils.PathSearch("rootVolume", v, nil)),
			"data_volumes": flattenV2ResourcePoolResourcesDataVolumes(utils.PathSearch("dataVolumes",
				v, make([]interface{}, 0)).([]interface{})),
			"volume_group_configs": flattenV2ResourcePoolResourcesVolumeGroupConfigs(utils.PathSearch("volumeGroupConfigs",
				v, make([]interface{}, 0)).([]interface{})),
			"volume_group_configs_origin": flattenV2ResourcePoolResourcesVolumeGroupConfigs(utils.PathSearch("volumeGroupConfigs",
				v, make([]interface{}, 0)).([]interface{})),
			"os":     flattenV2ResourcePoolResourcesOs(utils.PathSearch("os", v, nil)),
			"driver": flattenV2ResourcePoolResourcesDriver(utils.PathSearch("driver", v, nil)),
		}
	}

	return result
}

func flattenV2ResourcePoolSpecResourcesNetwork(network interface{}) []map[string]interface{} {
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

func flattenV2ResourcePoolResourcesCreatingStep(creatingStep interface{}) []map[string]interface{} {
	if creatingStep == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"step": utils.PathSearch("step", creatingStep, nil),
			"type": utils.PathSearch("type", creatingStep, nil),
		},
	}
}

func flattenV2ResourcePoolResourcesRootVolume(rootVolume interface{}) []map[string]interface{} {
	if rootVolume == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"volume_type": utils.PathSearch("volumeType", rootVolume, nil),
			"size":        utils.PathSearch("size", rootVolume, nil),
		},
	}
}

func flattenV2ResourcePoolResourcesDataVolumes(dataVolumes []interface{}) []map[string]interface{} {
	if len(dataVolumes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dataVolumes))
	for _, dataVolume := range dataVolumes {
		result = append(result, map[string]interface{}{
			"volume_type":   utils.PathSearch("volumeType", dataVolume, nil),
			"size":          utils.PathSearch("size", dataVolume, nil),
			"extend_params": utils.JsonToString(utils.PathSearch("extendParams", dataVolume, nil)),
			"count":         utils.PathSearch("count", dataVolume, nil),
		})
	}

	return result
}

func flattenV2ResourcePoolResourcesVolumeGroupConfigs(volumeGroupConfigs []interface{}) []map[string]interface{} {
	if len(volumeGroupConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(volumeGroupConfigs))
	for _, volumeGroupConfig := range volumeGroupConfigs {
		result = append(result, map[string]interface{}{
			"volume_group":     utils.PathSearch("volumeGroup", volumeGroupConfig, nil),
			"docker_thin_pool": utils.PathSearch("dockerThinPool", volumeGroupConfig, nil),
			"lvm_config": flattenV2ResourcePoolVolumeGroupConfigsLvmConfig(utils.PathSearch("lvmConfig",
				volumeGroupConfig, nil)),
			"types": utils.PathSearch("types", volumeGroupConfig, make([]interface{}, 0)),
		})
	}
	return result
}

func flattenV2ResourcePoolVolumeGroupConfigsLvmConfig(lvmConfig interface{}) []map[string]interface{} {
	if lvmConfig == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"lv_type": utils.PathSearch("lvType", lvmConfig, nil),
			"path":    utils.PathSearch("path", lvmConfig, nil),
		},
	}
}

func flattenV2ResourcePoolResourcesOs(osInfo interface{}) []map[string]interface{} {
	if osInfo == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":       utils.PathSearch("name", osInfo, nil),
			"image_id":   utils.PathSearch("imageId", osInfo, nil),
			"image_type": utils.PathSearch("iamgeType", osInfo, nil),
		},
	}
}

func flattenV2ResourcePoolResourcesDriver(driver interface{}) []map[string]interface{} {
	if driver == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"version": utils.PathSearch("version", driver, nil),
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
	if keyPairName == "" {
		return nil
	}

	return []map[string]interface{}{
		{
			"key_pair_name": keyPairName,
		},
	}
}

func resourceV2ResourcePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		resourcePoolName = d.Id()
	)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	updateRespBody, err := updateV2ResourcePool(client, d, resourcePoolName)
	if err != nil {
		return diag.Errorf("error updating Modelarts V2 resource pool: %s", err)
	}

	d.Set("resources_order_origin", v2RefreshResourcesOrderOrigin(d.GetRawConfig()))

	// Only when expanding prepaid type nodes, we need to determine the order status.
	// When expanding the nodes, getResourcePoolOrNodesBilingMode result means the charging mode of the nodes.
	if getResourcePoolOrNodesBilingMode(d) == billingModePrePaid && d.HasChange("spec.0.resources") {
		// Whenever any count in resouces changes, the order status needs to be determined.
		oldRaw, newRaw := d.GetChange("spec.0.resources")
		if isV2ResourcePoolAnyNodeScalling(oldRaw.([]interface{}), newRaw.([]interface{})) {
			orderId := utils.PathSearch(`metadata.annotations."os.modelarts/order.id"`, updateRespBody, nil)
			if orderId == nil {
				return diag.Errorf("error updating Modelarts resource pool: order ID is not found in API response")
			}

			bssClient, err := cfg.NewServiceClient("bssv2", region)
			if err != nil {
				return diag.Errorf("error creating BSS client: %s", err)
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

	if d.HasChange("spec.0.resources") {
		err = waitForV2ResourcePoolStateCompleted(ctx, client, resourcePoolName, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the Modelarts V2 resource pool (%s) update to complete: %s", d.Id(), err)
		}

		err = waitForV2ReourcePoolDriverStatusCompleted(ctx, client, resourcePoolName, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for the Modelarts V2 resource pool (%s) driver status to become running: %s", d.Id(), err)
		}
	}

	return resourceV2ResourcePoolRead(ctx, d, meta)
}

func updateV2ResourcePool(client *golangsdk.ServiceClient, d *schema.ResourceData, resourcePoolName string) (interface{}, error) {
	updatehttpUrl := "v2/{project_id}/pools/{pool_name}"
	updatePath := client.Endpoint + updatehttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{pool_name}", resourcePoolName)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/merge-patch+json"},
		JSONBody:         utils.RemoveNil(buildUpdateV2ResourcePoolBodyParams(d)),
	}

	resp, err := client.Request("PATC", updatePath, &updateOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func isV2ResourcePoolAnyNodeScalling(oldResource, newResource []interface{}) bool {
	for i, v := range newResource {
		oldCount := utils.PathSearch(fmt.Sprintf("[%d].count", i), oldResource, 0).(int)
		newCount := utils.PathSearch("count", v, 0).(int)
		if oldCount != newCount {
			return true
		}
	}
	return false
}

func isAnyNodePoolCountIncrease(oldResource, newRawResource []interface{}) bool {
	for i, v := range newRawResource {
		// Node pool not change, only count increase.
		oldCount := utils.PathSearch(fmt.Sprintf("[%d].count", i), oldResource, -1).(int)
		newCount := utils.PathSearch("count", v, 0).(int)
		if newCount > oldCount {
			return true
		}
	}
	return false
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
		"labels":      buildUpdateV2ResourcePoolMetaDataLabels(utils.StringToJson(d.Get("metadata.0.labels").(string)).(map[string]interface{})),
		"annotations": buildUpdateV2ResourcePoolMetaDataAnnotations(d),
	}
	return params
}

func buildUpdateV2ResourcePoolMetaDataLabels(labels map[string]interface{}) map[string]interface{} {
	// The allowed keys are the keys that can be updated.
	allowedKeys := []string{
		"os.modelarts/workspace.id",
	}
	return filterAllowedParamsFromMap(labels, allowedKeys)
}

func filterAllowedParamsFromMap(params map[string]interface{}, allowedKeys []string) map[string]interface{} {
	for key := range params {
		if !utils.StrSliceContains(allowedKeys, key) {
			delete(params, key)
		}
	}
	return params
}

func buildUpdateV2ResourcePoolMetaDataAnnotations(d *schema.ResourceData) map[string]interface{} {
	sharedKeys := []string{
		"os.modelarts/billing.mode",
		"os.modelarts/period.num",
		"os.modelarts/period.type",
		"os.modelarts/auto.renew",
		"os.modelarts/order.id",
		"os.modelarts/auto.pay",
		"os.modelarts/promotion.info",
	}
	params := make(map[string]interface{})
	if annotations, ok := d.GetOk("metadata.0.annotations"); ok {
		allowedKeys := []string{"os.modelarts/description", "os.modelarts.pool/drain.enabled"}
		allowedKeys = append(allowedKeys, sharedKeys...)
		params = filterAllowedParamsFromMap(utils.StringToJson(annotations.(string)).(map[string]interface{}), allowedKeys)
	}

	oldRaw, newRaw := d.GetChange("spec.0.resources")
	// If the node pools or nodes are not increased, delete the billing mode related parameters.
	if !isAnyNodePoolCountIncrease(oldRaw.([]interface{}), newRaw.([]interface{})) {
		params = deleteNotAllowedKeysFromMap(params, sharedKeys)
	}
	return params
}

func deleteNotAllowedKeysFromMap(params map[string]interface{}, keys []string) map[string]interface{} {
	for key := range params {
		if utils.StrSliceContains(keys, key) {
			delete(params, key)
		}
	}
	return params
}

func buildUpdateV2ResourcePoolSpec(d *schema.ResourceData) map[string]interface{} {
	specifications := d.Get("spec").([]interface{})
	if len(specifications) < 1 || specifications[0] == nil {
		return nil
	}

	oldResources, _ := d.GetChange("spec.0.resources")
	params := map[string]interface{}{
		"resources": buildUpdateV2ResourcePoolSpecResources(d.GetRawConfig(), oldResources.([]interface{})),
		"scope": utils.ValueIgnoreEmpty(utils.ExpandToStringListBySet(utils.PathSearch("[0].scope", specifications,
			schema.NewSet(schema.HashString, nil)).(*schema.Set))),
	}
	return params
}

func buildUpdateV2ResourcePoolSpecResources(rawConfig cty.Value, oldResources []interface{}) []map[string]interface{} {
	resources := v2GetConfigFileSpecResources(rawConfig)
	if resources == nil {
		return nil
	}

	newResources := resources.(cty.Value)
	result := make([]map[string]interface{}, newResources.LengthInt())
	for i, newResource := range newResources.AsValueSlice() {
		matchedOldResource := v2GetMatchedResourceFromConfigfile(newResource, oldResources)
		result[i] = map[string]interface{}{
			// Required parameters.
			"flavor": v2GetConfigFileStringValueByKey(newResource, "flavor"),
			"count":  v2GetConfigFileIntValueByKey(newResource, "count"),
			// Only optional parameters.
			"tags": buildUpdateV2ResourcePoolResourceTags(newResource),
			// Computed parameters.
			"nodePool": v2GetOldResourceStringValueByConfigFileKey(newResource, "node_pool", matchedOldResource),
			"maxCount": v2GetOldResourceIntValueByConfigFileKey(newResource, "max_count", matchedOldResource),
			"taints":   buildUpdateV2ResourcePoolResourceTaints(newResource, matchedOldResource),
			"labels":   buildUpdateV2ResourcePoolResourceLabels(newResource, "labels", matchedOldResource),
			"network":  utils.ValueIgnoreEmpty(buildUpdateV2ResourcePoolResourceNetwork(newResource, matchedOldResource)),
			"extendParams": buildV2ResourcePoolResourcesExtendParamsBodyParams(
				utils.PathSearch("extend_params", matchedOldResource, "{}").(string),
				v2GetConfigFileStringValueByKey(newResource, "extend_params"),
			),
			"creatingStep":       buildUpdateV2ResourcePoolResourceCreatingStep(newResource, matchedOldResource),
			"rootVolume":         buildUpdateV2ResourcePoolSpecResourceRootVolume(newResource, matchedOldResource),
			"dataVolumes":        buildUpdateV2ResourcePoolSpecResourceDataVolumes(newResource, matchedOldResource),
			"volumeGroupConfigs": buildUpdateV2ResourcePoolResourceVolumeGroupConfigs(newResource, matchedOldResource),
			// Internal computed parameters.
			"os":     utils.ValueIgnoreEmpty(buildUpdateV2ResourcePoolResourceOs(newResource, matchedOldResource)),
			"driver": utils.ValueIgnoreEmpty(buildUpdateV2ResourcePoolResourceDriver(newResource, matchedOldResource)),
			"azs":    buildUpdateV2ResourcePoolResourceAzs(newResource, matchedOldResource),
		}
	}
	return result
}

func buildUpdateV2ResourcePoolResourceTags(resourceElem cty.Value) interface{} {
	raw := v2GetConfigFileMapValueByKey(resourceElem, "tags")
	if raw == nil {
		return nil
	}

	result := make([]map[string]interface{}, 0)
	for k, v := range raw.(cty.Value).AsValueMap() {
		tagMap := map[string]interface{}{
			"key": k,
		}

		if v.Type() == cty.String && !v.IsNull() && v.IsKnown() {
			tagMap["value"] = v.AsString()
		}

		result = append(result, tagMap)
	}
	return result
}

func buildUpdateV2ResourcePoolResourceTaints(resourceElem cty.Value, oldResource interface{}) interface{} {
	raw := v2GetConfigFileListValueByKey(resourceElem, "taints")
	if raw == nil {
		return buildCreateV2ResourcePoolResourceTaints(utils.PathSearch("taints", oldResource,
			schema.NewSet(schema.HashString, nil)).(*schema.Set))
	}

	taints := raw.(cty.Value)
	result := make([]map[string]interface{}, taints.LengthInt())
	for i, taint := range taints.AsValueSlice() {
		result[i] = map[string]interface{}{
			"key":    taint.GetAttr("key").AsString(),
			"effect": taint.GetAttr("effect").AsString(),
			"value":  utils.ValueIgnoreEmpty(v2GetConfigFileStringValueByKey(taint, "value")),
		}
	}
	return result
}

func buildUpdateV2ResourcePoolResourceLabels(elem cty.Value, key string, oldResource interface{}) interface{} {
	raw := v2GetConfigFileMapValueByKey(elem, key)
	if raw == nil {
		return utils.ValueIgnoreEmpty(utils.PathSearch(key, oldResource, nil))
	}

	labels := make(map[string]interface{})
	for k, v := range raw.(cty.Value).AsValueMap() {
		labels[k] = v.AsString()
	}
	return utils.ValueIgnoreEmpty(labels)
}

func buildV2ResourcePoolResourcesExtendParamsBodyParams(oldExtendParams, newExtendParams string) map[string]interface{} {
	extendParams := utils.TryMapValueAnalysis(utils.StringToJson(oldExtendParams))
	if objExtendParams := utils.TryMapValueAnalysis(utils.StringToJson(newExtendParams)); len(objExtendParams) > 0 {
		for k, v := range objExtendParams {
			extendParams[k] = v
		}
	}
	return extendParams
}

func buildUpdateV2ResourcePoolResourceNetwork(resourceElem cty.Value, oldResource interface{}) interface{} {
	raw := v2GetConfigFileListValueByKey(resourceElem, "network")
	if raw == nil {
		return buildCreateV2ResourcePoolSpecResourceNetwork(utils.PathSearch("network", oldResource, make([]interface{}, 0)).([]interface{}))
	}

	network := raw.(cty.Value)
	return utils.RemoveNil(map[string]interface{}{
		// 待测试： 如果不指定vpc_id时，使用network.Index(cty.NumberIntVal(0)取值是否会报错
		"vpc":            utils.ValueIgnoreEmpty(v2GetConfigFileStringValueByKey(network.Index(cty.NumberIntVal(0)), "vpc")),
		"subnet":         utils.ValueIgnoreEmpty(v2GetConfigFileStringValueByKey(network.Index(cty.NumberIntVal(0)), "subnet")),
		"securityGroups": buildV2UpdateSpecResourceNetworkSecurityGroups(resourceElem, "security_groups", oldResource),
	})
}

func buildV2UpdateSpecResourceNetworkSecurityGroups(elem cty.Value, key string, oldResource interface{}) interface{} {
	raw := v2GetRawConfigSetValueByKey(elem, key)
	if raw == nil {
		return utils.ValueIgnoreEmpty(utils.ExpandToStringListBySet(utils.PathSearch(key, oldResource,
			schema.NewSet(schema.HashString, nil)).(*schema.Set)))
	}

	securityGroupIds := raw.(cty.Value)
	results := make([]string, securityGroupIds.LengthInt())
	for i, securityGroupId := range securityGroupIds.AsValueSlice() {
		results[i] = securityGroupId.AsString()
	}
	return utils.ValueIgnoreEmpty(securityGroupIds)
}

func buildUpdateV2ResourcePoolResourceCreatingStep(resourceElem cty.Value, oldResource interface{}) interface{} {
	raw := v2GetConfigFileListValueByKey(resourceElem, "creating_step")
	if raw == nil {
		return buildCreateV2ResourcePoolResourceCreatingStep(utils.PathSearch("creating_step", oldResource, make([]interface{}, 0)).([]interface{}))
	}

	creatingStep := raw.(cty.Value)
	return map[string]interface{}{
		"type": creatingStep.Index(cty.NumberIntVal(0)).GetAttr("type").AsString(),
		"step": v2GetConfigFileIntValueByKey(creatingStep.Index(cty.NumberIntVal(0)), "step"),
	}
}

func buildUpdateV2ResourcePoolSpecResourceRootVolume(resourceElem cty.Value, oldResource interface{}) map[string]interface{} {
	raw := v2GetConfigFileListValueByKey(resourceElem, "root_volume")
	if raw == nil {
		return buildCreateV2ResourcePoolResourceRootVolume(utils.PathSearch("root_volume", oldResource,
			make([]interface{}, 0)).([]interface{}))
	}

	rootVolume := raw.(cty.Value)
	return map[string]interface{}{
		"volumeType": v2GetConfigFileStringValueByKey(rootVolume.Index(cty.NumberIntVal(0)), "volume_type"),
		"size":       v2GetConfigFileStringValueByKey(rootVolume.Index(cty.NumberIntVal(0)), "size"),
	}
}

func buildUpdateV2ResourcePoolSpecResourceDataVolumes(resourceElem cty.Value, oldResource interface{}) []map[string]interface{} {
	raw := v2GetConfigFileListValueByKey(resourceElem, "data_volumes")
	if raw == nil {
		return buildCreateV2ResourcePoolResourceDataVolumes(utils.PathSearch("data_volumes", oldResource,
			make([]interface{}, 0)).([]interface{}))
	}

	dataVolume := raw.(cty.Value)
	result := make([]map[string]interface{}, dataVolume.LengthInt())
	for i, volume := range dataVolume.AsValueSlice() {
		result[i] = map[string]interface{}{
			"volumeType": v2GetConfigFileStringValueByKey(volume, "volume_type"),
			"size":       v2GetConfigFileStringValueByKey(volume, "size"),
			"extendParams": buildV2ResourcePoolResourcesExtendParamsBodyParams(
				utils.PathSearch(fmt.Sprintf("[%d].extend_params", i), oldResource, "").(string),
				v2GetConfigFileStringValueByKey(dataVolume, "extend_params"),
			),
			"count": utils.ValueIgnoreEmpty(v2GetConfigFileIntValueByKey(volume, "count")),
		}
	}

	return result
}

func buildUpdateV2ResourcePoolResourceVolumeGroupConfigs(resourceElem cty.Value, oldResource interface{}) interface{} {
	raw := v2GetRawConfigSetValueByKey(resourceElem, "volume_group_configs")
	oldVolumeGroupConfigs := utils.PathSearch("volume_group_configs", oldResource, schema.NewSet(schema.HashString, nil)).(*schema.Set).List()
	if raw == nil {
		return buildCreateV2ResourcePoolResourceVolumeGroupConfigs(oldVolumeGroupConfigs)
	}

	newVolumeGroupConfigs := buildUpdateResourcePoolVolumeGroupConfigs(raw.(cty.Value))
	result := make([]interface{}, len(oldVolumeGroupConfigs))
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

		if types := utils.PathSearch("types", newVolumeGroupConfig, make([]interface{}, 0)).([]interface{}); len(types) > 0 {
			elem["types"] = utils.ValueIgnoreEmpty(types)
		} else {
			elem["types"] = utils.ValueIgnoreEmpty(utils.PathSearch("types", volumeGroupConfig, make([]interface{}, 0)).([]interface{}))
		}

		result = append(result, elem)
	}
	return result
}

func buildUpdateResourcePoolVolumeGroupConfigs(volumeGroupConfigs cty.Value) []map[string]interface{} {
	result := make([]map[string]interface{}, volumeGroupConfigs.LengthInt())
	for i, volumeGroupConfigElem := range volumeGroupConfigs.AsValueSlice() {
		result[i] = map[string]interface{}{
			"volume_group":     volumeGroupConfigElem.GetAttr("volume_group").AsString(),
			"docker_thin_pool": v2GetConfigFileIntValueByKey(volumeGroupConfigElem, "docker_thin_pool"),
			"lvm_config":       buildV2UpdateResourceVolumeGroupConfigsLvmConfig(volumeGroupConfigElem),
			"types":            buildV2UpdateResourcePoolResourceVolumeGroupConfigsTypes(volumeGroupConfigElem),
		}
	}
	return result
}

func buildV2UpdateResourceVolumeGroupConfigsLvmConfig(resourceElem cty.Value) []interface{} {
	raw := v2GetConfigFileListValueByKey(resourceElem, "lvm_config")
	if raw == nil {
		return nil
	}

	lvmConfig := raw.(cty.Value)
	return []interface{}{
		map[string]interface{}{
			"lv_type": lvmConfig.Index(cty.NumberIntVal(0)).GetAttr("lv_type").AsString(),
			"path":    v2GetConfigFileStringValueByKey(lvmConfig.Index(cty.NumberIntVal(0)), "path"),
		},
	}
}

func buildV2UpdateResourcePoolResourceVolumeGroupConfigsTypes(elem cty.Value) []interface{} {
	raw := v2GetConfigFileListValueByKey(elem, "types")
	if raw == nil {
		return nil
	}

	types := raw.(cty.Value)
	result := make([]interface{}, types.LengthInt())
	for i, typeElem := range types.AsValueSlice() {
		result[i] = typeElem.AsString()
	}
	return result
}

func buildUpdateV2ResourcePoolResourceOs(resourceElem cty.Value, oldResource interface{}) interface{} {
	raw := v2GetConfigFileListValueByKey(resourceElem, "os")
	if raw == nil {
		return buildCreateV2ResourcePoolResourceOsInfo(utils.PathSearch("os", oldResource, make([]interface{}, 0)).([]interface{}))
	}

	os := raw.(cty.Value)
	return utils.RemoveNil(map[string]interface{}{
		"name":      utils.ValueIgnoreEmpty(getConfigFileStringValueByKey(os.Index(cty.NumberIntVal(0)), "name")),
		"imageId":   utils.ValueIgnoreEmpty(getConfigFileStringValueByKey(os.Index(cty.NumberIntVal(0)), "image_id")),
		"imageType": utils.ValueIgnoreEmpty(getConfigFileStringValueByKey(os.Index(cty.NumberIntVal(0)), "image_type")),
	})
}

func buildUpdateV2ResourcePoolResourceDriver(resourceElem cty.Value, oldResource interface{}) interface{} {
	raw := v2GetConfigFileListValueByKey(resourceElem, "driver")
	if raw == nil {
		return buildCreateV2ResourcePoolResourceDriver(utils.PathSearch("driver", oldResource, make([]interface{}, 0)).([]interface{}))
	}

	driver := raw.(cty.Value)
	return utils.RemoveNil(map[string]interface{}{
		"version": utils.ValueIgnoreEmpty(getConfigFileStringValueByKey(driver.Index(cty.NumberIntVal(0)), "version")),
	})
}

func buildUpdateV2ResourcePoolResourceAzs(resourceElem cty.Value, oldResource interface{}) interface{} {
	raw := v2GetRawConfigSetValueByKey(resourceElem, "azs")
	if raw == nil {
		return buildCreateV2ResourcePoolResourceAzs(utils.PathSearch("azs", oldResource, schema.NewSet(schema.HashString, nil)).(*schema.Set))
	}

	azs := raw.(cty.Value)
	result := make([]map[string]interface{}, azs.LengthInt())
	for i, az := range azs.AsValueSlice() {
		result[i] = map[string]interface{}{
			"az":    v2GetConfigFileStringValueByKey(az, "az"),
			"count": v2GetConfigFileIntValueByKey(az, "count"),
		}
	}
	return result
}

func v2GetConfigFileSpecResources(rawConfig cty.Value) interface{} {
	if rawConfig.IsNull() || !rawConfig.IsKnown() || !rawConfig.Type().IsObjectType() {
		return nil
	}

	raw := v2GetConfigFileListValueByKey(rawConfig, "spec")
	if raw == nil {
		return nil
	}

	specAttr := raw.(cty.Value)
	return specAttr.Index(cty.NumberIntVal(0)).GetAttr("resources")
}

// Get the list value by the key from the config file.
func v2GetConfigFileListValueByKey(elem cty.Value, key string) interface{} {
	if !elem.Type().HasAttribute(key) {
		return nil
	}

	raw := elem.GetAttr(key)
	if raw.IsNull() || !raw.IsKnown() || !raw.Type().IsListType() {
		return nil
	}

	return raw
}

// Get the set value by the key from the config file.
func v2GetRawConfigSetValueByKey(elem cty.Value, key string) interface{} {
	if !elem.Type().HasAttribute(key) {
		return nil
	}

	raw := elem.GetAttr(key)
	if raw.IsNull() || !raw.IsKnown() || !raw.Type().IsSetType() {
		return nil
	}

	return raw
}

// Get the map value by the key from the config file.
func v2GetConfigFileMapValueByKey(elem cty.Value, key string) interface{} {
	if !elem.Type().HasAttribute(key) {
		return nil
	}

	raw := elem.GetAttr(key)
	if !raw.Type().IsMapType() || raw.IsNull() || !raw.IsKnown() {
		return nil
	}
	return raw
}

// Get the string value by the key from the config file.
func v2GetConfigFileStringValueByKey(elem cty.Value, key string) string {
	if !elem.Type().HasAttribute(key) {
		return ""
	}

	raw := elem.GetAttr(key)
	if raw.Type() != cty.String || raw.IsNull() || !raw.IsKnown() {
		return ""
	}

	return raw.AsString()
}

// Get the int value by the key from the config file.
func v2GetConfigFileIntValueByKey(elem cty.Value, key string) int {
	if !elem.Type().HasAttribute(key) {
		return 0
	}

	raw := elem.GetAttr(key)
	if raw.Type() != cty.Number || raw.IsNull() || !raw.IsKnown() {
		return 0
	}

	rawValue, _ := raw.AsBigFloat().Int64()
	return int(rawValue)
}

// Get the old resource string value by the key from the config file.
func v2GetOldResourceStringValueByConfigFileKey(elem cty.Value, key string, oldResource interface{}) interface{} {
	raw := v2GetConfigFileStringValueByKey(elem, key)
	if raw == "" {
		return utils.ValueIgnoreEmpty(utils.PathSearch(key, oldResource, nil))
	}

	return raw
}

// Get the old resource int value by the key from the config file.
func v2GetOldResourceIntValueByConfigFileKey(elem cty.Value, key string, oldResource interface{}) interface{} {
	num := v2GetConfigFileIntValueByKey(elem, key)
	if num != 0 {
		return num
	}

	return utils.ValueIgnoreEmpty(utils.PathSearch(key, oldResource, nil))
}

func v2GetToStringCreatingStep(resourceElem cty.Value) string {
	creatingStep := v2GetConfigFileListValueByKey(resourceElem, "creating_step")
	if creatingStep == nil {
		return ""
	}

	creatingStepElem := creatingStep.(cty.Value).Index(cty.NumberIntVal(0))
	return utils.JsonToString(map[string]interface{}{
		"step": v2GetConfigFileIntValueByKey(creatingStepElem, "step"),
		"type": v2GetConfigFileStringValueByKey(creatingStepElem, "type"),
	})
}

func v2RefreshResourcesOrderOrigin(rawConfig cty.Value) []map[string]interface{} {
	specResources := v2GetConfigFileSpecResources(rawConfig)
	if specResources == nil {
		return nil
	}

	resources := specResources.(cty.Value)
	result := make([]map[string]interface{}, resources.LengthInt())
	for i, resourceElem := range resources.AsValueSlice() {
		result[i] = map[string]interface{}{
			"flavor":        v2GetConfigFileStringValueByKey(resourceElem, "flavor"),
			"node_pool":     v2GetConfigFileStringValueByKey(resourceElem, "node_pool"),
			"creating_step": v2GetToStringCreatingStep(resourceElem),
		}
	}

	return result
}

func v2GetMatchedResourceFromConfigfile(resourceElem cty.Value, oldResources []interface{}) interface{} {
	var (
		newFlavor        = v2GetConfigFileStringValueByKey(resourceElem, "flavor_id")
		newNodePool      = v2GetConfigFileStringValueByKey(resourceElem, "node_pool")
		newCreatringStep = v2GetToStringCreatingStep(resourceElem)
	)

	oldResource, _ := v2FindResourceByFlavorAndNodePoolAndCreatingStep(oldResources, newFlavor, newNodePool, newCreatringStep)
	return oldResource
}

func v2FindResourceByFlavorAndNodePoolAndCreatingStep(oldResources []interface{}, flavor string, nodePool string,
	creatingStep string) (interface{}, int) {
	for index, oldResource := range oldResources {
		var (
			oldNodePool     = utils.PathSearch("nodePool", oldResource, "").(string)
			oldFlavor       = utils.PathSearch("flavor", oldResource, "").(string)
			oldCreatingStep = utils.JsonToString(utils.PathSearch("creatingStep", oldResource, nil))
		)
		if nodePool != "" &&
			oldNodePool == nodePool &&
			oldFlavor == flavor &&
			oldCreatingStep == creatingStep {
			return oldResource, index
		}

		pattern := regexp.MustCompile(fmt.Sprintf(`^%s-default|$`, oldFlavor))
		if nodePool == "" &&
			pattern.MatchString(oldNodePool) &&
			oldFlavor == flavor &&
			oldCreatingStep == creatingStep {
			return oldResource, index
		}
	}

	return nil, -1
}

func resourceV2ResourcePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		resourcePoolName = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	bssClient, err := cfg.NewServiceClient("bssv2", region)
	if err != nil {
		return diag.Errorf("error creating BSS client: %s", err)
	}

	// When there is no node in the resource pool, the resource pool will be automatically deleted.
	_, err = getV2ResourcePoolByName(client, resourcePoolName)
	if _, ok := err.(golangsdk.ErrDefault404); ok {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting V2 resource pool (%s)", resourcePoolName))
	}

	// If there are nodes in the prepaid billing mode under the resource pool (pre-paid or post-paid), we must unsubscribe the nodes first.
	if err := unsubscribeV2PrePaidBillingNodes(ctx, client, bssClient, resourcePoolName, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error unsubscribing nodes under specified resource pool(%s): %s", resourcePoolName, err)
	}

	resourcePool, err := getV2ResourcePoolByName(client, resourcePoolName)
	if err != nil {
		return diag.Errorf("error querying Modelarts V2 resource pool: %s", err)
	}

	// chargingMode means the billing mode of the resource pool.
	chargingMode := utils.PathSearch(`metadata.annotations."os.modelarts/billing.mode"`, resourcePool, billingModePostPaid).(string)
	if chargingMode == billingModePrePaid {
		mainResourceId := utils.PathSearch(`metadata.labels."os.modelarts/resource.id"`, resourcePool, "").(string)
		if mainResourceId == "" {
			return diag.Errorf("error getting main resource ID from the resource pool(%s)", resourcePoolName)
		}

		if err := common.UnsubscribePrePaidResource(d, cfg, []string{mainResourceId}); err != nil {
			return diag.Errorf("error unsubscribing Modelarts V2 resource pool: %s", err)
		}
	} else {
		err := deleteResourcePoolByName(client, resourcePoolName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = waitForDeleteResourcePoolCompleted(ctx, client, resourcePoolName, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts V2 resource pool (%s) deletion to complete: %s", resourcePoolName, err)
	}

	return nil
}

func unsubscribeV2PrePaidBillingNodes(ctx context.Context, client, bssClient *golangsdk.ServiceClient, resourcePoolName string,
	timeout time.Duration) error {
	nodes, err := listV2ResourcePoolNodes(client, resourcePoolName)
	if err != nil {
		return fmt.Errorf("error querying nodes under specified resource pool (%s): %s", resourcePoolName, err)
	}

	// Obtain the node IDs list that are in the pre-paid billing mode.
	deleteNodeIds := utils.PathSearch(
		fmt.Sprintf(`[?metadata.annotations."os.modelarts/billing.mode"=='%s'].metadata.labels."os.modelarts/resource.id"`,
			billingModePrePaid),
		nodes, make([]interface{}, 0)).([]interface{})

	if len(deleteNodeIds) == 0 {
		return nil
	}

	// Unsubscribe the pre-paid billing nodes.
	err = cbc.UnsubscribePrePaidResources(bssClient, deleteNodeIds)
	if err != nil {
		return err
	}
	err = cbc.WaitForResourcesUnsubscribed(ctx, bssClient, deleteNodeIds, timeout)
	if err != nil {
		return fmt.Errorf("error waiting for all nodes to be unsubscribed: %s ", err)
	}

	err = waitForV2NodeBatchUnsubscribeCompleted(ctx, client, resourcePoolName, deleteNodeIds, timeout)
	if err != nil {
		return err
	}

	return nil
}

func deleteResourcePoolByName(client *golangsdk.ServiceClient, resourcePoolName string) error {
	deleteHttpUrl := "v2/{project_id}/pools/{pool_name}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{pool_name}", resourcePoolName)

	deleteResourcePoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err := client.Request("DELETE", deletePath, &deleteResourcePoolOpt)
	if err != nil {
		return fmt.Errorf("error deleting Modelarts resource pool: %s", err)
	}
	return nil
}

func waitForDeleteResourcePoolCompleted(ctx context.Context, client *golangsdk.ServiceClient, resourcePoolName string, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			res, err := getV2ResourcePoolByName(client, resourcePoolName)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return "resource_pool_not_exist", "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			return res, "PENDING", nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
