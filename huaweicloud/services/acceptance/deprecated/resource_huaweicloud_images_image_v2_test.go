package deprecated

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/imageservice/v2/images"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccImagesImageV2_basic(t *testing.T) {
	var image images.Image

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheckDeprecated(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("huaweicloud_images_image_v2.image_1", &image),
					resource.TestCheckResourceAttr(
						"huaweicloud_images_image_v2.image_1", "name", "Rancher TerraformAccTest"),
					resource.TestCheckResourceAttr(
						"huaweicloud_images_image_v2.image_1", "container_format", "bare"),
					resource.TestCheckResourceAttr(
						"huaweicloud_images_image_v2.image_1", "schema", "/v2/schemas/image"),
				),
			},
			{
				ResourceName:      "huaweicloud_images_image_v2.image_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"region",
					"local_file_path",
					"image_cache_path",
					"image_source_url",
				},
			},
		},
	})
}

func TestAccImagesImageV2_name(t *testing.T) {
	var image images.Image

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheckDeprecated(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_name_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("huaweicloud_images_image_v2.image_1", &image),
					resource.TestCheckResourceAttr(
						"huaweicloud_images_image_v2.image_1", "name", "Rancher TerraformAccTest"),
				),
			},
			{
				Config: testAccImagesImageV2_name_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("huaweicloud_images_image_v2.image_1", &image),
					resource.TestCheckResourceAttr(
						"huaweicloud_images_image_v2.image_1", "name", "TerraformAccTest Rancher"),
				),
			},
		},
	})
}

func TestAccImagesImageV2_tags(t *testing.T) {
	var image images.Image

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheckDeprecated(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_tags_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("huaweicloud_images_image_v2.image_1", &image),
					testAccCheckImagesImageV2HasTag("huaweicloud_images_image_v2.image_1", "foo"),
					testAccCheckImagesImageV2HasTag("huaweicloud_images_image_v2.image_1", "bar"),
					testAccCheckImagesImageV2TagCount("huaweicloud_images_image_v2.image_1", 2),
				),
			},
			{
				Config: testAccImagesImageV2_tags_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("huaweicloud_images_image_v2.image_1", &image),
					testAccCheckImagesImageV2HasTag("huaweicloud_images_image_v2.image_1", "foo"),
					testAccCheckImagesImageV2HasTag("huaweicloud_images_image_v2.image_1", "bar"),
					testAccCheckImagesImageV2HasTag("huaweicloud_images_image_v2.image_1", "baz"),
					testAccCheckImagesImageV2TagCount("huaweicloud_images_image_v2.image_1", 3),
				),
			},
			{
				Config: testAccImagesImageV2_tags_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("huaweicloud_images_image_v2.image_1", &image),
					testAccCheckImagesImageV2HasTag("huaweicloud_images_image_v2.image_1", "foo"),
					testAccCheckImagesImageV2HasTag("huaweicloud_images_image_v2.image_1", "baz"),
					testAccCheckImagesImageV2TagCount("huaweicloud_images_image_v2.image_1", 2),
				),
			},
		},
	})
}

func TestAccImagesImageV2_visibility(t *testing.T) {
	var image images.Image

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckDeprecated(t)
		},
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_visibility,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("huaweicloud_images_image_v2.image_1", &image),
					resource.TestCheckResourceAttr(
						"huaweicloud_images_image_v2.image_1", "visibility", "private"),
				),
			},
		},
	})
}

func TestAccImagesImageV2_timeout(t *testing.T) {
	var image images.Image

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheckDeprecated(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageV2Exists("huaweicloud_images_image_v2.image_1", &image),
				),
			},
		},
	})
}

func testAccCheckImagesImageV2Destroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	imageClient, err := cfg.ImageV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating IMS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_images_image_v2" {
			continue
		}

		_, err := images.Get(imageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("the image still exists, which ID is %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckImagesImageV2Exists(n string, image *images.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("the image %s not found", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		imageClient, err := cfg.ImageV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating IMS client: %s", err)
		}

		found, err := images.Get(imageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("the image is not found, which ID is %s", rs.Primary.ID)
		}

		*image = *found

		return nil
	}
}

func testAccCheckImagesImageV2HasTag(n, tag string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("the image %s not found", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		imageClient, err := cfg.ImageV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating IMS client: %s", err)
		}

		found, err := images.Get(imageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("the image is not found, which ID is %s", rs.Primary.ID)
		}

		for _, v := range found.Tags {
			if tag == v {
				return nil
			}
		}

		return fmt.Errorf("the tag is not found: %s", tag)
	}
}

func testAccCheckImagesImageV2TagCount(n string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("the image %s not found", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		imageClient, err := cfg.ImageV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating IMS client: %s", err)
		}

		found, err := images.Get(imageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("the image is not found, which ID is %s", rs.Primary.ID)
		}

		if len(found.Tags) != expected {
			return fmt.Errorf("expected %d tags, found %d", expected, len(found.Tags))
		}

		return nil
	}
}

var testAccImagesImageV2_basic = `
  resource "huaweicloud_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
      container_format = "bare"
      disk_format = "qcow2"
  }`

var testAccImagesImageV2_name_1 = `
  resource "huaweicloud_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
      container_format = "bare"
      disk_format = "qcow2"
  }`

var testAccImagesImageV2_name_2 = `
  resource "huaweicloud_images_image_v2" "image_1" {
      name   = "TerraformAccTest Rancher"
      image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
      container_format = "bare"
      disk_format = "qcow2"
  }`

var testAccImagesImageV2_tags_1 = `
  resource "huaweicloud_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
      container_format = "bare"
      disk_format = "qcow2"
      tags = ["foo","bar"]
  }`

var testAccImagesImageV2_tags_2 = `
  resource "huaweicloud_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
      container_format = "bare"
      disk_format = "qcow2"
      tags = ["foo","bar","baz"]
  }`

var testAccImagesImageV2_tags_3 = `
  resource "huaweicloud_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
      container_format = "bare"
      disk_format = "qcow2"
      tags = ["foo","baz"]
  }`

var testAccImagesImageV2_visibility = `
  resource "huaweicloud_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
      container_format = "bare"
      disk_format = "qcow2"
      visibility = "private"
  }`

var testAccImagesImageV2_timeout = `
  resource "huaweicloud_images_image_v2" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
      container_format = "bare"
      disk_format = "qcow2"

      timeouts {
        create = "10m"
      }
  }`
