package workspace

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceWorkspaceAppGroupFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("appstream", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
	}

	uri := "/v1/{project_id}/app-groups/{app_group_id}"
	uri = strings.ReplaceAll(uri, "{app_group_id}", state.Primary.ID)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

// Before running this test, please create a workspace APP server group with SESSION_DESKTOP_APP type.
func TestAccResourceAppGroup_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_app_group.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()

		appGroup interface{}
		rc       = acceptance.InitResourceCheck(
			resourceName,
			&appGroup,
			getResourceWorkspaceAppGroupFunc,
		)
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
				Config: testResourceWorkspaceAppGroup_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "SESSION_DESKTOP_APP"),
					resource.TestCheckResourceAttr(resourceName, "server_group_id", acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID),
					resource.TestCheckResourceAttr(resourceName, "description", "Created APP group by script"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testResourceWorkspaceAppGroup_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "server_group_id", acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testResourceWorkspaceAppGroup_basic_step1(name string, appType ...string) string {
	actAppType := "SESSION_DESKTOP_APP"
	if len(appType) > 0 {
		actAppType = appType[0]
	}

	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_group" "test" {
  server_group_id = "%[1]s"
  name            = "%[2]s"
  type            = "%[3]s"
  description     = "Created APP group by script"
}
`, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID, name, actAppType)
}

func testResourceWorkspaceAppGroup_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_group" "test" {
  server_group_id = "%[1]s"
  name            = "%[2]s"
}
`, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID, name)
}
