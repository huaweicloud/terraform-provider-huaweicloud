package nat

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/nat/v2/gateways"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API NAT GET /v2/{project_id}/nat_gateways
func DataSourcePublicGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePublicGatewayRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the public NAT gateway is located.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The public NAT gateway ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The public NAT gateway name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the public NAT gateway.",
			},
			"spec": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The specification of the public NAT gateway.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the VPC to which the public NAT gateway belongs.",
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "The network ID of the downstream interface (the next hop of the DVR) " +
					"of the public NAT gateway.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The current status of the public NAT gateway.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The enterprise project ID of the public NAT gateway.",
			},

			// deprecated
			"router_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use vpc_id instead",
			},
			"internal_network_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use subnet_id instead",
			},
		},
	}
}

func getVpdID(d *schema.ResourceData) string {
	if v, ok := d.GetOk("vpc_id"); ok {
		return v.(string)
	}
	return d.Get("router_id").(string)
}

func getSubnetID(d *schema.ResourceData) string {
	if v, ok := d.GetOk("subnet_id"); ok {
		return v.(string)
	}
	return d.Get("internal_network_id").(string)
}

func dataSourcePublicGatewayRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	natClient, err := cfg.NatGatewayClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	listOpts := gateways.ListOpts{
		ID:                  d.Get("id").(string),
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Spec:                d.Get("spec").(string),
		VpcId:               getVpdID(d),
		InternalNetworkId:   getSubnetID(d),
		Status:              d.Get("status").(string),
		EnterpriseProjectId: d.Get("enterprise_project_id").(string),
	}

	resp, err := gateways.List(natClient, listOpts)
	if err != nil {
		return diag.Errorf("error querying public NAT gateway list: %s", err)
	}

	if len(resp) < 1 {
		return diag.Errorf("your query returned no results, please change your search criteria and try again")
	}

	if len(resp) > 1 {
		return diag.Errorf("Your query returned more than one result, please try a more specific search criteria")
	}

	gateway := resp[0]
	log.Printf("[DEBUG] Retrieved public NAT gateway (%s): %#v", gateway.ID, gateway)
	d.SetId(gateway.ID)

	mErr := multierror.Append(nil,
		d.Set("name", gateway.Name),
		d.Set("description", gateway.Description),
		d.Set("spec", gateway.Spec),
		d.Set("vpc_id", gateway.RouterId),
		d.Set("subnet_id", gateway.InternalNetworkId),
		d.Set("status", gateway.Status),
		d.Set("enterprise_project_id", gateway.EnterpriseProjectId),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving data source fields of the public NAT gateway: %s", err)
	}
	return nil
}
