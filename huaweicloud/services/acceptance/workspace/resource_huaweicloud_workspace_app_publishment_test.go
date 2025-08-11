package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getApplicationFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("appstream", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
	}
	return workspace.GetApplicationByName(client, state.Primary.Attributes["app_group_id"], state.Primary.Attributes["name"])
}

// Before running this test, please create a workspace APP server group with COMMON_APP type.
func TestAccAppPublishment_basic(t *testing.T) {
	var (
		application  interface{}
		resourceName = "huaweicloud_workspace_app_publishment.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
		baseConfig   = testDataSourceAppGroups_base(name, "COMMON_APP")
	)
	rc := acceptance.InitResourceCheck(
		resourceName,
		&application,
		getApplicationFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAppPublishment_basic_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "app_group_id", "huaweicloud_workspace_app_group.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "3"),
					resource.TestCheckResourceAttr(resourceName, "execute_path", "C:\\Program Files\\Sandboxie\\Start.exe"),
					resource.TestCheckResourceAttr(resourceName, "sandbox_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "version", "19.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "publisher", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "work_path", "C:\\Program Files\\Sandboxie"),
					resource.TestCheckResourceAttr(resourceName, "command_param", "/box:DefaultBox \"7567\""),
					resource.TestCheckResourceAttr(resourceName, "description", "Created APP by script"),
					resource.TestCheckResourceAttr(resourceName, "icon_path", "C:\\Program Files\\Sandboxie\\Start.exe"),
					resource.TestCheckResourceAttr(resourceName, "icon_index", "0"),
					resource.TestCheckResourceAttr(resourceName, "status", "FORBIDDEN"),
					resource.TestMatchResourceAttr(resourceName, "published_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`))),
			},
			{
				Config: testAccAppPublishment_basic_step2(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "type", "3"),
					resource.TestCheckResourceAttr(resourceName, "execute_path", "C:\\Program Files\\7-Zip\\7zFM.exe"),
					resource.TestCheckResourceAttr(resourceName, "sandbox_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "version", ""),
					resource.TestCheckResourceAttr(resourceName, "publisher", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "work_path", ""),
					resource.TestCheckResourceAttr(resourceName, "command_param", ""),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "icon_path", "C:\\Program Files\\Sandboxie\\Start.exe"),
					resource.TestCheckResourceAttr(resourceName, "icon_index", "0"),
					resource.TestCheckResourceAttr(resourceName, "status", "NORMAL"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAppPublishmentImportState(resourceName),
			},
		},
	})
}

func testAccAppPublishment_basic_step1(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_publishment" "test" {
  app_group_id   = huaweicloud_workspace_app_group.test.id
  name           = "%[2]s"
  type           = 3
  execute_path   = "C:\\Program Files\\Sandboxie\\Start.exe"
  sandbox_enable = true
  version        = "19.0.0.0"
  publisher      = "terraform"
  work_path      = "C:\\Program Files\\Sandboxie"
  command_param  = "/box:DefaultBox \"7567\""
  description    = "Created APP by script"
  icon_path      = "C:\\Program Files\\Sandboxie\\Start.exe"
  icon_index     = 0
  status         = "FORBIDDEN"
}
`, baseConfig, name)
}

func testAccAppPublishment_basic_step2(baseConfig, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_publishment" "test" {
  app_group_id   = huaweicloud_workspace_app_group.test.id
  name           = "%[2]s"
  type           = 3
  execute_path   = "C:\\Program Files\\7-Zip\\7zFM.exe"
  sandbox_enable = false
  publisher      = "terraform"
  icon_path      = "C:\\Program Files\\Sandboxie\\Start.exe"
  icon_index     = 0
  status         = "NORMAL"
}
`, baseConfig, updateName)
}

func testAppPublishmentImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		appGroupId := rs.Primary.Attributes["app_group_id"]
		appName := rs.Primary.Attributes["name"]
		if appGroupId == "" || appName == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<app_group_id>/<name>', but got '%s/%s'",
				appGroupId, appName)
		}

		return fmt.Sprintf("%s/%s", appGroupId, appName), nil
	}
}
