package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getResourceWorkspaceAppGroupFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("appstream", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
	}

	uri := "/v1/{project_id}/app-groups"
	queryParam := map[string]any{
		"app_group_id": state.Primary.ID,
	}
	resp, err := httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(queryParam).
		Request().
		Data()

	groups := utils.PathSearch("items", resp, make([]interface{}, 0))
	if len(groups.([]interface{})) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return resp, err
}

func TestAccResourceWorkspaceAppGroup_basic(t *testing.T) {
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
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
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
					resource.TestCheckResourceAttrPair(resourceName, "server_group_id", "huaweicloud_workspace_app_server_group.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created APP group by script"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testResourceWorkspaceAppGroup_basic_step2(name, updateName),
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

func testResourceWorkspaceAppGroup_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_server_group" "test" {
  name             = "%[1]s"
  os_type          = "Windows"
  flavor_id        = "%[2]s"
  vpc_id           = "%[3]s"
  subnet_id        = "%[4]s"
  system_disk_type = "SAS"
  system_disk_size = 80
  is_vdi           = true
  app_type         = "SESSION_DESKTOP_APP"
  image_id         = "%[5]s"
  image_type       = "gold"
  image_product_id = "%[6]s"
}
`, name, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_AD_VPC_ID,
		acceptance.HW_WORKSPACE_AD_NETWORK_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID)
}

func testResourceWorkspaceAppGroup_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_group" "test" {
  server_group_id = huaweicloud_workspace_app_server_group.test.id
  name            = "%[2]s"
  type            = "SESSION_DESKTOP_APP"
  description     = "Created APP group by script"
}
`, testResourceWorkspaceAppGroup_base(name), name)
}

func testResourceWorkspaceAppGroup_basic_step2(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_group" "test" {
  name = "%[2]s"
}
`, testResourceWorkspaceAppGroup_base(name), updateName)
}
