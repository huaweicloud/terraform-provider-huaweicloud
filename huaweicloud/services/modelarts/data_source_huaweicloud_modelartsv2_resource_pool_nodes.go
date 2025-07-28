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

// @API ModelArts GET /v2/{project_id}/pools/{pool_name}/nodes
func DataSourceV2ResourcePoolNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2ResourcePoolNodesRead,

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
				Description: `The resource pool name to which the resource nodes belong.`,
			},
			"nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2ResourcePoolNodeSchema(),
				Description: `All queried resource nodes under a specified resource pool.`,
			},
		},
	}
}

func dataSourceV2ResourcePoolNodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2ResourcePoolNodeMetadataSchema(),
				Description: `The metadata information of the node.`,
			},
			"spec": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2ResourcePoolNodeSpecSchema(),
				Description: `The specification of the node.`,
			},
			"status": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2ResourcePoolNodeStatusSchema(),
				Description: `The status information of the node.`,
			},
		},
	}
}

func dataSourceV2ResourcePoolNodeMetadataSchema() *schema.Resource {
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

func dataSourceV2ResourcePoolNodeSpecSchema() *schema.Resource {
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
				Elem:        dataSourceV2ResourcePoolNodeSpecHostNetworkSchema(),
				Description: `The network configuration of the node.`,
			},
			"os": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2ResourcePoolNodeSpecOsSchema(),
				Description: `The OS information of the node.`,
			},
		},
	}
}

func dataSourceV2ResourcePoolNodeSpecHostNetworkSchema() *schema.Resource {
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
		},
	}
}

func dataSourceV2ResourcePoolNodeSpecOsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The image ID of the OS.`,
			},
		},
	}
}

func dataSourceV2ResourcePoolNodeStatusSchema() *schema.Resource {
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
				Elem:        dataSourceV2ResourcePoolNodeStatusDriverSchema(),
				Computed:    true,
				Description: `The driver configuration of the node.`,
			},
			"os": {
				Type:        schema.TypeList,
				Elem:        dataSourceV2ResourcePoolNodeStatusOsSchema(),
				Computed:    true,
				Description: `The OS information of the kubernetes node.`,
			},
			"plugins": {
				Type:        schema.TypeList,
				Elem:        dataSourceV2ResourcePoolNodeStatusPluginsSchema(),
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

func dataSourceV2ResourcePoolNodeStatusDriverSchema() *schema.Resource {
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

func dataSourceV2ResourcePoolNodeStatusOsSchema() *schema.Resource {
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

func dataSourceV2ResourcePoolNodeStatusPluginsSchema() *schema.Resource {
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

func listV2ResourcePoolNodes(client *golangsdk.ServiceClient, resourcePoolName string) ([]interface{}, error) {
	// Currently, these query parameters are not available:
	// + continue
	// + limit
	// Without page parameters, the service will returns all of nodes.
	httpUrl := "v2/{project_id}/pools/{pool_name}/nodes"

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{pool_name}", resourcePoolName)

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

func flattenV2ResourcePoolNodeMetadata(metadata interface{}) []map[string]interface{} {
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

func flattenV2ResourcePoolNodeSpecHostNetwork(os interface{}) []map[string]interface{} {
	if os == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"vpc":    utils.PathSearch("vpc", os, nil),
			"subnet": utils.PathSearch("subnet", os, nil),
		},
	}
}

func flattenV2ResourcePoolNodeSpecOs(osInfo interface{}) []map[string]interface{} {
	if osInfo == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"image_id": utils.PathSearch("imageId", osInfo, nil),
		},
	}
}

func flattenV2ResourcePoolNodeSpec(spec interface{}) []map[string]interface{} {
	if spec == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"flavor":        utils.PathSearch("flavor", spec, nil),
			"extend_params": utils.JsonToString(utils.PathSearch("extendParams", spec, nil)),
			"host_network":  flattenV2ResourcePoolNodeSpecHostNetwork(utils.PathSearch("hostNetwork", spec, nil)),
			"os":            flattenV2ResourcePoolNodeSpecOs(utils.PathSearch("os", spec, nil)),
		},
	}
}

func flattenV2ResourcePoolNodeDriver(driver interface{}) []map[string]interface{} {
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

func flattenV2ResourcePoolNodeStatusOs(osInfo interface{}) []map[string]interface{} {
	if osInfo == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name": utils.PathSearch("name", osInfo, nil),
		},
	}
}

func flattenV2ResourcePoolNodePlugins(plugins []interface{}) []map[string]interface{} {
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

func flattenV2ResourcePoolNodeStatus(status interface{}) []map[string]interface{} {
	if status == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"phase":               utils.PathSearch("phase", status, nil),
			"az":                  utils.PathSearch("az", status, nil),
			"driver":              flattenV2ResourcePoolNodeDriver(utils.PathSearch("driver", status, nil)),
			"os":                  flattenV2ResourcePoolNodeStatusOs(utils.PathSearch("os", status, nil)),
			"plugins":             flattenV2ResourcePoolNodePlugins(utils.PathSearch("plugins", status, make([]interface{}, 0)).([]interface{})),
			"private_ip":          utils.PathSearch("privateIp", status, nil),
			"resources":           utils.JsonToString(utils.PathSearch("resources", status, nil)),
			"available_resources": utils.JsonToString(utils.PathSearch("availableResources", status, nil)),
		},
	}
}

func flattenV2ResourcePoolNodes(nodes []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(nodes))

	for _, node := range nodes {
		result = append(result, map[string]interface{}{
			"metadata": flattenV2ResourcePoolNodeMetadata(utils.PathSearch("metadata", node, nil)),
			"spec":     flattenV2ResourcePoolNodeSpec(utils.PathSearch("spec", node, nil)),
			"status":   flattenV2ResourcePoolNodeStatus(utils.PathSearch("status", node, nil)),
		})
	}

	return result
}

func dataSourceV2ResourcePoolNodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		resourcePoolName = d.Get("resource_pool_name").(string)
	)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	nodes, err := listV2ResourcePoolNodes(client, resourcePoolName)
	if err != nil {
		return diag.Errorf("error querying node list for the specified resource pool (%s): %s",
			resourcePoolName, err,
		)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("nodes", flattenV2ResourcePoolNodes(nodes)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
