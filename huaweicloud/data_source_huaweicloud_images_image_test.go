package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

func TestAccHuaweiCloudImagesV2ImageDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudImagesV2ImageDataSource_ubuntu,
			},
			{
				Config: testAccHuaweiCloudImagesV2ImageDataSource_basic,
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
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudImagesV2ImageDataSource_ubuntu,
			},
			{
				Config: testAccHuaweiCloudImagesV2ImageDataSource_queryTag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.huaweicloud_images_image.test"),
				),
			},
			{
				Config: testAccHuaweiCloudImagesV2ImageDataSource_querySizeMin,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.huaweicloud_images_image.test"),
				),
			},
			{
				Config: testAccHuaweiCloudImagesV2ImageDataSource_querySizeMax,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.huaweicloud_images_image.test"),
				),
			},
			{
				Config: testAccHuaweiCloudImagesV2ImageDataSource_ubuntu,
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

var testAccHuaweiCloudImagesV2ImageDataSource_ubuntu = fmt.Sprintf(`
%s

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

`, testAccCompute_data, rName, rName)

var testAccHuaweiCloudImagesV2ImageDataSource_basic = fmt.Sprintf(`
%s

data "huaweicloud_images_image" "test" {
	most_recent = true
	name = huaweicloud_images_image.test.name
}
`, testAccHuaweiCloudImagesV2ImageDataSource_ubuntu)

var testAccHuaweiCloudImagesV2ImageDataSource_queryTag = fmt.Sprintf(`
%s

data "huaweicloud_images_image" "test" {
	most_recent = true
	visibility = "private"
	tag = "foo=bar"
}
`, testAccHuaweiCloudImagesV2ImageDataSource_ubuntu)

var testAccHuaweiCloudImagesV2ImageDataSource_querySizeMin = fmt.Sprintf(`
%s

data "huaweicloud_images_image" "test" {
	most_recent = true
	visibility = "private"
	size_min = "13000000"
}
`, testAccHuaweiCloudImagesV2ImageDataSource_ubuntu)

var testAccHuaweiCloudImagesV2ImageDataSource_querySizeMax = fmt.Sprintf(`
%s

data "huaweicloud_images_image" "test" {
	most_recent = true
	visibility = "private"
	size_max = "23000000"
}
`, testAccHuaweiCloudImagesV2ImageDataSource_ubuntu)
