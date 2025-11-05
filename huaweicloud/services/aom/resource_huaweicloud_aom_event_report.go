package aom

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

var aomEventReportNonUpdatableParams = []string{"events", "action", "enterprise_project_id"}

// @API AOM PUT /v2/{project_id}/push/events
func ResourceEventReport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventReportCreate,
		ReadContext:   resourceEventReportRead,
		UpdateContext: resourceEventReportUpdate,
		DeleteContext: resourceEventReportDelete,

		CustomizeDiff: config.FlexibleForceNew(aomEventReportNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the event source are located.`,
			},

			// Required parameters.
			"events": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        eventSchema(),
				Description: `The list of events or alarms to be reported.`,
			},

			// Optional parameters.
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The action type of the request.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0",
				Description: `The enterprise project ID to which the resource belongs.`,
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

func eventSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:        schema.TypeMap,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The detail of the event or alarm.`,
			},
			"starts_at": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The time when the event or alarm occurred, in UTC milliseconds timestamp.`,
			},
			"ends_at": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The time when the event or alarm was cleared, in UTC milliseconds timestamp.`,
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The automatic clearing time for expired alarms, in milliseconds.`,
			},
			"annotations": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The additional fields of the event or alarm, in JSON format.`,
			},
		},
	}
}

func buildEventReportEvents(events []interface{}) []map[string]interface{} {
	if len(events) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(events))
	for _, event := range events {
		result = append(result, map[string]interface{}{
			"metadata":    utils.PathSearch("metadata", event, make(map[string]string)),
			"starts_at":   utils.ValueIgnoreEmpty(utils.PathSearch("starts_at", event, 0).(int)),
			"ends_at":     utils.ValueIgnoreEmpty(utils.PathSearch("ends_at", event, 0).(int)),
			"timeout":     utils.ValueIgnoreEmpty(utils.PathSearch("timeout", event, 0).(int)),
			"annotations": utils.StringToJson(utils.PathSearch("annotations", event, "").(string), nil),
		})
	}

	return result
}

func buildEventReportBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"events": buildEventReportEvents(d.Get("events").([]interface{})),
	}
}

func buildEventReportQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("action"); ok {
		res = fmt.Sprintf("%s&action=%s", res, v.(string))
	}

	if len(res) > 1 {
		return "?" + res[:1]
	}
	return ""
}

func resourceEventReportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                 = meta.(*config.Config)
		region              = cfg.GetRegion(d)
		httpUrl             = "v2/{project_id}/push/events"
		enterpriseProjectId = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath += buildEventReportQueryParams(d)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":          "application/json",
			"enterprise_project_id": enterpriseProjectId,
		},
		OkCodes: []int{
			200, 204,
		},
		JSONBody: buildEventReportBodyParams(d),
	}

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error reporting error event: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceEventReportRead(ctx, d, meta)
}

func resourceEventReportRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEventReportUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEventReportDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for report error events. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
