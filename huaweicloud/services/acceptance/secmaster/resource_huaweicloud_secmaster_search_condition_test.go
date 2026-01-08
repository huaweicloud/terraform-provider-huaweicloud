package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/secmaster"
)

func getResourceSearchConditionFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return secmaster.GetSearchConditionInfo(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccResourceSearchCondition_basic(t *testing.T) {
	var (
		rName      = "huaweicloud_secmaster_search_condition.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceSearchConditionFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterDataspaceId(t)
			acceptance.TestAccPreCheckSecMasterPipeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSearchCondition_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "dataspace_id", acceptance.HW_SECMASTER_DATASPACE_ID),
					resource.TestCheckResourceAttr(rName, "pipe_id", acceptance.HW_SECMASTER_PIPE_ID),
					resource.TestCheckResourceAttr(rName, "condition_name", name),
					resource.TestCheckResourceAttr(rName, "query", "request_time=60"),
				),
			},
			{
				Config: testAccSearchCondition_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "condition_name", updateName),
					resource.TestCheckResourceAttr(rName, "query", "request_method=GET"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSearchConditionImportStateFunc(rName),
				ImportStateVerifyIgnore: []string{
					"dataspace_id",
				},
			},
		},
	})
}

func testAccSearchCondition_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_search_condition" "test" {
  workspace_id   = "%[1]s"
  dataspace_id   = "%[2]s"
  pipe_id        = "%[3]s"
  condition_name = "%[4]s"
  query          = "request_time=60"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_DATASPACE_ID, acceptance.HW_SECMASTER_PIPE_ID, name)
}

func testAccSearchCondition_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_search_condition" "test" {
  workspace_id   = "%[1]s"
  dataspace_id   = "%[2]s"
  pipe_id        = "%[3]s"
  condition_name = "%[4]s"
  query          = "request_method=GET"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_DATASPACE_ID, acceptance.HW_SECMASTER_PIPE_ID, name)
}

func testAccSearchConditionImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var workspaceId, conditionId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		workspaceId = rs.Primary.Attributes["workspace_id"]
		conditionId = rs.Primary.ID

		if workspaceId == "" || conditionId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s/%s'",
				workspaceId, conditionId)
		}

		return fmt.Sprintf("%s/%s", workspaceId, conditionId), nil
	}
}
