package drs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS POST /v5/{project_id}/connections
// @API DRS GET /v5/{project_id}/connections
// @API DRS PUT /v5/{project_id}/connections/{connection_id}
// @API DRS DELETE /v5/{project_id}/connections/{connection_id}
func ResourceDrsConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionCreate,
		ReadContext:   resourceConnectionRead,
		UpdateContext: resourceConnectionUpdate,
		DeleteContext: resourceConnectionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     connectionEndpointSchema(),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     resourceConnectionConfigSchema(),
			},
			"vpc": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     connectionVpcSchema(),
			},
			"ssl": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     connectionSslSchema(),
			},
			"cloud": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     connectionCloudSchema(),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func connectionEndpointSchema() *schema.Resource {
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
				Computed: true,
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
			"source_sharding": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     connectionSourceShardingSchema(),
			},
		},
	}
}

func connectionSourceShardingSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"endpoint_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_port": {
				Type:     schema.TypeString,
				Optional: true,
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
		},
	}
}

func resourceConnectionConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"driver_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func connectionVpcSchema() *schema.Resource {
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

func connectionSslSchema() *schema.Resource {
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

func connectionCloudSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"az_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildConnectionEndpoint(d *schema.ResourceData) map[string]interface{} {
	endpointRaw := d.Get("endpoint").([]interface{})
	if len(endpointRaw) == 0 {
		return nil
	}

	endpoint, ok := endpointRaw[0].(map[string]interface{})
	if !ok {
		return nil
	}

	result := map[string]interface{}{
		"endpoint_name": endpoint["endpoint_name"],
		"db_user":       endpoint["db_user"],
		"db_password":   endpoint["db_password"],
		"id":            utils.ValueIgnoreEmpty(endpoint["id"]),
		"ip":            utils.ValueIgnoreEmpty(endpoint["ip"]),
		"db_port":       utils.ValueIgnoreEmpty(endpoint["db_port"]),
		"instance_id":   utils.ValueIgnoreEmpty(endpoint["instance_id"]),
		"instance_name": utils.ValueIgnoreEmpty(endpoint["instance_name"]),
		"db_name":       utils.ValueIgnoreEmpty(endpoint["db_name"]),
	}

	if sourceSharding, ok := endpoint["source_sharding"].([]interface{}); ok && len(sourceSharding) > 0 {
		result["source_sharding"] = buildSourceSharding(sourceSharding)
	}

	return utils.RemoveNil(result)
}

func buildSourceSharding(sourceSharding []interface{}) []map[string]interface{} {
	shardings := make([]map[string]interface{}, 0, len(sourceSharding))
	for _, s := range sourceSharding {
		sharding, ok := s.(map[string]interface{})
		if !ok {
			continue
		}

		shardingMap := map[string]interface{}{
			"endpoint_name": sharding["endpoint_name"],
			"db_user":       sharding["db_user"],
			"db_password":   sharding["db_password"],
			"id":            utils.ValueIgnoreEmpty(sharding["id"]),
			"ip":            utils.ValueIgnoreEmpty(sharding["ip"]),
			"db_port":       utils.ValueIgnoreEmpty(sharding["db_port"]),
			"instance_id":   utils.ValueIgnoreEmpty(sharding["instance_id"]),
			"instance_name": utils.ValueIgnoreEmpty(sharding["instance_name"]),
			"db_name":       utils.ValueIgnoreEmpty(sharding["db_name"]),
		}
		shardings = append(shardings, utils.RemoveNil(shardingMap))
	}
	return shardings
}

func buildConnectionConfig(d *schema.ResourceData) map[string]interface{} {
	configRaw := d.Get("config").([]interface{})
	if len(configRaw) == 0 {
		return nil
	}

	cfg, ok := configRaw[0].(map[string]interface{})
	if !ok {
		return nil
	}

	result := map[string]interface{}{
		"driver_name": cfg["driver_name"],
	}
	return utils.RemoveNil(result)
}

func buildConnectionVpc(d *schema.ResourceData) map[string]interface{} {
	vpcRaw := d.Get("vpc").([]interface{})
	if len(vpcRaw) == 0 {
		return nil
	}

	vpc, ok := vpcRaw[0].(map[string]interface{})
	if !ok {
		return nil
	}

	result := map[string]interface{}{
		"vpc_id":            vpc["vpc_id"],
		"subnet_id":         vpc["subnet_id"],
		"security_group_id": vpc["security_group_id"],
	}
	return utils.RemoveNil(result)
}

func buildConnectionSsl(d *schema.ResourceData) map[string]interface{} {
	sslRaw := d.Get("ssl").([]interface{})
	if len(sslRaw) == 0 {
		return nil
	}

	ssl, ok := sslRaw[0].(map[string]interface{})
	if !ok {
		return nil
	}

	result := map[string]interface{}{
		"ssl_link":           ssl["ssl_link"],
		"ssl_cert_name":      utils.ValueIgnoreEmpty(ssl["ssl_cert_name"]),
		"ssl_cert_key":       utils.ValueIgnoreEmpty(ssl["ssl_cert_key"]),
		"ssl_cert_check_sum": utils.ValueIgnoreEmpty(ssl["ssl_cert_check_sum"]),
		"ssl_cert_password":  utils.ValueIgnoreEmpty(ssl["ssl_cert_password"]),
	}
	return utils.RemoveNil(result)
}

func buildConnectionCloud(d *schema.ResourceData) map[string]interface{} {
	cloudRaw := d.Get("cloud").([]interface{})
	if len(cloudRaw) == 0 {
		return nil
	}

	cloud, ok := cloudRaw[0].(map[string]interface{})
	if !ok {
		return nil
	}

	result := map[string]interface{}{
		"region":     utils.ValueIgnoreEmpty(cloud["region"]),
		"project_id": utils.ValueIgnoreEmpty(cloud["project_id"]),
		"az_code":    utils.ValueIgnoreEmpty(cloud["az_code"]),
	}
	return utils.RemoveNil(result)
}

func buildCreateConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                  d.Get("name"),
		"db_type":               d.Get("db_type"),
		"endpoint":              buildConnectionEndpoint(d),
		"description":           d.Get("description"),
		"vpc":                   buildConnectionVpc(d),
		"ssl":                   buildConnectionSsl(d),
		"config":                buildConnectionConfig(d),
		"cloud":                 buildConnectionCloud(d),
		"enterprise_project_id": utils.ValueIgnoreEmpty(d.Get("enterprise_project_id")),
	}

	return utils.RemoveNil(bodyParams)
}

func buildUpdateConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"db_type":     d.Get("db_type"),
		"endpoint":    buildConnectionEndpoint(d),
		"description": d.Get("description"),
		"vpc":         buildConnectionVpc(d),
		"ssl":         buildConnectionSsl(d),
		"config":      buildConnectionConfig(d),
		"cloud":       buildConnectionCloud(d),
	}

	return utils.RemoveNil(bodyParams)
}

func GetConnectionById(client *golangsdk.ServiceClient, connectionId string) (interface{}, error) {
	requestPath := client.Endpoint + "v5/{project_id}/connections"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = fmt.Sprintf("%s?connection_id=%s", requestPath, connectionId)

	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	connection := utils.PathSearch("connections|[0]", respBody, nil)
	if connection == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return connection, nil
}

func resourceConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/connections"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateConnectionBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating DRS connection: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	connectionId := utils.PathSearch("connection_id", respBody, "").(string)
	if connectionId == "" {
		return diag.Errorf("error creating DRS connection: connection_id is not found in response")
	}

	d.SetId(connectionId)

	return resourceConnectionRead(ctx, d, meta)
}

func resourceConnectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	connection, err := GetConnectionById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DRS connection")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", connection, nil)),
		d.Set("db_type", utils.PathSearch("db_type", connection, nil)),
		d.Set("endpoint", flattenConnectionEndpoint(d, utils.PathSearch("endpoint", connection, nil))),
		d.Set("description", utils.PathSearch("description", connection, nil)),
		d.Set("config", flattenConnectionConfig(utils.PathSearch("config", connection, nil))),
		d.Set("vpc", flattenConnectionVpc(utils.PathSearch("vpc", connection, nil))),
		d.Set("ssl", flattenConnectionSsl(utils.PathSearch("ssl", connection, nil))),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", connection, nil)),
		d.Set("create_time", utils.PathSearch("create_time", connection, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenConnectionEndpoint(d *schema.ResourceData, endpoint interface{}) []interface{} {
	if endpoint == nil {
		return nil
	}

	result := map[string]interface{}{
		"endpoint_name": utils.PathSearch("endpoint_name", endpoint, nil),
		"ip":            utils.ValueIgnoreEmpty(utils.PathSearch("ip", endpoint, nil)),
		"db_port":       utils.PathSearch("db_port", endpoint, nil),
		"db_user":       utils.PathSearch("db_user", endpoint, nil),
		"instance_id":   utils.ValueIgnoreEmpty(utils.PathSearch("instance_id", endpoint, nil)),
		"instance_name": utils.ValueIgnoreEmpty(utils.PathSearch("instance_name", endpoint, nil)),
		"db_name":       utils.ValueIgnoreEmpty(utils.PathSearch("db_name", endpoint, nil)),
		"id":            utils.ValueIgnoreEmpty(utils.PathSearch("id", endpoint, nil)),
	}

	if sourceSharding := utils.PathSearch("source_sharding", endpoint, nil); sourceSharding != nil {
		configSourceSharding := getConfigSourceSharding(d)
		shardings := flattenSourceShardings(sourceSharding.([]interface{}), configSourceSharding)
		result["source_sharding"] = shardings
	}

	return []interface{}{result}
}

func getConfigSourceSharding(d *schema.ResourceData) []interface{} {
	if endpoints, ok := d.GetOk("endpoint"); ok {
		endpointMap, ok := endpoints.([]interface{})[0].(map[string]interface{})
		if !ok {
			return make([]interface{}, 0)
		}
		if shard, ok := endpointMap["source_sharding"].([]interface{}); ok {
			return shard
		}
	}
	return make([]interface{}, 0)
}

func flattenSourceShardings(apiShardings []interface{}, configShardings []interface{}) []interface{} {
	if len(apiShardings) == 0 {
		return nil
	}

	shardings := make([]interface{}, 0, len(apiShardings))

	for i, s := range apiShardings {
		shard, ok := s.(map[string]interface{})
		if !ok {
			continue
		}

		shardingMap := flattenShardingBasicFields(shard)
		shardingMap["endpoint_name"] = getShardingEndpointName(shard, configShardings, i)
		shardings = append(shardings, shardingMap)
	}

	return shardings
}

func flattenShardingBasicFields(shard map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":            utils.ValueIgnoreEmpty(utils.PathSearch("id", shard, nil)),
		"ip":            utils.ValueIgnoreEmpty(utils.PathSearch("ip", shard, nil)),
		"db_port":       utils.ValueIgnoreEmpty(utils.PathSearch("db_port", shard, nil)),
		"db_user":       utils.PathSearch("db_user", shard, nil),
		"instance_id":   utils.ValueIgnoreEmpty(utils.PathSearch("instance_id", shard, nil)),
		"instance_name": utils.ValueIgnoreEmpty(utils.PathSearch("instance_name", shard, nil)),
		"db_name":       utils.ValueIgnoreEmpty(utils.PathSearch("db_name", shard, nil)),
	}
}

func getShardingEndpointName(shard map[string]interface{}, configShardings []interface{}, index int) interface{} {
	apiEndpointName := utils.PathSearch("endpoint_name", shard, nil)
	if apiEndpointName != nil && apiEndpointName != "" {
		return apiEndpointName
	}

	if index < len(configShardings) {
		if configShard, ok := configShardings[index].(map[string]interface{}); ok {
			if configEndpointName, ok := configShard["endpoint_name"].(string); ok && configEndpointName != "" {
				return configEndpointName
			}
		}
	}

	return nil
}

func flattenConnectionConfig(cfg interface{}) []interface{} {
	if cfg == nil {
		return nil
	}

	driverName := utils.PathSearch("driver_name", cfg, nil)
	if driverName == nil || driverName == "" {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"driver_name": driverName,
		},
	}
}

func flattenConnectionVpc(vpc interface{}) []interface{} {
	if vpc == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"vpc_id":            utils.PathSearch("vpc_id", vpc, nil),
			"subnet_id":         utils.PathSearch("subnet_id", vpc, nil),
			"security_group_id": utils.PathSearch("security_group_id", vpc, nil),
		},
	}
}

func flattenConnectionSsl(ssl interface{}) []interface{} {
	if ssl == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"ssl_link":           utils.PathSearch("ssl_link", ssl, nil),
			"ssl_cert_name":      utils.PathSearch("ssl_cert_name", ssl, nil),
			"ssl_cert_check_sum": utils.PathSearch("ssl_cert_check_sum", ssl, nil),
		},
	}
}

func resourceConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/connections/{connection_id}"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{connection_id}", d.Id())

	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateConnectionBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating DRS connection: %s", err)
	}

	return resourceConnectionRead(ctx, d, meta)
}

func resourceConnectionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/connections/{connection_id}"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{connection_id}", d.Id())

	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting DRS connection: %s", err)
	}

	return nil
}
