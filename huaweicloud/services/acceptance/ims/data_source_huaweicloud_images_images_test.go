package ims

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccImagesDataSource_basic(t *testing.T) {
	var (
		imageName      = "CentOS 7.4 64bit"
		dataSourceName = "data.huaweicloud_images_images.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesDataSource_public_queryByName(imageName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.name", imageName),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.status", "active"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.image_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.owner"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.os"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.architecture"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.os_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.file"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.schema"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.container_format"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.min_ram_mb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.max_ram_mb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.min_disk_gb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.disk_format"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.size_bytes"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.updated_at"),
				),
			},
			{
				Config: testAccImagesDataSource_public_queryByNameRegex("^CentOS 7.4"),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.status", "active"),
				),
			},
			{
				Config: testAccImagesDataSource_market(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.visibility", "market"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.status", "active"),
				),
			},
		},
	})
}

func TestAccImagesDataSource_private_basic(t *testing.T) {
	var (
		rName          = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
		dataSourceName = "data.huaweicloud_images_images.image_id_filter"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesDataSource_private_base(rName),
			},
			{
				Config: testAccImagesDataSource_private_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.visibility", "private"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.status", "active"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.image_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.owner"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.os"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.architecture"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.os_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.file"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.schema"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.container_format"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.min_ram_mb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.max_ram_mb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.min_disk_gb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.disk_format"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.data_origin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.size_bytes"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.active_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "images.0.updated_at"),
					resource.TestCheckOutput("is_image_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_flavor_id_filter_useful", "true"),
					resource.TestCheckOutput("is_tag_filter_useful", "true"),
					resource.TestCheckOutput("is_support_agent_list_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccImagesDataSource_public_queryByName(imageName string) string {
	return fmt.Sprintf(`
data "huaweicloud_images_images" "test" {
  name       = "%s"
  visibility = "public"
}
`, imageName)
}

func testAccImagesDataSource_public_queryByNameRegex(regexp string) string {
	return fmt.Sprintf(`
data "huaweicloud_images_images" "test" {
  architecture = "x86"
  name_regex   = "%s"
  visibility   = "public"
}
`, regexp)
}

func testAccImagesDataSource_market() string {
	return `
data "huaweicloud_images_images" "test" {
  os         = "CentOS"
  visibility = "market"
}
`
}

func testAccImagesDataSource_private_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_compute_instance" "test" {
  name       = "%[2]s"
  image_name = "Ubuntu 18.04 server 64bit"
  flavor_id  = data.huaweicloud_compute_flavors.test.ids[0]

  security_group_ids = [
    huaweicloud_networking_secgroup.test.id
  ]

  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_ims_ecs_system_image" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_compute_instance.test.id
  description = "created by Terraform AccTest"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccImagesDataSource_private_basic(rName string) string {
	return fmt.Sprintf(`
%s

# Filter using image ID.
locals {
  image_id = huaweicloud_ims_ecs_system_image.test.id
}

data "huaweicloud_images_images" "image_id_filter" {
  image_id = local.image_id
}

output "is_image_id_filter_useful" {
  value = length(data.huaweicloud_images_images.image_id_filter.images) > 0 && alltrue(
    [for v in data.huaweicloud_images_images.image_id_filter.images[*].id : v == local.image_id]
  )
}

# Filter using name.
locals {
  name = huaweicloud_ims_ecs_system_image.test.name
}

data "huaweicloud_images_images" "name_filter" {
  name = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_images_images.name_filter.images) > 0 && alltrue(
    [for v in data.huaweicloud_images_images.name_filter.images[*].name : v == local.name]
  )
}

# Filter using flavor_id.
locals {
  flavor_id = huaweicloud_compute_instance.test.flavor_id
}

data "huaweicloud_images_images" "flavor_id_filter" {
  flavor_id = local.flavor_id
}

output "is_flavor_id_filter_useful" {
  value = length(data.huaweicloud_images_images.flavor_id_filter.images) > 0
}

# Filter using tag.
data "huaweicloud_images_images" "tag_filter" {
  visibility = "private"
  tag        = "foo=bar"
}

output "is_tag_filter_useful" {
  value = length(data.huaweicloud_images_images.tag_filter.images) > 0
}

# Filter using __support_agent_list.
data "huaweicloud_images_images" "support_agent_list_filter" {
  visibility           = "private"
  __support_agent_list = "hss,ces"
}

output "is_support_agent_list_filter_useful" {
  value = length(data.huaweicloud_images_images.support_agent_list_filter.images) > 0
}

# Filter using non existent name.
data "huaweicloud_images_images" "not_found" {
  name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_images_images.not_found.images) == 0
}
`, testAccImagesDataSource_private_base(rName))
}
