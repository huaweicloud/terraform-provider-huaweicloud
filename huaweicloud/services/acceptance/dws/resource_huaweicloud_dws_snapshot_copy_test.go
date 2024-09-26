package dws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
)

func getCopiedSnapshotFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS client: %s", err)
	}

	return dws.GetSnapshotById(client, state.Primary.ID)
}

func TestAccSnapshotCopy_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_dws_snapshot_copy.test"
		name  = acceptance.RandomAccResourceNameWithDash()

		rc = acceptance.InitResourceCheck(
			rName,
			&obj,
			getCopiedSnapshotFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsAutomatedSnapshot(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSnapshotCopy_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "snapshot_id", acceptance.HW_DWS_AUTOMATED_SNAPSHOT_ID),
					resource.TestCheckResourceAttr(rName, "description", "Copying a snapshot by terraform script"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSnapshotCopyImportStateFunc(rName),
			},
		},
	})
}

func testAccSnapshotCopyImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		if rs.Primary.Attributes["snapshot_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<snapshot_id>/<id>', but got '%s/%s'",
				rs.Primary.Attributes["snapshot_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["snapshot_id"], rs.Primary.ID), nil
	}
}

func testSnapshotCopy_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_snapshot_copy" "test" {
  snapshot_id = "%[1]s"
  name        = "%[2]s"
  description = "Copying a snapshot by terraform script"
}
`, acceptance.HW_DWS_AUTOMATED_SNAPSHOT_ID, name)
}
