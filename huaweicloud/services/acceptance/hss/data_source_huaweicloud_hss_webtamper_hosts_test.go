package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWebTamperHosts_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_webtamper_hosts.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceWebTamperHosts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.#"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.os_bit"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.protect_status"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.rasp_protect_status"),

					resource.TestCheckOutput("is_host_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_protect_status_filter_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testDataSourceWebTamperHosts_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_hss_webtamper_hosts" "test" {
  depends_on = [huaweicloud_hss_webtamper_protection.test]
}

# Filter using host ID.
locals {
  host_id = data.huaweicloud_hss_webtamper_hosts.test.hosts[0].id
}

data "huaweicloud_hss_webtamper_hosts" "host_id_filter" {
  host_id = local.host_id
}

output "is_host_id_filter_useful" {
  value = length(data.huaweicloud_hss_webtamper_hosts.host_id_filter.hosts) > 0 && alltrue(
    [for v in data.huaweicloud_hss_webtamper_hosts.host_id_filter.hosts[*].id : v == local.host_id]
  )
}

# Filter using name
locals {
  name = data.huaweicloud_hss_webtamper_hosts.test.hosts[0].name
}

data "huaweicloud_hss_webtamper_hosts" "name_filter" {
  name = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_hss_webtamper_hosts.name_filter.hosts) > 0 && alltrue(
    [for v in data.huaweicloud_hss_webtamper_hosts.name_filter.hosts[*].name : v == local.name]
  )
}

# Filter using protect_status
locals {
  protect_status = data.huaweicloud_hss_webtamper_hosts.test.hosts[0].protect_status
}

data "huaweicloud_hss_webtamper_hosts" "protect_status_filter" {
  protect_status = local.protect_status
}

output "is_protect_status_filter_filter_useful" {
  value = length(data.huaweicloud_hss_webtamper_hosts.protect_status_filter.hosts) > 0 && alltrue(
    [for v in data.huaweicloud_hss_webtamper_hosts.protect_status_filter.hosts[*].protect_status : v == local.protect_status]
  )
}

# Filter using non existent name.
data "huaweicloud_hss_webtamper_hosts" "not_found" {
  name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_webtamper_hosts.not_found.hosts) == 0
}
`, testAccWebTamperProtection_basic())
}
