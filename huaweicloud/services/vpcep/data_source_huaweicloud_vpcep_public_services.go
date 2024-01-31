package vpcep

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/vpcep/v1/services"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API VPCEP GET /v1/{project_id}/vpc-endpoint-services/public
func DataSourceVPCEPPublicServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcepPublicRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_charge": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVpcepPublicRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	vpcepClient, err := cfg.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	listOpts := services.ListOpts{
		ServiceName: d.Get("service_name").(string),
		ID:          d.Get("service_id").(string),
	}

	allServices, err := services.ListAllPublics(vpcepClient, listOpts)
	if err != nil {
		return diag.Errorf("unable to retrieve VPC endpoint public services: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("services", flattenListVpcEndpointsServices(allServices)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListVpcEndpointsServices(allServices []services.PublicService) []map[string]interface{} {
	if allServices == nil {
		return nil
	}
	endpointServices := make([]map[string]interface{}, len(allServices))
	for i, v := range allServices {
		endpointServices[i] = map[string]interface{}{
			"id":           v.ID,
			"service_name": v.ServiceName,
			"service_type": v.ServiceType,
			"owner":        v.Owner,
			"is_charge":    v.IsChange,
		}
	}
	return endpointServices
}
