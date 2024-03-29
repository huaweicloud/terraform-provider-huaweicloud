package dds

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/dds/v3/flavors"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API DDS GET /v3/{project_id}/flavors
func DataSourceDDSFlavorV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDDSFlavorV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DDS-Community", "DDS-Enhanced",
				}, true),
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"mongos", "shard", "config", "replica", "single",
				}, true),
			},
			"vcpus": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spec_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vcpus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDDSFlavorV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	ddsClient, err := conf.DdsV3Client(region)
	if err != nil {
		return diag.Errorf("Error creating DDS client: %s", err)
	}

	listOpts := flavors.ListOpts{
		Region:     region,
		EngineName: d.Get("engine_name").(string),
	}

	pages, err := flavors.List(ddsClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to retrieve flavors: %s", err)
	}

	allFlavors, err := flavors.ExtractFlavors(pages)
	if err != nil {
		return diag.Errorf("Unable to extract flavors: %s", err)
	}

	flavorList := make([]map[string]interface{}, 0)
	filterType := d.Get("type").(string)
	filterVcpus := d.Get("vcpus").(string)
	filterMemory := d.Get("memory").(string)

	for _, item := range allFlavors {
		if filterFlavor(item, filterType, filterVcpus, filterMemory) {
			continue
		}

		flavor := map[string]interface{}{
			"spec_code": item.SpecCode,
			"type":      item.Type,
			"vcpus":     item.Vcpus,
			"memory":    item.Ram,
		}
		flavorList = append(flavorList, flavor)
	}

	log.Printf("[DEBUG] extract %d/%d flavors by filters.", len(flavorList), len(allFlavors))
	if len(flavorList) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	d.SetId("dds flavors")
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("flavors", flavorList),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("Error setting dds instance fields: %s", err)
	}

	return nil
}

func filterFlavor(item flavors.Flavor, flavorType, vcpus, memory string) bool {
	if flavorType != "" && flavorType != item.Type {
		return true
	}
	if vcpus != "" && vcpus != item.Vcpus {
		return true
	}
	if memory != "" && memory != item.Ram {
		return true
	}

	return false
}
