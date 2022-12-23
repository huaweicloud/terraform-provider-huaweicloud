package dc

import (
	"context"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/dc/v3/gateways"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceVirtualGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualGatewayCreate,
		ReadContext:   resourceVirtualGatewayRead,
		UpdateContext: resourceVirtualGatewayUpdate,
		DeleteContext: resourceVirtualGatewayDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa5\\w-.]*$"),
						"Only chinese and english letters, digits, hyphens (-), underscores (_) and dots (.) are "+
							"allowed."),
					validation.StringLenBetween(0, 64),
				),
				Description: "The name of the virtual gateway.",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[^<>]*$`),
						"The angle brackets (< and >) are not allowed."),
					validation.StringLenBetween(0, 128),
				),
				Description: "The description of the virtual gateway.",
			},
			"asn": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 4294967295),
				Description:  "The local BGP ASN of the virtual gateway.",
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
		EnterpriseProjectId: common.GetEnterpriseProjectID(d, cfg),
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
		return common.CheckDeletedDiag(d, err, "virtual gateway")
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
		return diag.Errorf("error deleting virtual gateway (%s): %s", gatewayId, err)
	}

	return nil
}
