package iotda

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	v5 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

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
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 256),
					validation.StringMatch(regexp.MustCompile(stringRegxp), stringFormatMsg),
				),
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
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 256),
				),
			},

			"select": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 500),
			},

			"where": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 500),
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
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP_FORWARDING",
								"DIS_FORWARDING",
								"OBS_FORWARDING",
								"AMQP_FORWARDING",
								"DMS_KAFKA_FORWARDING",
							}, false),
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
										ValidateFunc: validation.All(
											validation.StringLenBetween(0, 256),
											validation.StringMatch(regexp.MustCompile(`^[A-Za-z-_0-9/{}]*$`),
												"Only letters, digits, hyphens (-), underscores (_), slash (/)"+
													" and braces ({}) are allowed"),
											validation.StringDoesNotMatch(regexp.MustCompile(`^/|/$|//`),
												"Can not start or end with slashes (/), and cannot contain more than"+
													" two adjacent slashes (/)"),
										),
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

func ResourceDataForwardingRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	projectId := c.RegionProjectIDMap[region]
	createOpts := buildDataForwardingRuleCreateParams(d)
	log.Printf("[DEBUG] Create IoTDA data forwarding rule params: %#v", createOpts)

	resp, err := client.CreateRoutingRule(createOpts)
	if err != nil {
		return diag.Errorf("error creating IoTDA data forwarding rule: %s", err)
	}

	if resp.RuleId == nil {
		return diag.Errorf("error creating IoTDA data forwarding rule: id is not found in API response")
	}

	d.SetId(*resp.RuleId)
	m := d.Get("targets").(*schema.Set)
	// create action rule
	targets, err := buildActionTargets(m.List(), d.Id(), projectId)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, v := range targets {
		_, err := client.CreateRuleAction(&v)
		if err != nil {
			return diag.Errorf("error add targets to IoTDA data forwarding rule: %s", err)
		}
	}

	// enable forwarding
	if d.Get("enabled").(bool) {
		_, err = client.UpdateRoutingRule(&model.UpdateRoutingRuleRequest{
			RuleId: d.Id(),
			Body: &model.UpdateRuleReq{
				Active: utils.Bool(true),
			},
		})
		if err != nil {
			return diag.Errorf("error activing the IoTDA data forwarding rule: %s", err)
		}
	}

	return ResourceDataForwardingRuleRead(ctx, d, meta)
}

func ResourceDataForwardingRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	response, err := client.ShowRoutingRule(&model.ShowRoutingRuleRequest{RuleId: d.Id()})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA data forwarding rule")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", response.RuleName),
		d.Set("trigger", flattenTrigger(response.Subject)),
		d.Set("description", response.Description),
		d.Set("select", response.Select),
		d.Set("where", response.Where),
		d.Set("enabled", response.Active),
		d.Set("space_id", response.AppId),
		setTargetsToState(d, client, d.Id()),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func ResourceDataForwardingRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	projectId := c.RegionProjectIDMap[region]

	if d.HasChange("targets") {
		o, n := d.GetChange("targets")
		oldTargetSet := o.(*schema.Set)
		newTargetSet := n.(*schema.Set)

		// delete
		for _, v := range oldTargetSet.Difference(newTargetSet).List() {
			ruleAction := v.(map[string]interface{})
			_, err = client.DeleteRuleAction(&model.DeleteRuleActionRequest{ActionId: ruleAction["id"].(string)})
			if err != nil {
				return diag.Errorf("error updating targets of IoTDA data forwarding rule: %s", err)
			}
		}

		// add and update
		for _, v := range newTargetSet.List() {
			target := v.(map[string]interface{})
			channel := target["type"].(string)
			channelDetail, err := buildChannelDetail(target, channel, projectId)
			if err != nil {
				return diag.FromErr(err)
			}

			if id, ok := target["id"]; ok && len(id.(string)) > 0 {
				_, err = client.UpdateRuleAction(&model.UpdateRuleActionRequest{
					ActionId: id.(string),
					Body: &model.UpdateActionReq{
						Channel:       &channel,
						ChannelDetail: channelDetail,
					},
				})
				if err != nil {
					return diag.Errorf("error updating targets of IoTDA data forwarding rule: %s", err)
				}
			} else {
				_, err = client.CreateRuleAction(&model.CreateRuleActionRequest{
					Body: &model.AddActionReq{
						RuleId:        d.Id(),
						Channel:       channel,
						ChannelDetail: channelDetail,
					},
				})
				if err != nil {
					return diag.Errorf("error updating targets of IoTDA data forwarding rule: %s", err)
				}
			}

		}
	}

	// This update must be the last
	if d.HasChanges("name", "description", "select", "where", "enabled") {
		_, err = client.UpdateRoutingRule(&model.UpdateRoutingRuleRequest{
			RuleId: d.Id(),
			Body: &model.UpdateRuleReq{
				RuleName:    utils.String(d.Get("name").(string)),
				Description: utils.StringIgnoreEmpty(d.Get("description").(string)),
				Select:      utils.StringIgnoreEmpty(d.Get("select").(string)),
				Where:       utils.StringIgnoreEmpty(d.Get("where").(string)),
				Active:      utils.Bool(d.Get("enabled").(bool)),
			},
		})

		if err != nil {
			return diag.Errorf("error updating IoTDA data forwarding rule: %s", err)
		}
	}

	return ResourceDataForwardingRuleRead(ctx, d, meta)
}

func ResourceDataForwardingRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	//delete targets
	targets := d.Get("targets").(*schema.Set)
	for _, v := range targets.List() {
		ruleAction := v.(map[string]interface{})
		_, err = client.DeleteRuleAction(&model.DeleteRuleActionRequest{ActionId: ruleAction["id"].(string)})
		if err != nil {
			return diag.Errorf("error deleting targets of IoTDA data forwarding rule: %s", err)
		}
	}

	// delete data forwarding
	deleteOpts := &model.DeleteRoutingRuleRequest{RuleId: d.Id()}
	_, err = client.DeleteRoutingRule(deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting IoTDA data forwarding rule: %s", err)
	}

	return nil
}

func buildDataForwardingRuleCreateParams(d *schema.ResourceData) *model.CreateRoutingRuleRequest {
	triggers := strings.SplitN(d.Get("trigger").(string), ":", 2)
	req := model.CreateRoutingRuleRequest{
		Body: &model.AddRuleReq{
			RuleName:    utils.String(d.Get("name").(string)),
			Description: utils.StringIgnoreEmpty(d.Get("description").(string)),
			Select:      utils.StringIgnoreEmpty(d.Get("select").(string)),
			Where:       utils.StringIgnoreEmpty(d.Get("where").(string)),
			AppId:       utils.StringIgnoreEmpty(d.Get("space_id").(string)),
			Subject: &model.RoutingRuleSubject{
				Resource: triggers[0],
				Event:    triggers[1],
			},
		},
	}
	return &req
}

func buildActionTargets(raw []interface{}, ruleId, projectId string) ([]model.CreateRuleActionRequest, error) {
	rst := make([]model.CreateRuleActionRequest, len(raw))
	for i, v := range raw {
		target := v.(map[string]interface{})
		channel := target["type"].(string)
		channelDetail, err := buildChannelDetail(target, channel, projectId)
		if err != nil {
			return nil, err
		}
		rst[i] = model.CreateRuleActionRequest{
			Body: &model.AddActionReq{
				RuleId:        ruleId,
				Channel:       channel,
				ChannelDetail: channelDetail,
			},
		}
	}

	return rst, nil
}

func buildChannelDetail(target map[string]interface{}, channel, projectId string) (*model.ChannelDetail, error) {
	switch channel {
	case "HTTP_FORWARDING":
		forward := target["http_forwarding"].([]interface{})
		if len(forward) == 0 {
			return nil, fmt.Errorf("http_forwarding is Required when the target type is HTTP_FORWARDING")
		}
		f := forward[0].(map[string]interface{})
		d := model.ChannelDetail{
			HttpForwarding: &model.HttpForwarding{
				Url: f["url"].(string),
			},
		}
		return &d, nil

	case "DIS_FORWARDING":
		forward := target["dis_forwarding"].([]interface{})
		if len(forward) == 0 {
			return nil, fmt.Errorf("dis_forwarding is Required when the target type is DIS_FORWARDING")
		}
		f := forward[0].(map[string]interface{})
		projectIdStr := f["project_id"].(string)
		if projectIdStr == "" {
			projectIdStr = projectId
		}
		d := model.ChannelDetail{
			DisForwarding: &model.DisForwarding{
				RegionName: f["region"].(string),
				ProjectId:  projectIdStr,
				StreamId:   utils.String(f["stream_id"].(string)),
			},
		}
		return &d, nil

	case "OBS_FORWARDING":
		forward := target["obs_forwarding"].([]interface{})
		if len(forward) == 0 {
			return nil, fmt.Errorf("obs_forwarding is Required when the target type is OBS_FORWARDING")
		}
		f := forward[0].(map[string]interface{})
		projectIdStr := f["project_id"].(string)
		if projectIdStr == "" {
			projectIdStr = projectId
		}
		d := model.ChannelDetail{
			ObsForwarding: &model.ObsForwarding{
				RegionName: f["region"].(string),
				ProjectId:  projectIdStr,
				BucketName: f["bucket"].(string),
				FilePath:   utils.StringIgnoreEmpty(f["custom_directory"].(string)),
			},
		}
		return &d, nil

	case "AMQP_FORWARDING":
		forward := target["amqp_forwarding"].([]interface{})
		if len(forward) == 0 {
			return nil, fmt.Errorf("amqp_forwarding is Required when the target type is AMQP_FORWARDING")
		}
		f := forward[0].(map[string]interface{})
		d := model.ChannelDetail{
			AmqpForwarding: &model.AmqpForwarding{
				QueueName: f["queue_name"].(string),
			},
		}
		return &d, nil

	case "kafka_forwarding":
		forward := target["kafka_forwarding"].([]interface{})
		if len(forward) == 0 {
			return nil, fmt.Errorf("kafka_forwarding is Required when the target type is DMS_KAFKA_FORWARDING")
		}
		f := forward[0].(map[string]interface{})
		addressesRaw := f["addresses"].([]interface{})
		addresses := make([]model.NetAddress, len(addressesRaw))
		for i, item := range addressesRaw {
			item := item.(map[string]interface{})
			addresses[i] = model.NetAddress{
				Ip:     utils.String(item["ip"].(string)),
				Port:   utils.Int32(int32(item["port"].(int))),
				Domain: utils.String(item["domain"].(string)),
			}
		}

		projectIdStr := f["project_id"].(string)
		if projectIdStr == "" {
			projectIdStr = projectId
		}
		d := model.ChannelDetail{
			DmsKafkaForwarding: &model.DmsKafkaForwarding{
				RegionName: f["region"].(string),
				ProjectId:  projectIdStr,
				Topic:      f["topic"].(string),
				Username:   utils.String(f["user_name"].(string)),
				Password:   utils.String(f["password"].(string)),
				Addresses:  addresses,
			},
		}
		return &d, nil

	default:
		return nil, fmt.Errorf("the target type is %q is not support", channel)
	}
}

func setTargetsToState(d *schema.ResourceData, client *v5.IoTDAClient, id string) error {
	var rst []model.RoutingRuleAction
	var marker *string
	for {
		resp, err := client.ListRuleActions(&model.ListRuleActionsRequest{RuleId: utils.String(id), Marker: marker})
		if err != nil {
			return fmt.Errorf("error setting the targets: %s", err)
		}
		if resp.Actions == nil || len(*resp.Actions) == 0 {
			break
		}
		rst = append(rst, *resp.Actions...)
		marker = resp.Marker
	}

	return d.Set("targets", flattenTargets(rst))
}

func flattenTargets(s []model.RoutingRuleAction) []interface{} {
	rst := make([]interface{}, len(s))
	for i, v := range s {
		switch *v.Channel {
		case "HTTP_FORWARDING":
			if v.ChannelDetail != nil && v.ChannelDetail.HttpForwarding != nil {
				rst[i] = map[string]interface{}{
					"id":   v.ActionId,
					"type": v.Channel,
					"http_forwarding": []interface{}{
						map[string]interface{}{
							"url": v.ChannelDetail.HttpForwarding.Url,
						},
					},
				}
			}
		case "DIS_FORWARDING":
			if v.ChannelDetail != nil && v.ChannelDetail.DisForwarding != nil {
				rst[i] = map[string]interface{}{
					"id":   v.ActionId,
					"type": v.Channel,
					"dis_forwarding": []interface{}{
						map[string]interface{}{
							"region":     v.ChannelDetail.DisForwarding.RegionName,
							"project_id": v.ChannelDetail.DisForwarding.ProjectId,
							"stream_id":  v.ChannelDetail.DisForwarding.StreamId,
						},
					},
				}
			}
		case "OBS_FORWARDING":
			if v.ChannelDetail != nil && v.ChannelDetail.ObsForwarding != nil {
				rst[i] = map[string]interface{}{
					"id":   v.ActionId,
					"type": v.Channel,
					"obs_forwarding": []interface{}{
						map[string]interface{}{
							"region":           v.ChannelDetail.ObsForwarding.RegionName,
							"project_id":       v.ChannelDetail.ObsForwarding.ProjectId,
							"bucket":           v.ChannelDetail.ObsForwarding.BucketName,
							"custom_directory": v.ChannelDetail.ObsForwarding.FilePath,
						},
					},
				}
			}
		case "AMQP_FORWARDING":
			if v.ChannelDetail != nil && v.ChannelDetail.AmqpForwarding != nil {
				rst[i] = map[string]interface{}{
					"id":   v.ActionId,
					"type": v.Channel,
					"amqp_forwarding": []interface{}{
						map[string]interface{}{
							"queue_name": v.ChannelDetail.AmqpForwarding.QueueName,
						},
					},
				}
			}
		case "DMS_KAFKA_FORWARDING":
			if v.ChannelDetail != nil && v.ChannelDetail.DmsKafkaForwarding != nil {
				rst[i] = map[string]interface{}{
					"id":   v.ActionId,
					"type": v.Channel,
					"kafka_forwarding": []interface{}{
						map[string]interface{}{
							"region":     v.ChannelDetail.DmsKafkaForwarding.RegionName,
							"project_id": v.ChannelDetail.DmsKafkaForwarding.ProjectId,
							"topic":      v.ChannelDetail.DmsKafkaForwarding.Topic,
							"user_name":  v.ChannelDetail.DmsKafkaForwarding.Username,
							"addresses":  flattenAddress(v.ChannelDetail.DmsKafkaForwarding.Addresses),
						},
					},
				}
			}
		}
	}

	return rst
}

func flattenAddress(s []model.NetAddress) []interface{} {
	rst := make([]interface{}, len(s))
	for i, v := range s {
		var port *int
		if v.Port != nil {
			p := int(*v.Port)
			port = &p
		}
		rst[i] = map[string]interface{}{
			"ip":     v.Ip,
			"port":   port,
			"domain": v.Domain,
		}
	}
	return rst
}

func flattenTrigger(routingRuleSubject *model.RoutingRuleSubject) string {
	if routingRuleSubject != nil {
		return fmt.Sprintf("%s:%s", routingRuleSubject.Resource, routingRuleSubject.Event)
	}
	return ""
}
