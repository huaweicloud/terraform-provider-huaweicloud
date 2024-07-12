package nat

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/nat/v3/dnats"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API NAT POST /v3/{project_id}/private-nat/dnat-rules
// @API NAT GET /v3/{project_id}/private-nat/dnat-rules/{dnat_rule_id}
// @API NAT PUT /v3/{project_id}/private-nat/dnat-rules/{dnat_rule_id}
// @API NAT DELETE /v3/{project_id}/private-nat/dnat-rules/{dnat_rule_id}
func ResourcePrivateDnatRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateDnatRuleCreate,
		ReadContext:   resourcePrivateDnatRuleRead,
		UpdateContext: resourcePrivateDnatRuleUpdate,
		DeleteContext: resourcePrivateDnatRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the DNAT rule is located.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The private NAT gateway ID to which the DNAT rule belongs.",
			},
			"transit_ip_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the transit IP for private NAT.",
			},
			"transit_service_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The port of the transit IP.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The protocol type.",
			},
			"backend_interface_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The network interface ID of the transit IP for private NAT.",
				ExactlyOneOf: []string{"backend_private_ip"},
			},
			"backend_private_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The private IP address of the backend instance.",
			},
			"internal_service_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The port of the backend instance.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the DNAT rule.",
			},
			"backend_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of backend instance.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the DNAT rule.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the DNAT rule.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the enterprise project to which the private DNAT rule belongs.",
			},
		},
	}
}

func resourcePrivateDnatRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NatV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	opts := dnats.CreateOpts{
		GatewayId:           d.Get("gateway_id").(string),
		Protocol:            d.Get("protocol").(string),
		Description:         d.Get("description").(string),
		NetworkInterfaceId:  d.Get("backend_interface_id").(string),
		PrivateIpAddress:    d.Get("backend_private_ip").(string),
		InternalServicePort: strconv.Itoa(d.Get("internal_service_port").(int)),
		TransitIpId:         d.Get("transit_ip_id").(string),
		TransitServicePort:  strconv.Itoa(d.Get("transit_service_port").(int)),
	}

	log.Printf("[DEBUG] The create options of the DNAT rule (Private NAT) is: %#v", opts)
	resp, err := dnats.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating DNAT rule (Private NAT): %s", err)
	}
	d.SetId(resp.ID)

	return resourcePrivateDnatRuleRead(ctx, d, meta)
}

func resourcePrivateDnatRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NatV3Client(region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	resp, err := dnats.Get(client, d.Id())
	if err != nil {
		// If the private DNAT rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving private DNAT rule")
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("gateway_id", resp.GatewayId),
		d.Set("transit_ip_id", resp.TransitIpId),
		d.Set("description", resp.Description),
		d.Set("backend_interface_id", resp.NetworkInterfaceId),
		d.Set("protocol", resp.Protocol),
		d.Set("backend_private_ip", resp.PrivateIpAddress),
		d.Set("backend_type", resp.Type),
		d.Set("created_at", resp.CreatedAt),
		d.Set("updated_at", resp.UpdatedAt),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
	)

	// Parse the internal service port
	if internalServicePort, err := strconv.Atoi(resp.InternalServicePort); err != nil {
		mErr = multierror.Append(mErr, fmt.Errorf("invalid format for internal service port, want 'string', but '%T'",
			resp.InternalServicePort))
	} else if internalServicePort != 0 {
		mErr = multierror.Append(mErr, d.Set("internal_service_port", internalServicePort))
	}

	// Parse the transit service port
	if transitServicePort, err := strconv.Atoi(resp.TransitServicePort); err != nil {
		mErr = multierror.Append(mErr, fmt.Errorf("invalid format for transit service port, want 'string', but '%T'",
			resp.TransitServicePort))
	} else if transitServicePort != 0 {
		mErr = multierror.Append(mErr, d.Set("transit_service_port", transitServicePort))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving resource fields of the DNAT rule (Private NAT): %s", err)
	}
	return nil
}

func resourcePrivateDnatRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NatV3Client(region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	ruleId := d.Id()
	opts := dnats.UpdateOpts{
		Protocol:            d.Get("protocol").(string),
		Description:         utils.String(d.Get("description").(string)),
		InternalServicePort: strconv.Itoa(d.Get("internal_service_port").(int)),
		TransitIpId:         d.Get("transit_ip_id").(string),
		TransitServicePort:  strconv.Itoa(d.Get("transit_service_port").(int)),
	}
	if d.HasChange("backend_private_ip") {
		opts.PrivateIpAddress = d.Get("backend_private_ip").(string)
	} else if d.HasChange("backend_interface_id") {
		opts.NetworkInterfaceId = d.Get("backend_interface_id").(string)
	}

	_, err = dnats.Update(client, ruleId, opts)
	if err != nil {
		return diag.Errorf("error updating DNAT rule (Private NAT) (%s): %s", ruleId, err)
	}

	return resourcePrivateDnatRuleRead(ctx, d, meta)
}

func resourcePrivateDnatRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NatV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	ruleId := d.Id()
	err = dnats.Delete(client, ruleId)
	if err != nil {
		// If the private DNAT rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting private DNAT rule")
	}

	return nil
}
