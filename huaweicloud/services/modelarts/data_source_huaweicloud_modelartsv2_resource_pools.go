package modelarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ServiceStage GET /v2/{project_id}/pools
func DataSourceV2ResourcePools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2ResourcePoolsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the resource pools are located.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The workspace ID to which the resource pool belongs.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status of the resource pools to be queried.`,
			},
			"resource_pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metadata": {
							Type:        schema.TypeList,
							Elem:        dataV2ResourcePoolMetadataSchema(),
							Computed:    true,
							Description: `The metadata configuration of the resource pool.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the resource pool.`,
						},
						"scope": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: `The list of job types supported by the resource pool.`,
						},
						"resources": {
							Type:        schema.TypeList,
							Elem:        dataV2ResourcePoolResourceSchema(),
							Computed:    true,
							Description: `The list of resource specifications in the resource pool.`,
						},
						"network_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ModelArts network ID of the resource pool.`,
						},
						"prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The prefix of the user-defined node name of the resource pool.`,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the VPC to which the resource pool belongs.`,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The network ID of the subnet to which the resource pool belongs.`,
						},
						"clusters": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataV2ResourcePoolClustersSchema(),
							Description: `The list of the CCE clusters.`,
						},
						"user_login": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataV2ResourcePoolUserLoginSchema(),
							Description: `The user login info of the resource pool.`,
						},
						"workspace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The workspace ID of the resource pool.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the resource pool.`,
						},
						"charging_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The charging mode of the resource pool.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the resource pool.`,
						},
						"resource_pool_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource ID of the resource pool.`,
						},
					},
				},
				Description: "All application details.",
			},
		},
	}
}

func dataV2ResourcePoolMetadataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the resource pool.`,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The annotations of the resource pool.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels of the resource pool.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the resource pool, in RFC3339 format.`,
			},
		},
	}
}

func dataV2ResourcePoolResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource flavor ID.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of resources of the corresponding flavors.`,
			},
			"node_pool": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of resource pool nodes.`,
			},
			"max_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The max number of resources of the corresponding flavors.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the VPC to which the the resource pool nodes belong.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The network ID of a subnet to which the the resource pool nodes belong.`,
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The security group IDs to which the the resource pool nodes belong.`,
			},
			"azs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2ResourcePoolResourceAzsSchema(),
				Description: `The availability zones for the resource pool nodes.`,
			},
			"taints": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2ResourcePoolResourceTaintsSchema(),
				Description: `The taints added to resource pool nodes.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels of resource pool nodes.`,
			},
			"tags": common.TagsComputedSchema(
				`The key/value pairs to associate with the resource pool nodes.`,
			),
			"extend_params": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The extend parameters of the resource pool nodes.`,
			},
			"root_volume": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2ResourcePoolResourceRootVolumeSchema(),
				Description: `The root volume of the resource pool nodes.`,
			},
			"data_volumes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2ResourcePoolResourceDataVolumesSchema(),
				Description: `The list of data volumes of the resource pool nodes.`,
			},
			"volume_group_configs": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        dataV2ResourcePoolResourceVolumeGroupConfigsSchema(),
				Description: `The extend configurations of the volume groups.`,
			},
			"creating_step": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"step": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The creation step of the resource pool nodes.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the resource pool nodes.`,
						},
					},
				},
				Description: `The creation step configuration of the resource pool nodes.`,
			},
		},
	}
}

func dataV2ResourcePoolResourceAzsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"az": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The AZ name.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of nodes for the corresponding AZ.`,
			},
		},
	}
}

func dataV2ResourcePoolResourceTaintsSchema() *schema.Resource {
	return &schema.Resource{
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
	}
}

func dataV2ResourcePoolResourceRootVolumeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the root volume.`,
			},
			"size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The size of the root volume.`,
			},
		},
	}
}

func dataV2ResourcePoolResourceDataVolumesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the data volume.`,
			},
			"size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The size of the data volume.`,
			},
			"extend_params": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The extend parameters of the data volume.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The count of the current data volume configuration.`,
			},
		},
	}
}

func dataV2ResourcePoolResourceVolumeGroupConfigsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"volume_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the volume group.`,
			},
			"docker_thin_pool": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The percentage of container volumes to data volumes on resource pool nodes.`,
			},
			"lvm_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lv_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The LVM write mode.`,
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The volume mount path.`,
						},
					},
				},
				Description: `The configuration of the LVM management.`,
			},
			"types": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of storage types of the volume group.`,
			},
		},
	}
}

func dataV2ResourcePoolClustersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the CCE cluster that resource pool used.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the CCE cluster that resource pool used.`,
			},
		},
	}
}

func dataV2ResourcePoolUserLoginSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key_pair_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The key pair name of the login user.`,
			},
		},
	}
}

func buildV2ResourcePoolsQueryParams(d *schema.ResourceData) string {
	result := ""

	if workspaceId, ok := d.GetOk("workspace_id"); ok {
		result += fmt.Sprintf("&workspaceId=%v", workspaceId)
	}

	if status, ok := d.GetOk("status"); ok {
		result += fmt.Sprintf("&status=%v", status)
	}

	if result == "" {
		return result
	}
	return "?" + result[1:]
}

func listV2ResourcePools(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/pools"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildV2ResourcePoolsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenV2ResourcePoolMetadata(metadata interface{}) []map[string]interface{} {
	if metadata == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":        utils.PathSearch("name", metadata, nil),
			"annotations": utils.PathSearch("annotations", metadata, make(map[string]interface{})),
			"labels":      utils.PathSearch("labels", metadata, make(map[string]interface{})),
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("creationTimestamp",
				metadata, "").(string))/1000, false),
		},
	}
}

func flattenV2ResourcePoolResourceAzs(azList []interface{}) []map[string]interface{} {
	if len(azList) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(azList))
	for _, az := range azList {
		result = append(result, map[string]interface{}{
			"az":    utils.PathSearch("az", az, nil),
			"count": utils.PathSearch("count", az, nil),
		})
	}

	return result
}

func flattenV2ResourcePoolResourceTaints(taints []interface{}) []map[string]interface{} {
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

func flattenV2ResourcePoolResources(resources []interface{}) []map[string]interface{} {
	if len(resources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(resources))
	for _, resource := range resources {
		result = append(result, map[string]interface{}{
			"flavor_id":          utils.PathSearch("flavor", resource, nil),
			"count":              utils.PathSearch("count", resource, nil),
			"max_count":          utils.PathSearch("maxCount", resource, nil),
			"node_pool":          utils.PathSearch("nodePool", resource, nil),
			"vpc_id":             utils.PathSearch("network.vpc", resource, nil),
			"subnet_id":          utils.PathSearch("network.subnet", resource, nil),
			"security_group_ids": utils.PathSearch("network.securityGroups", resource, nil),
			"azs":                flattenV2ResourcePoolResourceAzs(utils.PathSearch("azs", resource, make([]interface{}, 0)).([]interface{})),
			"taints":             flattenV2ResourcePoolResourceTaints(utils.PathSearch("taints", resource, make([]interface{}, 0)).([]interface{})),
			"labels":             utils.PathSearch("labels", resource, nil),
			"tags":               flattenResourcePoolResourcesTags(resource),
			"extend_params":      utils.JsonToString(utils.PathSearch("extendParams", resource, nil)),
			"root_volume":        flattenResourcePoolResourcesRootVolume(utils.PathSearch("rootVolume", resource, nil)),
			"data_volumes": flattenResourcePoolResourcesDataVolumes(utils.PathSearch("dataVolumes",
				resource, make([]interface{}, 0)).([]interface{})),
			"volume_group_configs": flattenResourcePoolResourcesVolumeGroupConfigs(utils.PathSearch("volumeGroupConfigs",
				resource, make([]interface{}, 0)).([]interface{})),
			"creating_step": flattenResourcePoolResourcesCreatingStep(utils.PathSearch("creatingStep", resource, nil)),
		})
	}

	return result
}

func flattenV2ResourcePoolClusters(clusters []interface{}) []map[string]interface{} {
	if len(clusters) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(clusters))
	for _, cluster := range clusters {
		result = append(result, map[string]interface{}{
			"provider_id": utils.PathSearch("providerId", cluster, nil),
			"name":        utils.PathSearch("name", cluster, nil),
		})
	}

	return result
}

func flattenV2ResourcePoolUserLogin(userLogin interface{}) []map[string]interface{} {
	if userLogin == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"key_pair_name": utils.PathSearch("keyPairName", userLogin, nil),
		},
	}
}

func flattenV2ResourcePoolChargingMode(chargingMode string) interface{} {
	if chargingMode == "1" {
		return "prePaid"
	}

	return nil
}

func flattenV2ResourcePools(resourcePools []interface{}) []map[string]interface{} {
	if len(resourcePools) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(resourcePools))
	for _, resourcePool := range resourcePools {
		result = append(result, map[string]interface{}{
			"metadata": flattenV2ResourcePoolMetadata(utils.PathSearch("metadata", resourcePool, nil)),
			"name":     utils.PathSearch(`metadata.labels."os.modelarts/name"`, resourcePool, nil),
			"scope":    utils.PathSearch("spec.scope", resourcePool, nil),
			"resources": flattenV2ResourcePoolResources(utils.PathSearch("spec.resources",
				resourcePool, make([]interface{}, 0)).([]interface{})),
			"network_id": utils.PathSearch("spec.network.name", resourcePool, nil),
			"prefix":     utils.PathSearch(`metadata.labels."os.modelarts/node.prefix"`, resourcePool, nil),
			"vpc_id":     utils.PathSearch("spec.network.vpcId", resourcePool, nil),
			"subnet_id":  utils.PathSearch("spec.network.subnetId", resourcePool, nil),
			"clusters": flattenV2ResourcePoolClusters(utils.PathSearch("spec.clusters",
				resourcePool, make([]interface{}, 0)).([]interface{})),
			"user_login":   flattenV2ResourcePoolUserLogin(utils.PathSearch("spec.userLogin", resourcePool, nil)),
			"workspace_id": utils.PathSearch(`metadata.labels."os.modelarts/workspace.id"`, resourcePool, nil),
			"description":  utils.PathSearch(`metadata.annotations."os.modelarts/description"`, resourcePool, nil),
			"charging_mode": flattenV2ResourcePoolChargingMode(utils.PathSearch(`metadata.annotations."os.modelarts/billing.mode"`,
				resourcePool, "").(string)),
			"status":           utils.PathSearch("status.phase", resourcePool, nil),
			"resource_pool_id": utils.PathSearch(`metadata.labels."os.modelarts/resource.id"`, resourcePool, nil),
		})
	}

	return result
}

func dataSourceV2ResourcePoolsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resourcePools, err := listV2ResourcePools(client, d)
	if err != nil {
		return diag.Errorf("error getting resource pools: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("resource_pools", flattenV2ResourcePools(resourcePools)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
