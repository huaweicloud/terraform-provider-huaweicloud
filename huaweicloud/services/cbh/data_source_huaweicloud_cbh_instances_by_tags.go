package cbh

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

// @API CBH POST /v2/{project_id}/cbs/instance/filter
func DataSourceCbhInstancesByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceCbhInstancesByTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"without_any_tag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to query all resources without tags.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the tags to be included in the query.`,
				Elem:        instanceTagsRequestBodySchema(),
			},
			"tags_any": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the tags to be included in the query.`,
				Elem:        instanceTagsRequestBodySchema(),
			},
			"not_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the tags that are not included in the query.`,
				Elem:        instanceTagsRequestBodySchema(),
			},
			"not_tags_any": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the tags that are not included in the query.`,
				Elem:        instanceTagsRequestBodySchema(),
			},
			"sys_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the system tags to be included in the query.`,
				Elem:        instanceTagsRequestBodySchema(),
			},
			"matches": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the search fields, key is the field to be matched.`,
				Elem:        instanceMatchesRequestBodySchema(),
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of CBH instance.`,
				Elem:        instanceResourcesComputedSchema(),
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total number of records.`,
			},
		},
	}
}

func instanceTagsRequestBodySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the key of the tag.`,
			},
			"values": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the value list of the tag.`,
			},
		},
	}
}

func instanceMatchesRequestBodySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the key of the tag.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the value of the tag.`,
			},
		},
	}
}

func instanceResourcesDetailComputedSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the resource.`,
			},
			"server_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the server.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the instance.`,
			},
			"alter_permit": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `The instance whether can be expanded.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the enterprise project.`,
			},
			"period_num": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance period number.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance start time in timestamp format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance end time in timestamp format.`,
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance created time in UTC time format.`,
			},
			"upgrade_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The instance upgrade time in timestamp format.`,
			},
			"update": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance whether can be expanded.`,
			},
			"bastion_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance current version.`,
			},
			"az_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of availability zone information.`,
				Elem:        resourcesResourceDetailAzInfoElem(),
			},
			"status_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of status information.`,
				Elem:        resResDetStaInfoElem(),
			},
			"resource_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of resource information.`,
				Elem:        resResDetResInfoElem(),
			},
			"network": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of network information.`,
				Elem:        resourcesResourceDetailNetworkElem(),
			},
			"ha_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of high availability information.`,
				Elem:        resourcesResourceDetailHaInfoElem(),
			},
		},
	}
}

func instanceResourcesTagsComputedSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The key of the tag.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the tag.`,
			},
		},
	}
}

func instanceResourcesComputedSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the instance.`,
			},
			"resource_detail": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of resource details.`,
				Elem:        instanceResourcesDetailComputedSchema(),
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The tags of the resource.`,
				Elem:        instanceResourcesTagsComputedSchema(),
			},
			"sys_tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The system tags of the resource.`,
				Elem:        instanceResourcesTagsComputedSchema(),
			},
			"resource_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the resource.`,
			},
		},
	}
}

func resResDetResInfoElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"specification": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The specification of the instance.`,
			},
			"order_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The order ID of the instance.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the instance.`,
			},
			"data_disk_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The size of the data disk of the instance.`,
			},
			"disk_resource_id": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The ID of the data disk of the instance.`,
			},
		},
	}
}

func resResDetStaInfoElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the instance.`,
			},
			"task_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The task status of the instance.`,
			},
			"create_instance_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the instance in the process of creating an instance.`,
			},
			"instance_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the instance.`,
			},
			"instance_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the instance.`,
			},
			"fail_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The failure reason of the instance.`,
			},
		},
	}
}

func resourcesResourceDetailAzInfoElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The region where the instance is located.`,
			},
			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The availability zone ID of the instance.`,
			},
			"availability_zone_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The availability zone name of the instance.`,
			},
			"slave_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The availability zone ID of the slave instance.`,
			},
			"slave_zone_display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The availability zone name of the slave instance.`,
			},
		},
	}
}

func resourcesResourceDetailNetworkElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The floating IP of the instance.`,
			},
			"web_port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The port of the web interface of the instance.`,
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The elastic public IP of the instance.`,
			},
			"public_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the elastic public IP of the instance.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The private IP of the instance.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the VPC where the instance is located.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the subnet where the instance is located.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the security group to which the instance belongs.`,
			},
		},
	}
}

func resourcesResourceDetailHaInfoElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ha_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the primary-backup mode.`,
			},
			"instance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the instance.`,
			},
		},
	}
}

func buildInstanceQueryTagsBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	respArray := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		respArray = append(respArray, map[string]interface{}{
			"key":    rawMap["key"],
			"values": rawMap["values"],
		})
	}
	return respArray
}

func buildInstanceQueryMatchesBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	respArray := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		respArray = append(respArray, map[string]interface{}{
			"key":   rawMap["key"],
			"value": rawMap["value"],
		})
	}
	return respArray
}

func buildInstanceQueryBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"without_any_tag": d.Get("without_any_tag"),
		"tags":            buildInstanceQueryTagsBodyParams(d.Get("tags").([]interface{})),
		"tags_any":        buildInstanceQueryTagsBodyParams(d.Get("tags_any").([]interface{})),
		"not_tags":        buildInstanceQueryTagsBodyParams(d.Get("not_tags").([]interface{})),
		"not_tags_any":    buildInstanceQueryTagsBodyParams(d.Get("not_tags_any").([]interface{})),
		"sys_tags":        buildInstanceQueryTagsBodyParams(d.Get("sys_tags").([]interface{})),
		"matches":         buildInstanceQueryMatchesBodyParams(d.Get("matches").([]interface{})),
	}
}

// Invalid pagination parameters will prevent pagination from being used when querying data.
func resourceCbhInstancesByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cbh"
		httpUrl = "v2/{project_id}/cbs/instance/filter"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildInstanceQueryBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CBH instances by tags: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("resources", flattenResourcesAttributes(utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("total_count", utils.PathSearch("total_count", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenResourcesAttributes(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	restArray := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		resourceDetail := utils.PathSearch("resource_detail", v, nil)
		tags := utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})
		sysTags := utils.PathSearch("sys_tags", v, make([]interface{}, 0)).([]interface{})
		restArray = append(restArray, map[string]interface{}{
			"resource_id":     utils.PathSearch("resource_id", v, nil),
			"resource_detail": flattenResourcesDetailAttributes(resourceDetail),
			"tags":            flattenResourcesTagsAttributes(tags),
			"sys_tags":        flattenResourcesTagsAttributes(sysTags),
			"resource_name":   utils.PathSearch("resource_name", v, nil),
		})
	}

	return restArray
}

func flattenResourcesDetailAttributes(rawMap interface{}) []map[string]interface{} {
	if rawMap == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":                  utils.PathSearch("name", rawMap, nil),
			"server_id":             utils.PathSearch("server_id", rawMap, nil),
			"instance_id":           utils.PathSearch("instance_id", rawMap, nil),
			"alter_permit":          utils.PathSearch("alter_permit", rawMap, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", rawMap, nil),
			"period_num":            utils.PathSearch("period_num", rawMap, nil),
			"start_time":            utils.PathSearch("start_time", rawMap, nil),
			"end_time":              utils.PathSearch("end_time", rawMap, nil),
			"created_time":          utils.PathSearch("created_time", rawMap, nil),
			"upgrade_time":          utils.PathSearch("upgrade_time", rawMap, nil),
			"update":                utils.PathSearch("update", rawMap, nil),
			"bastion_version":       utils.PathSearch("bastion_version", rawMap, nil),
			"az_info":               flattenDetailAzInfoAttributes(utils.PathSearch("az_info", rawMap, nil)),
			"status_info":           flattenDetailStatusInfoAttributes(utils.PathSearch("status_info", rawMap, nil)),
			"resource_info":         flattenDetailResourceInfoAttributes(utils.PathSearch("resource_info", rawMap, nil)),
			"network":               flattenDetailNetworkAttributes(utils.PathSearch("network", rawMap, nil)),
			"ha_info":               flattenDetailHaInfoAttributes(utils.PathSearch("ha_info", rawMap, nil)),
		},
	}
}

func flattenDetailHaInfoAttributes(rawMap interface{}) []map[string]interface{} {
	if rawMap == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"ha_id":         utils.PathSearch("ha_id", rawMap, nil),
			"instance_type": utils.PathSearch("instance_type", rawMap, nil),
		},
	}
}

func flattenDetailNetworkAttributes(rawMap interface{}) []map[string]interface{} {
	if rawMap == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"vip":               utils.PathSearch("vip", rawMap, nil),
			"web_port":          utils.PathSearch("web_port", rawMap, nil),
			"public_ip":         utils.PathSearch("public_ip", rawMap, nil),
			"public_id":         utils.PathSearch("public_id", rawMap, nil),
			"private_ip":        utils.PathSearch("private_ip", rawMap, nil),
			"vpc_id":            utils.PathSearch("vpc_id", rawMap, nil),
			"subnet_id":         utils.PathSearch("subnet_id", rawMap, nil),
			"security_group_id": utils.PathSearch("security_group_id", rawMap, nil),
		},
	}
}

func flattenDetailResourceInfoAttributes(rawMap interface{}) []map[string]interface{} {
	if rawMap == nil {
		return nil
	}

	diskResourceID := utils.PathSearch("disk_resource_id", rawMap, make([]interface{}, 0)).([]interface{})
	return []map[string]interface{}{
		{
			"specification":    utils.PathSearch("specification", rawMap, nil),
			"order_id":         utils.PathSearch("order_id", rawMap, nil),
			"resource_id":      utils.PathSearch("resource_id", rawMap, nil),
			"data_disk_size":   utils.PathSearch("data_disk_size", rawMap, nil),
			"disk_resource_id": utils.ExpandToStringList(diskResourceID),
		},
	}
}

func flattenDetailStatusInfoAttributes(rawMap interface{}) []map[string]interface{} {
	if rawMap == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"status":                 utils.PathSearch("status", rawMap, nil),
			"task_status":            utils.PathSearch("task_status", rawMap, nil),
			"create_instance_status": utils.PathSearch("create_instance_status", rawMap, nil),
			"instance_status":        utils.PathSearch("instance_status", rawMap, nil),
			"instance_description":   utils.PathSearch("instance_description", rawMap, nil),
			"fail_reason":            utils.PathSearch("fail_reason", rawMap, nil),
		},
	}
}

func flattenDetailAzInfoAttributes(rawMap interface{}) []map[string]interface{} {
	if rawMap == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"region":                    utils.PathSearch("region", rawMap, nil),
			"zone":                      utils.PathSearch("zone", rawMap, nil),
			"availability_zone_display": utils.PathSearch("availability_zone_display", rawMap, nil),
			"slave_zone":                utils.PathSearch("slave_zone", rawMap, nil),
			"slave_zone_display":        utils.PathSearch("slave_zone_display", rawMap, nil),
		},
	}
}

func flattenResourcesTagsAttributes(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	restArray := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		restArray = append(restArray, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return restArray
}
