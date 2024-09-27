package ims

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccImageDataSource_basic(t *testing.T) {
	var (
		imageName      = "CentOS 7.4 64bit"
		dataSourceName = "data.huaweicloud_images_image.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImageDataSource_public_queryByName(imageName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", imageName),
					resource.TestCheckResourceAttr(dataSourceName, "protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "active"),
					resource.TestCheckResourceAttrSet(dataSourceName, "image_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "image_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "owner"),
					resource.TestCheckResourceAttrSet(dataSourceName, "os"),
					resource.TestCheckResourceAttrSet(dataSourceName, "architecture"),
					resource.TestCheckResourceAttrSet(dataSourceName, "os_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "file"),
					resource.TestCheckResourceAttrSet(dataSourceName, "schema"),
					resource.TestCheckResourceAttrSet(dataSourceName, "container_format"),
					resource.TestCheckResourceAttrSet(dataSourceName, "min_ram_mb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "max_ram_mb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "min_disk_gb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "disk_format"),
					resource.TestCheckResourceAttrSet(dataSourceName, "size_bytes"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "updated_at"),
				),
			},
			{
				Config: testAccImageDataSource_public_queryByNameRegex("^CentOS 7.4"),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "active"),
				),
			},
			{
				Config: testAccImageDataSource_market(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "visibility", "market"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "active"),
				),
			},
		},
	})
}

func TestAccImageDataSource_private(t *testing.T) {
	var (
		rName       = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
		byImageId   = "data.huaweicloud_images_image.image_id_filter"
		dcByImageId = acceptance.InitDataSourceCheck(byImageId)

		byName   = "data.huaweicloud_images_image.name_filter"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byFlavorId   = "data.huaweicloud_images_image.flavor_id_filter"
		dcByFlavorId = acceptance.InitDataSourceCheck(byFlavorId)

		byTag   = "data.huaweicloud_images_image.tag_filter"
		dcByTag = acceptance.InitDataSourceCheck(byTag)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImageDataSource_private_base(rName),
			},
			{
				Config: testAccImageDataSource_private_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dcByImageId.CheckResourceExists(),
					resource.TestCheckResourceAttr(byImageId, "name", rName),
					resource.TestCheckResourceAttr(byImageId, "visibility", "private"),
					resource.TestCheckResourceAttr(byImageId, "status", "active"),
					resource.TestCheckResourceAttrSet(byImageId, "image_id"),
					resource.TestCheckResourceAttrSet(byImageId, "image_type"),
					resource.TestCheckResourceAttrSet(byImageId, "owner"),
					resource.TestCheckResourceAttrSet(byImageId, "os"),
					resource.TestCheckResourceAttrSet(byImageId, "architecture"),
					resource.TestCheckResourceAttrSet(byImageId, "enterprise_project_id"),
					resource.TestCheckResourceAttrSet(byImageId, "os_version"),
					resource.TestCheckResourceAttrSet(byImageId, "file"),
					resource.TestCheckResourceAttrSet(byImageId, "schema"),
					resource.TestCheckResourceAttrSet(byImageId, "description"),
					resource.TestCheckResourceAttrSet(byImageId, "container_format"),
					resource.TestCheckResourceAttrSet(byImageId, "min_ram_mb"),
					resource.TestCheckResourceAttrSet(byImageId, "max_ram_mb"),
					resource.TestCheckResourceAttrSet(byImageId, "min_disk_gb"),
					resource.TestCheckResourceAttrSet(byImageId, "disk_format"),
					resource.TestCheckResourceAttrSet(byImageId, "data_origin"),
					resource.TestCheckResourceAttrSet(byImageId, "size_bytes"),
					resource.TestCheckResourceAttrSet(byImageId, "active_at"),
					resource.TestCheckResourceAttrSet(byImageId, "created_at"),
					resource.TestCheckResourceAttrSet(byImageId, "updated_at"),
					// Check whether filter parameter `name` is effective.
					dcByName.CheckResourceExists(),
					// Check whether filter parameter `flavor_id` is effective.
					dcByFlavorId.CheckResourceExists(),
					// Check whether filter parameter `tag` is effective.
					dcByTag.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccImageDataSource_public_queryByName(imageName string) string {
	return fmt.Sprintf(`
data "huaweicloud_images_image" "test" {
  name        = "%s"
  visibility  = "public"
  most_recent = true
}
`, imageName)
}

func testAccImageDataSource_public_queryByNameRegex(regexp string) string {
	return fmt.Sprintf(`
data "huaweicloud_images_image" "test" {
  name_regex   = "%s"
  visibility   = "public"
  architecture = "x86"
  most_recent  = true
}
`, regexp)
}

func testAccImageDataSource_market() string {
	return `
data "huaweicloud_images_image" "test" {
  visibility  = "market"
  os          = "CentOS"
  most_recent = true
}
`
}

func testAccImageDataSource_private_base(rName string) string {
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
  name               = "%[2]s"
  image_name         = "Ubuntu 18.04 server 64bit"
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

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

func testAccImageDataSource_private_basic(rName string) string {
	return fmt.Sprintf(`
%s

# Filter using image ID.
data "huaweicloud_images_image" "image_id_filter" {
  image_id = huaweicloud_ims_ecs_system_image.test.id
}

# Filter using name.
data "huaweicloud_images_image" "name_filter" {
  name        = huaweicloud_ims_ecs_system_image.test.name
  most_recent = true
}

# Filter using flavor_id.
data "huaweicloud_images_image" "flavor_id_filter" {
  flavor_id   = huaweicloud_compute_instance.test.flavor_id
  most_recent = true
}

# Filter using tag.
data "huaweicloud_images_image" "tag_filter" {
  most_recent = true
  visibility  = "private"
  tag         = "foo=bar"
}
`, testAccImageDataSource_private_base(rName))
}
