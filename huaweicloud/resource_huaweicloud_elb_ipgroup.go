package huaweicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/golangsdk/openstack/elb/v3/ipgroups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceIpGroupV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpGroupV3Create,
		Read:   resourceIpGroupV3Read,
		Update: resourceIpGroupV3Update,
		Delete: resourceIpGroupV3Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ip_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceIpGroupAddresses(d *schema.ResourceData) []ipgroups.IpListOpt {
	var IpList []ipgroups.IpListOpt
	ipListRaw := d.Get("ip_list").([]interface{})

	for _, v := range ipListRaw {
		ipList := v.(map[string]interface{})
		ipListOpts := ipgroups.IpListOpt{
			Ip:          ipList["ip"].(string),
			Description: ipList["description"].(string),
		}
		IpList = append(IpList, ipListOpts)
	}

	return IpList
}

func resourceIpGroupV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb v3 client: %s", err)
	}

	ipList := resourceIpGroupAddresses(d)
	desc := d.Get("description").(string)
	createOpts := ipgroups.CreateOpts{
		Name:                d.Get("name").(string),
		Description:         &desc,
		IpList:              &ipList,
		EnterpriseProjectID: GetEnterpriseProjectID(d, config),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	ig, err := ipgroups.Create(elbClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating IpGroup: %s", err)
	}
	d.SetId(ig.ID)

	return resourceIpGroupV3Read(d, meta)
}

func resourceIpGroupV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb v3 client: %s", err)
	}

	ig, err := ipgroups.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "ipgroup")
	}

	logp.Printf("[DEBUG] Retrieved ip group %s: %#v", d.Id(), ig)

	d.Set("name", ig.Name)
	d.Set("description", ig.Description)
	d.Set("region", GetRegion(d, config))

	ipList := make([]map[string]interface{}, len(ig.IpList))
	for i, ip := range ig.IpList {
		ipList[i] = map[string]interface{}{
			"ip":          ip.Ip,
			"description": ip.Description,
		}
	}
	d.Set("ip_list", ipList)

	return nil
}

func resourceIpGroupV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb v3 client: %s", err)
	}

	var updateOpts ipgroups.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("ip_list") {
		ipList := resourceIpGroupAddresses(d)
		updateOpts.IpList = &ipList
	}

	logp.Printf("[DEBUG] Updating ipgroup %s with options: %#v", d.Id(), updateOpts)
	_, err = ipgroups.Update(elbClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating HuaweiCloud elb ip group: %s", err)
	}

	return resourceIpGroupV3Read(d, meta)
}

func resourceIpGroupV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb v3 client: %s", err)
	}

	logp.Printf("[DEBUG] Deleting ip group %s", d.Id())
	if err = ipgroups.Delete(elbClient, d.Id()).ExtractErr(); err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud elb ip group: %s", err)
	}

	return nil
}
