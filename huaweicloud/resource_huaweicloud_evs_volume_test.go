package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/evs/v3/volumes"
)

func TestAccEvsVolume_basic(t *testing.T) {
	var volume volumes.Volume

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_evs_volume.test"
	rNameUpdate := rName + "-updated"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEvsVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsVolumeExists(resourceName, &volume),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "size", "12"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cascade",
				},
			},
			{
				Config: testAccEvsVolume_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsVolumeExists(resourceName, &volume),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "size", "20"),
				),
			},
		},
	})
}

func TestAccEvsVolume_tags(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_evs_volume.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEvsVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_tags(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsVolumeTags(resourceName, "foo", "bar"),
					testAccCheckEvsVolumeTags(resourceName, "key", "value"),
				),
			},
			{
				Config: testAccEvsVolume_tags_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsVolumeTags(resourceName, "foo2", "bar2"),
					testAccCheckEvsVolumeTags(resourceName, "key2", "value2"),
				),
			},
		},
	})
}

func TestAccEvsVolume_image(t *testing.T) {
	var volume volumes.Volume

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_evs_volume.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEvsVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_image(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsVolumeExists(resourceName, &volume),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccEvsVolume_withEpsId(t *testing.T) {
	var volume volumes.Volume

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_evs_volume.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEvsVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_epsID(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsVolumeExists(resourceName, &volume),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckEvsVolumeDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud evs storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_evs_volume" {
			continue
		}

		_, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Volume still exists")
		}
	}

	return nil
}

func testAccCheckEvsVolumeExists(n string, volume *volumes.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		blockStorageClient, err := config.BlockStorageV3Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud evs storage client: %s", err)
		}

		found, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Volume not found")
		}

		*volume = *found

		return nil
	}
}

func testAccCheckEvsVolumeTags(n string, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		blockStorageClient, err := config.BlockStorageV3Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud block storage client: %s", err)
		}

		found, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Volume not found")
		}

		if found.Tags == nil {
			return fmt.Errorf("No Tags")
		}

		for key, value := range found.Tags {
			if k != key {
				continue
			}

			if v == value {
				return nil
			}
			return fmt.Errorf("Bad value for %s: %s", k, value)
		}
		return fmt.Errorf("Tag not found: %s", k)
	}
}

func testAccEvsVolume_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name              = "%s"
  description       = "test volume"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 12
}
`, rName)
}

func testAccEvsVolume_update(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name              = "%s"
  description       = "test volume"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 20
}
`, rName)
}

func testAccEvsVolume_tags(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name              = "%s"
  description       = "test volume with tags"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 12

  tags = {
    foo = "bar"
	key = "value"
  }
}
`, rName)
}

func testAccEvsVolume_tags_update(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name              = "%s"
  description       = "test volume with tags"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 12

  tags = {
    foo2 = "bar2"
	key2 = "value2"
  }
}
`, rName)
}

func testAccEvsVolume_image(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_image_v2" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_evs_volume" "test" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 40
  image_id          = data.huaweicloud_images_image_v2.test.id
}
`, rName)
}

func testAccEvsVolume_epsID(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name                  = "%s"
  description           = "test volume for epsID"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  volume_type           = "SSD"
  size                  = 20
  enterprise_project_id = "%s"
}
`, rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}
