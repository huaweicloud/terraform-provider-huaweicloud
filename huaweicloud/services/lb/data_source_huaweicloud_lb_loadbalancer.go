package lb

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/elb/v2/loadbalancers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

// @API ELB GET /v2/{project_id}/elb/loadbalancers
// @API ELB GET /v2.0/{project_id}/loadbalancers/{id}/tags
func DataSourceELBV2Loadbalancer() *schema.Resource {
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
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vip_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceELBV2LoadbalancerRead(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	elbClient, err := cfg.LoadBalancerClient(region)
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud elb client %s", err)
	}
	listOpts := loadbalancers.ListOpts{
		Name:                d.Get("name").(string),
		ID:                  d.Get("id").(string),
		OperatingStatus:     d.Get("status").(string),
		Description:         d.Get("description").(string),
		VipAddress:          d.Get("vip_address").(string),
		VipSubnetID:         d.Get("vip_subnet_id").(string),
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d),
	}
	pages, err := loadbalancers.List(elbClient, listOpts).AllPages()
	if err != nil {
		return fmtp.Errorf("Unable to retrieve load balancers: %s", err)
	}
	lbList, err := loadbalancers.ExtractLoadBalancers(pages)
	if err != nil {
		return fmtp.Errorf("Unable to extract load balancers: %s", err)
	}
	if len(lbList) < 1 {
		return fmtp.Errorf("Your query returned no results, please change your search criteria and try again")
	}
	if len(lbList) > 1 {
		return fmtp.Errorf("Your query returned more than one result, please try a more specific search criteria")
	}

	lb := lbList[0]
	d.SetId(lb.ID)

	var publicIp string
	if len(lb.PublicIps) > 0 {
		publicIp = lb.PublicIps[0].PublicIpAddress
	}
	mErr := multierror.Append(
		d.Set("region", cfg.GetRegion(d)),
		d.Set("name", lb.Name),
		d.Set("status", lb.OperatingStatus),
		d.Set("description", lb.Description),
		d.Set("vip_address", lb.VipAddress),
		d.Set("vip_subnet_id", lb.VipSubnetID),
		d.Set("enterprise_project_id", lb.EnterpriseProjectID),
		d.Set("vip_port_id", lb.VipPortID),
		d.Set("public_ip", publicIp),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.Errorf("Error setting elb load balancer fields: %s", err)
	}

	// Get tags for v2.0 API
	elbV2Client, err := cfg.ElbV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb v2.0 client: %s", err)
	}
	resourceTags, err := tags.Get(elbV2Client, "loadbalancers", d.Id()).Extract()
	if err != nil {
		logp.Printf("[WARN] Error fetching tags of elb load balancer %s: %s", d.Id(), err)
	}
	tagmap := utils.TagsToMap(resourceTags.Tags)
	d.Set("tags", tagmap)

	return nil
}
