package nat

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

// @API NAT POST /v3/{project_id}/private-nat/gateways
// @API NAT GET /v3/{project_id}/private-nat/gateways/{gateway_id}
// @API NAT PUT /v3/{project_id}/private-nat/gateways/{gateway_id}
// @API NAT DELETE /v3/{project_id}/private-nat/gateways/{gateway_id}
// @API NAT POST /v3/{project_id}/private-nat-gateways/{resource_id}/tags/action
func ResourcePrivateGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateGatewayCreate,
		ReadContext:   resourcePrivateGatewayRead,
		UpdateContext: resourcePrivateGatewayUpdate,
		DeleteContext: resourcePrivateGatewayDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the private NAT gateway is located.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The network ID of the subnet to which the private NAT gateway belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the private NAT gateway.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the private NAT gateway.",
			},
			"spec": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The specification of the private NAT gateway.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the enterprise project to which the private NAT gateway belongs.",
			},
			"tags": common.TagsSchema(),
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the private NAT gateway.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the private NAT gateway.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the private NAT gateway.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the VPC to which the private NAT gateway belongs.",
			},
		},
	}
}

func buildDownLinkVpcs(d *schema.ResourceData) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"virsubnet_id": d.Get("subnet_id"),
		},
	}
}

func buildCreatePrivateGatewayBodyParams(d *schema.ResourceData, epsId string) map[string]interface{} {
	gatewayBodyParams := map[string]interface{}{
		"name":                  d.Get("name"),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"downlink_vpcs":         buildDownLinkVpcs(d),
		"spec":                  utils.ValueIgnoreEmpty(d.Get("spec")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(epsId),
		"tags":                  utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
	}

	return map[string]interface{}{
		"gateway": gatewayBodyParams,
	}
}

func resourcePrivateGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/private-nat/gateways"
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePrivateGatewayBodyParams(d, cfg.GetEnterpriseProjectID(d))),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating private NAT gateway: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	gatewayId := utils.PathSearch("gateway.id", respBody, "").(string)
	if gatewayId == "" {
		return diag.Errorf("error creating private NAT gateway: ID is not found in API response")
	}

	d.SetId(gatewayId)

	return resourcePrivateGatewayRead(ctx, d, meta)
}

func GetPrivateGateway(client *golangsdk.ServiceClient, gatewayId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/private-nat/gateways/{gateway_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{gateway_id}", gatewayId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func resourcePrivateGatewayRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	respBody, err := GetPrivateGateway(client, d.Id())
	if err != nil {
		// If the private NAT gateway does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving private NAT gateway")
	}
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("subnet_id", utils.PathSearch("gateway.downlink_vpcs[0].virsubnet_id", respBody, nil)),
		d.Set("name", utils.PathSearch("gateway.name", respBody, nil)),
		d.Set("description", utils.PathSearch("gateway.description", respBody, nil)),
		d.Set("spec", utils.PathSearch("gateway.spec", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("gateway.enterprise_project_id", respBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("gateway.tags", respBody, make([]interface{}, 0)))),
		d.Set("created_at", utils.PathSearch("gateway.created_at", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("gateway.updated_at", respBody, nil)),
		d.Set("status", utils.PathSearch("gateway.status", respBody, nil)),
		d.Set("vpc_id", utils.PathSearch("gateway.downlink_vpcs[0].vpc_id", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdatePrivateGatewayBodyParams(d *schema.ResourceData) map[string]interface{} {
	gatewayBodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": d.Get("description"),
		"spec":        utils.ValueIgnoreEmpty(d.Get("spec")),
	}

	return map[string]interface{}{
		"gateway": gatewayBodyParams,
	}
}

func resourcePrivateGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	if d.HasChanges("name", "description", "spec") {
		httpUrl := "v3/{project_id}/private-nat/gateways/{gateway_id}"
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{gateway_id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildUpdatePrivateGatewayBodyParams(d),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating private NAT gateway: %s", err)
		}
	}

	if d.HasChange("tags") {
		natClient, err := cfg.NatV3Client(region)
		if err != nil {
			return diag.Errorf("error creating NAT v3 client: %s", err)
		}

		err = utils.UpdateResourceTags(natClient, d, "private-nat-gateways", d.Id())
		if err != nil {
			return diag.Errorf("error updating tags of the private NAT gateway: %s", err)
		}
	}

	return resourcePrivateGatewayRead(ctx, d, meta)
}

func resourcePrivateGatewayDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/private-nat/gateways/{gateway_id}"
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{gateway_id}", d.Id())
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		// If the private NAT gateway does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting private NAT gateway")
	}

	return nil
}
