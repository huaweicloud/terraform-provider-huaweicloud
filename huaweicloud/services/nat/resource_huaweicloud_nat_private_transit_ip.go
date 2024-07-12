package nat

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/nat/v3/transitips"

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

func resourcePrivateTransitIpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NatV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	opts := transitips.CreateOpts{
		SubnetId:            d.Get("subnet_id").(string),
		IpAddress:           d.Get("ip_address").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		Tags:                utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}

	log.Printf("[DEBUG] The create options of the transit IP (Private NAT) is: %#v", opts)
	resp, err := transitips.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating transit IP (Private NAT): %s", err)
	}
	d.SetId(resp.ID)

	return resourcePrivateTransitIpRead(ctx, d, meta)
}

func resourcePrivateTransitIpRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	natClient, err := cfg.NatV3Client(region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	resp, err := transitips.Get(natClient, d.Id())
	if err != nil {
		// If the transit IP does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving transit IP (private NAT)")
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("subnet_id", resp.SubnetId),
		d.Set("ip_address", resp.IpAddress),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("tags", utils.TagsToMap(resp.Tags)),
		d.Set("gateway_id", resp.GatewayId),
		d.Set("network_interface_id", resp.NetworkInterfaceId),
		d.Set("created_at", resp.CreatedAt),
		d.Set("updated_at", resp.UpdatedAt),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving transit IP (Private NAT) fields: %s", err)
	}
	return nil
}

func resourcePrivateTransitIpUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	natClient, err := cfg.NatV3Client(region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	transitIpId := d.Id()
	err = utils.UpdateResourceTags(natClient, d, "transit-ips", transitIpId)
	if err != nil {
		return diag.Errorf("error updating tags of the transit IP (Private NAT) (%s): %s", transitIpId, err)
	}

	return resourcePrivateTransitIpRead(ctx, d, meta)
}

func resourcePrivateTransitIpDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NatV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	transitIpId := d.Id()
	err = transitips.Delete(client, transitIpId)
	if err != nil {
		// If the transit IP does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting transit IP (private NAT)")
	}

	return nil
}
