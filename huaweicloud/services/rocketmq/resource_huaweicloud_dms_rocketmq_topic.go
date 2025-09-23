package rocketmq

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RocketMQ POST /v2/{project_id}/instances/{instance_id}/topics
// @API RocketMQ DELETE /v2/{project_id}/instances/{instance_id}/topics/{topic}
// @API RocketMQ GET /v2/{project_id}/instances/{instance_id}/topics/{topic}
// @API RocketMQ PUT /v2/{project_id}/instances/{instance_id}/topics/{topic}
func ResourceDmsRocketMQTopic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRocketMQTopicCreate,
		UpdateContext: resourceDmsRocketMQTopicUpdate,
		ReadContext:   resourceDmsRocketMQTopicRead,
		DeleteContext: resourceDmsRocketMQTopicDelete,
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
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the topic.`,
			},
			"message_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Specifies the message type of the topic.",
			},
			"brokers": {
				Type:        schema.TypeList,
				Elem:        rocketMQTopicBrokerRefSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Specifies the list of associated brokers of the topic.",
			},
			"queue_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the number of queues.`,
			},
			"queues": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the queue info of the topic.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"broker": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the associated broker.`,
						},
						"queue_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `Specifies the number of the queues.`,
						},
					},
				},
			},
			"permission": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the permissions of the topic.`,
				ValidateFunc: validation.StringInSlice([]string{
					"all", "sub", "pub",
				}, false),
			},
			"total_read_queue_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the total number of read queues.`,
			},
			"total_write_queue_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the total number of write queues.`,
			},
		},
	}
}

func rocketMQTopicBrokerRefSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Indicates the name of the broker.`,
			},
			"read_queue_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the read queues number of the broker.`,
			},
			"write_queue_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the read queues number of the broker.`,
			},
		},
	}
	return &sc
}

func resourceDmsRocketMQTopicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createRocketmqTopic: create DMS rocketmq topic
	var (
		createRocketmqTopicHttpUrl = "v2/{project_id}/instances/{instance_id}/topics"
		createRocketmqTopicProduct = "dmsv2"
	)
	createRocketmqTopicClient, err := cfg.NewServiceClient(createRocketmqTopicProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQTopic Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createRocketmqTopicPath := createRocketmqTopicClient.Endpoint + createRocketmqTopicHttpUrl
	createRocketmqTopicPath = strings.ReplaceAll(createRocketmqTopicPath, "{project_id}", createRocketmqTopicClient.ProjectID)
	createRocketmqTopicPath = strings.ReplaceAll(createRocketmqTopicPath, "{instance_id}", instanceID)

	createRocketmqTopicOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createRocketmqTopicOpt.JSONBody = utils.RemoveNil(buildCreateRocketmqTopicBodyParams(d))
	createRocketmqTopicResp, err := createRocketmqTopicClient.Request("POST", createRocketmqTopicPath, &createRocketmqTopicOpt)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQTopic: %s", err)
	}

	createRocketmqTopicRespBody, err := utils.FlattenResponse(createRocketmqTopicResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createRocketmqTopicRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find topic ID from the API response")
	}
	d.SetId(instanceID + "/" + id)

	return resourceDmsRocketMQTopicUpdate(ctx, d, meta)
}

func buildCreateRocketmqTopicBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":         d.Get("name"),
		"message_type": utils.ValueIgnoreEmpty(d.Get("message_type")),
		"brokers":      utils.ValueIgnoreEmpty(buildCreateRocketmqTopicBrokersChildBody(d)),
		"queue_num":    utils.ValueIgnoreEmpty(d.Get("queue_num")),
		"queues":       utils.ValueIgnoreEmpty(buildCreateRocketmqTopicQueuesChildBody(d)),
		"permission":   utils.ValueIgnoreEmpty(d.Get("permission")),
	}
	return bodyParams
}

func buildCreateRocketmqTopicBrokersChildBody(d *schema.ResourceData) []string {
	rawParams := d.Get("brokers").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}
	params := make([]string, 0)
	for _, param := range rawParams {
		params = append(params, utils.PathSearch("name", param, "").(string))
	}
	return params
}

func buildCreateRocketmqTopicQueuesChildBody(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("queues").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}
	params := make([]map[string]interface{}, 0)
	for _, rawParam := range rawParams {
		param, ok := rawParam.(map[string]interface{})
		if !ok {
			continue
		}
		params = append(params, map[string]interface{}{
			"broker":    param["broker"],
			"queue_num": param["queue_num"],
		})
	}
	return params
}

func resourceDmsRocketMQTopicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateRocketmqTopicHasChanges := []string{
		"total_read_queue_num",
		"total_write_queue_num",
		"permission",
	}

	if d.HasChanges(updateRocketmqTopicHasChanges...) {
		// updateRocketmqTopic: update DMS rocketmq topic
		var (
			updateRocketmqTopicHttpUrl = "v2/{project_id}/instances/{instance_id}/topics/{topic}"
			updateRocketmqTopicProduct = "dmsv2"
		)
		updateRocketmqTopicClient, err := cfg.NewServiceClient(updateRocketmqTopicProduct, region)
		if err != nil {
			return diag.Errorf("error creating DmsRocketMQTopic Client: %s", err)
		}

		parts := strings.SplitN(d.Id(), "/", 2)
		if len(parts) != 2 {
			return diag.Errorf("invalid id format, must be <instance_id>/<topic>")
		}
		instanceID := parts[0]
		topic := parts[1]
		updateRocketmqTopicPath := updateRocketmqTopicClient.Endpoint + updateRocketmqTopicHttpUrl
		updateRocketmqTopicPath = strings.ReplaceAll(updateRocketmqTopicPath, "{project_id}", updateRocketmqTopicClient.ProjectID)
		updateRocketmqTopicPath = strings.ReplaceAll(updateRocketmqTopicPath, "{instance_id}", instanceID)
		updateRocketmqTopicPath = strings.ReplaceAll(updateRocketmqTopicPath, "{topic}", topic)

		updateRocketmqTopicOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		updateRocketmqTopicOpt.JSONBody = utils.RemoveNil(buildUpdateRocketmqTopicBodyParams(d, cfg))
		_, err = updateRocketmqTopicClient.Request("PUT", updateRocketmqTopicPath, &updateRocketmqTopicOpt)
		if err != nil {
			return diag.Errorf("error updating DmsRocketMQTopic: %s", err)
		}
	}
	return resourceDmsRocketMQTopicRead(ctx, d, meta)
}

func buildUpdateRocketmqTopicBodyParams(d *schema.ResourceData, _ *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"read_queue_num":  utils.ValueIgnoreEmpty(d.Get("total_read_queue_num")),
		"write_queue_num": utils.ValueIgnoreEmpty(d.Get("total_write_queue_num")),
		"permission":      utils.ValueIgnoreEmpty(d.Get("permission")),
	}
	return bodyParams
}

func resourceDmsRocketMQTopicRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getRocketmqTopic: query DMS rocketmq topic
	var (
		getRocketmqTopicHttpUrl = "v2/{project_id}/instances/{instance_id}/topics/{topic}"
		getRocketmqTopicProduct = "dmsv2"
	)
	getRocketmqTopicClient, err := cfg.NewServiceClient(getRocketmqTopicProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQTopic Client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<topic>")
	}
	instanceID := parts[0]
	topic := parts[1]
	getRocketmqTopicPath := getRocketmqTopicClient.Endpoint + getRocketmqTopicHttpUrl
	getRocketmqTopicPath = strings.ReplaceAll(getRocketmqTopicPath, "{project_id}", getRocketmqTopicClient.ProjectID)
	getRocketmqTopicPath = strings.ReplaceAll(getRocketmqTopicPath, "{instance_id}", instanceID)
	getRocketmqTopicPath = strings.ReplaceAll(getRocketmqTopicPath, "{topic}", topic)

	getRocketmqTopicOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRocketmqTopicResp, err := getRocketmqTopicClient.Request("GET", getRocketmqTopicPath, &getRocketmqTopicOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DmsRocketMQTopic")
	}

	getRocketmqTopicRespBody, err := utils.FlattenResponse(getRocketmqTopicResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", topic),
		d.Set("total_read_queue_num", utils.PathSearch("total_read_queue_num", getRocketmqTopicRespBody, nil)),
		d.Set("total_write_queue_num", utils.PathSearch("total_write_queue_num", getRocketmqTopicRespBody, nil)),
		d.Set("permission", utils.PathSearch("permission", getRocketmqTopicRespBody, nil)),
		d.Set("brokers", flattenGetRocketmqTopicResponseBodyBrokerRef(getRocketmqTopicRespBody)),
		d.Set("message_type", utils.PathSearch("message_type", getRocketmqTopicRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetRocketmqTopicResponseBodyBrokerRef(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("brokers", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":            utils.PathSearch("broker_name", v, nil),
			"read_queue_num":  utils.PathSearch("read_queue_num", v, nil),
			"write_queue_num": utils.PathSearch("write_queue_num", v, nil),
		})
	}
	return rst
}

func resourceDmsRocketMQTopicDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteRocketmqTopic: delete DMS rocketmq topic
	var (
		deleteRocketmqTopicHttpUrl = "v2/{project_id}/instances/{instance_id}/topics/{topic}"
		deleteRocketmqTopicProduct = "dmsv2"
	)
	deleteRocketmqTopicClient, err := cfg.NewServiceClient(deleteRocketmqTopicProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQTopic Client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<topic>")
	}
	instanceID := parts[0]
	topic := parts[1]
	deleteRocketmqTopicPath := deleteRocketmqTopicClient.Endpoint + deleteRocketmqTopicHttpUrl
	deleteRocketmqTopicPath = strings.ReplaceAll(deleteRocketmqTopicPath, "{project_id}", deleteRocketmqTopicClient.ProjectID)
	deleteRocketmqTopicPath = strings.ReplaceAll(deleteRocketmqTopicPath, "{instance_id}", instanceID)
	deleteRocketmqTopicPath = strings.ReplaceAll(deleteRocketmqTopicPath, "{topic}", topic)

	deleteRocketmqTopicOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteRocketmqTopicClient.Request("DELETE", deleteRocketmqTopicPath, &deleteRocketmqTopicOpt)
	if err != nil {
		return diag.Errorf("error deleting DmsRocketMQTopic: %s", err)
	}

	return nil
}
