package modelarts

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v2/{project_id}/pools/{pool_name}/nodepools
func DataSourceV2NodePools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2NodePoolsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the node pools are located.`,
			},
			"pool_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource pool name.`,
			},
			"node_pools": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolsSchema(),
				Description: `The list of the node pools.`,
			},
		},
	}
}

func dataSourceV2NodePoolsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolsMetadataSchema(),
				Description: `The metadata information of the node pool.`,
			},
			"spec": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolsSpecSchema(),
				Description: `The expectation information of the node pool.`,
			},
			"status": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolsStatusSchema(),
				Description: `The status information of the node pool.`,
			},
		},
	}
}

func dataSourceV2NodePoolsMetadataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the node pool.`,
			},
		},
	}
}

func dataSourceV2NodePoolsSpecSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolsSpecResourceSchema(),
				Description: `The list of resources in the node pool.`,
			},
		},
	}
}

func dataSourceV2NodePoolsSpecResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource flavor name.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The desired usage of the flavor.`,
			},
			"max_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The elastic usage of the resource flavor.`,
			},
			"azs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolsResourceAvailabilityZonesSchema(),
				Description: `The availability zones information of nodes in the resource pool.`,
			},
			"node_pool": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the node pool.`,
			},
			"taints": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolsSpecResourceTaintsSchema(),
				Description: `The taints information.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The Kubernetes label.`,
			},
			"tags": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The resource tags.`,
			},
			"network": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolsSpecResourceNetworkSchema(),
				Description: `The network configuration.`,
			},
			"extend_params": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolsSpecResourceExtendParamsSchema(),
				Description: `The custom configuration.`,
			},
			"creating_step": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolsSpecResourceCreatingStepSchema(),
				Description: `The information about batch creation.`,
			},
			"root_volume": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolsSpecResourceRootVolumeSchema(),
				Description: `The custom system disk information.`,
			},
			"data_volumes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolsSpecResourceDataVolumesSchema(),
				Description: `The custom data disks information.`,
			},
			"volume_group_configs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolsSpecResourceVolumeGroupConfigSchema(),
				Description: `The advanced disk configuration.`,
			},
			"os": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolsSpecOsSchema(),
				Description: `The OS image information.`,
			},
		},
	}
}

func dataSourceV2NodePoolsResourceAvailabilityZonesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"az": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The availability zone name.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of availability zone resource instances.`,
			},
		},
	}
}

func dataSourceV2NodePoolsSpecResourceTaintsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The taint key.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The taint value.`,
			},
			"effect": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The effect of the action.`,
			},
		},
	}
}

func dataSourceV2NodePoolsSpecResourceNetworkSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vpc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the VPC.`,
			},
			"subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the subnet.`,
			},
			"security_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The IDs of the security groups.`,
			},
		},
	}
}

func dataSourceV2NodePoolsSpecResourceExtendParamsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"docker_base_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The container image space size of a node.`,
			},
			"post_install": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The post-installation script.`,
			},
			"runtime": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The container runtime.`,
			},
			"label_policy_on_existing_nodes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Kubernetes label update policy of existing nodes.`,
			},
			"taint_policy_on_existing_nodes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Kubernetes taint update policy of existing nodes.`,
			},
			"tag_policy_on_existing_nodes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource tag update policy of existing nodes.`,
			},
			"x_parameter_plane_subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subnet ID used for data transmission on the parameter plane between physical clusters.`,
			},
			"node_pool_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the node pool specified by user.`,
			},
		},
	}
}

func dataSourceV2NodePoolsSpecResourceCreatingStepSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"step": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The step of a supernode.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The batch creation type.`,
			},
		},
	}
}

func dataSourceV2NodePoolsSpecResourceRootVolumeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The disk type.`,
			},
			"size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The disk size.`,
			},
		},
	}
}

func dataSourceV2NodePoolsSpecResourceDataVolumesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The disk type.`,
			},
			"size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The disk size.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of disks.`,
			},
			"extend_params": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolsSpecResourceDataVolumesExtendParamsSchema(),
				Description: `The custom disk configuration.`,
			},
		},
	}
}

func dataSourceV2NodePoolsSpecResourceDataVolumesExtendParamsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"volume_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the disk group.`,
			},
		},
	}
}

func dataSourceV2NodePoolsSpecResourceVolumeGroupConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"volume_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The disk group name.`,
			},
			"docker_thin_pool": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The percentage of container disks to data disks on nodes in a resource pool.`,
			},
			"lvm_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolsSpecResourceVolumeGroupConfigLvmConfigSchema(),
				Description: `The LVM configuration.`,
			},
			"types": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The storage type.`,
			},
		},
	}
}

func dataSourceV2NodePoolsSpecResourceVolumeGroupConfigLvmConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"lv_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The LVM write mode.`,
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The disk mount path.`,
			},
		},
	}
}

func dataSourceV2NodePoolsSpecOsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The OS name and version.`,
			},
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The OS image ID.`,
			},
			"image_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The OS image type.`,
			},
			"auto_match": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The automatic OS image matching configuration.`,
			},
		},
	}
}

func dataSourceV2NodePoolsStatusSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolsStatusResourcesSchema(),
				Description: `The resources in different states in the node pool.`,
			},
		},
	}
}

func dataSourceV2NodePoolsStatusResourcesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"creating": {
				Type:        schema.TypeList,
				Elem:        dataSourceV2NodePoolsStatusResourcesFlavorSchema(),
				Computed:    true,
				Description: `The number of resources that are being created.`,
			},
			"available": {
				Type:        schema.TypeList,
				Elem:        dataSourceV2NodePoolsStatusResourcesFlavorSchema(),
				Computed:    true,
				Description: `The number of available resources.`,
			},
			"abnormal": {
				Type:        schema.TypeList,
				Elem:        dataSourceV2NodePoolsStatusResourcesFlavorSchema(),
				Computed:    true,
				Description: `The number of abnormal resources.`,
			},
			"deleting": {
				Type:        schema.TypeList,
				Elem:        dataSourceV2NodePoolsStatusResourcesFlavorSchema(),
				Computed:    true,
				Description: `The number of resources that are being deleted.`,
			},
		},
	}
}

func dataSourceV2NodePoolsStatusResourcesFlavorSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource flavor ID.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of resource specification instances in the resource pool.`,
			},
			"max_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of elastic resource specification instances in the resource pool.`,
			},
			"azs": {
				Type:        schema.TypeList,
				Elem:        dataSourceV2NodePoolsStatusResourcesCountAZSchema(),
				Computed:    true,
				Description: `The AZ distribution of the resource specification instances to be created in the resource pool.`,
			},
			"node_pool": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The node pool ID.`,
			},
		},
	}
}

func dataSourceV2NodePoolsStatusResourcesCountAZSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"az": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The availability zone name.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of availability zone resource instances.`,
			},
		},
	}
}

func listV2NodePools(client *golangsdk.ServiceClient, poolName string) ([]interface{}, error) {
	httpUrl := "v2/{project_id}/pools/{pool_name}/nodepools"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{pool_name}", poolName)

	listFlavorsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	resp, err := client.Request("GET", listPath, &listFlavorsOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func dataSourceV2NodePoolsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		poolName = d.Get("pool_name").(string)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	pools, err := listV2NodePools(client, poolName)
	if err != nil {
		return diag.Errorf("error retrieving node pools: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("node_pools", flattenV2NodePools(pools)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenV2NodePools(pools []interface{}) []map[string]interface{} {
	if pools == nil {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(pools))
	for _, pool := range pools {
		result = append(result, map[string]interface{}{
			"metadata": flattenV2NodePoolsMetadata(utils.PathSearch("metadata", pool, nil)),
			"spec":     flattenV2NodePoolsSpec(utils.PathSearch("spec", pool, nil)),
			"status":   flattenV2NodePoolsStatus(utils.PathSearch("status", pool, nil)),
		})
	}

	return result
}

func flattenV2NodePoolsMetadata(metadata interface{}) []map[string]interface{} {
	if metadata == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name": utils.PathSearch("name", metadata, nil),
		},
	}
}

func flattenV2NodePoolsSpec(spec interface{}) []map[string]interface{} {
	if spec == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"resources": flattenV2NodePoolsSpecResources(utils.PathSearch("resources", spec, nil)),
		},
	}
}

func flattenV2NodePoolsSpecResources(resources interface{}) []map[string]interface{} {
	if resources == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"flavor":    utils.PathSearch("flavor", resources, nil),
			"count":     utils.PathSearch("count", resources, nil),
			"max_count": utils.PathSearch("maxCount", resources, nil),
			"azs": flattenV2NodePoolsSpecResourcesAvailabilityZones(
				utils.PathSearch("azs", resources, make([]interface{}, 0)).([]interface{})),
			"node_pool": utils.PathSearch("nodePool", resources, nil),
			"taints": flattenV2NodePoolsSpecResourcesTaints(
				utils.PathSearch("taints", resources, make([]interface{}, 0)).([]interface{})),
			"labels":        utils.PathSearch("labels", resources, nil),
			"tags":          utils.FlattenTagsToMap(utils.PathSearch("tags", resources, nil)),
			"network":       flattenV2NodePoolsSpecResourcesNetwork(utils.PathSearch("network", resources, nil)),
			"extend_params": flattenV2NodePoolsSpecResourcesExtendParams(utils.PathSearch("extendParams", resources, nil)),
			"creating_step": flattenV2NodePoolsSpecResourcesCreatingStep(utils.PathSearch("creatingStep", resources, nil)),
			"root_volume":   flattenV2NodePoolsSpecResourcesRootVolume(utils.PathSearch("rootVolume", resources, nil)),
			"data_volumes": flattenV2NodePoolsSpecResourcesDataVolumes(
				utils.PathSearch("dataVolumes", resources, make([]interface{}, 0)).([]interface{})),
			"volume_group_configs": flattenV2NodePoolsSpecResourcesVolumeGroupConfigs(
				utils.PathSearch("volumeGroupConfigs", resources, make([]interface{}, 0)).([]interface{})),
			"os": flattenV2NodePoolsSpecOs(utils.PathSearch("os", resources, nil)),
		},
	}
}

func flattenV2NodePoolsSpecResourcesAvailabilityZones(azsInfo []interface{}) []interface{} {
	if len(azsInfo) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(azsInfo))
	for _, v := range azsInfo {
		rst = append(rst, map[string]interface{}{
			"az":    utils.PathSearch("az", v, nil),
			"count": utils.PathSearch("count", v, nil),
		})
	}

	return rst
}

func flattenV2NodePoolsSpecResourcesTaints(taintsInfo []interface{}) []interface{} {
	if len(taintsInfo) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(taintsInfo))
	for _, v := range taintsInfo {
		rst = append(rst, map[string]interface{}{
			"key":    utils.PathSearch("key", v, nil),
			"value":  utils.PathSearch("value", v, nil),
			"effect": utils.PathSearch("effect", v, nil),
		})
	}

	return rst
}

func flattenV2NodePoolsSpecResourcesNetwork(networkInfo interface{}) []map[string]interface{} {
	if networkInfo == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"vpc":             utils.PathSearch("vpc", networkInfo, nil),
			"subnet":          utils.PathSearch("subnet", networkInfo, nil),
			"security_groups": utils.PathSearch("securityGroups", networkInfo, nil),
		},
	}
}

func flattenV2NodePoolsSpecResourcesExtendParams(extendParamsInfo interface{}) []map[string]interface{} {
	if extendParamsInfo == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"docker_base_size":               utils.PathSearch("dockerBaseSize", extendParamsInfo, nil),
			"post_install":                   utils.PathSearch("postInstall", extendParamsInfo, nil),
			"runtime":                        utils.PathSearch("runtime", extendParamsInfo, nil),
			"label_policy_on_existing_nodes": utils.PathSearch("labelPolicyOnExistingNodes", extendParamsInfo, nil),
			"taint_policy_on_existing_nodes": utils.PathSearch("taintPolicyOnExistingNodes", extendParamsInfo, nil),
			"tag_policy_on_existing_nodes":   utils.PathSearch("tagPolicyOnExistingNodes", extendParamsInfo, nil),
			"x_parameter_plane_subnet":       utils.PathSearch("xParameterPlaneSubnet", extendParamsInfo, nil),
			"node_pool_name":                 utils.PathSearch("nodePoolName", extendParamsInfo, nil),
		},
	}
}

func flattenV2NodePoolsSpecResourcesCreatingStep(creatingStep interface{}) []map[string]interface{} {
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

func flattenV2NodePoolsSpecResourcesRootVolume(rootVolume interface{}) []map[string]interface{} {
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

func flattenV2NodePoolsSpecResourcesDataVolumes(dataVolumes []interface{}) []interface{} {
	if len(dataVolumes) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataVolumes))
	for _, v := range dataVolumes {
		rst = append(rst, map[string]interface{}{
			"volume_type": utils.PathSearch("volumeType", v, nil),
			"size":        utils.PathSearch("size", v, nil),
			"count":       utils.PathSearch("count", v, nil),
			"extend_params": flattenV2NodePoolsSpecResourcesDataVolumesExtendParams(
				utils.PathSearch("extendParams", v, nil)),
		})
	}

	return rst
}

func flattenV2NodePoolsSpecResourcesDataVolumesExtendParams(extendParams interface{}) []map[string]interface{} {
	if extendParams == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"volume_group": utils.PathSearch("volumeGroup", extendParams, nil),
		},
	}
}

func flattenV2NodePoolsSpecResourcesVolumeGroupConfigs(volumeGroupConfigs []interface{}) []interface{} {
	if len(volumeGroupConfigs) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(volumeGroupConfigs))
	for _, v := range volumeGroupConfigs {
		result = append(result, map[string]interface{}{
			"volume_group":     utils.PathSearch("volumeGroup", v, nil),
			"docker_thin_pool": utils.PathSearch("dockerThinPool", v, nil),
			"lvm_config": flattenV2NodePoolsSpecResourcesVolumeGroupConfigLvmConfig(
				utils.PathSearch("lvmConfig", v, nil)),
			"types": utils.PathSearch("types", v, nil),
		})
	}

	return result
}

func flattenV2NodePoolsSpecResourcesVolumeGroupConfigLvmConfig(lvmConfig interface{}) []map[string]interface{} {
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

func flattenV2NodePoolsSpecOs(osInfo interface{}) []map[string]interface{} {
	if osInfo == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":       utils.PathSearch("name", osInfo, nil),
			"image_id":   utils.PathSearch("imageId", osInfo, nil),
			"image_type": utils.PathSearch("imageType", osInfo, nil),
			"auto_match": utils.PathSearch("autoMatch", osInfo, nil),
		},
	}
}

func flattenV2NodePoolsStatus(statusInfo interface{}) []map[string]interface{} {
	if statusInfo == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"resources": flattenV2NodePoolsStatusResources(utils.PathSearch("resources", statusInfo, nil)),
		},
	}
}

func flattenV2NodePoolsStatusResources(resourcesInfo interface{}) []map[string]interface{} {
	if resourcesInfo == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"creating":  flattenV2NodePoolsStatusResourcesFlavor(utils.PathSearch("creating", resourcesInfo, nil)),
			"available": flattenV2NodePoolsStatusResourcesFlavor(utils.PathSearch("available", resourcesInfo, nil)),
			"abnormal":  flattenV2NodePoolsStatusResourcesFlavor(utils.PathSearch("abnormal", resourcesInfo, nil)),
			"deleting":  flattenV2NodePoolsStatusResourcesFlavor(utils.PathSearch("deleting", resourcesInfo, nil)),
		},
	}
}

func flattenV2NodePoolsStatusResourcesFlavor(resourceFlavor interface{}) []map[string]interface{} {
	if resourceFlavor == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"flavor":    utils.PathSearch("flavor", resourceFlavor, nil),
			"count":     utils.PathSearch("count", resourceFlavor, nil),
			"max_count": utils.PathSearch("maxCount", resourceFlavor, nil),
			"azs": flattenV2NodePoolsStatusResourcesFlavorAvailabilityZones(
				utils.PathSearch("azs", resourceFlavor, make([]interface{}, 0)).([]interface{})),
			"node_pool": utils.PathSearch("nodePool", resourceFlavor, nil),
		},
	}
}

func flattenV2NodePoolsStatusResourcesFlavorAvailabilityZones(azsInfo []interface{}) []interface{} {
	if len(azsInfo) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(azsInfo))
	for _, v := range azsInfo {
		rst = append(rst, map[string]interface{}{
			"az":    utils.PathSearch("az", v, nil),
			"count": utils.PathSearch("count", v, nil),
		})
	}

	return rst
}
