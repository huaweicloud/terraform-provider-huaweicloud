package dds

import (
	"context"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/dds/v3/flavors"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

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
	config := meta.(*config.Config)
	ddsClient, err := config.DdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud DDS client: %s", err)
	}

	listOpts := flavors.ListOpts{
		Region:     config.GetRegion(d),
		EngineName: d.Get("engine_name").(string),
	}

	pages, err := flavors.List(ddsClient, listOpts).AllPages()
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve flavors: %s", err)
	}

	allFlavors, err := flavors.ExtractFlavors(pages)
	if err != nil {
		return fmtp.DiagErrorf("Unable to extract flavors: %s", err)
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

	logp.Printf("Extract %d/%d flavors by filters.", len(flavorList), len(allFlavors))
	if len(flavorList) < 1 {
		return fmtp.DiagErrorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	d.SetId("dds flavors")
	mErr := multierror.Append(
		d.Set("region", config.GetRegion(d)),
		d.Set("flavors", flavorList),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting dds instance fields: %s", err)
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
