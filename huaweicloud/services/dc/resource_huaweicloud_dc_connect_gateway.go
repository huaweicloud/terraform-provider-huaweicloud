package dc

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DC POST /v3/{project_id}/dcaas/connect-gateways
// @API DC GET /v3/{project_id}/dcaas/connect-gateways/{connect_gateway_id}
// @API DC PUT /v3/{project_id}/dcaas/connect-gateways/{connect_gateway_id}
// @API DC DELETE /v3/{project_id}/dcaas/connect-gateways/{connect_gateway_id}
func ResourceDcConnectGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcConnectGatewayCreate,
		ReadContext:   resourceDcConnectGatewayRead,
		UpdateContext: resourceDcConnectGatewayUpdate,
		DeleteContext: resourceDcConnectGatewayDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address_family": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_site": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bgp_asn": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"current_geip_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gcb_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway_site": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDcConnectGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/dcaas/connect-gateways"
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDcConnectGatewayBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DC connect gateway: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("connect_gateway.id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating DC connect gateway: ID is not found in API response")
	}

	d.SetId(id)

	return resourceDcConnectGatewayRead(ctx, d, meta)
}

func buildCreateDcConnectGatewayBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":           d.Get("name"),
		"description":    utils.ValueIgnoreEmpty(d.Get("description")),
		"address_family": utils.ValueIgnoreEmpty(d.Get("address_family")),
	}

	return map[string]interface{}{
		"connect_gateway": bodyParams,
	}
}

func resourceDcConnectGatewayRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/dcaas/connect-gateways/{connect_gateway_id}"
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{connect_gateway_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DC connect gateway")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("connect_gateway.name", getRespBody, nil)),
		d.Set("description", utils.PathSearch("connect_gateway.description", getRespBody, nil)),
		d.Set("address_family", utils.PathSearch("connect_gateway.address_family", getRespBody, nil)),
		d.Set("status", utils.PathSearch("connect_gateway.status", getRespBody, nil)),
		d.Set("access_site", utils.PathSearch("connect_gateway.access_site", getRespBody, nil)),
		d.Set("bgp_asn", utils.PathSearch("connect_gateway.bgp_asn", getRespBody, nil)),
		d.Set("current_geip_count", utils.PathSearch("connect_gateway.current_geip_count", getRespBody, nil)),
		d.Set("created_time", utils.PathSearch("connect_gateway.created_time", getRespBody, nil)),
		d.Set("updated_time", utils.PathSearch("connect_gateway.updated_time", getRespBody, nil)),
		d.Set("gcb_id", utils.PathSearch("connect_gateway.gcb_id", getRespBody, nil)),
		d.Set("gateway_site", utils.PathSearch("connect_gateway.gateway_site", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDcConnectGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/dcaas/connect-gateways/{connect_gateway_id}"
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	if d.HasChanges("name", "description", "address_family") {
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{connect_gateway_id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildUpdateDcConnectGatewayBodyParams(d),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating DC connect gateway: %s", err)
		}
	}

	return resourceDcConnectGatewayRead(ctx, d, meta)
}

func buildUpdateDcConnectGatewayBodyParams(d *schema.ResourceData) map[string]interface{} {
	gatewayBodyParams := map[string]interface{}{
		"name":           d.Get("name"),
		"description":    d.Get("description"),
		"address_family": d.Get("address_family"),
	}

	return map[string]interface{}{
		"connect_gateway": gatewayBodyParams,
	}
}

func resourceDcConnectGatewayDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/dcaas/connect-gateways/{connect_gateway_id}"
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{connect_gateway_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting DC connect gateway")
	}

	return nil
}
