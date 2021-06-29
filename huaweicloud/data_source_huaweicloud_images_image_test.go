package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccImsImageDataSource_basic(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.huaweicloud_images_image.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImageDataSource_ubuntu(rName),
			},
			{
				Config: testAccImsImageDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "protected", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "visibility", "private"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "active"),
				),
			},
		},
	})
}

func TestAccImsImageDataSource_testQueries(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.huaweicloud_images_image.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImageDataSource_ubuntu(rName),
			},
			{
				Config: testAccImsImageDataSource_queryTag(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID(dataSourceName),
				),
			},
			{
				Config: testAccImsImageDataSource_querySizeMin(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID(dataSourceName),
				),
			},
			{
				Config: testAccImsImageDataSource_querySizeMax(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID(dataSourceName),
				),
			},
		},
	})
}

func testAccCheckImagesV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find image data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Image data source ID not set")
		}

		return nil
	}
}

func testAccImsImageDataSource_ubuntu(rName string) string {
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
  description = "created by Terraform AccTest"

  tags = {
    foo = "bar"
    key = "value"
  }
}

`, rName, rName)
}

func testAccImsImageDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_images_image" "test" {
  most_recent = true
  name        = huaweicloud_images_image.test.name
}
`, testAccImsImageDataSource_ubuntu(rName))
}

func testAccImsImageDataSource_queryTag(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_images_image" "test" {
  most_recent = true
  visibility  = "private"
  tag         = "foo=bar"
}
`, testAccImsImageDataSource_ubuntu(rName))
}

func testAccImsImageDataSource_querySizeMin(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_images_image" "test" {
  most_recent = true
  visibility  = "private"
  size_min    = "13000000"
}
`, testAccImsImageDataSource_ubuntu(rName))
}

func testAccImsImageDataSource_querySizeMax(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_images_image" "test" {
  most_recent = true
  visibility  = "private"
  size_max    = "23000000"
}
`, testAccImsImageDataSource_ubuntu(rName))
}
