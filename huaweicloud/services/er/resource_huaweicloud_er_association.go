package er

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var associationNonUpdatableParams = []string{
	"instance_id",
	"route_table_id",
	"attachment_id",
}

// @API ER POST /v3/{project_id}/enterprise-router/{er_id}/route-tables/{route_table_id}/associate
// @API ER GET /v3/{project_id}/enterprise-router/{er_id}/route-tables/{route_table_id}/associations
// @API ER POST /v3/{project_id}/enterprise-router/{er_id}/route-tables/{route_table_id}/associations/{association_id}/change-route-policy
// @API ER POST /v3/{project_id}/enterprise-router/{er_id}/route-tables/{route_table_id}/disassociate
func ResourceAssociation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssociationCreate,
		ReadContext:   resourceAssociationRead,
		UpdateContext: resourceAssociationUpdate,
		DeleteContext: resourceAssociationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAssociationImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(associationNonUpdatableParams),

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
				Description: `The ID of the route table to which the association belongs.`,
			},
			"attachment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the attachment corresponding to the association.`,
			},

			// Optional parameters.
			"route_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"export_policy_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The export route policy ID.`,
						},
					},
				},
				Description: `The export route policy configuration.`,
			},

			// Attributes.
			"attachment_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the attachment corresponding to the association.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the association.`,
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

func buildAssociationRoutePolicy(routePolicies []interface{}) map[string]interface{} {
	if len(routePolicies) < 1 {
		return nil
	}

	routePolicy := routePolicies[0]
	return map[string]interface{}{
		"export_policy_id": utils.ValueIgnoreEmpty(utils.PathSearch("export_policy_id", routePolicy, nil)),
	}
}

func buildCreateAssociationBodyParams(attachmentId string, routePolicies []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"attachment_id": attachmentId,
		"route_policy":  buildAssociationRoutePolicy(routePolicies),
	}
}

func createAssociation(client *golangsdk.ServiceClient, instanceId, attachmentId, routeTableId string,
	routePolicies []interface{}) (interface{}, error) {
	httpUrl := "v3/{project_id}/enterprise-router/{er_id}/route-tables/{route_table_id}/associate"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{er_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{route_table_id}", routeTableId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildCreateAssociationBodyParams(attachmentId, routePolicies)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func listAssociations(client *golangsdk.ServiceClient, instanceId, routeTableId string) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/enterprise-router/{er_id}/route-tables/{route_table_id}/associations?limit={limit}"
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

		associations := utils.PathSearch("associations", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, associations...)
		nextMarker := utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" || nextMarker == marker {
			break
		}
		marker = nextMarker
	}

	return result, nil
}

func GetAssociationById(client *golangsdk.ServiceClient, instanceId, routeTableId, associationId string) (interface{}, error) {
	associations, err := listAssociations(client, instanceId, routeTableId)
	if err != nil {
		return nil, err
	}

	association := utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", associationId), associations, nil)
	if association == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte(fmt.Sprintf("the association (%s) does not exist", associationId)),
			},
		}
	}

	return association, nil
}

func associationStatusRefreshFunc(client *golangsdk.ServiceClient, instanceId, routeTableId, associationId string,
	targets []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetAssociationById(client, instanceId, routeTableId, associationId)
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

func resourceAssociationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		instanceId    = d.Get("instance_id").(string)
		attachmentId  = d.Get("attachment_id").(string)
		routeTableId  = d.Get("route_table_id").(string)
		routePolicies = d.Get("route_policy").([]interface{})
	)

	client, err := cfg.NewServiceClient("er", region)
	if err != nil {
		return diag.Errorf("error creating ER client: %s", err)
	}

	respBody, err := createAssociation(client, instanceId, attachmentId, routeTableId, routePolicies)
	if err != nil {
		return diag.Errorf("error creating the association to the route table: %s", err)
	}

	associationId := utils.PathSearch("association.id", respBody, "").(string)
	if associationId == "" {
		return diag.Errorf("unable to find the association ID from the API response")
	}
	d.SetId(associationId)

	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: associationStatusRefreshFunc(client, instanceId, routeTableId, d.Id(), []string{"available"}),
		Timeout: d.Timeout(schema.TimeoutCreate),
		// After the creation request is sent, it will briefly enter the pending status.
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the association (%s) to become available: %s", d.Id(), err)
	}

	return resourceAssociationRead(ctx, d, meta)
}

func flattenAssociationRoutePolicy(routePolicy map[string]interface{}) []map[string]interface{} {
	if len(routePolicy) < 1 {
		return nil
	}

	exportPolicyId := utils.PathSearch("export_policy_id", routePolicy, nil)
	if exportPolicyId == nil || exportPolicyId == "" {
		return nil
	}

	return []map[string]interface{}{
		{
			"export_policy_id": exportPolicyId,
		},
	}
}

func resourceAssociationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		instanceId    = d.Get("instance_id").(string)
		routeTableId  = d.Get("route_table_id").(string)
		associationId = d.Id()
	)

	client, err := cfg.NewServiceClient("er", region)
	if err != nil {
		return diag.Errorf("error creating ER client: %s", err)
	}

	respBody, err := GetAssociationById(client, instanceId, routeTableId, associationId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ER association")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("route_table_id", utils.PathSearch("route_table_id", respBody, nil)),
		d.Set("attachment_id", utils.PathSearch("attachment_id", respBody, nil)),
		d.Set("attachment_type", utils.PathSearch("resource_type", respBody, nil)),
		d.Set("route_policy", flattenAssociationRoutePolicy(utils.PathSearch("route_policy", respBody,
			make(map[string]interface{})).(map[string]interface{}))),
		d.Set("status", utils.PathSearch("state", respBody, nil)),
		// The time results are not the time in RF3339 format without milliseconds.
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
			utils.PathSearch("created_at", respBody, "").(string))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
			utils.PathSearch("updated_at", respBody, "").(string))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateAssociationRoutePolicyBodyParams(routePolicies []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"route_policy": buildAssociationRoutePolicy(routePolicies),
	}
}

func changeAssociationRoutePolicy(client *golangsdk.ServiceClient, instanceId, routeTableId, associationId string,
	routePolicies []interface{}) (interface{}, error) {
	httpUrl := "v3/{project_id}/enterprise-router/{er_id}/route-tables/{route_table_id}/associations/{association_id}/change-route-policy"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{er_id}", instanceId)
	updatePath = strings.ReplaceAll(updatePath, "{route_table_id}", routeTableId)
	updatePath = strings.ReplaceAll(updatePath, "{association_id}", associationId)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildUpdateAssociationRoutePolicyBodyParams(routePolicies)),
	}

	requestResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceAssociationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		instanceId    = d.Get("instance_id").(string)
		routeTableId  = d.Get("route_table_id").(string)
		associationId = d.Id()
		routePolicies = d.Get("route_policy").([]interface{})
	)

	client, err := cfg.NewServiceClient("er", region)
	if err != nil {
		return diag.Errorf("error creating ER client: %s", err)
	}

	_, err = changeAssociationRoutePolicy(client, instanceId, routeTableId, associationId, routePolicies)
	if err != nil {
		return diag.Errorf("error updating the route policy of association (%s): %s", associationId, err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      associationStatusRefreshFunc(client, instanceId, routeTableId, associationId, []string{"available"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the association (%s) route policy update to complete: %s", associationId, err)
	}

	return resourceAssociationRead(ctx, d, meta)
}

func deleteAssociation(client *golangsdk.ServiceClient, instanceId, attachmentId, routeTableId string) error {
	httpUrl := "v3/{project_id}/enterprise-router/{er_id}/route-tables/{route_table_id}/disassociate"
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

func resourceAssociationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		instanceId    = d.Get("instance_id").(string)
		attachmentId  = d.Get("attachment_id").(string)
		routeTableId  = d.Get("route_table_id").(string)
		associationId = d.Id()
	)

	client, err := cfg.NewServiceClient("er", region)
	if err != nil {
		return diag.Errorf("error creating ER client: %s", err)
	}

	err = deleteAssociation(client, instanceId, attachmentId, routeTableId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting association (%s)", associationId))
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: associationStatusRefreshFunc(client, instanceId, routeTableId, associationId, nil),
		Timeout: d.Timeout(schema.TimeoutDelete),
		// After the deletion request is sent, it will briefly enter the pending status.
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the association (%s) to be deleted: %s", associationId, err)
	}

	return nil
}

func resourceAssociationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
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
