package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getAppPolicyGroupFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("appstream", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
	}
	return workspace.GetAppGroupPolicy(client, state.Primary.Attributes["name"], state.Primary.ID)
}

// Before running this test, please create a workspace APP server group with SESSION_DESKTOP_APP type.
func TestAccAppPolicyGroup_basic(t *testing.T) {
	var (
		policyGroup  interface{}
		resourceName = "huaweicloud_workspace_app_policy_group.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&policyGroup,
		getAppPolicyGroupFunc,
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
				Config: testAccAppPolicyGroup_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.type", "APPGROUP"),
					resource.TestCheckResourceAttrPair(resourceName, "targets.0.id", "huaweicloud_workspace_app_group.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "targets.0.name", "huaweicloud_workspace_app_group.test", "name"),
				),
			},
			{
				Config: testAccAppPolicyGroup_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.type", "ALL"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.id", "default-apply-all-targets"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.name", "All-Targets"),
				),
			},
			{
				Config: testAccAppPolicyGroup_basic_step3(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "targets.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"priority", "policies",
				},
			},
		},
	})
}

func testAccAppPolicyGroup_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_policy_group" "test" {
  name        = "%[2]s"
  priority    = 1
  description = "Created by terraform script"

  targets {
    id   = huaweicloud_workspace_app_group.test.id
    name = huaweicloud_workspace_app_group.test.name
    type = "APPGROUP"
  }

  policies = jsonencode({
    "client": {
      "automatic_reconnection_interval" : 10,
      "session_persistence_time" : 120,
      "forbid_screen_capture" : true
    }
  })
}
`, testDataSourceAppGroups_base(name), name)
}

func testAccAppPolicyGroup_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_policy_group" "test" {
  name        = "%[1]s"
  priority    = 1
  description = "Updated by terraform script"

  targets {
    id   = "default-apply-all-targets"
    name = "All-Targets"
    type = "ALL"
  }

  policies = jsonencode({
    "client": {
      "automatic_reconnection_interval" : 5,
      "session_persistence_time" : 180,
      "forbid_screen_capture" : false
    }
  })
}
`, name)
}

func testAccAppPolicyGroup_basic_step3(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_policy_group" "test" {
  name     = "%[1]s"
  priority = 1

  policies = jsonencode({
    "client": {
      "automatic_reconnection_interval" : 5,
      "session_persistence_time" : 180,
      "forbid_screen_capture" : false
    }
  })
}
`, name)
}
