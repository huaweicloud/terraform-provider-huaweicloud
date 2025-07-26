package dc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/dc/v3/gateways"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DC DELETE /v3/{project_id}/dcaas/virtual-gateways/{gatewayId}
// @API DC GET /v3/{project_id}/dcaas/virtual-gateways/{gatewayId}
// @API DC PUT /v3/{project_id}/dcaas/virtual-gateways/{gatewayId}
// @API DC POST /v3/{project_id}/dcaas/virtual-gateways
func ResourceVirtualGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualGatewayCreate,
		ReadContext:   resourceVirtualGatewayRead,
		UpdateContext: resourceVirtualGatewayUpdate,
		DeleteContext: resourceVirtualGatewayDelete,

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
				Description: "The region where the virtual gateway is located.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the VPC connected to the virtual gateway.",
			},
			"local_ep_group": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "The list of IPv6 subnets from the virtual gateway to access cloud services, which is " +
					"usually the CIDR block of the VPC.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the virtual gateway.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the virtual gateway.",
			},
			"asn": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The local BGP ASN of the virtual gateway.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The enterprise project ID to which the virtual gateway belongs.",
			},
			// Attributes
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the virtual gateway.",
			},
			"tags": common.TagsSchema(),
		},
	}
}

func buildVirtualGatewayCreateOpts(d *schema.ResourceData, cfg *config.Config) gateways.CreateOpts {
	return gateways.CreateOpts{
		VpcId:               d.Get("vpc_id").(string),
		LocalEpGroup:        utils.ExpandToStringList(d.Get("local_ep_group").([]interface{})),
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		BgpAsn:              d.Get("asn").(int),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}
}

func resourceVirtualGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DcV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DC v3 client: %s", err)
	}

	opts := buildVirtualGatewayCreateOpts(d, cfg)
	resp, err := gateways.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating virtual gateway: %s", err)
	}
	d.SetId(resp.ID)

	// create tags
	if err := utils.CreateResourceTags(client, d, "dc-vgw", d.Id()); err != nil {
		return diag.Errorf("error setting tags of DC virtual gateway %s: %s", d.Id(), err)
	}

	return resourceVirtualGatewayRead(ctx, d, meta)
}

func resourceVirtualGatewayRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DcV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DC v3 client: %s", err)
	}

	gatewayId := d.Id()
	resp, err := gateways.Get(client, gatewayId)
	if err != nil {
		// When the gateway does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving DC virtual gateway")
	}

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("vpc_id", resp.VpcId),
		d.Set("local_ep_group", resp.LocalEpGroup),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("asn", resp.BgpAsn),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("status", resp.Status),
		utils.SetResourceTagsToState(d, client, "dc-vgw", d.Id()),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving virtual gateway fields: %s", err)
	}
	return nil
}

func resourceVirtualGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DcV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DC v3 client: %s", err)
	}

	var (
		gatewayId = d.Id()

		opts = gateways.UpdateOpts{
			Name:         d.Get("name").(string),
			Description:  utils.String(d.Get("description").(string)),
			LocalEpGroup: utils.ExpandToStringList(d.Get("local_ep_group").([]interface{})),
		}
	)
	_, err = gateways.Update(client, gatewayId, opts)
	if err != nil {
		return diag.Errorf("error updating virtual gateway (%s): %s", gatewayId, err)
	}

	// update tags
	tagErr := utils.UpdateResourceTags(client, d, "dc-vgw", d.Id())
	if tagErr != nil {
		return diag.Errorf("error updating tags of DC virtual gateway %s: %s", d.Id(), tagErr)
	}

	return resourceVirtualGatewayRead(ctx, d, meta)
}

func resourceVirtualGatewayDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DcV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DC v3 client: %s", err)
	}

	gatewayId := d.Id()
	err = gateways.Delete(client, gatewayId)
	if err != nil {
		// When the gateway does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting DC virtual gateway")
	}

	return nil
}
