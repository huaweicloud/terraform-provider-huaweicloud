package elb

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/elb/v3/ipgroups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceIpGroupV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpGroupV3Create,
		ReadContext:   resourceIpGroupV3Read,
		UpdateContext: resourceIpGroupV3Update,
		DeleteContext: resourceIpGroupV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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

func resourceIpGroupV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating elb v3 client: %s", err)
	}

	ipList := resourceIpGroupAddresses(d)
	desc := d.Get("description").(string)
	createOpts := ipgroups.CreateOpts{
		Name:                d.Get("name").(string),
		Description:         &desc,
		IpList:              &ipList,
		EnterpriseProjectID: config.GetEnterpriseProjectID(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	ig, err := ipgroups.Create(elbClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating IpGroup: %s", err)
	}
	d.SetId(ig.ID)

	return resourceIpGroupV3Read(ctx, d, meta)
}

func resourceIpGroupV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating elb v3 client: %s", err)
	}

	ig, err := ipgroups.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "ipgroup")
	}

	log.Printf("[DEBUG] Retrieved ip group %s: %#v", d.Id(), ig)

	mErr := multierror.Append(nil,
		d.Set("name", ig.Name),
		d.Set("description", ig.Description),
		d.Set("region", config.GetRegion(d)),
	)

	ipList := make([]map[string]interface{}, len(ig.IpList))
	for i, ip := range ig.IpList {
		ipList[i] = map[string]interface{}{
			"ip":          ip.Ip,
			"description": ip.Description,
		}
	}
	d.Set("ip_list", ipList)

	mErr = multierror.Append(mErr, d.Set("ip_list", ipList))

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Dedicated ELB ipgroup fields: %s", err)
	}

	return nil
}

func resourceIpGroupV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating elb v3 client: %s", err)
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

	log.Printf("[DEBUG] Updating ipgroup %s with options: %#v", d.Id(), updateOpts)
	_, err = ipgroups.Update(elbClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating elb ip group: %s", err)
	}

	return resourceIpGroupV3Read(ctx, d, meta)
}

func resourceIpGroupV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating elb v3 client: %s", err)
	}

	log.Printf("[DEBUG] Deleting ip group %s", d.Id())
	if err = ipgroups.Delete(elbClient, d.Id()).ExtractErr(); err != nil {
		return diag.Errorf("error deleting elb ip group: %s", err)
	}

	return nil
}
