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

func getEvsv5SnapshotResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "evs"
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating EVS client: %s", err)
	}

	return evs.GetEvsv5SnapshotDetail(client, state.Primary.ID)
}

func TestAccEvsv5Snapshot_basic(t *testing.T) {
	var (
		snapshot     interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_evsv5_snapshot.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&snapshot,
		getEvsv5SnapshotResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckEVSVolumeID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsv5Snapshot_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Daily backup"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "volume_id", acceptance.HW_EVS_VOLUME_ID),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "instant_access", "true"),
					resource.TestCheckResourceAttr(resourceName, "incremental", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "size"),
					resource.TestCheckResourceAttrSet(resourceName, "encrypted"),
					resource.TestCheckResourceAttrSet(resourceName, "category"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(resourceName, "snapshot_type"),
					resource.TestCheckResourceAttrSet(resourceName, "progress"),
				),
			},
			{
				Config: testAccEvsv5Snapshot_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated backup"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "volume_id", acceptance.HW_EVS_VOLUME_ID),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "instant_access", "false"),
					resource.TestCheckResourceAttr(resourceName, "incremental", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "size"),
					resource.TestCheckResourceAttrSet(resourceName, "encrypted"),
					resource.TestCheckResourceAttrSet(resourceName, "category"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(resourceName, "snapshot_type"),
					resource.TestCheckResourceAttrSet(resourceName, "progress"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEvsv5Snapshot_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_evsv5_snapshot" "test" {
  volume_id             = "%[1]s"
  name                  = "%[2]s"
  description           = "Daily backup"
  enterprise_project_id = "%[3]s"
  instant_access        = true
  incremental           = false

  tags = {
    foo = "bar"
  }
}
`, acceptance.HW_EVS_VOLUME_ID, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccEvsv5Snapshot_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_evsv5_snapshot" "test" {
  volume_id             = "%[1]s"
  name                  = "%[2]s-update"
  description           = "Updated backup"
  enterprise_project_id = "%[3]s"
  instant_access        = false
  incremental           = false

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, acceptance.HW_EVS_VOLUME_ID, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
