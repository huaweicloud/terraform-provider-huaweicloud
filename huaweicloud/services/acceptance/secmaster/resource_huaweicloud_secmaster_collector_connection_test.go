package secmaster

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/secmaster"
)

func getResourceCollectorConnectionFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return secmaster.GetCollectorConnectionById(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

// Due to the complexity of fields configuration and API behavior, only expected failure
// scenarios can be tested at present.
func TestAccResourceCollectorConnection_basic(t *testing.T) {
	var (
		rName  = "huaweicloud_secmaster_collector_connection.test"
		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceCollectorConnectionFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterTemplateId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceCollectorConnection_basic(),
				ExpectError: regexp.MustCompile(`参数不合法`),
			},
		},
	})
}

func testAccResourceCollectorConnection_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_collector_connection" "test" {
  workspace_id = "%[1]s"
  title        = "test_connection"
  name         = "test_connection"
  template_id  = "%[2]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_TEMPLATE_ID)
}
