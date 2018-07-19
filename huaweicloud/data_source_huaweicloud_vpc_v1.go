package huaweicloud

import (
	"fmt"
	"log"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/vpcs"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceVirtualPrivateCloudVpcV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVirtualPrivateCloudV1Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"shared": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"routes": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"nexthop": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVirtualPrivateCloudV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcClient, err := config.networkingV1Client(GetRegion(d, config))

	listOpts := vpcs.ListOpts{
		ID:     d.Get("id").(string),
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
		CIDR:   d.Get("cidr").(string),
	}

	refinedVpcs, err := vpcs.List(vpcClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve vpcs: %s", err)
	}

	if len(refinedVpcs) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedVpcs) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Vpc := refinedVpcs[0]

	var s []map[string]interface{}
	for _, route := range Vpc.Routes {
		mapping := map[string]interface{}{
			"destination": route.DestinationCIDR,
			"nexthop":     route.NextHop,
		}
		s = append(s, mapping)
	}

	log.Printf("[INFO] Retrieved Vpc using given filter %s: %+v", Vpc.ID, Vpc)
	d.SetId(Vpc.ID)

	d.Set("name", Vpc.Name)
	d.Set("cidr", Vpc.CIDR)
	d.Set("status", Vpc.Status)
	d.Set("id", Vpc.ID)
	d.Set("shared", Vpc.EnableSharedSnat)
	d.Set("region", GetRegion(d, config))
	if err := d.Set("routes", s); err != nil {
		return err
	}

	return nil
}
