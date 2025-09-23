package geminidb

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/geminidb/v3/flavors"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDBforNoSQL GET /v3.1/{project_id}/flavors
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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.GeminiDBV31Client(region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	listOpts := flavors.ListFlavorOpts{
		EngineName: "cassandra",
	}

	pages, err := flavors.List(client, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("unable to list flavors: %s", err)
	}

	allFlavors, err := flavors.ExtractFlavors(pages)
	if err != nil {
		return diag.Errorf("unable to extract flavors: %s", err)
	}

	filter := map[string]interface{}{
		"Vcpus":         d.Get("vcpus"),
		"Ram":           d.Get("memory"),
		"EngineVersion": d.Get("version"),
	}

	filterFlavors, err := utils.FilterSliceWithField(allFlavors.Flavors, filter)
	if err != nil {
		return diag.Errorf("filter GaussDB cassandra flavors failed: %s", err)
	}
	log.Printf("[DEBUG] filter %d GaussDB cassandra flavors from %d through options %v", len(filterFlavors), len(allFlavors.Flavors), filter)

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
	var mErr *multierror.Error
	mErr = multierror.Append(mErr,
		d.Set("flavors", flavorsToSet),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
