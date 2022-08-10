package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/blockstorage/v2/volumes"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccBlockStorageV2Volume_basic(t *testing.T) {
	var volume volumes.Volume

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV2Volume_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV2VolumeExists("huaweicloud_blockstorage_volume_v2.volume_1", &volume),
					testAccCheckBlockStorageV2VolumeMetadata(&volume, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"huaweicloud_blockstorage_volume_v2.volume_1", "name", "volume_1"),
					resource.TestCheckResourceAttr(
						"huaweicloud_blockstorage_volume_v2.volume_1", "size", "1"),
				),
			},
			{
				ResourceName:      "huaweicloud_blockstorage_volume_v2.volume_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cascade",
				},
			},
			{
				Config: testAccBlockStorageV2Volume_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV2VolumeExists("huaweicloud_blockstorage_volume_v2.volume_1", &volume),
					testAccCheckBlockStorageV2VolumeMetadata(&volume, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"huaweicloud_blockstorage_volume_v2.volume_1", "name", "volume_1-updated"),
					resource.TestCheckResourceAttr(
						"huaweicloud_blockstorage_volume_v2.volume_1", "size", "2"),
				),
			},
		},
	})
}

func TestAccBlockStorageV2Volume_online_resize(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV2Volume_online_resize,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"huaweicloud_blockstorage_volume_v2.volume_1", "size", "1"),
				),
			},
			{
				Config: testAccBlockStorageV2Volume_online_resize_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"huaweicloud_blockstorage_volume_v2.volume_1", "size", "2"),
				),
			},
		},
	})
}

func TestAccBlockStorageV2Volume_image(t *testing.T) {
	var volume volumes.Volume

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV2Volume_image,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV2VolumeExists("huaweicloud_blockstorage_volume_v2.volume_1", &volume),
					resource.TestCheckResourceAttr(
						"huaweicloud_blockstorage_volume_v2.volume_1", "name", "volume_1"),
				),
			},
		},
	})
}

func TestAccBlockStorageV2Volume_timeout(t *testing.T) {
	var volume volumes.Volume

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV2VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV2Volume_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV2VolumeExists("huaweicloud_blockstorage_volume_v2.volume_1", &volume),
				),
			},
		},
	})
}

func testAccCheckBlockStorageV2VolumeDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	blockStorageClient, err := config.BlockStorageV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud block storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_blockstorage_volume_v2" {
			continue
		}

		_, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Volume still exists")
		}
	}

	return nil
}

func testAccCheckBlockStorageV2VolumeExists(n string, volume *volumes.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		blockStorageClient, err := config.BlockStorageV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud block storage client: %s", err)
		}

		found, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Volume not found")
		}

		*volume = *found

		return nil
	}
}

func testAccCheckBlockStorageV2VolumeDoesNotExist(t *testing.T, n string, volume *volumes.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		blockStorageClient, err := config.BlockStorageV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud block storage client: %s", err)
		}

		_, err = volumes.Get(blockStorageClient, volume.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}

		return fmtp.Errorf("Volume still exists")
	}
}

func testAccCheckBlockStorageV2VolumeMetadata(
	volume *volumes.Volume, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if volume.Metadata == nil {
			return fmtp.Errorf("No metadata")
		}

		for key, value := range volume.Metadata {
			if k != key {
				continue
			}

			if v == value {
				return nil
			}

			return fmtp.Errorf("Bad value for %s: %s", k, value)
		}

		return fmtp.Errorf("Metadata not found: %s", k)
	}
}

const testAccBlockStorageV2Volume_basic = `
resource "huaweicloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  size = 1
}
`

const testAccBlockStorageV2Volume_update = `
resource "huaweicloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1-updated"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  size = 2
}
`

// NOTE: Volume size cannot be smaller than the image minDisk size.
var testAccBlockStorageV2Volume_image = fmt.Sprintf(`
resource "huaweicloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 40
  image_id = "%s"
}
`, HW_IMAGE_ID)

const testAccBlockStorageV2Volume_timeout = `
resource "huaweicloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  size = 1

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

var testAccBlockStorageV2Volume_online_resize = fmt.Sprintf(`
resource "huaweicloud_compute_instance_v2" "basic" {
  name            = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  network {
    uuid = "%s"
  }
}

resource "huaweicloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "test volume"
  availability_zone = "%s"
  size = 1
}

resource "huaweicloud_compute_volume_attach_v2" "va_1" {
  instance_id = "${huaweicloud_compute_instance_v2.basic.id}"
  volume_id   = "${huaweicloud_blockstorage_volume_v2.volume_1.id}"
}
`, HW_AVAILABILITY_ZONE, HW_NETWORK_ID, HW_AVAILABILITY_ZONE)

var testAccBlockStorageV2Volume_online_resize_update = fmt.Sprintf(`
resource "huaweicloud_compute_instance_v2" "basic" {
  name            = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  network {
    uuid = "%s"
  }
}

resource "huaweicloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "test volume"
  availability_zone = "%s"
  size = 2
}

resource "huaweicloud_compute_volume_attach_v2" "va_1" {
  instance_id = "${huaweicloud_compute_instance_v2.basic.id}"
  volume_id   = "${huaweicloud_blockstorage_volume_v2.volume_1.id}"
}
`, HW_AVAILABILITY_ZONE, HW_NETWORK_ID, HW_AVAILABILITY_ZONE)
