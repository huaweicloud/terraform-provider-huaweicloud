package drs

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var drsJobV5NonUpdatableParams = []string{
	"base_info",
	"base_info.*.name",
	"base_info.*.job_type",
	"base_info.*.multi_write",
	"base_info.*.engine_type",
	"base_info.*.job_direction",
	"base_info.*.task_type",
	"base_info.*.net_type",
	"base_info.*.charging_mode",
	"base_info.*.enterprise_project_id",
	"base_info.*.description",
	"base_info.*.start_time",
	"base_info.*.expired_days",
	"base_info.*.tags",
	"base_info.*.tags.*.key",
	"base_info.*.tags.*.value",
	"base_info.*.is_open_fast_clean",
	"source_endpoint",
	"source_endpoint.*.db_type",
	"source_endpoint.*.endpoint_type",
	"source_endpoint.*.endpoint_role",
	"source_endpoint.*.endpoint",
	"source_endpoint.*.endpoint.*.id",
	"source_endpoint.*.endpoint.*.endpoint_name",
	"source_endpoint.*.endpoint.*.ip",
	"source_endpoint.*.endpoint.*.db_port",
	"source_endpoint.*.endpoint.*.db_user",
	"source_endpoint.*.endpoint.*.db_password",
	"source_endpoint.*.endpoint.*.instance_id",
	"source_endpoint.*.endpoint.*.instance_name",
	"source_endpoint.*.endpoint.*.db_name",
	"source_endpoint.*.endpoint.*.source_sharding",
	"source_endpoint.*.cloud",
	"source_endpoint.*.cloud.*.region",
	"source_endpoint.*.cloud.*.project_id",
	"source_endpoint.*.cloud.*.az_code",
	"source_endpoint.*.vpc",
	"source_endpoint.*.vpc.*.vpc_id",
	"source_endpoint.*.vpc.*.subnet_id",
	"source_endpoint.*.vpc.*.security_group_id",
	"source_endpoint.*.config",
	"source_endpoint.*.config.*.is_target_readonly",
	"source_endpoint.*.config.*.node_num",
	"source_endpoint.*.ssl",
	"source_endpoint.*.ssl.*.ssl_link",
	"source_endpoint.*.ssl.*.ssl_cert_name",
	"source_endpoint.*.ssl.*.ssl_cert_key",
	"source_endpoint.*.ssl.*.ssl_cert_check_sum",
	"source_endpoint.*.ssl.*.ssl_cert_password",
	"source_endpoint.*.customized_dns",
	"source_endpoint.*.customized_dns.*.is_set_dns",
	"source_endpoint.*.customized_dns.*.set_dns_action",
	"source_endpoint.*.customized_dns.*.dns_ip",
	"target_endpoint",
	"target_endpoint.*.db_type",
	"target_endpoint.*.endpoint_type",
	"target_endpoint.*.endpoint_role",
	"target_endpoint.*.endpoint",
	"target_endpoint.*.endpoint.*.id",
	"target_endpoint.*.endpoint.*.endpoint_name",
	"target_endpoint.*.endpoint.*.ip",
	"target_endpoint.*.endpoint.*.db_port",
	"target_endpoint.*.endpoint.*.db_user",
	"target_endpoint.*.endpoint.*.db_password",
	"target_endpoint.*.endpoint.*.instance_id",
	"target_endpoint.*.endpoint.*.instance_name",
	"target_endpoint.*.endpoint.*.db_name",
	"target_endpoint.*.endpoint.*.source_sharding",
	"target_endpoint.*.cloud",
	"target_endpoint.*.cloud.*.region",
	"target_endpoint.*.cloud.*.project_id",
	"target_endpoint.*.cloud.*.az_code",
	"target_endpoint.*.vpc",
	"target_endpoint.*.vpc.*.vpc_id",
	"target_endpoint.*.vpc.*.subnet_id",
	"target_endpoint.*.vpc.*.security_group_id",
	"target_endpoint.*.config",
	"target_endpoint.*.config.*.is_target_readonly",
	"target_endpoint.*.config.*.node_num",
	"target_endpoint.*.ssl",
	"target_endpoint.*.ssl.*.ssl_link",
	"target_endpoint.*.ssl.*.ssl_cert_name",
	"target_endpoint.*.ssl.*.ssl_cert_key",
	"target_endpoint.*.ssl.*.ssl_cert_check_sum",
	"target_endpoint.*.ssl.*.ssl_cert_password",
	"target_endpoint.*.customized_dns",
	"target_endpoint.*.customized_dns.*.is_set_dns",
	"target_endpoint.*.customized_dns.*.set_dns_action",
	"target_endpoint.*.customized_dns.*.dns_ip",
	"period_order",
	"period_order.*.period_type",
	"period_order.*.period_num",
	"period_order.*.is_auto_renew",
	"node_info",
	"node_info.*.spec",
	"node_info.*.spec.*.node_type",
	"node_info.*.vpc",
	"node_info.*.vpc.*.vpc_id",
	"node_info.*.vpc.*.subnet_id",
	"node_info.*.vpc.*.custom_node_ip",
	"node_info.*.vpc.*.security_group_id",
	"node_info.*.base_info",
	"node_info.*.base_info.*.instance_type",
	"node_info.*.base_info.*.arch",
	"node_info.*.base_info.*.availability_zone",
	"node_info.*.base_info.*.status",
	"node_info.*.base_info.*.role",
	"public_ip_list",
	"public_ip_list.*.id",
	"public_ip_list.*.public_ip",
	"public_ip_list.*.type",
}

// @API DRS POST /v5/{project_id}/jobs
// @API DRS GET /v5/{project_id}/jobs/{job_id}
// @API DRS DELETE /v5/{project_id}/jobs/{job_id}
func ResourceDrsJobV5() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDrsJobV5Create,
		ReadContext:   resourceDrsJobV5Read,
		UpdateContext: resourceDrsJobV5Update,
		DeleteContext: resourceDrsJobV5Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(drsJobV5NonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"base_info": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     drsJobV5BaseInfoSchema(),
			},
			"source_endpoint": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     drsJobV5EndpointSchema(),
			},
			"target_endpoint": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     drsJobV5EndpointSchema(),
			},
			"node_info": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     drsJobV5NodeInfoSchema(),
			},
			"period_order": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     drsJobV5PeriodOrderSchema(),
			},
			"public_ip_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     drsJobV5PublicIpListSchema(),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func drsJobV5BaseInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"job_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"multi_write": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"engine_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"job_direction": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"net_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expired_days": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     drsJobV5TagsSchema(),
			},
			"is_open_fast_clean": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func drsJobV5TagsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func drsJobV5EndpointSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint_role": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     drsJobV5EndpointDetailSchema(),
			},
			"cloud": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     drsJobV5CloudSchema(),
			},
			"vpc": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     drsJobV5EndpointVpcSchema(),
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     drsJobV5ConfigSchema(),
			},
			"ssl": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     drsJobV5SslSchema(),
			},
			"customized_dns": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     drsJobV5CustomizedDNSSchema(),
			},
		},
	}
}

func drsJobV5EndpointDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"endpoint_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_user": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This field `source_sharding` has a circular reference, so I changed it to a JSON string.
			"source_sharding": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func drsJobV5CloudSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"az_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func drsJobV5EndpointVpcSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func drsJobV5ConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"is_target_readonly": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func drsJobV5SslSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ssl_link": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ssl_cert_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_cert_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"ssl_cert_check_sum": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_cert_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func drsJobV5CustomizedDNSSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"is_set_dns": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"set_dns_action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dns_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func drsJobV5PeriodOrderSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"period_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"period_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"is_auto_renew": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func drsJobV5NodeInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"spec": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     drsJobV5NodeInfoSpecSchema(),
			},
			"vpc": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     drsJobV5NodeInfoVpcSchema(),
			},
			"base_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     drsJobV5NodeInfoBaseInfoSchema(),
			},
		},
	}
}

func drsJobV5NodeInfoSpecSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func drsJobV5NodeInfoVpcSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"custom_node_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func drsJobV5NodeInfoBaseInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"arch": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func drsJobV5PublicIpListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func buildDrsJobV5BodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"job": map[string]interface{}{
			"base_info":       buildJobBaseInfoParams(d.Get("base_info").([]interface{})),
			"source_endpoint": buildEndpointParams(d.Get("source_endpoint").([]interface{})),
			"target_endpoint": buildEndpointParams(d.Get("target_endpoint").([]interface{})),
			"period_order":    buildPeriodOrderParams(d.Get("period_order").([]interface{})),
			"node_info":       buildNodeInfoParams(d.Get("node_info").([]interface{})),
			"public_ip_list":  buildPublicIpListParams(d.Get("public_ip_list").([]interface{})),
		},
	}
	return bodyParams
}

func buildJobBaseInfoParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"name":                  utils.ValueIgnoreEmpty(rawMap["name"]),
		"job_type":              utils.ValueIgnoreEmpty(rawMap["job_type"]),
		"multi_write":           utils.ValueIgnoreEmpty(rawMap["multi_write"]),
		"engine_type":           utils.ValueIgnoreEmpty(rawMap["engine_type"]),
		"job_direction":         utils.ValueIgnoreEmpty(rawMap["job_direction"]),
		"task_type":             utils.ValueIgnoreEmpty(rawMap["task_type"]),
		"net_type":              utils.ValueIgnoreEmpty(rawMap["net_type"]),
		"charging_mode":         utils.ValueIgnoreEmpty(rawMap["charging_mode"]),
		"enterprise_project_id": utils.ValueIgnoreEmpty(rawMap["enterprise_project_id"]),
		"description":           utils.ValueIgnoreEmpty(rawMap["description"]),
		"start_time":            utils.ValueIgnoreEmpty(rawMap["start_time"]),
		"expired_days":          utils.ValueIgnoreEmpty(rawMap["expired_days"]),
		"tags":                  buildTagsParams(rawMap["tags"].([]interface{})),
		"is_open_fast_clean":    utils.ValueIgnoreEmpty(rawMap["is_open_fast_clean"]),
	}
}

func buildTagsParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"key":   utils.ValueIgnoreEmpty(rawMap["key"]),
			"value": utils.ValueIgnoreEmpty(rawMap["value"]),
		})
	}

	return rst
}

func buildEndpointParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"db_type":        rawMap["db_type"],
			"endpoint_type":  rawMap["endpoint_type"],
			"endpoint_role":  rawMap["endpoint_role"],
			"endpoint":       buildEndpointDetailParams(rawMap["endpoint"].([]interface{})),
			"cloud":          buildCloudParams(rawMap["cloud"].([]interface{})),
			"vpc":            buildEndpointVpcParams(rawMap["vpc"].([]interface{})),
			"config":         buildConfigParams(rawMap["config"].([]interface{})),
			"ssl":            buildSslParams(rawMap["ssl"].([]interface{})),
			"customized_dns": buildCustomizedDNSParams(rawMap["customized_dns"].([]interface{})),
		})
	}

	return rst
}

func buildEndpointDetailParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"endpoint_name":   rawMap["endpoint_name"],
		"db_user":         rawMap["db_user"],
		"db_password":     rawMap["db_password"],
		"id":              utils.ValueIgnoreEmpty(rawMap["id"]),
		"ip":              utils.ValueIgnoreEmpty(rawMap["ip"]),
		"db_port":         utils.ValueIgnoreEmpty(rawMap["db_port"]),
		"instance_id":     utils.ValueIgnoreEmpty(rawMap["instance_id"]),
		"instance_name":   utils.ValueIgnoreEmpty(rawMap["instance_name"]),
		"db_name":         utils.ValueIgnoreEmpty(rawMap["db_name"]),
		"source_sharding": utils.StringToJsonArray(rawMap["source_sharding"].(string)),
	}
}

func buildCloudParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"region":     rawMap["region"],
		"project_id": rawMap["project_id"],
		"az_code":    utils.ValueIgnoreEmpty(rawMap["az_code"]),
	}
}

func buildEndpointVpcParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"vpc_id":            rawMap["vpc_id"],
		"subnet_id":         rawMap["subnet_id"],
		"security_group_id": utils.ValueIgnoreEmpty(rawMap["security_group_id"]),
	}
}

func buildConfigParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"is_target_readonly": utils.ValueIgnoreEmpty(rawMap["is_target_readonly"]),
		"node_num":           utils.ValueIgnoreEmpty(rawMap["node_num"]),
	}
}

func buildSslParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"ssl_link":           utils.ValueIgnoreEmpty(rawMap["ssl_link"]),
		"ssl_cert_name":      utils.ValueIgnoreEmpty(rawMap["ssl_cert_name"]),
		"ssl_cert_key":       utils.ValueIgnoreEmpty(rawMap["ssl_cert_key"]),
		"ssl_cert_check_sum": utils.ValueIgnoreEmpty(rawMap["ssl_cert_check_sum"]),
		"ssl_cert_password":  utils.ValueIgnoreEmpty(rawMap["ssl_cert_password"]),
	}
}

func buildCustomizedDNSParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"is_set_dns":     rawMap["is_set_dns"],
		"set_dns_action": rawMap["set_dns_action"],
		"dns_ip":         rawMap["dns_ip"],
	}
}

func buildPeriodOrderParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"period_type":   rawMap["period_type"],
		"period_num":    rawMap["period_num"],
		"is_auto_renew": rawMap["is_auto_renew"],
	}
}

func buildNodeInfoParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"spec":      buildNodeInfoSpecParams(rawMap["spec"].([]interface{})),
		"vpc":       buildNodeInfoVpcParams(rawMap["vpc"].([]interface{})),
		"base_info": buildNodeInfoBaseInfoParams(rawMap["base_info"].([]interface{})),
	}
}

func buildNodeInfoSpecParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"node_type": rawMap["node_type"],
	}
}

func buildNodeInfoVpcParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"vpc_id":            rawMap["vpc_id"],
		"subnet_id":         rawMap["subnet_id"],
		"custom_node_ip":    utils.ValueIgnoreEmpty(rawMap["custom_node_ip"]),
		"security_group_id": utils.ValueIgnoreEmpty(rawMap["security_group_id"]),
	}
}

func buildNodeInfoBaseInfoParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"instance_type":     rawMap["instance_type"],
		"arch":              rawMap["arch"],
		"availability_zone": rawMap["availability_zone"],
		"status":            utils.ValueIgnoreEmpty(rawMap["status"]),
		"role":              utils.ValueIgnoreEmpty(rawMap["role"]),
	}
}

func buildPublicIpListParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"id":        rawMap["id"],
			"public_ip": rawMap["public_ip"],
			"type":      rawMap["type"],
		})
	}

	return rst
}

func GetJobV5Detail(client *golangsdk.ServiceClient, jobId string) (interface{}, error) {
	requestPath := client.Endpoint + "v5/{project_id}/jobs/{job_id}?type=detail"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func waitingForJobV5Success(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			jobDetail, err := GetJobV5Detail(client, d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("job.status", jobDetail, "").(string)
			if status == "" {
				return nil, "ERROR", errors.New("status is not found in job detail API response")
			}

			if status == "CONFIGURATION" {
				return jobDetail, "COMPLETED", nil
			}

			if status == "CREATING" {
				return jobDetail, "PENDING", nil
			}

			return jobDetail, status, nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceDrsJobV5Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildDrsJobV5BodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating DRS job: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the ID from the API response")
	}
	d.SetId(jobId)

	if err := waitingForJobV5Success(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for DRS job (%s) to success: %s", d.Id(), err)
	}

	return resourceDrsJobV5Read(ctx, d, meta)
}

// The response structure of the query API differs significantly from the request structure of the creation API,
// so integration will not be implemented at this time.
func resourceDrsJobV5Read(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

// The API editing capabilities have many limitations, so integration is not currently being implemented.
func resourceDrsJobV5Update(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrsJobV5Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs/{job_id}"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "DRS.10000010"),
			"error deleting DRS job",
		)
	}

	return nil
}
