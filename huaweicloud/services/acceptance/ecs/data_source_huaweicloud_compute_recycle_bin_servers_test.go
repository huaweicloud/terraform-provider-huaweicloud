package ecs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEcsComputeRecycleBinServers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_compute_recycle_bin_servers.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckECSID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEcsComputeRecycleBinServers_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "servers.#"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_scheduler_hints.#"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_scheduler_hints.0.group"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_scheduler_hints.0.tenancy"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_scheduler_hints.0.dedicated_host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.sys_tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.sys_tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.sys_tags.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_ext_sts_task_state"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_ext_sts_power_state"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_srv_usg_terminated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.progress"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.metadata.%"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_ext_srv_attr_root_device_name"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.locked"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.flavor.#"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.flavor.0.ram"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.flavor.0.gpus.#"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.flavor.0.gpus.0.count"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.flavor.0.gpus.0.memory_mb"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.flavor.0.gpus.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.flavor.0.asic_accelerators.#"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.flavor.0.asic_accelerators.0.memory_mb"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.flavor.0.asic_accelerators.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.flavor.0.asic_accelerators.0.count"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.flavor.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.flavor.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.flavor.0.disk"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.flavor.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.user_id"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.config_drive"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.security_options.#"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.security_options.0.secure_boot_enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.security_options.0.tpm_enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.hypervisor.#"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.hypervisor.0.hypervisor_type"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.hypervisor.0.csd_hypervisor"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_ext_srv_attr_hypervisor_hostname"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_ext_srv_attr_reservation_id"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_ext_srv_attr_launch_index"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_ext_srv_attr_instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.security_groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.security_groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.security_groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_extended_volumes_volumes_attached.#"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_extended_volumes_volumes_attached.0.id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"servers.0.os_extended_volumes_volumes_attached.0.delete_on_termination"),
					resource.TestCheckResourceAttrSet(dataSource,
						"servers.0.os_extended_volumes_volumes_attached.0.boot_index"),
					resource.TestCheckResourceAttrSet(dataSource,
						"servers.0.os_extended_volumes_volumes_attached.0.device"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.cpu_options.#"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.cpu_options.0.hw_cpu_threads"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_ext_az_availability_zone"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.host_status"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_ext_srv_attr_ramdisk_id"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_ext_srv_attr_kernel_id"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.updated"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.auto_terminate_time"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_ext_sts_vm_state"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_ext_srv_attr_hostname"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_ext_srv_attr_host"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_dcf_disk_config"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_srv_usg_launched_at"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.os_ext_srv_attr_user_data"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.key_name"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.image.#"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.image.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.created"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.network_interfaces.#"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.network_interfaces.0.port_id"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.network_interfaces.0.primary"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.network_interfaces.0.ip_addresses"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.network_interfaces.0.ipv6_addresses"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.network_interfaces.0.association.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"servers.0.network_interfaces.0.association.0.public_ip_address"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.network_interfaces.0.subnet_id"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("all_tenants_filter_is_useful", "true"),
					resource.TestCheckOutput("availability_zone_filter_is_useful", "true"),
					resource.TestCheckOutput("expect_fields_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceEcsComputeRecycleBinServers_basic() string {
	return `
data "huaweicloud_compute_recycle_bin_servers" "test" {}

locals {
  name = data.huaweicloud_compute_recycle_bin_servers.test.servers[0].name
}
data "huaweicloud_compute_recycle_bin_servers" "name_filter" {
  name = data.huaweicloud_compute_recycle_bin_servers.test.servers[0].name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_compute_recycle_bin_servers.name_filter.servers) > 0 alltrue(
    [for v in data.huaweicloud_compute_recycle_bin_servers.name_filter.servers[*].name : v == local.name]
  )
}

data "huaweicloud_compute_recycle_bin_servers" "all_tenants_filter" {
  all_tenants = "1"
}
output "all_tenants_filter_is_useful" {
  value = length(data.huaweicloud_compute_recycle_bin_servers.all_tenants_filter.servers) > 0
}

locals {
  availability_zone = data.huaweicloud_compute_recycle_bin_servers.test.servers[0].os_ext_az_availability_zone
}
data "huaweicloud_compute_recycle_bin_servers" "availability_zone_filter" {
  availability_zone = data.huaweicloud_compute_recycle_bin_servers.test.servers[0].os_ext_az_availability_zone
}
output "availability_zone_filter_is_useful" {
  value = length(data.huaweicloud_compute_recycle_bin_servers.availability_zone_filter.servers) > 0 alltrue(
    [for v in data.huaweicloud_compute_recycle_bin_servers.availability_zone_filter.servers[*].availability_zone :
      v == local.availability_zone]
  )
}

data "huaweicloud_compute_recycle_bin_servers" "expect_fields_filter" {
  expect_fields = "volumes_attached"
}
output "expect_fields_filter_is_useful" {
  value = length(data.huaweicloud_compute_recycle_bin_servers.expect_fields_filter.servers) > 0
}
`
}
