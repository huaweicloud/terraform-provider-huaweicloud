package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/addons"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceCCEAddonV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCceAddonsV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"addon_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCceAddonsV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	cceClient, err := config.CceAddonV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("unable to create HuaweiCloud CCE client : %s", err)
	}

	listOpts := addons.ListOpts{
		Uid:               d.Get("addon_id").(string),
		AddonTemplateName: d.Get("template_name").(string),
		Version:           d.Get("version").(string),
		Status:            d.Get("status").(string),
	}

	refinedAddons, err := addons.List(cceClient, d.Get("cluster_id").(string), listOpts)

	if err != nil {
		return fmt.Errorf("unable to retrieve Addons: %s", err)
	}

	if len(refinedAddons) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	if len(refinedAddons) > 1 {
		return fmt.Errorf("your query returned more than one result, please try a more specific search criteria")
	}

	Addon := refinedAddons[0]

	log.Printf("[DEBUG] Retrieved Nodes using given filter %s: %+v", Addon.Metadata.Id, Addon)
	d.SetId(Addon.Metadata.Id)
	d.Set("addon_id", Addon.Metadata.Id)
	d.Set("template_name", Addon.Spec.AddonTemplateName)
	d.Set("version", Addon.Spec.Version)
	d.Set("status", Addon.Status.Status)
	d.Set("region", GetRegion(d, config))
	d.Set("description", Addon.Spec.Description)

	return nil
}
