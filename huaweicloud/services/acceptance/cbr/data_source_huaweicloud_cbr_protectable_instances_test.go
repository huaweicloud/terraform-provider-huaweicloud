package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceProtectableInstances_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceName()

		dataSourceServer = "data.huaweicloud_cbr_protectable_instances.test_server"
		dcServer         = acceptance.InitDataSourceCheck(dataSourceServer)

		dataSourceDisk = "data.huaweicloud_cbr_protectable_instances.test_disk"
		dcDisk         = acceptance.InitDataSourceCheck(dataSourceDisk)

		byResourceID   = "data.huaweicloud_cbr_protectable_instances.filter_by_resource_id"
		dcByResourceID = acceptance.InitDataSourceCheck(byResourceID)

		byName   = "data.huaweicloud_cbr_protectable_instances.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byStatus   = "data.huaweicloud_cbr_protectable_instances.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceProtectableInstances_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					// test server search attribute
					dcServer.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.children"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.detail"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.protectable.#"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.protectable.0.result"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.protectable.0.vault.#"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.protectable.0.vault.0.auto_bind"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.protectable.0.vault.0.auto_expand"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.protectable.0.vault.0.billing.#"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.protectable.0.vault.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.protectable.0.vault.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.protectable.0.vault.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.protectable.0.vault.0.locked"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.protectable.0.vault.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.protectable.0.vault.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.protectable.0.vault.0.provider_id"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.protectable.0.vault.0.resources.#"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.protectable.0.vault.0.smn_notify"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.protectable.0.vault.0.threshold"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.protectable.0.vault.0.user_id"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.size"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceServer, "instances.0.type"),

					// test disk search attribute
					dcDisk.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceDisk, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSourceDisk, "instances.0.children"),
					resource.TestCheckResourceAttrSet(dataSourceDisk, "instances.0.detail"),
					resource.TestCheckResourceAttrSet(dataSourceDisk, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceDisk, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceDisk, "instances.0.size"),
					resource.TestCheckResourceAttrSet(dataSourceDisk, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceDisk, "instances.0.type"),

					dcByResourceID.CheckResourceExists(),
					resource.TestCheckOutput("resource_id_filter_is_useful", "true"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceProtectableInstances_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[1]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  charging_mode      = "spot"
  spot_duration      = 2

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[1]s"
  type                  = "server"
  consistent_level      = "crash_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = "0"
  backup_name_prefix    = "test-prefix-"
  is_multi_az           = true

  resources {
    server_id = huaweicloud_compute_instance.test.id
  }
}
`, name)
}

func testDataSourceProtectableInstances_basic(name string) string {
	return fmt.Sprintf(`
%s

# There are some backfill field validations that require adding a repository before they can be performed.
data "huaweicloud_cbr_protectable_instances" "test_server" {
  depends_on = [huaweicloud_cbr_vault.test]

  server_id        = huaweicloud_compute_instance.test.id
  protectable_type = "server"
}

# After creating the server, there will be a disk that can be found.
data "huaweicloud_cbr_protectable_instances" "test_disk" {
  depends_on = [huaweicloud_compute_instance.test]

  protectable_type = "disk"
}

# Filter by resource_id
locals {
  resource_id = data.huaweicloud_cbr_protectable_instances.test_disk.instances[0].id
}

data "huaweicloud_cbr_protectable_instances" "filter_by_resource_id" {
  protectable_type = "disk"
  resource_id      = local.resource_id
}

locals {
  resource_id_filter_result = [
    for v in data.huaweicloud_cbr_protectable_instances.filter_by_resource_id.instances[*].id : v == local.resource_id
  ]
}

output "resource_id_filter_is_useful" {
  value = alltrue(local.resource_id_filter_result) && length(local.resource_id_filter_result) > 0
}

# Filter by name
locals {
  name = data.huaweicloud_cbr_protectable_instances.test_disk.instances[0].name
}

data "huaweicloud_cbr_protectable_instances" "filter_by_name" {
  protectable_type = "disk"
  name             = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_cbr_protectable_instances.filter_by_name.instances[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

# Filter by status
locals {
  status = data.huaweicloud_cbr_protectable_instances.test_disk.instances[0].status
}

data "huaweicloud_cbr_protectable_instances" "filter_by_status" {
  protectable_type = "disk"
  status           = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_cbr_protectable_instances.filter_by_status.instances[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}
`, testDataSourceProtectableInstances_base(name))
}
