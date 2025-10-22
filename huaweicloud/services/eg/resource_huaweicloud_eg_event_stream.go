package eg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

// @API EG POST /v1/{project_id}/eventstreamings
// @API EG GET /v1/{project_id}/eventstreamings/{eventstreaming_id}
// @API EG PUT /v1/{project_id}/eventstreamings/{eventstreaming_id}
// @API EG DELETE /v1/{project_id}/eventstreamings/{eventstreaming_id}
// @API EG POST /v1/{project_id}/eventstreamings/operate/{eventstreaming_id}
func ResourceEventStream() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventStreamCreate,
		ReadContext:   resourceEventStreamRead,
		UpdateContext: resourceEventStreamUpdate,
		DeleteContext: resourceEventStreamDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the event stream is located.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the event stream.",
			},
			"source": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The name of the event source type.",
						},
						"kafka": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							ExactlyOneOf: []string{
								"source.0.mobile_rocketmq",
								"source.0.community_rocketmq",
								"source.0.dms_rocketmq",
							},
							Description: "The event source configuration detail for DMS Kafka type.",
						},
						"mobile_rocketmq": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  "The event source configuration detail for mobile RocketMQ type.",
						},
						"community_rocketmq": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  "The event source configuration detail for community RocketMQ type.",
						},
						"dms_rocketmq": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  "The event source configuration detail for DMS RocketMQ type.",
						},
					},
				},
				MaxItems:    1,
				Description: "The source configuration of the event stream.",
			},
			"sink": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the event target type.",
						},
						"functiongraph": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							ExactlyOneOf: []string{
								"sink.0.kafka",
							},
							Description: "The event target configuration detail for FunctionGraph type.",
						},
						"kafka": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  "The event target configuration detail for DMS Kafka type.",
						},
					},
				},
				MaxItems:    1,
				Description: "The target configuration of the event stream.",
			},
			"rule_config": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transform": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The type of transform rule.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The rule content definition.",
									},
									"template": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The template definition of the rule content.",
									},
								},
							},
							MaxItems:    1,
							Description: "The configuration detail of the transform rule.",
						},
						"filter": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  "The configuration detail of the filter rule.",
						},
					},
				},
				MaxItems:    1,
				Description: "The rule configuration of the event stream.",
			},
			"option": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"thread_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The number of concurrent threads.",
						},
						"batch_window": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"count": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The number of items pushed in batches.",
									},
									"time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The number of retries.",
									},
									"interval": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The interval of the batch push.",
									},
								},
							},
							MaxItems:    1,
							Description: "The configuration of the batch push.",
						},
					},
				},
				MaxItems:    1,
				Description: "The runtime configuration of the event stream.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the event stream.",
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The desired running status of the event stream.",
			},
			// Attributes
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the event stream.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The (UTC) creation time of the event stream, in RFC3339 format.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The (UTC) update time of the event stream, in RFC3339 format.",
			},
		},
	}
}

func unmarshalJsonFormatParamster(paramName, paramVal string) map[string]interface{} {
	parseResult := make(map[string]interface{})
	err := json.Unmarshal([]byte(paramVal), &parseResult)
	if err != nil {
		log.Printf("[ERROR] Invalid type of the %s, not json format", paramName)
	}
	return parseResult
}

func marshalJsonFormatParamster(paramName string, paramVal interface{}) interface{} {
	jsonFilter, err := json.Marshal(paramVal)
	if err != nil {
		log.Printf("[ERROR] unable to convert the %s, not json format", paramName)
	}
	return string(jsonFilter)
}

func buildCreateEventStreamBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"source":      buildEventStreamSource(d.Get("source").([]interface{})),
		"sink":        buildEventStreamSink(d.Get("sink").([]interface{})),
		"rule_config": buildEventStreamRuleConfig(d.Get("rule_config").([]interface{})),
		"option":      buildEventStreamRunOption(d.Get("option").([]interface{})),
		"description": utils.ValueIgnoreEmpty(d.Get("description").(string)),
	}
}

func buildEventStreamSource(eventSources []interface{}) interface{} {
	if len(eventSources) < 1 || eventSources[0] == nil {
		return nil
	}
	newSource := eventSources[0]
	return map[string]interface{}{
		"name": utils.PathSearch("name", newSource, "").(string),
		"source_kafka": utils.ValueIgnoreEmpty(unmarshalJsonFormatParamster("Kafka configuration",
			utils.PathSearch("kafka", newSource, "").(string))),
		"source_mobile_rocketmq": utils.ValueIgnoreEmpty(unmarshalJsonFormatParamster("mobile RocketMQ configuration",
			utils.PathSearch("mobile_rocketmq", newSource, "").(string))),
		"source_community_rocketmq_kafka": utils.ValueIgnoreEmpty(unmarshalJsonFormatParamster("community RocketMQ configuration",
			utils.PathSearch("community_rocketmq", newSource, "").(string))),
		"source_dms_rocketmq": utils.ValueIgnoreEmpty(unmarshalJsonFormatParamster("DMS RocketMQ configuration",
			utils.PathSearch("dms_rocketmq", newSource, "").(string))),
	}
}

func buildEventStreamSink(eventSinks []interface{}) interface{} {
	if len(eventSinks) < 1 || eventSinks[0] == nil {
		return nil
	}
	newSink := eventSinks[0]
	return map[string]interface{}{
		"name": utils.PathSearch("name", newSink, ""),
		"sink_fg": utils.ValueIgnoreEmpty(unmarshalJsonFormatParamster("FunctionGraph configuration",
			utils.PathSearch("functiongraph", newSink, "").(string))),
		"sink_kafka": utils.ValueIgnoreEmpty(unmarshalJsonFormatParamster("DMS Kafka",
			utils.PathSearch("kafka", newSink, "").(string))),
	}
}

func buildEventStreamRuleConfig(ruleConfigs []interface{}) interface{} {
	if len(ruleConfigs) < 1 || ruleConfigs[0] == nil {
		return nil
	}
	ruleConfig := ruleConfigs[0]
	return map[string]interface{}{
		"transform": buildEventStreamTransform(utils.PathSearch("transform", ruleConfig, make([]interface{}, 0)).([]interface{})),
		"filter":    utils.ValueIgnoreEmpty(unmarshalJsonFormatParamster("filter rule", utils.PathSearch("filter", ruleConfig, "").(string))),
	}
}

func buildEventStreamTransform(transforms []interface{}) interface{} {
	if len(transforms) < 1 || transforms[0] == nil {
		return nil
	}
	transform := transforms[0]
	return map[string]interface{}{
		"type":     utils.PathSearch("type", transform, nil),
		"value":    utils.ValueIgnoreEmpty(utils.PathSearch("value", transform, nil)),
		"template": utils.ValueIgnoreEmpty(utils.PathSearch("template", transform, nil)),
	}
}

func buildEventStreamRunOption(runOptions []interface{}) interface{} {
	if len(runOptions) < 1 || runOptions[0] == nil {
		return nil
	}
	runOption := runOptions[0]
	return map[string]interface{}{
		"thread_num":   utils.PathSearch("thread_num", runOption, nil),
		"batch_window": buildEventStreamBatchWindow(utils.PathSearch("batch_window", runOption, make([]interface{}, 0)).([]interface{})),
	}
}

func buildEventStreamBatchWindow(batchWindows []interface{}) interface{} {
	if len(batchWindows) < 1 || batchWindows[0] == nil {
		return nil
	}
	batchWindow := batchWindows[0]
	return map[string]interface{}{
		"count":    utils.PathSearch("count", batchWindow, nil),
		"time":     utils.PathSearch("time", batchWindow, nil),
		"interval": utils.PathSearch("interval", batchWindow, nil),
	}
}

func buildOperateEventStream(action string) interface{} {
	return map[string]interface{}{
		"operation": action,
	}
}

func doActionForEventStream(client *golangsdk.ServiceClient, resourceId, action string) error {
	httpUrl := "v1/{project_id}/eventstreamings/operate/{eventstreaming_id}"
	doActionPath := client.Endpoint + httpUrl
	doActionPath = strings.ReplaceAll(doActionPath, "{project_id}", client.ProjectID)
	doActionPath = strings.ReplaceAll(doActionPath, "{eventstreaming_id}", resourceId)
	doActionOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	doActionOpts.JSONBody = buildOperateEventStream(action)
	_, err := client.Request("POST", doActionPath, &doActionOpts)
	if err != nil {
		return fmt.Errorf("failed to operate the event stream: %s", err)
	}
	return nil
}

func resourceEventStreamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/eventstreamings"
	)
	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createOpts.JSONBody = utils.RemoveNil(buildCreateEventStreamBodyParams(d))
	requestResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating event stream: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	streamId := utils.PathSearch("eventStreamingID", respBody, "").(string)
	if streamId == "" {
		return diag.Errorf("unable to find the stream ID from API response")
	}
	d.SetId(streamId)

	if action, ok := d.GetOk("action"); ok && action.(string) == "START" {
		err := doActionForEventStream(client, d.Id(), action.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceEventStreamRead(ctx, d, meta)
}

func flattenEventStreamRuleConfig(ruleConfig interface{}) []interface{} {
	if ruleConfig == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"transform": flattenEventStreamTransform(utils.PathSearch("transform", ruleConfig, nil)),
			"filter":    marshalJsonFormatParamster("filter rule", utils.PathSearch("filter", ruleConfig, nil)),
		},
	}
}

func flattenEventStreamTransform(transformResp interface{}) []interface{} {
	if transformResp == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"type":     utils.PathSearch("type", transformResp, nil),
			"value":    utils.PathSearch("value", transformResp, nil),
			"template": utils.PathSearch("template", transformResp, nil),
		},
	}
}

func flattenEventStreamRunOption(runOption interface{}) []interface{} {
	if runOption == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"thread_num":   utils.PathSearch("thread_num", runOption, nil),
			"batch_window": flattenEventStreamBatchWindow(utils.PathSearch("batch_window", runOption, nil)),
		},
	}
}

func flattenEventStreamBatchWindow(batchWindow interface{}) []interface{} {
	if batchWindow == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"count":    utils.PathSearch("count", batchWindow, nil),
			"time":     utils.PathSearch("time", batchWindow, nil),
			"interval": utils.PathSearch("interval", batchWindow, nil),
		},
	}
}

func parseEventStreamAction(status string) string {
	if status == "RUNNING" {
		return "START"
	}
	// Record the error message and return the 'PAUSE' status.
	if status == "ERROR" {
		log.Printf("The event stream is running abnormally")
	}
	return "PAUSE"
}

func resourceEventStreamRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v1/{project_id}/eventstreamings/{eventstreaming_id}"
		streamId = d.Id()
	)
	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{eventstreaming_id}", streamId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving event stream")
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("rule_config", flattenEventStreamRuleConfig(utils.PathSearch("rule_config", respBody, nil))),
		d.Set("option", flattenEventStreamRunOption(utils.PathSearch("option", respBody, nil))),
		d.Set("action", parseEventStreamAction(utils.PathSearch("status", respBody, "").(string))),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("created_at", utils.PathSearch("created_time", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_time", respBody, nil)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving EG event stream (%s) fields: %s", streamId, err)
	}
	return nil
}

func buildUpdateEventStreamBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"source":      buildEventStreamSource(d.Get("source").([]interface{})),
		"sink":        buildEventStreamSink(d.Get("sink").([]interface{})),
		"rule_config": buildEventStreamRuleConfig(d.Get("rule_config").([]interface{})),
		"option":      buildEventStreamRunOption(d.Get("option").([]interface{})),
		"description": d.Get("description").(string),
	}
}

func resourceEventStreamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		httpUrl        = "v1/{project_id}/eventstreamings/{eventstreaming_id}"
		streamId       = d.Id()
		newAction      = d.Get("action").(string)
		actionReverted = false // Whether the action of the event stream is restored to the default value through an update operation.
	)
	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	if d.HasChangeExcept("action") {
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{eventstreaming_id}", streamId)
		updateOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		}

		updateOpts.JSONBody = utils.RemoveNil(buildUpdateEventStreamBodyParams(d))
		_, err = client.Request("PUT", updatePath, &updateOpts)
		if err != nil {
			return diag.Errorf("error updating event stream (%s): %s", streamId, err)
		}
		// The operation has reverted to the 'PAUSE' value.
		actionReverted = true
	}

	if actionReverted && newAction != "PAUSE" || !d.HasChangeExcept("action") && d.HasChange("action") {
		err := doActionForEventStream(client, streamId, newAction)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceEventStreamRead(ctx, d, meta)
}

func resourceEventStreamDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v1/{project_id}/eventstreamings/{eventstreaming_id}"
		streamId = d.Id()
	)
	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{eventstreaming_id}", streamId)
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting event stream (%s): %s", streamId, err)
	}
	return nil
}
