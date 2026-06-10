package er

import (
	"context"
	"fmt"
	"strconv"
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

var propagationNonUpdatableParams = []string{
	"instance_id",
	"route_table_id",
	"attachment_id",
}

// @API ER POST /v3/{project_id}/enterprise-router/{er_id}/route-tables/{route_table_id}/enable-propagations
// @API ER GET /v3/{project_id}/enterprise-router/{er_id}/route-tables/{route_table_id}/propagations
// @API ER POST /v3/{project_id}/enterprise-router/{er_id}/route-tables/{route_table_id}/disable-propagations
func ResourcePropagation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePropagationCreate,
		ReadContext:   resourcePropagationRead,
		UpdateContext: resourcePropagationUpdate,
		DeleteContext: resourcePropagationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePropagationImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(propagationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the ER instance and route table are located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the ER instance to which the route table and the attachment belongs.`,
			},
			"route_table_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the route table to which the propagation belongs.`,
			},
			"attachment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the attachment corresponding to the propagation.`,
			},

			// Attributes.
			"attachment_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the attachment corresponding to the propagation.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the propagation.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time.`,
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

func createPropagation(client *golangsdk.ServiceClient, instanceId, attachmentId, routeTableId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/enterprise-router/{er_id}/route-tables/{route_table_id}/enable-propagations"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{er_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{route_table_id}", routeTableId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"attachment_id": attachmentId,
		},
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func listPropagations(client *golangsdk.ServiceClient, instanceId, routeTableId string) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/enterprise-router/{er_id}/route-tables/{route_table_id}/propagations?limit={limit}"
		limit   = 200
		marker  string
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{er_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{route_table_id}", routeTableId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		requestPath := listPath
		if marker != "" {
			requestPath += fmt.Sprintf("&marker=%s", marker)
		}

		requestResp, err := client.Request("GET", requestPath, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		propagations := utils.PathSearch("propagations", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, propagations...)
		nextMarker := utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" || nextMarker == marker {
			break
		}
		marker = nextMarker
	}

	return result, nil
}

// GetPropagationById queries propagation details from a specified route table using given parameters.
func GetPropagationById(client *golangsdk.ServiceClient, instanceId, routeTableId, propagationId string) (interface{}, error) {
	propagations, err := listPropagations(client, instanceId, routeTableId)
	if err != nil {
		return nil, err
	}

	propagation := utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", propagationId), propagations, nil)
	if propagation == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte(fmt.Sprintf("the propagation (%s) does not exist", propagationId)),
			},
		}
	}

	return propagation, nil
}

func propagationStatusRefreshFunc(client *golangsdk.ServiceClient, instanceId, routeTableId, propagationId string,
	targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetPropagationById(client, instanceId, routeTableId, propagationId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return "not_found", "COMPLETED", nil
			}

			return respBody, "ERROR", err
		}

		statusResp := utils.PathSearch("state", respBody, "").(string)
		if utils.StrSliceContains([]string{"failed"}, statusResp) {
			return respBody, "ERROR", fmt.Errorf("unexpect status (%s)", statusResp)
		}
		if utils.StrSliceContains(targets, statusResp) {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}

func resourcePropagationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		instanceId   = d.Get("instance_id").(string)
		attachmentId = d.Get("attachment_id").(string)
		routeTableId = d.Get("route_table_id").(string)
	)

	client, err := cfg.NewServiceClient("er", region)
	if err != nil {
		return diag.Errorf("error creating ER client: %s", err)
	}

	respBody, err := createPropagation(client, instanceId, attachmentId, routeTableId)
	if err != nil {
		return diag.Errorf("error creating the propagation to the route table: %s", err)
	}

	propagationId := utils.PathSearch("propagation.id", respBody, "").(string)
	if propagationId == "" {
		return diag.Errorf("unable to find the propagation ID from the API response")
	}
	d.SetId(propagationId)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: propagationStatusRefreshFunc(client, instanceId, routeTableId, d.Id(), []string{"available"}),
		Timeout: d.Timeout(schema.TimeoutCreate),
		// After the creation request is sent, it will briefly enter the pending status.
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the propagation (%s) to become available: %s", d.Id(), err)
	}

	return resourcePropagationRead(ctx, d, meta)
}

func resourcePropagationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		instanceId    = d.Get("instance_id").(string)
		routeTableId  = d.Get("route_table_id").(string)
		propagationId = d.Id()
	)

	client, err := cfg.NewServiceClient("er", region)
	if err != nil {
		return diag.Errorf("error creating ER client: %s", err)
	}

	respBody, err := GetPropagationById(client, instanceId, routeTableId, propagationId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ER propagation")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("route_table_id", utils.PathSearch("route_table_id", respBody, nil)),
		d.Set("attachment_id", utils.PathSearch("attachment_id", respBody, nil)),
		d.Set("attachment_type", utils.PathSearch("resource_type", respBody, nil)),
		d.Set("status", utils.PathSearch("state", respBody, nil)),
		// The time results are not the time in RF3339 format without milliseconds.
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
			utils.PathSearch("created_at", respBody, "").(string))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
			utils.PathSearch("updated_at", respBody, "").(string))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePropagationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func deletePropagation(client *golangsdk.ServiceClient, instanceId, attachmentId, routeTableId string) error {
	httpUrl := "v3/{project_id}/enterprise-router/{er_id}/route-tables/{route_table_id}/disable-propagations"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{er_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{route_table_id}", routeTableId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"attachment_id": attachmentId,
		},
	}

	_, err := client.Request("POST", deletePath, &deleteOpt)
	return err
}

func resourcePropagationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		instanceId    = d.Get("instance_id").(string)
		attachmentId  = d.Get("attachment_id").(string)
		routeTableId  = d.Get("route_table_id").(string)
		propagationId = d.Id()
	)

	client, err := cfg.NewServiceClient("er", region)
	if err != nil {
		return diag.Errorf("error creating ER client: %s", err)
	}

	err = deletePropagation(client, instanceId, attachmentId, routeTableId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting propagation (%s)", propagationId))
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: propagationStatusRefreshFunc(client, instanceId, routeTableId, propagationId, nil),
		Timeout: d.Timeout(schema.TimeoutDelete),
		// After the deletion request is sent, it will briefly enter the pending status.
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the propagation (%s) to be deleted: %s", propagationId, err)
	}

	return nil
}

func resourcePropagationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importId := d.Id()
	parts := strings.SplitN(importId, "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format for import ID, want '<instance_id>/<route_table_id>/<id>', but got '%s'", importId)
	}

	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("route_table_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
