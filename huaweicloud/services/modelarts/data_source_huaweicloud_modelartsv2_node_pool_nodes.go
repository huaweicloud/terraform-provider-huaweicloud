package modelarts

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v2/{project_id}/pools/{pool_name}/nodepools/{nodepool_name}/nodes
func DataSourceV2NodePoolNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2NodePoolNodesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the resource nodes are located.`,
			},
			"resource_pool_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource pool name to which the node pool belongs.`,
			},
			"node_pool_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The node pool name to which the resource nodes belongs.`,
			},
			"nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolNodeSchema(),
				Description: `All queried resource nodes under a specified node pool.`,
			},
		},
	}
}

func dataSourceV2NodePoolNodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolNodeMetadataSchema(),
				Description: `The metadata information of the node.`,
			},
			"spec": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolNodeSpecSchema(),
				Description: `The specification of the node.`,
			},
			"status": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolNodeStatusSchema(),
				Description: `The status information of the node.`,
			},
		},
	}
}

func dataSourceV2NodePoolNodeMetadataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the node.`,
			},
			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation timestamp of the node.`,
			},
			"labels": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The labels of the node, in JSON format.`,
			},
			"annotations": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The annotation configuration of the node, in JSON format.`,
			},
		},
	}
}

func dataSourceV2NodePoolNodeSpecSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor of the node.`,
			},
			"extend_params": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The extend parameters of the node, in JSON format.`,
			},
			"host_network": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolNodeSpecHostNetworkSchema(),
				Description: `The network configuration of the node.`,
			},
			"os": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2NodePoolNodeSpecOsSchema(),
				Description: `The OS information of the node.`,
			},
		},
	}
}

func dataSourceV2NodePoolNodeSpecHostNetworkSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vpc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The VPC ID to which the node belongs.`,
			},
			"subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subnet ID to which the node belongs.`,
			},
			"security_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The security group IDs that the node used.`,
			},
		},
	}
}

func dataSourceV2NodePoolNodeSpecOsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The OS name of the node.`,
			},
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The image ID of the OS.`,
			},
			"image_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The image type of the OS.`,
			},
		},
	}
}

func dataSourceV2NodePoolNodeStatusSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"phase": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current phase of the node.`,
			},
			"az": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The availability zone where the node is located.`,
			},
			"driver": {
				Type:        schema.TypeList,
				Elem:        dataSourceV2NodePoolNodeStatusDriverSchema(),
				Computed:    true,
				Description: `The driver configuration of the node.`,
			},
			"os": {
				Type:        schema.TypeList,
				Elem:        dataSourceV2NodePoolNodeStatusOsSchema(),
				Computed:    true,
				Description: `The OS information of the kubernetes node.`,
			},
			"plugins": {
				Type:        schema.TypeList,
				Elem:        dataSourceV2NodePoolNodeStatusPluginsSchema(),
				Computed:    true,
				Description: `The plugin configuration of the node.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The private IP address of the node.`,
			},
			"resources": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource detail of the node, in JSON format.`,
			},
			"available_resources": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The available resource detail of the node, in JSON format.`,
			},
		},
	}
}

func dataSourceV2NodePoolNodeStatusDriverSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"phase": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current phase of the driver.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the driver.`,
			},
		},
	}
}

func dataSourceV2NodePoolNodeStatusOsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The OS name of the kubernetes node.`,
			},
		},
	}
}

func dataSourceV2NodePoolNodeStatusPluginsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the plugin.`,
			},
			"phase": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current phase of the plugin.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the plugin.`,
			},
		},
	}
}

func listV2NodePoolNodes(client *golangsdk.ServiceClient, resourcePoolName, nodePoolName string) ([]interface{}, error) {
	// Currently, these query parameters are not available:
	// + continue
	// + limit
	// Without page parameters, the service will returns all of nodes.
	httpUrl := "v2/{project_id}/pools/{pool_name}/nodepools/{nodepool_name}/nodes"

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{pool_name}", resourcePoolName)
	listPath = strings.ReplaceAll(listPath, "{nodepool_name}", nodePoolName)

	listFlavorsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	requestResp, err := client.Request("GET", listPath, &listFlavorsOpt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenV2NodePoolNodeMetadata(metadata interface{}) []map[string]interface{} {
	if metadata == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":               utils.PathSearch("name", metadata, nil),
			"creation_timestamp": utils.PathSearch("creationTimestamp", metadata, nil),
			"labels":             utils.JsonToString(utils.PathSearch("labels", metadata, nil)),
			"annotations":        utils.JsonToString(utils.PathSearch("annotations", metadata, nil)),
		},
	}
}

func flattenV2NodePoolNodeSpecHostNetwork(os interface{}) []map[string]interface{} {
	if os == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"vpc":             utils.PathSearch("vpc", os, nil),
			"subnet":          utils.PathSearch("subnet", os, nil),
			"security_groups": utils.PathSearch("securityGroups", os, make([]interface{}, 0)).([]interface{}),
		},
	}
}

func flattenV2NodePoolNodeSpecOs(osInfo interface{}) []map[string]interface{} {
	if osInfo == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":       utils.PathSearch("name", osInfo, nil),
			"image_id":   utils.PathSearch("imageId", osInfo, nil),
			"image_type": utils.PathSearch("imageType", osInfo, nil),
		},
	}
}

func flattenV2NodePoolNodeSpec(spec interface{}) []map[string]interface{} {
	if spec == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"flavor":        utils.PathSearch("flavor", spec, nil),
			"extend_params": utils.JsonToString(utils.PathSearch("extendParams", spec, nil)),
			"host_network":  flattenV2NodePoolNodeSpecHostNetwork(utils.PathSearch("hostNetwork", spec, nil)),
			"os":            flattenV2NodePoolNodeSpecOs(utils.PathSearch("os", spec, nil)),
		},
	}
}

func flattenV2NodePoolNodeDriver(driver interface{}) []map[string]interface{} {
	if driver == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"phase":   utils.PathSearch("phase", driver, nil),
			"version": utils.PathSearch("version", driver, nil),
		},
	}
}

func flattenV2NodePoolNodeStatusOs(osInfo interface{}) []map[string]interface{} {
	if osInfo == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name": utils.PathSearch("name", osInfo, nil),
		},
	}
}

func flattenV2NodePoolNodePlugins(plugins []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(plugins))

	for _, plugin := range plugins {
		result = append(result, map[string]interface{}{
			"name":    utils.PathSearch("name", plugin, nil),
			"phase":   utils.PathSearch("phase", plugin, nil),
			"version": utils.PathSearch("version", plugin, nil),
		})
	}

	return result
}

func flattenV2NodePoolNodeStatus(status interface{}) []map[string]interface{} {
	if status == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"phase":               utils.PathSearch("phase", status, nil),
			"az":                  utils.PathSearch("az", status, nil),
			"driver":              flattenV2NodePoolNodeDriver(utils.PathSearch("driver", status, nil)),
			"os":                  flattenV2NodePoolNodeStatusOs(utils.PathSearch("os", status, nil)),
			"plugins":             flattenV2NodePoolNodePlugins(utils.PathSearch("plugins", status, make([]interface{}, 0)).([]interface{})),
			"private_ip":          utils.PathSearch("privateIp", status, nil),
			"resources":           utils.JsonToString(utils.PathSearch("resources", status, nil)),
			"available_resources": utils.JsonToString(utils.PathSearch("availableResources", status, nil)),
		},
	}
}

func flattenV2NodePoolNodes(nodes []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(nodes))

	for _, node := range nodes {
		result = append(result, map[string]interface{}{
			"metadata": flattenV2NodePoolNodeMetadata(utils.PathSearch("metadata", node, nil)),
			"spec":     flattenV2NodePoolNodeSpec(utils.PathSearch("spec", node, nil)),
			"status":   flattenV2NodePoolNodeStatus(utils.PathSearch("status", node, nil)),
		})
	}

	return result
}

func dataSourceV2NodePoolNodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		resourcePoolName = d.Get("resource_pool_name").(string)
		nodePoolName     = d.Get("node_pool_name").(string)
	)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	nodes, err := listV2NodePoolNodes(client, resourcePoolName, nodePoolName)
	if err != nil {
		return diag.Errorf("error querying node list for the specified node pool (%s) under the resource pool (%s): %s",
			nodePoolName, resourcePoolName, err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("nodes", flattenV2NodePoolNodes(nodes)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
