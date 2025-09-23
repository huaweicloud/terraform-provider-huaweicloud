package eg

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var eventBatchNonUpdatableParams = []string{"channel_id", "events"}

// @API EG POST /v1/{project_id}/channels/{channel_id}/events
func ResourceEventBatchAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventBatchActionCreate,
		ReadContext:   resourceEventBatchActionRead,
		UpdateContext: resourceEventBatchActionUpdate,
		DeleteContext: resourceEventBatchActionDelete,

		CustomizeDiff: config.FlexibleForceNew(eventBatchNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the event channel is located.`,
			},
			"channel_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the event channel.`,
			},
			"events": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        eventBatchActionEventSchema(),
				Description: `The list of events to be published.`,
			},
			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func eventBatchActionEventSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the event.`,
			},
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the event source.`,
			},
			"spec_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The CloudEvents protocol version.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the event.`,
			},
			"data_content_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The content type of the event data.`,
			},
			"data_schema": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The URI of the event data schema.`,
			},
			"data": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The payload content of the event, in JSON format.`,
			},
			"time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The time when the event occurred, in UTC format.`,
			},
			"subject": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The subject of the event.`,
			},
		},
	}
}

func buildEventBatchActionEvents(events []interface{}) []map[string]interface{} {
	if len(events) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(events))
	for _, event := range events {
		eventMap := map[string]interface{}{
			"id":              utils.PathSearch("id", event, nil),
			"source":          utils.PathSearch("source", event, nil),
			"specversion":     utils.PathSearch("spec_version", event, nil),
			"type":            utils.PathSearch("type", event, nil),
			"datacontenttype": utils.ValueIgnoreEmpty(utils.PathSearch("data_content_type", event, nil)),
			"data":            utils.StringToJson(utils.PathSearch("data", event, nil).(string)),
			"time":            utils.ValueIgnoreEmpty(utils.PathSearch("time", event, nil)),
		}
		if v := utils.PathSearch("data_schema", event, ""); v != "" {
			eventMap["data_schema"] = v
		}
		if v := utils.PathSearch("subject", event, ""); v != "" {
			eventMap["subject"] = v
		}
		result = append(result, eventMap)
	}

	return result
}

func buildEventBatchActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"events": buildEventBatchActionEvents(d.Get("events").([]interface{})),
	}
}

func checkEventBatchActionResult(respBody interface{}) error {
	responseEvents := utils.PathSearch("events", respBody, make([]interface{}, 0)).([]interface{})

	errorMessages := make([]string, 0)
	for _, respEvent := range responseEvents {
		if errorCode := utils.PathSearch("error_code", respEvent, "").(string); errorCode != "" {
			errorMsg := utils.PathSearch("error_msg", respEvent, "").(string)
			eventID := utils.PathSearch("event_id", respEvent, "").(string)
			errorMessages = append(errorMessages, fmt.Sprintf("Event %s failed: %s - %s", eventID, errorCode, errorMsg))
		}
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("failed to publish events: %s", strings.Join(errorMessages, "; "))
	}
	return nil
}

func resourceEventBatchActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/channels/{channel_id}/events"
	)

	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{channel_id}", d.Get("channel_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildEventBatchActionBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating EG event batch action: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	err = checkEventBatchActionResult(respBody)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceEventBatchActionRead(ctx, d, meta)
}

func resourceEventBatchActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEventBatchActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEventBatchActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for dispatch events to channel. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
