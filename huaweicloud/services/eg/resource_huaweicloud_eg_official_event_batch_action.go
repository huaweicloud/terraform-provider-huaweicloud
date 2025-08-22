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

var officialEventBatchNonUpdatableParams = []string{"source_name", "events"}

// @API EG POST /v1/{project_id}/official/sources/{source_name}/events
func ResourceOfficialEventBatchAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOfficialEventBatchActionCreate,
		ReadContext:   resourceOfficialEventBatchActionRead,
		UpdateContext: resourceOfficialEventBatchActionUpdate,
		DeleteContext: resourceOfficialEventBatchActionDelete,

		CustomizeDiff: config.FlexibleForceNew(officialEventBatchNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the official event source is located.`,
			},
			"source_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the official event source.`,
			},
			"events": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        officialEventBatchActionEventSchema(),
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

func officialEventBatchActionEventSchema() *schema.Resource {
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
				Description: `The name of the official event source.`,
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

func buildOfficialEventBatchActionEvents(events []interface{}) []map[string]interface{} {
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
			"dataschema":      utils.ValueIgnoreEmpty(utils.PathSearch("data_schema", event, nil)),
			"data":            utils.StringToJson(utils.PathSearch("data", event, nil).(string)),
			"time":            utils.ValueIgnoreEmpty(utils.PathSearch("time", event, nil)),
			"subject":         utils.ValueIgnoreEmpty(utils.PathSearch("subject", event, nil)),
		}
		result = append(result, eventMap)
	}

	return result
}

func buildOfficialEventBatchActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"events": buildOfficialEventBatchActionEvents(d.Get("events").([]interface{})),
	}
}

func checkOfficialEventBatchActionResult(respBody interface{}) error {
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
		return fmt.Errorf("failed to publish official events: %s", strings.Join(errorMessages, "; "))
	}
	return nil
}

func resourceOfficialEventBatchActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/official/sources/{source_name}/events"
	)

	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{source_name}", d.Get("source_name").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildOfficialEventBatchActionBodyParams(d),
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

	err = checkOfficialEventBatchActionResult(respBody)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceOfficialEventBatchActionRead(ctx, d, meta)
}

func resourceOfficialEventBatchActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOfficialEventBatchActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOfficialEventBatchActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for dispatch official events to channel. Deleting 
this resource will not clear the corresponding request record, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
