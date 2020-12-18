package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccHuaweiCloudImagesV2ImageDataSource_basic(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudImagesV2ImageDataSource_ubuntu(rName),
			},
			{
				Config: testAccHuaweiCloudImagesV2ImageDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.huaweicloud_images_image.test"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_images_image.test", "name", rName),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_images_image.test", "protected", "false"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_images_image.test", "visibility", "private"),
				),
			},
		},
	})
}

func TestAccHuaweiCloudImagesV2ImageDataSource_testQueries(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudImagesV2ImageDataSource_ubuntu(rName),
			},
			{
				Config: testAccHuaweiCloudImagesV2ImageDataSource_queryTag(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.huaweicloud_images_image.test"),
				),
			},
			{
				Config: testAccHuaweiCloudImagesV2ImageDataSource_querySizeMin(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.huaweicloud_images_image.test"),
				),
			},
			{
				Config: testAccHuaweiCloudImagesV2ImageDataSource_querySizeMax(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.huaweicloud_images_image.test"),
				),
			},
			{
				Config: testAccHuaweiCloudImagesV2ImageDataSource_ubuntu(rName),
			},
		},
	})
}

func testAccCheckImagesV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find image data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Image data source ID not set")
		}

		return nil
	}
}

func testAccHuaweiCloudImagesV2ImageDataSource_ubuntu(rName string) string {
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

resource "huaweicloud_compute_instance" "test" {
  name              = "%s"
  image_name        = "Ubuntu 18.04 server 64bit"
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  security_groups   = ["default"]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_images_image" "test" {
  name        = "%s"
  instance_id = huaweicloud_compute_instance.test.id
  description = "created by TerraformAccTest"

  tags = {
    foo = "bar"
    key = "value"
  }
}

`, rName, rName)
}

func testAccHuaweiCloudImagesV2ImageDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_images_image" "test" {
	most_recent = true
	name = huaweicloud_images_image.test.name
}
`, testAccHuaweiCloudImagesV2ImageDataSource_ubuntu(rName))
}

func testAccHuaweiCloudImagesV2ImageDataSource_queryTag(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_images_image" "test" {
	most_recent = true
	visibility = "private"
	tag = "foo=bar"
}
`, testAccHuaweiCloudImagesV2ImageDataSource_ubuntu(rName))
}

func testAccHuaweiCloudImagesV2ImageDataSource_querySizeMin(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_images_image" "test" {
	most_recent = true
	visibility = "private"
	size_min = "13000000"
}
`, testAccHuaweiCloudImagesV2ImageDataSource_ubuntu(rName))
}

func testAccHuaweiCloudImagesV2ImageDataSource_querySizeMax(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_images_image" "test" {
	most_recent = true
	visibility = "private"
	size_max = "23000000"
}
`, testAccHuaweiCloudImagesV2ImageDataSource_ubuntu(rName))
}
