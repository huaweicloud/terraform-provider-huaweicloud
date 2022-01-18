package css

import (
	"context"

	"github.com/chnsz/golangsdk/openstack/css/v1/cluster"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

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
					},
				},
			},
		},
	}
}

func dataSourceCssFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.CssV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating CSS V1 client: %s", err)
	}

	flavorsResp, err := cluster.ListFlavors(client)
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve CSS flavors: %s ", err)
	}

	allFlavors := flatternFlavors(flavorsResp)

	if len(allFlavors) < 1 {
		return fmtp.DiagErrorf("No data found. Please change your search criteria and try again.")
	}

	filter := map[string]interface{}{
		"Region":  region,
		"Type":    d.Get("type"),
		"Name":    d.Get("name"),
		"Version": d.Get("version"),
	}

	if v, ok := d.GetOk("cpu"); ok {
		filter["Cpu"] = v
	}

	if v, ok := d.GetOk("ram"); ok {
		filter["Ram"] = v
	}

	filterFlavors, err := utils.FilterSliceWithField(allFlavors, filter)
	if err != nil {
		return fmtp.DiagErrorf("filter CSS flavors failed: %s", err)
	}
	logp.Printf("filter %d CSS flavors from %d through options %v", len(filterFlavors), len(allFlavors), filter)

	if len(filterFlavors) < 1 {
		return fmtp.DiagErrorf("No data found. Please change your search criteria and try again.")
	}

	mErr := d.Set("flavors", buildFlavors(filterFlavors))
	if mErr != nil {
		return fmtp.DiagErrorf("set flavors err:%s", mErr)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return fmtp.DiagErrorf("unable to generate ID:%s", err)
	}

	d.SetId(uuid)
	return nil
}

func flatternFlavors(flavors *cluster.EsFlavorsResp) []flavor {
	var rst []flavor
	for _, v := range flavors.Versions {
		for _, f := range v.Flavors {
			newFlavor := flavor{
				Type:      v.Type,
				Version:   v.Version,
				Name:      f.Name,
				FlavorId:  f.FlavorId,
				Region:    f.Region,
				Ram:       f.Ram,
				Cpu:       f.Cpu,
				Diskrange: f.Diskrange,
			}

			rst = append(rst, newFlavor)
		}
	}

	return rst
}

func buildFlavors(flavors []interface{}) []map[string]interface{} {
	var rst []map[string]interface{}
	for _, v := range flavors {
		f := v.(flavor)
		newFlavor := make(map[string]interface{})

		newFlavor["id"] = f.FlavorId
		newFlavor["region"] = f.Region
		newFlavor["name"] = f.Name
		newFlavor["memory"] = f.Ram
		newFlavor["vcpus"] = f.Cpu
		newFlavor["disk_range"] = f.Diskrange
		newFlavor["type"] = f.Type
		newFlavor["version"] = f.Version

		rst = append(rst, newFlavor)
	}

	return rst
}

type flavor struct {
	Ram       int
	Cpu       int
	Name      string
	Region    string
	Diskrange string
	FlavorId  string
	Version   string
	Type      string
}
