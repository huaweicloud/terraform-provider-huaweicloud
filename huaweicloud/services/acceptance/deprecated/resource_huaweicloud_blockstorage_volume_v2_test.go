package deprecated

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/blockstorage/v2/volumes"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBlockStorageV2Volume_basic(t *testing.T) {
	var volume volumes.Volume

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBlockStorageV2VolumeDestroy,
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
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBlockStorageV2VolumeDestroy,
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
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBlockStorageV2VolumeDestroy,
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
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBlockStorageV2VolumeDestroy,
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
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	blockStorageClient, err := cfg.BlockStorageV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating block storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_blockstorage_volume_v2" {
			continue
		}

		_, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("the block storage volume still exists, which ID is %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckBlockStorageV2VolumeExists(n string, volume *volumes.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("the block storage volume %s not found", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID set for the block storage volume %s", n)
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		blockStorageClient, err := cfg.BlockStorageV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating block storage client: %s", err)
		}

		found, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("the block storage volume is not found, which ID is %s", rs.Primary.ID)
		}

		*volume = *found

		return nil
	}
}

func testAccCheckBlockStorageV2VolumeMetadata(
	volume *volumes.Volume, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if volume.Metadata == nil {
			return errors.New("no metadata")
		}

		for key, value := range volume.Metadata {
			if k != key {
				continue
			}

			if v == value {
				return nil
			}

			return fmt.Errorf("bad value for %s: %s", k, value)
		}

		return fmt.Errorf("metadata not found: %s", k)
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
`, acceptance.HW_IMAGE_ID)

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
  availability_zone = "%[1]s"
  network {
    uuid = "%[2]s"
  }
}

resource "huaweicloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "test volume"
  availability_zone = "%[1]s"
  size = 1
}

resource "huaweicloud_compute_volume_attach_v2" "va_1" {
  instance_id = "${huaweicloud_compute_instance_v2.basic.id}"
  volume_id   = "${huaweicloud_blockstorage_volume_v2.volume_1.id}"
}
`, acceptance.HW_AVAILABILITY_ZONE, acceptance.HW_NETWORK_ID)

var testAccBlockStorageV2Volume_online_resize_update = fmt.Sprintf(`
resource "huaweicloud_compute_instance_v2" "basic" {
  name            = "instance_1"
  security_groups = ["default"]
  availability_zone = "%[1]s"
  network {
    uuid = "%[2]s"
  }
}

resource "huaweicloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "test volume"
  availability_zone = "%[1]s"
  size = 2
}

resource "huaweicloud_compute_volume_attach_v2" "va_1" {
  instance_id = "${huaweicloud_compute_instance_v2.basic.id}"
  volume_id   = "${huaweicloud_blockstorage_volume_v2.volume_1.id}"
}
`, acceptance.HW_AVAILABILITY_ZONE, acceptance.HW_NETWORK_ID)
