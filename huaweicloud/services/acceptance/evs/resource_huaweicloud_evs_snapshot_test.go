package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/evs/v2/snapshots"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/evs"
)

func getSnapshotResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "evs"
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating EVS client: %s", err)
	}

	return evs.GetSnapshotDetail(client, state.Primary.ID)
}

func TestAccEvsSnapshot_basic(t *testing.T) {
	var (
		snapshot     snapshots.Snapshot
		rName        = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
		resourceName = "huaweicloud_evs_snapshot.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&snapshot,
		getSnapshotResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsSnapshotV2_basic(rName),
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
				Config: testAccEvsSnapshotV2_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "description", "Daily backup update"),
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

func testAccEvsSnapshotV2_base(rName string) string {
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

func testAccEvsSnapshotV2_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_snapshot" "test" {
  volume_id   = huaweicloud_evs_volume.test.id
  name        = "%[2]s"
  description = "Daily backup"
  metadata    = {
    foo = "bar"
    key = "value"
  }
}
`, testAccEvsSnapshotV2_base(rName), rName)
}

func testAccEvsSnapshotV2_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_snapshot" "test" {
  volume_id   = huaweicloud_evs_volume.test.id
  name        = "%[2]s-update"
  description = "Daily backup update"
  metadata    = {
    foo = "bar"
    key = "value"
  }
}
`, testAccEvsSnapshotV2_base(rName), rName)
}
