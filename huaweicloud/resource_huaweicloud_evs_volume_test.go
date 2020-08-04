package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/evs/v3/volumes"
)

func TestAccEvsStorageV3Volume_basic(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEvsStorageV3VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsStorageV3Volume_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsStorageV3VolumeExists("huaweicloud_evs_volume.volume_1", &volume),
					resource.TestCheckResourceAttr(
						"huaweicloud_evs_volume.volume_1", "name", "volume_1"),
				),
			},
			{
				Config: testAccEvsStorageV3Volume_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsStorageV3VolumeExists("huaweicloud_evs_volume.volume_1", &volume),
					resource.TestCheckResourceAttr(
						"huaweicloud_evs_volume.volume_1", "name", "volume_1-updated"),
				),
			},
		},
	})
}

func TestAccEvsStorageV3Volume_tags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEvsStorageV3VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsStorageV3Volume_tags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsStorageV3VolumeTags("huaweicloud_evs_volume.volume_tags", "foo", "bar"),
					testAccCheckEvsStorageV3VolumeTags("huaweicloud_evs_volume.volume_tags", "key", "value"),
				),
			},
			{
				Config: testAccEvsStorageV3Volume_tags_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsStorageV3VolumeTags("huaweicloud_evs_volume.volume_tags", "foo2", "bar2"),
					testAccCheckEvsStorageV3VolumeTags("huaweicloud_evs_volume.volume_tags", "key2", "value2"),
				),
			},
		},
	})
}

func TestAccEvsStorageV3Volume_image(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEvsStorageV3VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsStorageV3Volume_image,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsStorageV3VolumeExists("huaweicloud_evs_volume.volume_1", &volume),
					resource.TestCheckResourceAttr(
						"huaweicloud_evs_volume.volume_1", "name", "volume_1"),
				),
			},
		},
	})
}

func TestAccEvsStorageV3Volume_timeout(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEvsStorageV3VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsStorageV3Volume_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsStorageV3VolumeExists("huaweicloud_evs_volume.volume_1", &volume),
				),
			},
		},
	})
}

func testAccCheckEvsStorageV3VolumeDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	blockStorageClient, err := config.blockStorageV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud evs storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_evs_volume_v3" {
			continue
		}

		_, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Volume still exists")
		}
	}

	return nil
}

func testAccCheckEvsStorageV3VolumeExists(n string, volume *volumes.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		blockStorageClient, err := config.blockStorageV3Client(OS_REGION_NAME)
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

func testAccCheckEvsStorageV3VolumeTags(n string, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		blockStorageClient, err := config.blockStorageV3Client(OS_REGION_NAME)
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

var testAccEvsStorageV3Volume_basic = fmt.Sprintf(`
resource "huaweicloud_evs_volume" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  availability_zone = "%s"
  volume_type = "SAS"
  size = 12
}
`, OS_AVAILABILITY_ZONE)

var testAccEvsStorageV3Volume_update = fmt.Sprintf(`
resource "huaweicloud_evs_volume" "volume_1" {
  name = "volume_1-updated"
  description = "first test volume"
  availability_zone = "%s"
  volume_type = "SAS"
  size = 12
}
`, OS_AVAILABILITY_ZONE)

var testAccEvsStorageV3Volume_tags = fmt.Sprintf(`
resource "huaweicloud_evs_volume" "volume_tags" {
  name = "volume_tags"
  description = "test volume with tags"
  availability_zone = "%s"
  volume_type = "SAS"
  tags = {
    foo = "bar"
	key = "value"
  }
  size = 12
}
`, OS_AVAILABILITY_ZONE)

var testAccEvsStorageV3Volume_tags_update = fmt.Sprintf(`
resource "huaweicloud_evs_volume" "volume_tags" {
  name = "volume_tags-updated"
  description = "test volume with tags"
  availability_zone = "%s"
  volume_type = "SAS"
  tags = {
    foo2 = "bar2"
	key2 = "value2"
  }
  size = 12
}
`, OS_AVAILABILITY_ZONE)

var testAccEvsStorageV3Volume_image = fmt.Sprintf(`
resource "huaweicloud_evs_volume" "volume_1" {
  name = "volume_1"
  availability_zone = "%s"
  volume_type = "SAS"
  size = 12
  image_id = "%s"
}
`, OS_AVAILABILITY_ZONE, OS_IMAGE_ID)

var testAccEvsStorageV3Volume_timeout = fmt.Sprintf(`
resource "huaweicloud_evs_volume" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  availability_zone = "%s"
  size = 12
  volume_type = "SAS"
  device_type = "SCSI"
  timeouts {
    create = "10m"
    delete = "5m"
  }
}
`, OS_AVAILABILITY_ZONE)
