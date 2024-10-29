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
					resource.TestCheckResourceAttr(resourceName, "type", "COMMON_APP"),
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

func testResourceWorkspaceAppGroup_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_group" "test" {
  name            = "%[1]s"
  server_group_id = "%[2]s"
  description     = "Created APP group by script"
}
`, name, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID)
}

func testResourceWorkspaceAppGroup_basic_step2(updateName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_group" "test" {
  name = "%[1]s"
}
`, updateName)
}
