package geminidb

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/geminidb/v3/flavors"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDBforNoSQL GET /v3.1/{project_id}/flavors
func DataSourceGaussDBNoSQLFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"cassandra", "mongodb", "influxdb", "redis",
				}, false),
				Default: "cassandra",
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"availability_zone": {
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
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

// This method is used to filter all available gaussdb flavors by ram, vcpus, engine version and availability zone.
func filterAvailableFlavors(d *schema.ResourceData, flavorList []flavors.Flavor) ([]map[string]interface{}, []string, error) {
	var names []string
	result := make([]map[string]interface{}, 0)

	// build filter by ram, vcpus and engine version
	filter := map[string]interface{}{
		"EngineVersion": d.Get("engine_version").(string),
	}
	if ram, ok := d.GetOk("memory"); ok {
		filter["Ram"] = strconv.Itoa(ram.(int))
	}
	if vcpus, ok := d.GetOk("vcpus"); ok {
		filter["Vcpus"] = strconv.Itoa(vcpus.(int))
	}

	filterFlavors, err := utils.FilterSliceWithField(flavorList, filter)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("[DEBUG] the filter is %v and the result is %v", filter, filterFlavors)

	// filter by availability zone
	az := d.Get("availability_zone").(string)
	for _, item := range filterFlavors {
		flavor := item.(flavors.Flavor)

		var azList []string
		for k, v := range flavor.AzStatus {
			if v == "normal" && (az == "" || (az != "" && az == k)) {
				azList = append(azList, k)
			}
		}
		result = append(result, map[string]interface{}{
			"name":               flavor.SpecCode,
			"vcpus":              flavor.Vcpus,
			"memory":             flavor.Ram,
			"engine":             flavor.EngineName,
			"engine_version":     flavor.EngineVersion,
			"availability_zones": azList,
		})
		names = append(names, flavor.SpecCode)
	}

	log.Printf("[DEBUG] after filtering, the NoSQL flavors is %v", result)
	return result, names, nil
}

func dataSourceGaussDBFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.GeminiDBV31Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	opt := flavors.ListFlavorOpts{
		EngineName: d.Get("engine").(string),
	}
	pages, err := flavors.List(client, opt).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}
	flavorsResp, err := flavors.ExtractFlavors(pages)
	if err != nil {
		return diag.Errorf("unable to retrieve GuassDB NoSQL flavors: %s", err)
	}

	result, names, err := filterAvailableFlavors(d, flavorsResp.Flavors)
	if err != nil {
		return diag.Errorf("an error occurred while filtering the GaussDB NoSQL flavors: %s", err)
	}
	d.SetId(hashcode.Strings(names))
	var mErr *multierror.Error
	mErr = multierror.Append(mErr,
		d.Set("flavors", result),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
