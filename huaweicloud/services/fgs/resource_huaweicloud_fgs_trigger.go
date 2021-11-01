package fgs

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/fgs/v2/trigger"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const (
	timingTrigger = "TIMER"
	obsTrigger    = "OBS"
	smnTrigger    = "SMN"
	disTrigger    = "DIS"
	kafkaTrigger  = "KAFKA"

	obsEventCreated             = "ObjectCreated"
	obsEventPut                 = "Put"
	obsEventPost                = "Post"
	obsEventCopy                = "Copy"
	obsEventMultiUpload         = "CompleteMultipartUpload"
	obsEventRemoved             = "ObjectRemoved"
	obsEventDelete              = "Delete"
	obsEventDeleteWithoutMarker = "DeleteMarkerCreated"

	statusActive   = "ACTIVE"
	statusDisabled = "DISABLED"
)

func ResourceFunctionGraphTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceFunctionGraphTriggerCreate,
		Read:   resourceFunctionGraphTriggerRead,
		Update: resourceFunctionGraphTriggerUpdate,
		Delete: resourceFunctionGraphTriggerDelete,

		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"function_urn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					timingTrigger, obsTrigger, smnTrigger, disTrigger, kafkaTrigger,
				}, false),
			},
			// SMN trigger does not support status.
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					statusActive, statusDisabled,
				}, false),
			},
			"timer": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				MaxItems:     1,
				Elem:         timerSchemaResource(),
				ExactlyOneOf: []string{"obs", "smn", "dis", "kafka"},
			},
			"obs": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     obsSchemaResource(),
			},
			"smn": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     smnSchemaResource(),
			},
			"dis": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     disSchemaResource(),
			},
			"kafka": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     kafkaSchemaResource(),
			},
		},
	}
}

func timerSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: common.StandardVerifyWithHyphensAndStart(1, 64),
			},
			"schedule_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Rate", "Cron",
				}, false),
			},
			"schedule": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"additional_information": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func obsSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"events": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						obsEventCreated, obsEventPut, obsEventPost, obsEventCopy, obsEventMultiUpload, obsEventRemoved,
						obsEventDelete, obsEventDeleteWithoutMarker,
					}, false),
				},
				Set: schema.HashString,
			},
			"event_notification_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"suffix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func smnSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"topic_urn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func disSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"stream_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"starting_position": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TRIM_HORIZON", "LATEST",
				}, false),
			},
			"max_fetch_bytes": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1024, 4194304),
			},
			"pull_period": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(2, 60000),
			},
			"serial_enable": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func kafkaSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"batch_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      100,
				ValidateFunc: validation.IntBetween(1, 1000),
			},
			"topic_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}
func buildTimingEventData(d *schema.ResourceData) map[string]interface{} {
	event := make(map[string]interface{})

	event["name"] = d.Get("timer.0.name").(string)
	event["schedule"] = d.Get("timer.0.schedule").(string)
	event["schedule_type"] = d.Get("timer.0.schedule_type").(string)
	event["user_event"] = d.Get("timer.0.additional_information").(string)

	return event
}

func buildObsEventName(s *schema.Set) []string {
	result := make([]string, s.Len())
	for i, val := range s.List() {
		switch val.(string) {
		case obsEventPut, obsEventPost, obsEventCopy, obsEventMultiUpload:
			// The obs create events are format as 's3:ObjectCreated:{event}'
			// The events of 'ObjectCreated' are 'Put', 'Post', 'Copy' and 'CompleteMultipartUpload'.
			result[i] = fmt.Sprintf("s3:%s:%s", obsEventCreated, val.(string))
		case obsEventDelete, obsEventDeleteWithoutMarker:
			// The obs remove events are format as 's3:ObjectRemoved:{event}'
			// The events of 'ObjectRemoved' are 'Delete' and 'DeleteMarkerCreated'.
			result[i] = fmt.Sprintf("s3:%s:%s", obsEventRemoved, val.(string))
		default:
			// The obs events are format as 's3:ObjectCreated:*' or 's3:ObjectRemoved:*'
			result[i] = fmt.Sprintf("s3:%s:*", val.(string))
		}
	}
	return result
}

func buildObsEventData(d *schema.ResourceData) map[string]interface{} {
	event := make(map[string]interface{})

	obsEvents := d.Get("obs.0.events").(*schema.Set)
	event["bucket"] = d.Get("obs.0.bucket_name").(string)
	event["events"] = buildObsEventName(obsEvents)
	event["name"] = d.Get("obs.0.event_notification_name").(string)
	if prefix, ok := d.GetOk("obs.0.prefix"); ok {
		event["prefix"] = prefix.(string)
	}
	if suffix, ok := d.GetOk("obs.0.suffix"); ok {
		event["suffix"] = suffix.(string)
	}

	return event
}

func buildSmnEventData(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"topic_urn": d.Get("smn.0.topic_urn").(string),
	}
}

func buildDisEventData(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"stream_name":        d.Get("dis.0.stream_name").(string),
		"sharditerator_type": d.Get("dis.0.starting_position").(string),
		"max_fetch_bytes":    d.Get("dis.0.max_fetch_bytes").(int),
		"polling_interval":   d.Get("dis.0.pull_period").(int),
		"polling_unit":       "ms",
		"is_serial":          strconv.FormatBool(d.Get("dis.0.serial_enable").(bool)),
	}
}

func buildKafkaEventData(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"instance_id": d.Get("kafka.0.instance_id").(string),
		"batch_size":  d.Get("kafka.0.batch_size").(int),
		"topic_ids":   utils.ExpandToStringListBySet(d.Get("kafka.0.topic_ids").(*schema.Set)),
	}
}

func buildFunctionGraphTriggerParameters(d *schema.ResourceData, config *config.Config) (trigger.CreateOpts, error) {
	triggerType := d.Get("type").(string)

	opts := trigger.CreateOpts{
		TriggerTypeCode: triggerType,
		TriggerStatus:   d.Get("status").(string),
		EventTypeCode:   "MessageCreated",
	}

	switch triggerType {
	case timingTrigger:
		opts.EventData = buildTimingEventData(d)
	case obsTrigger:
		opts.EventData = buildObsEventData(d)
	case smnTrigger:
		opts.EventData = buildSmnEventData(d)
	case disTrigger:
		opts.EventData = buildDisEventData(d)
	case kafkaTrigger:
		opts.EventData = buildKafkaEventData(d)
	default:
		return opts, fmtp.Errorf("Currently, trigger type only support 'TIMER', 'OBS', 'SMN', 'DIS' and 'KAFKA'.")
	}
	return opts, nil
}

func resourceFunctionGraphTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.FgsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud FunctionGraph v2 client: %s", err)
	}

	opts, err := buildFunctionGraphTriggerParameters(d, config)
	if err != nil {
		return fmtp.Errorf("Error building create options of FunctionGraph: %s", err)
	}
	logp.Printf("[DEBUG] The create options is: %#v", opts)
	urn := d.Get("function_urn").(string)
	resp, err := trigger.Create(client, opts, urn).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud FunctionGraph trigger for function (%s): %s", urn, err)
	}
	d.SetId(resp.TriggerId)

	if resp.TriggerTypeCode == kafkaTrigger {
		// The defualt status of terraform DMS kafka trigger is 'ACTIVE'.
		if d.Get("status").(string) == "" {
			d.Set("status", statusActive)
		}
		// After creation, the status is 'DISABLED'. If we want an 'ACTIVE' kafka trigger, needs to update status.
		// Only the DMS kafka trigger cannot enter the target state immediately.
		if resp.Status != d.Get("status").(string) {
			if err = resourceFunctionGraphTriggerUpdate(d, meta); err != nil {
				return err
			}
		}
	}

	return resourceFunctionGraphTriggerRead(d, meta)
}

func setTimerEventData(d *schema.ResourceData, eventData map[string]interface{}) error {
	var info string
	if val, ok := eventData["additional_information"]; ok {
		info = val.(string)
	}
	result := []map[string]interface{}{
		{
			"name":                   eventData["name"],
			"schedule":               eventData["schedule"],
			"schedule_type":          eventData["schedule_type"],
			"additional_information": info,
		},
	}
	return d.Set("timer", result)
}

func makeObsEventNamesByResponse(events []interface{}) ([]interface{}, error) {
	result := make([]interface{}, len(events))
	regex := regexp.MustCompile(`^s3:(.*):(.*)$`)
	for i, val := range events {
		obsEvents := regex.FindAllStringSubmatch(val.(string), -1)
		if len(obsEvents) == 0 || len(obsEvents[0]) < 3 {
			return result, fmtp.Errorf("Wrong OBS event: %s", val)
		}
		// The obs events are format as 's3:ObjectCreated:*' or 's3:ObjectCreated:{event}'
		// The events of 'ObjectCreated' are 'Put', 'Post', 'Copy' and 'CompleteMultipartUpload'.
		// The events of 'ObjectRemoved' are 'Delete' and 'DeleteMarkerCreated'.
		if obsEvents[0][2] == "*" {
			result[i] = obsEvents[0][1] // 's3:{event}:*'
		} else {
			result[i] = obsEvents[0][2] // 's3:ObjectCreated:{event}' or 's3:ObjectRemoved:{event}'
		}
	}
	return result, nil
}

func setObsEventData(d *schema.ResourceData, resp *trigger.Trigger) error {
	eventData := resp.EventData
	events, err := makeObsEventNamesByResponse(eventData["events"].([]interface{}))
	if err != nil {
		return err
	}
	result := []map[string]interface{}{
		{
			"bucket_name": eventData["bucket"],
			"events":      events,
			"prefix":      eventData["prefix"],
			"suffix":      eventData["suffix"],
			// The value of event_notification_name is set in the trigger_id in response body.
			"event_notification_name": resp.TriggerId,
		},
	}
	return d.Set("obs", result)
}

func setSmnEventData(d *schema.ResourceData, eventData map[string]interface{}) error {
	result := []map[string]interface{}{
		{
			"topic_urn": eventData["topic_urn"],
		},
	}
	return d.Set("smn", result)
}

func setDisEventData(d *schema.ResourceData, eventData map[string]interface{}) error {
	isEnabled, err := strconv.ParseBool(eventData["is_serial"].(string))
	if err != nil {
		return err
	}
	result := []map[string]interface{}{
		{
			"stream_name":       eventData["stream_name"],
			"starting_position": eventData["sharditerator_type"],
			"max_fetch_bytes":   eventData["max_fetch_bytes"],
			"pull_period":       eventData["polling_interval"],
			"serial_enable":     isEnabled,
		},
	}
	return d.Set("dis", result)
}

func setKafkaEventData(d *schema.ResourceData, eventData map[string]interface{}) error {
	result := []map[string]interface{}{
		{
			"instance_id": eventData["instance_id"],
			"topic_ids":   eventData["topic_ids"],
			"batch_size":  eventData["batch_size"],
		},
	}
	return d.Set("kafka", result)
}

func setTriggerEventData(d *schema.ResourceData, resp *trigger.Trigger) error {
	switch resp.TriggerTypeCode {
	case timingTrigger:
		return setTimerEventData(d, resp.EventData)
	case obsTrigger:
		return setObsEventData(d, resp)
	case smnTrigger:
		return setSmnEventData(d, resp.EventData)
	case disTrigger:
		return setDisEventData(d, resp.EventData)
	case kafkaTrigger:
		return setKafkaEventData(d, resp.EventData)
	}
	return fmtp.Errorf("The type of trigger currently only support 'TIMER', 'OBS', 'SMN', 'DIS' and 'KAFKA'")
}

func setTriggerParamters(d *schema.ResourceData, resp *trigger.Trigger) error {
	mErr := multierror.Append(nil,
		d.Set("type", resp.TriggerTypeCode),
		d.Set("status", resp.Status),
		setTriggerEventData(d, resp),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func resourceFunctionGraphTriggerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.FgsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud FunctionGraph v2 client: %s", err)
	}

	urn := d.Get("function_urn").(string)
	pages, err := trigger.List(client, urn).AllPages()
	if err != nil {
		return fmtp.Errorf("Error retrieving FunctionGraph trigger: %s", err)
	}
	triggerList, err := trigger.ExtractList(pages)
	if len(triggerList) > 0 {
		for _, v := range triggerList {
			if v.TriggerId != d.Id() {
				continue
			}
			mErr := multierror.Append(nil,
				d.Set("region", config.GetRegion(d)),
				setTriggerParamters(d, &v),
			)
			if mErr.ErrorOrNil() != nil {
				return mErr
			}
			return nil
		}
	}
	return fmtp.Errorf("Unable to find the FunctionGraph trigger (%s) form function (%s): %s", d.Id(), urn, err)
}

func resourceFunctionGraphTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.FgsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud FunctionGraph v2 client: %s", err)
	}

	urn := d.Get("function_urn").(string)
	triggerType := d.Get("type").(string)
	targetStatus := d.Get("status").(string)
	opts := trigger.UpdateOpts{
		TriggerStatus: targetStatus,
	}
	err = trigger.Update(client, opts, urn, triggerType, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Updating HuaweiCloud FunctionGraph trigger failed: %s", err)
	}
	// After request send, check the cluster state and wait for it become running.
	stateConf := &resource.StateChangeConf{
		Target:       []string{targetStatus},
		Refresh:      triggerV2StateRefreshFunc(client, urn, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		//the system will recyle the cluster when creating failed
		return err
	}
	return resourceFunctionGraphTriggerRead(d, meta)
}

func triggerV2StateRefreshFunc(client *golangsdk.ServiceClient, urn, triggerId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pages, err := trigger.List(client, urn).AllPages()
		if err != nil {
			return nil, "DELETED", fmtp.Errorf("Error retrieving FunctionGraph trigger: %s", err)
		}
		triggerList, err := trigger.ExtractList(pages)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return triggerList, "DELETED", nil
			}
			return nil, "", err
		}
		if len(triggerList) == 0 {
			return nil, "DELETED", fmtp.Errorf("Unable to find the FunctionGraph trigger (%s) form function (%s): %s",
				triggerId, urn, err)
		}
		for _, v := range triggerList {
			if v.TriggerId == triggerId {
				return triggerList, v.Status, nil
			}
		}
		return triggerList, "DELETED", nil
	}
}

func resourceFunctionGraphTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.FgsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud FunctionGraph v2 client: %s", err)
	}

	urn := d.Get("function_urn").(string)
	triggerType := d.Get("type").(string)
	err = trigger.Delete(client, urn, triggerType, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud FunctionGraph trigger (%s) from the function (%s): %s",
			d.Id(), urn, err)
	}
	d.SetId("")
	return nil
}
