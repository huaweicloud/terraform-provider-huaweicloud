package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dayu/v1/connections"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
)

func getDataConnectionFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.DataArtsV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio v1 client: %s", err)
	}

	resp, err := connections.Get(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
	return resp, dataarts.ParseDataConnectionError(err)
}

func TestAccResourceDataConnection_basic(t *testing.T) {
	var (
		dataConnection connections.Connection
		resourceName   = "huaweicloud_dataarts_studio_data_connection.test"
		name           = acceptance.RandomAccResourceName()
		updateName     = acceptance.RandomAccResourceName()
		rc             = acceptance.InitResourceCheck(resourceName, &dataConnection, getDataConnectionFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataConnection_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(resourceName, "type", "DLI"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "env_type", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "config"),
				),
			},
			{
				Config: testAccDataConnection_basic(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(resourceName, "type", "DLI"),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "env_type", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "config"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDataConnectionImportFunc(resourceName),
			},
		},
	})
}

func testAccDataConnectionImportFunc(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", n, rs)
		}
		workspaceId := rs.Primary.Attributes["workspace_id"]
		name := rs.Primary.Attributes["name"]
		if workspaceId == "" || name == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<name>', but got '%s/%s'", workspaceId, name)
		}
		return fmt.Sprintf("%s/%s", workspaceId, name), nil
	}
}

func testAccDataConnection_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_studio_data_connection" "test" {
  workspace_id = "%[1]s"
  type         = "DLI"
  name         = "%[2]s"
  env_type     = 0
  config       = jsonencode({
    "cdm_property_enable": "false"
  })
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}
