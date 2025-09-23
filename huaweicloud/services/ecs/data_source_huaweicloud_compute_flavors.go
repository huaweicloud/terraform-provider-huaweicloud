package ecs

import (
	"context"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/flavors"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

// @API ECS GET /v1/{project_id}/cloudservers/flavors
func DataSourceEcsFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEcsFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"performance_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"generation": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cpu_core_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"memory_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"performance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"generation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_core_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceEcsFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	ecsClient, err := conf.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	listOpts := &flavors.ListOpts{
		AvailabilityZone: d.Get("availability_zone").(string),
	}

	pages, err := flavors.List(ecsClient, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}

	allFlavors, err := flavors.ExtractFlavors(pages)
	if err != nil {
		return diag.Errorf("unable to retrieve flavors: %s ", err)
	}

	cpu := d.Get("cpu_core_count").(int)
	mem := int64(d.Get("memory_size").(int)) * 1024
	pType := d.Get("performance_type").(string)
	gen := d.Get("generation").(string)
	sType := d.Get("storage_type").(string)

	var ids []string
	var filteredFlavors []map[string]interface{}
	for _, flavor := range allFlavors {
		vCpu, _ := strconv.Atoi(flavor.Vcpus)
		if cpu > 0 && vCpu != cpu {
			continue
		}

		if mem > 0 && flavor.Ram != mem {
			continue
		}

		if pType != "" && flavor.OsExtraSpecs.PerformanceType != pType {
			continue
		}

		if sType != "" && flavor.OsExtraSpecs.StorageType != sType {
			continue
		}

		if gen != "" && flavor.OsExtraSpecs.Generation != gen {
			continue
		}

		ids = append(ids, flavor.ID)

		flavorToSet := map[string]interface{}{
			"id":               flavor.ID,
			"performance_type": flavor.OsExtraSpecs.PerformanceType,
			"storage_type":     flavor.OsExtraSpecs.StorageType,
			"generation":       flavor.OsExtraSpecs.Generation,
			"cpu_core_count":   vCpu,
			"memory_size":      flavor.Ram / 1024,
		}
		filteredFlavors = append(filteredFlavors, flavorToSet)
	}

	d.SetId(hashcode.Strings(ids))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("ids", ids),
		d.Set("flavors", filteredFlavors),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
