package dws

import (
	"context"
	"strconv"
	"strings"

	"github.com/chnsz/golangsdk/openstack/dws/v2/flavors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func DataSourceDwsFlavlors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDwsFlavorsRead,

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
			"vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flavor_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vcpus": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volumetype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDwsFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.DwsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating DWS V2 client: %s", err)
	}

	resp, rErr := flavors.ListNodeTypes(client)

	if rErr != nil {
		return fmtp.DiagErrorf("query the node type of dws flavors failed: %s", err)
	}

	az := d.Get("availability_zone").(string)
	cpu := d.Get("vcpus").(int)
	mem := d.Get("memory").(int)

	flavors := []dwsFlavor{}
	//filter flavors by arguments
	for _, node := range resp.NodeTypes {
		nodeTmp := parseNodeDetail(node)

		if cpu > 0 && nodeTmp.vcpus != cpu {
			continue
		}

		if mem > 0 && nodeTmp.memory != mem {
			continue
		}

		if az != "" {
			if !utils.StrSliceContains(nodeTmp.availabilityZones, az) {
				continue
			}
			nodeTmp.availabilityZones = []string{az}
		}

		flavors = append(flavors, nodeTmp)
	}

	var ids []string
	var resultFlavors []map[string]interface{}
	for _, item := range flavors {
		resultFlavors = append(resultFlavors, item.flattenDwsFlavor()...)
		ids = append(ids, item.id)
	}

	if len(resultFlavors) < 1 {
		return fmtp.DiagErrorf("Your query returned no results. Please change your search criteria and try again.")
	}

	logp.Printf("[DEBUG] Value of resultFlavors: %#v", resultFlavors)

	d.SetId(hashcode.Strings(ids))
	d.Set("region", config.GetRegion(d))
	d.Set("flavors", resultFlavors)

	return nil
}

type dwsFlavor struct {
	id                string
	specCode          string
	vcpus             int
	memory            int
	volumetype        string
	size              int
	availabilityZones []string
}

func parseNodeDetail(node flavors.NodeType) dwsFlavor {
	nodeTmp := dwsFlavor{
		id:       node.Id,
		specCode: node.SpecName,
	}
	for _, v := range node.Detail {
		switch v.Type {
		case "vCPU":
			nodeTmp.vcpus, _ = strconv.Atoi(v.Value)
		case "LOCAL_DISK", "SSD":
			nodeTmp.size, _ = strconv.Atoi(v.Value)
			nodeTmp.volumetype = v.Type
		case "mem":
			nodeTmp.memory, _ = strconv.Atoi(v.Value)
		case "availableZones":
			nodeTmp.availabilityZones = strings.Split(v.Value, ",")
		}
	}
	return nodeTmp
}

func (flavor *dwsFlavor) flattenDwsFlavor() []map[string]interface{} {
	if flavor == nil {
		return nil
	}
	azLength := len(flavor.availabilityZones)
	if azLength == 0 {
		return nil
	}
	var rt []map[string]interface{}
	for _, availableZone := range flavor.availabilityZones {
		newFlavor := map[string]interface{}{
			"flavor_id":         flavor.specCode,
			"vcpus":             flavor.vcpus,
			"memory":            flavor.memory,
			"volumetype":        flavor.volumetype,
			"size":              flavor.size,
			"availability_zone": availableZone,
		}
		rt = append(rt, newFlavor)
	}
	return rt
}
