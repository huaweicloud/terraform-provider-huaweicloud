package vpc

import (
	"github.com/chnsz/golangsdk/openstack/networking/v1/vpcs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func DataSourceVpcV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcV1Read,

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

func dataSourceVpcV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud Vpc client: %s", err)
	}

	listOpts := vpcs.ListOpts{
		ID:     d.Get("id").(string),
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
		CIDR:   d.Get("cidr").(string),
	}

	epsID := config.GetEnterpriseProjectID(d)

	if epsID != "" {
		listOpts.EnterpriseProjectID = epsID
	}

	refinedVpcs, err := vpcs.List(vpcClient, listOpts)
	if err != nil {
		return fmtp.Errorf("Unable to retrieve vpcs: %s", err)
	}

	if len(refinedVpcs) < 1 {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedVpcs) > 1 {
		return fmtp.Errorf("Your query returned more than one result." +
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

	logp.Printf("[INFO] Retrieved Vpc using given filter %s: %+v", Vpc.ID, Vpc)
	d.SetId(Vpc.ID)

	d.Set("name", Vpc.Name)
	d.Set("cidr", Vpc.CIDR)
	d.Set("enterprise_project_id", Vpc.EnterpriseProjectID)
	d.Set("status", Vpc.Status)
	d.Set("id", Vpc.ID)
	d.Set("region", config.GetRegion(d))
	if err := d.Set("routes", s); err != nil {
		return err
	}

	return nil
}
