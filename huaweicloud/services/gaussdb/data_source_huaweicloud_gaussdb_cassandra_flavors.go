package gaussdb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/geminidb/v3/flavors"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceCassandraFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCassandraFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vcpus": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
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
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"az_status": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceCassandraFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.GeminiDBV31Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud GaussDB client: %s", err)
	}

	listOpts := flavors.ListFlavorOpts{
		EngineName: "cassandra",
	}

	pages, err := flavors.List(client, listOpts).AllPages()
	if err != nil {
		return fmtp.DiagErrorf("Unable to list flavors: %s", err)
	}

	allFlavors, err := flavors.ExtractFlavors(pages)
	if err != nil {
		return fmtp.DiagErrorf("Unable to extract flavors: %s", err)
	}

	filter := map[string]interface{}{
		"Vcpus":         d.Get("vcpus"),
		"Ram":           d.Get("memory"),
		"EngineVersion": d.Get("version"),
	}

	filterFlavors, err := utils.FilterSliceWithField(allFlavors.Flavors, filter)
	if err != nil {
		return fmtp.DiagErrorf("filter Gaussdb cassandra flavors failed: %s", err)
	}
	logp.Printf("filter %d Gaussdb cassandra flavors from %d through options %v", len(filterFlavors), len(allFlavors.Flavors), filter)

	var flavorsToSet []map[string]interface{}
	var flavorsIds []string

	for _, flavorInAll := range filterFlavors {
		flavor := flavorInAll.(flavors.Flavor)
		flavorToSet := map[string]interface{}{
			"vcpus":     flavor.Vcpus,
			"memory":    flavor.Ram,
			"name":      flavor.SpecCode,
			"version":   flavor.EngineVersion,
			"az_status": flavor.AzStatus,
		}

		flavorID := flavor.SpecCode
		flavorsIds = append(flavorsIds, flavorID)

		flavorsToSet = append(flavorsToSet, flavorToSet)
	}

	d.SetId(hashcode.Strings(flavorsIds))
	return diag.FromErr(d.Set("flavors", flavorsToSet))
}
