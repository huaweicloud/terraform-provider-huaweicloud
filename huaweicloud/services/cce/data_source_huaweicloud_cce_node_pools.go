package cce

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE GET /api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools
func DataSourceNodePools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNodePoolsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the CCE node pools are located.`,
			},

			// Required parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the CCE cluster.`,
			},

			// Optional parameters.
			"show_default_node_pool": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
				Description: `Whether to show the default node pool.`,
			},

			// Attributes.
			"node_pools": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataNodePoolsElem(),
				Description: `The list of CCE node pools that matched filter parameters.`,
			},
		},
	}
}

func dataNodePoolsElem() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the node pool.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the node pool.`,
			},
			"initial_node_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The initial number of nodes in the node pool.`,
			},
			"current_node_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The current number of nodes in the node pool.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor ID of the node pool.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the node pool.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The availability zone of the node pool.`,
			},
			"os": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The operating system of the node pool.`,
			},
			"key_pair": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The key pair name of the node pool.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subnet ID of the NIC.`,
			},
			"subnet_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The subnet ID list of the NIC.`,
			},
			"ecs_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ECS group ID of the node pool.`,
			},
			"max_pods": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum number of pods allowed on a node.`,
			},
			"extend_param": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The extended parameters of the node pool, in key/value format.`,
			},
			"extend_params":          dataNodePoolsExtendParamsElem(),
			"extension_scale_groups": nodePoolExtensionScaleGroupsSchema(),
			"scall_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether auto scaling is enabled.`,
			},
			"min_node_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The minimum number of nodes if auto scaling is enabled.`,
			},
			"max_node_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum number of nodes if auto scaling is enabled.`,
			},
			"scale_down_cooldown_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The interval between two scaling operations, in minutes.`,
			},
			"priority": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The weight of the node pool during scaling.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels of a Kubernetes node.`,
			},
			"tags": common.TagsComputedSchema(),
			"root_volume": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataNodePoolsVolumeElem(),
				Description: `The system disk configuration of the node pool.`,
			},
			"data_volumes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataNodePoolsVolumeElem(),
				Description: `The data disk configurations of the node pool.`,
			},
			"storage": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataNodePoolsStorageElem(),
				Description: `The disk initialization configuration of the node pool.`,
			},
			"taints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The key of the taint.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The value of the taint.`,
						},
						"effect": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The effect of the taint.`,
						},
					},
				},
				Description: `The taints configuration of the node pool.`,
			},
			"security_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The custom security group IDs of the node pool.`,
			},
			"pod_security_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The pod security group IDs of the node pool.`,
			},
			"initialized_conditions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The custom initialization flags of the node pool.`,
			},
			"label_policy_on_existing_nodes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The label policy on existing nodes.`,
			},
			"tag_policy_on_existing_nodes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The tag policy on existing nodes.`,
			},
			"taint_policy_on_existing_nodes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The taint policy on existing nodes.`,
			},
			"hostname_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The hostname type of the kubernetes node.`,
						},
					},
				},
				Description: `The hostname configuration of the kubernetes node.`,
			},
			"partition": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The partition to which the node belongs.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise project ID of the node pool.`,
			},
			"runtime": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The runtime of the node pool.`,
			},
			"billing_mode": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The billing mode of a node.`,
			},
			"period_unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The charging period unit of the node pool.`,
			},
			"period": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The charging period of the node pool.`,
			},
			"auto_renew": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Whether auto-renew is enabled.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the node pool.`,
			},
		},
	}

	return &sc
}

func dataNodePoolsExtendParamsElem() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: `The extended parameters of the node pool.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"max_pods": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: `The maximum number of pods allowed on a node.`,
				},
				"docker_base_size": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: `The available disk space of a single container on a node, in GB.`,
				},
				"preinstall": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `The script to be executed before installation.`,
				},
				"postinstall": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `The script to be executed after installation.`,
				},
				"node_image_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `The image ID used to create the node.`,
				},
				"node_multi_queue": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `The number of ENI queues.`,
				},
				"nic_threshold": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `The ENI pre-binding thresholds.`,
				},
				"agency_name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `The agency name of the node pool.`,
				},
				"kube_reserved_mem": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: `The reserved memory for Kubernetes components.`,
				},
				"system_reserved_mem": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: `The reserved memory for system components.`,
				},
				"security_reinforcement_type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `The security reinforcement type.`,
				},
				"market_type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `The market type of the spot node pool.`,
				},
				"spot_price": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `The highest price per hour for a spot node.`,
				},
			},
		},
	}
}

func dataNodePoolsStorageElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"selectors": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The disk selection configuration.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The selector name.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The storage type.`,
						},
						"match_label_size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The matched disk size.`,
						},
						"match_label_volume_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The EVS disk type.`,
						},
						"match_label_metadata_encrypted": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The disk encryption identifier.`,
						},
						"match_label_metadata_cmkid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The customer master key ID of an encrypted disk.`,
						},
						"match_label_count": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The number of disks to be selected.`,
						},
					},
				},
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The storage group configuration.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of a virtual storage group.`,
						},
						"cce_managed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the storage space is for kubernetes and runtime components.`,
						},
						"selector_names": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of selector names to match.`,
						},
						"virtual_spaces": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The detailed space configuration in a group.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The virtual space name.`,
									},
									"size": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The size of a virtual space.`,
									},
									"lvm_lv_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The LVM write mode.`,
									},
									"lvm_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The absolute path to which the disk is attached.`,
									},
									"runtime_lv_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The LVM write mode of runtime.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataNodePoolsVolumeElem() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The disk size in GB.`,
			},
			"volumetype": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The disk type.`,
			},
			"extend_params": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The disk expansion parameters.`,
			},
			"kms_key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The KMS key ID of the disk.`,
			},
			"dss_pool_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The DSS pool ID of the disk.`,
			},
			"iops": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The IOPS of the disk.`,
			},
			"throughput": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The throughput of the disk.`,
			},
			"hw_passthrough": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether passthrough is enabled.`,
			},
		},
	}

	return &sc
}

func buildGetNodePoolsQueryParams(d *schema.ResourceData) string {
	if showDefault, ok := d.GetOk("show_default_node_pool"); ok {
		return fmt.Sprintf("?showDefaultNodePool=%v", showDefault)
	}
	return ""
}

func listNodePools(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl  = "api/v3/projects/{project_id}/clusters/{cluster_id}/nodepools"
		listPath = client.Endpoint + httpUrl
	)

	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{cluster_id}", d.Get("cluster_id").(string))
	listPath += buildGetNodePoolsQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenNodePoolStorageVirtualSpaces(group interface{}) []map[string]interface{} {
	virtualSpacesRaw := utils.PathSearch("virtualSpaces", group, make([]interface{}, 0)).([]interface{})
	if len(virtualSpacesRaw) < 1 {
		return nil
	}

	virtualSpaces := make([]map[string]interface{}, 0, len(virtualSpacesRaw))
	for _, space := range virtualSpacesRaw {
		virtualSpaces = append(virtualSpaces, map[string]interface{}{
			"name":            utils.PathSearch("name", space, nil),
			"size":            utils.PathSearch("size", space, nil),
			"lvm_lv_type":     utils.PathSearch("lvmConfig.lvType", space, nil),
			"lvm_path":        utils.PathSearch("lvmConfig.path", space, nil),
			"runtime_lv_type": utils.PathSearch("runtimeConfig.lvType", space, nil),
		})
	}
	return virtualSpaces
}

func flattenNodePoolStorage(storage interface{}) []map[string]interface{} {
	if storage == nil {
		return nil
	}

	selectorsRaw := utils.PathSearch("storageSelectors", storage, make([]interface{}, 0)).([]interface{})
	groupsRaw := utils.PathSearch("storageGroups", storage, make([]interface{}, 0)).([]interface{})
	if len(selectorsRaw) < 1 && len(groupsRaw) < 1 {
		return nil
	}

	selectors := make([]map[string]interface{}, 0, len(selectorsRaw))
	for _, selector := range selectorsRaw {
		selectors = append(selectors, map[string]interface{}{
			"name":                           utils.PathSearch("name", selector, nil),
			"type":                           utils.PathSearch("storageType", selector, nil),
			"match_label_size":               utils.PathSearch("matchLabels.size", selector, nil),
			"match_label_volume_type":        utils.PathSearch("matchLabels.volumeType", selector, nil),
			"match_label_metadata_encrypted": utils.PathSearch("matchLabels.metadataEncrypted", selector, nil),
			"match_label_metadata_cmkid":     utils.PathSearch("matchLabels.metadataCmkid", selector, nil),
			"match_label_count":              utils.PathSearch("matchLabels.count", selector, nil),
		})
	}

	groups := make([]map[string]interface{}, 0, len(groupsRaw))
	for _, group := range groupsRaw {
		groups = append(groups, map[string]interface{}{
			"name":           utils.PathSearch("name", group, nil),
			"cce_managed":    utils.PathSearch("cceManaged", group, nil),
			"selector_names": utils.PathSearch("selectorNames", group, make([]interface{}, 0)),
			"virtual_spaces": flattenNodePoolStorageVirtualSpaces(group),
		})
	}

	return []map[string]interface{}{
		{
			"selectors": selectors,
			"groups":    groups,
		},
	}
}

func flattenNodePoolTags(userTags []interface{}) map[string]string {
	if len(userTags) < 1 {
		return nil
	}

	result := make(map[string]string, len(userTags))
	for _, tag := range userTags {
		key := utils.PathSearch("key", tag, "").(string)
		if key == "" {
			continue
		}
		result[key] = utils.PathSearch("value", tag, "").(string)
	}
	return result
}

func flattenNodePoolRootVolume(rootVolume interface{}) []map[string]interface{} {
	if rootVolume == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"size":           int(utils.PathSearch("size", rootVolume, float64(0)).(float64)),
			"volumetype":     utils.PathSearch("volumetype", rootVolume, nil),
			"extend_params":  utils.PathSearch("extendParam", rootVolume, nil),
			"hw_passthrough": utils.PathSearch(`"hw:passthrough"`, rootVolume, nil),
			"dss_pool_id":    utils.PathSearch("cluster_id", rootVolume, nil),
			"iops":           int(utils.PathSearch("iops", rootVolume, float64(0)).(float64)),
			"throughput":     int(utils.PathSearch("throughput", rootVolume, float64(0)).(float64)),
			"kms_key_id":     utils.PathSearch(`metadata."__system__cmkid"`, rootVolume, nil),
		},
	}
}

func flattenNodePoolDataVolumes(dataVolumes []interface{}) []map[string]interface{} {
	if len(dataVolumes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dataVolumes))
	for _, volume := range dataVolumes {
		result = append(result, map[string]interface{}{
			"size":           int(utils.PathSearch("size", volume, float64(0)).(float64)),
			"volumetype":     utils.PathSearch("volumetype", volume, nil),
			"extend_params":  utils.PathSearch("extendParam", volume, nil),
			"hw_passthrough": utils.PathSearch(`"hw:passthrough"`, volume, nil),
			"dss_pool_id":    utils.PathSearch("cluster_id", volume, nil),
			"iops":           int(utils.PathSearch("iops", volume, float64(0)).(float64)),
			"throughput":     int(utils.PathSearch("throughput", volume, float64(0)).(float64)),
			"kms_key_id":     utils.PathSearch(`metadata."__system__cmkid"`, volume, nil),
		})
	}
	return result
}

func flattenNodePoolHostnameConfig(hostnameConfig interface{}) []map[string]interface{} {
	if hostnameConfig == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"type": utils.PathSearch("type", hostnameConfig, nil),
		},
	}
}

func flattenNodePoolExtendParams(extendParams map[string]interface{}) []map[string]interface{} {
	if len(extendParams) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"max_pods":                    int(utils.PathSearch("maxPods", extendParams, float64(0)).(float64)),
			"docker_base_size":            int(utils.PathSearch("dockerBaseSize", extendParams, float64(0)).(float64)),
			"preinstall":                  utils.PathSearch(`"alpha.cce/preInstall"`, extendParams, nil),
			"postinstall":                 utils.PathSearch(`"alpha.cce/postInstall"`, extendParams, nil),
			"node_image_id":               utils.PathSearch(`"alpha.cce/NodeImageID"`, extendParams, nil),
			"node_multi_queue":            utils.PathSearch("nicMultiqueue", extendParams, nil),
			"nic_threshold":               utils.PathSearch("nicThreshold", extendParams, nil),
			"agency_name":                 utils.PathSearch("agency_name", extendParams, nil),
			"kube_reserved_mem":           int(utils.PathSearch("kubeReservedMem", extendParams, float64(0)).(float64)),
			"system_reserved_mem":         int(utils.PathSearch("systemReservedMem", extendParams, float64(0)).(float64)),
			"security_reinforcement_type": utils.PathSearch("securityReinforcementType", extendParams, nil),
			"market_type":                 utils.PathSearch("marketType", extendParams, nil),
			"spot_price":                  utils.PathSearch("spotPrice", extendParams, nil),
		},
	}
}

func flattenNodePoolTaints(taints []interface{}) []map[string]interface{} {
	if len(taints) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(taints))
	for _, taint := range taints {
		result = append(result, map[string]interface{}{
			"key":    utils.PathSearch("key", taint, nil),
			"value":  utils.PathSearch("value", taint, nil),
			"effect": utils.PathSearch("effect", taint, nil),
		})
	}
	return result
}

func flattenNodePoolExtensionScaleGroupsSpecAutoscaling(spec interface{}) []map[string]interface{} {
	autoscaling := utils.PathSearch("autoscaling", spec, nil)
	if autoscaling == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"extension_priority": int(utils.PathSearch("extensionPriority", autoscaling, float64(0)).(float64)),
			"max_node_count":     int(utils.PathSearch("maxNodeCount", autoscaling, float64(0)).(float64)),
			"min_node_count":     int(utils.PathSearch("minNodeCount", autoscaling, float64(0)).(float64)),
			"enable":             utils.PathSearch("enable", autoscaling, nil),
		},
	}
}

func flattenNodePoolExtensionScaleGroupsSpec(extensionScaleGroup interface{}) []map[string]interface{} {
	spec := utils.PathSearch("spec", extensionScaleGroup, nil)
	if spec == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"flavor":                             utils.PathSearch("flavor", spec, nil),
			"az":                                 utils.PathSearch("az", spec, nil),
			"capacity_reservation_specification": flattenExtensionScaleGroupsSpecCapacity(spec),
			"autoscaling":                        flattenNodePoolExtensionScaleGroupsSpecAutoscaling(spec),
		},
	}
}

func flattenNodePoolExtensionScaleGroups(groups []interface{}) []map[string]interface{} {
	if len(groups) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(groups))
	for _, group := range groups {
		result = append(result, map[string]interface{}{
			"metadata": flattenExtensionScaleGroupsMetadata(group),
			"spec":     flattenNodePoolExtensionScaleGroupsSpec(group),
		})
	}
	return result
}

func flattenNodePoolLabels(k8sTags interface{}) map[string]interface{} {
	tagsMap, ok := k8sTags.(map[string]interface{})
	if !ok || len(tagsMap) < 1 {
		return nil
	}

	result := make(map[string]interface{})
	for key, val := range tagsMap {
		if strings.Contains(key, "cce.cloud.com") {
			continue
		}
		result[key] = val
	}

	if len(result) < 1 {
		return nil
	}
	return result
}

func flattenNodePoolExtendParam(extendParam map[string]interface{}) map[string]string {
	if len(extendParam) < 1 {
		return nil
	}

	result := make(map[string]string, len(extendParam))
	for key, val := range extendParam {
		if strVal, ok := val.(string); ok {
			result[key] = strVal
			continue
		}
		result[key] = utils.JsonToString(val)
	}
	return result
}

func flattenNodePoolItem(nodePool interface{}) map[string]interface{} {
	extendParam := utils.PathSearch("spec.nodeTemplate.extendParam", nodePool, make(map[string]interface{})).(map[string]interface{})

	item := map[string]interface{}{
		"id":                       utils.PathSearch("metadata.uid", nodePool, nil),
		"name":                     utils.PathSearch("metadata.name", nodePool, nil),
		"type":                     utils.PathSearch("spec.type", nodePool, nil),
		"flavor_id":                utils.PathSearch("spec.nodeTemplate.flavor", nodePool, nil),
		"availability_zone":        utils.PathSearch("spec.nodeTemplate.az", nodePool, nil),
		"os":                       utils.PathSearch("spec.nodeTemplate.os", nodePool, nil),
		"billing_mode":             int(utils.PathSearch("spec.nodeTemplate.billingMode", nodePool, float64(0)).(float64)),
		"key_pair":                 utils.PathSearch("spec.nodeTemplate.login.sshKey", nodePool, nil),
		"scall_enable":             utils.PathSearch("spec.autoscaling.enable", nodePool, nil),
		"initial_node_count":       int(utils.PathSearch("spec.initialNodeCount", nodePool, float64(0)).(float64)),
		"current_node_count":       int(utils.PathSearch("status.currentNode", nodePool, float64(0)).(float64)),
		"min_node_count":           int(utils.PathSearch("spec.autoscaling.minNodeCount", nodePool, float64(0)).(float64)),
		"max_node_count":           int(utils.PathSearch("spec.autoscaling.maxNodeCount", nodePool, float64(0)).(float64)),
		"scale_down_cooldown_time": int(utils.PathSearch("spec.autoscaling.scaleDownCooldownTime", nodePool, float64(0)).(float64)),
		"priority":                 int(utils.PathSearch("spec.autoscaling.priority", nodePool, float64(0)).(float64)),
		"ecs_group_id":             utils.PathSearch("spec.nodeManagement.serverGroupReference", nodePool, nil),
		"storage":                  flattenNodePoolStorage(utils.PathSearch("spec.nodeTemplate.storage", nodePool, nil)),
		"security_groups":          utils.PathSearch("spec.customSecurityGroups", nodePool, make([]interface{}, 0)),
		"tags": flattenNodePoolTags(utils.PathSearch("spec.nodeTemplate.userTags",
			nodePool, make([]interface{}, 0)).([]interface{})),
		"status":      utils.PathSearch("status.phase", nodePool, nil),
		"root_volume": flattenNodePoolRootVolume(utils.PathSearch("spec.nodeTemplate.rootVolume", nodePool, nil)),
		"data_volumes": flattenNodePoolDataVolumes(utils.PathSearch("spec.nodeTemplate.dataVolumes",
			nodePool, make([]interface{}, 0)).([]interface{})),
		"initialized_conditions":         utils.PathSearch("spec.nodeTemplate.initializedConditions", nodePool, make([]interface{}, 0)),
		"label_policy_on_existing_nodes": utils.PathSearch("spec.labelPolicyOnExistingNodes", nodePool, nil),
		"tag_policy_on_existing_nodes":   utils.PathSearch("spec.userTagsPolicyOnExistingNodes", nodePool, nil),
		"taint_policy_on_existing_nodes": utils.PathSearch("spec.taintPolicyOnExistingNodes", nodePool, nil),
		"hostname_config":                flattenNodePoolHostnameConfig(utils.PathSearch("spec.nodeTemplate.hostnameConfig", nodePool, nil)),
		"enterprise_project_id":          utils.PathSearch("spec.nodeTemplate.serverEnterpriseProjectID", nodePool, nil),
		"subnet_id":                      utils.PathSearch("spec.nodeTemplate.nodeNicSpec.primaryNic.subnetId", nodePool, nil),
		"subnet_list":                    utils.PathSearch("spec.nodeTemplate.nodeNicSpec.primaryNic.subnetList", nodePool, make([]interface{}, 0)),
		"extend_params":                  flattenNodePoolExtendParams(extendParam),
		"taints": flattenNodePoolTaints(utils.PathSearch("spec.nodeTemplate.taints",
			nodePool, make([]interface{}, 0)).([]interface{})),
		"extension_scale_groups": flattenNodePoolExtensionScaleGroups(utils.PathSearch("spec.extensionScaleGroups",
			nodePool, make([]interface{}, 0)).([]interface{})),
		"period_unit":         utils.PathSearch("periodType", extendParam, nil),
		"period":              int(utils.PathSearch("periodNum", extendParam, float64(0)).(float64)),
		"auto_renew":          utils.PathSearch("isAutoRenew", extendParam, nil),
		"runtime":             utils.PathSearch("spec.nodeTemplate.runtime.name", nodePool, nil),
		"partition":           utils.PathSearch("spec.nodeTemplate.partition", nodePool, nil),
		"pod_security_groups": utils.PathSearch("spec.podSecurityGroups[*].id", nodePool, make([]interface{}, 0)),
		"labels":              flattenNodePoolLabels(utils.PathSearch("spec.nodeTemplate.k8sTags", nodePool, nil)),
		"extend_param":        flattenNodePoolExtendParam(extendParam),
		"max_pods":            int(utils.PathSearch("maxPods", extendParam, float64(0)).(float64)),
	}

	return item
}

func flattenNodePools(nodePools []interface{}) []map[string]interface{} {
	if len(nodePools) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(nodePools))
	for _, nodePool := range nodePools {
		result = append(result, flattenNodePoolItem(nodePool))
	}
	return result
}

func dataSourceNodePoolsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("cce", region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	nodePools, err := listNodePools(client, d)
	if err != nil {
		return diag.Errorf("error querying CCE node pools: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("node_pools", flattenNodePools(nodePools)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
