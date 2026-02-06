package rabbitmq

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RabbitMQ PUT /v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/queues
// @API RabbitMQ POST /v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/queues
// @API RabbitMQ GET /v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/queues/{queue}
func ResourceDmsRabbitmqQueue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRabbitmqQueueCreate,
		ReadContext:   resourceDmsRabbitmqQueueRead,
		DeleteContext: resourceDmsRabbitmqQueueDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceExchangeOrQueueImportState,
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
			"vhost": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auto_delete": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"durable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dead_letter_exchange": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dead_letter_routing_key": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"dead_letter_exchange"},
				ForceNew:     true,
			},
			"message_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"lazy_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"messages": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"consumers": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"consumer_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemeConsumerDetails(),
			},
			"queue_bindings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemeQueueBindings(),
			},
		},
	}
}

func schemeConsumerDetails() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"consumer_tag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"channel_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"ack_required": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"prefetch_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func schemeQueueBindings() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routing_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"properties_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDmsRabbitmqQueueCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	vhost := d.Get("vhost").(string)
	name := d.Get("name").(string)

	createHttpUrl := "v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/queues"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)
	createPath = strings.ReplaceAll(createPath, "{vhost}", vhost)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildRabbitmqQueueRequestBody(d)),
	}

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating queue: %s", err)
	}

	id := fmt.Sprintf("%s/%s/%s", instanceID, vhost, name)
	d.SetId(id)

	return resourceDmsRabbitmqQueueRead(ctx, d, cfg)
}

func buildRabbitmqQueueRequestBody(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                    d.Get("name"),
		"auto_delete":             d.Get("auto_delete"),
		"durable":                 utils.ValueIgnoreEmpty(d.Get("durable")),
		"dead_letter_exchange":    utils.ValueIgnoreEmpty(d.Get("dead_letter_exchange")),
		"dead_letter_routing_key": utils.ValueIgnoreEmpty(d.Get("dead_letter_routing_key")),
		"message_ttl":             utils.ValueIgnoreEmpty(d.Get("message_ttl")),
		"lazy_mode":               utils.ValueIgnoreEmpty(d.Get("lazy_mode")),
	}
	return bodyParams
}

func resourceDmsRabbitmqQueueRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	vhost := d.Get("vhost").(string)
	name := d.Get("name").(string)

	getRespBody, err := GetRabbitmqQueue(client, instanceID, vhost, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving the queue")
	}

	auguments := utils.PathSearch("arguments", getRespBody, nil)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("auto_delete", utils.PathSearch("auto_delete", getRespBody, nil)),
		d.Set("durable", utils.PathSearch("durable", getRespBody, nil)),
		d.Set("dead_letter_exchange", utils.PathSearch(`"x-dead-letter-exchange"`, auguments, nil)),
		d.Set("dead_letter_routing_key", utils.PathSearch(`"x-dead-letter-routing-key"`, auguments, nil)),
		d.Set("message_ttl", utils.PathSearch(`"x-message-ttl"`, auguments, nil)),
		d.Set("lazy_mode", utils.PathSearch(`"x-queue-mode"`, auguments, nil)),
		d.Set("messages", utils.PathSearch("messages", getRespBody, nil)),
		d.Set("consumers", utils.PathSearch("consumers", getRespBody, nil)),
		d.Set("policy", utils.PathSearch("policy", getRespBody, nil)),
		d.Set("consumer_details", flattenConsumerDetails(
			utils.PathSearch("consumer_details", getRespBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("queue_bindings", flattenQueueBingdings(
			utils.PathSearch("queue_bindings", getRespBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetRabbitmqQueue(client *golangsdk.ServiceClient, instanceID, vhost, name string) (interface{}, error) {
	getHttpUrl := "v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/queues/{queue}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)
	getPath = strings.ReplaceAll(getPath, "{vhost}", vhost)

	// queue name may have % or |
	if strings.Contains(name, "/") {
		replacedName := strings.ReplaceAll(name, "/", "__F_SLASH__")
		getPath = strings.ReplaceAll(getPath, "{queue}", url.PathEscape(replacedName))
	} else {
		getPath = strings.ReplaceAll(getPath, "{queue}", url.PathEscape(name))
	}

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func flattenConsumerDetails(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		rst = append(rst, map[string]interface{}{
			"consumer_tag": utils.PathSearch("consumer_tag", params, nil),
			"channel_details": flattenChannelDetails(utils.PathSearch("channel_details", params,
				make(map[string]interface{})).(map[string]interface{})),
			"ack_required":   utils.PathSearch("ack_required", params, nil),
			"prefetch_count": utils.PathSearch("prefetch_count", params, nil),
		})
	}

	return rst
}

func flattenChannelDetails(channelDetails map[string]interface{}) []map[string]interface{} {
	if len(channelDetails) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":            utils.PathSearch("name", channelDetails, nil),
			"number":          utils.PathSearch("number", channelDetails, nil),
			"user":            utils.PathSearch("user", channelDetails, nil),
			"connection_name": utils.PathSearch("connection_name", channelDetails, nil),
			"peer_host":       utils.PathSearch("peer_host", channelDetails, nil),
			"peer_port":       utils.PathSearch("peer_port", channelDetails, nil),
		},
	}
}

func flattenQueueBingdings(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		rst = append(rst, map[string]interface{}{
			"source":           utils.PathSearch("source", params, nil),
			"destination_type": utils.PathSearch("destination_type", params, nil),
			"destination":      utils.PathSearch("destination", params, nil),
			"routing_key":      utils.PathSearch("routing_key", params, nil),
			"properties_key":   utils.PathSearch("properties_key", params, nil),
		})
	}

	return rst
}

func resourceDmsRabbitmqQueueDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	vhost := d.Get("vhost").(string)
	name := d.Get("name").(string)

	deleteHttpUrl := "v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/queues"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceID)
	deletePath = strings.ReplaceAll(deletePath, "{vhost}", vhost)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		JSONBody: map[string]interface{}{
			"name": []string{name},
		},
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting queue")
	}

	return nil
}

func resourceExchangeOrQueueImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ",")
	if len(parts) != 3 {
		parts = strings.Split(d.Id(), "/")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid ID format, must be <instance_id>/<vhsot>/<name> or <instance_id>,<vhost>,<name>")
		}
	} else {
		// reform ID to be separated by slashes
		id := fmt.Sprintf("%s/%s/%s", parts[0], parts[1], parts[2])
		d.SetId(id)
	}

	d.Set("instance_id", parts[0])
	d.Set("vhost", parts[1])
	d.Set("name", parts[2])

	return []*schema.ResourceData{d}, nil
}
