package er

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/er/v3/routes"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ER POST /v3/{project_id}/enterprise-router/route-tables/{route_table_id}/static-routes
// @API ER GET /v3/{project_id}/enterprise-router/route-tables/{route_table_id}/static-routes/{route_id}
// @API ER PUT /v3/{project_id}/enterprise-router/route-tables/{route_table_id}/static-routes/{route_id}
// @API ER DELETE /v3/{project_id}/enterprise-router/route-tables/{route_table_id}/static-routes/{route_id}
func ResourceStaticRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStaticRouteCreate,
		UpdateContext: resourceStaticRouteUpdate,
		ReadContext:   resourceStaticRouteRead,
		DeleteContext: resourceStaticRouteDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceStaticRouteImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the static route and related route table are located.`,
			},
			"route_table_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the route table to which the static route belongs.`,
			},
			"destination": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The destination of the static route.`,
			},
			"attachment_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the corresponding attachment.`,
			},
			"is_blackhole": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Whether route is the black hole route.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the static route.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the static route.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the static route.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the static route.`,
			},
		},
	}
}

func buildStaticRouteCreateOpts(d *schema.ResourceData) routes.CreateOpts {
	return routes.CreateOpts{
		Destination:  d.Get("destination").(string),
		AttachmentId: d.Get("attachment_id").(string),
		IsBlackHole:  utils.Bool(d.Get("is_blackhole").(bool)),
	}
}

func staticRouteStatusRefreshFunc(client *golangsdk.ServiceClient, routeTableId, staticRouteId string,
	targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := routes.Get(client, routeTableId, staticRouteId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return "NOT_FOUND", "COMPLETED", nil
			}
			return "AN_ERROR_OCCURRED", "ERROR", err
		}

		if utils.IsStrContainsSliceElement(resp.Status, targets, false, true) {
			return resp, "COMPLETED", nil
		}
		return resp, "PENDING", nil
	}
}

func resourceStaticRouteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ErV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	var (
		routeTableId = d.Get("route_table_id").(string)
		opts         = buildStaticRouteCreateOpts(d)
	)
	resp, err := routes.Create(client, routeTableId, opts)
	if err != nil {
		return diag.Errorf("error creating static route: %s", err)
	}
	d.SetId(resp.ID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      staticRouteStatusRefreshFunc(client, routeTableId, d.Id(), []string{"available"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the create operation completed: %s", err)
	}

	return resourceStaticRouteRead(ctx, d, meta)
}

func resourceStaticRouteRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	var (
		routeTableId  = d.Get("route_table_id").(string)
		staticRouteId = d.Id()
	)

	resp, err := routes.Get(client, routeTableId, staticRouteId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "ER static route")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("destination", resp.Destination),
		d.Set("is_blackhole", resp.IsBlackHole),
		// Attributes
		d.Set("type", resp.Type),
		d.Set("status", resp.Status),
		// The time results are not the time in RF3339 format without milliseconds.
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(resp.CreatedAt)/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(resp.UpdatedAt)/1000, false)),
	)

	if len(resp.Attachments) > 0 && resp.Attachments[0].AttachmentId != "" {
		// If the static route is not a black hole route, set related VPC attachment ID.
		mErr = multierror.Append(mErr, d.Set("attachment_id", resp.Attachments[0].AttachmentId))
	} else {
		// Override 'attachment_id' while static route is the black hole route.
		mErr = multierror.Append(mErr, d.Set("attachment_id", nil))
	}

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving static route (%s) fields: %s", staticRouteId, mErr)
	}
	return nil
}

func resourceStaticRouteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	var (
		routeTableId  = d.Get("route_table_id").(string)
		staticRouteId = d.Id()
		opts          = routes.UpdateOpts{
			AttachmentId: d.Get("attachment_id").(string),
			IsBlackHole:  utils.Bool(d.Get("is_blackhole").(bool)),
		}
	)
	_, err = routes.Update(client, routeTableId, staticRouteId, opts)
	if err != nil {
		return diag.Errorf("error updating static route (%s): %s", staticRouteId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      staticRouteStatusRefreshFunc(client, routeTableId, d.Id(), []string{"available"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the update operation completed: %s", err)
	}

	return resourceStaticRouteRead(ctx, d, meta)
}

func resourceStaticRouteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	var (
		routeTableId  = d.Get("route_table_id").(string)
		staticRouteId = d.Id()
	)
	err = routes.Delete(client, routeTableId, staticRouteId)
	if err != nil {
		return diag.Errorf("error deleting static route (%s): %s", staticRouteId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      staticRouteStatusRefreshFunc(client, routeTableId, d.Id(), nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the delete operation completed: %s", err)
	}
	return nil
}

func resourceStaticRouteImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid format for import ID, want '<route_table_id>/<id>', but got '%s'", d.Id())
	}

	d.SetId(parts[1])
	if err := d.Set("route_table_id", parts[0]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}
