package dc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var globalGatewayNonUpdatableParams = []string{"bgp_asn", "enterprise_project_id", "tags"}

// @API DC POST /v3/{project_id}/dcaas/global-dc-gateways
// @API DC GET /v3/{project_id}/dcaas/global-dc-gateways/{global_dc_gateway_id}
// @API DC PUT /v3/{project_id}/dcaas/global-dc-gateways/{global_dc_gateway_id}
// @API DC DELETE /v3/{project_id}/dcaas/global-dc-gateways/{global_dc_gateway_id}
func ResourceDcGlobalGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcGlobalGatewayCreate,
		ReadContext:   resourceDcGlobalGatewayRead,
		UpdateContext: resourceDcGlobalGatewayUpdate,
		DeleteContext: resourceDcGlobalGatewayDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(globalGatewayNonUpdatableParams),
			config.MergeDefaultTags(),
		),

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
				Description: `Specifies the name of the global DC gateway.`,
			},
			// Field `description` can be edited to be empty, so the Computed attribute is not added.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the global DC gateway.`,
			},
			"address_family": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the IP address family of the global DC gateway.`,
			},
			// Fields `bgp_asn`, `enterprise_project_id`, and `tags` do not support editing.
			"bgp_asn": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the BGP ASN of the global DC gateway.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the enterprise project ID that the global DC gateway belongs to.`,
			},
			"tags": common.TagsSchema(),
			// Field `enable_force_new` is used internally to control forceNew.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},

			// Attributes
			"reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cause of the failure to create the global DC gateway.`,
			},
			"global_center_network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the central network that the global DC gateway is added to.`,
			},
			"location_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The location where the underlying device of the global DC gateway is deployed.`,
			},
			"locales": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The locale address description information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"en_us": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The region name in English.`,
						},
						"zh_cn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The region name in Chinese.`,
						},
					},
				},
			},
			"current_peer_link_count": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: `The number of peer links allowed on a global DC gateway, indicating the number of enterprise
routers that the global DC gateway can be attached to.`,
			},
			"available_peer_link_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of peer links that can be created for a global DC gateway.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the global DC gateway.`,
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the global DC gateway was created.`,
			},
			"updated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the global DC gateway was updated.`,
			},
			"all_tags": common.TagsComputedSchema(),
		},
	}
}

func buildCreateDcGlobalGatewayBodyParams(d *schema.ResourceData) map[string]interface{} {
	gatewayBodyParams := map[string]interface{}{
		"name":                  d.Get("name"),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"address_family":        utils.ValueIgnoreEmpty(d.Get("address_family")),
		"bgp_asn":               utils.ValueIgnoreEmpty(d.Get("bgp_asn")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(d.Get("enterprise_project_id")),
		"tags":                  utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
	}

	return map[string]interface{}{
		"global_dc_gateway": gatewayBodyParams,
	}
}

func resourceDcGlobalGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/dcaas/global-dc-gateways"
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDcGlobalGatewayBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating DC global gateway: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	gatewayID := utils.PathSearch("global_dc_gateway.id", respBody, "").(string)
	if gatewayID == "" {
		return diag.Errorf("error creating DC global gateway: ID is not found in API response")
	}

	d.SetId(gatewayID)

	return resourceDcGlobalGatewayRead(ctx, d, meta)
}

// Whether to configure the enterprise project ID when querying detailed information does not affect the query results.
// In order to maintain unity with the API, the query conditions for the enterprise project ID are added here.
func buildDcGlobalGatewayQueryParams(cfg *config.Config, d *schema.ResourceData) string {
	epsID := cfg.GetEnterpriseProjectID(d)
	if epsID == "" {
		return ""
	}

	return fmt.Sprintf("?enterprise_project_id=%s", epsID)
}

func flattenLocalesAttribute(localeResp interface{}) []map[string]interface{} {
	if localeResp == nil {
		return nil
	}

	localAttribute := map[string]interface{}{
		"en_us": utils.PathSearch("en_us", localeResp, nil),
		"zh_cn": utils.PathSearch("zh_cn", localeResp, nil),
	}

	return []map[string]interface{}{localAttribute}
}

func resourceDcGlobalGatewayRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/dcaas/global-dc-gateways/{global_dc_gateway_id}"
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{global_dc_gateway_id}", d.Id())
	requestPath += buildDcGlobalGatewayQueryParams(cfg, d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	// When the resource does not exist, the query API returns status code `404`.
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DC global gateway")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	gatewayResp := utils.PathSearch("global_dc_gateway", respBody, nil)
	if gatewayResp == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", gatewayResp, nil)),
		d.Set("description", utils.PathSearch("description", gatewayResp, nil)),
		d.Set("address_family", utils.PathSearch("address_family", gatewayResp, nil)),
		d.Set("bgp_asn", utils.PathSearch("bgp_asn", gatewayResp, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", gatewayResp, nil)),
		d.Set("tags", utils.FlattenSameKeyTagsToMap(d, utils.PathSearch("tags", gatewayResp, nil))),
		d.Set("reason", utils.PathSearch("reason", gatewayResp, nil)),
		d.Set("global_center_network_id", utils.PathSearch("global_center_network_id", gatewayResp, nil)),
		d.Set("location_name", utils.PathSearch("location_name", gatewayResp, nil)),
		d.Set("locales", flattenLocalesAttribute(utils.PathSearch("locales", gatewayResp, nil))),
		d.Set("current_peer_link_count", utils.PathSearch("current_peer_link_count", gatewayResp, nil)),
		d.Set("available_peer_link_count", utils.PathSearch("available_peer_link_count", gatewayResp, nil)),
		d.Set("status", utils.PathSearch("status", gatewayResp, nil)),
		d.Set("created_time", utils.PathSearch("created_time", gatewayResp, nil)),
		d.Set("updated_time", utils.PathSearch("updated_time", gatewayResp, nil)),
		d.Set("all_tags", utils.FlattenTagsToMap(utils.PathSearch("tags", gatewayResp, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateDcGlobalGatewayBodyParams(d *schema.ResourceData) map[string]interface{} {
	gatewayBodyParams := map[string]interface{}{
		"name":           d.Get("name"),
		"description":    d.Get("description"),
		"address_family": d.Get("address_family"),
	}

	return map[string]interface{}{
		"global_dc_gateway": gatewayBodyParams,
	}
}

func resourceDcGlobalGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	if d.HasChanges("name", "description", "address_family") {
		requestPath := client.Endpoint + "v3/{project_id}/dcaas/global-dc-gateways/{global_dc_gateway_id}"
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestPath = strings.ReplaceAll(requestPath, "{global_dc_gateway_id}", d.Id())
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildUpdateDcGlobalGatewayBodyParams(d),
		}

		_, err = client.Request("PUT", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating DC global gateway: %s", err)
		}
	}

	return resourceDcGlobalGatewayRead(ctx, d, meta)
}

func resourceDcGlobalGatewayDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/dcaas/global-dc-gateways/{global_dc_gateway_id}"
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{global_dc_gateway_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting DC global gateway: %s", err)
	}

	return nil
}
