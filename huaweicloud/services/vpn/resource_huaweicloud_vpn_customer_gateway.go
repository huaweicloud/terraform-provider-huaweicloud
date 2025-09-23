// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product VPN
// ---------------------------------------------------------------

package vpn

import (
	"context"
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

// @API VPN POST /v5/{project_id}/customer-gateways
// @API VPN DELETE /v5/{project_id}/customer-gateways/{id}
// @API VPN GET /v5/{project_id}/customer-gateways/{id}
// @API VPN PUT /v5/{project_id}/customer-gateways/{id}
// @API VPN POST /v5/{project_id}/{resource_type}/{resource_id}/tags/create
// @API VPN POST /v5/{project_id}/{resource_type}/{resource_id}/tags/delete
func ResourceCustomerGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomerGatewayCreate,
		UpdateContext: resourceCustomerGatewayUpdate,
		ReadContext:   resourceCustomerGatewayRead,
		DeleteContext: resourceCustomerGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The customer gateway name.`,
			},
			"id_value": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The identifier of a customer gateway.`,
					utils.SchemaDescInput{
						Required: true,
					},
				),
				ConflictsWith: []string{"ip", "route_mode"},
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: utils.SchemaDesc(
					`The IP address of the customer gateway.`,
					utils.SchemaDescInput{
						Deprecated: true,
					}),
				ConflictsWith: []string{"id_value", "id_type"},
			},
			"route_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "bgp",
				ForceNew: true,
				Description: utils.SchemaDesc(
					`The route mode of the customer gateway.`,
					utils.SchemaDescInput{
						Deprecated: true,
					}),
				ValidateFunc: validation.StringInSlice([]string{
					"static", "bgp",
				}, false),
				ConflictsWith: []string{"id_value", "id_type"},
			},
			"asn": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     65000,
				ForceNew:    true,
				Description: `The BGP ASN number of the customer gateway, the default value is 65000.`,
			},
			"id_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "ip",
				ForceNew:      true,
				Description:   `The identifier type of a customer gateway.`,
				ConflictsWith: []string{"ip", "route_mode"},
			},
			"certificate_content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"certificate_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"serial_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signature_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issuer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_updatable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The create time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
		},
	}
}

func resourceCustomerGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createCustomerGateway: Create a VPN Customer Gateway.
	var (
		createCustomerGatewayHttpUrl = "v5/{project_id}/customer-gateways"
		createCustomerGatewayProduct = "vpn"
	)
	createCustomerGatewayClient, err := cfg.NewServiceClient(createCustomerGatewayProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	createCustomerGatewayPath := createCustomerGatewayClient.Endpoint + createCustomerGatewayHttpUrl
	createCustomerGatewayPath = strings.ReplaceAll(createCustomerGatewayPath, "{project_id}", createCustomerGatewayClient.ProjectID)

	createCustomerGatewayOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createCustomerGatewayOpt.JSONBody = utils.RemoveNil(buildCreateCustomerGatewayBodyParams(d))
	createCustomerGatewayResp, err := createCustomerGatewayClient.Request("POST", createCustomerGatewayPath, &createCustomerGatewayOpt)
	if err != nil {
		return diag.Errorf("error creating VPN customer gateway: %s", err)
	}

	createCustomerGatewayRespBody, err := utils.FlattenResponse(createCustomerGatewayResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("customer_gateway.id", createCustomerGatewayRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating VPN customer gateway: ID is not found in API response")
	}
	d.SetId(id)

	return resourceCustomerGatewayRead(ctx, d, meta)
}

func buildCreateCustomerGatewayBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"customer_gateway": buildCreateCustomerGatewayCustomerGatewayChildBody(d),
	}
}

func buildCreateCustomerGatewayCustomerGatewayChildBody(d *schema.ResourceData) map[string]interface{} {
	_, ipOk := d.GetOk("ip")
	params := map[string]interface{}{
		"name":    d.Get("name"),
		"bgp_asn": utils.ValueIgnoreEmpty(d.Get("asn")),
		"tags":    utils.ValueIgnoreEmpty(utils.ExpandResourceTags(d.Get("tags").(map[string]interface{}))),
	}

	if ipOk {
		params["ip"] = d.Get("ip")
		params["route_mode"] = utils.ValueIgnoreEmpty(d.Get("route_mode"))
	} else {
		params["id_value"] = d.Get("id_value")
		params["id_type"] = utils.ValueIgnoreEmpty(d.Get("id_type"))
	}

	if certificateContent, ok := d.GetOk("certificate_content"); ok {
		params["ca_certificate"] = map[string]interface{}{
			"content": certificateContent,
		}
	}
	return params
}

func resourceCustomerGatewayRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getCustomerGateway: Query the VPN customer gateway detail
	var (
		getCustomerGatewayHttpUrl = "v5/{project_id}/customer-gateways/{id}"
		getCustomerGatewayProduct = "vpn"
	)
	getCustomerGatewayClient, err := cfg.NewServiceClient(getCustomerGatewayProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	getCustomerGatewayPath := getCustomerGatewayClient.Endpoint + getCustomerGatewayHttpUrl
	getCustomerGatewayPath = strings.ReplaceAll(getCustomerGatewayPath, "{project_id}", getCustomerGatewayClient.ProjectID)
	getCustomerGatewayPath = strings.ReplaceAll(getCustomerGatewayPath, "{id}", d.Id())

	getCustomerGatewayOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getCustomerGatewayResp, err := getCustomerGatewayClient.Request("GET", getCustomerGatewayPath, &getCustomerGatewayOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VPN customer gateway")
	}

	getCustomerGatewayRespBody, err := utils.FlattenResponse(getCustomerGatewayResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("customer_gateway.name", getCustomerGatewayRespBody, nil)),
		d.Set("id_value", utils.PathSearch("customer_gateway.id_value", getCustomerGatewayRespBody, nil)),
		d.Set("asn", utils.PathSearch("customer_gateway.bgp_asn", getCustomerGatewayRespBody, nil)),
		d.Set("id_type", utils.PathSearch("customer_gateway.id_type", getCustomerGatewayRespBody, nil)),
		d.Set("certificate_id", utils.PathSearch("customer_gateway.ca_certificate.id", getCustomerGatewayRespBody, nil)),
		d.Set("serial_number", utils.PathSearch("customer_gateway.ca_certificate.serial_number",
			getCustomerGatewayRespBody, nil)),
		d.Set("signature_algorithm", utils.PathSearch("customer_gateway.ca_certificate.signature_algorithm",
			getCustomerGatewayRespBody, nil)),
		d.Set("issuer", utils.PathSearch("customer_gateway.ca_certificate.issuer",
			getCustomerGatewayRespBody, nil)),
		d.Set("subject", utils.PathSearch("customer_gateway.ca_certificate.subject",
			getCustomerGatewayRespBody, nil)),
		d.Set("expire_time", utils.PathSearch("customer_gateway.ca_certificate.expire_time",
			getCustomerGatewayRespBody, nil)),
		d.Set("is_updatable", utils.PathSearch("customer_gateway.ca_certificate.is_updatable",
			getCustomerGatewayRespBody, nil)),
		d.Set("created_at", utils.PathSearch("customer_gateway.created_at", getCustomerGatewayRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("customer_gateway.updated_at", getCustomerGatewayRespBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("customer_gateway.tags", getCustomerGatewayRespBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCustomerGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	updateCustomerGatewayClient, err := cfg.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	updateCustomerGatewayHasChanges := []string{
		"name",
		"certificate_content",
	}

	if d.HasChanges(updateCustomerGatewayHasChanges...) {
		// updateCustomerGateway: Update the configuration of VPN customer gateway
		updateCustomerGatewayHttpUrl := "v5/{project_id}/customer-gateways/{id}"

		updateCustomerGatewayPath := updateCustomerGatewayClient.Endpoint + updateCustomerGatewayHttpUrl
		updateCustomerGatewayPath = strings.ReplaceAll(updateCustomerGatewayPath, "{project_id}", updateCustomerGatewayClient.ProjectID)
		updateCustomerGatewayPath = strings.ReplaceAll(updateCustomerGatewayPath, "{id}", d.Id())

		updateCustomerGatewayOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		updateCustomerGatewayOpt.JSONBody = utils.RemoveNil(buildUpdateCustomerGatewayBodyParams(d))
		_, err = updateCustomerGatewayClient.Request("PUT", updateCustomerGatewayPath, &updateCustomerGatewayOpt)
		if err != nil {
			return diag.Errorf("error updating VPN customer gateway: %s", err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		tagErr := updateTags(updateCustomerGatewayClient, d, "customer-gateway", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of VPN customer gateway (%s): %s", d.Id(), tagErr)
		}
	}
	return resourceCustomerGatewayRead(ctx, d, meta)
}

func buildUpdateCustomerGatewayBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"customer_gateway": buildUpdateCustomerGatewayCustomerGatewayChildBody(d),
	}
	return bodyParams
}

func buildUpdateCustomerGatewayCustomerGatewayChildBody(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"name": utils.ValueIgnoreEmpty(d.Get("name")),
	}
	if certificateContent, ok := d.GetOk("certificate_content"); ok {
		params["ca_certificate"] = map[string]interface{}{
			"content": certificateContent,
		}
	}
	return params
}

func resourceCustomerGatewayDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteCustomerGateway: Delete an existing VPN customer Gateway
	var (
		deleteCustomerGatewayHttpUrl = "v5/{project_id}/customer-gateways/{id}"
		deleteCustomerGatewayProduct = "vpn"
	)
	deleteCustomerGatewayClient, err := cfg.NewServiceClient(deleteCustomerGatewayProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	deleteCustomerGatewayPath := deleteCustomerGatewayClient.Endpoint + deleteCustomerGatewayHttpUrl
	deleteCustomerGatewayPath = strings.ReplaceAll(deleteCustomerGatewayPath, "{project_id}", deleteCustomerGatewayClient.ProjectID)
	deleteCustomerGatewayPath = strings.ReplaceAll(deleteCustomerGatewayPath, "{id}", d.Id())

	deleteCustomerGatewayOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteCustomerGatewayClient.Request("DELETE", deleteCustomerGatewayPath, &deleteCustomerGatewayOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting VPN customer gateway")
	}

	return nil
}
