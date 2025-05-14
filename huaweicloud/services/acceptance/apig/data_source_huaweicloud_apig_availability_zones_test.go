package apig

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAvailabilityZones_basic(t *testing.T) {
	dataSource := "data.huaweicloud_apig_availability_zones.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAvailabilityZones_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "availability_zones.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.code"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.port"),
					resource.TestCheckOutput("is_specs_set", "true"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.local_name.0.en_us"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.local_name.0.zh_cn"),
				),
			},
		},
	})
}

const testAccDataSourceAvailabilityZones_basic = `
data "huaweicloud_apig_availability_zones" "test" {}

output "is_specs_set" {
  value = length(try(data.huaweicloud_apig_availability_zones.test.availability_zones[0].specs, {})) > 0
}
`
