package elb

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/elb/v3/ipgroups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API ELB POST /v3/{project_id}/elb/ipgroups
// @API ELB GET /v3/{project_id}/elb/ipgroups/{ipgroup_id}
// @API ELB PUT /v3/{project_id}/elb/ipgroups/{ipgroup_id}
// @API ELB DELETE /v3/{project_id}/elb/ipgroups/{ipgroup_id}
func ResourceIpGroupV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpGroupV3Create,
		ReadContext:   resourceIpGroupV3Read,
		UpdateContext: resourceIpGroupV3Update,
		DeleteContext: resourceIpGroupV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"listener_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIpGroupAddresses(d *schema.ResourceData) []ipgroups.IpListOpt {
	var ipLists []ipgroups.IpListOpt
	ipListRaw := d.Get("ip_list").([]interface{})

	for _, v := range ipListRaw {
		ipList := v.(map[string]interface{})
		ipListOpt := ipgroups.IpListOpt{
			Ip:          ipList["ip"].(string),
			Description: ipList["description"].(string),
		}
		ipLists = append(ipLists, ipListOpt)
	}

	return ipLists
}

func resourceIpGroupV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	ipList := resourceIpGroupAddresses(d)
	desc := d.Get("description").(string)
	createOpts := ipgroups.CreateOpts{
		Name:                d.Get("name").(string),
		Description:         &desc,
		IpList:              &ipList,
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d),
	}

	log.Printf("[DEBUG] Create ELB IP Group options: %#v", createOpts)
	ipGroup, err := ipgroups.Create(elbClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating IpGroup: %s", err)
	}
	d.SetId(ipGroup.ID)

	return resourceIpGroupV3Read(ctx, d, meta)
}

func resourceIpGroupV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	ipGroup, err := ipgroups.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "ipgroup")
	}

	log.Printf("[DEBUG] Retrieved ip group %s: %#v", d.Id(), ipGroup)

	ipList := make([]map[string]interface{}, len(ipGroup.IpList))
	for i, ip := range ipGroup.IpList {
		ipList[i] = map[string]interface{}{
			"ip":          ip.Ip,
			"description": ip.Description,
		}
	}

	listenerIDs := make([]string, 0)
	for _, listener := range ipGroup.Listeners {
		listenerIDs = append(listenerIDs, listener.ID)
	}

	mErr := multierror.Append(nil,
		d.Set("name", ipGroup.Name),
		d.Set("description", ipGroup.Description),
		d.Set("region", cfg.GetRegion(d)),
		d.Set("ip_list", ipList),
		d.Set("listener_ids", listenerIDs),
		d.Set("created_at", ipGroup.CreatedAt),
		d.Set("updated_at", ipGroup.UpdatedAt),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Dedicated ELB ipgroup fields: %s", err)
	}

	return nil
}

func resourceIpGroupV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
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
		return diag.Errorf("error updating ELB ip group: %s", err)
	}

	return resourceIpGroupV3Read(ctx, d, meta)
}

func resourceIpGroupV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	log.Printf("[DEBUG] Deleting ip group %s", d.Id())
	if err = ipgroups.Delete(elbClient, d.Id()).ExtractErr(); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting ELB ip group")
	}

	return nil
}
