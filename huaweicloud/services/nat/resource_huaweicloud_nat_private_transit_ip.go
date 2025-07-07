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

// @API NAT POST /v3/{project_id}/private-nat/transit-ips
// @API NAT GET /v3/{project_id}/private-nat/transit-ips/{transit_ip_id}
// @API NAT DELETE /v3/{project_id}/private-nat/transit-ips/{transit_ip_id}
// @API NAT POST /v3/{project_id}/transit-ips/{resource_id}/tags/action
func ResourcePrivateTransitIp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateTransitIpCreate,
		ReadContext:   resourcePrivateTransitIpRead,
		UpdateContext: resourcePrivateTransitIpUpdate,
		DeleteContext: resourcePrivateTransitIpDelete,

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
				Description: "The region where the transit IP is located.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the transit subnet to which the transit IP belongs.",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The IP address of the transit subnet.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the enterprise project to which the transit IP belongs.",
			},
			"tags": common.TagsSchema(),
			"network_interface_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network interface ID of the transit IP for private NAT.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the transit IP.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the private NAT gateway to which the transit IP belongs.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the transit IP for private NAT.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the transit IP for private NAT.",
			},
		},
	}
}

func buildCreatePrivateTransitIpBodyParams(d *schema.ResourceData, epsId string) map[string]interface{} {
	transitIpBodyParams := map[string]interface{}{
		"virsubnet_id":          d.Get("subnet_id"),
		"ip_address":            utils.ValueIgnoreEmpty(d.Get("ip_address")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(epsId),
		"tags":                  utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
	}

	return map[string]interface{}{
		"transit_ip": transitIpBodyParams,
	}
}

func resourcePrivateTransitIpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/private-nat/transit-ips"
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePrivateTransitIpBodyParams(d, cfg.GetEnterpriseProjectID(d))),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating transit IP (private NAT): %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	transitIpId := utils.PathSearch("transit_ip.id", respBody, "").(string)
	if transitIpId == "" {
		return diag.Errorf("error creating transit IP (private NAT): ID is not found in API response")
	}

	d.SetId(transitIpId)

	return resourcePrivateTransitIpRead(ctx, d, meta)
}

func GetTransitIp(client *golangsdk.ServiceClient, transitIpId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/private-nat/transit-ips/{transit_ip_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{transit_ip_id}", transitIpId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func resourcePrivateTransitIpRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	respBody, err := GetTransitIp(client, d.Id())
	if err != nil {
		// If the transit IP does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving transit IP (private NAT)")
	}
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("subnet_id", utils.PathSearch("transit_ip.virsubnet_id", respBody, nil)),
		d.Set("ip_address", utils.PathSearch("transit_ip.ip_address", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("transit_ip.enterprise_project_id", respBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("transit_ip.tags", respBody, make([]interface{}, 0)))),
		d.Set("gateway_id", utils.PathSearch("transit_ip.gateway_id", respBody, nil)),
		d.Set("network_interface_id", utils.PathSearch("transit_ip.network_interface_id", respBody, nil)),
		d.Set("status", utils.PathSearch("transit_ip.status", respBody, nil)),
		d.Set("created_at", utils.PathSearch("transit_ip.created_at", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("transit_ip.updated_at", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePrivateTransitIpUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NatV3Client(region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	if d.HasChange("tags") {
		err = utils.UpdateResourceTags(client, d, "transit-ips", d.Id())
		if err != nil {
			return diag.Errorf("error updating tags of the transit IP (Private NAT): %s", err)
		}
	}

	return resourcePrivateTransitIpRead(ctx, d, meta)
}

func resourcePrivateTransitIpDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/private-nat/transit-ips/{transit_ip_id}"
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{transit_ip_id}", d.Id())
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		// If the transit IP does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting transit IP (private NAT)")
	}

	return nil
}
