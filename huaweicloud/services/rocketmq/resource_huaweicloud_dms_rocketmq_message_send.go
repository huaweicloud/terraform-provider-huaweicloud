package rocketmq

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var rocketmqMessageSendNonUpdatableParams = []string{
	"instance_id",
	"topic",
	"body",
	"property_list",
}

// @API RocketMQ POST /v2/{engine}/{project_id}/instances/{instance_id}/messages
func ResourceRocketMQMessageSend() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRocketMQMessageSendCreate,
		ReadContext:   resourceRocketMQMessageSendRead,
		UpdateContext: resourceRocketMQMessageSendUpdate,
		DeleteContext: resourceRocketMQMessageSendDelete,

		CustomizeDiff: config.FlexibleForceNew(rocketmqMessageSendNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the message to be sent is located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the RocketMQ instance.`,
			},
			"topic": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the topic to send the message.`,
			},
			"body": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The content of the message to be sent.`,
			},
			"property_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The list of message properties.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the property.`,
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The value of the property.`,
						},
					},
				},
			},
			"msg_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the message that was sent.`,
			},
			"queue_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The queue ID of the message.`,
			},
			"queue_offset": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The queue offset of the message.`,
			},
			"broker_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The broker name of the message.`,
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildRocketMQMessageSendBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"topic":         d.Get("topic").(string),
		"body":          d.Get("body").(string),
		"property_list": buildRocketMQMessageSendPropertyList(d.Get("property_list").([]interface{})),
	}
}

func buildRocketMQMessageSendPropertyList(properties []interface{}) []map[string]interface{} {
	if len(properties) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(properties))
	for i, property := range properties {
		result[i] = map[string]interface{}{
			"name":  utils.PathSearch("name", property, nil),
			"value": utils.PathSearch("value", property, nil),
		}
	}
	return result
}

func resourceRocketMQMessageSendCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v2/rocketmq/{project_id}/instances/{instance_id}/messages"
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildRocketMQMessageSendBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("unable to send message to RocketMQ topic: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate resource ID: %s", err)
	}
	d.SetId(randUUID)

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("msg_id", utils.PathSearch("msg_id", respBody, nil)),
		d.Set("queue_id", utils.PathSearch("queue_id", respBody, nil)),
		d.Set("queue_offset", utils.PathSearch("queue_offset", respBody, nil)),
		d.Set("broker_name", utils.PathSearch("broker_name", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceRocketMQMessageSendRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRocketMQMessageSendUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRocketMQMessageSendDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to send message to RocketMQ topic. Deleting this resource will
not clear the corresponding message record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
