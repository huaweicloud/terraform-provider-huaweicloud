package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataBackups_basic(t *testing.T) {
	var (
		randName = acceptance.RandomAccResourceNameWithDash()

		dataSource = "data.huaweicloud_cbr_backups.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byCheckpointID   = "data.huaweicloud_cbr_backups.filter_by_checkpoint_id"
		dcByCheckpointID = acceptance.InitDataSourceCheck(byCheckpointID)

		byImageType   = "data.huaweicloud_cbr_backups.filter_by_image_type"
		dcByImageType = acceptance.InitDataSourceCheck(byImageType)

		byIncremental   = "data.huaweicloud_cbr_backups.filter_by_incremental"
		dcByIncremental = acceptance.InitDataSourceCheck(byIncremental)

		byName   = "data.huaweicloud_cbr_backups.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byResourceAz   = "data.huaweicloud_cbr_backups.filter_by_resource_az"
		dcByResourceAz = acceptance.InitDataSourceCheck(byResourceAz)

		byResourceID   = "data.huaweicloud_cbr_backups.filter_by_resource_id"
		dcByResourceID = acceptance.InitDataSourceCheck(byResourceID)

		byResourceName   = "data.huaweicloud_cbr_backups.filter_by_resource_name"
		dcByResourceName = acceptance.InitDataSourceCheck(byResourceName)

		byResourceType   = "data.huaweicloud_cbr_backups.filter_by_resource_type"
		dcByResourceType = acceptance.InitDataSourceCheck(byResourceType)

		byStatus   = "data.huaweicloud_cbr_backups.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byVault   = "data.huaweicloud_cbr_backups.filter_by_vault_id"
		dcByVault = acceptance.InitDataSourceCheck(byVault)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataBackups_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "backups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.checkpoint_id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.extend_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.image_type"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.incremental"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.protected_at"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.provider_id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.resource_az"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.resource_size"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.vault_id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.version"),

					dcByCheckpointID.CheckResourceExists(),
					resource.TestCheckOutput("checkpoint_id_filter_is_useful", "true"),

					dcByImageType.CheckResourceExists(),
					resource.TestCheckOutput("image_type_filter_is_useful", "true"),

					dcByIncremental.CheckResourceExists(),
					resource.TestCheckOutput("incremental_filter_is_useful", "true"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcByResourceAz.CheckResourceExists(),
					resource.TestCheckOutput("resource_az_filter_is_useful", "true"),

					dcByResourceID.CheckResourceExists(),
					resource.TestCheckOutput("resource_id_filter_is_useful", "true"),

					dcByResourceName.CheckResourceExists(),
					resource.TestCheckOutput("resource_name_filter_is_useful", "true"),

					dcByResourceType.CheckResourceExists(),
					resource.TestCheckOutput("resource_type_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),

					dcByVault.CheckResourceExists(),
					resource.TestCheckOutput("vault_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataBackups_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  data_disks {
    type = "SAS"
    size = "10"
  }
}

resource "huaweicloud_cbr_vault" "test" {
  name             = "%[2]s"
  type             = "server"
  consistent_level = "app_consistent"
  protection_type  = "backup"
  size             = 200
}

resource "huaweicloud_images_image" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_compute_instance.test.id
  vault_id    = huaweicloud_cbr_vault.test.id
}
`, common.TestBaseComputeResources(name), name)
}

// To avoid exceeding the API request frequency limit, each datasource adds `depends_on`.
func testAccDataBackups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cbr_backups" "test" {
  depends_on = [huaweicloud_images_image.test]
}

# Filter by checkpoint_id
locals {
  checkpoint_id = data.huaweicloud_cbr_backups.test.backups[0].checkpoint_id
}

data "huaweicloud_cbr_backups" "filter_by_checkpoint_id" {
  depends_on    = [data.huaweicloud_cbr_backups.test]
  checkpoint_id = local.checkpoint_id
}

locals {
  checkpoint_id_filter_result = [
    for v in data.huaweicloud_cbr_backups.filter_by_checkpoint_id.backups[*].checkpoint_id : v == local.checkpoint_id
  ]
}

output "checkpoint_id_filter_is_useful" {
  value = alltrue(local.checkpoint_id_filter_result) && length(local.checkpoint_id_filter_result) > 0
}

# Filter by image_type
locals {
  image_type = data.huaweicloud_cbr_backups.test.backups[0].image_type
}

data "huaweicloud_cbr_backups" "filter_by_image_type" {
  depends_on = [data.huaweicloud_cbr_backups.filter_by_checkpoint_id]
  image_type = local.image_type
}

locals {
  image_type_filter_result = [
    for v in data.huaweicloud_cbr_backups.filter_by_image_type.backups[*].image_type : v == local.image_type
  ]
}

output "image_type_filter_is_useful" {
  value = alltrue(local.image_type_filter_result) && length(local.image_type_filter_result) > 0
}

# Filter by incremental
locals {
  incremental = data.huaweicloud_cbr_backups.test.backups[0].incremental
}

data "huaweicloud_cbr_backups" "filter_by_incremental" {
  depends_on  = [data.huaweicloud_cbr_backups.filter_by_image_type]
  incremental = local.incremental
}

locals {
  incremental_filter_result = [
    for v in data.huaweicloud_cbr_backups.filter_by_incremental.backups[*].incremental : v == local.incremental
  ]
}

output "incremental_filter_is_useful" {
  value = alltrue(local.incremental_filter_result) && length(local.incremental_filter_result) > 0
}

# Filter by name
locals {
  name = data.huaweicloud_cbr_backups.test.backups[0].name
}

data "huaweicloud_cbr_backups" "filter_by_name" {
  depends_on = [data.huaweicloud_cbr_backups.filter_by_incremental]
  name       = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_cbr_backups.filter_by_name.backups[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

# Filter by resource_az
locals {
  resource_az = data.huaweicloud_cbr_backups.test.backups[0].resource_az
}

data "huaweicloud_cbr_backups" "filter_by_resource_az" {
  depends_on  = [data.huaweicloud_cbr_backups.filter_by_name]
  resource_az = local.resource_az
}

locals {
  resource_az_filter_result = [
    for v in data.huaweicloud_cbr_backups.filter_by_resource_az.backups[*].resource_az : v == local.resource_az
  ]
}

output "resource_az_filter_is_useful" {
  value = alltrue(local.resource_az_filter_result) && length(local.resource_az_filter_result) > 0
}

# Filter by resource_id
locals {
  resource_id = data.huaweicloud_cbr_backups.test.backups[0].resource_id
}

data "huaweicloud_cbr_backups" "filter_by_resource_id" {
  depends_on  = [data.huaweicloud_cbr_backups.filter_by_resource_az]
  resource_id = local.resource_id
}

locals {
  resource_id_filter_result = [
    for v in data.huaweicloud_cbr_backups.filter_by_resource_id.backups[*].resource_id : v == local.resource_id
  ]
}

output "resource_id_filter_is_useful" {
  value = alltrue(local.resource_id_filter_result) && length(local.resource_id_filter_result) > 0
}

# Filter by resource_name
locals {
  resource_name = data.huaweicloud_cbr_backups.test.backups[0].resource_name
}

data "huaweicloud_cbr_backups" "filter_by_resource_name" {
  depends_on    = [data.huaweicloud_cbr_backups.filter_by_resource_id]
  resource_name = local.resource_name
}

locals {
  resource_name_filter_result = [
    for v in data.huaweicloud_cbr_backups.filter_by_resource_name.backups[*].resource_name : v == local.resource_name
  ]
}

output "resource_name_filter_is_useful" {
  value = alltrue(local.resource_name_filter_result) && length(local.resource_name_filter_result) > 0
}

# Filter by resource_type
locals {
  resource_type = data.huaweicloud_cbr_backups.test.backups[0].resource_type
}

data "huaweicloud_cbr_backups" "filter_by_resource_type" {
  depends_on    = [data.huaweicloud_cbr_backups.filter_by_resource_name]
  resource_type = local.resource_type
}

locals {
  resource_type_filter_result = [
    for v in data.huaweicloud_cbr_backups.filter_by_resource_type.backups[*].resource_type : v == local.resource_type
  ]
}

output "resource_type_filter_is_useful" {
  value = alltrue(local.resource_type_filter_result) && length(local.resource_type_filter_result) > 0
}

# Filter by status
locals {
  status = data.huaweicloud_cbr_backups.test.backups[0].status
}

data "huaweicloud_cbr_backups" "filter_by_status" {
  depends_on = [data.huaweicloud_cbr_backups.filter_by_resource_type]
  status     = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_cbr_backups.filter_by_status.backups[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

# Filter by vault_id
locals {
  vault_id = data.huaweicloud_cbr_backups.test.backups[0].vault_id
}

data "huaweicloud_cbr_backups" "filter_by_vault_id" {
  depends_on = [data.huaweicloud_cbr_backups.filter_by_status]
  vault_id   = local.vault_id
}

locals {
  vault_id_filter_result = [
    for v in data.huaweicloud_cbr_backups.filter_by_vault_id.backups[*].vault_id : v == local.vault_id
  ]
}

output "vault_id_filter_is_useful" {
  value = alltrue(local.vault_id_filter_result) && length(local.vault_id_filter_result) > 0
}
`, testAccDataBackups_base(name), name)
}
