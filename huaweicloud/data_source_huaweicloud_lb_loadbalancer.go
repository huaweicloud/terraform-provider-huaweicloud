package huaweicloud

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/elb/v2/loadbalancers"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func dataSourceELBV2Loadbalancer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceELBV2LoadbalancerRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ONLINE", "FROZEN",
				}, true),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vip_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vip_port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceELBV2LoadbalancerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	elbClient, err := config.LoadBalancerClient(region)
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud elb client %s", err)
	}
	// Client for getting tags
	elbV2Client, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud elb v2.0 client: %s", err)
	}
	listOpts := loadbalancers.ListOpts{
		Name:                d.Get("name").(string),
		ID:                  d.Get("id").(string),
		OperatingStatus:     d.Get("status").(string),
		Description:         d.Get("description").(string),
		VipAddress:          d.Get("vip_address").(string),
		VipSubnetID:         d.Get("vip_subnet_id").(string),
		EnterpriseProjectID: GetEnterpriseProjectID(d, config),
	}
	pages, err := loadbalancers.List(elbClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to retrieve loadbalancers: %s", err)
	}
	lbList, err := loadbalancers.ExtractLoadBalancers(pages)
	if err != nil {
		return fmt.Errorf("Unable to extract loadbalancers: %s", err)
	}

	if len(lbList) < 1 {
		return fmt.Errorf("Your query returned no results, Please change your search criteria and try again")
	}

	if len(lbList) > 1 {
		return fmt.Errorf("Your query returned more than one result, Please try a more specific search criteria")
	}

	lb := lbList[0]
	d.SetId(lb.ID)

	mErr := multierror.Append(
		d.Set("region", GetRegion(d, config)),
		d.Set("name", lb.Name),
		d.Set("status", lb.OperatingStatus),
		d.Set("description", lb.Description),
		d.Set("vip_address", lb.VipAddress),
		d.Set("vip_subnet_id", lb.VipSubnetID),
		d.Set("enterprise_project_id", lb.EnterpriseProjectID),
		d.Set("vip_port_id", lb.VipPortID),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmt.Errorf("Error setting elb loadbalancer fields: %s", err)
	}

	// Get tags
	resourceTags, err := tags.Get(elbV2Client, "loadbalancers", d.Id()).Extract()
	if err != nil {
		fmt.Errorf("Error fetching tags of elb loadbalancer: %s", err)
	}
	tagmap := utils.TagsToMap(resourceTags.Tags)
	d.Set("tags", tagmap)

	return nil
}
