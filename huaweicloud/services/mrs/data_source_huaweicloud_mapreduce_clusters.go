// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product MRS
// ---------------------------------------------------------------

package mrs

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

// @API MRS GET /v1.1/{project_id}/cluster_infos
func DataSourceMrsClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceMrsClustersRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of cluster.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status of cluster.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The enterprise project ID used to query clusters in a specified enterprise project.`,
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `You can search for a cluster by its tags.`,
			},
			"clusters": {
				Type:        schema.TypeList,
				Elem:        mrsClustersClustersSchema(),
				Computed:    true,
				Description: `The list of clusters.`,
			},
		},
	}
}

func mrsClustersClustersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Cluster ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Cluster name.`,
			},
			"master_node_num": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Number of Master nodes deployed in a cluster.`,
			},
			"core_node_num": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Number of Core nodes deployed in a cluster.`,
			},
			"total_node_num": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Total number of nodes deployed in a cluster.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Cluster status.`,
			},
			"billing_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Cluster billing mode.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `VPC ID.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Subnet ID.`,
			},
			"duration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Cluster subscription duration.`,
			},
			"fee": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Cluster creation fee, which is automatically calculated.`,
			},
			"hadoop_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Hadoop version.`,
			},
			"master_node_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Instance specifications of a Master node.`,
			},
			"core_node_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Instance specifications of a Core node.`,
			},
			"component_list": {
				Type:        schema.TypeList,
				Elem:        mrsClustersClustersComponentListSchema(),
				Computed:    true,
				Description: `Component list.`,
			},
			"external_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `External IP address.`,
			},
			"external_alternate_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Backup external IP address.`,
			},
			"internal_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Internal IP address.`,
			},
			"deployment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Cluster deployment ID.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Cluster description.`,
			},
			"order_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Cluster creation order ID.`,
			},
			"master_node_product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Product ID of a Master node.`,
			},
			"master_node_spec_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specification ID of a Master node.`,
			},
			"core_node_product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Product ID of a Core node.`,
			},
			"core_node_spec_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specification ID of a Core node.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The AZ.`,
			},
			"vnc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `URI for remotely logging in to an ECS.`,
			},
			"volume_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Disk storage space.`,
			},
			"volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Disk type.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Enterprise project ID.`,
			},
			"type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Cluster type.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Security group ID.`,
			},
			"slave_security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Security group ID of a non-Master node.`,
			},
			"stage_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Cluster progress description.`,
			},
			"safe_mode": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Running mode of an MRS cluster.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Cluster version.`,
			},
			"node_public_cert_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Name of the key file.`,
			},
			"master_node_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `IP address of a Master node.`,
			},
			"private_ip_first": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Preferred private IP address.`,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"log_collection": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Whether to collect logs when cluster installation fails.`,
			},
			"task_node_groups": {
				Type:        schema.TypeList,
				Elem:        mrsClustersClustersNodeGroupSchema(),
				Computed:    true,
				Description: `List of Task nodes.`,
			},
			"node_groups": {
				Type:        schema.TypeList,
				Elem:        mrsClustersClustersNodeGroupSchema(),
				Computed:    true,
				Description: `List of Master, Core and Task nodes.`,
			},
			"master_data_volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Data disk storage type of the Master node.`,
			},
			"master_data_volume_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Data disk storage space of the Master node`,
			},
			"master_data_volume_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Number of data disks of the Master node`,
			},
			"core_data_volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Data disk storage type of the Core node.`,
			},
			"core_data_volume_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Data disk storage space of the Core node.`,
			},
			"core_data_volume_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Number of data disks of the Core node.`,
			},
			"period_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Whether the subscription type is yearly or monthly.`,
			},
			"scale": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Status of node changes.`,
			},
			"eip_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Unique ID of the cluster EIP.`,
			},
			"eip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `IPv4 address of the cluster EIP.`,
			},
			"eipv6_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `IPv6 address of the cluster EIP.`,
			},
			"mrs_ecs_default_agency": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The default agency name bound to the cluster node.`,
			},
		},
	}
	return &sc
}

func mrsClustersClustersComponentListSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"component_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Component ID.`,
			},
			"component_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Component name.`,
			},
			"component_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Component version.`,
			},
			"component_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Component description.`,
			},
		},
	}
	return &sc
}

func mrsClustersClustersNodeGroupSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Node group name.`,
			},
			"node_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Number of nodes in a node group.`,
			},
			"node_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Instance specifications of a node group.`,
			},
			"node_spec_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Instance specification ID of a node group.`,
			},
			"node_product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Instance product ID of a node group.`,
			},
			"vm_product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `VM product ID of a node group.`,
			},
			"vm_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `VM specification code of a node group.`,
			},
			"root_volume_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Root disk storage space of a node group.`,
			},
			"root_volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Root disk storage type of a node group.`,
			},
			"root_volume_product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Root disk product ID of a node group.`,
			},
			"root_volume_resource_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Root disk specification code of a node group.`,
			},
			"root_volume_resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `System disk product type of a node group.`,
			},
			"data_volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Data disk storage type of a node group.`,
			},
			"data_volume_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Number of data disks of a node group.`,
			},
			"data_volume_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Data disk storage space of a node group.`,
			},
			"data_volume_product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Data disk product ID of a node group.`,
			},
			"data_volume_resource_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Data disk specification code of a node group.`,
			},
			"data_volume_resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Data disk product type of a node group.`,
			},
		},
	}
	return &sc
}

func buildClustersQueryParams(d *schema.ResourceData) string {
	res := "?pageSize=100"

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&clusterName=%v", res, v)
	}

	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&clusterState=%v", res, v)
	}

	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterpriseProjectId=%v", res, v)
	}

	if v, ok := d.GetOk("tags"); ok {
		res = fmt.Sprintf("%s&tags=%v", res, v)
	}

	return res
}

func listClusters(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl     = "v1.1/{project_id}/cluster_infos"
		currentPage = 1
		result      = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildClustersQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithPageNum := listPath + fmt.Sprintf("&currentPage=%d", currentPage)
		requestResp, err := client.Request("GET", listPathWithPageNum, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		clusters := utils.PathSearch("clusters", respBody, make([]interface{}, 0)).([]interface{})
		if len(clusters) < 1 {
			break
		}
		result = append(result, clusters...)
		currentPage++
	}

	return result, nil
}

func resourceMrsClustersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("mrs", region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	cluster, err := listClusters(client, d)
	if err != nil {
		return diag.Errorf("error retrieving MRS clusters: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("clusters", flattenGetClustersResponseBodyClusters(cluster)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetClustersResponseBodyClusters(clusters []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(clusters))
	for _, v := range clusters {
		rst = append(rst, map[string]interface{}{
			"id":                       utils.PathSearch("clusterId", v, nil),
			"name":                     utils.PathSearch("clusterName", v, nil),
			"master_node_num":          utils.PathSearch("masterNodeNum", v, nil),
			"core_node_num":            utils.PathSearch("coreNodeNum", v, nil),
			"total_node_num":           utils.PathSearch("totalNodeNum", v, nil),
			"status":                   utils.PathSearch("clusterState", v, nil),
			"billing_type":             utils.PathSearch("billingType", v, nil),
			"vpc_id":                   utils.PathSearch("vpcId", v, nil),
			"subnet_id":                utils.PathSearch("subnetId", v, nil),
			"duration":                 utils.PathSearch("duration", v, nil),
			"fee":                      utils.PathSearch("fee", v, nil),
			"hadoop_version":           utils.PathSearch("hadoopVersion", v, nil),
			"master_node_size":         utils.PathSearch("masterNodeSize", v, nil),
			"core_node_size":           utils.PathSearch("coreNodeSize", v, nil),
			"component_list":           flattenClustersComponentList(v),
			"external_ip":              utils.PathSearch("externalIp", v, nil),
			"external_alternate_ip":    utils.PathSearch("externalAlternateIp", v, nil),
			"internal_ip":              utils.PathSearch("internalIp", v, nil),
			"deployment_id":            utils.PathSearch("deploymentId", v, nil),
			"description":              utils.PathSearch("remark", v, nil),
			"order_id":                 utils.PathSearch("orderId", v, nil),
			"master_node_product_id":   utils.PathSearch("masterNodeProductId", v, nil),
			"master_node_spec_id":      utils.PathSearch("masterNodeSpecId", v, nil),
			"core_node_product_id":     utils.PathSearch("coreNodeProductId", v, nil),
			"core_node_spec_id":        utils.PathSearch("coreNodeSpecId", v, nil),
			"availability_zone":        utils.PathSearch("availabilityZoneId", v, nil),
			"vnc":                      utils.PathSearch("vnc", v, nil),
			"volume_size":              utils.PathSearch("volumeSize", v, nil),
			"volume_type":              utils.PathSearch("volumeType", v, nil),
			"enterprise_project_id":    utils.PathSearch("enterpriseProjectId", v, nil),
			"type":                     utils.PathSearch("clusterType", v, nil),
			"security_group_id":        utils.PathSearch("securityGroupsId", v, nil),
			"slave_security_group_id":  utils.PathSearch("slaveSecurityGroupsId", v, nil),
			"stage_desc":               utils.PathSearch("stageDesc", v, nil),
			"safe_mode":                utils.PathSearch("safeMode", v, nil),
			"version":                  utils.PathSearch("clusterVersion", v, nil),
			"node_public_cert_name":    utils.PathSearch("nodePublicCertName", v, nil),
			"master_node_ip":           utils.PathSearch("masterNodeIp", v, nil),
			"private_ip_first":         utils.PathSearch("privateIpFirst", v, nil),
			"tags":                     flattenTags(utils.PathSearch("tags", v, "").(string)),
			"log_collection":           utils.PathSearch("logCollection", v, nil),
			"task_node_groups":         flattenClustersNodeGroups(v, "taskNodeGroups"),
			"node_groups":              flattenClustersNodeGroups(v, "nodeGroups"),
			"master_data_volume_type":  utils.PathSearch("masterDataVolumeType", v, nil),
			"master_data_volume_size":  utils.PathSearch("masterDataVolumeSize", v, nil),
			"master_data_volume_count": utils.PathSearch("masterDataVolumeCount", v, nil),
			"core_data_volume_type":    utils.PathSearch("coreDataVolumeType", v, nil),
			"core_data_volume_size":    utils.PathSearch("coreDataVolumeSize", v, nil),
			"core_data_volume_count":   utils.PathSearch("coreDataVolumeCount", v, nil),
			"period_type":              utils.PathSearch("periodType", v, nil),
			"scale":                    utils.PathSearch("scale", v, nil),
			"eip_id":                   utils.PathSearch("eipId", v, nil),
			"eip_address":              utils.PathSearch("eipAddress", v, nil),
			"eipv6_address":            utils.PathSearch("eipv6Address", v, nil),
			"mrs_ecs_default_agency":   utils.PathSearch("mrsEcsDefaultAgency", v, nil),
		})
	}
	return rst
}

func flattenClustersComponentList(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("componentList", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"component_id":      utils.PathSearch("componentId", v, nil),
			"component_name":    utils.PathSearch("componentName", v, nil),
			"component_version": utils.PathSearch("componentVersion", v, nil),
			"component_desc":    utils.PathSearch("componentDesc", v, nil),
		})
	}
	return rst
}

func flattenClustersNodeGroups(resp interface{}, key string) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch(key, resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"group_name":                     utils.PathSearch("groupName", v, nil),
			"node_num":                       utils.PathSearch("nodeNum", v, nil),
			"node_size":                      utils.PathSearch("nodeSize", v, nil),
			"node_spec_id":                   utils.PathSearch("nodeSpecId", v, nil),
			"node_product_id":                utils.PathSearch("nodeProductId", v, nil),
			"vm_product_id":                  utils.PathSearch("vmProductId", v, nil),
			"vm_spec_code":                   utils.PathSearch("vmSpecCode", v, nil),
			"root_volume_size":               utils.PathSearch("rootVolumeSize", v, nil),
			"root_volume_type":               utils.PathSearch("rootVolumeType", v, nil),
			"root_volume_product_id":         utils.PathSearch("rootVolumeProductId", v, nil),
			"root_volume_resource_spec_code": utils.PathSearch("rootVolumeResourceSpecCode", v, nil),
			"root_volume_resource_type":      utils.PathSearch("rootVolumeResourceType", v, nil),
			"data_volume_type":               utils.PathSearch("dataVolumeType", v, nil),
			"data_volume_count":              utils.PathSearch("dataVolumeCount", v, nil),
			"data_volume_size":               utils.PathSearch("dataVolumeSize", v, nil),
			"data_volume_product_id":         utils.PathSearch("dataVolumeProductId", v, nil),
			"data_volume_resource_spec_code": utils.PathSearch("dataVolumeResourceSpecCode", v, nil),
			"data_volume_resource_type":      utils.PathSearch("dataVolumeResourceType", v, nil),
		})
	}
	return rst
}
