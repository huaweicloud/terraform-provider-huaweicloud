package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSnapshotRollBack_basic(t *testing.T) {
	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy
	// method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// One-time action resource do not need to be checked and no processing is performed in the Read method.
				Config: testAccSnapshotRollBack_basic(),
			},
		},
	})
}

func testAccSnapshotRollBack_base() string {
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name              = "%[1]s"
  description       = "Created by acc test"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 12
}

resource "huaweicloud_evs_snapshot" "test" {
  volume_id   = huaweicloud_evs_volume.test.id
  name        = "%[1]s"
  description = "Daily backup"
  metadata    = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testAccSnapshotRollBack_basic() string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_evs_snapshot_rollback" "test" {
  volume_id   = huaweicloud_evs_volume.test.id
  snapshot_id = huaweicloud_evs_snapshot.test.id
  name        = huaweicloud_evs_volume.test.name
}
`, testAccSnapshotRollBack_base())
}
