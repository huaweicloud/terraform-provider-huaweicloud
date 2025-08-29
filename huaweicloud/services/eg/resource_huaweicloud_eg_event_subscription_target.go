package eg

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

var eventSubscriptionTargetNonUpdateParams = []string{"subscription_id", "name", "provider_type"}

// @API EG POST /v1/{project_id}/subscriptions/{subscription_id}/targets
// @API EG GET /v1/{project_id}/subscriptions/{subscription_id}/targets/{target_id}
// @API EG PUT /v1/{project_id}/subscriptions/{subscription_id}/targets/{target_id}
// @API EG DELETE /v1/{project_id}/subscriptions/{subscription_id}/targets/{target_id}
func ResourceEventSubscriptionTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventSubscriptionTargetCreate,
		ReadContext:   resourceEventSubscriptionTargetRead,
		UpdateContext: resourceEventSubscriptionTargetUpdate,
		DeleteContext: resourceEventSubscriptionTargetDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceEventSubscriptionTargetImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(eventSubscriptionTargetNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the event subscription target is located.",
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the event subscription.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the event subscription target.",
			},
			"provider_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The provider type of the event subscription target.",
			},
			"key_transform": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        keyTransformSchema(),
				Description: "The transform configuration for event data transformation.",
			},
			"connection_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The connection ID used by the event subscription target.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the enterprise project.",
			},
			"detail": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The configuration details of the event subscription target, in JSON format.",
			},
			"kafka_detail": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        kafkaDetailSchema(),
				Description: "The Kafka target configuration details.",
			},
			"smn_detail": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        smnDetailSchema(),
				Description: "The SMN target configuration details.",
			},
			"eg_detail": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        egDetailSchema(),
				Description: "The EG channel target configuration details.",
			},
			"apigw_detail": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        apigwDetailSchema(),
				Description: "The APIGW target configuration details.",
			},
			"retry_times": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "The number of retry times for the event subscription target.",
			},
			"dead_letter_queue": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        deadLetterQueueSchema(),
				Description: "The dead letter queue configuration of the event subscription target.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The create time, in UTC format.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time, in UTC format.",
			},
		},
	}
}

func keyTransformSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of transform rule.",
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The value of the transform rule.",
			},
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The template definition for VARIABLE type transform rules.",
			},
		},
	}
}

func kafkaDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"topic": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The topic name of the Kafka instance.",
			},
			"key_transform": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        keyTransformSchema(),
				Description: "The transform configuration of the Kafka messages.",
			},
		},
	}
}

func smnDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"urn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URN of the SMN topic.",
			},
			"agency_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The agency name for cross-account access.",
			},
			"key_transform": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        keyTransformSchema(),
				Description: "The subject transform configuration of the SMN messages.",
			},
		},
	}
}

func egDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"target_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The target project ID of the EG channel.",
			},
			"target_channel_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The target channel ID of the EG channel.",
			},
			"target_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The target region of the EG channel.",
			},
			"agency_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The agency name of cross-account access.",
			},
			"cross_region": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether this is a cross-region EG channel target.",
			},
			"cross_account": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether this is a cross-account EG channel target.",
			},
		},
	}
}

func apigwDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL of the APIGW endpoint.",
			},
			"invocation_http_parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        invocationHttpParametersSchema(),
				Description: "The HTTP parameters for the APIGW invocation.",
			},
		},
	}
}

func invocationHttpParametersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"header_parameters": {
				Type:        schema.TypeList,
				Elem:        headerParameterSchema(),
				Optional:    true,
				Description: "The header parameters for the HTTP request.",
			},
		},
	}
}

func headerParameterSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The key of the header parameter.",
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The value of the header parameter.",
			},
			"is_value_secret": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the header parameter value is secret.",
			},
		},
	}
}

func deadLetterQueueSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of dead letter queue.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance ID of the dead letter queue.",
			},
			"connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The connection ID of the dead letter queue.",
			},
			"topic": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The topic name of the dead letter queue.",
			},
		},
	}
}

func buildKeyTransform(keyTransformList []interface{}) map[string]interface{} {
	if len(keyTransformList) == 0 {
		return nil
	}

	return map[string]interface{}{
		"type":     utils.PathSearch("type", keyTransformList[0], nil),
		"value":    utils.ValueIgnoreEmpty(utils.PathSearch("value", keyTransformList[0], nil)),
		"template": utils.ValueIgnoreEmpty(utils.PathSearch("template", keyTransformList[0], nil)),
	}
}

func buildKafkaDetail(kafkaDetailList []interface{}) map[string]interface{} {
	if len(kafkaDetailList) == 0 {
		return nil
	}

	return map[string]interface{}{
		"topic":        utils.PathSearch("topic", kafkaDetailList[0], nil),
		"keyTransform": utils.ValueIgnoreEmpty(utils.PathSearch("key_transform", kafkaDetailList[0], nil)),
	}
}

func buildSmnDetail(smnDetailList []interface{}) map[string]interface{} {
	if len(smnDetailList) == 0 {
		return nil
	}

	return map[string]interface{}{
		"urn":               utils.PathSearch("urn", smnDetailList[0], nil),
		"agency_name":       utils.PathSearch("agency_name", smnDetailList[0], nil),
		"subject_transform": utils.ValueIgnoreEmpty(utils.PathSearch("key_transform", smnDetailList[0], nil)),
	}
}

func buildEgDetail(egDetailList []interface{}) map[string]interface{} {
	if len(egDetailList) == 0 {
		return nil
	}

	return map[string]interface{}{
		"target_project_id": utils.PathSearch("target_project_id", egDetailList[0], nil),
		"target_channel_id": utils.PathSearch("target_channel_id", egDetailList[0], nil),
		"target_region":     utils.PathSearch("target_region", egDetailList[0], nil),
		"agency_name":       utils.PathSearch("agency_name", egDetailList[0], nil),
		"cross_region":      utils.ValueIgnoreEmpty(utils.PathSearch("cross_region", egDetailList[0], nil)),
		"cross_account":     utils.ValueIgnoreEmpty(utils.PathSearch("cross_account", egDetailList[0], nil)),
	}
}

func buildHeaderParameters(headerParameters []interface{}) []interface{} {
	if len(headerParameters) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(headerParameters))
	for _, item := range headerParameters {
		headerParam := map[string]interface{}{
			"key":             utils.ValueIgnoreEmpty(utils.PathSearch("key", item, nil)),
			"value":           utils.ValueIgnoreEmpty(utils.PathSearch("value", item, nil)),
			"is_value_secret": utils.ValueIgnoreEmpty(utils.PathSearch("is_value_secret", item, nil)),
		}
		result = append(result, headerParam)
	}
	return result
}

func buildInvocationHttpParameters(invocationHttpParameters interface{}) map[string]interface{} {
	if invocationHttpParameters == nil {
		return nil
	}

	return map[string]interface{}{
		"header_parameters": buildHeaderParameters(
			utils.PathSearch("header_parameters", invocationHttpParameters, make([]interface{}, 0)).([]interface{})),
	}
}

func buildApigwDetail(apigwDetailList []interface{}) map[string]interface{} {
	if len(apigwDetailList) == 0 {
		return nil
	}

	return map[string]interface{}{
		"url": utils.PathSearch("url", apigwDetailList[0], nil),
		"invocation_http_parameters": utils.ValueIgnoreEmpty(
			buildInvocationHttpParameters(utils.PathSearch("invocation_http_parameters", apigwDetailList[0], nil))),
	}
}

func buildDeadLetterQueue(deadLetterQueueList []interface{}) map[string]interface{} {
	if len(deadLetterQueueList) == 0 {
		return nil
	}

	return map[string]interface{}{
		"type":          utils.ValueIgnoreEmpty(utils.PathSearch("type", deadLetterQueueList[0], nil)),
		"instance_id":   utils.ValueIgnoreEmpty(utils.PathSearch("instance_id", deadLetterQueueList[0], nil)),
		"connection_id": utils.ValueIgnoreEmpty(utils.PathSearch("connection_id", deadLetterQueueList[0], nil)),
		"topic":         utils.ValueIgnoreEmpty(utils.PathSearch("topic", deadLetterQueueList[0], nil)),
	}
}

func buildEventSubscriptionTargetBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":              d.Get("name"),
		"provider_type":     d.Get("provider_type"),
		"transform":         buildKeyTransform(d.Get("key_transform").([]interface{})),
		"connection_id":     utils.ValueIgnoreEmpty(d.Get("connection_id")),
		"detail":            utils.StringToJson(d.Get("detail").(string)),
		"kafka_detail":      utils.ValueIgnoreEmpty(buildKafkaDetail(d.Get("kafka_detail").([]interface{}))),
		"smn_detail":        utils.ValueIgnoreEmpty(buildSmnDetail(d.Get("smn_detail").([]interface{}))),
		"eg_detail":         utils.ValueIgnoreEmpty(buildEgDetail(d.Get("eg_detail").([]interface{}))),
		"apigw_detail":      utils.ValueIgnoreEmpty(buildApigwDetail(d.Get("apigw_detail").([]interface{}))),
		"retry_times":       utils.ValueIgnoreEmpty(d.Get("retry_times")),
		"dead_letter_queue": utils.ValueIgnoreEmpty(buildDeadLetterQueue(d.Get("dead_letter_queue").([]interface{}))),
	}
}

func buildEventSubscriptionTargetQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}
	return res
}

func resourceEventSubscriptionTargetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		httpUrl        = "v1/{project_id}/subscriptions/{subscription_id}/targets"
		region         = cfg.GetRegion(d)
		subscriptionId = d.Get("subscription_id").(string)
	)

	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{subscription_id}", subscriptionId)
	createPath += buildEventSubscriptionTargetQueryParams(d)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildEventSubscriptionTargetBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating EG event subscription target: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	targetId := utils.PathSearch("id", respBody, "").(string)
	if targetId == "" {
		return diag.Errorf("unable to find the target ID from the API response")
	}
	d.SetId(targetId)

	return resourceEventSubscriptionTargetRead(ctx, d, meta)
}

// GetEventSubscriptionTargetById is a method is used to get the event subscription target
func GetEventSubscriptionTargetById(client *golangsdk.ServiceClient, subscriptionId, targetId string) (interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/subscriptions/{subscription_id}/targets/{target_id}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{subscription_id}", subscriptionId)
	getPath = strings.ReplaceAll(getPath, "{target_id}", targetId)

	opt := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", getPath, opt)
	if err != nil {
		return nil, err
	}

	target, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return target, nil
}

func flattenKeyTransform(keyTransform interface{}) []interface{} {
	if keyTransform == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"type":     utils.PathSearch("type", keyTransform, nil),
			"value":    utils.ValueIgnoreEmpty(utils.PathSearch("value", keyTransform, nil)),
			"template": utils.ValueIgnoreEmpty(utils.PathSearch("template", keyTransform, nil)),
		},
	}
}

func flattenKafkaDetail(kafkaDetail interface{}) []interface{} {
	if kafkaDetail == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"topic":         utils.PathSearch("topic", kafkaDetail, nil),
			"key_transform": utils.ValueIgnoreEmpty(flattenKeyTransform(utils.PathSearch("keyTransform", kafkaDetail, nil))),
		},
	}
}

func flattenSmnDetail(smnDetail interface{}) []interface{} {
	if smnDetail == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"urn":           utils.PathSearch("urn", smnDetail, nil),
			"agency_name":   utils.PathSearch("agency_name", smnDetail, nil),
			"key_transform": utils.ValueIgnoreEmpty(flattenKeyTransform(utils.PathSearch("subject_transform", smnDetail, nil))),
		},
	}
}

func flattenEgDetail(egDetail interface{}) []interface{} {
	if egDetail == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"target_project_id": utils.PathSearch("target_project_id", egDetail, nil),
			"target_channel_id": utils.PathSearch("target_channel_id", egDetail, nil),
			"target_region":     utils.PathSearch("target_region", egDetail, nil),
			"agency_name":       utils.PathSearch("agency_name", egDetail, nil),
			"cross_region":      utils.ValueIgnoreEmpty(utils.PathSearch("cross_region", egDetail, nil)),
			"cross_account":     utils.ValueIgnoreEmpty(utils.PathSearch("cross_account", egDetail, nil)),
		},
	}
}

func flattenHeaderParameters(parameters []interface{}) []interface{} {
	if parameters == nil {
		return nil
	}

	result := make([]interface{}, 0, len(parameters))
	for _, item := range parameters {
		result = append(result, map[string]interface{}{
			"key":             utils.ValueIgnoreEmpty(utils.PathSearch("key", item, nil)),
			"value":           utils.ValueIgnoreEmpty(utils.PathSearch("value", item, nil)),
			"is_value_secret": utils.ValueIgnoreEmpty(utils.PathSearch("is_value_secret", item, nil)),
		})
	}
	return result
}

func flattenInvocationHttpParameters(parametersObj interface{}) interface{} {
	if parametersObj == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"header_parameters": utils.ValueIgnoreEmpty(
				flattenHeaderParameters(utils.PathSearch("header_parameters", parametersObj, make([]interface{}, 0)).([]interface{}))),
		},
	}
}

func flattenApigwDetail(apigwDetail interface{}) []interface{} {
	if apigwDetail == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"url": utils.PathSearch("url", apigwDetail, nil),
			"invocation_http_parameters": utils.ValueIgnoreEmpty(
				flattenInvocationHttpParameters(utils.PathSearch("invocation_http_parameters", apigwDetail, nil))),
		},
	}
}

func flattenDeadLetterQueue(deadLetterQueue interface{}) []interface{} {
	if deadLetterQueue == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"type":          utils.ValueIgnoreEmpty(utils.PathSearch("type", deadLetterQueue, nil)),
			"instance_id":   utils.ValueIgnoreEmpty(utils.PathSearch("instance_id", deadLetterQueue, nil)),
			"connection_id": utils.ValueIgnoreEmpty(utils.PathSearch("connection_id", deadLetterQueue, nil)),
			"topic":         utils.ValueIgnoreEmpty(utils.PathSearch("topic", deadLetterQueue, nil)),
		},
	}
}

func resourceEventSubscriptionTargetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		subscriptionId = d.Get("subscription_id").(string)
		targetId       = d.Id()
	)

	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	target, err := GetEventSubscriptionTargetById(client, subscriptionId, targetId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving event subscription target")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", target, nil)),
		d.Set("provider_type", utils.PathSearch("provider_type", target, nil)),
		d.Set("key_transform", flattenKeyTransform(utils.PathSearch("transform", target, nil))),
		d.Set("connection_id", utils.ValueIgnoreEmpty(utils.PathSearch("connection_id", target, nil))),
		d.Set("detail", utils.ValueIgnoreEmpty(utils.JsonToString(utils.PathSearch("detail", target, nil)))),
		d.Set("kafka_detail", utils.ValueIgnoreEmpty(flattenKafkaDetail(utils.PathSearch("kafka_detail", target, nil)))),
		d.Set("smn_detail", utils.ValueIgnoreEmpty(flattenSmnDetail(utils.PathSearch("smn_detail", target, nil)))),
		d.Set("eg_detail", utils.ValueIgnoreEmpty(flattenEgDetail(utils.PathSearch("eg_detail", target, nil)))),
		d.Set("apigw_detail", utils.ValueIgnoreEmpty(flattenApigwDetail(utils.PathSearch("apigw_detail", target, nil)))),
		d.Set("retry_times", utils.ValueIgnoreEmpty(utils.PathSearch("retry_times", target, nil))),
		d.Set("dead_letter_queue", utils.ValueIgnoreEmpty(flattenDeadLetterQueue(utils.PathSearch("dead_letter_queue", target, nil)))),
		d.Set("created_at", utils.PathSearch("created_time", target, nil)),
		d.Set("updated_at", utils.PathSearch("updated_time", target, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEventSubscriptionTargetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		subscriptionId = d.Get("subscription_id").(string)
		targetId       = d.Id()
		httpUrl        = "v1/{project_id}/subscriptions/{subscription_id}/targets/{target_id}"
	)

	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{subscription_id}", subscriptionId)
	updatePath = strings.ReplaceAll(updatePath, "{target_id}", targetId)
	updatePath += buildEventSubscriptionTargetQueryParams(d)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildEventSubscriptionTargetBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating EG event subscription target: %s", err)
	}

	return resourceEventSubscriptionTargetRead(ctx, d, meta)
}

func resourceEventSubscriptionTargetDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		httpUrl        = "v1/{project_id}/subscriptions/{subscription_id}/targets/{target_id}"
		subscriptionId = d.Get("subscription_id").(string)
		targetId       = d.Id()
	)

	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{subscription_id}", subscriptionId)
	deletePath = strings.ReplaceAll(deletePath, "{target_id}", targetId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting EG event subscription target")
	}

	return nil
}

func resourceEventSubscriptionTargetImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<subscription_id>/<id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("subscription_id", parts[0])
}
