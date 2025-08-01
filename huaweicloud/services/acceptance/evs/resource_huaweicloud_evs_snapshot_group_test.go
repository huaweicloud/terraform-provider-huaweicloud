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

func getSnapshotGroupResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "evs"
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating EVS client: %s", err)
	}

	return evs.GetSnapshotGroupDetail(client, state.Primary.ID)
}

func TestAccEvsSnapshotGroup_basic(t *testing.T) {
	var (
		snapshotGroup interface{}
		rName         = acceptance.RandomAccResourceName()
		resourceName  = "huaweicloud_evs_snapshot_group.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&snapshotGroup,
		getSnapshotGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckEVSServerID(t)
			acceptance.TestAccPreCheckEVSVolumeID(t) // The volume of id should attach to server id which used above
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsSnapshotGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Daily group backup"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "server_id", acceptance.HW_EVS_SERVER_ID),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "instant_access", "false"),
					resource.TestCheckResourceAttr(resourceName, "incremental", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccEvsSnapshotGroup_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "server_id", acceptance.HW_EVS_SERVER_ID),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "instant_access", "false"),
					resource.TestCheckResourceAttr(resourceName, "incremental", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"volume_ids", "instant_access", "incremental",
				},
			},
		},
	})
}

func testAccEvsSnapshotGroup_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_evs_snapshot_group" "test" {
  server_id      = "%[1]s"
  volume_ids     = ["%[2]s"]
  instant_access = false
  name           = "%[3]s"
  description    = "Daily group backup"
  tags = {
    foo = "bar"
  }
  enterprise_project_id = "%[4]s"
  incremental           = false
}
`, acceptance.HW_EVS_SERVER_ID, acceptance.HW_EVS_VOLUME_ID, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccEvsSnapshotGroup_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_evs_snapshot_group" "test" {
  server_id      = "%[1]s"
  volume_ids     = ["%[2]s"]
  instant_access = false
  name           = "%[3]s-update"
  description    = ""
  tags = {
    foo = "bar"
    key = "value"
  }
  enterprise_project_id = "%[4]s"
  incremental           = false
}
`, acceptance.HW_EVS_SERVER_ID, acceptance.HW_EVS_VOLUME_ID, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
