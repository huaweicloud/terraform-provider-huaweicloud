package fgs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	triggerStatusActive   = "ACTIVE"
	triggerStatusDisabled = "DISABLED"
)

// @API FunctionGraph POST /v2/{project_id}/fgs/triggers/{function_urn}
// @API FunctionGraph GET /v2/{project_id}/fgs/triggers/{function_urn}/{trigger_type_code}/{trigger_id}
// @API FunctionGraph PUT /v2/{project_id}/fgs/triggers/{function_urn}/{trigger_type_code}/{trigger_id}
// @API FunctionGraph DELETE /v2/{project_id}/fgs/triggers/{function_urn}/{trigger_type_code}/{trigger_id}
func ResourceFunctionTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFunctionTriggerCreate,
		ReadContext:   resourceFunctionTriggerRead,
		UpdateContext: resourceFunctionTriggerUpdate,
		DeleteContext: resourceFunctionTriggerDelete,

		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceFunctionTriggermportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the function trigger is located.`,
			},
			"function_urn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The function URN to which the function trigger belongs.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The type of the function trigger.`,
			},
			"event_data": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  `The detailed configuration of the function trigger event, in JSON format.`,
			},
			// INFO:
			// + Currently, only some triggers support setting the **DISABLED** value, such as `TIMER`, `DDS`, `DMS`,
			//   `KAFKA` and `LTS`.
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					triggerStatusActive, triggerStatusDisabled,
				}, false),
				Description: `The expected status of the function trigger.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the function trigger.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the function trigger.`,
			},
		},
	}
}

func buildCreateFunctionTriggerBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"trigger_type_code": d.Get("type"),
		"event_data":        utils.StringToJson(d.Get("event_data").(string)),
	}
}

func resourceFunctionTriggerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v2/{project_id}/fgs/triggers/{function_urn}"
		functionUrn = d.Get("function_urn").(string)
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{function_urn}", functionUrn)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateFunctionTriggerBodyParams(d)),
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating function trigger: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	triggerId := utils.PathSearch("trigger_id", respBody, "").(string)
	if triggerId == "" {
		return diag.Errorf("unable to find the function trigger ID from the API response")
	}
	d.SetId(triggerId)

	return resourceFunctionTriggerRead(ctx, d, meta)
}

func waitForFunctionTriggerStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: functionTriggerStatusRefreshFunc(client, d.Get("function_urn").(string), d.Get("type").(string),
			d.Id(), []string{"ACTIVE", "DISABLED"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func functionTriggerStatusRefreshFunc(client *golangsdk.ServiceClient, functionUrn, triggerType, triggerId string,
	targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetTriggerById(client, functionUrn, triggerType, triggerId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
				return "Resource Not Found", "COMPLETED", nil
			}
			return respBody, "ERROR", err
		}

		status := utils.PathSearch("trigger_status", respBody, "").(string)
		// The values of the trigger status only 'ACTIVE' and 'DISABLED', and does not include abnormal status.
		if utils.StrSliceContains(targets, status) {
			return respBody, "COMPLETED", nil
		}
		return respBody, "PENDING", nil
	}
}

func GetTriggerById(client *golangsdk.ServiceClient, functionUrn, triggerType, triggerId string) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/fgs/triggers/{function_urn}/{trigger_type_code}/{trigger_id}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{function_urn}", functionUrn)
	getPath = strings.ReplaceAll(getPath, "{trigger_type_code}", triggerType)
	getPath = strings.ReplaceAll(getPath, "{trigger_id}", triggerId)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, common.ConvertExpected500ErrInto404Err(err, "error_code", "FSS.0500")
	}
	return utils.FlattenResponse(requestResp)
}

func resourceFunctionTriggerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		functionUrn = d.Get("function_urn").(string)
		triggerType = d.Get("type").(string)
		triggerId   = d.Id()
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	respBody, err := GetTriggerById(client, functionUrn, triggerType, triggerId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying function trigger (%s)", triggerId))
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("type", utils.PathSearch("trigger_type_code", respBody, nil)),
		d.Set("status", utils.PathSearch("trigger_status", respBody, nil)),
		d.Set("event_data", utils.JsonToString(utils.PathSearch("event_data", respBody, nil))),
		d.Set("created_at", utils.PathSearch("created_time", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("last_updated_time", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateFunctionTriggerBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"trigger_status": d.Get("status"),
		"event_data":     utils.StringToJson(d.Get("event_data").(string)),
	}
}

func resourceFunctionTriggerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v2/{project_id}/fgs/triggers/{function_urn}/{trigger_type_code}/{trigger_id}"
		functionUrn = d.Get("function_urn").(string)
		triggerType = d.Get("type").(string)
		triggerId   = d.Id()
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{function_urn}", functionUrn)
	updatePath = strings.ReplaceAll(updatePath, "{trigger_type_code}", triggerType)
	updatePath = strings.ReplaceAll(updatePath, "{trigger_id}", triggerId)
	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateFunctionTriggerBodyParams(d)),
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return diag.Errorf("error updating function trigger: %s", err)
	}

	err = waitForFunctionTriggerStatusCompleted(ctx, client, d)
	if err != nil {
		diag.Errorf("error waiting for the function trigger (%s) status to become available: %s", triggerId, err)
	}
	return nil
}

func resourceFunctionTriggerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v2/{project_id}/fgs/triggers/{function_urn}/{trigger_type_code}/{trigger_id}"
		functionUrn = d.Get("function_urn").(string)
		triggerType = d.Get("type").(string)
		triggerId   = d.Id()
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{function_urn}", functionUrn)
	deletePath = strings.ReplaceAll(deletePath, "{trigger_type_code}", triggerType)
	deletePath = strings.ReplaceAll(deletePath, "{trigger_id}", triggerId)
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", "FSS.0401"), // Function not found.
			"error deleting function trigger")
	}

	err = waitForFunctionTriggerDeleted(ctx, client, d)
	if err != nil {
		diag.Errorf("error waiting for the function trigger (%s) status to become deleted: %s", triggerId, err)
	}
	return nil
}

func waitForFunctionTriggerDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      functionTriggerStatusRefreshFunc(client, d.Get("function_urn").(string), d.Get("type").(string), d.Id(), nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceFunctionTriggermportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	var (
		importId = d.Id()
		parts    = strings.Split(importId, "/")
	)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid resource ID format for function trigger, want '<function_urn>/<type>/<id>', but got '%s'", importId)
	}
	d.SetId(parts[2])
	mErr := multierror.Append(
		d.Set("function_urn", parts[0]),
		d.Set("type", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
