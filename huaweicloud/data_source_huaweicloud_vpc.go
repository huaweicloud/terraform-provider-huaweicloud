package huaweicloud

import (
	"fmt"
	"log"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/vpcs"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceVirtualPrivateCloudVpcV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVirtualPrivateCloudV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nexthop": {
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
	vpcClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud Vpc client: %s", err)
	}

	listOpts := vpcs.ListOpts{
		ID:     d.Get("id").(string),
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
		CIDR:   d.Get("cidr").(string),
	}

	epsID := GetEnterpriseProjectID(d, config)

	if epsID != "" {
		listOpts.EnterpriseProjectID = epsID
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
	d.Set("enterprise_project_id", Vpc.EnterpriseProjectID)
	d.Set("status", Vpc.Status)
	d.Set("id", Vpc.ID)
	d.Set("shared", Vpc.EnableSharedSnat)
	d.Set("region", GetRegion(d, config))
	if err := d.Set("routes", s); err != nil {
		return err
	}

	return nil
}
