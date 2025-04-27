package lts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHosts_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_lts_hosts.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byIds   = "data.huaweicloud_lts_hosts.filter_by_ids"
		dcByIds = acceptance.InitDataSourceCheck(byIds)

		byNames   = "data.huaweicloud_lts_hosts.filter_by_names"
		dcByNames = acceptance.InitDataSourceCheck(byNames)

		byIps   = "data.huaweicloud_lts_hosts.filter_by_ips"
		dcByIps = acceptance.InitDataSourceCheck(byIps)

		byStatus   = "data.huaweicloud_lts_hosts.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byVersion   = "data.huaweicloud_lts_hosts.filter_by_version"
		dcByVersion = acceptance.InitDataSourceCheck(byVersion)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLTSHostGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHosts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "hosts.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByIds.CheckResourceExists(),
					resource.TestCheckOutput("is_ids_filter_useful", "true"),
					dcByNames.CheckResourceExists(),
					resource.TestCheckOutput("is_names_filter_useful", "true"),
					dcByIps.CheckResourceExists(),
					resource.TestCheckOutput("is_ips_filter_useful", "true"),
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					dcByVersion.CheckResourceExists(),
					resource.TestCheckOutput("is_version_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byIds, "hosts.0.host_type"),
					resource.TestMatchResourceAttr(byIds, "hosts.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataSourceHosts_basic() string {
	return fmt.Sprintf(`
locals {
  host_ids = split(",", "%[1]s")
}

data "huaweicloud_lts_hosts" "test" {}

# Filter by host ID
data "huaweicloud_lts_hosts" "filter_by_ids" {
  host_id_list = local.host_ids
}

locals {
  ids_filter_result = [
    for v in data.huaweicloud_lts_hosts.filter_by_ids.hosts[*].host_id : v if contains(local.host_ids, v)
  ]
}

output "is_ids_filter_useful" {
  value = join(",", sort(local.ids_filter_result)) == join(",", sort(local.host_ids))
}

# Filter by host name
locals {
  host_names = data.huaweicloud_lts_hosts.filter_by_ids.hosts[*].host_name
}

data "huaweicloud_lts_hosts" "filter_by_names" {
  filter {
    host_name_list = local.host_names
  }
}

locals {
  names_filter_result = [
    for v in data.huaweicloud_lts_hosts.filter_by_names.hosts[*].host_name : v if contains(local.host_names, v)
  ]
}

output "is_names_filter_useful" {
  value = join(",", sort(local.names_filter_result)) == join(",", sort(local.host_names))
}

# Filter by host IP
locals {
  host_ips = data.huaweicloud_lts_hosts.filter_by_ids.hosts[*].host_ip
}

data "huaweicloud_lts_hosts" "filter_by_ips" {
  filter {
    host_ip_list = local.host_ips
  }
}

locals {
  ips_filter_result = [
    for v in data.huaweicloud_lts_hosts.filter_by_ips.hosts[*].host_ip : v if contains(local.host_ips, v)
  ]
}

output "is_ips_filter_useful" {
  value = join(",", sort(local.ips_filter_result)) == join(",", sort(local.host_ips))
}

# Filter by host status
locals {
  host_status = try(data.huaweicloud_lts_hosts.test.hosts[0].host_status, "")
}

data "huaweicloud_lts_hosts" "filter_by_status" {
  filter {
    host_status = local.host_status
  }
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_lts_hosts.filter_by_status.hosts[*].host_status : v == local.host_status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by host version
locals {
  host_version = try(data.huaweicloud_lts_hosts.test.hosts[0].host_version, "")
}

data "huaweicloud_lts_hosts" "filter_by_version" {
  filter {
    host_version = local.host_version
  }
}

locals {
  version_filter_result = [
    for v in data.huaweicloud_lts_hosts.filter_by_version.hosts[*].host_version : v == local.host_version
  ]
}

output "is_version_filter_useful" {
  value = length(local.version_filter_result) > 0 && alltrue(local.version_filter_result)
}
`, acceptance.HW_LTS_HOST_IDS)
}
