package iec

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/iec/v1/firewalls"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IEC GET /v1/firewalls
func DataSourceNetworkACL() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkACLRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"id"},
			},

			// Computed
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Computed but always be empty due to the API response
			"networks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "schema: Internal",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"inbound_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"outbound_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceNetworkACLRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	listOpts := firewalls.ListOpts{
		ID:   d.Get("id").(string),
		Name: d.Get("name").(string),
	}

	log.Printf("[DEBUG] query firewall using given filter: %+v", listOpts)
	allFWs, err := firewalls.List(iecClient, listOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to retrieve firewall: %s", err)
	}

	total := len(allFWs.Firewalls)
	if total < 1 {
		return diag.Errorf("your query returned no results. Please change your search criteria and try again")
	}
	if total > 1 {
		return diag.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	fwGroup := allFWs.Firewalls[0]
	log.Printf("[DEBUG] retrieved IEC firewall %s: %+v", fwGroup.ID, fwGroup)

	d.SetId(fwGroup.ID)

	mErr := multierror.Append(nil,
		d.Set("name", fwGroup.Name),
		d.Set("status", fwGroup.Status),
		d.Set("description", fwGroup.Description),
	)

	// currently, the following attributes are empty due to the API response
	networkList := make([]map[string]interface{}, 0, len(fwGroup.Subnets))
	for _, val := range fwGroup.Subnets {
		subnet := make(map[string]interface{})
		subnet["vpc_id"] = val.VpcID
		subnet["subnet_id"] = val.ID
		networkList = append(networkList, subnet)
	}
	mErr = multierror.Append(mErr,
		d.Set("networks", networkList),
		d.Set("inbound_rules", getFirewallRuleIDs(fwGroup.IngressFWPolicy)),
		d.Set("outbound_rules", getFirewallRuleIDs(fwGroup.EgressFWPolicy)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
