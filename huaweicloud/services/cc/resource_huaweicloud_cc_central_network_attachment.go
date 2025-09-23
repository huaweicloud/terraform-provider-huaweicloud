// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CC
// ---------------------------------------------------------------

package cc

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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CC POST /v3/{domain_id}/gcn/central-network/{central_network_id}/gdgw-attachments
// @API CC GET /v3/{domain_id}/gcn/central-network/{central_network_id}/gdgw-attachments/{gdgw_attachment_id}
// @API CC PUT /v3/{domain_id}/gcn/central-network/{central_network_id}/gdgw-attachments/{gdgw_attachment_id}
// @API CC DELETE /v3/{domain_id}/gcn/central-network/{central_network_id}/attachments/{attachment_id}
func ResourceCentralNetworkAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCentralNetworkAttachmentCreate,
		UpdateContext: resourceCentralNetworkAttachmentUpdate,
		ReadContext:   resourceCentralNetworkAttachmentRead,
		DeleteContext: resourceCentralNetworkAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCentralNetworkAttachmentImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"central_network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the attachment.`,
			},
			"enterprise_router_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The enterprise router ID.`,
			},
			"enterprise_router_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The project ID to which the enterprise router belongs.`,
			},
			"enterprise_router_region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The region ID to which the enterprise router belongs.`,
			},
			"global_dc_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The global DC gateway ID.`,
			},
			"global_dc_gateway_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The project ID to which the global DC gateway belongs.`,
			},
			"global_dc_gateway_region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The region ID to which the global DC gateway belongs.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the attachment.`,
			},
			"central_network_plane_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The central network plane ID.`,
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Central network attachment status.`,
			},
		},
	}
}

func resourceCentralNetworkAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createCentralNetworkAttachmentHttpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/gdgw-attachments"
		createCentralNetworkAttachmentProduct = "cc"
	)
	createCentralNetworkAttachmentClient, err := cfg.NewServiceClient(createCentralNetworkAttachmentProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	createCentralNetworkAttachmentPath := createCentralNetworkAttachmentClient.Endpoint + createCentralNetworkAttachmentHttpUrl
	createCentralNetworkAttachmentPath = strings.ReplaceAll(createCentralNetworkAttachmentPath, "{domain_id}", cfg.DomainID)
	createCentralNetworkAttachmentPath = strings.ReplaceAll(createCentralNetworkAttachmentPath, "{central_network_id}",
		d.Get("central_network_id").(string))

	createCentralNetworkAttachmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createCentralNetworkAttachmentOpt.JSONBody = utils.RemoveNil(buildCreateCentralNetworkAttachmentBodyParams(d))
	createCentralNetworkAttachmentResp, err := createCentralNetworkAttachmentClient.Request("POST",
		createCentralNetworkAttachmentPath, &createCentralNetworkAttachmentOpt)
	if err != nil {
		return diag.Errorf("error creating central network attachment: %s", err)
	}

	createCentralNetworkAttachmentRespBody, err := utils.FlattenResponse(createCentralNetworkAttachmentResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("central_network_gdgw_attachment.id", createCentralNetworkAttachmentRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating central network attachment: ID is not found in API response")
	}
	d.SetId(id)

	err = centralNetworkAttachmentWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the central network attachment (%s) creation to complete: %s", d.Id(), err)
	}

	return resourceCentralNetworkAttachmentRead(ctx, d, meta)
}

func buildCreateCentralNetworkAttachmentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"central_network_gdgw_attachment": map[string]interface{}{
			"name":                         d.Get("name"),
			"enterprise_router_id":         d.Get("enterprise_router_id"),
			"enterprise_router_project_id": d.Get("enterprise_router_project_id"),
			"enterprise_router_region_id":  d.Get("enterprise_router_region_id"),
			"global_dc_gateway_id":         d.Get("global_dc_gateway_id"),
			"global_dc_gateway_project_id": d.Get("global_dc_gateway_project_id"),
			"global_dc_gateway_region_id":  d.Get("global_dc_gateway_region_id"),
			"description":                  utils.ValueIgnoreEmpty(d.Get("description")),
			"central_network_plane_id":     utils.ValueIgnoreEmpty(d.Get("central_network_plane_id")),
		},
	}
	return bodyParams
}

func centralNetworkAttachmentWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			var (
				centralNetworkAttachmentWaitingHttpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/gdgw-attachments/{id}"
				centralNetworkAttachmentWaitingProduct = "cc"
			)
			centralNetworkAttachmentWaitingClient, err := cfg.NewServiceClient(centralNetworkAttachmentWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating CC client: %s", err)
			}

			centralNetworkAttachmentWaitingPath := centralNetworkAttachmentWaitingClient.Endpoint + centralNetworkAttachmentWaitingHttpUrl
			centralNetworkAttachmentWaitingPath = strings.ReplaceAll(centralNetworkAttachmentWaitingPath, "{domain_id}", cfg.DomainID)
			centralNetworkAttachmentWaitingPath = strings.ReplaceAll(centralNetworkAttachmentWaitingPath, "{central_network_id}",
				d.Get("central_network_id").(string))
			centralNetworkAttachmentWaitingPath = strings.ReplaceAll(centralNetworkAttachmentWaitingPath, "{id}", d.Id())

			centralNetworkAttachmentWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
				MoreHeaders: map[string]string{"Content-Type": "application/json"},
			}

			centralNetworkAttachmentWaitingResp, err := centralNetworkAttachmentWaitingClient.Request("GET",
				centralNetworkAttachmentWaitingPath, &centralNetworkAttachmentWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			centralNetworkAttachmentWaitingRespBody, err := utils.FlattenResponse(centralNetworkAttachmentWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`central_network_gdgw_attachment.state`, centralNetworkAttachmentWaitingRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `central_network_gdgw_attachment.state`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"AVAILABLE",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return centralNetworkAttachmentWaitingRespBody, "COMPLETED", nil
			}

			pendingStatus := []string{
				"CREATING",
				"UPDATING",
				"DELETING",
			}
			if utils.StrSliceContains(pendingStatus, status) {
				return centralNetworkAttachmentWaitingRespBody, "PENDING", nil
			}

			return centralNetworkAttachmentWaitingRespBody, status, nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceCentralNetworkAttachmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getCentralNetworkAttachmentHttpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/gdgw-attachments/{id}"
		getCentralNetworkAttachmentProduct = "cc"
	)
	getCentralNetworkAttachmentClient, err := cfg.NewServiceClient(getCentralNetworkAttachmentProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	getCentralNetworkAttachmentPath := getCentralNetworkAttachmentClient.Endpoint + getCentralNetworkAttachmentHttpUrl
	getCentralNetworkAttachmentPath = strings.ReplaceAll(getCentralNetworkAttachmentPath, "{domain_id}", cfg.DomainID)
	getCentralNetworkAttachmentPath = strings.ReplaceAll(getCentralNetworkAttachmentPath, "{central_network_id}",
		d.Get("central_network_id").(string))
	getCentralNetworkAttachmentPath = strings.ReplaceAll(getCentralNetworkAttachmentPath, "{id}", d.Id())

	getCentralNetworkAttachmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getCentralNetworkAttachmentResp, err := getCentralNetworkAttachmentClient.Request("GET",
		getCentralNetworkAttachmentPath, &getCentralNetworkAttachmentOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving central network attachment")
	}

	respBody, err := utils.FlattenResponse(getCentralNetworkAttachmentResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("central_network_gdgw_attachment.name", respBody, nil)),
		d.Set("enterprise_router_id", utils.PathSearch("central_network_gdgw_attachment.enterprise_router_id", respBody, nil)),
		d.Set("enterprise_router_project_id", utils.PathSearch("central_network_gdgw_attachment.enterprise_router_project_id", respBody, nil)),
		d.Set("enterprise_router_region_id", utils.PathSearch("central_network_gdgw_attachment.enterprise_router_region_id", respBody, nil)),
		d.Set("global_dc_gateway_id", utils.PathSearch("central_network_gdgw_attachment.global_dc_gateway_id", respBody, nil)),
		d.Set("global_dc_gateway_project_id", utils.PathSearch("central_network_gdgw_attachment.global_dc_gateway_project_id", respBody, nil)),
		d.Set("global_dc_gateway_region_id", utils.PathSearch("central_network_gdgw_attachment.global_dc_gateway_region_id", respBody, nil)),
		d.Set("central_network_plane_id", utils.PathSearch("central_network_gdgw_attachment.central_network_plane_id", respBody, nil)),
		d.Set("state", utils.PathSearch("central_network_gdgw_attachment.state", respBody, nil)),
		d.Set("central_network_id", utils.PathSearch("central_network_gdgw_attachment.central_network_id", respBody, nil)),
		d.Set("description", utils.PathSearch("central_network_gdgw_attachment.description", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCentralNetworkAttachmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateCentralNetworkAttachmentChanges := []string{
		"name",
		"description",
	}

	if d.HasChanges(updateCentralNetworkAttachmentChanges...) {
		var (
			updateCentralNetworkAttachmentHttpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/gdgw-attachments/{id}"
			updateCentralNetworkAttachmentProduct = "cc"
		)
		updateCentralNetworkAttachmentClient, err := cfg.NewServiceClient(updateCentralNetworkAttachmentProduct, region)
		if err != nil {
			return diag.Errorf("error creating CC client: %s", err)
		}

		updateCentralNetworkAttachmentPath := updateCentralNetworkAttachmentClient.Endpoint + updateCentralNetworkAttachmentHttpUrl
		updateCentralNetworkAttachmentPath = strings.ReplaceAll(updateCentralNetworkAttachmentPath, "{domain_id}", cfg.DomainID)
		updateCentralNetworkAttachmentPath = strings.ReplaceAll(updateCentralNetworkAttachmentPath, "{central_network_id}",
			d.Get("central_network_id").(string))
		updateCentralNetworkAttachmentPath = strings.ReplaceAll(updateCentralNetworkAttachmentPath, "{id}", d.Id())

		updateCentralNetworkAttachmentOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/json"},
		}

		updateCentralNetworkAttachmentOpt.JSONBody = utils.RemoveNil(buildUpdateCentralNetworkAttachmentBodyParams(d))
		_, err = updateCentralNetworkAttachmentClient.Request("PUT", updateCentralNetworkAttachmentPath,
			&updateCentralNetworkAttachmentOpt)
		if err != nil {
			return diag.Errorf("error updating central network attachment: %s", err)
		}

		err = centralNetworkAttachmentWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the central network attachment (%s) update to complete: %s", d.Id(), err)
		}
	}

	return resourceCentralNetworkAttachmentRead(ctx, d, meta)
}

func buildUpdateCentralNetworkAttachmentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"central_network_gdgw_attachment": map[string]interface{}{
			"name":        utils.ValueIgnoreEmpty(d.Get("name")),
			"description": utils.ValueIgnoreEmpty(d.Get("description")),
		},
	}
	return bodyParams
}

func resourceCentralNetworkAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteCentralNetworkAttachmentHttpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/attachments/{id}"
		deleteCentralNetworkAttachmentProduct = "cc"
	)
	deleteCentralNetworkAttachmentClient, err := cfg.NewServiceClient(deleteCentralNetworkAttachmentProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	deleteCentralNetworkAttachmentPath := deleteCentralNetworkAttachmentClient.Endpoint + deleteCentralNetworkAttachmentHttpUrl
	deleteCentralNetworkAttachmentPath = strings.ReplaceAll(deleteCentralNetworkAttachmentPath, "{domain_id}", cfg.DomainID)
	deleteCentralNetworkAttachmentPath = strings.ReplaceAll(deleteCentralNetworkAttachmentPath, "{central_network_id}",
		d.Get("central_network_id").(string))
	deleteCentralNetworkAttachmentPath = strings.ReplaceAll(deleteCentralNetworkAttachmentPath, "{id}", d.Id())

	deleteCentralNetworkAttachmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteCentralNetworkAttachmentClient.Request("DELETE", deleteCentralNetworkAttachmentPath,
		&deleteCentralNetworkAttachmentOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting central network attachment")
	}

	err = centralNetworkAttachmentDeleteWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the central network attachment (%s) deletion to complete: %s", d.Id(), err)
	}

	return nil
}

func centralNetworkAttachmentDeleteWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			// centralNetworkAttachmentDeleteWaiting: missing operation notes
			var (
				deleteWaitingHttpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/gdgw-attachments/{id}"
				deleteWaitingProduct = "cc"
			)
			centralNetworkAttachmentDeleteWaitingClient, err := cfg.NewServiceClient(deleteWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating CC client: %s", err)
			}

			deleteWaitingPath := centralNetworkAttachmentDeleteWaitingClient.Endpoint + deleteWaitingHttpUrl
			deleteWaitingPath = strings.ReplaceAll(deleteWaitingPath, "{domain_id}", cfg.DomainID)
			deleteWaitingPath = strings.ReplaceAll(deleteWaitingPath, "{central_network_id}",
				d.Get("central_network_id").(string))
			deleteWaitingPath = strings.ReplaceAll(deleteWaitingPath, "{id}", d.Id())

			deleteWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
				MoreHeaders: map[string]string{"Content-Type": "application/json"},
			}

			deleteWaitingResp, err := centralNetworkAttachmentDeleteWaitingClient.Request("GET", deleteWaitingPath, &deleteWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return deleteWaitingResp, "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			centralNetworkAttachmentDeleteWaitingRespBody, err := utils.FlattenResponse(deleteWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`central_network_gdgw_attachment.state`, centralNetworkAttachmentDeleteWaitingRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `central_network_gdgw_attachment.state`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"DELETED",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return centralNetworkAttachmentDeleteWaitingRespBody, "COMPLETED", nil
			}

			pendingStatus := []string{
				"DELETING",
				"FREEZING",
				"UNFREEZING",
				"RECOVERING",
			}
			if utils.StrSliceContains(pendingStatus, status) {
				return centralNetworkAttachmentDeleteWaitingRespBody, "PENDING", nil
			}

			return centralNetworkAttachmentDeleteWaitingRespBody, status, nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceCentralNetworkAttachmentImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <central_network_id>/<id>")
	}

	d.Set("central_network_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
