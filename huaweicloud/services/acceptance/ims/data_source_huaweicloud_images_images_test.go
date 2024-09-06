package ims

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccImsImagesDataSource_basic(t *testing.T) {
	imageName := "CentOS 7.4 64bit"
	dataSourceName := "data.huaweicloud_images_images.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImagesDataSource_publicName(imageName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.name", imageName),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.status", "active"),
				),
			},
			{
				Config: testAccImsImagesDataSource_osVersion(imageName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.status", "active"),
				),
			},
			{
				Config: testAccImsImagesDataSource_nameRegex("^CentOS 7.4"),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.protected", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.visibility", "public"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.status", "active"),
				),
			},
			{
				Config: testAccImsImagesDataSource_market(),
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

func TestAccImsImagesDataSource_testQueries(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.huaweicloud_images_images.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImagesDataSource_base(rName),
			},
			{
				Config: testAccImsImagesDataSource_queryName(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.protected", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.visibility", "private"),
					resource.TestCheckResourceAttr(dataSourceName, "images.0.status", "active"),
				),
			},
			{
				Config: testAccImsImagesDataSource_queryTag(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccImsImagesDataSource_publicName(imageName string) string {
	return fmt.Sprintf(`
data "huaweicloud_images_images" "test" {
  name       = "%s"
  visibility = "public"
}
`, imageName)
}

func testAccImsImagesDataSource_nameRegex(regexp string) string {
	return fmt.Sprintf(`
data "huaweicloud_images_images" "test" {
  architecture = "x86"
  name_regex   = "%s"
  visibility   = "public"
}
`, regexp)
}

func testAccImsImagesDataSource_osVersion(osVersion string) string {
	return fmt.Sprintf(`
data "huaweicloud_images_images" "test" {
  architecture = "x86"
  os_version   = "%s"
  visibility   = "public"
}
`, osVersion)
}

func testAccImsImagesDataSource_market() string {
	return `
data "huaweicloud_images_images" "test" {
  os         = "CentOS"
  visibility = "market"
}
`
}

func testAccImsImagesDataSource_base(rName string) string {
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

func testAccImsImagesDataSource_queryName(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_images_images" "test" {
  name = huaweicloud_ims_ecs_system_image.test.name
}
`, testAccImsImagesDataSource_base(rName))
}

func testAccImsImagesDataSource_queryTag(rName string) string {
	return fmt.Sprintf(`
%s
data "huaweicloud_images_images" "test" {
  visibility = "private"
  tag        = "foo=bar"
}
`, testAccImsImagesDataSource_base(rName))
}
