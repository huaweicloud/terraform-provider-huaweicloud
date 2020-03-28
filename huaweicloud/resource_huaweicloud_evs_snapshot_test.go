package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/evs/v2/snapshots"
)

func TestAccEvsSnapshotV2_basic(t *testing.T) {
	var snapshot snapshots.Snapshot

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEvsSnapshotV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsSnapshotV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsSnapshotV2Exists("huaweicloud_evs_snapshot.snapshot_1", &snapshot),
					resource.TestCheckResourceAttr(
						"huaweicloud_evs_snapshot.snapshot_1", "name", "snapshot_acc"),
					resource.TestCheckResourceAttr(
						"huaweicloud_evs_snapshot.snapshot_1", "description", "Daily backup"),
					resource.TestCheckResourceAttr(
						"huaweicloud_evs_snapshot.snapshot_1", "status", "available"),
				),
			},
		},
	})
}

func testAccCheckEvsSnapshotV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	evsClient, err := config.blockStorageV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud EVS storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_evs_snapshot" {
			continue
		}

		_, err := snapshots.Get(evsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("EVS snapshot still exists")
		}
	}

	return nil
}

func testAccCheckEvsSnapshotV2Exists(n string, sp *snapshots.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		evsClient, err := config.blockStorageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Huaweicloud EVS storage client: %s", err)
		}

		found, err := snapshots.Get(evsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("EVS snapshot not found")
		}

		*sp = *found

		return nil
	}
}

const testAccEvsSnapshotV2_basic = `
resource "huaweicloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_acc"
  description = "volume for snapshot testing"
  size = 40
  cascade = true
}

resource "huaweicloud_evs_snapshot" "snapshot_1" {
  volume_id = huaweicloud_blockstorage_volume_v2.volume_1.id
  name = "snapshot_acc"
  description = "Daily backup"
}
`
