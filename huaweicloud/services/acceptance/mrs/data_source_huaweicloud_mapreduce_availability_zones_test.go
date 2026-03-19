package mrs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAvailabilityZones_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_mapreduce_availability_zones.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byScope   = "data.huaweicloud_mapreduce_availability_zones.filter_by_scope"
		dcByScope = acceptance.InitDataSourceCheck(byScope)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAvailabilityZones_basic,
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "default_az_code"),
					resource.TestCheckResourceAttrSet(all, "support_physical_az_group"),
					resource.TestMatchResourceAttr(all, "available_zones.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "available_zones.0.id"),
					resource.TestCheckResourceAttrSet(all, "available_zones.0.az_id"),
					resource.TestCheckResourceAttrSet(all, "available_zones.0.az_code"),
					resource.TestCheckResourceAttrSet(all, "available_zones.0.az_name"),
					resource.TestCheckResourceAttrSet(all, "available_zones.0.status"),
					resource.TestCheckResourceAttrSet(all, "available_zones.0.region_id"),
					resource.TestCheckResourceAttrSet(all, "available_zones.0.az_type"),
					resource.TestCheckResourceAttrSet(all, "available_zones.0.az_category"),
					resource.TestCheckResourceAttrSet(all, "available_zones.0.charge_policy"),
					resource.TestCheckResourceAttr(all, "available_zones.0.az_tags.#", "1"),
					resource.TestCheckResourceAttrSet(all, "available_zones.0.az_tags.0.public_border_group"),
					// az_group_id, az_tags.mode, az_tags.alias currently are empty, not checked.
					// Filter by 'scope' parameter.
					dcByScope.CheckResourceExists(),
					resource.TestMatchResourceAttr(byScope, "available_zones.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
		},
	})
}

const testAccDataAvailabilityZones_basic = `
#Without any filter parameters.
data "huaweicloud_mapreduce_availability_zones" "test" {}

# Filter by 'scope' parameter.
data "huaweicloud_mapreduce_availability_zones" "filter_by_scope" {
  scope = "Center"
}
`
