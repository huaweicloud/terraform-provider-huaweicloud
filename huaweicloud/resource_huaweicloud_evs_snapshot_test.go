package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/evs/v2/snapshots"
)

func TestAccEvsSnapshotV2_basic(t *testing.T) {
	var snapshot snapshots.Snapshot

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_evs_snapshot.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEvsSnapshotV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsSnapshotV2_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEvsSnapshotV2Exists(resourceName, &snapshot),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Daily backup"),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
				),
			},
		},
	})
}

func testAccCheckEvsSnapshotV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	evsClient, err := config.BlockStorageV2Client(HW_REGION_NAME)
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
		evsClient, err := config.BlockStorageV2Client(HW_REGION_NAME)
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

func testAccEvsSnapshotV2_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_evs_snapshot" "test" {
  volume_id   = huaweicloud_evs_volume.test.id
  name        = "%s"
  description = "Daily backup"
}
`, testAccEvsStorageV3Volume_basic(rName), rName)
}
