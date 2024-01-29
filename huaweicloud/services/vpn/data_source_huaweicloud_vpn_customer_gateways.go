package vpn

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/vpn/v5/customer_gateways"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPN GET /v5/{project_id}/customer-gateways
func DataSourceVpnCustomerGateways() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceCustomerGatewaysRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"customer_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"route_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"asn": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"customer_gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"asn": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ca_certificate": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     caCertificateSchema(),
						},
					},
				},
			},
		},
	}
}

func caCertificateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
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
		},
	}
	return &sc
}

func datasourceCustomerGatewaysRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.VpnV5Client(region)
	if err != nil {
		return diag.Errorf("error creating vpn v5 client: %s", err)
	}
	opts := customer_gateways.ListOpts{}
	allCustomerGateway, err := customer_gateways.List(client, opts)
	if err != nil {
		return diag.Errorf("unable to list customer gateways: %s ", err)
	}

	log.Printf("[DEBUG] retrieved VPN customer gateways: %#v", allCustomerGateway)
	filter := map[string]interface{}{
		"ID":        d.Get("customer_gateway_id"),
		"Name":      d.Get("name").(string),
		"BGPAsn":    d.Get("asn").(int),
		"RouteMode": d.Get("route_mode").(string),
		"IP":        d.Get("ip").(string),
	}

	filterCustomerGateways, err := utils.FilterSliceWithField(allCustomerGateway, filter)
	if err != nil {
		return diag.Errorf("filter customer gateways failed: %s", err)
	}

	var customerGateways []map[string]interface{}
	for _, item := range filterCustomerGateways {
		customerGateway := item.(customer_gateways.CustomerGateway)
		customerGateways = append(customerGateways, flattenSourceCustomerGateway(customerGateway))
	}

	uuidStr, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuidStr)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("customer_gateways", customerGateways),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSourceCustomerGateway(gateway customer_gateways.CustomerGateway) map[string]interface{} {
	certificate := gateway.CaCertificate

	caCertificate := []map[string]interface{}{
		{
			"serial_number":       certificate.SerialNumber,
			"signature_algorithm": certificate.SignatureAlgorithm,
			"issuer":              certificate.Issuer,
			"subject":             certificate.Subject,
			"expire_time":         certificate.ExpireTime,
			"is_updatable":        certificate.IsUpdatable,
		},
	}

	resourceCustomerGateway := map[string]interface{}{
		"id":             gateway.ID,
		"name":           gateway.Name,
		"ip":             gateway.IP,
		"route_mode":     gateway.RouteMode,
		"asn":            gateway.BGPAsn,
		"id_type":        gateway.IDType,
		"id_value":       gateway.IDValue,
		"created_at":     gateway.CreatedAt,
		"updated_at":     gateway.UpdatedAt,
		"ca_certificate": caCertificate,
	}
	return resourceCustomerGateway
}
