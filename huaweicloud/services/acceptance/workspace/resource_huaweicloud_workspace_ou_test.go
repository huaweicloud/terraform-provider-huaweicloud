package workspace

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getOuFun(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("workspace", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace client: %s", err)
	}

	return workspace.GetOuByName(client, state.Primary.ID)
}

func parseOUNames() (ouName, ouUpdateName string) {
	ouNames := strings.Split(acceptance.HW_WORKSPACE_AD_SERVER_OU_NAMES, ",")
	if len(ouNames) < 2 {
		return "", ""
	}

	ouName = ouNames[0]
	ouUpdateName = ouNames[1]
	return
}

// Before using this test, please ensure that the Workspace service has been registered to the AD domain.
// OUs already created on the AD server and does not exist in the Workspace service.
func TestAccOu_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_workspace_ou.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getOuFun)

		ouName, ouUpdateName = parseOUNames()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceOUNames(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOu_step1(ouName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "ou_name", ouName),
					resource.TestCheckResourceAttrPair(rName, "domain",
						"data.huaweicloud_workspace_service.test", "ad_domain.0.name"),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
					// Check attributes.
					resource.TestCheckResourceAttrSet(rName, "ou_dn"),
					resource.TestCheckResourceAttrSet(rName, "domain_id"),
				),
			},
			{
				Config: testAccOu_step2(ouUpdateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "ou_name", ouUpdateName),
					resource.TestCheckResourceAttr(rName, "description", ""),
				),
			},
			// Import by ID.
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Import by name.
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOuImportStateFunc(rName),
			},
		},
	})
}

func testAccOu_step1(ouName string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_ou" "test" {
  ou_name     = "%[1]s"
  domain      = try(data.huaweicloud_workspace_service.test.ad_domain[0].name, "NOT_FOUND")
  description = "Created by terraform script"
}
`, ouName)
}

func testAccOu_step2(ouUpdateName string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_ou" "test" {
  ou_name = "%[1]s"
  domain  = try(data.huaweicloud_workspace_service.test.ad_domain[0].name, "NOT_FOUND")
}
`, ouUpdateName)
}

func testAccOuImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		ouName := rs.Primary.Attributes["ou_name"]
		if ouName == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<name>', but got '%s'", ouName)
		}
		return ouName, nil
	}
}
