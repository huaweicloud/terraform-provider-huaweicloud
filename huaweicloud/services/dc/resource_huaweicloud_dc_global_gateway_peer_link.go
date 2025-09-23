package dc

import (
	"context"
	"errors"
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

var peerLinkNonUpdatableParams = []string{"global_dc_gateway_id", "peer_site", "peer_site.*.gateway_id",
	"peer_site.*.project_id", "peer_site.*.region_id"}

// @API DC POST /v3/{project_id}/dcaas/global-dc-gateways/{global_dc_gateway_id}/peer-links
// @API DC GET /v3/{project_id}/dcaas/global-dc-gateways/{global_dc_gateway_id}/peer-links/{peer_link_id}
// @API DC PUT /v3/{project_id}/dcaas/global-dc-gateways/{global_dc_gateway_id}/peer-links/{peer_link_id}
// @API DC DELETE /v3/{project_id}/dcaas/global-dc-gateways/{global_dc_gateway_id}/peer-links/{peer_link_id}
func ResourceDcGlobalGatewayPeerLink() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcGlobalGatewayPeerLinkCreate,
		ReadContext:   resourceDcGlobalGatewayPeerLinkRead,
		UpdateContext: resourceDcGlobalGatewayPeerLinkUpdate,
		DeleteContext: resourceDcGlobalGatewayPeerLinkDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDcGlobalGatewayPeerLinkImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(peerLinkNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the global DC gateway peer link.`,
			},
			// Fields `global_dc_gateway_id` and `peer_site` do not support editing.
			"global_dc_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the global DC gateway ID.`,
			},
			"peer_site": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Specifies the site of the peer link.`,
				Elem:        peerSiteSchema(),
			},
			// Field `description` can be edited to be empty, so the Computed attribute is not added.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of global DC gateway peer link.`,
			},
			// Field `enable_force_new` is used internally to control forceNew.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cause of the failure to add the peer link.`,
			},
			"bandwidth_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The bandwidth information.`,
				Elem:        bandwidthInfoSchema(),
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the peer link.`,
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the peer link was added.`,
			},
			"updated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the peer link was updated.`,
			},
			"create_owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cloud service where the peer link is used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the instance associated with the peer link.`,
			},
		},
	}
}

func peerSiteSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of enterprise router that the global DC gateway is attached to.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the project ID of the enterprise router that the global DC gateway is attached to.`,
			},
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the region ID of the enterprise router that the global DC gateway is attached to.`,
			},
			"link_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The connection ID of the peer gateway at the peer site.`,
			},
			"site_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The site information of the global DC gateway.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the peer gateway.`,
			},
		},
	}
}

func bandwidthInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bandwidth_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The bandwidth size.`,
			},
			"gcb_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The global connection bandwidth ID.`,
			},
		},
	}
}

func buildPeerLinkPeerSiteBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawArray := d.Get("peer_site").([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"gateway_id": rawMap["gateway_id"],
		"project_id": rawMap["project_id"],
		"region_id":  rawMap["region_id"],
	}
}

func buildCreateDcGlobalGatewayPeerLinkBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"peer_site":   buildPeerLinkPeerSiteBodyParams(d),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return map[string]interface{}{
		"peer_link": bodyParams,
	}
}

func GetPeerLinkDetail(client *golangsdk.ServiceClient, gatewayID, id string) (interface{}, error) {
	requestPath := client.Endpoint + "v3/{project_id}/dcaas/global-dc-gateways/{global_dc_gateway_id}/peer-links/{peer_link_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{global_dc_gateway_id}", gatewayID)
	requestPath = strings.ReplaceAll(requestPath, "{peer_link_id}", id)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func waitingForPeerLinkActive(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetPeerLinkDetail(client, d.Get("global_dc_gateway_id").(string), d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("peer_link.status", respBody, "").(string)
			if status == "" {
				return nil, "ERROR", errors.New("the status in detail API response is empty")
			}

			if status == "ACTIVE" {
				return respBody, "COMPLETED", nil
			}

			if status == "ERROR" {
				reason := utils.PathSearch("peer_link.reason", respBody, "").(string)
				return nil, "ERROR", fmt.Errorf("the status in detail API response is ERROR, cause: %s", reason)
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceDcGlobalGatewayPeerLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v3/{project_id}/dcaas/global-dc-gateways/{global_dc_gateway_id}/peer-links"
		product   = "dc"
		gatewayID = d.Get("global_dc_gateway_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{global_dc_gateway_id}", gatewayID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDcGlobalGatewayPeerLinkBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating DC global gateway peer link: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("peer_link.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating DC global gateway peer link: ID is not found in API response")
	}
	d.SetId(id)
	if err := waitingForPeerLinkActive(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for DC global gateway peer link (%s) creation to active: %s", id, err)
	}

	return resourceDcGlobalGatewayPeerLinkRead(ctx, d, meta)
}

func resourceDcGlobalGatewayPeerLinkRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		product   = "dc"
		gatewayID = d.Get("global_dc_gateway_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	respBody, err := GetPeerLinkDetail(client, gatewayID, d.Id())
	// When the resource does not exist, the query API returns status code `404`.
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DC global gateway peer link")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("peer_link.name", respBody, nil)),
		d.Set("global_dc_gateway_id", utils.PathSearch("peer_link.global_dc_gateway_id", respBody, nil)),
		d.Set("peer_site", flattenPeerSiteAttribute(respBody)),
		d.Set("description", utils.PathSearch("peer_link.description", respBody, nil)),
		d.Set("reason", utils.PathSearch("peer_link.reason", respBody, nil)),
		d.Set("bandwidth_info", flattenBandwidthInfoAttribute(respBody)),
		d.Set("status", utils.PathSearch("peer_link.status", respBody, nil)),
		d.Set("created_time", utils.PathSearch("peer_link.created_time", respBody, nil)),
		d.Set("updated_time", utils.PathSearch("peer_link.updated_time", respBody, nil)),
		d.Set("create_owner", utils.PathSearch("peer_link.create_owner", respBody, nil)),
		d.Set("instance_id", utils.PathSearch("peer_link.instance_id", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBandwidthInfoAttribute(respBody interface{}) []interface{} {
	bandwidthInfo := utils.PathSearch("peer_link.bandwidth_info", respBody, nil)
	if bandwidthInfo == nil {
		return nil
	}

	rawMap := map[string]interface{}{
		"bandwidth_size": utils.PathSearch("bandwidth_size", bandwidthInfo, nil),
		"gcb_id":         utils.PathSearch("gcb_id", bandwidthInfo, nil),
	}
	return []interface{}{rawMap}
}

func flattenPeerSiteAttribute(respBody interface{}) []interface{} {
	peerSite := utils.PathSearch("peer_link.peer_site", respBody, nil)
	if peerSite == nil {
		return nil
	}

	rawMap := map[string]interface{}{
		"gateway_id": utils.PathSearch("gateway_id", peerSite, nil),
		"project_id": utils.PathSearch("project_id", peerSite, nil),
		"region_id":  utils.PathSearch("region_id", peerSite, nil),
		"link_id":    utils.PathSearch("link_id", peerSite, nil),
		"site_code":  utils.PathSearch("site_code", peerSite, nil),
		"type":       utils.PathSearch("type", peerSite, nil),
	}
	return []interface{}{rawMap}
}

func buildUpdateDcGlobalGatewayPeerLinkBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
	}

	return map[string]interface{}{
		"peer_link": bodyParams,
	}
}

func updateDcGlobalGatewayPeerLink(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v3/{project_id}/dcaas/global-dc-gateways/{global_dc_gateway_id}/peer-links/{peer_link_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{global_dc_gateway_id}", d.Get("global_dc_gateway_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{peer_link_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateDcGlobalGatewayPeerLinkBodyParams(d),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func resourceDcGlobalGatewayPeerLinkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	if d.HasChanges("name", "description") {
		if err := updateDcGlobalGatewayPeerLink(client, d); err != nil {
			return diag.Errorf("error updating DC global gateway peer link: %s", err)
		}

		if err := waitingForPeerLinkActive(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for DC global gateway peer link (%s) update to active: %s", d.Id(), err)
		}
	}

	return resourceDcGlobalGatewayPeerLinkRead(ctx, d, meta)
}

func waitingForPeerLinkDelete(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetPeerLinkDetail(client, d.Get("global_dc_gateway_id").(string), d.Id())
			if err != nil {
				var errDefault404 golangsdk.ErrDefault404
				if errors.As(err, &errDefault404) {
					return "success deleted", "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceDcGlobalGatewayPeerLinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v3/{project_id}/dcaas/global-dc-gateways/{global_dc_gateway_id}/peer-links/{peer_link_id}"
		product   = "dc"
		gatewayID = d.Get("global_dc_gateway_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{global_dc_gateway_id}", gatewayID)
	requestPath = strings.ReplaceAll(requestPath, "{peer_link_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting DC global gateway peer link: %s", err)
	}

	if err := waitingForPeerLinkDelete(ctx, client, d, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for DC global gateway peer link (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}

func resourceDcGlobalGatewayPeerLinkImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<global_dc_gateway_id>/<id>',"+
			" but got '%s'", d.Id())
	}
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("global_dc_gateway_id", parts[0])
}
