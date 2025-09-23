package deprecated

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	dedicatedgroups "github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apigroups"
	dedicatedenvs "github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/environments"
	"github.com/chnsz/golangsdk/openstack/apigw/shared/v1/environments"
	"github.com/chnsz/golangsdk/openstack/apigw/shared/v1/groups"
	"github.com/chnsz/golangsdk/openstack/fgs/v2/trigger"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	timingTrigger        = "TIMER"
	obsTrigger           = "OBS"
	smnTrigger           = "SMN"
	disTrigger           = "DIS"
	kafkaTrigger         = "KAFKA"
	apigTrigger          = "APIG"
	dedicatedApigTrigger = "DEDICATEDGATEWAY"
	ltsTrigger           = "LTS"

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

// @API FunctionGraph POST /v2/{project_id}/fgs/triggers/{function_urn}
// @API FunctionGraph GET /v2/{project_id}/fgs/triggers/{function_urn}
// @API FunctionGraph PUT /v2/{project_id}/fgs/triggers/{function_urn}/{trigger_type_code}/{trigger_id}
// @API FunctionGraph DELETE /v2/{project_id}/fgs/triggers/{function_urn}/{trigger_type_code}/{trigger_id}
func ResourceFunctionGraphTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFunctionGraphTriggerCreate,
		ReadContext:   resourceFunctionGraphTriggerRead,
		UpdateContext: resourceFunctionGraphTriggerUpdate,
		DeleteContext: resourceFunctionGraphTriggerDelete,

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
					timingTrigger, obsTrigger, smnTrigger, disTrigger, kafkaTrigger, apigTrigger, dedicatedApigTrigger, ltsTrigger,
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
				ExactlyOneOf: []string{"obs", "smn", "dis", "kafka", "apig", "lts"},
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
			"apig": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     apigSchemaResource(),
			},
			"lts": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     ltsSchemaResource(),
			},
		},
	}
}

func timerSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"pull_period": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
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
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ForceNew:     true,
				RequiredWith: []string{"kafka.0.user_name"},
			},
			"batch_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  100,
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

func apigSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"api_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"env_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_authentication": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"IAM", "APP", "NONE",
				}, false),
				Default: "IAM",
			},
			"request_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP", "HTTPS",
				}, false),
				Default: "HTTPS",
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  5000,
			},
		},
	}
}

func ltsSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_topic_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
		"instance_id":    d.Get("kafka.0.instance_id").(string),
		"kafka_user":     d.Get("kafka.0.user_name").(string),
		"kafka_password": d.Get("kafka.0.password").(string),
		"batch_size":     d.Get("kafka.0.batch_size").(int),
		"topic_ids":      utils.ExpandToStringListBySet(d.Get("kafka.0.topic_ids").(*schema.Set)),
	}
}

func buildLtsEventData(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"log_group_id": d.Get("lts.0.log_group_id").(string),
		"log_topic_id": d.Get("lts.0.log_topic_id").(string),
	}
}

// Obtain environment ID and sub-domain of shared APIG.
func getSharedApigSubDomainAndEnvId(d *schema.ResourceData, cfg *config.Config) (envId string, subDomain string, errorMessage error) {
	apigwClient, err := cfg.ApiGatewayV1Client(cfg.GetRegion(d))
	if err != nil {
		return envId, subDomain, fmt.Errorf("error creating shared APIG v1.0 client: %s", err)
	}

	envName := d.Get("apig.0.env_name").(string)
	// Obtain environment information.
	envOpt := environments.ListOpts{
		EnvName: envName,
	}
	envList, err := environments.List(apigwClient, envOpt)
	if err != nil {
		return envId, subDomain, fmt.Errorf("unable to obtain the environment list: %s", err)
	}
	if len(envList) == 0 {
		return envId, subDomain, fmt.Errorf("there is no environment named %s: %s", envName, err)
	}
	envId = envList[0].Id

	// Obtain group information.
	groupId := d.Get("apig.0.group_id").(string)
	groupResp, err := groups.Get(apigwClient, groupId).Extract()
	if err != nil {
		return envId, subDomain, fmt.Errorf("unable to obtain the APIG group (%s): %s", groupId, err)
	}
	subDomain = groupResp.SlDomain

	return
}

// Obtain environment ID and sub-domain of dedicated APIG.
func getDedicatedApigSubDomainAndEnvId(d *schema.ResourceData, cfg *config.Config) (envId string, subDomain string, errorMessage error) {
	apigwClient, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return envId, subDomain, fmt.Errorf("error creating dedicated APIG v2 client: %s", err)
	}

	instanceId := d.Get("apig.0.instance_id").(string)
	envName := d.Get("apig.0.env_name").(string)
	// Obtain environment information.
	envOpt := dedicatedenvs.ListOpts{
		Name: envName,
	}
	pages, err := dedicatedenvs.List(apigwClient, instanceId, envOpt).AllPages()
	if err != nil {
		return envId, subDomain, fmt.Errorf("error getting environment list: %s", err)
	}
	envList, err := dedicatedenvs.ExtractEnvironments(pages)
	if err != nil {
		return envId, subDomain, fmt.Errorf("unable to retrieve the response to list: %s", err)
	}
	if len(envList) == 0 {
		return envId, subDomain, fmt.Errorf("there is no environment named %s: %s", envName, err)
	}
	envId = envList[0].Id

	// Obtain group information.
	groupId := d.Get("apig.0.group_id").(string)
	groupResp, err := dedicatedgroups.Get(apigwClient, instanceId, groupId).Extract()
	if err != nil {
		return envId, subDomain, fmt.Errorf("unable to obtain the APIG group (%s): %s", groupId, err)
	}
	subDomain = groupResp.Subdomain

	return
}

func buildApigEventData(d *schema.ResourceData, cfg *config.Config) (map[string]interface{}, error) {
	// Common configuration
	result := map[string]interface{}{
		"env_name":     d.Get("apig.0.env_name").(string),
		"group_id":     d.Get("apig.0.group_id").(string),
		"protocol":     d.Get("apig.0.request_protocol").(string),
		"auth":         d.Get("apig.0.security_authentication").(string),
		"name":         d.Get("apig.0.api_name").(string),
		"path":         fmt.Sprintf("/%s", d.Get("apig.0.api_name").(string)), // Use API name as path.
		"backend_type": "FUNCTION",
		"match_mode":   "SWA",
		"req_method":   "ANY",
		"type":         1,
		"func_info": map[string]interface{}{
			"timeout": d.Get("apig.0.timeout").(int),
		},
	}

	var envId, subDomain string
	var err error
	// The different between the shared APIG and the dedicated APIG is whether the instance ID is set.
	if instanceId, ok := d.GetOk("apig.0.instance_id"); ok {
		result["instance_id"] = instanceId
		envId, subDomain, err = getDedicatedApigSubDomainAndEnvId(d, cfg)
		if err != nil {
			return result, err
		}
	} else {
		envId, subDomain, err = getSharedApigSubDomainAndEnvId(d, cfg)
		if err != nil {
			return result, err
		}
	}
	result["env_id"] = envId
	result["sl_domain"] = subDomain

	return result, nil
}

func buildFunctionGraphTriggerParameters(d *schema.ResourceData, cfg *config.Config) (trigger.CreateOpts, error) {
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
	case ltsTrigger:
		opts.EventData = buildLtsEventData(d)
	case apigTrigger, dedicatedApigTrigger:
		eventData, err := buildApigEventData(d, cfg)
		if err != nil {
			return opts, err
		}
		opts.EventData = eventData
	default:
		return opts, fmt.Errorf("Currently, trigger type only support 'TIMER', 'OBS', 'SMN', 'DIS', 'KAFKA', 'APIG', 'LTS' " +
			"and 'DEDICATEDGATEWAY'.")
	}
	return opts, nil
}

func resourceFunctionGraphTriggerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	opts, err := buildFunctionGraphTriggerParameters(d, cfg)
	if err != nil {
		return diag.Errorf("error building create options of FunctionGraph: %s", err)
	}
	log.Printf("[DEBUG] The create options is: %#v", opts)
	urn := d.Get("function_urn").(string)
	resp, err := trigger.Create(client, opts, urn).Extract()
	if err != nil {
		return diag.Errorf("error creating FunctionGraph trigger for function (%s): %s", urn, err)
	}
	d.SetId(resp.TriggerId)

	if resp.TriggerTypeCode == kafkaTrigger {
		// The default status of terraform DMS kafka trigger is 'ACTIVE'.
		if d.Get("status").(string) == "" {
			d.Set("status", statusActive)
		}
		// After creation, the status is 'DISABLED'. If we want an 'ACTIVE' kafka trigger, needs to update status.
		// Only the DMS kafka trigger cannot enter the target state immediately.
		if resp.Status != d.Get("status").(string) {
			err := resourceFunctionGraphTriggerUpdate(ctx, d, meta)
			if err != nil {
				return err
			}
		}
	}

	return resourceFunctionGraphTriggerRead(ctx, d, meta)
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
			return result, fmt.Errorf("wrong OBS event: %s", val)
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
			"user_name":   eventData["kafka_user"],
			"password":    eventData["kafka_password"],
			"topic_ids":   eventData["topic_ids"],
			"batch_size":  eventData["batch_size"],
		},
	}
	return d.Set("kafka", result)
}

func setLtsEventData(d *schema.ResourceData, eventData map[string]interface{}) error {
	result := []map[string]interface{}{
		{
			"log_group_id": eventData["log_group_id"],
			"log_topic_id": eventData["log_topic_id"],
		},
	}
	return d.Set("lts", result)
}

func setApigEventData(d *schema.ResourceData, eventData map[string]interface{}) error {
	result := make([]map[string]interface{}, 1)
	funcInfo := eventData["func_info"].(map[string]interface{})
	apigInfo := map[string]interface{}{
		"group_id":                eventData["group_id"],
		"api_name":                eventData["api_name"],
		"env_name":                eventData["env_name"],
		"security_authentication": eventData["auth"],
		"request_protocol":        eventData["protocol"],
		"timeout":                 funcInfo["timeout"],
	}
	if instanceId, ok := eventData["instance_id"]; ok {
		apigInfo["instance_id"] = instanceId
	}
	return d.Set("apig", result)
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
	case apigTrigger, dedicatedApigTrigger:
		return setApigEventData(d, resp.EventData)
	case ltsTrigger:
		return setLtsEventData(d, resp.EventData)
	}
	return fmt.Errorf("the type of trigger currently only support 'TIMER', 'OBS', 'SMN', 'DIS', 'KAFKA', 'APIG', 'LTS' and " +
		"'DEDICATEDGATEWAY'.")
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

func resourceFunctionGraphTriggerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	urn := d.Get("function_urn").(string)
	pages, err := trigger.List(client, urn).AllPages()
	if err != nil {
		return common.CheckDeletedDiag(d, parseRequestError(err), "error retrieving FunctionGraph trigger")
	}
	triggerList, _ := trigger.ExtractList(pages)
	if len(triggerList) > 0 {
		for _, v := range triggerList {
			if v.TriggerId != d.Id() {
				continue
			}
			v := v
			mErr := multierror.Append(nil,
				d.Set("region", cfg.GetRegion(d)),
				setTriggerParamters(d, &v),
			)
			if mErr.ErrorOrNil() != nil {
				return diag.Errorf("error setting Trigger Parameters: %s", mErr.ErrorOrNil())
			}
			return nil
		}
	}
	return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte(fmt.Sprintf("unable to find the FunctionGraph trigger (%s) from function (%s), the trigger "+
				"has been deleted", d.Id(), urn)),
		},
	}, "error retrieving FunctionGraph trigger")
}

func resourceFunctionGraphTriggerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	urn := d.Get("function_urn").(string)
	triggerType := d.Get("type").(string)
	targetStatus := d.Get("status").(string)
	opts := trigger.UpdateOpts{
		TriggerStatus: targetStatus,
	}
	err = trigger.Update(client, opts, urn, triggerType, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("updating FunctionGraph trigger failed: %s", err)
	}
	// After request send, check the cluster state and wait for it become running.
	stateConf := &resource.StateChangeConf{
		Target:       []string{targetStatus},
		Refresh:      triggerV2StateRefreshFunc(client, urn, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		// the system will recyle the cluster when creating failed
		return diag.Errorf("wait for state Context failed: %s", err)
	}
	return resourceFunctionGraphTriggerRead(ctx, d, meta)
}

func triggerV2StateRefreshFunc(client *golangsdk.ServiceClient, urn, triggerId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pages, err := trigger.List(client, urn).AllPages()
		if err != nil {
			return nil, "DELETED", fmt.Errorf("error retrieving FunctionGraph trigger: %s", err)
		}
		triggerList, err := trigger.ExtractList(pages)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return triggerList, "DELETED", nil
			}
			return nil, "", err
		}
		if len(triggerList) == 0 {
			return nil, "DELETED", fmt.Errorf("unable to find the FunctionGraph trigger (%s) form function (%s): %s",
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

func resourceFunctionGraphTriggerDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	urn := d.Get("function_urn").(string)
	triggerType := d.Get("type").(string)
	err = trigger.Delete(client, urn, triggerType, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting FunctionGraph trigger (%s) from the function (%s): %s",
			d.Id(), urn, err)
	}
	return nil
}

func parseRequestError(respErr error) error {
	var apiErr trigger.Error
	if errCode, ok := respErr.(golangsdk.ErrDefault500); ok && errCode.Body != nil {
		pErr := json.Unmarshal(errCode.Body, &apiErr)
		if pErr == nil && apiErr.Code == "FSS.0500" && apiErr.Message == "Error getting associated function" {
			return golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Body: []byte("the related function and this trigger has been deleted"),
				},
			}
		}
	}
	return respErr
}
