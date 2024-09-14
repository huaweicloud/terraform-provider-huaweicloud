package nat

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/nat/v3/snats"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API NAT POST /v3/{project_id}/private-nat/snat-rules
// @API NAT GET /v3/{project_id}/private-nat/snat-rules/{snat_rule_id}
// @API NAT PUT /v3/{project_id}/private-nat/snat-rules/{snat_rule_id}
// @API NAT DELETE /v3/{project_id}/private-nat/snat-rules/{snat_rule_id}
func ResourcePrivateSnatRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateSnatRuleCreate,
		ReadContext:   resourcePrivateSnatRuleRead,
		UpdateContext: resourcePrivateSnatRuleUpdate,
		DeleteContext: resourcePrivateSnatRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the SNAT rule is located.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The private NAT gateway ID to which the SNAT rule belongs.",
			},
			"transit_ip_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the transit IP associated with SNAT rule.",
			},
			"cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"subnet_id"},
				Description:  "The CIDR block of the match rule.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The subnet ID of the match rule.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the SNAT rule.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the SNAT rule.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the SNAT rule.",
			},
			"transit_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The address of the transit IP",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the enterprise project to which the private SNAT rule belongs.",
			},
		},
	}
}

func resourcePrivateSnatRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NatV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	opts := snats.CreateOpts{
		GatewayId:    d.Get("gateway_id").(string),
		TransitIpIds: []string{d.Get("transit_ip_id").(string)},
		Cidr:         d.Get("cidr").(string),
		SubnetId:     d.Get("subnet_id").(string),
		Description:  d.Get("description").(string),
	}

	log.Printf("[DEBUG] The create options of the SNAT rule (Private NAT) is: %#v", opts)
	resp, err := snats.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating SNAT rule (Private NAT): %s", err)
	}
	d.SetId(resp.ID)

	return resourcePrivateSnatRuleRead(ctx, d, meta)
}

func resourcePrivateSnatRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NatV3Client(region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	resp, err := snats.Get(client, d.Id())
	if err != nil {
		// If the private SNAT rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving private SNAT rule")
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("gateway_id", resp.GatewayId),
		d.Set("transit_ip_id", utils.PathSearch("[0].ID", resp.TransitIpAssociations, nil)),
		d.Set("description", resp.Description),
		d.Set("subnet_id", resp.SubnetId),
		d.Set("cidr", resp.Cidr),
		d.Set("created_at", resp.CreatedAt),
		d.Set("updated_at", resp.UpdatedAt),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("transit_ip_address", utils.PathSearch("[0].Address", resp.TransitIpAssociations, nil)),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving resource fields of the private SNAT rule: %s", err)
	}
	return nil
}

func resourcePrivateSnatRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NatV3Client(region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	ruleId := d.Id()
	opts := snats.UpdateOpts{
		TransitIpIds: []string{d.Get("transit_ip_id").(string)},
		Description:  utils.String(d.Get("description").(string)),
	}

	_, err = snats.Update(client, ruleId, opts)
	if err != nil {
		return diag.Errorf("error updating private SNAT rule (%s): %s", ruleId, err)
	}

	return resourcePrivateSnatRuleRead(ctx, d, meta)
}

func resourcePrivateSnatRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NatV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	ruleId := d.Id()
	err = snats.Delete(client, ruleId)
	if err != nil {
		// If the private SNAT rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting private SNAT rule")
	}

	return nil
}
