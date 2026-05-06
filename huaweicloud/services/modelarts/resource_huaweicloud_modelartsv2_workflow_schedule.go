package modelarts

import (
	"context"
	"fmt"
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

var (
	v2WorkflowScheduleNonUpdatableParams = []string{
		"workflow_id",
		"policies",
		"policies.*.on_failure",
		"policies.*.on_running",
	}
	v2WorkflowScheduleNotFoundErrCodes = []string{
		"ModelArts.7512", // The workflow does not exist.
		"ModelArts.7525", // The workflow schedule does not exist.
	}
)

// @API ModelArts POST /v2/{project_id}/workflows/{workflow_id}/schedules
// @API ModelArts GET /v2/{project_id}/workflows/{workflow_id}/schedules/{schedule_id}
// @API ModelArts PUT /v2/{project_id}/workflows/{workflow_id}/schedules/{schedule_id}
// @API ModelArts DELETE /v2/{project_id}/workflows/{workflow_id}/schedules/{schedule_id}
func ResourceV2WorkflowSchedule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2WorkflowScheduleCreate,
		ReadContext:   resourceV2WorkflowScheduleRead,
		UpdateContext: resourceV2WorkflowScheduleUpdate,
		DeleteContext: resourceV2WorkflowScheduleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV2WorkflowScheduleImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(v2WorkflowScheduleNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the workflow schedule is located.`,
			},

			// Required parameters.
			"workflow_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workflow to which the schedule configuration belongs.`,
			},
			"content": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  `The content of the workflow schedule, in JSON format.`,
			},
			"policies": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"on_failure": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Description: utils.SchemaDesc(
								`The policy action when the workflow execution fails.`,
								utils.SchemaDescInput{
									Computed: true,
								},
							),
						},
						"on_running": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Description: utils.SchemaDesc(
								`The policy action when the workflow is already running.`,
								utils.SchemaDescInput{
									Computed: true,
								},
							),
						},
					},
				},
				Description: utils.SchemaDesc(
					`The scheduling policies of the workflow schedule.`,
					utils.SchemaDescInput{
						Computed: true,
					},
				),
			},

			// Attributes.
			"enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the workflow schedule is enabled.`,
			},
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user ID that created the workflow schedule.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the workflow schedule, in RFC3339 format.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func buildV2WorkflowSchedulePolicies(policies []interface{}) map[string]interface{} {
	if len(policies) < 1 {
		return map[string]interface{}{
			"on_failure": "retry",
			"on_running": "cancel",
		}
	}

	return map[string]interface{}{
		"on_failure": utils.PathSearch("on_failure", policies[0], nil),
		"on_running": utils.PathSearch("on_running", policies[0], nil),
	}
}

func buildV2WorkflowScheduleCreateBodyParams(scheduleContent string, policies []interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		// Fixed values.
		"action":   "run",
		"type":     "time",
		"policies": buildV2WorkflowSchedulePolicies(policies),
		// Required parameters.
		"content": utils.StringToJson(scheduleContent),
	}
	return bodyParams
}

func createV2WorkflowSchedule(client *golangsdk.ServiceClient, workflowId, scheduleContent string, policies []interface{}) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/workflows/{workflow_id}/schedules"
	)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workflow_id}", workflowId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildV2WorkflowScheduleCreateBodyParams(scheduleContent, policies)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceV2WorkflowScheduleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		workflowId      = d.Get("workflow_id").(string)
		scheduleContent = d.Get("content").(string)
		policies        = d.Get("policies").([]interface{})
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resp, err := createV2WorkflowSchedule(client, workflowId, scheduleContent, policies)
	if err != nil {
		return diag.Errorf("error creating ModelArts workflow schedule: %s", err)
	}

	scheduleId := utils.PathSearch("uuid", resp, "").(string)
	if scheduleId == "" {
		return diag.Errorf("unable to find the ModelArts workflow schedule ID from the API response")
	}
	d.SetId(scheduleId)

	return resourceV2WorkflowScheduleRead(ctx, d, meta)
}

func flattenV2WorkflowSchedulePolicies(policies map[string]interface{}) []map[string]interface{} {
	if len(policies) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"on_failure": utils.PathSearch("on_failure", policies, nil),
			"on_running": utils.PathSearch("on_running", policies, nil),
		},
	}
}

func GetV2WorkflowScheduleById(client *golangsdk.ServiceClient, workflowId, scheduleId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/workflows/{workflow_id}/schedules/{schedule_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workflow_id}", workflowId)
	getPath = strings.ReplaceAll(getPath, "{schedule_id}", scheduleId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceV2WorkflowScheduleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		workflowId = d.Get("workflow_id").(string)
		scheduleId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resp, err := GetV2WorkflowScheduleById(client, workflowId, scheduleId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", v2WorkflowScheduleNotFoundErrCodes...),
			fmt.Sprintf("error retrieving ModelArts workflow schedule (%s)", scheduleId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("content", utils.JsonToString(utils.PathSearch("content", resp, nil))),
		d.Set("enable", utils.PathSearch("enable", resp, nil)),
		d.Set("policies", flattenV2WorkflowSchedulePolicies(utils.PathSearch("policies", resp,
			make(map[string]interface{})).(map[string]interface{}))),
		d.Set("user_id", utils.PathSearch("user_id", resp, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at",
			resp, "").(string))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildV2WorkflowScheduleUpdateBodyParams(scheduleContent string, isEnable bool) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"content": utils.StringToJson(scheduleContent),
		"enable":  isEnable,
	}
	return bodyParams
}

func updateV2WorkflowScheduleContent(client *golangsdk.ServiceClient, workflowId, scheduleId, scheduleContent string, isEnable bool) error {
	httpUrl := "v2/{project_id}/workflows/{workflow_id}/schedules/{schedule_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{workflow_id}", workflowId)
	updatePath = strings.ReplaceAll(updatePath, "{schedule_id}", scheduleId)

	opt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildV2WorkflowScheduleUpdateBodyParams(scheduleContent, isEnable)),
	}

	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func resourceV2WorkflowScheduleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		workflowId = d.Get("workflow_id").(string)
		scheduleId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	if d.HasChange("content") {
		err := updateV2WorkflowScheduleContent(client, workflowId, scheduleId, d.Get("content").(string), d.Get("enable").(bool))
		if err != nil {
			return diag.Errorf("error updating ModelArts workflow schedule content: %s", err)
		}
	}

	return resourceV2WorkflowScheduleRead(ctx, d, meta)
}

func deleteV2WorkflowSchedule(client *golangsdk.ServiceClient, workflowId, scheduleId string) error {
	httpUrl := "v2/{project_id}/workflows/{workflow_id}/schedules/{schedule_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workflow_id}", workflowId)
	deletePath = strings.ReplaceAll(deletePath, "{schedule_id}", scheduleId)

	opt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &opt)
	return err
}

func resourceV2WorkflowScheduleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		workflowId = d.Get("workflow_id").(string)
		scheduleId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = deleteV2WorkflowSchedule(client, workflowId, scheduleId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", v2WorkflowScheduleNotFoundErrCodes...),
			fmt.Sprintf("error deleting ModelArts workflow schedule (%s)", scheduleId))
	}

	return nil
}

func resourceV2WorkflowScheduleImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workflow_id>/<id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("workflow_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
