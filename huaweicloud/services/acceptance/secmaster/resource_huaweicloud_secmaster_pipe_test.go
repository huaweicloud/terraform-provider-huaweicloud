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

func getResourcePipeFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return secmaster.GetPipeInfo(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccResourcePipe_basic(t *testing.T) {
	var (
		rName    = "huaweicloud_secmaster_pipe.test"
		pipeName = acceptance.RandomAccResourceName()

		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourcePipeFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterDataspaceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePipe_basic(pipeName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "pipe_name", pipeName),
					resource.TestCheckResourceAttr(rName, "shards", "3"),
					resource.TestCheckResourceAttr(rName, "storage_period", "30"),
					resource.TestCheckResourceAttrSet(rName, "pipe_type"),
					resource.TestCheckResourceAttrSet(rName, "dataspace_name"),
					resource.TestCheckResourceAttrSet(rName, "process_status"),
					resource.TestCheckResourceAttrSet(rName, "create_by"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_by"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
					resource.TestCheckResourceAttrSet(rName, "mapping"),
					resource.TestCheckResourceAttrSet(rName, "timestamp_field"),
				),
			},
			{
				Config: testAccResourcePipe_update(pipeName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "pipe_name", pipeName),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
					resource.TestCheckResourceAttr(rName, "shards", "2"),
					resource.TestCheckResourceAttr(rName, "storage_period", "20"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPipeImportStateFunc(rName),
			},
		},
	})
}

func testAccResourcePipe_basic(pipeName string) string {
	return fmt.Sprintf(`

resource "huaweicloud_secmaster_pipe" "test" {
  workspace_id    = "%[1]s"
  dataspace_id    = "%[2]s"
  pipe_name       = "%[3]s"
  shards          = 3
  storage_period  = 30
  description     = "test description"
  timestamp_field = "timestamp"

  mapping = jsonencode({
    id = {
      is_chinese_exist = true
      properties       = {}
      type             = "text"
    }
    name = {
      is_chinese_exist = false
      properties       = {}
      type             = "text"
    }
  })
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_DATASPACE_ID, pipeName)
}

func testAccResourcePipe_update(pipeName string) string {
	return fmt.Sprintf(`

resource "huaweicloud_secmaster_pipe" "test" {
  workspace_id    = "%[1]s"
  dataspace_id    = "%[2]s"
  pipe_name       = "%[3]s"
  shards          = 2
  storage_period  = 20
  description     = "test description update"
  timestamp_field = "timestamp"

  mapping = jsonencode({
    id = {
      is_chinese_exist = true
      properties       = {}
      type             = "text"
    }
    name = {
      is_chinese_exist = false
      properties       = {}
      type             = "text"
    }
  })
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_DATASPACE_ID, pipeName)
}

func testAccPipeImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var workspaceId, pipeId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rName)
		}

		workspaceId = rs.Primary.Attributes["workspace_id"]
		pipeId = rs.Primary.ID

		if workspaceId == "" || pipeId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<pipe_id>', but got '%s/%s'",
				workspaceId, pipeId)
		}

		return fmt.Sprintf("%s/%s", workspaceId, pipeId), nil
	}
}
