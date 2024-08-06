package css

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/css/v1/cluster"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CSS GET /v1.0/{project_id}/es-flavors
func DataSourceCssFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCssFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ess", "ess-cold", "ess-master", "ess-client"}, false),
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
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
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vcpus": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zones": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCssFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	flavorsResp, err := cluster.ListFlavors(client)
	if err != nil {
		return diag.Errorf("unable to retrieve CSS flavors: %s ", err)
	}

	allFlavors := flatternFlavors(flavorsResp)

	if len(allFlavors) < 1 {
		return diag.Errorf("no data found, please change your search criteria and try again")
	}

	filter := map[string]interface{}{
		"Region":  region,
		"Type":    d.Get("type"),
		"Name":    d.Get("name"),
		"Version": d.Get("version"),
	}

	if v, ok := d.GetOk("vcpus"); ok {
		filter["CPU"] = v
	}

	if v, ok := d.GetOk("memory"); ok {
		filter["RAM"] = v
	}

	filterFlavors, err := utils.FilterSliceWithField(allFlavors, filter)
	if err != nil {
		return diag.Errorf("filter CSS flavors failed: %s", err)
	}
	log.Printf("[DEBUG] filter %d CSS flavors from %d through options %v", len(filterFlavors), len(allFlavors), filter)

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("flavors", buildFlavors(filterFlavors)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flatternFlavors(flavors *cluster.EsFlavorsResp) []flavor {
	var rst []flavor
	for _, v := range flavors.Versions {
		for _, f := range v.Flavors {
			newFlavor := flavor{
				Type:              v.Type,
				Version:           v.Version,
				Name:              f.Name,
				FlavorId:          f.FlavorId,
				Region:            f.Region,
				RAM:               f.Ram,
				CPU:               f.Cpu,
				Diskrange:         f.Diskrange,
				AvailabilityZones: f.AvailableAZ,
			}

			rst = append(rst, newFlavor)
		}
	}

	return rst
}

func buildFlavors(flavors []interface{}) []map[string]interface{} {
	if len(flavors) < 1 {
		return nil
	}

	var rst []map[string]interface{}
	for _, v := range flavors {
		f := v.(flavor)
		newFlavor := make(map[string]interface{})

		newFlavor["id"] = f.FlavorId
		newFlavor["region"] = f.Region
		newFlavor["name"] = f.Name
		newFlavor["memory"] = f.RAM
		newFlavor["vcpus"] = f.CPU
		newFlavor["disk_range"] = f.Diskrange
		newFlavor["type"] = f.Type
		newFlavor["version"] = f.Version
		newFlavor["availability_zones"] = f.AvailabilityZones

		rst = append(rst, newFlavor)
	}

	return rst
}

type flavor struct {
	RAM               int
	CPU               int
	Name              string
	Region            string
	Diskrange         string
	FlavorId          string
	Version           string
	Type              string
	AvailabilityZones string
}
