package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/evs"
)

func getV3SnapshotResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "evs"
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating EVS client: %s", err)
	}

	return evs.GetV3SnapshotDetail(client, state.Primary.ID)
}

func TestAccV3Snapshot_basic(t *testing.T) {
	var (
		snapshot     interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_evsv3_snapshot.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&snapshot,
		getV3SnapshotResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3Snapshot_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Daily backup"),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "metadata.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "size"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccV3Snapshot_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"metadata", "force",
				},
			},
		},
	})
}

func testAccV3Snapshot_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name              = "%s"
  description       = "Created by acc test"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 12
}
`, rName)
}

func testAccV3Snapshot_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evsv3_snapshot" "test" {
  volume_id   = huaweicloud_evs_volume.test.id
  name        = "%[2]s"
  description = "Daily backup"

  metadata = {
    foo = "bar"
    key = "value"
  }
}
`, testAccV3Snapshot_base(rName), rName)
}

func testAccV3Snapshot_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evsv3_snapshot" "test" {
  volume_id   = huaweicloud_evs_volume.test.id
  name        = "%[2]s-update"
  description = ""

  metadata = {
    foo = "bar"
    key = "value"
  }
}
`, testAccV3Snapshot_base(rName), rName)
}
