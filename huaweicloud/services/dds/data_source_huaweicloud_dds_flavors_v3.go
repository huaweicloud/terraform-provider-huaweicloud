package dds

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/dds/v3/flavors"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API DDS GET /v3.1/{project_id}/flavors
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
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
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
						"engine_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
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
						"az_status": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"engine_versions": {
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

func dataSourceDDSFlavorV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	ddsClient, err := conf.NewServiceClient("ddsv31", region)
	if err != nil {
		return diag.Errorf("Error creating DDS v3.1 client: %s", err)
	}

	listOpts := flavors.ListOpts{
		Region:        region,
		EngineName:    d.Get("engine_name").(string),
		EngineVersion: d.Get("engine_version").(string),
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
			"engine_name":     item.EngineName,
			"spec_code":       item.SpecCode,
			"type":            item.Type,
			"vcpus":           item.Vcpus,
			"memory":          item.Ram,
			"az_status":       item.AzStatus,
			"engine_versions": item.EngineVersions,
		}
		flavorList = append(flavorList, flavor)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("flavors", flavorList),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("Error setting DDS instance fields: %s", err)
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
