package dataarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/service/instances
func DataSourceDataServiceInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataServiceInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the exclusive clusters are located.`,
			},

			// Parameters in request header
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the exclusive clusters belong.`,
			},

			// Query arguments
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The exclusive cluster name to be queried.`,
			},
			"create_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The create user name of the exclusive clusters to be queried.`,
			},

			// Attributes
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataInstancesElem(),
				Description: `All exclusive clusters that match the filter parameters.`,
			},
		},
	}
}

func dataInstancesElem() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the exclusive cluster, in UUID format.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the exclusive cluster.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the exclusive cluster.`,
			},
			"external_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The external IP address of the exclusive cluster.`,
			},
			"intranet_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The intranet IP address of the exclusive cluster.`,
			},
			"intranet_address_ipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The intranet IPv6 address of the exclusive cluster.`,
			},
			"public_zone_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The public zone ID of the exclusive cluster.`,
			},
			"public_zone_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The public zone name of the exclusive cluster.`,
			},
			"private_zone_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The private zone ID of the exclusive cluster.`,
			},
			"private_zone_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The private zone name of the exclusive cluster.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise project ID to which the exclusive cluster belongs.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The create time of the exclusive cluster, in RFC3339 format.`,
			},
			"create_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The create user of the exclusive cluster.`,
			},
			"current_namespace_publish_api_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of the published APIs in the current namespace.`,
			},
			"all_namespace_publish_api_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of the published APIs.`,
			},
			"api_publishable_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The API quota of the exclusive cluster.`,
			},
			"deletable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the exclusive cluster can be deleted.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the exclusive cluster.`,
			},
			"flavor": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataInstanceFlavorSchema(),
				Description: `The flavor of the exclusive cluster.`,
			},
			"gateway_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the exclusive cluster.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The availability zone where the exclusive cluster is located.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The VPC ID to which the exclusive cluster belongs.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subnet ID to which the exclusive cluster belongs.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The security group ID associated to the exclusive cluster.`,
			},
			"node_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The node number of the exclusive cluster.`,
			},
			"nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataInstanceNodeSchema(),
				Description: `The list of instance nodes.`,
			},
		},
	}
	return &sc
}

func dataInstanceFlavorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor name.`,
			},
			"disk_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of the disk size.`,
			},
			"vcpus": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of CPU cores in the flavor.`,
			},
			"memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The memory size in the flavor, in GB.`,
			},
		},
	}
	return &sc
}

func dataInstanceNodeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The node ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The node name.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The private IP address of the node.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the node.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The create time of the node, in RFC3339 format.`,
			},
			"create_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The create user of the node.`,
			},
			"gateway_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the node.`,
			},
		},
	}
	return &sc
}

func buildDataServiceInstancesQueryParams(d *schema.ResourceData) string {
	res := ""
	if apiName, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, apiName)
	}
	if apiType, ok := d.GetOk("create_user"); ok {
		res = fmt.Sprintf("%s&create_user=%v", res, apiType)
	}
	return res
}

func getDataServiceInstances(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl     = "v1/{project_id}/service/instances?limit=100"
		workspaceId = d.Get("workspace_id").(string)
		offset      = 0
		result      = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildDataServiceInstancesQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
			"Dlm-Type":     "EXCLUSIVE",
		},
	}

	for {
		listPathWithOffsset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffsset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		instances := utils.PathSearch("instances", respBody, make([]interface{}, 0)).([]interface{})
		if len(instances) < 1 {
			break
		}
		result = append(result, instances...)
		offset += len(instances)
	}

	return result, nil
}

func flattenInstanceFlavor(flavor map[string]interface{}) []map[string]interface{} {
	if len(flavor) < 1 {
		return nil
	}
	return []map[string]interface{}{
		{
			"id":        utils.PathSearch("id", flavor, nil),
			"name":      utils.PathSearch("name", flavor, nil),
			"disk_size": utils.PathSearch("disk", flavor, nil),
			"vcpus":     utils.PathSearch("cpu", flavor, nil),
			"memory":    utils.PathSearch("mem", flavor, nil),
		},
	}
}

func flattenInstanceNodes(nodes []interface{}) []interface{} {
	if len(nodes) < 1 {
		return nil
	}
	result := make([]interface{}, 0, len(nodes))
	for _, node := range nodes {
		result = append(result, map[string]interface{}{
			"id":              utils.PathSearch("id", node, nil),
			"name":            utils.PathSearch("name", node, nil),
			"private_ip":      utils.PathSearch("private_ip", node, nil),
			"status":          utils.PathSearch("status", node, nil),
			"gateway_version": utils.PathSearch("gateway_version", node, nil),
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_time",
				node, "").(string))/1000, false),
			"create_user": utils.PathSearch("create_user", node, nil),
		})
	}
	return result
}

func flattenInstances(instances []interface{}) []interface{} {
	result := make([]interface{}, 0, len(instances))
	for _, instance := range instances {
		result = append(result, map[string]interface{}{
			"id":                                utils.PathSearch("id", instance, nil),
			"name":                              utils.PathSearch("name", instance, nil),
			"description":                       utils.PathSearch("description", instance, nil),
			"external_address":                  utils.PathSearch("external_address", instance, nil),
			"intranet_address":                  utils.PathSearch("intranet_address", instance, nil),
			"intranet_address_ipv6":             utils.PathSearch("intranet_address_ipv6", instance, nil),
			"public_zone_id":                    utils.PathSearch("public_zone_id", instance, nil),
			"public_zone_name":                  utils.PathSearch("public_zone_name", instance, nil),
			"private_zone_id":                   utils.PathSearch("private_zone_id", instance, nil),
			"private_zone_name":                 utils.PathSearch("private_zone_name", instance, nil),
			"enterprise_project_id":             utils.PathSearch("enterprise_project_id", instance, nil),
			"current_namespace_publish_api_num": utils.PathSearch("current_namespace_publish_api_num", instance, nil),
			"all_namespace_publish_api_num":     utils.PathSearch("all_namespace_publish_api_num", instance, nil),
			"api_publishable_num":               utils.PathSearch("api_publishable_num", instance, nil),
			"deletable":                         utils.PathSearch("deletable", instance, nil),
			"status":                            utils.PathSearch("instance_status", instance, nil),
			"node_num":                          utils.PathSearch("node_num", instance, nil),
			"gateway_version":                   utils.PathSearch("gateway_version", instance, nil),
			"availability_zone":                 utils.PathSearch("availability_zone", instance, nil),
			"vpc_id":                            utils.PathSearch("vpc_id", instance, nil),
			"subnet_id":                         utils.PathSearch("subnet_id", instance, nil),
			"security_group_id":                 utils.PathSearch("security_group_id", instance, nil),
			"created_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time",
				instance, float64(0)).(float64))/1000, false),
			"create_user": utils.PathSearch("create_user", instance, nil),
			"flavor": flattenInstanceFlavor(utils.PathSearch("flavor",
				instance, make(map[string]interface{})).(map[string]interface{})),
			"nodes": flattenInstanceNodes(utils.PathSearch("nodes",
				instance, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func dataSourceDataServiceInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	instances, err := getDataServiceInstances(client, d)
	if err != nil {
		return diag.Errorf("error getting Data Service exclusive instances for DataArts Studio: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", flattenInstances(instances)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
