package iotda

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA GET /v5/iot/{project_id}/routing-rule/actions
// @API IoTDA POST /v5/iot/{project_id}/routing-rule/actions
// @API IoTDA DELETE /v5/iot/{project_id}/routing-rule/rules/{rule_id}
// @API IoTDA GET /v5/iot/{project_id}/routing-rule/rules/{rule_id}
// @API IoTDA PUT /v5/iot/{project_id}/routing-rule/rules/{rule_id}
// @API IoTDA POST /v5/iot/{project_id}/routing-rule/rules
// @API IoTDA DELETE /v5/iot/{project_id}/routing-rule/actions/{action_id}
// @API IoTDA PUT /v5/iot/{project_id}/routing-rule/actions/{action_id}
func ResourceDataForwardingRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceDataForwardingRuleCreate,
		UpdateContext: ResourceDataForwardingRuleUpdate,
		DeleteContext: ResourceDataForwardingRuleDelete,
		ReadContext:   ResourceDataForwardingRuleRead,
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

			"trigger": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"device:create",
					"device:delete",
					"device:update",
					"device.status:update",
					"device.property:report",
					"device.message:report",
					"device.message.status:update",
					"batchtask:update",
					"product:create",
					"product:delete",
					"product:update",
					"device.command.status:update",
				}, false),
			},

			"enabled": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"targets"},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"where": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"targets": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 10,
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})

					if m["id"] != nil {
						buf.WriteString(m["id"].(string))
					}

					return hashcode.String(buf.String())
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},

						"http_forwarding": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},

						"dis_forwarding": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:     schema.TypeString,
										Required: true,
									},

									"stream_id": {
										Type:     schema.TypeString,
										Required: true,
									},

									"project_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},

						"obs_forwarding": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:     schema.TypeString,
										Required: true,
									},

									"bucket": {
										Type:     schema.TypeString,
										Required: true,
									},

									"project_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"custom_directory": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},

						"amqp_forwarding": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"queue_name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},

						"kafka_forwarding": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:     schema.TypeString,
										Required: true,
									},

									"addresses": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:     schema.TypeInt,
													Required: true,
												},

												"ip": {
													Type:     schema.TypeString,
													Optional: true,
												},

												"domain": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},

									"topic": {
										Type:     schema.TypeString,
										Required: true,
									},

									"project_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"user_name": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"password": {
										Type:      schema.TypeString,
										Optional:  true,
										Sensitive: true,
									},
								},
							},
						},
						"fgs_forwarding": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"func_urn": {
										Type:     schema.TypeString,
										Required: true,
									},
									"func_name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildCreateDataForwardingRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	rst := map[string]interface{}{
		"rule_name":   d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"select":      utils.ValueIgnoreEmpty(d.Get("select")),
		"where":       utils.ValueIgnoreEmpty(d.Get("where")),
		"app_id":      utils.ValueIgnoreEmpty(d.Get("space_id")),
	}

	triggers := strings.SplitN(d.Get("trigger").(string), ":", 2)
	if len(triggers) == 2 {
		rst["subject"] = map[string]interface{}{
			"resource": triggers[0],
			"event":    triggers[1],
		}
	}

	return rst
}

func createRoutingRule(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	requestPath := client.Endpoint + "v5/iot/{project_id}/routing-rule/rules"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDataForwardingRuleBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func buildChannelDetailHttpForwardingBodyParams(rawMap map[string]interface{}) (map[string]interface{}, error) {
	forward := rawMap["http_forwarding"].([]interface{})
	if len(forward) == 0 {
		return nil, errors.New("http_forwarding is Required when the target type is HTTP_FORWARDING")
	}

	fMap := forward[0].(map[string]interface{})
	return map[string]interface{}{
		"http_forwarding": map[string]interface{}{
			"url": fMap["url"],
		},
	}, nil
}

func buildChannelDetailProjectId(fMap map[string]interface{}, projectId string) string {
	if rawValue := fMap["project_id"].(string); rawValue != "" {
		return rawValue
	}

	return projectId
}

func buildChannelDetailDisForwardingBodyParams(rawMap map[string]interface{}, projectId string) (map[string]interface{}, error) {
	forward := rawMap["dis_forwarding"].([]interface{})
	if len(forward) == 0 {
		return nil, errors.New("dis_forwarding is Required when the target type is DIS_FORWARDING")
	}

	fMap := forward[0].(map[string]interface{})
	return map[string]interface{}{
		"dis_forwarding": map[string]interface{}{
			"region_name": fMap["region"],
			"project_id":  buildChannelDetailProjectId(fMap, projectId),
			"stream_id":   fMap["stream_id"],
		},
	}, nil
}

func buildChannelDetailObsForwardingBodyParams(rawMap map[string]interface{}, projectId string) (map[string]interface{}, error) {
	forward := rawMap["obs_forwarding"].([]interface{})
	if len(forward) == 0 {
		return nil, errors.New("obs_forwarding is Required when the target type is OBS_FORWARDING")
	}

	fMap := forward[0].(map[string]interface{})
	return map[string]interface{}{
		"obs_forwarding": map[string]interface{}{
			"region_name": fMap["region"],
			"project_id":  buildChannelDetailProjectId(fMap, projectId),
			"bucket_name": fMap["bucket"],
			"file_path":   utils.ValueIgnoreEmpty(fMap["custom_directory"]),
		},
	}, nil
}

func buildChannelDetailAmqpForwardingBodyParams(rawMap map[string]interface{}) (map[string]interface{}, error) {
	forward := rawMap["amqp_forwarding"].([]interface{})
	if len(forward) == 0 {
		return nil, errors.New("amqp_forwarding is Required when the target type is AMQP_FORWARDING")
	}

	fMap := forward[0].(map[string]interface{})
	return map[string]interface{}{
		"amqp_forwarding": map[string]interface{}{
			"queue_name": fMap["queue_name"],
		},
	}, nil
}

func buildChannelDetailKafkaForwardingBodyParams(rawMap map[string]interface{}, projectId string) (map[string]interface{}, error) {
	forward := rawMap["kafka_forwarding"].([]interface{})
	if len(forward) == 0 {
		return nil, errors.New("kafka_forwarding is Required when the target type is DMS_KAFKA_FORWARDING")
	}

	fMap := forward[0].(map[string]interface{})
	addressesRaw := fMap["addresses"].([]interface{})
	addresses := make([]interface{}, len(addressesRaw))
	for i, item := range addressesRaw {
		itemMap := item.(map[string]interface{})
		addresses[i] = map[string]interface{}{
			"ip":     itemMap["ip"],
			"port":   itemMap["port"],
			"domain": itemMap["domain"],
		}
	}

	return map[string]interface{}{
		"dms_kafka_forwarding": map[string]interface{}{
			"region_name": fMap["region"],
			"project_id":  buildChannelDetailProjectId(fMap, projectId),
			"topic":       fMap["topic"],
			"username":    fMap["user_name"],
			"password":    fMap["password"],
			"addresses":   addresses,
		},
	}, nil
}

func buildChannelDetailFgsForwardingBodyParams(rawMap map[string]interface{}) (map[string]interface{}, error) {
	forward := rawMap["fgs_forwarding"].([]interface{})
	if len(forward) == 0 {
		return nil, errors.New("fgs_forwarding is Required when the target type is FUNCTIONGRAPH_FORWARDING")
	}

	fMap := forward[0].(map[string]interface{})
	return map[string]interface{}{
		"functiongraph_forwarding": map[string]interface{}{
			"func_urn":  fMap["func_urn"],
			"func_name": fMap["func_name"],
		},
	}, nil
}

func buildActionChannelDetailBodyParams(rawMap map[string]interface{}, projectId string) (map[string]interface{}, error) {
	switch rawMap["type"].(string) {
	case "HTTP_FORWARDING":
		return buildChannelDetailHttpForwardingBodyParams(rawMap)
	case "DIS_FORWARDING":
		return buildChannelDetailDisForwardingBodyParams(rawMap, projectId)
	case "OBS_FORWARDING":
		return buildChannelDetailObsForwardingBodyParams(rawMap, projectId)
	case "AMQP_FORWARDING":
		return buildChannelDetailAmqpForwardingBodyParams(rawMap)
	case "DMS_KAFKA_FORWARDING":
		return buildChannelDetailKafkaForwardingBodyParams(rawMap, projectId)
	case "FUNCTIONGRAPH_FORWARDING":
		return buildChannelDetailFgsForwardingBodyParams(rawMap)
	}

	return nil, fmt.Errorf("the target type %s is not support", rawMap["type"].(string))
}

func buildCreateRuleActionBodyParams(rawMap map[string]interface{}, ruleId, projectId string) (map[string]interface{}, error) {
	channelDetailBodyParam, err := buildActionChannelDetailBodyParams(rawMap, projectId)
	if err != nil {
		return nil, err
	}

	rst := map[string]interface{}{
		"rule_id":        ruleId,
		"channel":        rawMap["type"],
		"channel_detail": channelDetailBodyParam,
	}

	return rst, nil
}

func createRuleAction(client *golangsdk.ServiceClient, targetMap map[string]interface{}, ruleId string) error {
	bodyParams, err := buildCreateRuleActionBodyParams(targetMap, ruleId, client.ProjectID)
	if err != nil {
		return err
	}

	requestPath := client.Endpoint + "v5/iot/{project_id}/routing-rule/actions"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(bodyParams),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	return err
}

func createRuleActions(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	targets := d.Get("targets").(*schema.Set).List()
	for _, v := range targets {
		rawMap := v.(map[string]interface{})
		if err := createRuleAction(client, rawMap, d.Id()); err != nil {
			return err
		}
	}

	return nil
}

func enableRoutingRule(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v5/iot/{project_id}/routing-rule/rules/{rule_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{rule_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"active": true,
		},
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func ResourceDataForwardingRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	respBody, err := createRoutingRule(client, d)
	if err != nil {
		return diag.Errorf("error creating IoTDA data forwarding rule: %s", err)
	}

	ruleId := utils.PathSearch("rule_id", respBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("error creating IoTDA data forwarding rule: ID is not found in API response")
	}

	d.SetId(ruleId)

	if _, ok := d.GetOk("targets"); ok {
		if err := createRuleActions(client, d); err != nil {
			return diag.Errorf("error creating IoTDA data forwarding rule actions in creation operation: %s", err)
		}
	}

	if d.Get("enabled").(bool) {
		if err := enableRoutingRule(client, d); err != nil {
			return diag.Errorf("error activating the IoTDA data forwarding rule: %s", err)
		}
	}

	return ResourceDataForwardingRuleRead(ctx, d, meta)
}

func flattenTriggerAttribute(respBody interface{}) string {
	subject := utils.PathSearch("subject", respBody, nil)
	if subject == nil {
		return ""
	}

	resource := utils.PathSearch("resource", subject, "").(string)
	event := utils.PathSearch("event", subject, "").(string)
	return fmt.Sprintf("%s:%s", resource, event)
}

func buildRuleActionsQueryParams(ruleID, marker string) string {
	rst := fmt.Sprintf("?rule_id=%s", ruleID)
	if marker != "" {
		rst += fmt.Sprintf("&marker=%s", marker)
	}

	return rst
}

func queryRuleActions(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v5/iot/{project_id}/routing-rule/actions"
		allActions []interface{}
		marker     string
	)

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithMarker := requestPath + buildRuleActionsQueryParams(d.Id(), marker)
		resp, err := client.Request("GET", requestPathWithMarker, &requestOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		actions := utils.PathSearch("actions", respBody, make([]interface{}, 0)).([]interface{})
		if len(actions) == 0 {
			break
		}

		allActions = append(allActions, actions...)
		marker = utils.PathSearch("marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return allActions, nil
}

func queryDataForwardingRules(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	requestPath := client.Endpoint + "v5/iot/{project_id}/routing-rule/rules/{rule_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{rule_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func flattenTargetsHttpForwarding(actionId, channel string, httpForwarding interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":   actionId,
		"type": channel,
		"http_forwarding": []interface{}{
			map[string]interface{}{
				"url": utils.PathSearch("url", httpForwarding, nil),
			},
		},
	}
}

func flattenTargetsDisForwarding(actionId, channel string, disForwarding interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":   actionId,
		"type": channel,
		"dis_forwarding": []interface{}{
			map[string]interface{}{
				"region":     utils.PathSearch("region_name", disForwarding, nil),
				"project_id": utils.PathSearch("project_id", disForwarding, nil),
				"stream_id":  utils.PathSearch("stream_id", disForwarding, nil),
			},
		},
	}
}

func flattenTargetsObsForwarding(actionId, channel string, obsForwarding interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":   actionId,
		"type": channel,
		"obs_forwarding": []interface{}{
			map[string]interface{}{
				"region":           utils.PathSearch("region_name", obsForwarding, nil),
				"project_id":       utils.PathSearch("project_id", obsForwarding, nil),
				"bucket":           utils.PathSearch("bucket_name", obsForwarding, nil),
				"custom_directory": utils.PathSearch("file_path", obsForwarding, nil),
			},
		},
	}
}

func flattenTargetsAmqpForwarding(actionId, channel string, amqpForwarding interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":   actionId,
		"type": channel,
		"amqp_forwarding": []interface{}{
			map[string]interface{}{
				"queue_name": utils.PathSearch("queue_name", amqpForwarding, nil),
			},
		},
	}
}

func flattenDmsKafkaForwardingAddresses(addresses []interface{}) []interface{} {
	rst := make([]interface{}, len(addresses))
	for i, v := range addresses {
		rst[i] = map[string]interface{}{
			"ip":     utils.PathSearch("ip", v, nil),
			"port":   utils.PathSearch("port", v, nil),
			"domain": utils.PathSearch("domain", v, nil),
		}
	}

	return rst
}

func flattenTargetsDmsKafkaForwarding(actionId, channel string, dmsKafkaForwarding interface{}) map[string]interface{} {
	addresses := utils.PathSearch("addresses", dmsKafkaForwarding, make([]interface{}, 0)).([]interface{})
	return map[string]interface{}{
		"id":   actionId,
		"type": channel,
		"kafka_forwarding": []interface{}{
			map[string]interface{}{
				"region":     utils.PathSearch("region_name", dmsKafkaForwarding, nil),
				"project_id": utils.PathSearch("project_id", dmsKafkaForwarding, nil),
				"topic":      utils.PathSearch("topic", dmsKafkaForwarding, nil),
				"user_name":  utils.PathSearch("username", dmsKafkaForwarding, nil),
				"addresses":  flattenDmsKafkaForwardingAddresses(addresses),
			},
		},
	}
}

func flattenTargetsFunctionGraphForwarding(actionId, channel string, functionGraphForwarding interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":   actionId,
		"type": channel,
		"fgs_forwarding": []interface{}{
			map[string]interface{}{
				"func_urn":  utils.PathSearch("func_urn", functionGraphForwarding, nil),
				"func_name": utils.PathSearch("func_name", functionGraphForwarding, nil),
			},
		},
	}
}

func flattenTargetsAttribute(allActions []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(allActions))
	for _, v := range allActions {
		channel := utils.PathSearch("channel", v, "").(string)
		channelDetail := utils.PathSearch("channel_detail", v, nil)
		if channel == "" || channelDetail == nil {
			continue
		}

		actionId := utils.PathSearch("action_id", v, "").(string)
		switch channel {
		case "HTTP_FORWARDING":
			httpForwarding := utils.PathSearch("http_forwarding", channelDetail, nil)
			if httpForwarding != nil {
				rst = append(rst, flattenTargetsHttpForwarding(actionId, channel, httpForwarding))
			}
		case "DIS_FORWARDING":
			disForwarding := utils.PathSearch("dis_forwarding", channelDetail, nil)
			if disForwarding != nil {
				rst = append(rst, flattenTargetsDisForwarding(actionId, channel, disForwarding))
			}
		case "OBS_FORWARDING":
			obsForwarding := utils.PathSearch("obs_forwarding", channelDetail, nil)
			if obsForwarding != nil {
				rst = append(rst, flattenTargetsObsForwarding(actionId, channel, obsForwarding))
			}
		case "AMQP_FORWARDING":
			amqpForwarding := utils.PathSearch("amqp_forwarding", channelDetail, nil)
			if amqpForwarding != nil {
				rst = append(rst, flattenTargetsAmqpForwarding(actionId, channel, amqpForwarding))
			}
		case "DMS_KAFKA_FORWARDING":
			dmsKafkaForwarding := utils.PathSearch("dms_kafka_forwarding", channelDetail, nil)
			if dmsKafkaForwarding != nil {
				rst = append(rst, flattenTargetsDmsKafkaForwarding(actionId, channel, dmsKafkaForwarding))
			}
		case "FUNCTIONGRAPH_FORWARDING":
			functionGraphForwarding := utils.PathSearch("functiongraph_forwarding", channelDetail, nil)
			if functionGraphForwarding != nil {
				rst = append(rst, flattenTargetsFunctionGraphForwarding(actionId, channel, functionGraphForwarding))
			}
		}
	}

	return rst
}

func ResourceDataForwardingRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	respBody, err := queryDataForwardingRules(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA data forwarding rule")
	}

	allActions, err := queryRuleActions(client, d)
	if err != nil {
		return diag.Errorf("error retrieving IoTDA rule actions: %s", err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("rule_name", respBody, nil)),
		d.Set("trigger", flattenTriggerAttribute(respBody)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("select", utils.PathSearch("select", respBody, nil)),
		d.Set("where", utils.PathSearch("where", respBody, nil)),
		d.Set("enabled", utils.PathSearch("active", respBody, nil)),
		d.Set("space_id", utils.PathSearch("app_id", respBody, nil)),
		d.Set("targets", flattenTargetsAttribute(allActions)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func deleteRuleAction(client *golangsdk.ServiceClient, actionId string) error {
	requestPath := client.Endpoint + "v5/iot/{project_id}/routing-rule/actions/{action_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{action_id}", actionId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("DELETE", requestPath, &requestOpt)
	return err
}

func deleteRuleActions(client *golangsdk.ServiceClient, delRaws []interface{}) error {
	for _, v := range delRaws {
		ruleAction := v.(map[string]interface{})
		err := deleteRuleAction(client, ruleAction["id"].(string))
		if err != nil {
			// When the action ID does not exist, the API response status code is `404`.
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				continue
			}

			return fmt.Errorf("error deleting IoTDA data forwarding rule action: %s", err)
		}
	}

	return nil
}

func buildUpdateRuleActionBodyParams(rawMap map[string]interface{}, projectId string) (map[string]interface{}, error) {
	channelDetailBodyParam, err := buildActionChannelDetailBodyParams(rawMap, projectId)
	if err != nil {
		return nil, err
	}

	rst := map[string]interface{}{
		"channel":        rawMap["type"],
		"channel_detail": channelDetailBodyParam,
	}

	return rst, nil
}

func updateRuleAction(client *golangsdk.ServiceClient, targetMap map[string]interface{}, actionId string) error {
	bodyParams, err := buildUpdateRuleActionBodyParams(targetMap, client.ProjectID)
	if err != nil {
		return err
	}

	requestPath := client.Endpoint + "v5/iot/{project_id}/routing-rule/actions/{action_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{action_id}", actionId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(bodyParams),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	return err
}

func createOrUpdateRuleActions(client *golangsdk.ServiceClient, raws []interface{}, ruleId string) error {
	for _, v := range raws {
		rawMap := v.(map[string]interface{})
		if id, ok := rawMap["id"]; ok && len(id.(string)) > 0 {
			if err := updateRuleAction(client, rawMap, id.(string)); err != nil {
				return fmt.Errorf("error updating IoTDA data forwarding rule action: %s", err)
			}
		} else {
			if err := createRuleAction(client, rawMap, ruleId); err != nil {
				return fmt.Errorf("error creating IoTDA data forwarding rule action in update operation: %s", err)
			}
		}
	}

	return nil
}

func buildUpdateRoutingRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"rule_name":   d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"select":      utils.ValueIgnoreEmpty(d.Get("select")),
		"where":       utils.ValueIgnoreEmpty(d.Get("where")),
		"active":      d.Get("enabled"),
	}
}

func updateRoutingRule(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v5/iot/{project_id}/routing-rule/rules/{rule_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{rule_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateRoutingRuleBodyParams(d)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func ResourceDataForwardingRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	if d.HasChange("targets") {
		o, n := d.GetChange("targets")
		oSet := o.(*schema.Set)
		nSet := n.(*schema.Set)

		if err := deleteRuleActions(client, oSet.Difference(nSet).List()); err != nil {
			return diag.FromErr(err)
		}

		if err := createOrUpdateRuleActions(client, nSet.List(), d.Id()); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("name", "description", "select", "where", "enabled") {
		if err := updateRoutingRule(client, d); err != nil {
			return diag.Errorf("error updating IoTDA data forwarding rule: %s", err)
		}
	}

	return ResourceDataForwardingRuleRead(ctx, d, meta)
}

func deleteRoutingRule(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v5/iot/{project_id}/routing-rule/rules/{rule_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{rule_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("DELETE", requestPath, &requestOpt)
	return err
}

func ResourceDataForwardingRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	targets := d.Get("targets").(*schema.Set)
	if err := deleteRuleActions(client, targets.List()); err != nil {
		return diag.FromErr(err)
	}

	if err := deleteRoutingRule(client, d); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IoTDA data forwarding rule")
	}

	return nil
}
