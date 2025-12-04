package secmaster

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getOperationConnectionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/soc/assetcredentials/{asset_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", state.Primary.Attributes["workspace_id"])
	requestPath = strings.ReplaceAll(requestPath, "{asset_id}", state.Primary.ID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SecMaster operation connection: %s", err)
	}

	return utils.FlattenResponse(resp)
}

func TestAccOperationConnection_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_secmaster_operation_connection.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOperationConnectionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterComponentId(t)
			acceptance.TestAccPreCheckSecMasterComponentVersionID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testOperationConnection_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "component_id", acceptance.HW_SECMASTER_COMPONENT_ID),
					resource.TestCheckResourceAttr(rName, "component_version_id", acceptance.HW_SECMASTER_COMPONENT_VERSION_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttrSet(rName, "project_id"),
					resource.TestCheckResourceAttrSet(rName, "component_name"),
					resource.TestCheckResourceAttrSet(rName, "type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "enabled"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "creator_id"),
					resource.TestCheckResourceAttrSet(rName, "creator_name"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
					resource.TestCheckResourceAttrSet(rName, "project_id"),
				),
			},
			{
				Config: testOperationConnection_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testOperationConnectionImportState(rName),
			},
		},
	})
}

func testOperationConnection_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_operation_connection" "test" {
  workspace_id         = "%[1]s"
  component_id         = "%[2]s"
  component_version_id = "%[3]s"
  config               = "{\"connection_type\":\"other\"}"
  name                 = "%[4]s"
  description          = "test description"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_COMPONENT_ID,
		acceptance.HW_SECMASTER_COMPONENT_VERSION_ID, name)
}

func testOperationConnection_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_operation_connection" "test" {
  workspace_id         = "%[1]s"
  component_id         = "%[2]s"
  component_version_id = "%[3]s"
  config               = "{\"connection_type\":\"other\"}"
  name                 = "%[4]s_update"
  description          = "test description update"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_COMPONENT_ID,
		acceptance.HW_SECMASTER_COMPONENT_VERSION_ID, name)
}

func testOperationConnectionImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		if rs.Primary.Attributes["workspace_id"] == "" {
			return "", fmt.Errorf("attribute (workspace_id) of resource (%s) not found", name)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["workspace_id"], rs.Primary.ID), nil
	}
}
