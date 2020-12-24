package huaweicloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/vpcep/v1/services"
)

func DataSourceVPCEPPublicServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcepPublicRead,

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

func dataSourceVpcepPublicRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	vpcepClient, err := config.VPCEPClient(region)
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud VPC endpoint client: %s", err)
	}

	listOpts := services.ListOpts{
		ServiceName: d.Get("service_name").(string),
		ID:          d.Get("service_id").(string),
	}

	allServices, err := services.ListPublic(vpcepClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve vpc endpoint public services: %s", err)
	}

	if len(allServices) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	d.SetId(allServices[0].ID)
	d.Set("region", region)
	services := make([]map[string]interface{}, len(allServices))
	for i, v := range allServices {
		services[i] = map[string]interface{}{
			"id":           v.ID,
			"service_name": v.ServiceName,
			"service_type": v.ServiceType,
			"owner":        v.Owner,
			"is_charge":    v.IsChange,
		}
	}
	if err := d.Set("services", services); err != nil {
		return err
	}

	return nil
}
