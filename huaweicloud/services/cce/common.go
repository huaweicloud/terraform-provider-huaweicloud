package cce

import (
	"log"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func resourceNodeExtendParamsSchema(conflictList []string) *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: conflictList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"max_pods": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"docker_base_size": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"preinstall": {
					Type:      schema.TypeString,
					Optional:  true,
					StateFunc: utils.DecodeHashAndHexEncode,
				},
				"postinstall": {
					Type:      schema.TypeString,
					Optional:  true,
					StateFunc: utils.DecodeHashAndHexEncode,
				},
				"node_image_id": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"node_multi_queue": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"nic_threshold": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"agency_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"kube_reserved_mem": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"system_reserved_mem": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"security_reinforcement_type": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"market_type": {
					Type:     schema.TypeString,
					Optional: true,
					Description: utils.SchemaDesc(
						"",
						utils.SchemaDescInput{
							Internal: true,
						},
					),
				},
				"spot_price": {
					Type:     schema.TypeString,
					Optional: true,
					Description: utils.SchemaDesc(
						"",
						utils.SchemaDescInput{
							Internal: true,
						},
					),
				},
			},
		},
	}
}

func resourceNodePoolExtendParamsSchema(conflictList []string) *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		Computed:      true,
		ConflictsWith: conflictList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"max_pods": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"docker_base_size": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"preinstall": {
					Type:      schema.TypeString,
					Optional:  true,
					StateFunc: utils.DecodeHashAndHexEncode,
					Computed:  true,
				},
				"postinstall": {
					Type:      schema.TypeString,
					Optional:  true,
					StateFunc: utils.DecodeHashAndHexEncode,
					Computed:  true,
				},
				"node_image_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"node_multi_queue": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"nic_threshold": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"agency_name": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"kube_reserved_mem": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"system_reserved_mem": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"security_reinforcement_type": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"market_type": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					Description: utils.SchemaDesc(
						"",
						utils.SchemaDescInput{
							Internal: true,
						},
					),
				},
				"spot_price": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					Description: utils.SchemaDesc(
						"",
						utils.SchemaDescInput{
							Internal: true,
						},
					),
				},
			},
		},
	}
}

func buildResourceNodeExtendParam(d *schema.ResourceData) map[string]interface{} {
	extendParam := make(map[string]interface{})
	if v, ok := d.GetOk("extend_param"); ok {
		for key, val := range v.(map[string]interface{}) {
			extendParam[key] = val.(string)
		}
		if v, ok := extendParam["periodNum"]; ok {
			periodNum, err := strconv.Atoi(v.(string))
			if err != nil {
				log.Printf("[WARNING] PeriodNum %s invalid, Type conversion error: %s", v.(string), err)
			}
			extendParam["periodNum"] = periodNum
		}
	}

	if v, ok := d.GetOk("ecs_performance_type"); ok {
		extendParam["ecs:performancetype"] = v.(string)
	}
	if v, ok := d.GetOk("max_pods"); ok {
		extendParam["maxPods"] = v.(int)
	}
	if v, ok := d.GetOk("order_id"); ok {
		extendParam["orderID"] = v.(string)
	}
	if v, ok := d.GetOk("product_id"); ok {
		extendParam["productID"] = v.(string)
	}
	if v, ok := d.GetOk("public_key"); ok {
		extendParam["publicKey"] = v.(string)
	}
	if v, ok := d.GetOk("preinstall"); ok {
		extendParam["alpha.cce/preInstall"] = utils.TryBase64EncodeString(v.(string))
	}
	if v, ok := d.GetOk("postinstall"); ok {
		extendParam["alpha.cce/postInstall"] = utils.TryBase64EncodeString(v.(string))
	}

	return extendParam
}

func buildResourceNodeExtendParams(extendParamsRaw []interface{}) map[string]interface{} {
	if len(extendParamsRaw) != 1 {
		return nil
	}

	if extendParams, ok := extendParamsRaw[0].(map[string]interface{}); ok {
		res := map[string]interface{}{
			"maxPods":                   utils.ValueIgnoreEmpty(extendParams["max_pods"]),
			"dockerBaseSize":            utils.ValueIgnoreEmpty(extendParams["docker_base_size"]),
			"alpha.cce/preInstall":      utils.ValueIgnoreEmpty(utils.TryBase64EncodeString(extendParams["preinstall"].(string))),
			"alpha.cce/postInstall":     utils.ValueIgnoreEmpty(utils.TryBase64EncodeString(extendParams["postinstall"].(string))),
			"alpha.cce/NodeImageID":     utils.ValueIgnoreEmpty(extendParams["node_image_id"]),
			"nicMultiqueue":             utils.ValueIgnoreEmpty(extendParams["node_multi_queue"]),
			"nicThreshold":              utils.ValueIgnoreEmpty(extendParams["nic_threshold"]),
			"agency_name":               utils.ValueIgnoreEmpty(extendParams["agency_name"]),
			"kubeReservedMem":           utils.ValueIgnoreEmpty(extendParams["kube_reserved_mem"]),
			"systemReservedMem":         utils.ValueIgnoreEmpty(extendParams["system_reserved_mem"]),
			"marketType":                utils.ValueIgnoreEmpty(extendParams["market_type"]),
			"spotPrice":                 utils.ValueIgnoreEmpty(extendParams["spot_price"]),
			"securityReinforcementType": utils.ValueIgnoreEmpty(extendParams["security_reinforcement_type"]),
		}

		return res
	}

	return nil
}

func buildExtendParams(d *schema.ResourceData) map[string]interface{} {
	res := make(map[string]interface{})
	extendParam := buildResourceNodeExtendParam(d)
	extendParams := buildResourceNodeExtendParams(d.Get("extend_params").([]interface{}))

	// defaults to use extend_params
	if len(extendParam) != 0 {
		for k, v := range extendParam {
			res[k] = v
		}
	} else {
		for k, v := range extendParams {
			res[k] = v
		}
	}

	// assemble the charge info
	var isPrePaid bool
	var billingMode int

	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		isPrePaid = true
	}
	if v, ok := d.GetOk("billing_mode"); ok {
		billingMode = v.(int)
	}
	if isPrePaid || billingMode == 1 {
		res["chargingMode"] = 1
		res["isAutoRenew"] = "false"
		res["isAutoPay"] = common.GetAutoPay(d)
	}

	if v, ok := d.GetOk("period_unit"); ok {
		res["periodType"] = v.(string)
	}
	if v, ok := d.GetOk("period"); ok {
		res["periodNum"] = v.(int)
	}
	if v, ok := d.GetOk("auto_renew"); ok {
		res["isAutoRenew"] = v.(string)
	}

	return utils.RemoveNil(res)
}

func flattenExtendParams(extendParams map[string]interface{}) []map[string]interface{} {
	if len(extendParams) == 0 {
		return nil
	}

	res := []map[string]interface{}{
		{
			"max_pods":                    utils.PathSearch("maxPods", extendParams, nil),
			"docker_base_size":            utils.PathSearch("dockerBaseSize", extendParams, nil),
			"preinstall":                  utils.PathSearch("alpha.cce/preInstall", extendParams, nil),
			"postinstall":                 utils.PathSearch("alpha.cce/postInstall", extendParams, nil),
			"node_image_id":               utils.PathSearch("alpha.cce/NodeImageID", extendParams, nil),
			"node_multi_queue":            utils.PathSearch("nicMultiqueue", extendParams, nil),
			"nic_threshold":               utils.PathSearch("nicThreshold", extendParams, nil),
			"agency_name":                 utils.PathSearch("agency_name", extendParams, nil),
			"kube_reserved_mem":           utils.PathSearch("kubeReservedMem", extendParams, nil),
			"system_reserved_mem":         utils.PathSearch("systemReservedMem", extendParams, nil),
			"security_reinforcement_type": utils.PathSearch("securityReinforcementType", extendParams, nil),
			"market_type":                 utils.PathSearch("marketType", extendParams, nil),
			"spot_price":                  utils.PathSearch("spotPrice", extendParams, nil),
		},
	}

	return res
}

func resourceNodeRootVolume() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"size": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"volumetype": {
					Type:     schema.TypeString,
					Required: true,
				},
				"extend_params": {
					Type:     schema.TypeMap,
					Optional: true,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"kms_key_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"dss_pool_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"iops": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"throughput": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},

				// Internal parameters
				"hw_passthrough": {
					Type:        schema.TypeBool,
					Optional:    true,
					Computed:    true,
					Description: "schema: Internal",
				},

				// Deprecated parameters
				"extend_param": {
					Type:       schema.TypeString,
					Optional:   true,
					Deprecated: "use extend_params instead",
				},
			},
		},
	}
}

func resourceNodeDataVolume() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Description: utils.SchemaDesc("", utils.SchemaDescInput{
			Required: true,
		}),
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"size": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"volumetype": {
					Type:     schema.TypeString,
					Required: true,
				},
				"extend_params": {
					Type:     schema.TypeMap,
					Optional: true,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"kms_key_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"dss_pool_id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"iops": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"throughput": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},

				// Internal parameters
				"hw_passthrough": {
					Type:        schema.TypeBool,
					Optional:    true,
					Computed:    true,
					Description: "schema: Internal",
				},

				// Deprecated parameters
				"extend_param": {
					Type:       schema.TypeString,
					Optional:   true,
					Deprecated: "use extend_params instead",
				},
			},
		},
	}
}

func buildResourceNodeRootVolume(d *schema.ResourceData) nodes.VolumeSpec {
	var root nodes.VolumeSpec
	volumeRaw := d.Get("root_volume").([]interface{})
	if len(volumeRaw) == 1 {
		rawMap := volumeRaw[0].(map[string]interface{})
		root.Size = rawMap["size"].(int)
		root.VolumeType = rawMap["volumetype"].(string)
		root.HwPassthrough = rawMap["hw_passthrough"].(bool)
		root.ExtendParam = rawMap["extend_params"].(map[string]interface{})
		root.Iops = rawMap["iops"].(int)
		root.Throughput = rawMap["throughput"].(int)

		if rawMap["kms_key_id"].(string) != "" {
			metadata := nodes.VolumeMetadata{
				SystemEncrypted: "1",
				SystemCmkid:     rawMap["kms_key_id"].(string),
			}
			root.Metadata = &metadata
		}

		if rawMap["dss_pool_id"].(string) != "" {
			root.ClusterID = rawMap["dss_pool_id"].(string)
			root.ClusterType = "dss"
		}
	}

	return root
}

func buildResourceNodeDataVolume(d *schema.ResourceData) []nodes.VolumeSpec {
	volumeRaw := d.Get("data_volumes").([]interface{})
	volumes := make([]nodes.VolumeSpec, len(volumeRaw))
	for i, raw := range volumeRaw {
		rawMap := raw.(map[string]interface{})
		volumes[i] = nodes.VolumeSpec{
			Size:          rawMap["size"].(int),
			VolumeType:    rawMap["volumetype"].(string),
			HwPassthrough: rawMap["hw_passthrough"].(bool),
			ExtendParam:   rawMap["extend_params"].(map[string]interface{}),
			Iops:          rawMap["iops"].(int),
			Throughput:    rawMap["throughput"].(int),
		}
		if rawMap["kms_key_id"].(string) != "" {
			metadata := nodes.VolumeMetadata{
				SystemEncrypted: "1",
				SystemCmkid:     rawMap["kms_key_id"].(string),
			}
			volumes[i].Metadata = &metadata
		}

		if rawMap["dss_pool_id"].(string) != "" {
			volumes[i].ClusterID = rawMap["dss_pool_id"].(string)
			volumes[i].ClusterType = "dss"
		}
	}
	return volumes
}

func flattenResourceNodeRootVolume(d *schema.ResourceData, rootVolume nodes.VolumeSpec) []map[string]interface{} {
	res := []map[string]interface{}{
		{
			"size":           rootVolume.Size,
			"volumetype":     rootVolume.VolumeType,
			"hw_passthrough": rootVolume.HwPassthrough,
			"extend_param":   "",
			"dss_pool_id":    rootVolume.ClusterID,
			"iops":           rootVolume.Iops,
			"throughput":     rootVolume.Throughput,
		},
	}

	orignRootVolume := buildResourceNodeRootVolume(d)
	if !reflect.DeepEqual(orignRootVolume, nodes.VolumeSpec{}) {
		orignExtendParams := orignRootVolume.ExtendParam
		extendParams := make(map[string]interface{})

		for k := range orignExtendParams {
			if value, ok := rootVolume.ExtendParam[k]; ok {
				extendParams[k] = value
			}
		}

		res[0]["extend_params"] = extendParams
	} else {
		res[0]["extend_params"] = rootVolume.ExtendParam
	}

	if rootVolume.Metadata != nil {
		res[0]["kms_key_id"] = rootVolume.Metadata.SystemCmkid
	}

	return res
}

func flattenResourceNodeDataVolume(d *schema.ResourceData, dataVolumes []nodes.VolumeSpec) []map[string]interface{} {
	if len(dataVolumes) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(dataVolumes))
	orignDataVolumes := buildResourceNodeDataVolume(d)
	if len(orignDataVolumes) == len(dataVolumes) {
		for i, v := range dataVolumes {
			res[i] = map[string]interface{}{
				"size":           v.Size,
				"volumetype":     v.VolumeType,
				"hw_passthrough": v.HwPassthrough,
				"extend_param":   "",
				"dss_pool_id":    v.ClusterID,
				"iops":           v.Iops,
				"throughput":     v.Throughput,
			}

			orignExtendParams := orignDataVolumes[i].ExtendParam
			extendParams := make(map[string]interface{})

			for k := range orignExtendParams {
				if value, ok := v.ExtendParam[k]; ok {
					extendParams[k] = value
				}
			}

			res[i]["extend_params"] = extendParams

			if v.Metadata != nil {
				res[i]["kms_key_id"] = v.Metadata.SystemCmkid
			}
		}
	} else {
		for i, v := range dataVolumes {
			res[i] = map[string]interface{}{
				"size":           v.Size,
				"volumetype":     v.VolumeType,
				"hw_passthrough": v.HwPassthrough,
				"extend_param":   "",
				"dss_pool_id":    v.ClusterID,
				"extend_params":  v.ExtendParam,
				"iops":           v.Iops,
				"throughput":     v.Throughput,
			}

			if v.Metadata != nil {
				res[i]["kms_key_id"] = v.Metadata.SystemCmkid
			}
		}
	}

	return res
}

func resourceNodeStorageSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"selectors": {
					Type:     schema.TypeList,
					Required: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Required: true,
							},
							"type": {
								Type:     schema.TypeString,
								Optional: true,
								Default:  "evs",
							},
							"match_label_size": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"match_label_volume_type": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"match_label_metadata_encrypted": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"match_label_metadata_cmkid": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"match_label_count": {
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
				},
				"groups": {
					Type:     schema.TypeList,
					Required: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Required: true,
							},
							"cce_managed": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"selector_names": {
								Type:     schema.TypeList,
								Required: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
							"virtual_spaces": {
								Type:     schema.TypeList,
								Required: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"name": {
											Type:     schema.TypeString,
											Required: true,
										},
										"size": {
											Type:     schema.TypeString,
											Required: true,
										},
										"lvm_lv_type": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"lvm_path": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"runtime_lv_type": {
											Type:     schema.TypeString,
											Optional: true,
										},
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

func buildResourceNodeStorage(d *schema.ResourceData) *nodes.StorageSpec {
	v, ok := d.GetOk("storage")
	if !ok {
		return nil
	}

	var storageSpec nodes.StorageSpec
	storageSpecRaw := v.([]interface{})
	storageSpecRawMap := storageSpecRaw[0].(map[string]interface{})
	storageSelectorSpecRaw := storageSpecRawMap["selectors"].([]interface{})
	storageGroupSpecRaw := storageSpecRawMap["groups"].([]interface{})

	var selectors []nodes.StorageSelectorsSpec
	for _, s := range storageSelectorSpecRaw {
		sMap := s.(map[string]interface{})
		selector := nodes.StorageSelectorsSpec{
			Name:        sMap["name"].(string),
			StorageType: sMap["type"].(string),
			MatchLabels: nodes.MatchLabelsSpec{
				Size:              sMap["match_label_size"].(string),
				VolumeType:        sMap["match_label_volume_type"].(string),
				MetadataEncrypted: sMap["match_label_metadata_encrypted"].(string),
				MetadataCmkid:     sMap["match_label_metadata_cmkid"].(string),
				Count:             sMap["match_label_count"].(string),
			},
		}
		selectors = append(selectors, selector)
	}
	storageSpec.StorageSelectors = selectors

	var groups []nodes.StorageGroupsSpec
	for _, g := range storageGroupSpecRaw {
		gMap := g.(map[string]interface{})
		group := nodes.StorageGroupsSpec{
			Name:          gMap["name"].(string),
			CceManaged:    gMap["cce_managed"].(bool),
			SelectorNames: utils.ExpandToStringList(gMap["selector_names"].([]interface{})),
		}

		virtualSpacesRaw := gMap["virtual_spaces"].([]interface{})
		virtualSpaces := make([]nodes.VirtualSpacesSpec, 0, len(virtualSpacesRaw))
		for _, v := range virtualSpacesRaw {
			virtualSpaceMap := v.(map[string]interface{})
			virtualSpace := nodes.VirtualSpacesSpec{
				Name: virtualSpaceMap["name"].(string),
				Size: virtualSpaceMap["size"].(string),
			}

			if virtualSpaceMap["lvm_lv_type"].(string) != "" {
				lvmConfig := nodes.LVMConfigSpec{
					LvType: virtualSpaceMap["lvm_lv_type"].(string),
					Path:   virtualSpaceMap["lvm_path"].(string),
				}
				virtualSpace.LVMConfig = &lvmConfig
			}

			if virtualSpaceMap["runtime_lv_type"].(string) != "" {
				runtimeConfig := nodes.RuntimeConfigSpec{
					LvType: virtualSpaceMap["runtime_lv_type"].(string),
				}
				virtualSpace.RuntimeConfig = &runtimeConfig
			}

			virtualSpaces = append(virtualSpaces, virtualSpace)
		}
		group.VirtualSpaces = virtualSpaces

		groups = append(groups, group)
	}

	storageSpec.StorageGroups = groups
	return &storageSpec
}

func flattenResourceNodeStorage(storageRaw *nodes.StorageSpec) []map[string]interface{} {
	if storageRaw == nil {
		return nil
	}

	storageSelectorsRaw := storageRaw.StorageSelectors
	storageSelectors := make([]map[string]interface{}, len(storageSelectorsRaw))
	for i, s := range storageSelectorsRaw {
		storageSelector := map[string]interface{}{
			"name": s.Name,
			"type": s.StorageType,
		}

		if s.MatchLabels != (nodes.MatchLabelsSpec{}) {
			storageSelector["match_label_size"] = s.MatchLabels.Size
			storageSelector["match_label_volume_type"] = s.MatchLabels.VolumeType
			storageSelector["match_label_metadata_encrypted"] = s.MatchLabels.MetadataEncrypted
			storageSelector["match_label_metadata_cmkid"] = s.MatchLabels.MetadataCmkid
			storageSelector["match_label_count"] = s.MatchLabels.Count
		}
		storageSelectors[i] = storageSelector
	}

	storageGroupsRaw := storageRaw.StorageGroups
	storageGroups := make([]map[string]interface{}, len(storageGroupsRaw))
	for i, v := range storageGroupsRaw {
		storageGroup := map[string]interface{}{
			"name":           v.Name,
			"cce_managed":    v.CceManaged,
			"selector_names": v.SelectorNames,
		}

		virtualSpaces := make([]map[string]interface{}, len(v.VirtualSpaces))
		for k, s := range v.VirtualSpaces {
			virtualSpace := map[string]interface{}{
				"name": s.Name,
				"size": s.Size,
			}

			if s.LVMConfig != nil {
				virtualSpace["lvm_lv_type"] = s.LVMConfig.LvType
				virtualSpace["lvm_path"] = s.LVMConfig.Path
			}
			if s.RuntimeConfig != nil {
				virtualSpace["runtime_lv_type"] = s.RuntimeConfig.LvType
			}

			virtualSpaces[k] = virtualSpace
		}
		storageGroup["virtual_spaces"] = virtualSpaces

		storageGroups[i] = storageGroup
	}

	return []map[string]interface{}{
		{
			"selectors": storageSelectors,
			"groups":    storageGroups,
		},
	}
}

func flattenResourceNodeTaints(taints []nodes.TaintSpec) []map[string]interface{} {
	if len(taints) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(taints))

	for i, v := range taints {
		res[i] = map[string]interface{}{
			"key":    utils.PathSearch("key", v, nil),
			"value":  utils.PathSearch("value", v, nil),
			"effect": utils.PathSearch("effect", v, nil),
		}
	}

	return res
}

func schemaChargingMode(conflicts []string) *schema.Schema {
	resourceSchema := schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ValidateFunc: validation.StringInSlice([]string{
			"prePaid", "postPaid",
		}, false),
		ConflictsWith: conflicts,
	}

	return &resourceSchema
}

func schemaPeriodUnit(conflicts []string) *schema.Schema {
	resourceSchema := schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		RequiredWith: []string{"period"},
		ValidateFunc: validation.StringInSlice([]string{
			"month", "year",
		}, false),
		ConflictsWith: conflicts,
	}

	return &resourceSchema
}

func schemaPeriod(conflicts []string) *schema.Schema {
	resourceSchema := schema.Schema{
		Type:          schema.TypeInt,
		Optional:      true,
		RequiredWith:  []string{"period_unit"},
		ValidateFunc:  validation.IntBetween(1, 9),
		ConflictsWith: conflicts,
	}

	return &resourceSchema
}

func schemaAutoRenew(conflicts []string) *schema.Schema {
	resourceSchema := schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ValidateFunc: validation.StringInSlice([]string{
			"true", "false",
		}, false),
		ConflictsWith: conflicts,
	}

	return &resourceSchema
}

func schemaAutoRenewComputed(conflicts []string) *schema.Schema {
	resourceSchema := schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ValidateFunc: validation.StringInSlice([]string{
			"true", "false",
		}, false),
		ConflictsWith: conflicts,
	}

	return &resourceSchema
}

func schemaAutoPay(conflicts []string) *schema.Schema {
	resourceSchema := schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ValidateFunc: validation.StringInSlice([]string{
			"true", "false",
		}, false),
		ConflictsWith: conflicts,
		Deprecated:    "Deprecated",
	}

	return &resourceSchema
}
