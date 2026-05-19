package taurusdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapFlavors_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_htap_flavors.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTaurusDBHtapFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.instance_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.ram"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.az_status.%"),
					resource.TestCheckOutput("filter_by_version_is_useful", "true"),
					resource.TestCheckOutput("filter_by_specCode_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceTaurusDBHtapFlavors_basic() string {
	return `
data "huaweicloud_taurusdb_htap_flavors" "test" {
  engine_name            = "star-rocks"
  availability_zone_mode = "single"
}

output "flavors" {
  value = data.huaweicloud_taurusdb_htap_flavors.test.flavors
}

data "huaweicloud_taurusdb_htap_flavors" "filter_by_version" {
  engine_name            = "star-rocks"
  availability_zone_mode = "single"
  version_name           = "3.1.6.0"

  depends_on = [data.huaweicloud_taurusdb_htap_flavors.test]
}

output "filter_by_version_is_useful" {
  value = length(data.huaweicloud_taurusdb_htap_flavors.filter_by_version.flavors) > 0 && alltrue(
    [for flavor in data.huaweicloud_taurusdb_htap_flavors.filter_by_version.flavors : flavor.version_name == "3.1.6.0"]
  )
}

data "huaweicloud_taurusdb_htap_flavors" "filter_by_specCode" {
  engine_name            = "star-rocks"
  availability_zone_mode = "single"
  spec_code              = "gaussdb.sr-be.xlarge.x86.4"

  depends_on = [data.huaweicloud_taurusdb_htap_flavors.test]
}

output "filter_by_specCode_is_useful" {
  value = length(data.huaweicloud_taurusdb_htap_flavors.filter_by_specCode.flavors) > 0 && alltrue(
    [for flavor in data.huaweicloud_taurusdb_htap_flavors.filter_by_specCode.flavors : flavor.spec_code == "gaussdb.sr-be.xlarge.x86.4"]
  )
}
`
}
