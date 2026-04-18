package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
)

func getSecurityWorkspaceQueueAssociateResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	workspaceId := state.Primary.Attributes["workspace_id"]
	queueName := state.Primary.Attributes["queue_name"]
	return dataarts.GetSecurityWorkspaceAssociatedQueueByName(client, workspaceId, queueName)
}

func TestAccResourceSecurityWorkspaceQueueAssociate_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_dataarts_security_workspace_queue_associate.test"

		rc = acceptance.InitResourceCheck(rName, &obj, getSecurityWorkspaceQueueAssociateResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsConnectionID(t)
			acceptance.TestAccPreCheckDataArtsRelatedDliQueueName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityWorkspaceQueueAssociate_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "connection_id", acceptance.HW_DATAARTS_CONNECTION_ID),
					resource.TestCheckResourceAttr(rName, "source_type", "dli"),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform"),
					resource.TestCheckResourceAttrSet(rName, "queue_type"),
					resource.TestCheckResourceAttrSet(rName, "connection_name"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "create_user"),
					resource.TestCheckResourceAttrSet(rName, "queue_attribute"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSecurityWorkspaceQueueAssociate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_security_workspace_queue_associate" "test" {
  workspace_id  = "%[1]s"
  source_type   = "dli"
  queue_name    = "%[2]s"
  connection_id = "%[3]s"
  description   = "Created by terraform"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_DLI_QUEUE_NAME, acceptance.HW_DATAARTS_CONNECTION_ID)
}
