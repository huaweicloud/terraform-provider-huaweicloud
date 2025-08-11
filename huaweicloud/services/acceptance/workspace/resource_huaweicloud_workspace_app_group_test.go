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

func TestAccResourceAppGroup_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_app_group.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
		baseConfig   = testResourceWorkspaceAppGroup_base(name)

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
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceWorkspaceAppGroup_basic_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "SESSION_DESKTOP_APP"),
					resource.TestCheckResourceAttrPair(resourceName, "server_group_id", "huaweicloud_workspace_app_server_group.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created APP group by script"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testResourceWorkspaceAppGroup_basic_step2(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "server_group_id", ""),
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

func testResourceWorkspaceAppGroup_base(name string, appType ...string) string {
	actAppType := "SESSION_DESKTOP_APP"
	if len(appType) > 0 {
		actAppType = appType[0]
	}

	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_app_server_group" "test" {
  name             = "%[1]s"
  app_type         = "%[2]s"
  os_type          = "Windows"
  flavor_id        = "%[3]s"
  vpc_id           = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id        = data.huaweicloud_workspace_service.test.network_ids[0]
  system_disk_type = "SAS"
  system_disk_size = 80
  is_vdi           = true
  image_id         = "%[4]s"
  image_type       = "gold"
  image_product_id = "%[5]s"
}
`, name, actAppType,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID)
}

func testResourceWorkspaceAppGroup_basic_step1(baseConfig, name string, appType ...string) string {
	actAppType := "SESSION_DESKTOP_APP"
	if len(appType) > 0 {
		actAppType = appType[0]
	}

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_group" "test" {
  server_group_id = huaweicloud_workspace_app_server_group.test.id
  name            = "%[2]s"
  type            = "%[3]s"
  description     = "Created APP group by script"
}
`, baseConfig, name, actAppType)
}

func testResourceWorkspaceAppGroup_basic_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_group" "test" {
  name = "%[2]s"
}
`, baseConfig, name)
}
