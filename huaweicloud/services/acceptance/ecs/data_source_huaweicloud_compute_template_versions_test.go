package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEcsComputeTemplateVersions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_compute_template_versions.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEcsComputeTemplateVersions_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "launch_template_versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "launch_template_versions.0.launch_template_id"),
					resource.TestCheckResourceAttrSet(dataSource, "launch_template_versions.0.version_id"),
					resource.TestCheckResourceAttrSet(dataSource, "launch_template_versions.0.version_number"),
					resource.TestCheckResourceAttrSet(dataSource, "launch_template_versions.0.version_description"),
					resource.TestCheckResourceAttrSet(dataSource, "launch_template_versions.0.template_data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "launch_template_versions.0.template_data.0.flavor_id"),
					resource.TestCheckResourceAttrSet(dataSource, "launch_template_versions.0.template_data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "launch_template_versions.0.template_data.0.description"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.availability_zone_id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.auto_recovery"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.os_profile.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.os_profile.0.key_name"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.os_profile.0.user_data"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.os_profile.0.iam_agency_name"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.os_profile.0.enable_monitoring_service"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.security_group_ids.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.network_interfaces.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.network_interfaces.0.virsubnet_id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.network_interfaces.0.attachment.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.network_interfaces.0.attachment.0.device_index"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.block_device_mappings.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.block_device_mappings.0.volume_type"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.block_device_mappings.0.volume_size"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.block_device_mappings.0.source_id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.block_device_mappings.0.source_type"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.block_device_mappings.0.encrypted"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.block_device_mappings.0.cmk_id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.block_device_mappings.0.attachment.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.block_device_mappings.0.attachment.0.boot_index"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.block_device_mappings.0.attachment.0.delete_on_termination"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.market_options.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.market_options.0.market_type"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.market_options.0.spot_options.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.market_options.0.spot_options.0.spot_price"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.market_options.0.spot_options.0.block_duration_minutes"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.market_options.0.spot_options.0.instance_interruption_behavior"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.internet_access.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.internet_access.0.publicip.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.internet_access.0.publicip.0.publicip_type"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.internet_access.0.publicip.0.charging_mode"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.internet_access.0.publicip.0.bandwidth.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.internet_access.0.publicip.0.bandwidth.0.share_type"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.internet_access.0.publicip.0.bandwidth.0.size"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.internet_access.0.publicip.0.bandwidth.0.charge_mode"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.internet_access.0.publicip.0.bandwidth.0.id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.internet_access.0.publicip.0.delete_on_termination"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.metadata.%"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.tag_options.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.tag_options.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.tag_options.0.tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource,
						"launch_template_versions.0.template_data.0.tag_options.0.tags.0.value"),
				),
			},
		},
	})
}

func testDataSourceEcsComputeTemplateVersions_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_compute_template_versions" "test" {
  depends_on = [huaweicloud_compute_template.test]
}

locals {
  test_versions = data.huaweicloud_compute_template_versions.test.launch_template_versions
}
data "huaweicloud_compute_template_versions" "launch_template_id_filter" {
  depends_on = [huaweicloud_compute_template.test]

  launch_template_id = local.test_versions[0].launch_template_id
}
locals {
  id_filter_template_versions = data.huaweicloud_compute_template_versions.launch_template_id_filter.launch_template_versions
  launch_template_id          = local.test_versions[0].launch_template_id
}
output "launch_template_id_filter_is_useful" {
  value = length(local.id_filter_template_versions) > 0 && alltrue(
  [for v in local.id_filter_template_versions[*].launch_template_id : v == local.launch_template_id]
  )
}

data "huaweicloud_compute_template_versions" "flavor_id_filter" {
  depends_on = [huaweicloud_compute_template.test]

  flavor_id = local.test_versions[0].template_data[0].flavor_id
}
locals {
  flavor_id_filter_template_versions = data.huaweicloud_compute_template_versions.flavor_id_filter.launch_template_versions
  flavor_id                          = local.test_versions[0].template_data[0].flavor_id
}
output "flavor_id_filter_is_useful" {
  value = length(local.flavor_id_filter_template_versions) > 0 && alltrue(
  [for v in local.flavor_id_filter_template_versions[*].template_data[0].flavor_id : v == local.flavor_id]
  )
}

data "huaweicloud_compute_template_versions" "version_filter" {
  depends_on = [huaweicloud_compute_template.test]

  version = [local.test_versions[0].version_number]
}
locals {
  version_filter_template_versions = data.huaweicloud_compute_template_versions.version_filter.launch_template_versions
  version_number                   = local.test_versions[0].version_number
}
output "version_filter_is_useful" {
  value = length(local.version_filter_template_versions) > 0 && alltrue(
  [for v in local.version_filter_template_versions[*].version_number : v == local.version_number]
  )
}
`, testAccComputeTemplate_basic(name))
}
