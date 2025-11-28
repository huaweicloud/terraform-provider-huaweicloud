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
	return workspace.GetAppGroupPolicyById(client, state.Primary.ID)
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
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
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

// A server group can only be associated with one application group.
func testAccAppPolicyGroup_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_app_server_group" "test" {
  name             = "%[1]s"
  os_type          = "Windows"
  flavor_id        = "%[2]s"
  vpc_id           = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id        = data.huaweicloud_workspace_service.test.network_ids[0]
  system_disk_type = "SAS"
  system_disk_size = 90
  app_type         = "SESSION_DESKTOP_APP"
  is_vdi           = true
  image_id         = "%[3]s"
  image_type       = "gold"
  image_product_id = "%[4]s"
}

resource "huaweicloud_workspace_app_group" "test" {
  server_group_id = huaweicloud_workspace_app_server_group.test.id
  name            = "%[1]s"
  type            = "SESSION_DESKTOP_APP"
  description     = "Created by terraform script"
}
`, name, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID)
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
`, testAccAppPolicyGroup_base(name), name)
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
